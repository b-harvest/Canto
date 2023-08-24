package govshuttle

import (
	"github.com/Canto-Network/Canto/v7/x/govshuttle/keeper"
	"github.com/Canto-Network/Canto/v7/x/govshuttle/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	"github.com/ethereum/go-ethereum/common"
)

// InitGenesis initializes the capability module's state from a provided genesis
// state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, accountKeeper authkeeper.AccountKeeper, genState types.GenesisState) {
	// this line is used by starport scaffolding # genesis/module/init
	k.SetParams(ctx, genState.Params)
	k.SetPort(ctx, common.HexToAddress(genState.PortAddress))
	if acc := accountKeeper.GetModuleAccount(ctx, types.ModuleName); acc == nil {
		panic("the govshuttle module account has not been set")
	}

}

// ExportGenesis returns the capability module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	portAddress, _ := k.GetPort(ctx)
	genesis := types.NewGenesisState(k.GetParams(ctx), portAddress)

	// this line is used by starport scaffolding # genesis/module/export

	return genesis
}
