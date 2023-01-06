package keeper_test

import (
	"strconv"

	"github.com/Canto-Network/Canto-Testnet-v2/v1/x/liquidstaking/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	tmrand "github.com/tendermint/tendermint/libs/rand"
)

type (
	FuncSetId = func(sdk.Context, uint64)
	FuncGetId = func(sdk.Context) uint64
)

var (
	nativeTokenChunkSize = sdk.NewInt(500)
	liquidTokenChunkSize = sdk.NewInt(500)
)

func (suite *KeeperTestSuite) testStoreIdInitialized(funcGetId FuncGetId) {
	expected := uint64(0)
	ret := funcGetId(suite.ctx)
	suite.Require().Equal(expected, ret)
}

func (suite *KeeperTestSuite) testStoreIdIncrement(funcSetId FuncSetId, funcGetId FuncGetId, start, end, step uint64) {
	suite.testStoreIdInitialized(funcGetId)
	for expected := start; expected < end; expected += step {
		funcSetId(suite.ctx, expected)
		ret := funcGetId(suite.ctx)
		suite.Require().Equal(expected, ret)
	}
}

func (suite *KeeperTestSuite) TestStoreLastAliveChunkId() {
	seed := tmrand.Int63n(10000)
	tmrand.Seed(seed)
	suite.Run("Initialized Value must be 0", func() {
		suite.testStoreIdInitialized(suite.keeper.GetLastAliveChunkId)
	})
	suite.SetupTest()
	suite.Run("Increment by one", func() {
		suite.testStoreIdIncrement(suite.keeper.SetLastAliveChunkId, suite.keeper.GetLastAliveChunkId, 0, 50, 1)
	})
	suite.SetupTest()
	suite.Run("Increment by random: "+strconv.FormatInt(seed, 10), func() {
		suite.testStoreIdIncrement(suite.keeper.SetLastAliveChunkId, suite.keeper.GetLastAliveChunkId, 0, 100, uint64(tmrand.Intn(4)+1))
	})
}

func (suite *KeeperTestSuite) TestStoreLastUnbondingChunkId() {
	seed := tmrand.Int63n(10000)
	tmrand.Seed(seed)
	suite.Run("Initialized Value must be 0", func() {
		suite.testStoreIdInitialized(suite.keeper.GetLastUnbondingChunkId)
	})
	suite.SetupTest()
	suite.Run("Increment by one", func() {
		suite.testStoreIdIncrement(suite.keeper.SetLastUnbondingChunkId, suite.keeper.GetLastUnbondingChunkId, 0, 50, 1)
	})
	suite.SetupTest()
	suite.Run("Increment by random: "+strconv.FormatInt(seed, 10), func() {
		suite.testStoreIdIncrement(suite.keeper.SetLastUnbondingChunkId, suite.keeper.GetLastUnbondingChunkId, 0, 100, uint64(tmrand.Intn(4)+1))
	})
}

func (suite *KeeperTestSuite) TestStoreLastInsuranceBidId() {
	seed := tmrand.Int63n(10000)
	tmrand.Seed(seed)
	suite.Run("Initialized Value must be 0", func() {
		suite.testStoreIdInitialized(suite.keeper.GetLastInsuranceBidId)
	})
	suite.SetupTest()
	suite.Run("Increment by one", func() {
		suite.testStoreIdIncrement(suite.keeper.SetLastInsuranceBidId, suite.keeper.GetLastInsuranceBidId, 0, 50, 1)
	})
	suite.SetupTest()
	suite.Run("Increment by random: "+strconv.FormatInt(seed, 10), func() {
		suite.testStoreIdIncrement(suite.keeper.SetLastInsuranceBidId, suite.keeper.GetLastInsuranceBidId, 0, 100, uint64(tmrand.Intn(4)+1))
	})
}

