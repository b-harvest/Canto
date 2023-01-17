package keeper_test

import (
	"math"
	"strconv"

	"github.com/Canto-Network/Canto-Testnet-v2/v1/x/liquidstaking/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	tmrand "github.com/tendermint/tendermint/libs/rand"
)

func (suite *KeeperTestSuite) TestKeeperLiquidStake() {
	const iteration = 1000
	seed := tmrand.Int63n(10000)
	tmrand.Seed(seed)

	suite.Run("liquid staking negative: fails with empty account, seed:"+strconv.FormatInt(seed, 10), func() {
		tokenAmount := generateRandomTokenAmount()
		liquidStaker := generateRandomAccount()
		_, err := suite.keeper.LiquidStake(suite.ctx, liquidStaker, tokenAmount)
		suite.Require().Error(err)
	})

	suite.Run("liquid staking negative: fails with not enough token amount", func() {
		bondDenom := suite.app.StakingKeeper.BondDenom(suite.ctx)
		tokenAmount := generateRandomTokenAmount().AddRaw(100)
		liquidStaker := accountWithCoins(generateRandomAccount(), suite.app.BankKeeper, suite.ctx, sdk.NewCoins(sdk.NewCoin(bondDenom, tokenAmount.SubRaw(50))))
		_, err := suite.keeper.LiquidStake(suite.ctx, liquidStaker, tokenAmount)
		suite.Require().Error(err)
	})

	suite.Run("liquid staking positive: multiple times", func() {
		for i := 0; i < iteration; i++ {
			bondDenom := suite.app.StakingKeeper.BondDenom(suite.ctx)
			tokenAmount := generateRandomTokenAmount()
			liquidStaker := accountWithCoins(generateRandomAccount(), suite.app.BankKeeper, suite.ctx, sdk.NewCoins(sdk.NewCoin(bondDenom, tokenAmount)))
			suite.Require().True(tokenAmount.Equal(suite.app.BankKeeper.GetBalance(suite.ctx, liquidStaker, bondDenom).Amount))

			reqId, err := suite.keeper.LiquidStake(suite.ctx, liquidStaker, tokenAmount)
			suite.Require().NoError(err)
			_, found := suite.keeper.GetChunkBondRequest(suite.ctx, reqId)
			suite.Require().True(found)
			suite.Require().True(sdk.ZeroInt().Equal(suite.app.BankKeeper.GetBalance(suite.ctx, liquidStaker, bondDenom).Amount))
		}
		allReqs := suite.keeper.GetAllChunkBondRequests(suite.ctx)
		suite.Require().Equal(iteration, len(allReqs))
	})
}

func (suite *KeeperTestSuite) TestKeeperLiquidUnstake() {
	const iteration = 1000
	seed := tmrand.Int63n(10000)
	tmrand.Seed(seed)
	chunkSize := suite.keeper.GetParams(suite.ctx).ChunkSize
	liquidBondDenom := suite.keeper.LiquidBondDenom(suite.ctx)

	suite.Run("liquid unstaking negative: fails with empty account, seed:"+strconv.FormatInt(seed, 10), func() {
		liquidUnstaker := generateRandomAccount()
		numChunkUnbond := tmrand.Int63n(math.MaxInt64)
		_, err := suite.keeper.LiquidUnstake(suite.ctx, liquidUnstaker, uint64(numChunkUnbond))
		suite.Require().Error(err)
	})

	suite.Run("liquid unstaking negative: fails with not enough token amount", func() {
		numChunkUnbond := tmrand.Int63n(math.MaxInt64)
		liquidTokenAmount := chunkSize.MulRaw(numChunkUnbond)
		liquidUnstaker := accountWithCoins(generateRandomAccount(), suite.app.BankKeeper, suite.ctx, sdk.NewCoins(sdk.NewCoin(liquidBondDenom, liquidTokenAmount.SubRaw(50))))

		_, err := suite.keeper.LiquidUnstake(suite.ctx, liquidUnstaker, uint64(numChunkUnbond))
		suite.Require().Error(err)
	})

	suite.Run("liquid unstaking positive:  multiple times", func() {
		for i := 0; i < iteration; i++ {
			numChunkUnbond := tmrand.Int63n(math.MaxInt64)
			liquidTokenAmount := chunkSize.MulRaw(numChunkUnbond)
			liquidUnstaker := accountWithCoins(generateRandomAccount(), suite.app.BankKeeper, suite.ctx, sdk.NewCoins(sdk.NewCoin(liquidBondDenom, liquidTokenAmount)))
			reqId, err := suite.keeper.LiquidUnstake(suite.ctx, liquidUnstaker, uint64(numChunkUnbond))
			suite.Require().NoError(err)
			_, found := suite.keeper.GetChunkUnbondRequest(suite.ctx, reqId)
			suite.Require().True(found)
			suite.Require().True(sdk.ZeroInt().Equal(suite.app.BankKeeper.GetBalance(suite.ctx, liquidUnstaker, liquidBondDenom).Amount))
		}
		allReqs := suite.keeper.GetAllChunkUnbondRequests(suite.ctx)
		suite.Require().Equal(iteration, len(allReqs))
	})
}

