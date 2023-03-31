package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func BeginBlocker(ctx sdk.Context, k Keeper) {
	// TODO: Need to define the sequence of logic in BeginBlocker spec first.
	// TODO: rank paired chunks at epoch
	// TODO: cover slashing case: paired chunk case

	// TODO: should we panic here? or continue as much as possible?
	// Need to check references (e.g. crescent BeginBlocker, EndBlocker)
	err := k.coverUnpairingForUnstakeCase(ctx)
	if err != nil {
		panic(err)
	}
}
