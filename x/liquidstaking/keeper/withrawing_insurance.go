package keeper

import (
	"encoding/binary"
	"github.com/Canto-Network/Canto/v6/x/liquidstaking/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) SetWithdrawingInsurance(ctx sdk.Context, withdrawingInsurance types.WithdrawingInsurance) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixWithdrawingInsurance)
	withdrawingInsuranceId := make([]byte, 8)
	binary.LittleEndian.PutUint64(withdrawingInsuranceId, withdrawingInsurance.InsuranceId)
	bz := k.cdc.MustMarshal(&withdrawingInsurance)
	store.Set(withdrawingInsuranceId, bz)
}

func (k Keeper) GetWithdrawingInsurance(ctx sdk.Context, id uint64) (withdrawingInsurance types.WithdrawingInsurance, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixWithdrawingInsurance)
	withdrawingInsuranceId := make([]byte, 8)
	binary.LittleEndian.PutUint64(withdrawingInsuranceId, id)
	bz := store.Get(withdrawingInsuranceId)
	if bz == nil {
		return withdrawingInsurance, false
	}
	k.cdc.MustUnmarshal(bz, &withdrawingInsurance)
	return withdrawingInsurance, true
}

func (k Keeper) DeleteWithdrawingInsurance(ctx sdk.Context, id uint64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixWithdrawingInsurance)
	withdrawingInsuranceId := make([]byte, 8)
	binary.LittleEndian.PutUint64(withdrawingInsuranceId, id)
	store.Delete(withdrawingInsuranceId)
}

func (k Keeper) GetWithdrawingInsurances(ctx sdk.Context) []types.WithdrawingInsurance {
	var withdrawingInsurances []types.WithdrawingInsurance

	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.KeyPrefixWithdrawingInsurance)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var withdrawingInsurance types.WithdrawingInsurance
		k.cdc.MustUnmarshal(iterator.Value(), &withdrawingInsurance)

		withdrawingInsurances = append(withdrawingInsurances, withdrawingInsurance)
	}

	return withdrawingInsurances
}
