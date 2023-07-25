package keeper_test

import (
	"github.com/Canto-Network/Canto/v6/x/liquidstaking/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (suite *KeeperTestSuite) TestGetNetAmountState() {
	env := suite.setupLiquidStakeTestingEnv(testingEnvOptions{
		desc:                  "",
		numVals:               3,
		fixedValFeeRate:       TenPercentFeeRate,
		valFeeRates:           nil,
		fixedPower:            1,
		powers:                nil,
		numInsurances:         3,
		fixedInsuranceFeeRate: TenPercentFeeRate,
		insuranceFeeRates:     nil,
		numPairedChunks:       3,
		fundingAccountBalance: types.ChunkSize.MulRaw(40),
	})

	suite.ctx = suite.advanceHeight(suite.ctx, 100, "delegation rewards are accumulated")
	nas := suite.app.LiquidStakingKeeper.GetNetAmountState(suite.ctx)

	cachedCtx, _ := suite.ctx.CacheContext()
	suite.app.DistrKeeper.WithdrawDelegationRewards(cachedCtx, env.pairedChunks[0].DerivedAddress(), env.insurances[0].GetValidator())
	delReward := suite.app.BankKeeper.GetBalance(cachedCtx, env.pairedChunks[0].DerivedAddress(), suite.denom)
	totalDelReward := delReward.Amount.MulRaw(int64(len(env.pairedChunks)))
	suite.Equal("8999964000143999250000", totalDelReward.String())

	// Calc TotalRemainingRewards manually
	rest := totalDelReward.ToDec().Mul(sdk.OneDec().Sub(TenPercentFeeRate))
	remaining := rest.Mul(sdk.OneDec().Sub(nas.FeeRate))
	result := remaining.Mul(sdk.OneDec().Sub(nas.DiscountRate))
	suite.Equal("7595237306532985264014.570227833838100000", result.String())
	suite.Equal(result.String(), nas.TotalRemainingRewards.String())
	suite.True(
		totalDelReward.GT(nas.TotalRemainingRewards.TruncateInt()),
		"total del reward should be greater than total remaining rewards",
	)
}
