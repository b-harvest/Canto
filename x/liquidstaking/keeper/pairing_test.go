package keeper_test

import (
	"strconv"

	"github.com/Canto-Network/Canto-Testnet-v2/v1/x/liquidstaking/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	tmrand "github.com/tendermint/tendermint/libs/rand"
)

type StateBuildRecord struct {
	AliveChunk             map[uint64]types.AliveChunk
	ChunkBondRequest       map[uint64]types.ChunkBondRequest
	ChunkUnbondRequest     map[uint64]types.ChunkUnbondRequest
	InsuranceBid           map[uint64]types.InsuranceBid
	InsuranceUnbondRequest map[uint64]types.InsuranceUnbondRequest
}

func newStateBuildRecord() StateBuildRecord {
	ret := StateBuildRecord{
		AliveChunk:             make(map[uint64]types.AliveChunk),
		ChunkBondRequest:       make(map[uint64]types.ChunkBondRequest),
		ChunkUnbondRequest:     make(map[uint64]types.ChunkUnbondRequest),
		InsuranceBid:           make(map[uint64]types.InsuranceBid),
		InsuranceUnbondRequest: make(map[uint64]types.InsuranceUnbondRequest),
	}
	return ret
}

type NumState struct {
	insuranceBid           int
	insuranceUnbondRequest int
	chunkBondRequest       int
	chunkUnbondRequest     int
	aliveChunks            int
	insuranceUnbonded      int
	chunkUnbonded          int
}

func (suite *KeeperTestSuite) generateRandomNumState() (ret NumState) {
	params := suite.keeper.GetParams(suite.ctx)
	maxAliveChunks := int(params.MaxAliveChunk.Int64())
	ret.aliveChunks = tmrand.Intn(maxAliveChunks + 1)
	ret.chunkBondRequest = tmrand.Intn(30)
	ret.chunkUnbondRequest = func() int {
		if (ret.aliveChunks + ret.chunkBondRequest) == 0 {
			return 0
		}
		return tmrand.Intn(ret.aliveChunks + ret.chunkBondRequest)
	}()
	ret.insuranceBid = tmrand.Intn(30)
	ret.insuranceUnbondRequest = func() int {
		if ret.aliveChunks == 0 {
			return 0
		}
		return tmrand.Intn(ret.aliveChunks)
	}()
	return
}

func getStateFromKeeper(suite *KeeperTestSuite) types.State {
	return types.State{
		InsuranceBids:           suite.keeper.GetAllInsuranceBids(suite.ctx),
		InsuranceUnbondRequests: suite.keeper.GetAllInsuranceUnbondRequests(suite.ctx),
		ChunkBondRequests:       suite.keeper.GetAllChunkBondRequests(suite.ctx),
		ChunkUnbondRequests:     suite.keeper.GetAllChunkUnbondRequests(suite.ctx),
		AliveChunks:             suite.keeper.GetAllAliveChunks(suite.ctx),
	}
}

var (
	valAddrList = []string{}
)

