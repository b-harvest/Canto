package keeper

import (
	"github.com/Canto-Network/Canto-Testnet-v2/v1/x/liquidstaking/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) GetLastChunkBondRequestId(ctx sdk.Context) types.ChunkBondRequestId {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.KeyPrefixChunkBondRequestId)
	if bz == nil {
		return 0
	}

	return sdk.BigEndianToUint64(bz)
}

func (k Keeper) SetLastChunkBondRequestId(ctx sdk.Context, id types.ChunkBondRequestId) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.KeyPrefixChunkBondRequestId, sdk.Uint64ToBigEndian(id))
}

func (k Keeper) GetLastChunkUnbondRequestId(ctx sdk.Context) types.ChunkUnbondRequestId {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.KeyPrefixChunkUnbondRequestId)
	if bz == nil {
		return 0
	}

	return sdk.BigEndianToUint64(bz)
}

func (k Keeper) SetLastChunkUnbondRequestId(ctx sdk.Context, id types.ChunkUnbondRequestId) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.KeyPrefixChunkUnbondRequestId, sdk.Uint64ToBigEndian(id))
}

func (k Keeper) GetLastInsuranceBidId(ctx sdk.Context) types.InsuranceBidId {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.KeyPrefixInsuranceBidId)
	if bz == nil {
		return 0
	}

	return sdk.BigEndianToUint64(bz)
}

func (k Keeper) SetLastInsuranceBidId(ctx sdk.Context, id types.InsuranceBidId) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.KeyPrefixInsuranceBidId, sdk.Uint64ToBigEndian(id))
}

func (k Keeper) GetLastAliveChunkId(ctx sdk.Context) types.AliveChunkId {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.KeyPrefixAliveChunkId)
	if bz == nil {
		return 0
	}

	return sdk.BigEndianToUint64(bz)
}

func (k Keeper) SetLastAliveChunkId(ctx sdk.Context, id types.AliveChunkId) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.KeyPrefixAliveChunkId, sdk.Uint64ToBigEndian(id))
}

func (k Keeper) GetLastUnbondingChunkId(ctx sdk.Context) types.UnbondingChunkId {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.KeyPrefixAliveChunkId)
	if bz == nil {
		return 0
	}

	return sdk.BigEndianToUint64(bz)
}

func (k Keeper) SetLastUnbondingChunkId(ctx sdk.Context, id types.UnbondingChunkId) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.KeyPrefixAliveChunkId, sdk.Uint64ToBigEndian(id))
}

func (k Keeper) SetChunkBondRequest(ctx sdk.Context, req types.ChunkBondRequest) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&req)
	store.Set(types.GetChunkBondRequestKey(req.Id), bz)
}

func (k Keeper) GetChunkBondRequest(ctx sdk.Context, id types.ChunkBondRequestId) (req types.ChunkBondRequest, found bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetChunkBondRequestKey(id))
	if bz == nil {
		return req, false
	}

	k.cdc.MustUnmarshal(bz, &req)
	return req, true
}

func (k Keeper) DeleteChunkBondRequest(ctx sdk.Context, req types.ChunkBondRequest) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetChunkBondRequestKey(req.Id))
}

func (k Keeper) SetChunkUnbondRequest(ctx sdk.Context, req types.ChunkUnbondRequest) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&req)
	store.Set(types.GetChunkUnbondRequestKey(req.Id), bz)
}

func (k Keeper) GetChunkUnbondRequest(ctx sdk.Context, id types.ChunkUnbondRequestId) (req types.ChunkUnbondRequest, found bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetChunkUnbondRequestKey(id))
	if bz == nil {
		return req, false
	}

	k.cdc.MustUnmarshal(bz, &req)
	return req, true
}

func (k Keeper) DeleteChunkUnbondRequest(ctx sdk.Context, req types.ChunkUnbondRequest) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetChunkUnbondRequestKey(req.Id))
}

func (k Keeper) SetInsuranceBid(ctx sdk.Context, bid types.InsuranceBid) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&bid)
	store.Set(types.GetInsuranceBidKey(bid.Id), bz)
}

func (k Keeper) GetInsuranceBid(ctx sdk.Context, id types.InsuranceBidId) (bid types.InsuranceBid, found bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetInsuranceBidKey(id))
	if bz == nil {
		return bid, false
	}

	k.cdc.MustUnmarshal(bz, &bid)
	return bid, true
}

func (k Keeper) DeleteInsuranceBid(ctx sdk.Context, bid types.InsuranceBid) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetInsuranceBidKey(bid.Id))
}

func (k Keeper) SetInsuranceUnbondRequest(ctx sdk.Context, req types.InsuranceUnbondRequest) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&req)
	store.Set(types.GetInsuranceUnbondRequestKey(req.AliveChunkId), bz)
}

func (k Keeper) GetInsuranceUnbondRequest(ctx sdk.Context, id types.AliveChunkId) (req types.InsuranceUnbondRequest, found bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetInsuranceUnbondRequestKey(id))
	if bz == nil {
		return req, false
	}

	k.cdc.MustUnmarshal(bz, &req)
	return req, true
}

