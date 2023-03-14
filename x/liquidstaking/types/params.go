package types

import (
	fmt "fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"strings"
)

// Parameter store key
var (
	ParamStoreKeyLiquidBondDenom = []byte("LiquidBondDenom")
	DefaultLiquidBondDenom       = "lscanto"
)

var _ paramtypes.ParamSet = &Params{}

// ParamKeyTable returns the parameter key table.
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// NewParams creates a new Params object
func NewParams(
	liquidBondDenom string,
) Params {
	return Params{
		LiquidBondDenom: liquidBondDenom,
	}
}

func DefaultParams() Params {
	return Params{
		LiquidBondDenom: DefaultLiquidBondDenom,
	}
}

func validateLiquidBondDenom(i interface{}) error {
	v, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if strings.TrimSpace(v) == "" {
		return fmt.Errorf("liquid bond denom cannot be blank")
	}

	if err := sdk.ValidateDenom(v); err != nil {
		return err
	}
	return nil
}

// ParamSetPairs returns the parameter set pairs.
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(ParamStoreKeyLiquidBondDenom, &p.LiquidBondDenom, validateLiquidBondDenom),
	}
}

func (p Params) Validate() error {
	for _, v := range []struct {
		value     interface{}
		validator func(interface{}) error
	}{
		{p.LiquidBondDenom, validateLiquidBondDenom},
	} {
		if err := v.validator(v.value); err != nil {
			return err
		}
	}
	return nil
}
