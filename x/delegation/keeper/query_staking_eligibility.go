package keeper

import (
	"context"

	"cosmos-weighted-governance-sdk/x/delegation/types"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (q queryServer) StakingEligibility(ctx context.Context, req *types.QueryStakingEligibilityRequest) (*types.QueryStakingEligibilityResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	if req.Address == "" {
		return nil, status.Error(codes.InvalidArgument, "address cannot be empty")
	}

	// Check staking eligibility using the keeper method
	eligibility, err := q.k.CheckStakingEligibility(ctx, req.Address)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return eligibility, nil
}