func (suite *KeeperTestSuite) TestStoreLastChunkBondRequestId() {
	seed := tmrand.Int63n(10000)
	tmrand.Seed(seed)
	suite.Run("Initialized Value must be 0", func() {
		suite.testStoreIdInitialized(suite.keeper.GetLastChunkBondRequestId)
	})
	suite.SetupTest()
	suite.Run("Increment by one", func() {
		suite.testStoreIdIncrement(suite.keeper.SetLastChunkBondRequestId, suite.keeper.GetLastChunkBondRequestId, 0, 50, 1)
	})
	suite.SetupTest()
	suite.Run("Increment by random: "+strconv.FormatInt(seed, 10), func() {
		suite.testStoreIdIncrement(suite.keeper.SetLastChunkBondRequestId, suite.keeper.GetLastChunkBondRequestId, 0, 100, uint64(tmrand.Intn(4)+1))
	})
}

func (suite *KeeperTestSuite) TestStoreLastChunkUnbondRequestId() {
	seed := tmrand.Int63n(10000)
	tmrand.Seed(seed)
	suite.Run("Initialized Value must be 0", func() {
		suite.testStoreIdInitialized(suite.keeper.GetLastChunkUnbondRequestId)
	})
	suite.SetupTest()
	suite.Run("Increment by one", func() {
		suite.testStoreIdIncrement(suite.keeper.SetLastChunkUnbondRequestId, suite.keeper.GetLastChunkUnbondRequestId, 0, 50, 1)
	})
	suite.SetupTest()
	suite.Run("Increment by random: "+strconv.FormatInt(seed, 10), func() {
		suite.testStoreIdIncrement(suite.keeper.SetLastChunkUnbondRequestId, suite.keeper.GetLastChunkUnbondRequestId, 0, 100, uint64(tmrand.Intn(4)+1))
	})
}

func (suite *KeeperTestSuite) TestStoreAliveChunk() {
	seed := tmrand.Int63n(10000)
	tmrand.Seed(seed)
	_, found := suite.keeper.GetAliveChunk(suite.ctx, types.AliveChunkId(0))
	suite.Require().Equal(false, found)
	chunkMap := make(map[types.AliveChunkId]types.AliveChunk)
	const numIteration = 100

	suite.Run("Get/Set alive chunk with random seed: "+strconv.FormatInt(seed, 10), func() {
		count := 0
		for i := 0; i < numIteration; i++ {
			id := types.AliveChunkId(tmrand.Uint64())
			if _, has := chunkMap[id]; has {
				count++
				suite.Require().Less(count, 10)
				continue
			}

			expected := generateRandomAliveChunk(suite, id)
			chunkMap[id] = expected
			suite.keeper.SetAliveChunk(suite.ctx, expected)
			ret, found := suite.keeper.GetAliveChunk(suite.ctx, id)
			suite.Require().Equal(true, found)
			suite.Require().Equal(expected, ret)
		}
	})
	suite.Run("Delete Alive chunk with random keys", func() {
		count := 0
		for i := 0; i < numIteration; i++ {
			id := types.AliveChunkId(tmrand.Uint64())
			if _, has := chunkMap[id]; has {
				count++
				suite.Require().Less(count, 10)
				continue
			}
			_, found := suite.keeper.GetAliveChunk(suite.ctx, id)
			suite.Require().Equal(false, found)
			suite.keeper.DeleteAliveChunk(suite.ctx, id)
			_, found = suite.keeper.GetAliveChunk(suite.ctx, id)
			suite.Require().Equal(false, found)
		}
	})
	suite.Run("Delete Alive chunk with existing keys", func() {
		count := 0
		for _, expected := range chunkMap {
			id := expected.Id
			_, found := suite.keeper.GetAliveChunk(suite.ctx, id)
			suite.Require().Equal(true, found)
			suite.keeper.DeleteAliveChunk(suite.ctx, id)
			_, found = suite.keeper.GetAliveChunk(suite.ctx, id)
			suite.Require().Equal(false, found)
			count++
		}
		suite.Require().Equal(len(chunkMap), count)
	})
}

