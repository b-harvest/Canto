package liquidstaking

import (
	"github.com/Canto-Network/Canto/v6/x/liquidstaking/keeper"
	"github.com/Canto-Network/Canto/v6/x/liquidstaking/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// InitGenesis initializes the capability module's state from a provided genesis
// state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	k.SetParams(ctx, genState.Params)

}

// ExportGenesis returns the capability module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesisState()
	genesis.Params = k.GetParams(ctx)
	genesis.Epoch = k.GetEpoch(ctx)
	genesis.LastChunkId = k.GetLastChunkId(ctx)
	genesis.LastInsuranceId = k.GetLastInsuranceId(ctx)
	genesis.Chunks = k.GetChunks(ctx)
	genesis.Insurances = k.GetInsurances(ctx)
	genesis.WithdrawingInsurances = k.GetWithdrawingInsurances(ctx)
	genesis.LiquidUnstakeUnbondingDelegationInfos = k.GetLiquidUnstakeUnbondingDelegationInfos(ctx)

	return genesis
}
