package keeper

import (
	"context"
	"fmt"

	"cosmos-weighted-governance-sdk/x/voting/types"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	govtypesv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// WeightedTallyHandler is a custom tally handler that applies voting multipliers
type WeightedTallyHandler struct {
	k              Keeper
	stakingKeeper  govtypes.StakingKeeper
}

// NewWeightedTallyHandler creates a new weighted tally handler
func NewWeightedTallyHandler(k Keeper, sk govtypes.StakingKeeper) WeightedTallyHandler {
	return WeightedTallyHandler{
		k:             k,
		stakingKeeper: sk,
	}
}

// Tally calculates the weighted voting results
func (h WeightedTallyHandler) Tally(ctx context.Context, proposal govtypesv1.Proposal) (govtypesv1.TallyResult, error) {
	results := govtypesv1.NewTallyResult()
	
	// Get all votes for the proposal
	// Note: In a real implementation, you would need access to the gov keeper's vote storage
	// This is a simplified version showing the concept
	
	// For demonstration, we'll walk through all voter roles and apply multipliers
	// In production, you'd integrate with the actual governance vote storage
	
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	
	// Get total bonded tokens from staking keeper for percentage calculations
	totalBonded, err := h.stakingKeeper.TotalBondedTokens(ctx)
	if err != nil {
		return results, fmt.Errorf("failed to get total bonded tokens: %w", err)
	}

	// Walk through all voter roles to apply multipliers
	// Note: This is a simplified approach - in production you'd process actual votes
	err = h.k.VoterRole.Walk(ctx, nil, func(id uint64, role types.VoterRole) (bool, error) {
		// Get the voter's address
		voterAddr, err := h.k.addressCodec.StringToBytes(role.Address)
		if err != nil {
			return false, nil // Skip invalid addresses
		}

		// Get delegation info to determine voting power
		delegations, err := h.stakingKeeper.GetDelegatorDelegations(ctx, voterAddr, 100)
		if err != nil {
			return false, nil // Skip if can't get delegations
		}

		// Calculate base voting power from delegations
		votingPower := math.LegacyZeroDec()
		for _, delegation := range delegations {
			validator, err := h.stakingKeeper.GetValidator(ctx, delegation.GetValidatorAddr())
			if err != nil {
				continue
			}
			tokens := validator.TokensFromShares(delegation.GetShares())
			votingPower = votingPower.Add(tokens)
		}

		// Apply multiplier
		multiplier, err := math.LegacyNewDecFromStr(role.Multiplier)
		if err != nil {
			multiplier = math.LegacyOneDec()
		}
		
		weightedPower := votingPower.Mul(multiplier)

		// For demonstration, we'll assume all votes are "Yes"
		// In production, you'd check the actual vote option
		results.YesCount = results.YesCount.Add(weightedPower.TruncateInt())

		return false, nil // Continue iteration
	})

	if err != nil {
		return results, fmt.Errorf("failed to tally votes: %w", err)
	}

	return results, nil
}