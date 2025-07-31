package keeper_test

import (
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"

	"cosmos-weighted-governance-sdk/x/voting/keeper"
	"cosmos-weighted-governance-sdk/x/voting/types"
)

func TestVoterRoleMsgServerCreate(t *testing.T) {
	f := initFixture(t)
	srv := keeper.NewMsgServerImpl(f.keeper)

	// gov authority
	creator, err := f.addressCodec.BytesToString(f.keeper.GetAuthority())
	require.NoError(t, err)

	// happy path test
	resp, err := srv.CreateVoterRole(f.ctx, &types.MsgCreateVoterRole{
		Creator:    creator,
		Address:    "cosmos1wd5kwmn9wfqkgerjta047h6lta047h6lta047h6lta047wfl63q",
		Role:       "validator",
		Multiplier: "1.5",
		AddedAt:    1234567890,
		AddedBy:    creator,
	})
	require.NoError(t, err)
	require.Equal(t, uint64(1), resp.Id)
}

func TestVoterRoleMsgServerUpdate(t *testing.T) {
	f := initFixture(t)
	srv := keeper.NewMsgServerImpl(f.keeper)

	creator, err := f.addressCodec.BytesToString(f.keeper.GetAuthority())
	require.NoError(t, err)

	unauthorizedAddr, err := f.addressCodec.BytesToString([]byte("unauthorizedAddr___________"))
	require.NoError(t, err)

	// create one to update later
	_, err = srv.CreateVoterRole(f.ctx, &types.MsgCreateVoterRole{
		Creator:    creator,
		Address:    "cosmos1wd5kwmn9wfqkgerjta047h6lta047h6lta047h6lta047wfl63q",
		Role:       "validator", 
		Multiplier: "1.5",
		AddedAt:    1234567890,
		AddedBy:    creator,
	})
	require.NoError(t, err)

	tests := []struct {
		desc    string
		request *types.MsgUpdateVoterRole
		err     error
	}{
		{
			desc:    "invalid address",
			request: &types.MsgUpdateVoterRole{Creator: "invalid"},
			err:     sdkerrors.ErrInvalidAddress,
		},
		{
			desc:    "unauthorized",
			request: &types.MsgUpdateVoterRole{Creator: unauthorizedAddr},
			err:     types.ErrInvalidSigner,
		},
		{
			desc:    "key not found",
			request: &types.MsgUpdateVoterRole{Creator: creator, Id: 10},
			err:     sdkerrors.ErrKeyNotFound,
		},
		{
			desc: "completed",
			request: &types.MsgUpdateVoterRole{
				Creator:    creator,
				Id:         1,
				Address:    "cosmos1wd5kwmn9wfqkgerjta047h6lta047h6lta047h6lta047wfl63q",
				Role:       "core_contributor",
				Multiplier: "2.0",
				AddedAt:    1234567890,
				AddedBy:    creator,
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			_, err = srv.UpdateVoterRole(f.ctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestVoterRoleMsgServerDelete(t *testing.T) {
	f := initFixture(t)
	srv := keeper.NewMsgServerImpl(f.keeper)

	// Use governance authority
	creator, err := f.addressCodec.BytesToString(f.keeper.GetAuthority())
	require.NoError(t, err)

	unauthorizedAddr, err := f.addressCodec.BytesToString([]byte("unauthorizedAddr___________"))
	require.NoError(t, err)

	_, err = srv.CreateVoterRole(f.ctx, &types.MsgCreateVoterRole{
		Creator:    creator,
		Address:    "cosmos1wd5kwmn9wfqkgerjta047h6lta047h6lta047h6lta047wfl63q",
		Role:       "validator",
		Multiplier: "1.5",
		AddedAt:    1234567890,
		AddedBy:    creator,
	})
	require.NoError(t, err)

	tests := []struct {
		desc    string
		request *types.MsgDeleteVoterRole
		err     error
	}{
		{
			desc:    "invalid address",
			request: &types.MsgDeleteVoterRole{Creator: "invalid"},
			err:     sdkerrors.ErrInvalidAddress,
		},
		{
			desc:    "unauthorized",
			request: &types.MsgDeleteVoterRole{Creator: unauthorizedAddr},
			err:     types.ErrInvalidSigner,
		},
		{
			desc:    "key not found",
			request: &types.MsgDeleteVoterRole{Creator: creator, Id: 10},
			err:     sdkerrors.ErrKeyNotFound,
		},
		{
			desc:    "completed",
			request: &types.MsgDeleteVoterRole{Creator: creator, Id: 1},
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			_, err = srv.DeleteVoterRole(f.ctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
