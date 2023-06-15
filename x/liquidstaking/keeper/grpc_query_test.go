package keeper_test

import (
	"github.com/Canto-Network/Canto/v6/x/liquidstaking/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (suite *KeeperTestSuite) TestGRPCParams() {
	resp, err := suite.app.LiquidStakingKeeper.Params(sdk.WrapSDKContext(suite.ctx), &types.QueryParamsRequest{})
	suite.Require().NoError(err)
	suite.Require().Equal(suite.app.LiquidStakingKeeper.GetParams(suite.ctx), resp.Params)
}

func (suite *KeeperTestSuite) TestGRPCEpoch() {
	resp, err := suite.app.LiquidStakingKeeper.Epoch(sdk.WrapSDKContext(suite.ctx), &types.QueryEpochRequest{})
	suite.Require().NoError(err)
	suite.Require().Equal(suite.app.LiquidStakingKeeper.GetEpoch(suite.ctx), resp.Epoch)
}

func (suite *KeeperTestSuite) TestGRPCChunks() {
	suite.setupLiquidStakeTestingEnv(testingEnvOptions{
		desc:                  "",
		numVals:               3,
		fixedValFeeRate:       tenPercentFeeRate,
		valFeeRates:           nil,
		fixedPower:            1,
		powers:                nil,
		numInsurances:         3,
		fixedInsuranceFeeRate: tenPercentFeeRate,
		insuranceFeeRates:     nil,
		numPairedChunks:       3,
		fundingAccountBalance: types.ChunkSize.MulRaw(1000),
	})

	for _, tc := range []struct {
		name      string
		req       *types.QueryChunksRequest
		expectErr bool
		postRun   func(response *types.QueryChunksResponse)
	}{
		{
			"nil request",
			nil,
			true,
			nil,
		},
		{
			"query all",
			&types.QueryChunksRequest{},
			false,
			func(response *types.QueryChunksResponse) {
				suite.Require().Len(response.Chunks, 3)
			},
		},
		{
			"query only paired chunks",
			&types.QueryChunksRequest{
				Status: types.CHUNK_STATUS_PAIRED,
			},
			false,
			func(response *types.QueryChunksResponse) {
				suite.Require().Len(response.Chunks, 3)
			},
		},
		{
			"query only pairing chunks",
			&types.QueryChunksRequest{
				Status: types.CHUNK_STATUS_PAIRING,
			},
			false,
			func(response *types.QueryChunksResponse) {
				suite.Require().Len(response.Chunks, 0)
			},
		},
	} {
		suite.Run(tc.name, func() {
			resp, err := suite.app.LiquidStakingKeeper.Chunks(sdk.WrapSDKContext(suite.ctx), tc.req)
			if tc.expectErr {
				suite.Require().Error(err)
				return
			}
			suite.Require().NoError(err)
			if tc.postRun != nil {
				tc.postRun(resp)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestGRPCChunk() {
	env := suite.setupLiquidStakeTestingEnv(testingEnvOptions{
		desc:                  "",
		numVals:               3,
		fixedValFeeRate:       tenPercentFeeRate,
		valFeeRates:           nil,
		fixedPower:            1,
		powers:                nil,
		numInsurances:         3,
		fixedInsuranceFeeRate: tenPercentFeeRate,
		insuranceFeeRates:     nil,
		numPairedChunks:       3,
		fundingAccountBalance: types.ChunkSize.MulRaw(1000),
	})

	for _, tc := range []struct {
		name      string
		req       *types.QueryChunkRequest
		expectErr bool
		postRun   func(response *types.QueryChunkResponse)
	}{
		{
			"nil request",
			nil,
			true,
			nil,
		},
		{
			"invalid request",
			&types.QueryChunkRequest{},
			true,
			nil,
		},
		{
			"query by id",
			&types.QueryChunkRequest{
				Id: 1,
			},
			false,
			func(response *types.QueryChunkResponse) {
				chunk := env.pairedChunks[0]
				suite.True(chunk.Equal(response.Chunk))
			},
		},
	} {
		suite.Run(tc.name, func() {
			resp, err := suite.app.LiquidStakingKeeper.Chunk(sdk.WrapSDKContext(suite.ctx), tc.req)
			if tc.expectErr {
				suite.Require().Error(err)
				return
			}
			suite.Require().NoError(err)
			if tc.postRun != nil {
				tc.postRun(resp)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestGRPCInsurances() {
	suite.setupLiquidStakeTestingEnv(testingEnvOptions{
		desc:                  "",
		numVals:               3,
		fixedValFeeRate:       tenPercentFeeRate,
		valFeeRates:           nil,
		fixedPower:            1,
		powers:                nil,
		numInsurances:         5,
		fixedInsuranceFeeRate: tenPercentFeeRate,
		insuranceFeeRates:     nil,
		numPairedChunks:       3,
		fundingAccountBalance: types.ChunkSize.MulRaw(1000),
	})

	for _, tc := range []struct {
		name      string
		req       *types.QueryInsurancesRequest
		expectErr bool
		postRun   func(response *types.QueryInsurancesResponse)
	}{
		{
			"nil request",
			nil,
			true,
			nil,
		},
		{
			"query all",
			&types.QueryInsurancesRequest{},
			false,
			func(response *types.QueryInsurancesResponse) {
				suite.Require().Len(response.Insurances, 5)
			},
		},
		{
			"query only paired chunks",
			&types.QueryInsurancesRequest{
				Status: types.INSURANCE_STATUS_PAIRED,
			},
			false,
			func(response *types.QueryInsurancesResponse) {
				suite.Require().Len(response.Insurances, 3)
			},
		},
		{
			"query only pairing chunks",
			&types.QueryInsurancesRequest{
				Status: types.INSURANCE_STATUS_PAIRING,
			},
			false,
			func(response *types.QueryInsurancesResponse) {
				suite.Require().Len(response.Insurances, 2)
			},
		},
	} {
		suite.Run(tc.name, func() {
			resp, err := suite.app.LiquidStakingKeeper.Insurances(sdk.WrapSDKContext(suite.ctx), tc.req)
			if tc.expectErr {
				suite.Require().Error(err)
				return
			}
			suite.Require().NoError(err)
			if tc.postRun != nil {
				tc.postRun(resp)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestGRPCInsurance() {
	env := suite.setupLiquidStakeTestingEnv(testingEnvOptions{
		desc:                  "",
		numVals:               3,
		fixedValFeeRate:       tenPercentFeeRate,
		valFeeRates:           nil,
		fixedPower:            1,
		powers:                nil,
		numInsurances:         3,
		fixedInsuranceFeeRate: tenPercentFeeRate,
		insuranceFeeRates:     nil,
		numPairedChunks:       3,
		fundingAccountBalance: types.ChunkSize.MulRaw(1000),
	})

	for _, tc := range []struct {
		name      string
		req       *types.QueryInsuranceRequest
		expectErr bool
		postRun   func(response *types.QueryInsuranceResponse)
	}{
		{
			"nil request",
			nil,
			true,
			nil,
		},
		{
			"invalid request",
			&types.QueryInsuranceRequest{},
			true,
			nil,
		},
		{
			"query by id",
			&types.QueryInsuranceRequest{
				Id: 1,
			},
			false,
			func(response *types.QueryInsuranceResponse) {
				suite.True(env.insurances[0].Equal(response.Insurance))
			},
		},
	} {
		suite.Run(tc.name, func() {
			resp, err := suite.app.LiquidStakingKeeper.Insurance(sdk.WrapSDKContext(suite.ctx), tc.req)
			if tc.expectErr {
				suite.Require().Error(err)
				return
			}
			suite.Require().NoError(err)
			if tc.postRun != nil {
				tc.postRun(resp)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestGRPCWithdrawInsuranceRequests() {
	env := suite.setupLiquidStakeTestingEnv(testingEnvOptions{
		desc:                  "",
		numVals:               3,
		fixedValFeeRate:       tenPercentFeeRate,
		valFeeRates:           nil,
		fixedPower:            1,
		powers:                nil,
		numInsurances:         5,
		fixedInsuranceFeeRate: tenPercentFeeRate,
		insuranceFeeRates:     nil,
		numPairedChunks:       3,
		fundingAccountBalance: types.ChunkSize.MulRaw(1000),
	})
	// 3 providers requests withdraw.
	// 3 withdraw insurance requests will be queued.
	for i := 0; i < 3; i++ {
		suite.app.LiquidStakingKeeper.DoWithdrawInsurance(
			suite.ctx,
			types.NewMsgWithdrawInsurance(
				env.providers[i].String(),
				env.insurances[i].Id,
			),
		)
	}
	for _, tc := range []struct {
		name      string
		req       *types.QueryWithdrawInsuranceRequestsRequest
		expectErr bool
		postRun   func(response *types.QueryWithdrawInsuranceRequestsResponse)
	}{
		{
			"nil request",
			nil,
			true,
			nil,
		},
		{
			"query all",
			&types.QueryWithdrawInsuranceRequestsRequest{},
			false,
			func(response *types.QueryWithdrawInsuranceRequestsResponse) {
				// Only paired or unpairing insurances can have withdraw requests.
				suite.Require().Len(response.WithdrawInsuranceRequests, 3)
			},
		},
	} {
		suite.Run(tc.name, func() {
			resp, err := suite.app.LiquidStakingKeeper.WithdrawInsuranceRequests(sdk.WrapSDKContext(suite.ctx), tc.req)
			if tc.expectErr {
				suite.Require().Error(err)
				return
			}
			suite.Require().NoError(err)
			if tc.postRun != nil {
				tc.postRun(resp)
			}
		})
	}

}
