package keeper

import (
	"encoding/binary"
	"github.com/Canto-Network/Canto/v6/x/liquidstaking/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) SetUnpairingForRepairingInfo(ctx sdk.Context, unpairingForRepairingInfo types.UnpairingForRepairingInfo) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixUnpairingForRepairingInfo)
	chunkId := make([]byte, 8)
	binary.LittleEndian.PutUint64(chunkId, unpairingForRepairingInfo.ChunkId)
	bz := k.cdc.MustMarshal(&unpairingForRepairingInfo)
	store.Set(chunkId, bz)
}

func (k Keeper) GetUnpairingForRepairingInfo(ctx sdk.Context, chunkId uint64) (unpairingForRepairingInfo types.UnpairingForRepairingInfo, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixUnpairingForRepairingInfo)
	chunkIdBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(chunkIdBytes, chunkId)
	bz := store.Get(chunkIdBytes)
	if bz == nil {
		return unpairingForRepairingInfo, false
	}
	k.cdc.MustUnmarshal(bz, &unpairingForRepairingInfo)
	return unpairingForRepairingInfo, true
}
