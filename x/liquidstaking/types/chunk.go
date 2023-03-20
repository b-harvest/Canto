package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	MaxPairedChunks = 10
	ChunkSize       = 5000000 // 5M acanto
)

func NewChunk(id uint64) Chunk {
	return Chunk{
		Id:          id,
		InsuranceId: 0, // Not yet assigned
		Status:      CHUNK_STATUS_PAIRING,
	}
}

func (c *Chunk) DerivedAddress() sdk.AccAddress {
	return DeriveAddress(ModuleName, fmt.Sprintf("chunk%d", c.Id))
}
