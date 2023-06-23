package types_test

import (
	"github.com/Canto-Network/Canto/v6/x/liquidstaking/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"
	"github.com/tendermint/tendermint/crypto"
	"testing"
)

type chunksTestSuite struct {
	suite.Suite
}

func TestChunksTestSuite(t *testing.T) {
	suite.Run(t, new(chunksTestSuite))
}

func (suite *chunksTestSuite) TestDerivedAddress() {
	c := types.NewChunk(1)
	suite.Equal(
		sdk.AccAddress(crypto.AddressHash([]byte("liquidstakingchunk1"))).String(),
		c.DerivedAddress().String(),
	)
	suite.Equal(
		"A88056CA3B6E75677FD17A846C255361E3B8DA20",
		c.DerivedAddress().String(),
	)
}

func (suite *chunksTestSuite) TestEqual() {
	c1 := types.NewChunk(1)

	c2 := c1
	c2.Id = 2
	suite.False(c1.Equal(c2))

	c2 = c1
	suite.True(c1.Equal(c2))
	c2.PairedInsuranceId = 2
	suite.False(c1.Equal(c2))

	c1 = c2
	c2.UnpairingInsuranceId = 2
	suite.False(c1.Equal(c2))

	c1 = c2
	c2.Status = types.CHUNK_STATUS_UNPAIRING
	suite.False(c1.Equal(c2))
}

func (suite *chunksTestSuite) TestSetStatus() {
	c := types.NewChunk(1)
	suite.Equal(types.CHUNK_STATUS_PAIRING, c.Status)
	c.SetStatus(types.CHUNK_STATUS_PAIRED)
	suite.Equal(types.CHUNK_STATUS_PAIRED, c.Status)
}

func (suite *chunksTestSuite) TestValidate() {
	c := types.NewChunk(2)
	suite.NoError(c.Validate(2))
	suite.Error(c.Validate(1))
	c.Status = types.CHUNK_STATUS_UNSPECIFIED
	suite.Error(c.Validate(2))
}

func (suite *chunksTestSuite) TestHasPairedInsurance() {
	c := types.NewChunk(1)
	suite.False(c.HasPairedInsurance())
	c.PairedInsuranceId = 1
	suite.True(c.HasPairedInsurance())
}