func (suite *KeeperTestSuite) TestStoreGetAllAliveChunks() {
	seed := tmrand.Int63n(10000)
	tmrand.Seed(seed)

	expectedMap := make(map[uint64]types.AliveChunk)

	for i := 0; i < 1000; i++ {
		id := generateRandomId()
		expected := generateRandomAliveChunk(suite, id)
		suite.keeper.SetAliveChunk(suite.ctx, expected)
		expectedMap[id] = expected
	}
	aliveChunks := suite.keeper.GetAllAliveChunks(suite.ctx)
	suite.Require().Equal(len(expectedMap), len(aliveChunks))
	for _, actual := range aliveChunks {
		expected, found := expectedMap[actual.Id]
		suite.Require().True(found)
		suite.Require().Equal(expected, actual)
	}
}

func (suite *KeeperTestSuite) TestStoreUnbondingChunk() {
	seed := tmrand.Int63n(10000)
	tmrand.Seed(seed)
	_, found := suite.keeper.GetUnbondingChunk(suite.ctx, types.UnbondingChunkId(0))
	suite.Require().Equal(false, found)
	chunkMap := make(map[types.UnbondingChunkId]types.UnbondingChunk)
	const numIteration = 100

	suite.Run("Get/Set alive chunk with random seed: "+strconv.FormatInt(seed, 10), func() {
		count := 0
		for i := 0; i < numIteration; i++ {
			id := generateRandomId()
			if _, has := chunkMap[id]; has {
				count++
				suite.Require().Less(count, 10)
				continue
			}

			expected := generateRandomUnbondingChunk(suite, id)
			chunkMap[id] = expected
			suite.keeper.SetUnbondingChunk(suite.ctx, expected)
			ret, found := suite.keeper.GetUnbondingChunk(suite.ctx, id)
			suite.Require().Equal(true, found)
			suite.Require().Equal(expected, ret)
		}
	})
	suite.Run("Delete Unbonding chunk with random keys", func() {
		count := 0
		for i := 0; i < numIteration; i++ {
			id := generateRandomId()
			if _, has := chunkMap[id]; has {
				count++
				suite.Require().Less(count, 10)
				continue
			}
			_, found := suite.keeper.GetUnbondingChunk(suite.ctx, id)
			suite.Require().Equal(false, found)
			suite.keeper.DeleteUnbondingChunk(suite.ctx, id)
			_, found = suite.keeper.GetUnbondingChunk(suite.ctx, id)
			suite.Require().Equal(false, found)
		}
	})
	suite.Run("Delete Unbonding chunk with existing keys", func() {
		count := 0
		for _, expected := range chunkMap {
			id := expected.Id
			_, found := suite.keeper.GetUnbondingChunk(suite.ctx, id)
			suite.Require().Equal(true, found)
			suite.keeper.DeleteUnbondingChunk(suite.ctx, id)
			_, found = suite.keeper.GetUnbondingChunk(suite.ctx, id)
			suite.Require().Equal(false, found)
			count++
		}
		suite.Require().Equal(len(chunkMap), count)
	})
}

