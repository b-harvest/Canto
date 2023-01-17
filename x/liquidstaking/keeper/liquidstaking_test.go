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

	expected := make(map[string]types.ChunkBondRequest)
	suite.Run("liquid staking positive: multiple times", func() {
		for i := 0; i < iteration; i++ {
			bondDenom := suite.app.StakingKeeper.BondDenom(suite.ctx)
			tokenAmount := generateRandomTokenAmount()
			liquidStaker := accountWithCoins(generateRandomAccount(), suite.app.BankKeeper, suite.ctx, sdk.NewCoins(sdk.NewCoin(bondDenom, tokenAmount)))
			suite.Require().True(tokenAmount.Equal(suite.app.BankKeeper.GetBalance(suite.ctx, liquidStaker, bondDenom).Amount))

			reqId, err := suite.keeper.LiquidStake(suite.ctx, liquidStaker, tokenAmount)

			suite.Require().NoError(err)
			req, found := suite.keeper.GetChunkBondRequest(suite.ctx, reqId)
			expected[liquidStaker.String()] = req
			suite.Require().True(found)
			suite.Require().True(sdk.ZeroInt().Equal(suite.app.BankKeeper.GetBalance(suite.ctx, liquidStaker, bondDenom).Amount))
		}
		allReqs := suite.keeper.GetAllChunkBondRequests(suite.ctx)
		suite.Require().Equal(len(expected), len(allReqs))
	})

	suite.Run("cancel liquid stake positive", func() {
		for addr, req := range expected {
			liquidStaker, err := sdk.AccAddressFromBech32(addr)
			suite.Require().NoError(err)
			_, err = suite.keeper.CancelLiquidStaking(suite.ctx, liquidStaker, req.Id)
			suite.Require().NoError(err)
			_, found := suite.keeper.GetChunkBondRequest(suite.ctx, req.Id)
			suite.Require().False(found)
		}
		allReqs := suite.keeper.GetAllChunkBondRequests(suite.ctx)
		suite.Require().Equal(0, len(allReqs))
	})

	suite.Run("cancel liquid stake negative: address mismatch", func() {
		bondDenom := suite.app.StakingKeeper.BondDenom(suite.ctx)
		tokenAmount := generateRandomTokenAmount()
		liquidStaker := accountWithCoins(generateRandomAccount(), suite.app.BankKeeper, suite.ctx, sdk.NewCoins(sdk.NewCoin(bondDenom, tokenAmount)))
		wrongStaker := sdk.AccAddress{}
		for i := 0; i < 10; i++ {
			wrongStaker = generateRandomAccount()
			if !liquidStaker.Equals(wrongStaker) {
				break
			}
		}
		suite.Require().False(liquidStaker.Equals(wrongStaker))
		reqId, err := suite.keeper.LiquidStake(suite.ctx, liquidStaker, tokenAmount)
		suite.Require().NoError(err)
		_, err = suite.keeper.CancelLiquidStaking(suite.ctx, wrongStaker, reqId)
		suite.Require().Error(err)
		_, found := suite.keeper.GetChunkBondRequest(suite.ctx, reqId)
		suite.Require().True(found)
		_, err = suite.keeper.CancelLiquidStaking(suite.ctx, liquidStaker, reqId)
		suite.Require().NoError(err)
		_, found = suite.keeper.GetChunkBondRequest(suite.ctx, reqId)
		suite.Require().False(found)
	})
}

