package types

import sdk "github.com/cosmos/cosmos-sdk/types"

func NativeTokenToLiquidStakeToken(nativeTokenAmount, lsTokenTotalSupplyAmount sdk.Int, netAmount sdk.Dec) (lsTokenAmount sdk.Int) {
	return lsTokenTotalSupplyAmount.ToDec().
		QuoTruncate(netAmount.TruncateDec()).
		MulTruncate(nativeTokenAmount.ToDec()).
		TruncateInt()
}
