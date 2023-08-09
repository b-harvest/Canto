package types_test

import (
	"testing"

	"github.com/Canto-Network/Canto/v6/x/liquidstaking/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"
)

type insuranceStateTestSuite struct {
	suite.Suite
}

func TestInsuranceStateTestSuite(t *testing.T) {
	suite.Run(t, new(insuranceStateTestSuite))
}

func (suite *insuranceStateTestSuite) TestEqual() {
	is := types.InsuranceState{
		TotalInsuranceTokens:               sdk.ZeroInt(),
		TotalPairedInsuranceTokens:         sdk.ZeroInt(),
		TotalUnpairingInsuranceTokens:      sdk.ZeroInt(),
		TotalRemainingInsuranceCommissions: sdk.ZeroDec(),
	}
	cpy := is
	suite.True(is.Equal(cpy))

	cpy = is
	cpy.TotalInsuranceTokens = is.TotalInsuranceTokens.Add(sdk.OneInt())
	suite.False(
		is.Equal(cpy),
		"total insurance tokens should affect equality",
	)

	cpy = is
	cpy.TotalPairedInsuranceTokens = is.TotalPairedInsuranceTokens.Add(sdk.OneInt())
	suite.False(
		is.Equal(cpy),
		"total paired insurance tokens should affect equality",
	)

	cpy = is
	cpy.TotalUnpairingInsuranceTokens = is.TotalUnpairingInsuranceTokens.Add(sdk.OneInt())
	suite.False(
		is.Equal(cpy),
		"total unpairing insurance tokens should affect equality",
	)

	cpy = is
	cpy.TotalRemainingInsuranceCommissions = is.TotalRemainingInsuranceCommissions.Add(sdk.OneDec())
	suite.False(
		is.Equal(cpy),
		"total remaining insurance commissions should affect equality",
	)

	cpy = is
}

func (suite *insuranceStateTestSuite) TestIsZeroState() {
	is := types.InsuranceState{
		TotalInsuranceTokens:               sdk.ZeroInt(),
		TotalPairedInsuranceTokens:         sdk.ZeroInt(),
		TotalUnpairingInsuranceTokens:      sdk.ZeroInt(),
		TotalRemainingInsuranceCommissions: sdk.ZeroDec(),
	}
	suite.True(is.IsZeroState())

	cpy := is

	cpy = is
	cpy.TotalInsuranceTokens = is.TotalInsuranceTokens.Add(sdk.OneInt())
	suite.True(
		cpy.IsZeroState(),
		"total insurance tokens should not affect zero state",
	)
}

func (suite *insuranceStateTestSuite) TestString() {
	is := types.InsuranceState{
		TotalInsuranceTokens:               sdk.NewInt(1),
		TotalPairedInsuranceTokens:         sdk.NewInt(1),
		TotalUnpairingInsuranceTokens:      sdk.NewInt(1),
		TotalRemainingInsuranceCommissions: sdk.NewDec(1),
	}
	suite.Equal(
		`InsuranceState:
	TotalInsuranceTokens:       1
	TotalPairedInsuranceTokens: 1
    TotalUnpairingInsuranceTokens: 1
    TotalRemainingInsuranceCommissions: 1.000000000000000000
`,
		is.String(),
	)
}