func (k Keeper) DeleteInsuranceUnbondRequest(ctx sdk.Context, req types.InsuranceUnbondRequest) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetInsuranceUnbondRequestKey(req.AliveChunkId))
}

func (k Keeper) SetAliveChunk(ctx sdk.Context, aliveChunk types.AliveChunk) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&aliveChunk)
	store.Set(types.GetAliveChunkKey(aliveChunk.Id), bz)
}

func (k Keeper) GetAliveChunk(ctx sdk.Context, id types.AliveChunkId) (aliveChunk types.AliveChunk, found bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetAliveChunkKey(id))
	if bz == nil {
		return aliveChunk, false
	}

	k.cdc.MustUnmarshal(bz, &aliveChunk)
	return aliveChunk, true
}

func (k Keeper) DeleteAliveChunk(ctx sdk.Context, aliveChunk types.AliveChunk) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetAliveChunkKey(aliveChunk.Id))
}

func (k Keeper) SetUnbondingChunk(ctx sdk.Context, unbondingChunk types.UnbondingChunk) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&unbondingChunk)
	store.Set(types.GetUnbondingChunkKey(unbondingChunk.Id), bz)
}

func (k Keeper) GetUnbondingChunk(ctx sdk.Context, id types.UnbondingChunkId) (unbondingChunk types.UnbondingChunk, found bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetUnbondingChunkKey(id))
	if bz == nil {
		return unbondingChunk, false
	}

	k.cdc.MustUnmarshal(bz, &unbondingChunk)
	return unbondingChunk, true
}

func (k Keeper) DeleteUnbondingChunk(ctx sdk.Context, unbondingChunk types.UnbondingChunk) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetUnbondingChunkKey(unbondingChunk.Id))
}

func (k Keeper) iterateAliveChunks(ctx sdk.Context, cb func(aliveChunk types.AliveChunk) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.KeyPrefixAliveChunk)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var aliveChunk types.AliveChunk
		k.cdc.MustUnmarshal(iter.Value(), &aliveChunk)
		if cb(aliveChunk) {
			break
		}
	}
}

func (k Keeper) GetAllAliveChunks(ctx sdk.Context) (ret types.AliveChunks) {
	k.iterateAliveChunks(ctx, func(aliveChunk types.AliveChunk) (stop bool) {
		ret = append(ret, aliveChunk)
		return false
	})
	return
}

func (k Keeper) iterateChunkBondRequests(ctx sdk.Context, cb func(req types.ChunkBondRequest) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.KeyPrefixChunkBondRequest)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var req types.ChunkBondRequest
		k.cdc.MustUnmarshal(iter.Value(), &req)
		if cb(req) {
			break
		}
	}
}

func (k Keeper) GetAllChunkBondRequests(ctx sdk.Context) (ret types.ChunkBondRequests) {
	k.iterateChunkBondRequests(ctx, func(req types.ChunkBondRequest) (stop bool) {
		ret = append(ret, req)
		return false
	})
	return
}

func (k Keeper) iterateChunkUnbondRequests(ctx sdk.Context, cb func(req types.ChunkUnbondRequest) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.KeyPrefixChunkUnbondRequest)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var req types.ChunkUnbondRequest
		k.cdc.MustUnmarshal(iter.Value(), &req)
		if cb(req) {
			break
		}
	}
}

func (k Keeper) GetAllChunkUnbondRequests(ctx sdk.Context) (ret types.ChunkUnbondRequests) {
	k.iterateChunkUnbondRequests(ctx, func(req types.ChunkUnbondRequest) (stop bool) {
		ret = append(ret, req)
		return false
	})
	return
}

func (k Keeper) iterateInsuranceBids(ctx sdk.Context, cb func(req types.InsuranceBid) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.KeyPrefixInsuranceBid)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var req types.InsuranceBid
		k.cdc.MustUnmarshal(iter.Value(), &req)
		if cb(req) {
			break
		}
	}
}

func (k Keeper) GetAllInsuranceBids(ctx sdk.Context) (ret types.InsuranceBids) {
	k.iterateInsuranceBids(ctx, func(req types.InsuranceBid) (stop bool) {
		ret = append(ret, req)
		return false
	})
	return
}

func (k Keeper) iterateInsuranceUnbondRequests(ctx sdk.Context, cb func(req types.InsuranceUnbondRequest) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.KeyPrefixInsuranceUnbondRequest)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var req types.InsuranceUnbondRequest
		k.cdc.MustUnmarshal(iter.Value(), &req)
		if cb(req) {
			break
		}
	}
}

func (k Keeper) GetAllInsuranceUnbondRequests(ctx sdk.Context) (ret types.InsuranceUnbondRequests) {
	k.iterateInsuranceUnbondRequests(ctx, func(req types.InsuranceUnbondRequest) (stop bool) {
		ret = append(ret, req)
		return false
	})
	return
}