func (suite *KeeperTestSuite) TestKeeperBidInsurance() {
	const iteration = 1000
	seed := tmrand.Int63n(10000)
	tmrand.Seed(seed)

	bondDenom := suite.app.StakingKeeper.BondDenom(suite.ctx)
	params := suite.keeper.GetParams(suite.ctx)

	//TODO: get correct chunk size in bond denom(not liquid bond denom)
	tokenAmountForEachChunk := params.ChunkSize
	minimumInsurancePercentage := params.MinInsurancePercentage

	getMinLiquidTokenAmount := func() sdk.Int {
		return minimumInsurancePercentage.Add(sdk.NewDec(100)).MulInt(tokenAmountForEachChunk).TruncateInt()
	}

	suite.Run("bid insurance negative: fails with empty token amount, seed:"+strconv.FormatInt(seed, 10), func() {
		tokenAmount := getMinLiquidTokenAmount()
		insurer := generateRandomAccount()
		validator := generateRandomValidatorAccount(suite)
		_, err := suite.keeper.BidInsurance(suite.ctx, insurer, validator, tokenAmount, minimumInsurancePercentage)
		suite.Require().Error(err)
	})

	suite.Run("bid insurance negative: fails with not enough token amount", func() {
		tokenAmount := getMinLiquidTokenAmount().SubRaw(1)
		insurer := accountWithCoins(generateRandomAccount(), suite.app.BankKeeper, suite.ctx, sdk.NewCoins(sdk.NewCoin(bondDenom, tokenAmount.SubRaw(50))))
		validator := generateRandomValidatorAccount(suite)
		_, err := suite.keeper.BidInsurance(suite.ctx, insurer, validator, tokenAmount, minimumInsurancePercentage)
		suite.Require().Error(err)
	})

	suite.Run("bid insurance negative: fails with not enough insurance fee rate", func() {
		tokenAmount := getMinLiquidTokenAmount().SubRaw(1)
		insurer := accountWithCoins(generateRandomAccount(), suite.app.BankKeeper, suite.ctx, sdk.NewCoins(sdk.NewCoin(bondDenom, tokenAmount.AddRaw(50))))
		validator := generateRandomValidatorAccount(suite)
		_, err := suite.keeper.BidInsurance(suite.ctx, insurer, validator, tokenAmount, sdk.ZeroDec())
		suite.Require().Error(err)
	})

	suite.Run("bid insurance positive: multiple times", func() {
		for i := 0; i < iteration; i++ {
			tokenBufferAmount := tmrand.Int63n(100)
			tokenAmount := getMinLiquidTokenAmount().AddRaw(tokenBufferAmount)
			insurer := accountWithCoins(generateRandomAccount(), suite.app.BankKeeper, suite.ctx, sdk.NewCoins(sdk.NewCoin(bondDenom, tokenAmount)))
			validator := generateRandomValidatorAccount(suite)
			reqId, err := suite.keeper.BidInsurance(suite.ctx, insurer, validator, tokenAmount, minimumInsurancePercentage)
			suite.Require().NoError(err)
			_, found := suite.keeper.GetInsuranceBid(suite.ctx, reqId)
			suite.Require().True(found)
		}
		allReqs := suite.keeper.GetAllInsuranceBids(suite.ctx)
		suite.Require().Equal(iteration, len(allReqs))
	})
}

func (suite *KeeperTestSuite) TestKeeperUnbondInsurance() {
	const iteration = 1000
	seed := tmrand.Int63n(10000)
	tmrand.Seed(seed)

	suite.Run("unbond insurance negative: unbonding alive chunk is not found, seed:"+strconv.FormatInt(seed, 10), func() {
		insurer := generateRandomAccount()
		id := generateRandomId()
		_, err := suite.keeper.UnbondInsurance(suite.ctx, insurer, id)
		suite.Require().Error(err)
	})

	suite.Run("unbond insurance multiple times", func() {
		for i := 0; i < iteration; i++ {
			insurer := generateRandomAccount()
			id := generateRandomId()
			suite.keeper.SetAliveChunk(suite.ctx, types.AliveChunk{Id: id})

			reqId, err := suite.keeper.UnbondInsurance(suite.ctx, insurer, id)
			suite.Require().NoError(err)
			_, found := suite.keeper.GetInsuranceUnbondRequest(suite.ctx, reqId)
			suite.Require().True(found)
		}
		allReqs := suite.keeper.GetAllInsuranceUnbondRequests(suite.ctx)
		suite.Require().Equal(iteration, len(allReqs))
	})
}
