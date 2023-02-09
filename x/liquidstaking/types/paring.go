package types

import sdk "github.com/cosmos/cosmos-sdk/types"

type Rankable interface {
	GetInsuranceFeeRate() sdk.Dec
	GetValidatorAddress() string
}

type Rankables []Rankable

type AliveChunks []AliveChunk
type AliveChunkId = uint64

type UnbondingChunks []UnbondingChunk
type UnbondingChunkId = uint64

type ChunkUnbondRequestedAliveChunk struct {
	AliveChunk
	Address string
}
type ChunkUnbondRequestedAliveChunks []ChunkUnbondRequestedAliveChunk

type InsuranceUnbondRequestedAliveChunk struct {
	AliveChunk
	InsuranceProviderAddress string
}

type InsuranceUnbondRequestedAliveChunks []InsuranceUnbondRequestedAliveChunk

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

func NewInsuranceUnbondRequest(insuranceProviderAddr string, aliveChunkId uint64) InsuranceUnbondRequest {
	return InsuranceUnbondRequest{
		InsuranceProviderAddress: insuranceProviderAddr,
		AliveChunkId:             aliveChunkId,
	}
}

type State struct {
	InsuranceBids           InsuranceBids
	InsuranceUnbondRequests InsuranceUnbondRequests
	ChunkBondRequests       ChunkBondRequests
	ChunkUnbondRequests     ChunkUnbondRequests
	AliveChunks             AliveChunks
	InsuranceUnbonded       InsuranceUnbondRequestedAliveChunks
	ChunkUnbonded           ChunkUnbondRequestedAliveChunks
}

var _ Rankable = &AliveChunk{}

func (aliveChunk *AliveChunk) GetInsuranceFeeRate() sdk.Dec {
	return aliveChunk.InsuranceFeeRate
}

type DelegationState struct {
	SortedValidators LiquidValidators
	DelegationMap    map[string]sdk.Int
}

func (state *DelegationState) ChangeDelegation(from, to string, amount sdk.Int) {
	if len(from) > 0 {
		state.DelegationMap[from] =
			state.DelegationMap[from].Sub(amount)
	}
	if len(to) > 0 {
		state.DelegationMap[to] =
			state.DelegationMap[to].Add(amount)
	}
}
