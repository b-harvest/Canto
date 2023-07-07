package liquidstaking

import (
	"github.com/Canto-Network/Canto/v6/x/liquidstaking/keeper"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func BeginBlocker(ctx sdk.Context, k keeper.Keeper) {
	if k.IsEpochReached(ctx) {
		k.CoverRedelegationPenalty(ctx)
	}
}

func EndBlocker(ctx sdk.Context, k keeper.Keeper) {
	if k.IsEpochReached(ctx) {
		k.DistributeReward(ctx)
		k.DeleteMaturedRedelegationInfos(ctx)
		k.CoverSlashingAndHandleMatureUnbondings(ctx)
		k.HandleQueuedLiquidUnstakes(ctx)
		k.HandleUnprocessedQueuedLiquidUnstakes(ctx)
		k.HandleQueuedWithdrawInsuranceRequests(ctx)
		newlyRankedInInsurances, rankOutInsurances, err := k.RankInsurances(ctx)
		if err != nil {
			panic(err)
		}
		if err = k.RePairRankedInsurances(ctx, newlyRankedInInsurances, rankOutInsurances); err != nil {
			panic(err)
		}
		k.IncrementEpoch(ctx)
	}
}
