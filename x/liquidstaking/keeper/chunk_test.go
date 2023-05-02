package keeper_test

import "github.com/Canto-Network/Canto/v6/x/liquidstaking/types"

// Sets a bunch of chunks in the store and then get and ensure that each of them
// match up with what is stored on stack vs keeper
func (suite *KeeperTestSuite) TestChunkSetGet() {
	numberChunks := 10
	chunks := GenerateChunks(numberChunks)
	for _, chunk := range chunks {
		suite.app.LiquidStakingKeeper.SetChunk(suite.ctx, chunk)
	}

	for _, chunk := range chunks {
		id := chunk.Id
		status := chunk.Status
		pairedInsuranceId := chunk.PairedInsuranceId
		unpairingInsuranceId := chunk.UnpairingInsuranceId
		// Get chunk from the store
		result, found := suite.app.LiquidStakingKeeper.GetChunk(suite.ctx, id)

		// Validation
		suite.Require().True(found)
		suite.Require().Equal(result.Id, id)
		suite.Require().Equal(result.Status, status)
		suite.Require().Equal(result.PairedInsuranceId, pairedInsuranceId)
		suite.Require().Equal(result.UnpairingInsuranceId, unpairingInsuranceId)
	}
}

func (suite *KeeperTestSuite) TestDeleteChunk() {
	numberChunks := 10
	chunks := GenerateChunks(numberChunks)
	for _, chunk := range chunks {
		suite.app.LiquidStakingKeeper.SetChunk(suite.ctx, chunk)
	}

	for _, chunk := range chunks {
		id := chunk.Id
		// Get chunk from the store
		result, found := suite.app.LiquidStakingKeeper.GetChunk(suite.ctx, id)

		// Validation
		suite.Require().True(found)
		suite.Require().Equal(result.Id, id)

		// Delete chunk from the store
		suite.app.LiquidStakingKeeper.DeleteChunk(suite.ctx, id)

		// Get chunk from the store
		result, found = suite.app.LiquidStakingKeeper.GetChunk(suite.ctx, id)

		// Validation
		suite.Require().False(found)
		suite.Require().Equal(result.Id, uint64(0))
	}
}

func (suite *KeeperTestSuite) TestLastChunkIdSetGet() {
	// Set LastChunkId and retrieve it
	id := uint64(10)
	suite.app.LiquidStakingKeeper.SetLastChunkId(suite.ctx, id)

	result := suite.app.LiquidStakingKeeper.GetLastChunkId(suite.ctx)
	suite.Require().Equal(result, id)
}

// Creates a bunch of chunks
func GenerateChunks(number int) []types.Chunk {
	chunks := make([]types.Chunk, number)

	for i := 0; i < number; i++ {
		chunks[i] = types.NewChunk(uint64(i))
	}
	return chunks
}
