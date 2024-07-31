package app

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	tmstrings "github.com/tendermint/tendermint/libs/strings"
)

// EmptyAppOptions is a stub implementing AppOptions
type EmptyAppOptions struct{}

// Get implements AppOptions
func (eao EmptyAppOptions) Get(o string) interface{} {
	return nil
}

type RawValidator struct {
	Moniker string `yaml:"Moniker"`
	Address string `yaml:"Address"`
	//BalAmount    string `yaml:"BalAmount"`
	//StakeAmount  string `yaml:"StakeAmount"`
	ValidatorKey string `yaml:"ValidatorKey"`
	Mnemonic     string `yaml:"Mnemonic"`
}

type RawValidatorList []RawValidator

type Validator struct {
	Address        string
	VotingPower    sdk.Coins
	SelfDelegation sdk.Coin
	PublicKeyStr   string
	Moniker        string
}

type ValidatorKey struct {
	Type string `json:"@type"`
	Key  string `json:"key"`
}

type ValidatorList []Validator

func NewValidator(address string, votingPower sdk.Coins, selfDelegation sdk.Coin, pkStr string, moniker string) Validator {
	return Validator{
		Address:        address,
		VotingPower:    votingPower,
		SelfDelegation: selfDelegation,
		PublicKeyStr:   pkStr,
		Moniker:        moniker,
	}
}

func (v Validator) GetAddress() sdk.AccAddress {
	acc, err := sdk.AccAddressFromBech32(v.Address)
	if err != nil {
		panic(err)
	}
	return acc
}

func (v Validator) GetValAddress() sdk.ValAddress {
	acc, err := sdk.AccAddressFromBech32(v.Address)
	if err != nil {
		panic(err)
	}
	return sdk.ValAddress(acc)
}

func (v Validator) GetPubKey(codec codec.Codec) cryptotypes.PubKey {
	var pk cryptotypes.PubKey
	if err := codec.UnmarshalInterfaceJSON([]byte(v.PublicKeyStr), &pk); err != nil {
		panic(err)
	}

	return pk
}

func (v Validator) GetSelfDelegation() sdk.Coin {
	return v.SelfDelegation
}

func (v Validator) NewMsgCreateValidator(codec codec.Codec) (*stakingtypes.MsgCreateValidator, error) {
	msg, err := stakingtypes.NewMsgCreateValidator(
		v.GetValAddress(),
		v.GetPubKey(codec),
		v.GetSelfDelegation(),
		stakingtypes.Description{Moniker: v.Moniker},
		stakingtypes.CommissionRates{
			Rate:          sdk.OneDec(),
			MaxRate:       sdk.OneDec(),
			MaxChangeRate: sdk.OneDec(),
		},
		sdk.OneInt(),
	)
	if err != nil {
		return nil, err
	}
	return msg, nil
}

// CreateValidator defines a method for creating a new validator
func (v Validator) CreateValidator(ctx sdk.Context, k *stakingkeeper.Keeper, codec codec.Codec) error {
	msg, err := v.NewMsgCreateValidator(codec)
	if err != nil {
		return err
	}

	valAddr, err := sdk.ValAddressFromBech32(msg.ValidatorAddress)
	if err != nil {
		return err
	}

	// check to see if the pubkey or sender has been registered before
	if _, found := k.GetValidator(ctx, valAddr); found {
		return stakingtypes.ErrValidatorOwnerExists
	}

	pk, ok := msg.Pubkey.GetCachedValue().(cryptotypes.PubKey)
	if !ok {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidType, "Expecting cryptotypes.PubKey, got %T", pk)
	}

	if _, found := k.GetValidatorByConsAddr(ctx, sdk.GetConsAddress(pk)); found {
		return stakingtypes.ErrValidatorPubKeyExists
	}

	if _, err := msg.Description.EnsureLength(); err != nil {
		return err
	}

	cp := ctx.ConsensusParams()
	if cp != nil && cp.Validator != nil {
		if !tmstrings.StringInSlice(pk.Type(), cp.Validator.PubKeyTypes) {
			return sdkerrors.Wrapf(
				stakingtypes.ErrValidatorPubKeyTypeNotSupported,
				"got: %s, expected: %s", pk.Type(), cp.Validator.PubKeyTypes,
			)
		}
	}

	validator, err := stakingtypes.NewValidator(valAddr, pk, msg.Description)
	if err != nil {
		return err
	}
	commission := stakingtypes.NewCommissionWithTime(
		msg.Commission.Rate, msg.Commission.MaxRate,
		msg.Commission.MaxChangeRate, ctx.BlockHeader().Time,
	)

	validator, err = validator.SetInitialCommission(commission)
	if err != nil {
		return err
	}

	delegatorAddress, err := sdk.AccAddressFromBech32(msg.DelegatorAddress)
	if err != nil {
		return err
	}

	validator.MinSelfDelegation = msg.MinSelfDelegation

	k.SetValidator(ctx, validator)
	k.SetValidatorByConsAddr(ctx, validator)
	k.SetNewValidatorByPowerIndex(ctx, validator)

	// call the after-creation hook
	k.AfterValidatorCreated(ctx, validator.GetOperator())

	// move coins from the msg.Address account to a (self-delegation) delegator account
	// the validator account and global shares are updated within here
	// NOTE source will always be from a wallet which are unbonded
	_, err = k.Delegate(ctx, delegatorAddress, msg.Value.Amount, stakingtypes.Unbonded, validator, true)
	if err != nil {
		return err
	}
	return nil
}
