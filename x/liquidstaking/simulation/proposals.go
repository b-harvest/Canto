package simulation

import (
	"github.com/cosmos/cosmos-sdk/x/staking"
	"math/rand"

	"github.com/Canto-Network/Canto/v6/app/params"
	inflationkeeper "github.com/Canto-Network/Canto/v6/x/inflation/keeper"
	"github.com/Canto-Network/Canto/v6/x/liquidstaking/keeper"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	"github.com/cosmos/cosmos-sdk/x/simulation"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
)

const (
	OpWeightSimulateUpdateDynamicFeeRateProposal = "op_weight_simulate_update_dynamic_fee_rate_proposal"
	OpWeightSimulateUpdateMaximumDiscountRate    = "op_weight_simulate_update_maximum_discount_rate"
	OpWeightSimulateAdvanceEpoch                 = "op_weight_simulate_advance_epoch"
)

func ProposalContents(
	k keeper.Keeper,
	ak authkeeper.AccountKeeper,
	bk bankkeeper.BaseKeeper,
	sk stakingkeeper.Keeper,
	dk distrkeeper.Keeper,
	ik inflationkeeper.Keeper,
) []simtypes.WeightedProposalContent {
	return []simtypes.WeightedProposalContent{
		simulation.NewWeightedProposalContent(
			OpWeightSimulateUpdateDynamicFeeRateProposal,
			params.DefaultWeightUpdateDynamicFeeRateProposal,
			SimulateUpdateDynamicFeeRateProposal(k),
		),
		simulation.NewWeightedProposalContent(
			OpWeightSimulateUpdateMaximumDiscountRate,
			params.DefaultWeightUpdateMaximumDiscountRate,
			SimulateUpdateMaximumDiscountRate(k),
		),
		simulation.NewWeightedProposalContent(
			OpWeightSimulateAdvanceEpoch,
			params.DefaultWeightAdvanceEpoch,
			SimulateAdvanceEpoch(k, ak, bk, sk, dk, ik),
		),
	}
}

// SimulateUpdateDynamicFeeRateProposal generates random update dynamic fee rate param change proposal content.
func SimulateUpdateDynamicFeeRateProposal(k keeper.Keeper) simtypes.ContentSimulatorFn {
	return func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) simtypes.Content {
		params := k.GetParams(ctx)
		params.DynamicFeeRate = genDynamicFeeRate(r)
		k.SetParams(ctx, params)
		return nil
	}
}

// SimulateUpdateMaximumDiscountRate generates random update maximum discount rate param change proposal content.
func SimulateUpdateMaximumDiscountRate(k keeper.Keeper) simtypes.ContentSimulatorFn {
	return func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) simtypes.Content {
		params := k.GetParams(ctx)
		params.MaximumDiscountRate = genMaximumDiscountRate(r)
		k.SetParams(ctx, params)
		return nil
	}
}

func SimulateAdvanceEpoch(
	k keeper.Keeper,
	ak authkeeper.AccountKeeper,
	bk bankkeeper.BaseKeeper,
	sk stakingkeeper.Keeper,
	dk distrkeeper.Keeper,
	ik inflationkeeper.Keeper,
) simtypes.ContentSimulatorFn {
	return func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) simtypes.Content {
		// BEGIN BLOCK
		lsmEpoch := k.GetEpoch(ctx)
		if ctx.BlockHeight() <= lsmEpoch.StartHeight {
			// already advanced epoch
			return nil
		}
		ctx = ctx.WithBlockTime(lsmEpoch.StartTime.Add(lsmEpoch.Duration))
		staking.BeginBlocker(ctx, sk)

		// mimic the begin block logic of epoch module
		// currently epoch module use hooks when begin block and inflation module
		// implemented that hook, so actual logic is in inflation module.
		{
			_, found := ik.GetEpochMintProvision(ctx)
			if !found {
				panic("epoch mint provision not found")
			}
			mintedCoin := sdk.NewCoin(sdk.DefaultBondDenom, sdk.TokensFromConsensusPower(100, sdk.DefaultPowerReduction))
			_, _, err := ik.MintAndAllocateInflation(ctx, mintedCoin)
			if err != nil {
				panic(err)
			}
		}

		//feeCollector := ak.GetModuleAccount(ctx, authtypes.FeeCollectorName)
		// mimic the begin block logic of distribution module
		//{
		//	feeCollectorCoins := bk.GetAllBalances(ctx, feeCollector.GetAddress())
		//	feeCollected := sdk.NewDecCoinsFromCoins(feeCollectorCoins...)
		//	remaining := feeCollected
		//
		//	// mimic distribution.BeginBlock (AllocateTokens, get rewards from feeCollector, AllocateTokensToValidator, add remaining to feePool)
		//	err := bk.SendCoinsFromModuleToModule(ctx, authtypes.FeeCollectorName, distrtypes.ModuleName, feeCollectorCoins)
		//	if err != nil {
		//		panic(err)
		//	}
		//	totalPower := int64(0)
		//	sk.IterateBondedValidatorsByPower(ctx, func(index int64, validator stakingtypes.ValidatorI) (stop bool) {
		//		consPower := validator.GetConsensusPower(sk.PowerReduction(ctx))
		//		totalPower = totalPower + consPower
		//		return false
		//	})
		//	feePool := dk.GetFeePool(ctx)
		//	if totalPower != 0 {
		//		sk.IterateBondedValidatorsByPower(ctx, func(index int64, validator stakingtypes.ValidatorI) (stop bool) {
		//			consPower := validator.GetConsensusPower(sk.PowerReduction(ctx))
		//			powerFraction := sdk.NewDec(consPower).QuoTruncate(sdk.NewDec(totalPower))
		//			rewards := feeCollected.MulDecTruncate(powerFraction)
		//			dk.AllocateTokensToValidator(ctx, validator, rewards)
		//			remaining.Sub(rewards)
		//			return false
		//		})
		//	}
		//	feePool.CommunityPool = feePool.CommunityPool.Add(remaining...)
		//	dk.SetFeePool(ctx, feePool)
		//}
		k.CoverRedelegationPenalty(ctx)

		// END BLOCK
		ctx = ctx.WithBlockTime(lsmEpoch.StartTime.Add(lsmEpoch.Duration))

		staking.EndBlocker(ctx, sk)
		// mimic liquidstaking endblocker except increasing epoch
		{
			k.DistributeReward(ctx)
			k.CoverSlashingAndHandleMatureUnbondings(ctx)
			k.RemoveDeletableRedelegationInfos(ctx)
			k.HandleQueuedLiquidUnstakes(ctx)
			k.HandleUnprocessedQueuedLiquidUnstakes(ctx)
			k.HandleQueuedWithdrawInsuranceRequests(ctx)
			newlyRankedInInsurances, rankOutInsurances := k.RankInsurances(ctx)
			k.RePairRankedInsurances(ctx, newlyRankedInInsurances, rankOutInsurances)
			k.IncrementEpoch(ctx)
		}

		return nil

	}
}
