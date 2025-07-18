package keeper

import (
	"context"
	"errors"

	"cosmos-weighted-governance-sdk/x/voting/types"

	"cosmossdk.io/collections"
)

// InitGenesis initializes the module's state from a provided genesis state.
func (k Keeper) InitGenesis(ctx context.Context, genState types.GenesisState) error {
	if err := k.Port.Set(ctx, genState.PortId); err != nil {
		return err
	}
	for _, elem := range genState.VoterRoleList {
		if err := k.VoterRole.Set(ctx, elem.Id, elem); err != nil {
			return err
		}
	}

	if err := k.VoterRoleSeq.Set(ctx, genState.VoterRoleCount); err != nil {
		return err
	}

	return k.Params.Set(ctx, genState.Params)
}

// ExportGenesis returns the module's exported genesis.
func (k Keeper) ExportGenesis(ctx context.Context) (*types.GenesisState, error) {
	var err error

	genesis := types.DefaultGenesis()
	genesis.Params, err = k.Params.Get(ctx)
	if err != nil {
		return nil, err
	}
	genesis.PortId, err = k.Port.Get(ctx)
	if err != nil && !errors.Is(err, collections.ErrNotFound) {
		return nil, err
	}
	err = k.VoterRole.Walk(ctx, nil, func(key uint64, elem types.VoterRole) (bool, error) {
		genesis.VoterRoleList = append(genesis.VoterRoleList, elem)
		return false, nil
	})
	if err != nil {
		return nil, err
	}

	genesis.VoterRoleCount, err = k.VoterRoleSeq.Peek(ctx)
	if err != nil {
		return nil, err
	}

	return genesis, nil
}
