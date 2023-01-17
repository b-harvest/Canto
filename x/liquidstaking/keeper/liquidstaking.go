package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/Canto-Network/Canto-Testnet-v2/v1/x/liquidstaking/types"
)

func (k Keeper) LiquidBondDenom(ctx sdk.Context) (res string) {
	k.paramstore.Get(ctx, types.ParamStoreKeyLiquidBondDenom, &res)
	return
}

func (k Keeper) NewLiquidStakingState() types.LiquidStakingState {
	// TODO: calc mint rate
	return types.LiquidStakingState{
		MintRate: sdk.OneDec(),
	}
}

func (k Keeper) LiquidStake(
	ctx sdk.Context, liquidStaker sdk.AccAddress, tokenAmount sdk.Int) (types.ChunkBondRequestId, error) {
	chunkSize := k.GetParams(ctx).ChunkSize

	liquidStakingState := k.NewLiquidStakingState()
	mintTokenAmount, err := types.NativeTokenToLiquidToken(liquidStakingState, tokenAmount)
	if err != nil {
		return 0, err
	}
	if mintTokenAmount.LT(chunkSize) {
		return 0, types.ErrInvalidTokenAmount
	}

	// TODO: check speculation. for now, just deposit coins from liquidStaker into module
	bondDenom := k.stk.BondDenom(ctx)
	stake := sdk.NewCoin(bondDenom, tokenAmount)
	if err := k.bk.SendCoinsFromAccountToModule(ctx, liquidStaker, types.ModuleName, sdk.NewCoins(stake)); err != nil {
		return 0, err
	}

	id := k.GetLastChunkBondRequestId(ctx) + 1
	req := types.ChunkBondRequest{
		Id:          id,
		Address:     liquidStaker.String(),
		TokenAmount: tokenAmount,
	}

	k.SetChunkBondRequest(ctx, req)
	k.SetLastChunkBondRequestId(ctx, id)
	return id, nil
}

func (k Keeper) LiquidUnstake(ctx sdk.Context, liquidUnstaker sdk.AccAddress, numChunkUnbond uint64) (types.ChunkUnbondRequestId, error) {
	chunkSize := k.GetParams(ctx).ChunkSize

	// TODO: check speculation. for now, just deposit liquid coins from liquidUnstaker into module
	liquidBondDenom := k.LiquidBondDenom(ctx)
	liquidStake := sdk.NewCoin(liquidBondDenom, chunkSize)
	if err := k.bk.SendCoinsFromAccountToModule(ctx, liquidUnstaker, types.ModuleName, sdk.NewCoins(liquidStake)); err != nil {
		return 0, err
	}

	id := k.GetLastChunkUnbondRequestId(ctx) + 1
	req := types.ChunkUnbondRequest{
		Id:             id,
		Address:        liquidUnstaker.String(),
		NumChunkUnbond: numChunkUnbond,
	}

	k.SetChunkUnbondRequest(ctx, req)
	k.SetLastChunkUnbondRequestId(ctx, id)
	return id, nil
}
