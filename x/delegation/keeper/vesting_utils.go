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
	// Convert address string to AccAddress
	accAddr, err := k.authKeeper.AddressCodec().StringToBytes(address)
	if err != nil {
		return &types.QueryStakingEligibilityResponse{
			IsEligible: false,
			Reason:     "invalid address format",
			IsVesting:  false,
		}, nil
	}

	// Get account from auth keeper
	account := k.authKeeper.GetAccount(ctx, accAddr)
	if account == nil {
		return &types.QueryStakingEligibilityResponse{
			IsEligible: false,
			Reason:     "account not found",
			IsVesting:  false,
		}, nil
	}

	// Check if account is a vesting account
	vestingAcc, isVesting := account.(types.VestingAccount)
	if !isVesting {
		// Non-vesting accounts are always eligible
		return &types.QueryStakingEligibilityResponse{
			IsEligible: true,
			Reason:     "non-vesting account",
			IsVesting:  false,
		}, nil
	}

	// For vesting accounts, check vesting status
	blockTime := sdk.UnwrapSDKContext(ctx).BlockTime()
	
	// Get vesting and vested amounts
	vestedCoins := vestingAcc.GetVestedCoins(blockTime)
	vestingCoins := vestingAcc.GetVestingCoins(blockTime)
	originalVesting := vestingAcc.GetOriginalVesting()

	// Calculate amounts (assuming we're dealing with the native token)
	var vestedAmount, vestingAmount int64
	
	// Get the native token denomination (usually "stake" or similar)
	nativeDenom := "stake" // This should be configurable or determined from app state
	
	if vestedCoins.AmountOf(nativeDenom).IsPositive() {
		vestedAmount = vestedCoins.AmountOf(nativeDenom).Int64()
	}
	
	if vestingCoins.AmountOf(nativeDenom).IsPositive() {
		vestingAmount = vestingCoins.AmountOf(nativeDenom).Int64()
	}

	// Check if all tokens are vested
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

	// If tokens are still vesting, restrict staking
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
	eligibility, err := k.CheckStakingEligibility(ctx, delegatorAddr)
	if err != nil {
		return err
	}

	if !eligibility.IsEligible {
		return fmt.Errorf("staking not allowed: %s", eligibility.Reason)
	}

	// Additional validation for vesting accounts
	if eligibility.IsVesting {
		// Could add more sophisticated logic here, such as:
		// - Only allow staking of vested tokens
		// - Restrict staking amounts based on vesting schedule
		// - etc.
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