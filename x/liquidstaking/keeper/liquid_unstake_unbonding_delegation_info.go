package keeper

import (
	"encoding/binary"

	"github.com/Canto-Network/Canto/v6/x/liquidstaking/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) SetLiquidUnstakeUnbondingDelegationInfo(ctx sdk.Context, liquidUnstakeUnbondingDelegationInfo types.LiquidUnstakeUnbondingDelegationInfo) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixLiquidUnstakeUnbondingDelegationInfo)
	chunkId := make([]byte, 8)
	binary.LittleEndian.PutUint64(chunkId, liquidUnstakeUnbondingDelegationInfo.ChunkId)
	bz := k.cdc.MustMarshal(&liquidUnstakeUnbondingDelegationInfo)
	store.Set(chunkId, bz)
}

func (k Keeper) GetLiquidUnstakeUnbondingDelegationInfo(ctx sdk.Context, chunkId uint64) (liquidUnstakeUnbondingDelegationInfo types.LiquidUnstakeUnbondingDelegationInfo, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixLiquidUnstakeUnbondingDelegationInfo)
	chunkIdBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(chunkIdBytes, chunkId)
	bz := store.Get(chunkIdBytes)
	if bz == nil {
		return liquidUnstakeUnbondingDelegationInfo, false
	}
	k.cdc.MustUnmarshal(bz, &liquidUnstakeUnbondingDelegationInfo)
	return liquidUnstakeUnbondingDelegationInfo, true
}

func (k Keeper) DeleteLiquidUnstakeUnbondingDelegationInfo(ctx sdk.Context, chunkId uint64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixLiquidUnstakeUnbondingDelegationInfo)
	chunkIdBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(chunkIdBytes, chunkId)
	store.Delete(chunkIdBytes)
}

func (k Keeper) IterateAllLiquidUnstakeUnbondingDelegationInfos(ctx sdk.Context, cb func(liquidUnstakeUnbondingDelegationInfo types.LiquidUnstakeUnbondingDelegationInfo) (stop bool, err error)) error {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixLiquidUnstakeUnbondingDelegationInfo)
	iterator := store.Iterator(nil, nil)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var liquidUnstakeUnbondingDelegationInfo types.LiquidUnstakeUnbondingDelegationInfo
		k.cdc.MustUnmarshal(iterator.Value(), &liquidUnstakeUnbondingDelegationInfo)
		stop, err := cb(liquidUnstakeUnbondingDelegationInfo)
		if err != nil {
			return err
		}
		if stop {
			break
		}
	}
	return nil
}
