package types_test

import (
	"testing"

	"cosmos-weighted-governance-sdk/x/voting/types"

	"github.com/stretchr/testify/require"
)

func TestGenesisState_Validate(t *testing.T) {
	tests := []struct {
		desc     string
		genState *types.GenesisState
		valid    bool
	}{
		{
			desc:     "default is valid",
			genState: types.DefaultGenesis(),
			valid:    true,
		},
		{
			desc: "valid genesis state",
			genState: &types.GenesisState{
				PortId: types.PortID,
				Params: types.Params{
					MaxVoterRolesPerAddress: 1,
					RoleCreationCooldown:    300,
				},
				VoterRoleList: []types.VoterRole{{Id: 0}, {Id: 1}}, VoterRoleCount: 2,
			}, valid: true,
		}, {
			desc: "duplicated voterRole",
			genState: &types.GenesisState{
				VoterRoleList: []types.VoterRole{
					{
						Id: 0,
					},
					{
						Id: 0,
					},
				},
			},
			valid: false,
		}, {
			desc: "invalid voterRole count",
			genState: &types.GenesisState{
				VoterRoleList: []types.VoterRole{
					{
						Id: 1,
					},
				},
				VoterRoleCount: 0,
			},
			valid: false,
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			err := tc.genState.Validate()
			if tc.valid {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}
