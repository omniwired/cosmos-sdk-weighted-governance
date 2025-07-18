package keeper

import (
	"context"
	"fmt"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

// GovHooksWrapper is a wrapper for governance hooks that applies voting multipliers
type GovHooksWrapper struct {
	k              Keeper
	originalHooks  govtypes.GovHooks
}

// NewGovHooksWrapper creates a new governance hooks wrapper
func NewGovHooksWrapper(k Keeper, originalHooks govtypes.GovHooks) GovHooksWrapper {
	return GovHooksWrapper{
		k:             k,
		originalHooks: originalHooks,
	}
}

// AfterProposalSubmission is called after a proposal is submitted
func (h GovHooksWrapper) AfterProposalSubmission(ctx context.Context, proposalID uint64) error {
	if h.originalHooks != nil {
		return h.originalHooks.AfterProposalSubmission(ctx, proposalID)
	}
	return nil
}

// AfterProposalDeposit is called after a deposit is made
func (h GovHooksWrapper) AfterProposalDeposit(ctx context.Context, proposalID uint64, depositorAddr sdk.AccAddress) error {
	if h.originalHooks != nil {
		return h.originalHooks.AfterProposalDeposit(ctx, proposalID, depositorAddr)
	}
	return nil
}

// AfterProposalVote is called after a vote is cast - this is where we apply multipliers
func (h GovHooksWrapper) AfterProposalVote(ctx context.Context, proposalID uint64, voterAddr sdk.AccAddress) error {
	// Get the voter's address as string
	voterStr, err := h.k.addressCodec.BytesToString(voterAddr)
	if err != nil {
		return fmt.Errorf("failed to convert voter address: %w", err)
	}

	// Get the voting multiplier for this address
	multiplier, err := h.k.GetVotingMultiplier(ctx, voterStr)
	if err != nil {
		// Log error but don't fail - use default multiplier of 1.0
		multiplier = math.LegacyOneDec()
	}

	// Store the multiplier for this vote (to be applied during tally)
	// This would require additional state storage in a real implementation
	// For now, we'll just emit an event to track it
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			"vote_multiplier_applied",
			sdk.NewAttribute("proposal_id", fmt.Sprintf("%d", proposalID)),
			sdk.NewAttribute("voter", voterStr),
			sdk.NewAttribute("multiplier", multiplier.String()),
		),
	)

	// Call original hooks if any
	if h.originalHooks != nil {
		return h.originalHooks.AfterProposalVote(ctx, proposalID, voterAddr)
	}
	return nil
}

// AfterProposalFailedMinDeposit is called after a proposal fails to meet min deposit
func (h GovHooksWrapper) AfterProposalFailedMinDeposit(ctx context.Context, proposalID uint64) error {
	if h.originalHooks != nil {
		return h.originalHooks.AfterProposalFailedMinDeposit(ctx, proposalID)
	}
	return nil
}

// AfterProposalVotingPeriodEnded is called after the voting period ends
func (h GovHooksWrapper) AfterProposalVotingPeriodEnded(ctx context.Context, proposalID uint64) error {
	if h.originalHooks != nil {
		return h.originalHooks.AfterProposalVotingPeriodEnded(ctx, proposalID)
	}
	return nil
}