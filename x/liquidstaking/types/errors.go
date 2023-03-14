package types

import sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

var (
	ErrInsufficientBalance = sdkerrors.Register(ModuleName, 30000, "insufficient balance error")
)
