package simulation_test

import (
	cantoapp "github.com/Canto-Network/Canto/v6/app/params"
	"github.com/Canto-Network/Canto/v6/x/liquidstaking/simulation"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	"math/rand"
	"testing"
	"time"
)

func TestProposalContents(t *testing.T) {
	app, ctx := createTestApp(false)

	s := rand.NewSource(1)
	r := rand.New(s)

	accounts := getTestingAccounts(t, r, app, ctx, 10)

	getTestingValidator0(t, app, ctx, accounts)
	getTestingValidator1(t, app, ctx, accounts)

	// begin a new block
	blockTime := time.Now().UTC()
	app.BeginBlock(abci.RequestBeginBlock{
		Header: tmproto.Header{
			Height:  app.LastBlockHeight() + 1,
			AppHash: app.LastCommitID().Hash,
			Time:    blockTime,
		},
	})
	app.EndBlock(abci.RequestEndBlock{Height: app.LastBlockHeight() + 1})

	// change type of app.BankKeeper to BaseKeeper
	baseKeeper, ok := app.BankKeeper.(bankkeeper.BaseKeeper)
	require.True(t, ok)

	// execute ProposalContents function
	weightedProposalContent := simulation.ProposalContents(
		app.LiquidStakingKeeper,
		app.AccountKeeper,
		baseKeeper,
		app.StakingKeeper,
		app.DistrKeeper,
		app.InflationKeeper,
	)
	require.Len(t, weightedProposalContent, 3)

	w0 := weightedProposalContent[0]
	w1 := weightedProposalContent[1]
	w2 := weightedProposalContent[2]

	// tests w0 interface:
	require.Equal(t, simulation.OpWeightSimulateUpdateDynamicFeeRateProposal, w0.AppParamsKey())
	require.Equal(t, cantoapp.DefaultWeightUpdateDynamicFeeRateProposal, w0.DefaultWeight())

	// tests w1 interface:
	require.Equal(t, simulation.OpWeightSimulateUpdateMaximumDiscountRate, w1.AppParamsKey())
	require.Equal(t, cantoapp.DefaultWeightUpdateMaximumDiscountRate, w1.DefaultWeight())

	// tests w2 interface:
	require.Equal(t, simulation.OpWeightSimulateAdvanceEpoch, w2.AppParamsKey())
	require.Equal(t, cantoapp.DefaultWeightAdvanceEpoch, w2.DefaultWeight())

	content0 := w0.ContentSimulatorFn()(r, ctx, accounts)
	require.Nil(t, content0)

	content1 := w1.ContentSimulatorFn()(r, ctx, accounts)
	require.Nil(t, content1)

	content2 := w2.ContentSimulatorFn()(r, ctx, accounts)
	require.Nil(t, content2)
}
