package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/Canto-Network/Canto/v7/x/coinswap/types"
)

// InitGenesis initializes the coinswap module's state from a given genesis state.
func (k Keeper) InitGenesis(ctx sdk.Context, genState types.GenesisState) {
	if err := types.ValidateGenesis(genState); err != nil {
		panic(fmt.Errorf("panic for ValidateGenesis,%v", err))
	}

	// init to prevent nil slice, []types.WhitelistedValidator(nil)
	//if genState.Params.WhitelistedDenoms == nil || len(genState.Params.WhitelistedDenoms) == 0 {
	//	genState.Params.WhitelistedDenoms = []string{}
	//}
	k.SetParams(ctx, genState.Params)
	k.SetStandardDenom(ctx, genState.StandardDenom)
	k.setSequence(ctx, genState.Sequence)
	for _, pool := range genState.Pool {
		k.setPool(ctx, &pool)
	}
}

// ExportGenesis returns the coinswap module's genesis state.
func (k Keeper) ExportGenesis(ctx sdk.Context) types.GenesisState {

	params := k.GetParams(ctx)
	// init to prevent nil slice, []types.WhitelistedValidator(nil)
	//if params.WhitelistedDenoms == nil || len(params.WhitelistedDenoms) == 0 {
	//	params.WhitelistedDenoms = []string{}
	//}
	return types.GenesisState{
		Params:        params,
		StandardDenom: k.GetStandardDenom(ctx),
		Pool:          k.GetAllPools(ctx),
		Sequence:      k.getSequence(ctx),
	}
}
