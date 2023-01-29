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

func NativeTokenToLiquidToken(state LiquidStakingState, nativeTokenAmount sdk.Int) (sdk.Int, error) {
	// TODO: calc
	return nativeTokenAmount, nil
}

func LiquidTokenToNativeToken(state LiquidStakingState, liquidTokenAmount sdk.Int) (sdk.Int, error) {
	// TODO: calc
	return liquidTokenAmount, nil
}

var _ Rankable = &InsuranceBid{}

func (bid *InsuranceBid) GetInsuranceFeeRate() sdk.Dec {
	return bid.InsuranceFeeRate
}
