package keeper

import (
	"github.com/Canto-Network/Canto/v6/x/liquidstaking/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func BeginBlocker(ctx sdk.Context, k Keeper) {
	var unbondingDelegationInfos []types.LiquidUnstakeUnbondingDelegationInfo
	err := k.IterateAllLiquidUnstakeUnbondingDelegationInfos(ctx, func(liquidUnstakeUnbondingDelegationInfo types.LiquidUnstakeUnbondingDelegationInfo) (bool, error) {
		unbondingDelegationInfos = append(unbondingDelegationInfos, liquidUnstakeUnbondingDelegationInfo)
		return false, nil
	})
	if err != nil {
		panic(err)
	}

	// For all completed unboding delegation infos
	for _, unbondingDelegationInfo := range unbondingDelegationInfos {
		chunk, found := k.GetChunk(ctx, unbondingDelegationInfo.ChunkId)
		if !found {
			panic(types.ErrNotFoundChunk.Error())
		}

		// get insurance from chunk.InsuranceId
		insurance, found := k.GetInsurance(ctx, chunk.InsuranceId)
		if !found {
			panic(types.ErrNotFoundInsurance.Error())
		}

		// get unbonding delegation using staking keeper
		unbondingDelegation, found := k.stakingKeeper.GetUnbondingDelegation(
			ctx, chunk.DerivedAddress(),
			insurance.GetValidator(),
		)
		if !found {
			panic(types.ErrNotFoundUnbondingDelegation.Error())
		}

		// check if chunk got damaged during unbonding
		// for all entries of unbondingDelegation
		for _, entry := range unbondingDelegation.Entries {
			if entry.CompletionTime.Equal(unbondingDelegationInfo.CompletionTime) &&
				entry.InitialBalance.Equal(types.ChunkSize) {
				// unbonding chunk got damaged?
				diff := entry.InitialBalance.Sub(entry.Balance)
				if diff.IsPositive() {
					// chunk got damaged, insurance should cover it
				}
			}

		}
		// if entry is not completed
		// if entry.IsMature(ctx.BlockTime()) {
		// then chunk got damaged
		// }

		// Check if chunk got damaged or not

		k.DeleteChunk(ctx, chunk.Id)
	}
}
