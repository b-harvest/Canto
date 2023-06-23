package types_test

import (
	"github.com/Canto-Network/Canto/v6/x/liquidstaking/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"
	"testing"
)

type netAmountTestSuite struct {
	suite.Suite
}

func TestNetAmountTestSuite(t *testing.T) {
	suite.Run(t, new(netAmountTestSuite))
}

func (suite *netAmountTestSuite) TestCalcNetAmount() {
	nas := types.NetAmountState{
		TotalChunksBalance:          sdk.ZeroInt(),
		TotalLiquidTokens:           sdk.MustNewDecFromStr("250000000000000000000000").TruncateInt(),
		TotalUnbondingChunksBalance: sdk.MustNewDecFromStr("250000000000000000000000").TruncateInt(),
		TotalRemainingRewards:       sdk.MustNewDecFromStr("160000000000000000000"),
	}
	suite.Equal(
		"500320000000000000000000.000000000000000000",
		nas.CalcNetAmount(sdk.MustNewDecFromStr("160000000000000000000").TruncateInt()).String(),
	)
}

func (suite *netAmountTestSuite) TestCalcMintRate() {
	nas := types.NetAmountState{
		LsTokensTotalSupply: sdk.MustNewDecFromStr("500000000000000000000000").TruncateInt(),
		NetAmount:           sdk.MustNewDecFromStr("503320000000000000000000"),
	}
	suite.Equal(
		"0.993403798776126519",
		nas.CalcMintRate().String(),
	)

	nas.NetAmount = sdk.ZeroDec()
	suite.Equal(
		"0.000000000000000000",
		nas.CalcMintRate().String(),
	)
}
