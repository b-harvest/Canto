package simulation

import (
	"github.com/Canto-Network/Canto/v6/app/params"
	"github.com/Canto-Network/Canto/v6/x/liquidstaking"
	"github.com/Canto-Network/Canto/v6/x/liquidstaking/keeper"
	"github.com/Canto-Network/Canto/v6/x/liquidstaking/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"
	"math/rand"
)

const (
	OpWeightSimulateUpdateDynamicFeeRateProposal = "op_weight_simulate_update_dynamic_fee_rate_proposal"
	OpWeightSimulateAdvanceEpoch                 = "op_weight_simulate_advance_epoch"
)

func ProposalContents(k keeper.Keeper, am liquidstaking.AppModule) []simtypes.WeightedProposalContent {
	return []simtypes.WeightedProposalContent{
		simulation.NewWeightedProposalContent(
			OpWeightSimulateUpdateDynamicFeeRateProposal,
			params.DefaultWeightUpdateDynamicFeeRateProposal,
			SimulateUpdateDynamicFeeRateProposal(k),
		),
		simulation.NewWeightedProposalContent(
			OpWeightSimulateAdvanceEpoch,
			params.DefaultWeightAdvanceEpoch,
			SimulateAdvanceEpoch(am),
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

func SimulateAdvanceEpoch(am liquidstaking.AppModule) simtypes.ContentSimulatorFn {
	return func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) simtypes.Content {
		am.AdvanceEpochBeginBlock(ctx)
		am.AdvanceEpochEndBlock(ctx)
		return nil
	}
}