func (suite *KeeperTestSuite) TestStoreChunkBondRequest() {
	seed := tmrand.Int63n(10000)
	tmrand.Seed(seed)
	_, found := suite.keeper.GetChunkBondRequest(suite.ctx, types.ChunkBondRequestId(0))
	suite.Require().Equal(false, found)
	chunkMap := make(map[types.ChunkBondRequestId]types.ChunkBondRequest)
	const numIteration = 100

	suite.Run("Get/Set alive chunk with random seed: "+strconv.FormatInt(seed, 10), func() {
		count := 0
		for i := 0; i < numIteration; i++ {
			id := generateRandomId()
			if _, has := chunkMap[id]; has {
				count++
				suite.Require().Less(count, 10)
				continue
			}

			expected := generateRandomChunkBondRequest(id)
			chunkMap[id] = expected
			suite.keeper.SetChunkBondRequest(suite.ctx, expected)
			ret, found := suite.keeper.GetChunkBondRequest(suite.ctx, id)
			suite.Require().Equal(true, found)
			suite.Require().Equal(expected, ret)
		}
	})
	suite.Run("Delete Alive chunk with random keys", func() {
		count := 0
		for i := 0; i < numIteration; i++ {
			id := generateRandomId()
			if _, has := chunkMap[id]; has {
				count++
				suite.Require().Less(count, 10)
				continue
			}
			_, found := suite.keeper.GetChunkBondRequest(suite.ctx, id)
			suite.Require().Equal(false, found)
			suite.keeper.DeleteChunkBondRequest(suite.ctx, id)
			_, found = suite.keeper.GetChunkBondRequest(suite.ctx, id)
			suite.Require().Equal(false, found)
		}
	})
	suite.Run("Delete Alive chunk with existing keys", func() {
		count := 0
		for _, expected := range chunkMap {
			id := expected.Id
			_, found := suite.keeper.GetChunkBondRequest(suite.ctx, id)
			suite.Require().Equal(true, found)
			suite.keeper.DeleteChunkBondRequest(suite.ctx, id)
			_, found = suite.keeper.GetChunkBondRequest(suite.ctx, id)
			suite.Require().Equal(false, found)
			count++
		}
		suite.Require().Equal(len(chunkMap), count)
	})
}

func (suite *KeeperTestSuite) TestStoreGetAllChunkBondRequests() {
	seed := tmrand.Int63n(10000)
	tmrand.Seed(seed)

	expectedMap := make(map[uint64]types.ChunkBondRequest)

	for i := 0; i < 1000; i++ {
		id := generateRandomId()
		expected := generateRandomChunkBondRequest(id)
		suite.keeper.SetChunkBondRequest(suite.ctx, expected)
		expectedMap[id] = expected
	}
	reqs := suite.keeper.GetAllChunkBondRequests(suite.ctx)
	suite.Require().Equal(len(expectedMap), len(reqs))
	for _, actual := range reqs {
		expected, found := expectedMap[actual.Id]
		suite.Require().True(found)
		suite.Require().Equal(expected, actual)
	}
}

func (suite *KeeperTestSuite) TestStoreChunkUnbondRequest() {
	seed := tmrand.Int63n(10000)
	tmrand.Seed(seed)
	_, found := suite.keeper.GetChunkUnbondRequest(suite.ctx, types.ChunkUnbondRequestId(0))
	suite.Require().Equal(false, found)
	chunkMap := make(map[types.ChunkUnbondRequestId]types.ChunkUnbondRequest)
	const numIteration = 100

	suite.Run("Get/Set alive chunk with random seed: "+strconv.FormatInt(seed, 10), func() {
		count := 0
		for i := 0; i < numIteration; i++ {
			id := generateRandomId()
			if _, has := chunkMap[id]; has {
				count++
				suite.Require().Less(count, 10)
				continue
			}

			expected := generateRandomChunkUnbondRequest(id)
			chunkMap[id] = expected
			suite.keeper.SetChunkUnbondRequest(suite.ctx, expected)
			ret, found := suite.keeper.GetChunkUnbondRequest(suite.ctx, id)
			suite.Require().Equal(true, found)
			suite.Require().Equal(expected, ret)
		}
	})
	suite.Run("Delete Alive chunk with random keys", func() {
		count := 0
		for i := 0; i < numIteration; i++ {
			id := generateRandomId()
			if _, has := chunkMap[id]; has {
				count++
				suite.Require().Less(count, 10)
				continue
			}
			_, found := suite.keeper.GetChunkUnbondRequest(suite.ctx, id)
			suite.Require().Equal(false, found)
			suite.keeper.DeleteChunkUnbondRequest(suite.ctx, id)
			_, found = suite.keeper.GetChunkUnbondRequest(suite.ctx, id)
			suite.Require().Equal(false, found)
		}
	})
	suite.Run("Delete Alive chunk with existing keys", func() {
		count := 0
		for _, expected := range chunkMap {
			id := expected.Id
			_, found := suite.keeper.GetChunkUnbondRequest(suite.ctx, id)
			suite.Require().Equal(true, found)
			suite.keeper.DeleteChunkUnbondRequest(suite.ctx, id)
			_, found = suite.keeper.GetChunkUnbondRequest(suite.ctx, id)
			suite.Require().Equal(false, found)
			count++
		}
		suite.Require().Equal(len(chunkMap), count)
	})
}

