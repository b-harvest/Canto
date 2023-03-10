package types

import (
	"fmt"
	"time"

	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// Parameter store key
var (
	ParamStoreKeyEnableOnboarding      = []byte("EnableOnboarding")
	ParamStoreKeyPacketTimeoutDuration = []byte("PacketTimeoutDuration")
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
	enableOnboarding bool, timeoutDuration time.Duration,
) Params {
	return Params{
		EnableOnboarding:      enableOnboarding,
		PacketTimeoutDuration: timeoutDuration,
	}
}

// DefaultParams defines the default params for the onboarding module
func DefaultParams() Params {
	return Params{
		EnableOnboarding:      true,
		PacketTimeoutDuration: DefaultPacketTimeoutDuration,
	}
}

// ParamSetPairs returns the parameter set pairs.
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(ParamStoreKeyEnableOnboarding, &p.EnableOnboarding, validateBool),
		paramtypes.NewParamSetPair(ParamStoreKeyPacketTimeoutDuration, &p.PacketTimeoutDuration, validateDuration),
	}
}

func validateBool(i interface{}) error {
	_, ok := i.(bool)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	return nil
}

func validateDuration(i interface{}) error {
	duration, ok := i.(time.Duration)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if duration < 0 {
		return fmt.Errorf("packet timout duration cannot be negative")
	}

	return nil
}

// Validate checks that the fields have valid values
func (p Params) Validate() error {
	if err := validateDuration(p.PacketTimeoutDuration); err != nil {
		return err
	}

	return validateBool(p.EnableOnboarding)
}
