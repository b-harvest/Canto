package keeper

import (
	"context"
	"github.com/Canto-Network/Canto/v6/x/liquidstaking/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ types.MsgServer = &Keeper{}

func (k Keeper) LiquidStake(goCtx context.Context, msg *types.MsgLiquidStake) (*types.MsgLiquidStakeResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	if !k.bankKeeper.HasBalance(ctx, sdk.AccAddress(msg.DelegatorAddress), msg.Amount) {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "")
	}

	panic("implement me")
}

func (k Keeper) LiquidUnstake(goCtx context.Context, msg *types.MsgLiquidUnstake) (*types.MsgLiquidUnstakeResponse, error) {
	//ctx := sdk.UnwrapSDKContext(goCtx)
	panic("implement me")
}

func (k Keeper) InsuranceProvide(goCtx context.Context, msg *types.MsgInsuranceProvide) (*types.MsgInsuranceProvideResponse, error) {
	//ctx := sdk.UnwrapSDKContext(goCtx)
	panic("implement me")
}

func (k Keeper) CancelInsuranceProvide(goCtx context.Context, msg *types.MsgCancelInsuranceProvide) (*types.MsgCancelInsuranceProvideResponse, error) {
	//ctx := sdk.UnwrapSDKContext(goCtx)
	panic("implement me")
}

func (k Keeper) DepositInsurance(goCtx context.Context, msg *types.MsgDepositInsurance) (*types.MsgDepositInsuranceResponse, error) {
	//ctx := sdk.UnwrapSDKContext(goCtx)
	panic("implement me")
}

func (k Keeper) WithdrawInsurance(goCtx context.Context, msg *types.MsgWithdrawInsurance) (*types.MsgWithdrawInsuranceResponse, error) {
	//ctx := sdk.UnwrapSDKContext(goCtx)
	panic("implement me")
}

func (k Keeper) WithdrawInsuranceCommission(goCtx context.Context, msg *types.MsgWithdrawInsuranceCommission) (*types.MsgWithdrawInsuranceCommissionResponse, error) {
	//ctx := sdk.UnwrapSDKContext(goCtx)
	panic("implement me")
}
