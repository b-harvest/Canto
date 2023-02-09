package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/Canto-Network/Canto-Testnet-v2/v1/x/liquidstaking/types"
)

// GetParams returns the total set of liquidstaking parameters.
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	k.paramstore.GetParamSet(ctx, &params)
	return params
}

// SetParams sets the liquidstaking parameters to the param space.
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramstore.SetParamSet(ctx, &params)
}
