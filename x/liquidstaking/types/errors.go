package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// Sentinel errors for the liquidstaking module.
var (
	ErrNonActiveLiquidValidators = sdkerrors.Register(ModuleName, 2, "given liquid validator is not active")
	ErrInvalidTokenAmount        = sdkerrors.Register(ModuleName, 3, "given token amount is invalid")
	ErrInvalidDenom              = sdkerrors.Register(ModuleName, 4, "given denom is invalid")
	ErrInvalidInsuranceId        = sdkerrors.Register(ModuleName, 5, "given insurance ID is invalid")
	ErrInvalidAliveChunkId       = sdkerrors.Register(ModuleName, 6, "given alive chunk ID is invalid")
)
