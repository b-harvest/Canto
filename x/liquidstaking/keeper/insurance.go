package keeper

import (
	"github.com/Canto-Network/Canto/v6/x/liquidstaking/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	gogotypes "github.com/gogo/protobuf/types"
)

func (k Keeper) SetInsurance(ctx sdk.Context, insurance types.Insurance) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&insurance)
	store.Set(types.GetInsuranceKey(insurance.Id), bz)
}

func (k Keeper) GetInsurance(ctx sdk.Context, id uint64) (insurance types.Insurance, found bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetInsuranceKey(id))
	if bz == nil {
		return insurance, false
	}
	k.cdc.MustUnmarshal(bz, &insurance)
	return insurance, true
}

func (k Keeper) DeleteInsurance(ctx sdk.Context, id uint64) {
	store := ctx.KVStore(k.storeKey)
	insurance, _ := k.GetInsurance(ctx, id)
	store.Delete(types.GetInsuranceKey(insurance.Id))
}

func (k Keeper) getPairingInsurances(ctx sdk.Context) (
	pairingInsurances []types.Insurance,
	validatorMap map[string]stakingtypes.Validator,
) {
	validatorMap = make(map[string]stakingtypes.Validator)
	k.IterateAllInsurances(ctx, func(insurance types.Insurance) bool {
		if insurance.Status != types.INSURANCE_STATUS_PAIRING {
			return false
		}
		if _, ok := validatorMap[insurance.ValidatorAddress]; !ok {
			validator, found := k.stakingKeeper.GetValidator(ctx, insurance.GetValidator())
			err := k.IsValidValidator(ctx, validator, found)
			if err != nil {
				return false
			}
			validatorMap[insurance.ValidatorAddress] = validator
		}
		pairingInsurances = append(pairingInsurances, insurance)
		return false
	})
	return
}

func (k Keeper) IterateAllInsurances(ctx sdk.Context, cb func(insurance types.Insurance) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.KeyPrefixInsurance)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var insurance types.Insurance
		k.cdc.MustUnmarshal(iterator.Value(), &insurance)

		stop := cb(insurance)
		if stop {
			break
		}
	}
}

func (k Keeper) GetAllInsurances(ctx sdk.Context) (insurances []types.Insurance) {
	insurances = []types.Insurance{}
	k.IterateAllInsurances(ctx, func(insurance types.Insurance) (stop bool) {
		insurances = append(insurances, insurance)
		return false
	})
	return
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
	id := k.GetLastInsuranceId(ctx) + 1
	k.SetLastInsuranceId(ctx, id)
	return id
}
