package types

type LiquidValidators []LiquidValidator
type LiquidValidatorStates []LiquidValidatorState

type ChunkBondRequests []ChunkBondRequest
type ChunkBondRequestId = uint64

type ChunkUnbondRequests []ChunkUnbondRequest
type ChunkUnbondRequestId = uint64

type InsuranceBids []InsuranceBid
type InsuranceBidId = uint64

type InsuranceUnbondRequests []InsuranceUnbondRequest

// type InsuranceUnbondRequestId = uint64

func NewUnbondingChunk(aliveChunk AliveChunk) UnbondingChunk {
	return UnbondingChunk{
		Id:                       aliveChunk.Id,
		ValidatorAddress:         aliveChunk.ValidatorAddress,
		InsuranceProviderAddress: aliveChunk.InsuranceProviderAddress,
		TokenAmount:              aliveChunk.TokenAmount,
		InsuranceAmount:          aliveChunk.InsuranceAmount,
	}
}
