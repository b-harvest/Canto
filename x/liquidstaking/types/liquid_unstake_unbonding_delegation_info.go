package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"time"
)

func NewLiquidUnstakeUnbondingDelegationInfo(chunkId uint64, delegatorAddress, validatorAddress string, burnAmount sdk.Coin, completionTime time.Time) LiquidUnstakeUnbondingDelegationInfo {
	return LiquidUnstakeUnbondingDelegationInfo{
		ChunkId:          chunkId,
		DelegatorAddress: delegatorAddress,
		ValidatorAddress: validatorAddress,
		BurnAmount:       burnAmount,
		CompletionTime:   completionTime,
	}
}

func (info *LiquidUnstakeUnbondingDelegationInfo) GetDelegator() sdk.AccAddress {
	return sdk.MustAccAddressFromBech32(info.DelegatorAddress)
}

func (info *LiquidUnstakeUnbondingDelegationInfo) GetValidator() sdk.ValAddress {
	valAddr, _ := sdk.ValAddressFromBech32(info.ValidatorAddress)
	return valAddr
}
