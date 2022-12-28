package types

import (
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

func NewParams() Params {
	return Params{}
}

func DefaultParams() Params {
	return NewParams()
}

func (params *Params) Validate() error {
	return nil
}

// func (params Params) String() string {
// 	out, _ := yaml.Marshal(params)
// 	return string(out)
// }

var _ paramtypes.ParamSet = (*Params)(nil)

func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

func (params *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{}
}
