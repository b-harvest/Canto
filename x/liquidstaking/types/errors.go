package types

import sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

var (
	ErrInsufficientBalance        = sdkerrors.Register(ModuleName, 30000, "insufficient balance error")
	ErrMaxPairedChunkSizeExceeded = sdkerrors.Register(ModuleName, 30001, "reached maximum limit of paired chunk so cannot accept any more chunks.")
	ErrNoPairingInsurance         = sdkerrors.Register(ModuleName, 30002, "pairing insurance must exist to accept liquid stake request.")
	ErrInvalidAmount              = sdkerrors.Register(ModuleName, 30003, "amount of coin must be greater than or equal to 5M acanto.")
)
