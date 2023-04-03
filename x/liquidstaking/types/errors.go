package types

import sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

var (
	ErrMaxPairedChunkSizeExceeded  = sdkerrors.Register(ModuleName, 30001, "reached maximum limit of paired chunk so cannot accept any more chunks.")
	ErrNoPairingInsurance          = sdkerrors.Register(ModuleName, 30002, "pairing insurance must exist to accept liquid stake request.")
	ErrInvalidAmount               = sdkerrors.Register(ModuleName, 30003, "amount of coin must be greater than or equal to 5M acanto.")
	ErrValidatorNotFound           = sdkerrors.Register(ModuleName, 30004, "validator not found")
	ErrTombstonedValidator         = sdkerrors.Register(ModuleName, 30005, "validator is tombstoned")
	ErrPairingInsuranceNotFound    = sdkerrors.Register(ModuleName, 30006, "pairing insurance not found")
	ErrNotProviderOfInsurance      = sdkerrors.Register(ModuleName, 30007, "not provider of insurance")
	ErrNotFoundInsurance           = sdkerrors.Register(ModuleName, 30008, "insurance not found")
	ErrInvalidUnstakeAmount        = sdkerrors.Register(ModuleName, 30009, "unstake amount must be a multiple of the chunk size")
	ErrNoPairedChunk               = sdkerrors.Register(ModuleName, 30010, "no paired chunk")
	ErrNotFoundChunk               = sdkerrors.Register(ModuleName, 30011, "chunk not found")
	ErrNotFoundUnbondingDelegation = sdkerrors.Register(ModuleName, 30012, "unbonding delegation not found")
	ErrInvalidChunkStatus          = sdkerrors.Register(ModuleName, 30013, "invalid chunk status")
	ErrNotFoundDelegation          = sdkerrors.Register(ModuleName, 30014, "delegation not found")
)