func (suite *KeeperTestSuite) randomizeCurrentState(
	numState NumState,
) (ret types.State, record StateBuildRecord) {
	params := suite.keeper.GetParams(suite.ctx)
	maxAliveChunks := params.MaxAliveChunk
	suite.Require().True(maxAliveChunks.GTE(sdk.NewInt(int64(numState.aliveChunks))))
	suite.Require().LessOrEqual(numState.insuranceUnbondRequest, numState.aliveChunks)
	record = newStateBuildRecord()
	bondDenom := suite.app.StakingKeeper.BondDenom(suite.ctx)
	liquidStakingDenom := suite.keeper.LiquidStakingDenom(suite.ctx)
	chunkSize := params.ChunkSize

	for i := 0; i < numState.aliveChunks; i++ {
		id := generateUniqueId(suite, record.AliveChunk)
		chunkBondRequest := generateRandomChunkBondRequestWithTokenAmount(generateUniqueId(suite, record.ChunkBondRequest), nativeTokenChunkSize)
		insuranceBid := generateRandomInsuranceBid(suite, generateUniqueId(suite, record.InsuranceBid))

		aliveChunk := types.NewAliveChunk(id, chunkBondRequest, insuranceBid)
		record.AliveChunk[id] = aliveChunk
		suite.keeper.SetAliveChunk(suite.ctx, aliveChunk)

		depositCoinsIntoModule(
			suite,
			sdk.NewCoin(bondDenom, aliveChunk.TokenAmount.Add(aliveChunk.InsuranceAmount)),
		)
		valAddr, err := sdk.ValAddressFromBech32(aliveChunk.ValidatorAddress)
		suite.Require().NoError(err)
		val, found := suite.app.StakingKeeper.GetValidator(suite.ctx, valAddr)
		suite.Require().True(found)
		valAddrStr := valAddr.String()
		suite.Require().Equal(val.OperatorAddress, valAddrStr)

		err = suite.keeper.DelegateTokenAmount(suite.ctx, valAddrStr, aliveChunk.TokenAmount)
		suite.Require().NoError(err)

		valAddrList = append(valAddrList, valAddrStr)
	}
	for i := 0; i < numState.chunkBondRequest; i++ {
		id := generateUniqueId(suite, record.ChunkBondRequest)
		record.ChunkBondRequest[id] = generateRandomChunkBondRequestWithTokenAmount(id, nativeTokenChunkSize)
		suite.keeper.SetChunkBondRequest(suite.ctx, record.ChunkBondRequest[id])

		err := suite.app.BankKeeper.MintCoins(suite.ctx, types.ModuleName,
			sdk.NewCoins(sdk.NewCoin(bondDenom, record.ChunkBondRequest[id].TokenAmount)))
		suite.Require().NoError(err)
		depositCoinsIntoModule(suite, sdk.NewCoin(bondDenom, record.ChunkBondRequest[id].TokenAmount))
	}
	for i := 0; i < numState.chunkUnbondRequest; {
		tmp := tmrand.Intn(numState.chunkUnbondRequest-i) + 1

		id := generateUniqueId(suite, record.ChunkUnbondRequest)
		chunkUnbondRequest := generateRandomChunkUnbondRequest(id)
		chunkUnbondRequest.NumChunkUnbond = uint64(tmp)
		suite.keeper.SetChunkUnbondRequest(suite.ctx, chunkUnbondRequest)
		record.ChunkUnbondRequest[id] = chunkUnbondRequest
		i += tmp

		depositCoinsIntoModule(suite, sdk.NewCoin(liquidStakingDenom, chunkSize.MulRaw(int64(tmp))))
	}
	for i := 0; i < numState.insuranceBid; i++ {
		id := generateUniqueId(suite, record.InsuranceBid)
		insuranceBid := generateRandomInsuranceBid(suite, id)
		record.InsuranceBid[id] = insuranceBid
		suite.keeper.SetInsuranceBid(suite.ctx, insuranceBid)
		suite.keeper.SetInsuranceBidIndex(suite.ctx, insuranceBid)

		depositCoinsIntoModule(suite, sdk.NewCoin(bondDenom, insuranceBid.InsuranceAmount))
	}

	// hack to ensure last counter is greater than existing ids
	suite.keeper.SetLastAliveChunkId(suite.ctx, 50000)
	suite.keeper.SetLastChunkBondRequestId(suite.ctx, 50000)
	suite.keeper.SetLastChunkUnbondRequestId(suite.ctx, 50000)
	suite.keeper.SetLastInsuranceBidId(suite.ctx, 50000)
	suite.keeper.SetLastUnbondingChunkId(suite.ctx, 50000)
	used := make(map[uint64]bool)
	for i := 0; i < numState.insuranceUnbondRequest; i++ {
		insuranceUnbondRequest := func() types.InsuranceUnbondRequest {
			for _, aliveChunk := range record.AliveChunk {
				if _, found := used[aliveChunk.Id]; !found {
					used[aliveChunk.Id] = true
					return types.InsuranceUnbondRequest{
						InsuranceProviderAddress: aliveChunk.InsuranceProviderAddress,
						AliveChunkId:             aliveChunk.Id,
					}
				}
			}
			suite.Require().Fail("fail to generate insurance unbond request")
			return types.InsuranceUnbondRequest{}
		}()
		record.InsuranceUnbondRequest[insuranceUnbondRequest.AliveChunkId] = insuranceUnbondRequest
		suite.keeper.SetInsuranceUnbondRequest(suite.ctx, insuranceUnbondRequest)
		suite.keeper.SetInsuranceUnbondRequestIndex(suite.ctx, insuranceUnbondRequest)
	}
	ret = getStateFromKeeper(suite)
	return
}

func getExpectedNumResolveUnbondingQueues(currentNumState NumState) (ret NumState) {
	ret = currentNumState

	// handle insurance Unbond Request
	ret.aliveChunks -= currentNumState.insuranceUnbondRequest
	ret.insuranceUnbonded += currentNumState.insuranceUnbondRequest
	ret.insuranceUnbondRequest = 0

	// compensate Chunk Bond/Unbond request
	numCanceledChunkRequest := getMin(ret.chunkBondRequest, ret.chunkUnbondRequest)
	ret.chunkBondRequest -= numCanceledChunkRequest
	ret.chunkUnbondRequest -= numCanceledChunkRequest

	// handle Chunk Unbond request
	ret.chunkUnbonded = getMin(ret.chunkUnbondRequest, ret.aliveChunks)
	ret.aliveChunks -= ret.chunkUnbonded
	ret.chunkUnbondRequest -= ret.chunkUnbonded

	ret.insuranceUnbonded -= ret.chunkUnbondRequest
	ret.chunkUnbonded += ret.chunkUnbondRequest
	ret.chunkUnbondRequest = 0
	return
}

