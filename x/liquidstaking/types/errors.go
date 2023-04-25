package types

import sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

var (
	ErrMaxPairedChunkSizeExceeded           = sdkerrors.Register(ModuleName, 30001, "reached maximum limit of paired chunk so cannot accept any more chunks.")
	ErrNoPairingInsurance                   = sdkerrors.Register(ModuleName, 30002, "pairing insurance must exist to accept liquid stake request.")
	ErrInvalidAmount                        = sdkerrors.Register(ModuleName, 30003, "amount of coin must be multiple of the chunk size")
	ErrTombstonedValidator                  = sdkerrors.Register(ModuleName, 30005, "validator is tombstoned")
	ErrInvalidValidatorStatus               = sdkerrors.Register(ModuleName, 30006, "invalid validator status")
	ErrPairingInsuranceNotFound             = sdkerrors.Register(ModuleName, 30007, "pairing insurance not found")
	ErrNotProviderOfInsurance               = sdkerrors.Register(ModuleName, 30008, "not provider of insurance")
	ErrNotFoundInsurance                    = sdkerrors.Register(ModuleName, 30009, "insurance not found")
	ErrNoPairedChunk                        = sdkerrors.Register(ModuleName, 30010, "no paired chunk")
	ErrNotFoundChunk                        = sdkerrors.Register(ModuleName, 30011, "chunk not found")
	ErrInvalidChunkStatus                   = sdkerrors.Register(ModuleName, 30012, "invalid chunk status")
	ErrInvalidInsuranceStatus               = sdkerrors.Register(ModuleName, 30013, "invalid insurance status")
	ErrExceedAvailableChunks                = sdkerrors.Register(ModuleName, 30014, "exceed available chunks")
	ErrInvalidBondDenom                     = sdkerrors.Register(ModuleName, 30015, "invalid bond denom")
	ErrInvalidLiquidBondDenom               = sdkerrors.Register(ModuleName, 30016, "invalid liquid bond denom")
	ErrNotInWithdrawableStatus              = sdkerrors.Register(ModuleName, 30017, "insurance is not in withdrawable status")
	ErrUnpairingChunkHavePairedChunk        = sdkerrors.Register(ModuleName, 30018, "unpairing chunk cannot have paired chunk")
	ErrUnbondingDelegationNotRemoved        = sdkerrors.Register(ModuleName, 30019, "unbonding delegation not removed")
	ErrNotFoundUnpairingForUnstakeChunkInfo = sdkerrors.Register(ModuleName, 30020, "unstake chunk info not found")
	ErrNotFoundDelegation                   = sdkerrors.Register(ModuleName, 30021, "delegation not found")
	ErrNotFoundValidator                    = sdkerrors.Register(ModuleName, 30022, "validator not found")
)
