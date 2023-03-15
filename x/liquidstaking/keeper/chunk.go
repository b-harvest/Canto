package keeper

import (
	"encoding/binary"
	"github.com/Canto-Network/Canto/v6/x/liquidstaking/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	gogotypes "github.com/gogo/protobuf/types"
)

func (k Keeper) SetChunk(ctx sdk.Context, chunk types.Chunk) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixChunk)
	chunkId := make([]byte, 8)
	binary.LittleEndian.PutUint64(chunkId, chunk.Id)
	bz := k.cdc.MustMarshal(&chunk)
	store.Set(chunkId, bz)
}

func (k Keeper) GetChunk(ctx sdk.Context, id uint64) (chunk types.Chunk, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixChunk)
	chunkId := make([]byte, 8)
	binary.LittleEndian.PutUint64(chunkId, id)
	bz := store.Get(chunkId)
	if bz == nil {
		return chunk, false
	}
	k.cdc.MustUnmarshal(bz, &chunk)
	return chunk, true
}

func (k Keeper) DeleteChunk(ctx sdk.Context, id uint64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixChunk)
	chunkId := make([]byte, 8)
	binary.LittleEndian.PutUint64(chunkId, id)
	store.Delete(chunkId)
}

func (k Keeper) GetChunks(ctx sdk.Context) []types.Chunk {
	var chunks []types.Chunk

	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.KeyPrefixChunk)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var chunk types.Chunk
		k.cdc.MustUnmarshal(iterator.Value(), &chunk)

		chunks = append(chunks, chunk)
	}

	return chunks
}

func (k Keeper) SetLastChunkId(ctx sdk.Context, id uint64) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&gogotypes.UInt64Value{Value: id})
	store.Set(types.KeyPrefixLastChunkId, bz)
}

func (k Keeper) GetLastChunkId(ctx sdk.Context) (id uint64) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.KeyPrefixLastChunkId)
	if bz == nil {
		id = 0
	} else {
		var val gogotypes.UInt64Value
		k.cdc.MustUnmarshal(bz, &val)
		id = val.GetValue()
	}
	return
}

func (k Keeper) getNextChunkIdWithUpdate(ctx sdk.Context) uint64 {
	id := k.GetLastChunkId(ctx) + 1
	k.SetLastChunkId(ctx, id+1)
	return id
}
