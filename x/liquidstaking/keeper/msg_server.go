package keeper

import (
	"context"

	"github.com/Canto-Network/Canto-Testnet-v2/v1/x/liquidstaking/types"
)

var _ types.MsgServer = msgServer{}

type msgServer struct {
	Keeper
}

func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

func (k msgServer) LiquidStaking(
	goCtx context.Context,
	msg *types.MsgLiquidStaking,
) (*types.MsgLiquidStakingResponse, error) {
	return &types.MsgLiquidStakingResponse{}, nil
}

func (k msgServer) CancelLiquidStaking(
	goCtx context.Context,
	msg *types.MsgCancelLiquidStaking,
) (*types.MsgCancelLiquidStakingResponse, error) {
	return &types.MsgCancelLiquidStakingResponse{}, nil
}

func (k msgServer) LiquidUnstaking(
	goCtx context.Context,
	msg *types.MsgLiquidUnstaking,
) (*types.MsgLiquidUnstakingResponse, error) {
	return &types.MsgLiquidUnstakingResponse{}, nil
}

func (k msgServer) CancelLiquidUnstaking(
	goCtx context.Context,
	msg *types.MsgCancelLiquidUnstaking,
) (*types.MsgCancelLiquidUnstakingResponse, error) {
	return &types.MsgCancelLiquidUnstakingResponse{}, nil
}

func (k msgServer) BidInsurance(
	goCtx context.Context,
	msg *types.MsgBidInsurance,
) (*types.MsgBidInsuranceResponse, error) {
	return &types.MsgBidInsuranceResponse{}, nil
}

func (k msgServer) CancelInsuranceBid(
	goCtx context.Context,
	msg *types.MsgCancelInsuranceBid,
) (*types.MsgCancelInsuranceBidResponse, error) {
	return &types.MsgCancelInsuranceBidResponse{}, nil
}

func (k msgServer) UnbondInsurance(
	goCtx context.Context,
	msg *types.MsgUnbondInsurance,
) (*types.MsgUnbondInsuranceResponse, error) {
	return &types.MsgUnbondInsuranceResponse{}, nil
}

func (k msgServer) CancelInsuranceUnbond(
	goCtx context.Context,
	msg *types.MsgCancelInsuranceUnbond,
) (*types.MsgCancelInsuranceUnbondResponse, error) {
	return &types.MsgCancelInsuranceUnbondResponse{}, nil
}
