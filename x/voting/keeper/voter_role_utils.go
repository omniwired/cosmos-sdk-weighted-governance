package keeper

import (
	"context"
	"errors"
	"fmt"

	"cosmos-weighted-governance-sdk/x/voting/types"

	"cosmossdk.io/collections"
	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// GetVoterRoleByAddress retrieves a voter role by address
func (k Keeper) GetVoterRoleByAddress(ctx context.Context, address string) (*types.VoterRole, error) {
	var foundRole *types.VoterRole

	err := k.VoterRole.Walk(ctx, nil, func(key uint64, role types.VoterRole) (bool, error) {
		if role.Address == address {
			foundRole = &role
			return true, nil // stop iteration
		}
		return false, nil // continue iteration
	})

	if err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "failed to search voter roles")
	}

	if foundRole == nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrKeyNotFound, "voter role not found for address")
	}

	return foundRole, nil
}

// GetVotingMultiplier returns the voting multiplier for a given address
func (k Keeper) GetVotingMultiplier(ctx context.Context, address string) (math.LegacyDec, error) {
	role, err := k.GetVoterRoleByAddress(ctx, address)
	if err != nil {
		// If no role found, return default multiplier of 1.0
		if errors.Is(err, collections.ErrNotFound) {
			return math.LegacyOneDec(), nil
		}
		return math.LegacyDec{}, err
	}

	// Parse the multiplier string to decimal
	multiplier, err := math.LegacyNewDecFromStr(role.Multiplier)
	if err != nil {
		return math.LegacyDec{}, errorsmod.Wrap(sdkerrors.ErrInvalidRequest,
			fmt.Sprintf("invalid multiplier format: %s", role.Multiplier))
	}

	return multiplier, nil
}

// ValidateVoterRole validates voter role parameters
func (k Keeper) ValidateVoterRole(address, role, multiplier string) error {
	// Validate address format
	if _, err := k.addressCodec.StringToBytes(address); err != nil {
		return errorsmod.Wrap(sdkerrors.ErrInvalidAddress, fmt.Sprintf("invalid address: %s", err))
	}

	// Validate role type
	validRoles := map[string]bool{
		"core_contributor":  true,
		"validator":         true,
		"community_member":  true,
		"strategic_partner": true,
	}

	if !validRoles[role] {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest,
			fmt.Sprintf("invalid role: %s. Valid roles are: core_contributor, validator, community_member, strategic_partner", role))
	}

	// Validate multiplier
	multiplierDec, err := math.LegacyNewDecFromStr(multiplier)
	if err != nil {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest,
			fmt.Sprintf("invalid multiplier format: %s", multiplier))
	}

	// Check multiplier bounds (0.1 to 10.0)
	minMultiplier := math.LegacyNewDecWithPrec(1, 1) // 0.1
	maxMultiplier := math.LegacyNewDec(10)           // 10.0

	if multiplierDec.LT(minMultiplier) || multiplierDec.GT(maxMultiplier) {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest,
			fmt.Sprintf("multiplier must be between 0.1 and 10.0, got: %s", multiplier))
	}

	return nil
}

// GetDefaultMultipliers returns the default multipliers for each role type
func (k Keeper) GetDefaultMultipliers() map[string]string {
	return map[string]string{
		"core_contributor":  "2.0",
		"validator":         "1.5",
		"community_member":  "1.0",
		"strategic_partner": "1.8",
	}
}

// HasVoterRole checks if an address has a voter role
func (k Keeper) HasVoterRole(ctx context.Context, address string) bool {
	_, err := k.GetVoterRoleByAddress(ctx, address)
	return err == nil
}

// ListVoterRolesByRole returns all voter roles of a specific role type
func (k Keeper) ListVoterRolesByRole(ctx context.Context, role string) ([]types.VoterRole, error) {
	var roleList []types.VoterRole

	err := k.VoterRole.Walk(ctx, nil, func(key uint64, voterRole types.VoterRole) (bool, error) {
		if voterRole.Role == role {
			roleList = append(roleList, voterRole)
		}
		return false, nil // continue iteration
	})

	if err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "failed to list voter roles by role")
	}

	return roleList, nil
}

// GetVoterRoleStats returns statistics about voter roles
func (k Keeper) GetVoterRoleStats(ctx context.Context) (map[string]int, error) {
	stats := make(map[string]int)

	err := k.VoterRole.Walk(ctx, nil, func(key uint64, voterRole types.VoterRole) (bool, error) {
		stats[voterRole.Role]++
		return false, nil // continue iteration
	})

	if err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "failed to get voter role stats")
	}

	return stats, nil
}

// CountRolesForAddress counts the number of roles assigned to a specific address
func (k Keeper) CountRolesForAddress(ctx context.Context, address string) uint32 {
	var count uint32
	
	k.VoterRole.Walk(ctx, nil, func(id uint64, vr types.VoterRole) (stop bool, err error) {
		if vr.Address == address {
			count++
		}
		return false, nil
	})
	
	return count
}
