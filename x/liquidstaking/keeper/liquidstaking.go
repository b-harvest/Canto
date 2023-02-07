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

func (k Keeper) CancelLiquidStaking(
	ctx sdk.Context, liquidStaker sdk.AccAddress, id types.ChunkBondRequestId) (interface{}, error) {
	req, found := k.GetChunkBondRequest(ctx, id)
	if !found {
		return nil, types.ErrInvalidChunkBondRequestId
	}
	requesterAddress, err := sdk.AccAddressFromBech32(req.Address)
	if err != nil {
		return nil, err
	}
	if !liquidStaker.Equals(requesterAddress) {
		return nil, types.ErrInvalidRequesterAddress
	}

	// TODO: check speculation.
	bondDenom := k.stk.BondDenom(ctx)
	stake := sdk.NewCoin(bondDenom, req.TokenAmount)
	if err := k.bk.SendCoinsFromModuleToAccount(ctx, types.ModuleName, liquidStaker, sdk.NewCoins(stake)); err != nil {
		return 0, err
	}

	k.DeleteChunkBondRequest(ctx, req)
	return nil, nil
}

func (k Keeper) LiquidUnstake(ctx sdk.Context, liquidUnstaker sdk.AccAddress, numChunkUnbond uint64) (types.ChunkUnbondRequestId, error) {
	chunkSize := k.GetParams(ctx).ChunkSize

	// TODO: check speculation. for now, just deposit liquid coins from liquidUnstaker into module
	liquidBondDenom := k.LiquidBondDenom(ctx)
	liquidStake := sdk.NewCoin(liquidBondDenom, chunkSize.Mul(sdk.NewIntFromUint64(numChunkUnbond)))
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

func (k Keeper) CancelLiquidUnstaking(
	ctx sdk.Context, liquidUnstaker sdk.AccAddress, id types.ChunkUnbondRequestId) (interface{}, error) {
	req, found := k.GetChunkUnbondRequest(ctx, id)
	if !found {
		return nil, types.ErrInvalidChunkBondRequestId
	}
	requesterAddress, err := sdk.AccAddressFromBech32(req.Address)
	if err != nil {
		return nil, err
	}
	if !liquidUnstaker.Equals(requesterAddress) {
		return nil, types.ErrInvalidRequesterAddress
	}
	// TODO: check speculation. for now, just deposit liquid coins from liquidUnstaker into module
	chunkSize := k.GetParams(ctx).ChunkSize
	liquidBondDenom := k.LiquidBondDenom(ctx)
	liquidStake := sdk.NewCoin(liquidBondDenom, chunkSize.Mul(sdk.NewIntFromUint64(req.NumChunkUnbond)))
	if err := k.bk.SendCoinsFromModuleToAccount(ctx, types.ModuleName, liquidUnstaker, sdk.NewCoins(liquidStake)); err != nil {
		return 0, err
	}

	k.DeleteChunkUnbondRequest(ctx, req)
	return nil, nil
}

func (k Keeper) getMinimumInsuranceAmount(ctx *sdk.Context) sdk.Int {
	params := k.GetParams(*ctx)
	// TODO: calc correct minimum amount
	return params.MinInsurancePercentage.Add(sdk.NewDec(100)).MulInt(params.ChunkSize).TruncateInt()
}

func (k Keeper) BidInsurance(
	ctx sdk.Context,
	insurer sdk.AccAddress,
	val sdk.ValAddress,
	amount sdk.Int,
	insuranceFeeRate sdk.Dec,
) (types.InsuranceBidId, error) {
	minimumInsuranceAmount := k.getMinimumInsuranceAmount(&ctx)
	if minimumInsuranceAmount.GT(amount) {
		return 0, types.ErrInvalidTokenAmount
	}

	// TODO: check speculation. for now, just deposit coins from insurer into module
	bondDenom := k.stk.BondDenom(ctx)
	if err := k.bk.SendCoinsFromAccountToModule(ctx, insurer, types.ModuleName, sdk.NewCoins(sdk.NewCoin(bondDenom, amount))); err != nil {
		return 0, err
	}
	id := k.GetLastInsuranceBidId(ctx) + 1
	bid := types.InsuranceBid{
		Id:                       id,
		ValidatorAddress:         val.String(),
		InsuranceProviderAddress: insurer.String(),
		InsuranceAmount:          amount,
		InsuranceFeeRate:         insuranceFeeRate,
	}

	k.SetLastInsuranceBidId(ctx, id)
	k.SetInsuranceBid(ctx, bid)
	k.SetInsuranceBidIndex(ctx, bid)

	return id, nil
}

func (k Keeper) CancelInsuranceBid(
	ctx sdk.Context, insurer sdk.AccAddress, id types.InsuranceBidId) (interface{}, error) {
	bid, found := k.GetInsuranceBid(ctx, id)
	if !found {
		return nil, types.ErrInvalidInsuranceBidId
	}

	requesterAddress, err := sdk.AccAddressFromBech32(bid.InsuranceProviderAddress)
	if err != nil {
		return nil, err
	}
	if !insurer.Equals(requesterAddress) {
		return nil, types.ErrInvalidRequesterAddress
	}
	k.DeleteInsuranceBid(ctx, bid)
	return nil, nil
}

func (k Keeper) UnbondInsurance(
	ctx sdk.Context,
	insurer sdk.AccAddress,
	aliveChunkId types.AliveChunkId,
) (types.AliveChunkId, error) {
	if _, found := k.GetAliveChunk(ctx, aliveChunkId); !found {
		return 0, types.ErrInvalidAliveChunkId
	}
	req := types.InsuranceUnbondRequest{
		InsuranceProviderAddress: insurer.String(),
		AliveChunkId:             aliveChunkId,
	}
	k.SetInsuranceUnbondRequest(ctx, req)
	k.SetInsuranceUnbondRequestIndex(ctx, req)

	return aliveChunkId, nil
}

func (k Keeper) CancelInsuranceUnbond(
	ctx sdk.Context, insurer sdk.AccAddress, id types.AliveChunkId) (interface{}, error) {
	req, found := k.GetInsuranceUnbondRequest(ctx, id)
	if !found {
		return nil, types.ErrInvalidAliveChunkId
	}

	requesterAddress, err := sdk.AccAddressFromBech32(req.InsuranceProviderAddress)
	if err != nil {
		return nil, err
	}
	if !insurer.Equals(requesterAddress) {
		return nil, types.ErrInvalidRequesterAddress
	}
	k.DeleteInsuranceUnbondRequest(ctx, req)
	return nil, nil
}
