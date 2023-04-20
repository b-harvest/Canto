package keeper

import (
	"time"

	"github.com/Canto-Network/Canto/v6/x/liquidstaking/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) GetLUQueueTimeSlice(ctx sdk.Context, timestamp time.Time) []types.LiquidUnstakeEntry {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetPendingLiquidStakeTimeKey(timestamp))
	if bz == nil {
		return []types.LiquidUnstakeEntry{}
	}
	var pendingLiquidUnstake types.PendingLiquidUnstake
	k.cdc.MustUnmarshal(bz, &pendingLiquidUnstake)
	return pendingLiquidUnstake.Entries
}

func (k Keeper) SetLUQueueTimeSlice(ctx sdk.Context, timestamp time.Time, entries []types.LiquidUnstakeEntry) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&types.PendingLiquidUnstake{Entries: entries})
	store.Set(types.GetPendingLiquidStakeTimeKey(timestamp), bz)
}

func (k Keeper) InsertLUQueue(
	ctx sdk.Context,
	timestamp time.Time,
	chunkId uint64,
	delegatorAddress sdk.AccAddress,
	escrowedLsTokens sdk.Coin,
) {
	entry := types.LiquidUnstakeEntry{
		ChunkId:          chunkId,
		DelegatorAddress: delegatorAddress.String(),
		EscrowedLstokens: escrowedLsTokens,
	}
	timeSlice := k.GetLUQueueTimeSlice(ctx, timestamp)
	if len(timeSlice) == 0 {
		k.SetLUQueueTimeSlice(ctx, timestamp, []types.LiquidUnstakeEntry{entry})
	} else {
		timeSlice = append(timeSlice, entry)
		k.SetLUQueueTimeSlice(ctx, timestamp, timeSlice)
	}
}

func (k Keeper) LUQueueIterator(ctx sdk.Context, endTime time.Time) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return store.Iterator(
		types.KeyPrefixLiquidUnstakeQueueKey,
		sdk.InclusiveEndBytes(types.GetPendingLiquidStakeTimeKey(endTime)),
	)
}

func (k Keeper) DequeueAllLiquidUnstakeEntry(ctx sdk.Context, currTime time.Time) (entries []types.LiquidUnstakeEntry) {
	store := ctx.KVStore(k.storeKey)

	iterator := k.LUQueueIterator(ctx, currTime)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var pendingLiquidUnstake types.PendingLiquidUnstake
		k.cdc.MustUnmarshal(iterator.Value(), &pendingLiquidUnstake)
		entries = append(entries, pendingLiquidUnstake.Entries...)
		store.Delete(iterator.Key())
	}
	return
}

func (k Keeper) GetAllPendingLiquidUnstake(ctx sdk.Context) []types.PendingLiquidUnstake {
	var pendingLiquidUnstakes []types.PendingLiquidUnstake
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.KeyPrefixLiquidUnstakeQueueKey)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var pendingLiquidUnstake types.PendingLiquidUnstake
		k.cdc.MustUnmarshal(iterator.Value(), &pendingLiquidUnstake)
		pendingLiquidUnstakes = append(pendingLiquidUnstakes, pendingLiquidUnstake)
	}
	return pendingLiquidUnstakes
}
