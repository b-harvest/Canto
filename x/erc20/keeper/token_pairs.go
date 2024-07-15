package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"

	"github.com/Canto-Network/Canto/v7/x/erc20/types"
)

// GetTokenPairs - get all registered token tokenPairs
func (k Keeper) GetTokenPairs(ctx sdk.Context) []types.TokenPair {
	tokenPairs := []types.TokenPair{}

	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.KeyPrefixTokenPair)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var tokenPair types.TokenPair
		k.cdc.MustUnmarshal(iterator.Value(), &tokenPair)

		tokenPairs = append(tokenPairs, tokenPair)
	}

	return tokenPairs
}

// GetTokenPairID returns the pair id from either of the registered tokens.
func (k Keeper) GetTokenPairID(ctx sdk.Context, token string) []byte {
	if common.IsHexAddress(token) {
		addr := common.HexToAddress(token)
		return k.GetTokenPairIdByERC20Addr(ctx, addr)
	}
	return k.GetTokenPairIdByDenom(ctx, token)
}

// GetTokenPair - get registered token pair from the identifier
func (k Keeper) GetTokenPair(ctx sdk.Context, id []byte) (types.TokenPair, bool) {
	if id == nil {
		return types.TokenPair{}, false
	}

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixTokenPair)
	var tokenPair types.TokenPair
	bz := store.Get(id)
	if len(bz) == 0 {
		return types.TokenPair{}, false
	}

	k.cdc.MustUnmarshal(bz, &tokenPair)
	return tokenPair, true
}

// SetTokenPair stores a token pair
func (k Keeper) SetTokenPair(ctx sdk.Context, tokenPair types.TokenPair) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixTokenPair)
	key := tokenPair.GetID()
	bz := k.cdc.MustMarshal(&tokenPair)
	store.Set(key, bz)
}

// DeleteTokenPair removes a token pair.
func (k Keeper) DeleteTokenPair(ctx sdk.Context, tokenPair types.TokenPair) {
	id := tokenPair.GetID()
	k.deleteTokenPair(ctx, id)
	k.deleteTokenPairIdByERC20Addr(ctx, tokenPair.GetERC20Contract())
	k.deleteTokenPairIdByDenom(ctx, tokenPair.Denom)
}

// deleteTokenPair deletes the token pair for the given id
func (k Keeper) deleteTokenPair(ctx sdk.Context, id []byte) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixTokenPair)
	store.Delete(id)
}

// GetTokenPairIdByERC20Addr returns the token pair id for the given address
func (k Keeper) GetTokenPairIdByERC20Addr(ctx sdk.Context, erc20 common.Address) []byte {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixTokenPairByERC20Address)
	return store.Get(erc20.Bytes())
}

// GetTokenPairIdByDenom returns the token pair id for the given denomination
func (k Keeper) GetTokenPairIdByDenom(ctx sdk.Context, denom string) []byte {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixTokenPairByDenom)
	return store.Get([]byte(denom))
}

// SetTokenPairIdByERC20Addr sets the token pair id for the given address
func (k Keeper) SetTokenPairIdByERC20Addr(ctx sdk.Context, erc20 common.Address, id []byte) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixTokenPairByERC20Address)
	store.Set(erc20.Bytes(), id)
}

// deleteTokenPairIdByERC20Addr deletes the token pair id for the given address
func (k Keeper) deleteTokenPairIdByERC20Addr(ctx sdk.Context, erc20 common.Address) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixTokenPairByERC20Address)
	store.Delete(erc20.Bytes())
}

// SetTokenPairIdByDenom sets the token pair id for the denomination
func (k Keeper) SetTokenPairIdByDenom(ctx sdk.Context, denom string, id []byte) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixTokenPairByDenom)
	store.Set([]byte(denom), id)
}

// deleteTokenPairIdByDenom deletes the token pair id for the given denom
func (k Keeper) deleteTokenPairIdByDenom(ctx sdk.Context, denom string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixTokenPairByDenom)
	store.Delete([]byte(denom))
}

// IsTokenPairRegistered - check if registered token tokenPair is registered
func (k Keeper) IsTokenPairRegistered(ctx sdk.Context, id []byte) bool {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixTokenPair)
	return store.Has(id)
}

// IsERC20Registered check if registered ERC20 token is registered
func (k Keeper) IsERC20Registered(ctx sdk.Context, erc20 common.Address) bool {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixTokenPairByERC20Address)
	return store.Has(erc20.Bytes())
}

// IsDenomRegistered check if registered coin denom is registered
func (k Keeper) IsDenomRegistered(ctx sdk.Context, denom string) bool {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixTokenPairByDenom)
	return store.Has([]byte(denom))
}

func (k Keeper) GetAllTokenPairDenomIndexes(ctx sdk.Context) []types.TokenPairDenomIndex {
	var idxs []types.TokenPairDenomIndex
	k.IterateTokenPairDenomIndex(ctx, func(denom string, id []byte) (stop bool) {
		idx := types.TokenPairDenomIndex{
			Denom:       denom,
			TokenPairId: id,
		}
		idxs = append(idxs, idx)
		return false
	})
	return idxs
}

func (k Keeper) IterateTokenPairDenomIndex(ctx sdk.Context, cb func(denom string, id []byte) (stop bool)) {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixTokenPairByDenom)
	iter := prefixStore.Iterator(nil, nil)
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		denom := string(iter.Key())
		id := iter.Value()
		if cb(denom, id) {
			break
		}
	}
}

func (k Keeper) GetAllTokenPairERC20AddressIndexes(ctx sdk.Context) []types.TokenPairERC20AddressIndex {
	var idxs []types.TokenPairERC20AddressIndex
	k.IterateTokenPairERC20AddressIndex(ctx, func(erc20Addr common.Address, id []byte) (stop bool) {
		idx := types.TokenPairERC20AddressIndex{
			Erc20Address: erc20Addr.Bytes(),
			TokenPairId:  id,
		}
		idxs = append(idxs, idx)
		return false
	})
	return idxs
}

func (k Keeper) IterateTokenPairERC20AddressIndex(ctx sdk.Context, cb func(erc20Addr common.Address, id []byte) (stop bool)) {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixTokenPairByERC20Address)
	iter := prefixStore.Iterator(nil, nil)
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		erc20Addr := common.BytesToAddress(iter.Key())
		id := iter.Value()
		if cb(erc20Addr, id) {
			break
		}
	}
}
