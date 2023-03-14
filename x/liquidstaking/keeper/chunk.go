package keeper

import (
	"encoding/binary"
	"github.com/Canto-Network/Canto/v6/x/liquidstaking/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

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

func (k Keeper) CreateChunk(ctx sdk.Context, chunk types.Chunk) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixChunk)
	chunkId := make([]byte, 8)
	binary.LittleEndian.PutUint64(chunkId, chunk.Id)
	bz := k.cdc.MustMarshal(&chunk)
	store.Set(chunkId, bz)
}
