package keeper

import (
	"github.com/Canto-Network/Canto/v6/x/liquidstaking/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) GetLiquidUnstakeUnbondingDelegationInfos(ctx sdk.Context) []types.LiquidUnstakeUnbondingDelegationInfo {
	var liquidUnstakeUnbondingDelegations []types.LiquidUnstakeUnbondingDelegationInfo

	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.KeyPrefixLiquidUnstakeUnbondingDelegation)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var liquidUnstakeUnbondingDelegationInfo types.LiquidUnstakeUnbondingDelegationInfo
		k.cdc.MustUnmarshal(iterator.Value(), &liquidUnstakeUnbondingDelegationInfo)

		liquidUnstakeUnbondingDelegations = append(liquidUnstakeUnbondingDelegations, liquidUnstakeUnbondingDelegationInfo)
	}

	return liquidUnstakeUnbondingDelegations
}
