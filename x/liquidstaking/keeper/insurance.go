package keeper

import (
	"encoding/binary"
	"github.com/Canto-Network/Canto/v6/x/liquidstaking/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	gogotypes "github.com/gogo/protobuf/types"
)

func (k Keeper) SetInsurance(ctx sdk.Context, insurance types.Insurance) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixInsurance)
	insuranceId := make([]byte, 8)
	binary.LittleEndian.PutUint64(insuranceId, insurance.Id)
	bz := k.cdc.MustMarshal(&insurance)
	store.Set(insuranceId, bz)
}

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

func (k Keeper) DeleteInsurance(ctx sdk.Context, id uint64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixInsurance)
	insuranceId := make([]byte, 8)
	binary.LittleEndian.PutUint64(insuranceId, id)
	store.Delete(insuranceId)
}

func (k Keeper) GetInsurances(ctx sdk.Context) []types.Insurance {
	var insurances []types.Insurance

	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.KeyPrefixInsurance)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var insurance types.Insurance
		k.cdc.MustUnmarshal(iterator.Value(), &insurance)

		insurances = append(insurances, insurance)
	}

	return insurances
}

func (k Keeper) SetLastInsuranceId(ctx sdk.Context, id uint64) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&gogotypes.UInt64Value{Value: id})
	store.Set(types.KeyPrefixLastInsuranceId, bz)
}

func (k Keeper) GetLastInsuranceId(ctx sdk.Context) (id uint64) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.KeyPrefixLastInsuranceId)
	if bz == nil {
		id = 0
	} else {
		var val gogotypes.UInt64Value
		k.cdc.MustUnmarshal(bz, &val)
		id = val.GetValue()
	}
	return
}

func (k Keeper) getNextInsuranceIdWithUpdate(ctx sdk.Context) uint64 {
	id := k.GetLastInsuranceId(ctx)
	id++
	k.SetLastInsuranceId(ctx, id)
	return id
}
