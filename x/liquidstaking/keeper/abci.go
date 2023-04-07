package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func BeginBlocker(ctx sdk.Context, k Keeper) {
	// TODO: should we panic here? or continue as much as possible?
	// Need to check references (e.g. crescent BeginBlocker, EndBlocker)
	k.coverUnpairingForUnstakeCase(ctx)
	k.coverUnpairingForRepairingCase(ctx)
	k.coverVulnerableInsurances(ctx)
}