func (suite *KeeperTestSuite) TestKeeperLiquidUnstake() {
	const iteration = 1000
	seed := tmrand.Int63n(10000)
	tmrand.Seed(seed)
	chunkSize := suite.keeper.GetParams(suite.ctx).ChunkSize
	liquidBondDenom := suite.keeper.LiquidBondDenom(suite.ctx)
	expected := make(map[string]types.ChunkUnbondRequest)

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
			req, found := suite.keeper.GetChunkUnbondRequest(suite.ctx, reqId)
			expected[liquidUnstaker.String()] = req
			suite.Require().True(found)
			suite.Require().True(sdk.ZeroInt().Equal(suite.app.BankKeeper.GetBalance(suite.ctx, liquidUnstaker, liquidBondDenom).Amount))
		}
		allReqs := suite.keeper.GetAllChunkUnbondRequests(suite.ctx)
		suite.Require().Equal(len(expected), len(allReqs))
	})

	suite.Run("cancel liquid unstake positive", func() {
		for addr, req := range expected {
			liquidStaker, err := sdk.AccAddressFromBech32(addr)
			suite.Require().NoError(err)
			_, err = suite.keeper.CancelLiquidUnstaking(suite.ctx, liquidStaker, req.Id)
			suite.Require().NoError(err)
			_, found := suite.keeper.GetChunkUnbondRequest(suite.ctx, req.Id)
			suite.Require().False(found)
		}
		allReqs := suite.keeper.GetAllChunkUnbondRequests(suite.ctx)
		suite.Require().Equal(0, len(allReqs))
	})

	suite.Run("cancel liquid unstake negative: address mismatch", func() {
		numChunkUnbond := tmrand.Int63n(math.MaxInt64)
		liquidTokenAmount := chunkSize.MulRaw(numChunkUnbond)
		liquidUnstaker := accountWithCoins(generateRandomAccount(), suite.app.BankKeeper, suite.ctx, sdk.NewCoins(sdk.NewCoin(liquidBondDenom, liquidTokenAmount)))
		wrongStaker := sdk.AccAddress{}
		for i := 0; i < 10; i++ {
			wrongStaker = generateRandomAccount()
			if !liquidUnstaker.Equals(wrongStaker) {
				break
			}
		}
		suite.Require().False(liquidUnstaker.Equals(wrongStaker))
		reqId, err := suite.keeper.LiquidUnstake(suite.ctx, liquidUnstaker, uint64(numChunkUnbond))
		suite.Require().NoError(err)
		_, err = suite.keeper.CancelLiquidUnstaking(suite.ctx, wrongStaker, reqId)
		suite.Require().Error(err)
		_, found := suite.keeper.GetChunkUnbondRequest(suite.ctx, reqId)
		suite.Require().True(found)
		_, err = suite.keeper.CancelLiquidUnstaking(suite.ctx, liquidUnstaker, reqId)
		suite.Require().NoError(err)
		_, found = suite.keeper.GetChunkUnbondRequest(suite.ctx, reqId)
		suite.Require().False(found)
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

	expected := make(map[string]types.InsuranceBid)
	suite.Run("bid insurance positive: multiple times", func() {
		for i := 0; i < iteration; i++ {
			tokenBufferAmount := tmrand.Int63n(100)
			tokenAmount := getMinLiquidTokenAmount().AddRaw(tokenBufferAmount)
			insurer := accountWithCoins(generateRandomAccount(), suite.app.BankKeeper, suite.ctx, sdk.NewCoins(sdk.NewCoin(bondDenom, tokenAmount)))
			validator := generateRandomValidatorAccount(suite)
			reqId, err := suite.keeper.BidInsurance(suite.ctx, insurer, validator, tokenAmount, minimumInsurancePercentage)
			suite.Require().NoError(err)
			req, found := suite.keeper.GetInsuranceBid(suite.ctx, reqId)
			expected[insurer.String()] = req
			suite.Require().True(found)
		}
		allReqs := suite.keeper.GetAllInsuranceBids(suite.ctx)
		suite.Require().Equal(len(expected), len(allReqs))
	})

	suite.Run("cancel bid insurance positive", func() {
		for addr, req := range expected {
			insurer, err := sdk.AccAddressFromBech32(addr)
			suite.Require().NoError(err)
			_, err = suite.keeper.CancelInsuranceBid(suite.ctx, insurer, req.Id)
			suite.Require().NoError(err)
			_, found := suite.keeper.GetInsuranceBid(suite.ctx, req.Id)
			suite.Require().False(found)
		}
		allReqs := suite.keeper.GetAllInsuranceBids(suite.ctx)
		suite.Require().Equal(0, len(allReqs))
	})

	suite.Run("cancel bid insurance negative: address mismatch", func() {
		tokenBufferAmount := tmrand.Int63n(100)
		tokenAmount := getMinLiquidTokenAmount().AddRaw(tokenBufferAmount)
		insurer := accountWithCoins(generateRandomAccount(), suite.app.BankKeeper, suite.ctx, sdk.NewCoins(sdk.NewCoin(bondDenom, tokenAmount)))
		validator := generateRandomValidatorAccount(suite)
		wrongInsurer := sdk.AccAddress{}
		for i := 0; i < 10; i++ {
			wrongInsurer = generateRandomAccount()
			if !insurer.Equals(wrongInsurer) {
				break
			}
		}
		suite.Require().False(insurer.Equals(wrongInsurer))
		reqId, err := suite.keeper.BidInsurance(suite.ctx, insurer, validator, tokenAmount, minimumInsurancePercentage)
		suite.Require().NoError(err)
		_, err = suite.keeper.CancelInsuranceBid(suite.ctx, wrongInsurer, reqId)
		suite.Require().Error(err)
		_, found := suite.keeper.GetInsuranceBid(suite.ctx, reqId)
		suite.Require().True(found)
		_, err = suite.keeper.CancelInsuranceBid(suite.ctx, insurer, reqId)
		suite.Require().NoError(err)
		_, found = suite.keeper.GetInsuranceBid(suite.ctx, reqId)
		suite.Require().False(found)
	})
}

func (suite *KeeperTestSuite) TestKeeperUnbondInsurance() {
	const iteration = 1000
	seed := tmrand.Int63n(10000)
	tmrand.Seed(seed)
	used := make(map[uint64]bool)

	suite.Run("unbond insurance negative: unbonding alive chunk is not found, seed:"+strconv.FormatInt(seed, 10), func() {
		insurer := generateRandomAccount()
		id := generateUniqueId(suite, used)
		used[id] = true
		_, err := suite.keeper.UnbondInsurance(suite.ctx, insurer, id)
		suite.Require().Error(err)
	})

	expected := make(map[string]types.InsuranceUnbondRequest)
	suite.Run("unbond insurance multiple times", func() {
		for i := 0; i < iteration; i++ {
			insurer := generateRandomAccount()
			id := generateUniqueId(suite, used)
			used[id] = true
			suite.keeper.SetAliveChunk(suite.ctx, types.AliveChunk{Id: id})

			reqId, err := suite.keeper.UnbondInsurance(suite.ctx, insurer, id)
			suite.Require().NoError(err)
			req, found := suite.keeper.GetInsuranceUnbondRequest(suite.ctx, reqId)
			expected[insurer.String()] = req
			suite.Require().True(found)
		}
		allReqs := suite.keeper.GetAllInsuranceUnbondRequests(suite.ctx)
		suite.Require().Equal(len(expected), len(allReqs))
	})

	suite.Run("cancel insurance unbond request positive", func() {
		for addr, req := range expected {
			insurer, err := sdk.AccAddressFromBech32(addr)
			suite.Require().NoError(err)
			_, err = suite.keeper.CancelInsuranceUnbond(suite.ctx, insurer, req.AliveChunkId)
			suite.Require().NoError(err)
			_, found := suite.keeper.GetInsuranceUnbondRequest(suite.ctx, req.AliveChunkId)

			suite.Require().False(found)
		}
		allReqs := suite.keeper.GetAllInsuranceUnbondRequests(suite.ctx)
		suite.Require().Equal(0, len(allReqs))
	})

	suite.Run("cancel insurance unbond negative: address mismatch", func() {
		insurer := generateRandomAccount()

		id := generateUniqueId(suite, used)
		used[id] = true
		suite.keeper.SetAliveChunk(suite.ctx, types.AliveChunk{Id: id})
		reqId, err := suite.keeper.UnbondInsurance(suite.ctx, insurer, id)
		suite.Require().NoError(err)
		_, found := suite.keeper.GetInsuranceUnbondRequest(suite.ctx, reqId)
		suite.Require().True(found)

		wrongInsurer := sdk.AccAddress{}
		for i := 0; i < 10; i++ {
			wrongInsurer = generateRandomAccount()
			if !insurer.Equals(wrongInsurer) {
				break
			}
		}
		suite.Require().False(insurer.Equals(wrongInsurer))
		_, err = suite.keeper.CancelInsuranceUnbond(suite.ctx, wrongInsurer, reqId)
		suite.Require().Error(err)
		_, found = suite.keeper.GetInsuranceUnbondRequest(suite.ctx, reqId)
		suite.Require().True(found)
		_, err = suite.keeper.CancelInsuranceUnbond(suite.ctx, insurer, reqId)
		suite.Require().NoError(err)
		_, found = suite.keeper.GetInsuranceUnbondRequest(suite.ctx, reqId)
		suite.Require().False(found)
	})
}
