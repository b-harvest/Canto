package keeper

import (
	"github.com/Canto-Network/Canto/v6/x/liquidstaking/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

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
