package types

type AliveChunks []AliveChunk
type AliveChunkId = uint64

type UnbondingChunks []UnbondingChunk
type UnbondingChunkId = uint64

func NewAliveChunk(id AliveChunkId, chunkBondRequest ChunkBondRequest, insuranceBid InsuranceBid) AliveChunk {
	return AliveChunk{
		Id:                       id,
		ValidatorAddress:         insuranceBid.ValidatorAddress,
		InsuranceProviderAddress: insuranceBid.InsuranceProviderAddress,
		TokenAmount:              chunkBondRequest.TokenAmount,
		InsuranceAmount:          insuranceBid.InsuranceAmount,
		InsuranceFeeRate:         insuranceBid.InsuranceFeeRate,
	}
}

func NewUnbondingChunk(aliveChunk AliveChunk) UnbondingChunk {
	return UnbondingChunk{
		Id:                       aliveChunk.Id,
		ValidatorAddress:         aliveChunk.ValidatorAddress,
		InsuranceProviderAddress: aliveChunk.InsuranceProviderAddress,
		TokenAmount:              aliveChunk.TokenAmount,
		InsuranceAmount:          aliveChunk.InsuranceAmount,
	}
}
