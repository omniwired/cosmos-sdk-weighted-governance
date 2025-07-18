package keeper_test

import (
	"context"
	"strconv"
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"cosmos-weighted-governance-sdk/x/voting/keeper"
	"cosmos-weighted-governance-sdk/x/voting/types"
)

func createNVoterRole(keeper keeper.Keeper, ctx context.Context, n int) []types.VoterRole {
	items := make([]types.VoterRole, n)
	for i := range items {
		iu := uint64(i)
		items[i].Id = iu
		items[i].Address = strconv.Itoa(i)
		items[i].Role = strconv.Itoa(i)
		items[i].Multiplier = strconv.Itoa(i)
		items[i].AddedAt = int64(i)
		items[i].AddedBy = strconv.Itoa(i)
		_ = keeper.VoterRole.Set(ctx, iu, items[i])
		_ = keeper.VoterRoleSeq.Set(ctx, iu)
	}
	return items
}

func TestVoterRoleQuerySingle(t *testing.T) {
	f := initFixture(t)
	qs := keeper.NewQueryServerImpl(f.keeper)
	msgs := createNVoterRole(f.keeper, f.ctx, 2)
	tests := []struct {
		desc     string
		request  *types.QueryGetVoterRoleRequest
		response *types.QueryGetVoterRoleResponse
		err      error
	}{
		{
			desc:     "First",
			request:  &types.QueryGetVoterRoleRequest{Id: msgs[0].Id},
			response: &types.QueryGetVoterRoleResponse{VoterRole: msgs[0]},
		},
		{
			desc:     "Second",
			request:  &types.QueryGetVoterRoleRequest{Id: msgs[1].Id},
			response: &types.QueryGetVoterRoleResponse{VoterRole: msgs[1]},
		},
		{
			desc:    "KeyNotFound",
			request: &types.QueryGetVoterRoleRequest{Id: uint64(len(msgs))},
			err:     sdkerrors.ErrKeyNotFound,
		},
		{
			desc: "InvalidRequest",
			err:  status.Error(codes.InvalidArgument, "invalid request"),
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			response, err := qs.GetVoterRole(f.ctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				require.EqualExportedValues(t, tc.response, response)
			}
		})
	}
}

func TestVoterRoleQueryPaginated(t *testing.T) {
	f := initFixture(t)
	qs := keeper.NewQueryServerImpl(f.keeper)
	msgs := createNVoterRole(f.keeper, f.ctx, 5)

	request := func(next []byte, offset, limit uint64, total bool) *types.QueryAllVoterRoleRequest {
		return &types.QueryAllVoterRoleRequest{
			Pagination: &query.PageRequest{
				Key:        next,
				Offset:     offset,
				Limit:      limit,
				CountTotal: total,
			},
		}
	}
	t.Run("ByOffset", func(t *testing.T) {
		step := 2
		for i := 0; i < len(msgs); i += step {
			resp, err := qs.ListVoterRole(f.ctx, request(nil, uint64(i), uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.VoterRole), step)
			require.Subset(t, msgs, resp.VoterRole)
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(msgs); i += step {
			resp, err := qs.ListVoterRole(f.ctx, request(next, 0, uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.VoterRole), step)
			require.Subset(t, msgs, resp.VoterRole)
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		resp, err := qs.ListVoterRole(f.ctx, request(nil, 0, 0, true))
		require.NoError(t, err)
		require.Equal(t, len(msgs), int(resp.Pagination.Total))
		require.EqualExportedValues(t, msgs, resp.VoterRole)
	})
	t.Run("InvalidRequest", func(t *testing.T) {
		_, err := qs.ListVoterRole(f.ctx, nil)
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, "invalid request"))
	})
}
