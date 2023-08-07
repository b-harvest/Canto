package simulation

import (
	"encoding/json"
	"fmt"
	"math/rand"

	"github.com/Canto-Network/Canto/v6/x/liquidstaking/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
)

// Simulation parameter constants
const (
	dynamicFeeRate = "dynamic_fee_rate"
)

func genDynamicFeeRate(r *rand.Rand) types.DynamicFeeRate {
	//maxFeeRate := types.RandomDec(r, sdk.MustNewDecFromStr("0.4"), sdk.MustNewDecFromStr("0.5"))
	//
	//r0 := types.RandomDec(r, sdk.ZeroDec(), sdk.MustNewDecFromStr("0.01"))
	//slope1 := types.RandomDec(r, sdk.ZeroDec(), sdk.MustNewDecFromStr("0.3"))
	//slope2 := maxFeeRate.Sub(slope1).Sub(r0)
	//
	//uSoftCap := types.RandomDec(r, sdk.MustNewDecFromStr("0.05"), sdk.MustNewDecFromStr("0.06"))
	//uOptimal := types.RandomDec(r, sdk.MustNewDecFromStr("0.07"), sdk.MustNewDecFromStr("0.09"))
	//uHardCap := types.RandomDec(r, sdk.MustNewDecFromStr("0.1"), types.SecurityCap)

	return types.DynamicFeeRate{
		R0:         types.DefaultR0,
		USoftCap:   types.DefaultUSoftCap,
		UHardCap:   types.DefaultUHardCap,
		UOptimal:   types.DefaultUOptimal,
		Slope1:     types.DefaultSlope1,
		Slope2:     types.DefaultSlope2,
		MaxFeeRate: types.DefaultMaxFee,
	}
}

func genMaximumDiscountRate(r *rand.Rand) sdk.Dec {
	return types.RandomDec(r, sdk.ZeroDec(), sdk.MustNewDecFromStr("0.09"))
}

func RandomizedGenState(simState *module.SimulationState) {
	genesis := types.DefaultGenesisState()
	genesis.Epoch.StartTime = simState.GenTimestamp

	simState.AppParams.GetOrGenerate(
		simState.Cdc, dynamicFeeRate, &genesis.Params.DynamicFeeRate, simState.Rand,
		func(r *rand.Rand) { genesis.Params.DynamicFeeRate = genDynamicFeeRate(r) },
	)

	bz, _ := json.MarshalIndent(&genesis, "", " ")
	fmt.Printf("Selected randomly generated liquidstaking parameters:\n%s\n", bz)
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(genesis)
}
