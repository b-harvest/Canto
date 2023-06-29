package types_test

import (
	"testing"
	"time"

	"github.com/Canto-Network/Canto/v6/x/liquidstaking/types"
	"github.com/stretchr/testify/suite"
)

type redelegationInfoTestSuite struct {
	suite.Suite
}

func TestRedelegationInfoTestSuite(t *testing.T) {
	suite.Run(t, new(redelegationInfoTestSuite))
}

func (suite *redelegationInfoTestSuite) TestNewRedelegationInfo() {
	c := types.NewChunk(1)
	t := time.Now()
	info := types.NewRedelegationInfo(c.Id, t)
	suite.Equal(c.Id, info.ChunkId)
	suite.True(t.Equal(info.CompletionTime))
}

func (suite *redelegationInfoTestSuite) TestValidate() {
	c := types.NewChunk(1)
	t := time.Now()
	info := types.NewRedelegationInfo(c.Id, t)
	chunkMap := map[uint64]types.Chunk{
		c.Id: c,
	}
	suite.NoError(info.Validate(chunkMap))

	delete(chunkMap, c.Id)
	suite.Error(info.Validate(chunkMap))
}
