package keeper

import (
	"context"

	"github.com/Canto-Network/Canto-Testnet-v2/v1/x/liquidstaking/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ types.QueryServer = Querier{}

type Querier struct {
	Keeper
}

func NewQuerier(k Keeper) Querier {
	return Querier{Keeper: k}
}

func (k Querier) Params(
	c context.Context,
	req *types.QueryParamsRequest,
) (*types.QueryParamsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}
	ctx := sdk.UnwrapSDKContext(c)
	params := k.GetParams(ctx)
	return &types.QueryParamsResponse{Params: params}, nil
}

func (k Querier) LiquidValidators(
	c context.Context,
	req *types.QueryLiquidValidatorsRequest,
) (*types.QueryLiquidValidatorsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}
	_ = sdk.UnwrapSDKContext(c)
	liquidValidators := types.LiquidValidatorStates{}
	return &types.QueryLiquidValidatorsResponse{LiquidValidators: liquidValidators}, nil
}

func (k Querier) AliveChunks(
	c context.Context,
	req *types.QueryAliveChunksRequest,
) (*types.QueryAliveChunksResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}
	_ = sdk.UnwrapSDKContext(c)
	return &types.QueryAliveChunksResponse{}, nil
}

func (k Querier) UnbondingChunks(
	c context.Context,
	req *types.QueryUnbondingChunksRequest,
) (*types.QueryUnbondingChunksResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}
	_ = sdk.UnwrapSDKContext(c)
	return &types.QueryUnbondingChunksResponse{}, nil
}

func (k Querier) ChunkBondRequests(
	c context.Context,
	req *types.QueryChunkBondRequests,
) (*types.QueryChunkBondRequestsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}
	_ = sdk.UnwrapSDKContext(c)
	return &types.QueryChunkBondRequestsResponse{}, nil
}

func (k Querier) ChunkUnbondRequests(
	c context.Context,
	req *types.QueryChunkUnbondRequests,
) (*types.QueryChunkUnbondRequestsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}
	_ = sdk.UnwrapSDKContext(c)
	return &types.QueryChunkUnbondRequestsResponse{}, nil
}

func (k Querier) InsuranceBids(
	c context.Context,
	req *types.QueryInsuranceBidRequest,
) (*types.QueryInsuranceBidResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}
	_ = sdk.UnwrapSDKContext(c)
	return &types.QueryInsuranceBidResponse{}, nil
}

func (k Querier) LiquidStakingState(
	c context.Context,
	req *types.QueryLiquidStakingStateRequest,
) (*types.QueryLiquidStakingStateResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}
	_ = sdk.UnwrapSDKContext(c)
	return &types.QueryLiquidStakingStateResponse{}, nil
}
