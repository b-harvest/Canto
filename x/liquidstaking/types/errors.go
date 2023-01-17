package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// Sentinel errors for the liquidstaking module.
var (
	ErrNonActiveLiquidValidators = sdkerrors.Register(ModuleName, 2, "given liquid validator is not active")
	ErrInvalidTokenAmount        = sdkerrors.Register(ModuleName, 3, "given token amount is invalid")
)
