package keeper

import (
	gogotypes "github.com/gogo/protobuf/types"

	"github.com/Canto-Network/Canto/v6/x/liquidstaking/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

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