func (suite *KeeperTestSuite) TestStoreGetAllChunkUnbondRequests() {
	seed := tmrand.Int63n(10000)
	tmrand.Seed(seed)

	expectedMap := make(map[uint64]types.ChunkUnbondRequest)

	for i := 0; i < 1000; i++ {
		id := generateRandomId()
		expected := generateRandomChunkUnbondRequest(id)
		suite.keeper.SetChunkUnbondRequest(suite.ctx, expected)
		expectedMap[id] = expected
	}
	reqs := suite.keeper.GetAllChunkUnbondRequests(suite.ctx)
	suite.Require().Equal(len(expectedMap), len(reqs))
	for _, actual := range reqs {
		expected, found := expectedMap[actual.Id]
		suite.Require().True(found)
		suite.Require().Equal(expected, actual)
	}
}

func (suite *KeeperTestSuite) TestStoreInsuranceBid() {
	seed := tmrand.Int63n(10000)
	tmrand.Seed(seed)
	_, found := suite.keeper.GetInsuranceBid(suite.ctx, types.InsuranceBidId(0))
	suite.Require().Equal(false, found)
	expectedMap := make(map[types.InsuranceBidId]types.InsuranceBid)
	const numIteration = 100

	suite.Run("Get/Set insurance with random seed: "+strconv.FormatInt(seed, 10), func() {
		count := 0
		for i := 0; i < numIteration; i++ {
			id := generateRandomId()
			if _, has := expectedMap[id]; has {
				count++
				suite.Require().Less(count, 10)
				continue
			}

			expected := generateRandomInsuranceBid(suite, id)
			expectedMap[id] = expected
			suite.keeper.SetInsuranceBid(suite.ctx, expected)
			ret, found := suite.keeper.GetInsuranceBid(suite.ctx, id)
			suite.Require().Equal(true, found)
			suite.Require().Equal(expected, ret)
		}
	})
	suite.Run("Delete insurance with random seed: "+strconv.FormatInt(seed, 10), func() {
		count := 0
		for i := 0; i < numIteration; i++ {
			id := generateRandomId()
			if _, has := expectedMap[id]; has {
				count++
				suite.Require().Less(count, 10)
				continue
			}
			_, found := suite.keeper.GetInsuranceBid(suite.ctx, id)
			suite.Require().Equal(false, found)
			suite.keeper.DeleteInsuranceBid(suite.ctx, id)
			_, found = suite.keeper.GetInsuranceBid(suite.ctx, id)
			suite.Require().Equal(false, found)
		}
	})
	suite.Run("Delete insurance with existing keys", func() {
		count := 0
		for _, expected := range expectedMap {
			id := expected.Id
			_, found := suite.keeper.GetInsuranceBid(suite.ctx, id)
			suite.Require().Equal(true, found)
			suite.keeper.DeleteInsuranceBid(suite.ctx, id)
			_, found = suite.keeper.GetInsuranceBid(suite.ctx, id)
			suite.Require().Equal(false, found)
			count++
		}
		suite.Require().Equal(len(expectedMap), count)
	})
}