func (suite *KeeperTestSuite) TestKeeperResolveUnbondingQueues() {
	seed := tmrand.Int63n(10000)
	tmrand.Seed(seed)

	currentNumState := suite.generateRandomNumState()
	currentState, _ := suite.randomizeCurrentState(currentNumState)

	suite.Require().Equal(currentNumState.aliveChunks, len(currentState.AliveChunks))
	suite.Require().Equal(currentNumState.chunkBondRequest, len(currentState.ChunkBondRequests))
	suite.Require().Equal(currentNumState.chunkUnbondRequest, func() (ret int) {
		for _, req := range currentState.ChunkUnbondRequests {
			ret += int(req.NumChunkUnbond)
		}
		return
	}())
	suite.Require().Equal(currentNumState.insuranceBid, len(currentState.InsuranceBids))
	suite.Require().Equal(currentNumState.insuranceUnbondRequest, len(currentState.InsuranceUnbondRequests))

	suite.Run("simple check seed: "+strconv.FormatInt(seed, 10), func() {
		expectedNumState := getExpectedNumResolveUnbondingQueues(currentNumState)
		newState := suite.keeper.ResolveUnbondingQueues(suite.ctx, currentState)
		suite.Require().NotNil(newState)

		suite.Equal(expectedNumState.aliveChunks, len(newState.AliveChunks))
		suite.Equal(expectedNumState.chunkBondRequest, len(newState.ChunkBondRequests))
		suite.Equal(expectedNumState.chunkUnbondRequest, func() (ret int) {
			for _, req := range newState.ChunkUnbondRequests {
				ret += int(req.NumChunkUnbond)
			}
			return
		}())
		suite.Equal(expectedNumState.insuranceBid, len(newState.InsuranceBids))
		suite.Equal(expectedNumState.insuranceUnbondRequest, len(newState.InsuranceUnbondRequests))
		suite.Equal(expectedNumState.insuranceUnbonded, len(newState.InsuranceUnbonded))
		suite.Equal(expectedNumState.chunkUnbonded, len(newState.ChunkUnbonded))
	})
}

func getExpectedNumUnpairPair(currentNumState NumState, maxAliveChunk int) (ret NumState) {
	ret = currentNumState

	// pair chunk bond request with insurance bid
	remainingMaximumAliveChunks := maxAliveChunk - ret.aliveChunks
	remainingAvailableChunks := ret.chunkBondRequest + ret.insuranceUnbonded
	newAliveChunks := getMin(remainingAvailableChunks, ret.insuranceBid, remainingMaximumAliveChunks)
	reused := getMin(ret.insuranceUnbonded, newAliveChunks)
	ret.insuranceUnbonded -= reused
	ret.chunkBondRequest -= (newAliveChunks - reused)
	ret.insuranceBid -= newAliveChunks
	ret.aliveChunks += newAliveChunks
	return
}

func getExpectedNumPairing(currentNumState NumState, maxAliveChunk int) (ret NumState) {
	ret = getExpectedNumResolveUnbondingQueues(currentNumState)
	ret = getExpectedNumUnpairPair(ret, maxAliveChunk)
	return
}

func (suite *KeeperTestSuite) TestKeeperPairChunkAndInsurance() {
	seed := tmrand.Int63n(10000)
	tmrand.Seed(seed)

	currentNumState := suite.generateRandomNumState()
	currentState, _ := suite.randomizeCurrentState(currentNumState)

	suite.Require().Equal(currentNumState.aliveChunks, len(currentState.AliveChunks))
	suite.Require().Equal(currentNumState.chunkBondRequest, len(currentState.ChunkBondRequests))
	suite.Require().Equal(currentNumState.chunkUnbondRequest, func() (ret int) {
		for _, req := range currentState.ChunkUnbondRequests {
			ret += int(req.NumChunkUnbond)
		}
		return
	}())
	suite.Require().Equal(currentNumState.insuranceBid, len(currentState.InsuranceBids))
	suite.Require().Equal(currentNumState.insuranceUnbondRequest, len(currentState.InsuranceUnbondRequests))

	suite.Run("simple check seed: "+strconv.FormatInt(seed, 10), func() {
		params := suite.keeper.GetParams(suite.ctx)
		maxAliveChunk := int(params.MaxAliveChunk.Int64())
		expectedNumState := getExpectedNumPairing(currentNumState, maxAliveChunk)
		err := suite.keeper.PairChunkAndInsurance(suite.ctx)
		suite.Require().NoError(err)

		newState := getStateFromKeeper(suite)

		suite.Equal(expectedNumState.aliveChunks, len(newState.AliveChunks))
		suite.Equal(expectedNumState.chunkBondRequest, len(newState.ChunkBondRequests))
		suite.Equal(expectedNumState.chunkUnbondRequest, func() (ret int) {
			for _, req := range newState.ChunkUnbondRequests {
				ret += int(req.NumChunkUnbond)
			}
			return
		}())
		suite.Equal(expectedNumState.insuranceBid, len(newState.InsuranceBids))
		suite.Equal(expectedNumState.insuranceUnbondRequest, len(newState.InsuranceUnbondRequests))
	})
}
