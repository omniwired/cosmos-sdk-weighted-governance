package keeper_test

import (
	"testing"

	"cosmos-weighted-governance-sdk/x/voting/types"

	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params:         types.DefaultParams(),
		PortId:         types.PortID,
		VoterRoleList:  []types.VoterRole{{Id: 0}, {Id: 1}},
		VoterRoleCount: 2,
	}
	f := initFixture(t)
	err := f.keeper.InitGenesis(f.ctx, genesisState)
	require.NoError(t, err)
	got, err := f.keeper.ExportGenesis(f.ctx)
	require.NoError(t, err)
	require.NotNil(t, got)

	require.Equal(t, genesisState.PortId, got.PortId)
	require.EqualExportedValues(t, genesisState.Params, got.Params)
	require.EqualExportedValues(t, genesisState.VoterRoleList, got.VoterRoleList)
	require.Equal(t, genesisState.VoterRoleCount, got.VoterRoleCount)

}
