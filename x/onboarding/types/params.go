package types

import (
	"fmt"
	"time"

	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// Parameter store key
var (
	ParamStoreKeyEnableOnboarding = []byte("EnableOnboarding")
)

// DefaultPacketTimeoutDuration defines the default packet timeout for outgoing
// IBC transfers
var DefaultPacketTimeoutDuration = 4 * time.Hour

var _ paramtypes.ParamSet = &Params{}

// ParamKeyTable returns the parameter key table.
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// NewParams creates a new Params instance
func NewParams(
	enableOnboarding bool,
) Params {
	return Params{
		EnableOnboarding: enableOnboarding,
	}
}

// DefaultParams defines the default params for the onboarding module
func DefaultParams() Params {
	return Params{
		EnableOnboarding: true,
	}
}

// ParamSetPairs returns the parameter set pairs.
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(ParamStoreKeyEnableOnboarding, &p.EnableOnboarding, validateBool),
	}
}

func validateBool(i interface{}) error {
	_, ok := i.(bool)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	return nil
}

// Validate checks that the fields have valid values
func (p Params) Validate() error {
	return validateBool(p.EnableOnboarding)
}
