package types

import sdk "github.com/cosmos/cosmos-sdk/types"

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

func NativeTokenToLiquidToken(state LiquidStakingState, nativeTokenAmount sdk.Int) (sdk.Int, error) {
	// TODO: calc
	return nativeTokenAmount, nil
}

func LiquidTokenToNativeToken(state LiquidStakingState, liquidTokenAmount sdk.Int) (sdk.Int, error) {
	// TODO: calc
	return liquidTokenAmount, nil
}
