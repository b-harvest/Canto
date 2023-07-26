package simulation

import (
	"fmt"
	"math/rand"

	"github.com/Canto-Network/Canto/v6/app/params"
	inflationkeeper "github.com/Canto-Network/Canto/v6/x/inflation/keeper"
	inflationtypes "github.com/Canto-Network/Canto/v6/x/inflation/types"
	"github.com/Canto-Network/Canto/v6/x/liquidstaking/keeper"
	"github.com/Canto-Network/Canto/v6/x/liquidstaking/types"
	"github.com/armon/go-metrics"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/cosmos/cosmos-sdk/x/simulation"
	"github.com/cosmos/cosmos-sdk/x/staking"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	ethermint "github.com/evmos/ethermint/types"
)

const (
	OpWeightSimulateUpdateDynamicFeeRateProposal = "op_weight_simulate_update_dynamic_fee_rate_proposal"
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
			OpWeightSimulateAdvanceEpoch,
			params.DefaultWeightAdvanceEpoch,
			SimulateAdvanceEpoch(k, ak, bk, sk, dk, ik),
		),
	}
}

// SimulateUpdateDynamicFeeRateProposal generates random update dynamic fee rate param change proposal content.
func SimulateUpdateDynamicFeeRateProposal(k keeper.Keeper) simtypes.ContentSimulatorFn {
	return func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) simtypes.Content {
		feeRate := genDynamicFeeRate(r)
		k.SetParams(ctx, types.Params{DynamicFeeRate: feeRate})
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
		bondDenom := sk.BondDenom(ctx)
		lsmEpoch := k.GetEpoch(ctx)
		ctx = ctx.WithBlockTime(lsmEpoch.StartTime.Add(lsmEpoch.Duration))
		staking.BeginBlocker(ctx, sk)

		// mimic the begin block logic of epoch module
		// currently epoch module use hooks when begin block and inflation module
		// implemented that hook, so actual logic is in inflation module.
		{
			epochMintProvision, found := ik.GetEpochMintProvision(ctx)
			if !found {
				panic("epoch mint provision not found")
			}
			// mintedCoin := sdk.NewCoin(inflationParams.MintDenom, epochMintProvision.TruncateInt())
			mintedCoin := sdk.NewCoin(sdk.DefaultBondDenom, sdk.TokensFromConsensusPower(100, ethermint.PowerReduction))
			staking, communityPool, err := ik.MintAndAllocateInflation(ctx, mintedCoin)
			if err != nil {
				panic(err)
			}
			defer func() {
				if mintedCoin.Amount.IsInt64() {
					telemetry.IncrCounterWithLabels(
						[]string{"inflation", "allocate", "total"},
						float32(mintedCoin.Amount.Int64()),
						[]metrics.Label{telemetry.NewLabel("denom", mintedCoin.Denom)},
					)
				}
				if staking.AmountOf(mintedCoin.Denom).IsInt64() {
					telemetry.IncrCounterWithLabels(
						[]string{"inflation", "allocate", "staking", "total"},
						float32(staking.AmountOf(mintedCoin.Denom).Int64()),
						[]metrics.Label{telemetry.NewLabel("denom", mintedCoin.Denom)},
					)
				}
				if communityPool.AmountOf(mintedCoin.Denom).IsInt64() {
					telemetry.IncrCounterWithLabels(
						[]string{"inflation", "allocate", "community_pool", "total"},
						float32(communityPool.AmountOf(mintedCoin.Denom).Int64()),
						[]metrics.Label{telemetry.NewLabel("denom", mintedCoin.Denom)},
					)
				}
			}()

			ctx.EventManager().EmitEvent(
				sdk.NewEvent(
					inflationtypes.EventTypeMint,
					sdk.NewAttribute(inflationtypes.AttributeEpochNumber, fmt.Sprintf("%d", -1)),
					sdk.NewAttribute(inflationtypes.AttributeKeyEpochProvisions, epochMintProvision.String()),
					sdk.NewAttribute(sdk.AttributeKeyAmount, mintedCoin.Amount.String()),
				),
			)
		}

		feeCollector := ak.GetModuleAccount(ctx, authtypes.FeeCollectorName)
		// mimic the begin block logic of distribution module
		{
			feeCollectorBalance := bk.SpendableCoins(ctx, feeCollector.GetAddress())
			rewardsToBeDistributed := feeCollectorBalance.AmountOf(bondDenom)

			// mimic distribution.BeginBlock (AllocateTokens, get rewards from feeCollector, AllocateTokensToValidator, add remaining to feePool)
			err := bk.SendCoinsFromModuleToModule(ctx, authtypes.FeeCollectorName, distrtypes.ModuleName, feeCollectorBalance)
			if err != nil {
				panic(err)
			}
			totalRewards := sdk.ZeroDec()
			totalPower := int64(0)
			sk.IterateBondedValidatorsByPower(ctx, func(index int64, validator stakingtypes.ValidatorI) (stop bool) {
				consPower := validator.GetConsensusPower(sk.PowerReduction(ctx))
				totalPower = totalPower + consPower
				return false
			})
			if totalPower != 0 {
				sk.IterateBondedValidatorsByPower(ctx, func(index int64, validator stakingtypes.ValidatorI) (stop bool) {
					consPower := validator.GetConsensusPower(sk.PowerReduction(ctx))
					powerFraction := sdk.NewDec(consPower).QuoTruncate(sdk.NewDec(totalPower))
					reward := rewardsToBeDistributed.ToDec().MulTruncate(powerFraction)
					dk.AllocateTokensToValidator(ctx, validator, sdk.DecCoins{{Denom: bondDenom, Amount: reward}})
					totalRewards = totalRewards.Add(reward)
					return false
				})
			}
			remaining := rewardsToBeDistributed.ToDec().Sub(totalRewards)
			feePool := dk.GetFeePool(ctx)
			feePool.CommunityPool = feePool.CommunityPool.Add(sdk.DecCoins{
				{Denom: bondDenom, Amount: remaining}}...)
			dk.SetFeePool(ctx, feePool)
		}
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
