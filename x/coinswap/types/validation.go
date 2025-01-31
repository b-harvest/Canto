package types

import (
	fmt "fmt"
	"strings"

	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// ValidateInput verifies whether the given input is legal
func ValidateInput(input Input) error {
	if !(input.Coin.IsValid() && input.Coin.IsPositive()) {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidCoins, "invalid input (%s)", input.Coin.String())
	}

	if strings.HasPrefix(input.Coin.Denom, LptTokenPrefix) {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "invalid input denom, should not begin with (%s)", LptTokenPrefix)
	}

	if _, err := sdk.AccAddressFromBech32(input.Address); err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid input address (%s)", err)
	}
	return nil
}

// ValidateOutput verifies whether the given output is legal
func ValidateOutput(output Output) error {
	if !(output.Coin.IsValid() && output.Coin.IsPositive()) {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidCoins, "invalid output (%s)", output.Coin.String())
	}

	if strings.HasPrefix(output.Coin.Denom, LptTokenPrefix) {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "invalid output denom, should not begin with (%s)", LptTokenPrefix)
	}

	if _, err := sdk.AccAddressFromBech32(output.Address); err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid output address (%s)", err)
	}
	return nil
}

// ValidateDeadline verifies whether the given deadline is legal
func ValidateDeadline(deadline int64) error {
	if deadline <= 0 {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("deadline %d must be greater than 0", deadline))
	}
	return nil
}

// ValidateMaxToken verifies whether the maximum token is legal
func ValidateMaxToken(maxToken sdk.Coin) error {
	if !(maxToken.IsValid() && maxToken.IsPositive()) {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidCoins, "invalid maxToken (%s)", maxToken.String())
	}

	if strings.HasPrefix(maxToken.Denom, LptTokenPrefix) {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "max token must be non-liquidity token")
	}
	return nil
}

// ValidateExactStandardAmt verifies whether the standard token amount is legal
func ValidateExactStandardAmt(standardAmt sdkmath.Int) error {
	if !standardAmt.IsPositive() {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "standard token amount must be positive")
	}
	return nil
}

// ValidateMinLiquidity verifies whether the minimum liquidity is legal
func ValidateMinLiquidity(minLiquidity sdkmath.Int) error {
	if minLiquidity.IsNegative() {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "minimum liquidity can not be negative")
	}
	return nil
}

// ValidateMinToken verifies whether the minimum token amount is legal
func ValidateMinToken(minToken sdkmath.Int) error {
	if minToken.IsNegative() {
		return errorsmod.Wrap(sdkerrors.ErrInvalidCoins, "minimum token amount can not be negative")
	}
	return nil
}

// ValidateWithdrawLiquidity verifies whether the given liquidity is legal
func ValidateWithdrawLiquidity(liquidity sdk.Coin) error {
	if !liquidity.IsValid() || !liquidity.IsPositive() {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidCoins, "invalid withdrawLiquidity (%s)", liquidity.String())
	}

	if err := ValidateLptDenom(liquidity.Denom); err != nil {
		return err
	}
	return nil
}

// ValidateMinStandardAmt verifies whether the minimum standard amount is legal
func ValidateMinStandardAmt(minStandardAmt sdkmath.Int) error {
	if minStandardAmt.IsNegative() {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("minimum standard token amount %s can not be negative", minStandardAmt.String()))
	}
	return nil
}

// ValidateLptDenom returns nil if the Liquidity pool token denom is valid
func ValidateLptDenom(lptDenom string) error {
	if _, err := ParseLptDenom(lptDenom); err != nil {
		return errorsmod.Wrap(ErrInvalidDenom, lptDenom)
	}
	return nil
}
