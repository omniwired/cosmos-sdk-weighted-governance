package keeper

import (
	"context"
	"fmt"
	"time"

	"cosmos-weighted-governance-sdk/x/delegation/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// CheckStakingEligibility checks if an account is eligible to stake
func (k Keeper) CheckStakingEligibility(ctx context.Context, address string) (*types.QueryStakingEligibilityResponse, error) {
	accAddr, err := k.authKeeper.AddressCodec().StringToBytes(address)
	if err != nil {
		return &types.QueryStakingEligibilityResponse{
			IsEligible: false,
			Reason:     "invalid address format",
			IsVesting:  false,
		}, nil
	}

	account := k.authKeeper.GetAccount(ctx, accAddr)
	if account == nil {
		return &types.QueryStakingEligibilityResponse{
			IsEligible: false,
			Reason:     "account not found",
			IsVesting:  false,
		}, nil
	}

	vestingAcc, isVesting := account.(types.VestingAccount)
	if !isVesting {
		// regular accounts can stake whatever
		return &types.QueryStakingEligibilityResponse{
			IsEligible: true,
			Reason:     "non-vesting account",
			IsVesting:  false,
		}, nil
	}

	// vesting accounts need special handling
	blockTime := sdk.UnwrapSDKContext(ctx).BlockTime()
	
	vestedCoins := vestingAcc.GetVestedCoins(blockTime)
	vestingCoins := vestingAcc.GetVestingCoins(blockTime)
	originalVesting := vestingAcc.GetOriginalVesting()

	var vestedAmount, vestingAmount int64
	
	params, err := k.Params.Get(ctx)
	if err != nil {
		return &types.QueryStakingEligibilityResponse{
			IsEligible: false,
			Reason:     "failed to get module params",
			IsVesting:  true,
		}, nil
	}
	
	stakeDenom := params.StakeDenom
	
	if vestedCoins.AmountOf(stakeDenom).IsPositive() {
		vestedAmount = vestedCoins.AmountOf(stakeDenom).Int64()
	}
	
	if vestingCoins.AmountOf(stakeDenom).IsPositive() {
		vestingAmount = vestingCoins.AmountOf(stakeDenom).Int64()
	}

	// check if fully vested
	allVested := vestingCoins.IsZero() || 
		vestedCoins.IsAllGTE(originalVesting)

	if allVested {
		return &types.QueryStakingEligibilityResponse{
			IsEligible:    true,
			Reason:        "all tokens are vested",
			IsVesting:     true,
			VestedAmount:  vestedAmount,
			VestingAmount: vestingAmount,
		}, nil
	}

	// nope, still vesting
	return &types.QueryStakingEligibilityResponse{
		IsEligible:    false,
		Reason:        "tokens are still vesting - staking restricted",
		IsVesting:     true,
		VestedAmount:  vestedAmount,
		VestingAmount: vestingAmount,
	}, nil
}

// IsVestingAccount checks if an account is a vesting account
func (k Keeper) IsVestingAccount(ctx context.Context, address string) bool {
	accAddr, err := k.authKeeper.AddressCodec().StringToBytes(address)
	if err != nil {
		return false
	}

	account := k.authKeeper.GetAccount(ctx, accAddr)
	if account == nil {
		return false
	}

	_, isVesting := account.(types.VestingAccount)
	return isVesting
}

// ValidateStakingTransaction validates a staking transaction for vesting restrictions
func (k Keeper) ValidateStakingTransaction(ctx context.Context, delegatorAddr string, amount sdk.Coin) error {
	accAddr, err := k.authKeeper.AddressCodec().StringToBytes(delegatorAddr)
	if err != nil {
		return fmt.Errorf("invalid delegator address: %s", err)
	}

	account := k.authKeeper.GetAccount(ctx, accAddr)
	if account == nil {
		return fmt.Errorf("account not found")
	}

	vestingAcc, isVesting := account.(types.VestingAccount)
	if !isVesting {
		return nil // regular accounts have no restrictions
	}

	params, err := k.Params.Get(ctx)
	if err != nil {
		return fmt.Errorf("failed to get module params: %s", err)
	}

	if amount.Denom != params.StakeDenom {
		return nil // not staking native token, who cares
	}

	blockTime := sdk.UnwrapSDKContext(ctx).BlockTime()
	vestedCoins := vestingAcc.GetVestedCoins(blockTime)
	
	vestedAmount := vestedCoins.AmountOf(params.StakeDenom)
	
	// trying to stake more than they've vested? nice try
	if amount.Amount.GT(vestedAmount) {
		return fmt.Errorf("cannot stake unvested tokens: requested %s, vested %s %s", 
			amount.Amount.String(), vestedAmount.String(), params.StakeDenom)
	}

	return nil
}

// GetVestingInfo returns detailed vesting information for an account
func (k Keeper) GetVestingInfo(ctx context.Context, address string) (*VestingInfo, error) {
	accAddr, err := k.authKeeper.AddressCodec().StringToBytes(address)
	if err != nil {
		return nil, err
	}

	account := k.authKeeper.GetAccount(ctx, accAddr)
	if account == nil {
		return nil, fmt.Errorf("account not found")
	}

	vestingAcc, isVesting := account.(types.VestingAccount)
	if !isVesting {
		return &VestingInfo{
			IsVesting:     false,
			IsFullyVested: true,
		}, nil
	}

	blockTime := sdk.UnwrapSDKContext(ctx).BlockTime()
	vestedCoins := vestingAcc.GetVestedCoins(blockTime)
	vestingCoins := vestingAcc.GetVestingCoins(blockTime)
	originalVesting := vestingAcc.GetOriginalVesting()

	return &VestingInfo{
		IsVesting:       true,
		IsFullyVested:   vestingCoins.IsZero(),
		VestedCoins:     vestedCoins,
		VestingCoins:    vestingCoins,
		OriginalVesting: originalVesting,
		BlockTime:       blockTime,
	}, nil
}

// VestingInfo holds vesting information for an account
type VestingInfo struct {
	IsVesting       bool
	IsFullyVested   bool
	VestedCoins     sdk.Coins
	VestingCoins    sdk.Coins
	OriginalVesting sdk.Coins
	BlockTime       time.Time
}