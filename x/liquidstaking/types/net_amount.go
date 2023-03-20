package types

import sdk "github.com/cosmos/cosmos-sdk/types"

func (nas NetAmountState) CalcNetAmount() sdk.Dec {
	// TODO: Add reward module account's balance
	return nas.TotalChunksBalance.Add(nas.TotalLiquidTokens).
		Add(nas.TotalUnbondingBalance).ToDec().
		Add(nas.TotalRemainingRewards)
}

func (nas NetAmountState) CalcMintRate() sdk.Dec {
	if nas.NetAmount.IsNil() || !nas.NetAmount.IsPositive() {
		return sdk.ZeroDec()
	}
	return nas.LsTokensTotalSupply.ToDec().QuoTruncate(nas.NetAmount)
}
