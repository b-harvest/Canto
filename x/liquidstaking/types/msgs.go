package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	TypeMsgLiquidStake         = "liquid_stake"
	TypeMsgLiquidUnstake       = "liquid_unstake"
	TypeMsgRegisterInsurance   = "register_insurance"
	TypeMsgUnregisterInsurance = "unregister_insurance"
)

var (
	_ sdk.Msg = &MsgLiquidStake{}
	_ sdk.Msg = &MsgLiquidUnstake{}
	_ sdk.Msg = &MsgRegisterInsurance{}
	_ sdk.Msg = &MsgUnregisterInsurance{}
)

// Route returns the name of the module
func (msg MsgLiquidStake) Route() string { return RouterKey }

// Type returns the the action
func (msg MsgLiquidStake) Type() string { return TypeMsgLiquidStake }

// ValidateBasic runs stateless checks on the message
func (msg MsgLiquidStake) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.DelegatorAddress); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid delegator address %q: %v", msg.DelegatorAddress, err)
	}
	if msg.NumChunks == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "the number of staking chunk must not be 0")
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg *MsgLiquidStake) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgLiquidStake) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.DelegatorAddress)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr}
}

func (msg MsgLiquidStake) GetDelegator() sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.DelegatorAddress)
	if err != nil {
		panic(err)
	}
	return addr
}

// Route returns the name of the module
func (msg MsgLiquidUnstake) Route() string { return RouterKey }

// Type returns the the action
func (msg MsgLiquidUnstake) Type() string { return TypeMsgLiquidUnstake }

// ValidateBasic runs stateless checks on the message
func (msg MsgLiquidUnstake) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.DelegatorAddress); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid delegator address %q: %v", msg.DelegatorAddress, err)
	}
	if msg.NumChunks == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "the number of unstaking chunk must not be 0")
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg *MsgLiquidUnstake) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgLiquidUnstake) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.DelegatorAddress)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr}
}

func (msg MsgLiquidUnstake) GetDelegator() sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.DelegatorAddress)
	if err != nil {
		panic(err)
	}
	return addr
}

// Route returns the name of the module
func (msg MsgRegisterInsurance) Route() string { return RouterKey }

// Type returns the the action
func (msg MsgRegisterInsurance) Type() string { return TypeMsgRegisterInsurance }

// ValidateBasic runs stateless checks on the message
func (msg MsgRegisterInsurance) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.InsurerAddress); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid insurer address %q: %v", msg.InsurerAddress, err)
	}
	if msg.Amount.IsZero() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "staking insurance amount must not be zero")
	}
	return msg.Amount.Validate()
}

// GetSignBytes encodes the message for signing
func (msg *MsgRegisterInsurance) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgRegisterInsurance) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.InsurerAddress)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr}
}

func (msg MsgRegisterInsurance) GetInsurer() sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.InsurerAddress)
	if err != nil {
		panic(err)
	}
	return addr
}

func (msg MsgRegisterInsurance) GetValidator() sdk.ValAddress {
	addr, err := sdk.ValAddressFromBech32(msg.ValidatorAddress)
	if err != nil {
		panic(err)
	}
	return addr
}

// Route returns the name of the module
func (msg MsgUnregisterInsurance) Route() string { return RouterKey }

// Type returns the the action
func (msg MsgUnregisterInsurance) Type() string { return TypeMsgUnregisterInsurance }

// ValidateBasic runs stateless checks on the message
func (msg MsgUnregisterInsurance) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.InsurerAddress); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid insurer address %q: %v", msg.InsurerAddress, err)
	}
	// TODO: need to validate insurance_id
	return nil
}

// GetSignBytes encodes the message for signing
func (msg *MsgUnregisterInsurance) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgUnregisterInsurance) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.InsurerAddress)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr}
}

func (msg MsgUnregisterInsurance) GetInsurer() sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.InsurerAddress)
	if err != nil {
		panic(err)
	}
	return addr
}