func (suite *KeeperTestSuite) TestStoreGetAllInsuranceBids() {
	seed := tmrand.Int63n(10000)
	tmrand.Seed(seed)

	expectedMap := make(map[uint64]types.InsuranceBid)

	for i := 0; i < 1000; i++ {
		id := generateRandomId()
		expected := generateRandomInsuranceBid(suite, id)
		suite.keeper.SetInsuranceBid(suite.ctx, expected)
		expectedMap[id] = expected
	}
	reqs := suite.keeper.GetAllInsuranceBids(suite.ctx)
	suite.Require().Equal(len(expectedMap), len(reqs))
	for _, actual := range reqs {
		expected, found := expectedMap[actual.Id]
		suite.Require().True(found)
		suite.Require().Equal(expected, actual)
	}
}

func (suite *KeeperTestSuite) TestStoreInsuranceUnbondRequest() {
	seed := tmrand.Int63n(10000)
	tmrand.Seed(seed)
	_, found := suite.keeper.GetInsuranceUnbondRequest(suite.ctx, types.AliveChunkId(0))
	suite.Require().Equal(false, found)
	expectedMap := make(map[types.AliveChunkId]types.InsuranceUnbondRequest)
	const numIteration = 100

	suite.Run("Get/Set insurance with random seed: "+strconv.FormatInt(seed, 10), func() {
		count := 0
		for i := 0; i < numIteration; i++ {
			id := generateRandomId()
			if _, has := expectedMap[id]; has {
				count++
				suite.Require().Less(count, 10)
				continue
			}

			expected := generateRandomInsuranceUnbondRequest(id)
			expectedMap[id] = expected
			suite.keeper.SetInsuranceUnbondRequest(suite.ctx, expected)
			ret, found := suite.keeper.GetInsuranceUnbondRequest(suite.ctx, id)
			suite.Require().Equal(true, found)
			suite.Require().Equal(expected, ret)
		}
	})
	suite.Run("Delete insurance with random seed: "+strconv.FormatInt(seed, 10), func() {
		count := 0
		for i := 0; i < numIteration; i++ {
			id := generateRandomId()
			if _, has := expectedMap[id]; has {
				count++
				suite.Require().Less(count, 10)
				continue
			}
			_, found := suite.keeper.GetInsuranceUnbondRequest(suite.ctx, id)
			suite.Require().Equal(false, found)
			suite.keeper.DeleteInsuranceUnbondRequest(suite.ctx, id)
			_, found = suite.keeper.GetInsuranceUnbondRequest(suite.ctx, id)
			suite.Require().Equal(false, found)
		}
	})
	suite.Run("Delete insurance with existing keys", func() {
		count := 0
		for _, expected := range expectedMap {
			id := expected.AliveChunkId
			_, found := suite.keeper.GetInsuranceUnbondRequest(suite.ctx, id)
			suite.Require().Equal(true, found)
			suite.keeper.DeleteInsuranceUnbondRequest(suite.ctx, id)
			_, found = suite.keeper.GetInsuranceUnbondRequest(suite.ctx, id)
			suite.Require().Equal(false, found)
			count++
		}
		suite.Require().Equal(len(expectedMap), count)
	})
}

func (suite *KeeperTestSuite) TestStoreGetAllInsuranceUnbondRequests() {
	seed := tmrand.Int63n(10000)
	tmrand.Seed(seed)

	expectedMap := make(map[uint64]types.InsuranceUnbondRequest)

	for i := 0; i < 1000; i++ {
		id := generateRandomId()
		expected := generateRandomInsuranceUnbondRequest(id)
		suite.keeper.SetInsuranceUnbondRequest(suite.ctx, expected)
		expectedMap[id] = expected
	}
	reqs := suite.keeper.GetAllInsuranceUnbondRequests(suite.ctx)
	suite.Require().Equal(len(expectedMap), len(reqs))
	for _, actual := range reqs {
		expected, found := expectedMap[actual.AliveChunkId]
		suite.Require().True(found)
		suite.Require().Equal(expected, actual)
	}
}
