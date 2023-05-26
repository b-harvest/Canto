package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

var (
	KeyR0         = []byte("R0")
	KeyUSoftCap   = []byte("USoftCap")
	KeyUHardCap   = []byte("UHardCap")
	KeyUOptimal   = []byte("UOptimal")
	KeySlope1     = []byte("Slope1")
	KeySlope2     = []byte("Slope2")
	KeyMaxFeeRate = []byte("MaxFeeRate")

	DefaultR0       = sdk.ZeroDec()
	DefaultUSoftCap = sdk.MustNewDecFromStr("0.05")
	DefaultUHardCap = sdk.MustNewDecFromStr("0.1")
	DefaultUOptimal = sdk.MustNewDecFromStr("0.09")
	DefaultSlope1   = sdk.MustNewDecFromStr("0.1")
	DefaultSlope2   = sdk.MustNewDecFromStr("0.4")
	DefaultMaxFee   = sdk.MustNewDecFromStr("0.5")
)

var _ paramtypes.ParamSet = &Params{}

// ParamKeyTable returns the parameter key table.
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// NewParams creates a new Params object
func NewParams(
	r0, uSoftCap, uHardCap, uOptimal, slope1, slope2, maxFeeRate sdk.Dec,
) Params {
	return Params{
		R0:         r0,
		USoftCap:   uSoftCap,
		UHardCap:   uHardCap,
		UOptimal:   uOptimal,
		Slope1:     slope1,
		Slope2:     slope2,
		MaxFeeRate: maxFeeRate,
	}
}

func DefaultParams() Params {
	return NewParams(DefaultR0, DefaultUSoftCap, DefaultUHardCap, DefaultUOptimal, DefaultSlope1, DefaultSlope2, DefaultMaxFee)
}

// ParamSetPairs returns the parameter set pairs.
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyR0, &p.R0, validateR0),
		paramtypes.NewParamSetPair(KeyUSoftCap, &p.USoftCap, validateUSoftCap),
		paramtypes.NewParamSetPair(KeyUHardCap, &p.UHardCap, validateUHardCap),
		paramtypes.NewParamSetPair(KeyUOptimal, &p.UOptimal, validateUOptimal),
		paramtypes.NewParamSetPair(KeySlope1, &p.Slope1, validateSlope1),
		paramtypes.NewParamSetPair(KeySlope2, &p.Slope2, validateSlope2),
		paramtypes.NewParamSetPair(KeyMaxFeeRate, &p.MaxFeeRate, validateMaxFeeRate),
	}
}

func (p Params) Validate() error {
	for _, v := range []struct {
		value     interface{}
		validator func(interface{}) error
	}{
		{p.R0, validateR0},
		{p.USoftCap, validateUSoftCap},
		{p.UHardCap, validateUHardCap},
		{p.UOptimal, validateUOptimal},
		{p.Slope1, validateSlope1},
		{p.Slope2, validateSlope2},
		{p.MaxFeeRate, validateMaxFeeRate},
	} {
		if err := v.validator(v.value); err != nil {
			return err
		}
		// validate dynamic fee model
		if !p.USoftCap.LT(p.UOptimal) {
			return fmt.Errorf("uSoftCap should be less than uOptimal")
		}
		if !p.UOptimal.LT(p.UHardCap) {
			return fmt.Errorf("uOptimal should be less than uHardCap")
		}
		if !p.R0.Add(p.Slope1).Add(p.Slope2).LTE(p.MaxFeeRate) {
			return fmt.Errorf("r0 + slope1 + slope2 should not exceeds max fee rate")
		}
	}
	return nil
}

func validateR0(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNil() {
		return fmt.Errorf("r0 should not be nil")
	}

	if v.IsNegative() {
		return fmt.Errorf("r0 should not be negative")
	}

	if v.GT(sdk.OneDec()) {
		return fmt.Errorf("r0 should not be greater than 1")
	}

	return nil
}

func validateUSoftCap(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNil() {
		return fmt.Errorf("uSoftCap should not be nil")
	}

	if v.IsNegative() {
		return fmt.Errorf("uSoftCap should not be negative")
	}

	if v.GT(sdk.OneDec()) {
		return fmt.Errorf("uSoftCap should not be greater than 1")
	}

	return nil
}

func validateUHardCap(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNil() {
		return fmt.Errorf("uHardCap should not be nil")
	}

	if v.IsNegative() {
		return fmt.Errorf("uHardCap should not be negative")
	}

	if v.GT(sdk.OneDec()) {
		return fmt.Errorf("uHardCap should not be greater than 1")
	}

	return nil
}

func validateUOptimal(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNil() {
		return fmt.Errorf("uOptimal should not be nil")
	}

	if v.IsNegative() {
		return fmt.Errorf("uOptimal should not be negative")
	}

	if v.GT(sdk.OneDec()) {
		return fmt.Errorf("uOptimal should not be greater than 1")
	}

	return nil
}

func validateSlope1(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNil() {
		return fmt.Errorf("slope1 should not be nil")
	}

	if v.IsNegative() {
		return fmt.Errorf("slope1 should not be negative")
	}

	if v.GT(sdk.OneDec()) {
		return fmt.Errorf("slope1 should not be greater than 1")
	}

	return nil
}

func validateSlope2(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNil() {
		return fmt.Errorf("slope2 should not be nil")
	}

	if v.IsNegative() {
		return fmt.Errorf("slope2 should not be negative")
	}

	if v.GT(sdk.OneDec()) {
		return fmt.Errorf("slope2 should not be greater than 1")
	}

	return nil
}

func validateMaxFeeRate(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNil() {
		return fmt.Errorf("maxFeeRate should not be nil")
	}

	if v.IsNegative() {
		return fmt.Errorf("maxFeeRate should not be negative")
	}

	if v.GT(sdk.OneDec()) {
		return fmt.Errorf("maxFeeRate should not be greater than 1")
	}

	return nil
}
