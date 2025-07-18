package keeper

import (
	"context"
	"errors"

	"cosmos-weighted-governance-sdk/x/voting/types"

	"cosmossdk.io/collections"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (q queryServer) ListVoterRole(ctx context.Context, req *types.QueryAllVoterRoleRequest) (*types.QueryAllVoterRoleResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	voterRoles, pageRes, err := query.CollectionPaginate(
		ctx,
		q.k.VoterRole,
		req.Pagination,
		func(_ uint64, value types.VoterRole) (types.VoterRole, error) {
			return value, nil
		},
	)

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllVoterRoleResponse{VoterRole: voterRoles, Pagination: pageRes}, nil
}

func (q queryServer) GetVoterRole(ctx context.Context, req *types.QueryGetVoterRoleRequest) (*types.QueryGetVoterRoleResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	voterRole, err := q.k.VoterRole.Get(ctx, req.Id)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return nil, sdkerrors.ErrKeyNotFound
		}

		return nil, status.Error(codes.Internal, "internal error")
	}

	return &types.QueryGetVoterRoleResponse{VoterRole: voterRole}, nil
}

// QueryVoterRoleByAddress queries voter role by address
func (q queryServer) QueryVoterRoleByAddress(ctx context.Context, address string) (*types.VoterRole, error) {
	return q.k.GetVoterRoleByAddress(ctx, address)
}

// QueryVotingMultiplier queries voting multiplier for an address
func (q queryServer) QueryVotingMultiplier(ctx context.Context, address string) (string, error) {
	multiplier, err := q.k.GetVotingMultiplier(ctx, address)
	if err != nil {
		return "", err
	}
	return multiplier.String(), nil
}

// QueryVoterRolesByRole queries all voter roles of a specific role type
func (q queryServer) QueryVoterRolesByRole(ctx context.Context, role string) ([]types.VoterRole, error) {
	return q.k.ListVoterRolesByRole(ctx, role)
}

// QueryVoterRoleStats queries statistics about voter roles
func (q queryServer) QueryVoterRoleStats(ctx context.Context) (map[string]int, error) {
	return q.k.GetVoterRoleStats(ctx)
}
