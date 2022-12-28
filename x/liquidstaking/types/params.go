package types

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

var (
	ParamStoreKeyLiquidBondDenom        = []byte("ParamStoreKeyLiquidBondDenom")
	ParamStoreKeyMinInsurancePercentage = []byte("ParamStoreKeyMinInsurancePercentage")
	ParamStoreKeyChunkSize              = []byte("ParamStoreKeyChunkSize")
	ParamStoreKeyMaxAliveChunk          = []byte("ParamStoreKeyMaxAliveChunk")

	DefaultLiquidBondDenom        = "lsstake"
	DefaultMinInsurancePercentage = sdk.NewDec(10)
	DefaultChunkSize              = sdk.NewInt(500)
	DefaultMaxAliveChunk          = sdk.NewInt(10)
)

func NewParams(liquidBondDenom string, minInsurancePercentage sdk.Dec, chunkSize, maxAliveChunk sdk.Int) Params {
	return Params{
		LiquidBondDenom:        liquidBondDenom,
		MinInsurancePercentage: minInsurancePercentage,
		ChunkSize:              chunkSize,
		MaxAliveChunk:          maxAliveChunk,
	}
}

func DefaultParams() Params {
	return NewParams(DefaultLiquidBondDenom, DefaultMinInsurancePercentage, DefaultChunkSize, DefaultMaxAliveChunk)
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

func validateMinInsurancePercentage(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	maximumSlashablePercentage := sdk.ZeroDec()
	if v.LT(maximumSlashablePercentage) {
		return fmt.Errorf("minInsurancePercentage should be larger than maximum slashable percentage within an epoch")
	}
	return nil
}

func validateChunkSize(i interface{}) error {
	v, ok := i.(sdk.Int)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if !v.IsPositive() {
		return fmt.Errorf("chunk size should be larger than 0")
	}
	return nil
}

func validateMaxAliveChunk(i interface{}) error {
	v, ok := i.(sdk.Int)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if !v.IsPositive() {
		return fmt.Errorf("maxAliveChunk should be larger than 0")
	}
	return nil
}

func (params *Params) Validate() error {
	if err := validateMinInsurancePercentage(params.MinInsurancePercentage); err != nil {
		return err
	}
	if err := validateChunkSize(params.ChunkSize); err != nil {
		return err
	}
	if err := validateMaxAliveChunk(params.MaxAliveChunk); err != nil {
		return err
	}
	return nil
}

var _ paramtypes.ParamSet = (*Params)(nil)

func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

func (params *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(ParamStoreKeyLiquidBondDenom, &params.LiquidBondDenom, validateLiquidBondDenom),
		paramtypes.NewParamSetPair(ParamStoreKeyMinInsurancePercentage, &params.MinInsurancePercentage, validateMinInsurancePercentage),
		paramtypes.NewParamSetPair(ParamStoreKeyChunkSize, &params.ChunkSize, validateChunkSize),
		paramtypes.NewParamSetPair(ParamStoreKeyMaxAliveChunk, &params.MaxAliveChunk, validateMaxAliveChunk),
	}
}
