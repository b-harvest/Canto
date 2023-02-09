package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	TypeMsgLiquidStaking         = "liquid_staking"
	TypeMsgCancelLiquidStaking   = "cancel_liquid_staking"
	TypeMsgLiquidUnstaking       = "liquid_unstaking"
	TypeMsgCancelLiquidUnstaking = "cancel_liquid_unstaking"
	TypeMsgBidInsurance          = "bid_insurance"
	TypeMsgCancelInsuranceBid    = "cancel_insurance_bid"
	TypeMsgUnbondInsurance       = "unbond_insurance"
	TypeMsgCancelInsuranceUnbond = "cancel_insurance_unbond"
)

var (
	_ sdk.Msg = &MsgLiquidStaking{}
	_ sdk.Msg = &MsgCancelLiquidStaking{}
	_ sdk.Msg = &MsgLiquidUnstaking{}
	_ sdk.Msg = &MsgCancelLiquidUnstaking{}
	_ sdk.Msg = &MsgBidInsurance{}
	_ sdk.Msg = &MsgCancelInsuranceBid{}
	_ sdk.Msg = &MsgUnbondInsurance{}
	_ sdk.Msg = &MsgCancelInsuranceUnbond{}
)

// Route returns the name of the module
func (msg MsgLiquidStaking) Route() string { return RouterKey }

// Type returns the the action
func (msg MsgLiquidStaking) Type() string { return TypeMsgLiquidStaking }

// ValidateBasic runs stateless checks on the message
func (msg MsgLiquidStaking) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.RequesterAddress); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid requester address %q: %v", msg.RequesterAddress, err)
	}
	if !msg.TokenAmount.IsPositive() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "token amount should be greater than 0")
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg *MsgLiquidStaking) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgLiquidStaking) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.RequesterAddress)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr}
}

func (msg MsgLiquidStaking) GetRequester() sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.RequesterAddress)
	if err != nil {
		panic(err)
	}
	return addr
}

// Route returns the name of the module
func (msg MsgCancelLiquidStaking) Route() string { return RouterKey }

// Type returns the the action
func (msg MsgCancelLiquidStaking) Type() string { return TypeMsgCancelLiquidStaking }

// ValidateBasic runs stateless checks on the message
func (msg MsgCancelLiquidStaking) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.RequesterAddress); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid requester address %q: %v", msg.RequesterAddress, err)
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg *MsgCancelLiquidStaking) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgCancelLiquidStaking) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.RequesterAddress)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr}
}

// Route returns the name of the module
func (msg MsgLiquidUnstaking) Route() string { return RouterKey }

// Type returns the the action
func (msg MsgLiquidUnstaking) Type() string { return TypeMsgLiquidUnstaking }

// ValidateBasic runs stateless checks on the message
func (msg MsgLiquidUnstaking) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.RequesterAddress); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid requester address %q: %v", msg.RequesterAddress, err)
	}
	if msg.NumChunkUnstake == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "the number of unstaking chunk must not be 0")
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg *MsgLiquidUnstaking) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgLiquidUnstaking) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.RequesterAddress)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr}
}

func (msg MsgLiquidUnstaking) GetRequester() sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.RequesterAddress)
	if err != nil {
		panic(err)
	}
	return addr
}

// Route returns the name of the module
func (msg MsgCancelLiquidUnstaking) Route() string { return RouterKey }

// Type returns the the action
func (msg MsgCancelLiquidUnstaking) Type() string { return TypeMsgCancelLiquidUnstaking }

// ValidateBasic runs stateless checks on the message
func (msg MsgCancelLiquidUnstaking) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.RequesterAddress); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid requester address %q: %v", msg.RequesterAddress, err)
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg *MsgCancelLiquidUnstaking) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgCancelLiquidUnstaking) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.RequesterAddress)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr}
}

// Route returns the name of the module
func (msg MsgBidInsurance) Route() string { return RouterKey }

// Type returns the the action
func (msg MsgBidInsurance) Type() string { return TypeMsgBidInsurance }

// ValidateBasic runs stateless checks on the message
func (msg MsgBidInsurance) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.RequesterAddress); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid requester address %q: %v", msg.RequesterAddress, err)
	}
	if _, err := sdk.ValAddressFromBech32(msg.ValidatorAddress); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid validator address %q: %v", msg.ValidatorAddress, err)
	}
	if msg.InsuranceAmount.IsZero() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "staking insurance amount must not be zero")
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg *MsgBidInsurance) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgBidInsurance) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.RequesterAddress)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr}
}

func (msg MsgBidInsurance) GetInsurer() sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.RequesterAddress)
	if err != nil {
		panic(err)
	}
	return addr
}

// Route returns the name of the module
func (msg MsgCancelInsuranceBid) Route() string { return RouterKey }

// Type returns the the action
func (msg MsgCancelInsuranceBid) Type() string { return TypeMsgCancelInsuranceBid }

// ValidateBasic runs stateless checks on the message
func (msg MsgCancelInsuranceBid) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.RequesterAddress); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid requester address %q: %v", msg.RequesterAddress, err)
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg *MsgCancelInsuranceBid) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgCancelInsuranceBid) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.RequesterAddress)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr}
}

func (msg MsgCancelInsuranceBid) GetInsurer() sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.RequesterAddress)
	if err != nil {
		panic(err)
	}
	return addr
}

func (msg MsgBidInsurance) GetValidator() sdk.ValAddress {
	addr, err := sdk.ValAddressFromBech32(msg.ValidatorAddress)
	if err != nil {
		panic(err)
	}
	return addr
}

// Route returns the name of the module
func (msg MsgUnbondInsurance) Route() string { return RouterKey }

// Type returns the the action
func (msg MsgUnbondInsurance) Type() string { return TypeMsgUnbondInsurance }

// ValidateBasic runs stateless checks on the message
func (msg MsgUnbondInsurance) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.RequesterAddress); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid insurer address %q: %v", msg.RequesterAddress, err)
	}
	// TODO: need to validate insurance_id
	return nil
}

// GetSignBytes encodes the message for signing
func (msg *MsgUnbondInsurance) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgUnbondInsurance) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.RequesterAddress)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr}
}

func (msg MsgUnbondInsurance) GetInsurer() sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.RequesterAddress)
	if err != nil {
		panic(err)
	}
	return addr
}

// Route returns the name of the module
func (msg MsgCancelInsuranceUnbond) Route() string { return RouterKey }

// Type returns the the action
func (msg MsgCancelInsuranceUnbond) Type() string { return TypeMsgCancelInsuranceUnbond }

// ValidateBasic runs stateless checks on the message
func (msg MsgCancelInsuranceUnbond) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.RequesterAddress); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid insurer address %q: %v", msg.RequesterAddress, err)
	}
	// TODO: need to validate insurance_id
	return nil
}

// GetSignBytes encodes the message for signing
func (msg *MsgCancelInsuranceUnbond) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgCancelInsuranceUnbond) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.RequesterAddress)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr}
}

func (msg MsgCancelInsuranceUnbond) GetInsurer() sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.RequesterAddress)
	if err != nil {
		panic(err)
	}
	return addr
}
