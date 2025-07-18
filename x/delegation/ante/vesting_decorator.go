package ante

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"cosmos-weighted-governance-sdk/x/delegation/keeper"
)

// VestingDelegationDecorator checks if vesting accounts are attempting to delegate unvested tokens
type VestingDelegationDecorator struct {
	dk keeper.Keeper
}

// NewVestingDelegationDecorator creates a new VestingDelegationDecorator
func NewVestingDelegationDecorator(dk keeper.Keeper) VestingDelegationDecorator {
	return VestingDelegationDecorator{
		dk: dk,
	}
}

// AnteHandle checks if the transaction contains staking messages and validates vesting constraints
func (vdd VestingDelegationDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {
	// Skip validation during simulation
	if simulate {
		return next(ctx, tx, simulate)
	}

	// Check all messages in the transaction
	for _, msg := range tx.GetMsgs() {
		// Check if this is a delegation message
		switch msg := msg.(type) {
		case *stakingtypes.MsgDelegate:
			if err := vdd.validateDelegation(ctx, msg.DelegatorAddress, msg.Amount); err != nil {
				return ctx, err
			}
		case *stakingtypes.MsgBeginRedelegate:
			// For redelegation, we need to check if the source delegation can be moved
			if err := vdd.validateDelegation(ctx, msg.DelegatorAddress, msg.Amount); err != nil {
				return ctx, err
			}
		case *stakingtypes.MsgUndelegate:
			// Undelegation is typically allowed, but we can add checks if needed
			continue
		}
	}

	return next(ctx, tx, simulate)
}

// validateDelegation checks if the delegation is allowed based on vesting status
func (vdd VestingDelegationDecorator) validateDelegation(ctx sdk.Context, delegatorAddr string, amount sdk.Coin) error {
	// Use the keeper's validation method
	err := vdd.dk.ValidateStakingTransaction(ctx, delegatorAddr, amount)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "vesting validation failed: %s", err.Error())
	}

	return nil
}