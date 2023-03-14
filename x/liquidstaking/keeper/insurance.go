package keeper

import (
	"encoding/binary"
	"github.com/Canto-Network/Canto/v6/x/liquidstaking/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) GetInsurance(ctx sdk.Context, id uint64) (insurance types.Insurance, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixInsurance)
	insuranceId := make([]byte, 8)
	binary.LittleEndian.PutUint64(insuranceId, id)
	bz := store.Get(insuranceId)
	if bz == nil {
		return insurance, false
	}
	k.cdc.MustUnmarshal(bz, &insurance)
	return insurance, true
}
