package types_test

import (
	"github.com/Canto-Network/Canto/v6/x/liquidstaking/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ethermint "github.com/evmos/ethermint/types"
	"github.com/stretchr/testify/suite"
	"testing"
)

type dynamicFeeRateTestSuite struct {
	suite.Suite
}

func TestDynamicFeeRateTestSuite(t *testing.T) {
	suite.Run(t, new(dynamicFeeRateTestSuite))
}

func (suite *dynamicFeeRateTestSuite) TestCalcFormulaBetweenSoftCapAndOptimal() {
	for _, tc := range []struct {
		name       string
		setupParam func(params *types.DynamicFeeRate)
		u          sdk.Dec
		expected   string
	}{
		{
			"default params and u = 6%(=0.06)",
			func(params *types.DynamicFeeRate) {},
			sdk.MustNewDecFromStr("0.06"),
			"0.025000000000000000",
		},
	} {
		suite.Run(tc.name, func() {
			params := types.DefaultParams().DynamicFeeRate
			tc.setupParam(&params)
			suite.Equal(
				tc.expected,
				types.CalcFormulaBetweenSoftCapAndOptimal(
					params.R0, params.USoftCap, params.UOptimal, params.Slope1, tc.u,
				).String(),
			)
			suite.Equal(
				tc.expected,
				types.CalcDynamicFeeRate(tc.u, params).String(),
			)
		})
	}
}

func (suite *dynamicFeeRateTestSuite) TestCalcFormulaUpperOptimal() {
	for _, tc := range []struct {
		name       string
		setupParam func(params *types.DynamicFeeRate)
		u          sdk.Dec
		expected   string
	}{
		{
			"default params and u = 10%(=0.1)",
			func(params *types.DynamicFeeRate) {},
			sdk.MustNewDecFromStr("0.1"),
			"0.500000000000000000",
		},
	} {
		suite.Run(tc.name, func() {
			params := types.DefaultParams().DynamicFeeRate
			tc.setupParam(&params)
			suite.Equal(
				tc.expected,
				types.CalcFormulaUpperOptimal(
					params.R0, params.UOptimal, params.UHardCap, params.Slope1, params.Slope2, tc.u,
				).String(),
			)
			suite.Equal(
				tc.expected,
				types.CalcDynamicFeeRate(tc.u, params).String(),
			)
		})
	}
}

func (suite *dynamicFeeRateTestSuite) TestGetAvailableChunkSlots() {
	for _, tc := range []struct {
		name           string
		setupParam     func(params *types.DynamicFeeRate)
		u              sdk.Dec
		totalSupplyAmt sdk.Int
		expected       string
	}{
		{
			"(Normal) default params, u = 6%(=0.06), and total supply = 10B",
			func(params *types.DynamicFeeRate) {},
			sdk.MustNewDecFromStr("0.06"),
			sdk.TokensFromConsensusPower(1_000_000_000, ethermint.PowerReduction),
			sdk.NewInt(160).String(),
		},
		{
			"(Normal) default params, u = 9%(=0.09), and total supply = 10.5B",
			func(params *types.DynamicFeeRate) {},
			sdk.MustNewDecFromStr("0.09"),
			sdk.TokensFromConsensusPower(1_050_000_000, ethermint.PowerReduction),
			sdk.NewInt(42).String(),
		},
		{
			"(Abnormal) hardcap = 5%(=0.05), u = 6%(=0.06), and total supply = 10B",
			func(params *types.DynamicFeeRate) {
				params.UHardCap = sdk.MustNewDecFromStr("0.05")
			},
			sdk.MustNewDecFromStr("0.06"),
			sdk.TokensFromConsensusPower(1_000_000_000, ethermint.PowerReduction),
			sdk.ZeroInt().String(),
		},
	} {
		suite.Run(tc.name, func() {
			params := types.DefaultParams().DynamicFeeRate
			tc.setupParam(&params)
			suite.Equal(
				tc.expected,
				types.GetAvailableChunkSlots(tc.u, params.UHardCap, tc.totalSupplyAmt).String(),
			)
		})
	}
}
