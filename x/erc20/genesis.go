package erc20

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	"github.com/ethereum/go-ethereum/common"

	"github.com/Canto-Network/Canto/v7/x/erc20/keeper"
	"github.com/Canto-Network/Canto/v7/x/erc20/types"
)

// InitGenesis import module genesis
func InitGenesis(
	ctx sdk.Context,
	k keeper.Keeper,
	accountKeeper authkeeper.AccountKeeper,
	data types.GenesisState,
) {
	k.SetParams(ctx, data.Params)

	// ensure erc20 module account is set on genesis
	if acc := accountKeeper.GetModuleAccount(ctx, types.ModuleName); acc == nil {
		// NOTE: shouldn't occur
		panic("the erc20 module account has not been set")
	}

	// create token pairs by id
	tokenPairById := map[string]types.TokenPair{}
	for _, pair := range data.TokenPairs {
		id := pair.GetID()
		tokenPairById[string(id)] = pair
	}

	// set token pairs and indexes first
	// multiple contracts at the same denom can exist,
	// but only one which is in indexes are valid.
	for _, idx := range data.DenomIndexes {
		id := idx.GetTokenPairId()
		k.SetTokenPairIdByDenom(ctx, idx.Denom, id)
		k.SetTokenPair(ctx, tokenPairById[string(id)])
		delete(tokenPairById, string(id))
	}
	for _, idx := range data.Erc20AddressIndexes {
		id := idx.GetTokenPairId()
		k.SetTokenPairIdByERC20Addr(ctx, common.BytesToAddress(idx.Erc20Address), id)
		k.SetTokenPair(ctx, tokenPairById[string(id)])
		delete(tokenPairById, string(id))
	}

	// set remaining token pairs (if any left which is not in indexes)
	for _, pair := range tokenPairById {
		k.SetTokenPair(ctx, pair)
	}
}

// ExportGenesis export module status
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	return &types.GenesisState{
		Params:              k.GetParams(ctx),
		TokenPairs:          k.GetTokenPairs(ctx),
		DenomIndexes:        k.GetAllTokenPairDenomIndexes(ctx),
		Erc20AddressIndexes: k.GetAllTokenPairERC20AddressIndexes(ctx),
	}
}
