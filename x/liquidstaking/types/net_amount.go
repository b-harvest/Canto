package types

import sdk "github.com/cosmos/cosmos-sdk/types"

func (nas NetAmountState) CalcNetAmount(rewardPoolBalance sdk.Int) sdk.Dec {
	// TODO: Add reward module account's balance
	return rewardPoolBalance.Add(nas.TotalChunksBalance).
		Add(nas.TotalLiquidTokens).
		Add(nas.TotalUnbondingBalance).ToDec().
		Add(nas.TotalRemainingRewards)
}

func (nas NetAmountState) CalcMintRate() sdk.Dec {
	if nas.NetAmount.IsNil() || !nas.NetAmount.IsPositive() {
		return sdk.ZeroDec()
	}
	return nas.LsTokensTotalSupply.ToDec().QuoTruncate(nas.NetAmount)
}
