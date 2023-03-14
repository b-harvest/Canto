package keeper

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	vestexported "github.com/cosmos/cosmos-sdk/x/auth/vesting/exported"
	transfertypes "github.com/cosmos/ibc-go/v3/modules/apps/transfer/types"
	channeltypes "github.com/cosmos/ibc-go/v3/modules/core/04-channel/types"
	host "github.com/cosmos/ibc-go/v3/modules/core/24-host"
	"github.com/cosmos/ibc-go/v3/modules/core/exported"

	"github.com/Canto-Network/Canto/v6/ibc"
	"github.com/Canto-Network/Canto/v6/x/onboarding/types"
)

// OnRecvPacket performs an IBC receive callback. It returns the tokens that
// users transferred to their Cosmos secp256k1 address instead of the Ethereum
// ethsecp256k1 address. The expected behavior is as follows:
//
// First transfer from authorized source chain:
//  - sends back IBC tokens which originated from the source chain
//  - sends over all canto native tokens
// Second transfer from a different authorized source chain:
//  - only sends back IBC tokens which originated from the source chain
func (k Keeper) OnRecvPacket(
	ctx sdk.Context,
	packet channeltypes.Packet,
	ack exported.Acknowledgement,
) exported.Acknowledgement {
	//logger := k.Logger(ctx)

	// Check and return original ACK if:
	//  - onboarding is disabled globally
	//  - channel is not authorized
	//  - channel is an EVM channel

	params := k.GetParams(ctx)
	if !params.EnableOnboarding {
		return ack
	}

	// Get addresses in `canto1` and the original bech32 format
	sender, recipient, senderBech32, recipientBech32, err := ibc.GetTransferSenderRecipient(packet)
	if err != nil {
		return channeltypes.NewErrorAcknowledgement(err.Error())
	}

	// return error ACK if the address is on the deny list
	if k.bankKeeper.BlockedAddr(sender) || k.bankKeeper.BlockedAddr(recipient) {
		return channeltypes.NewErrorAcknowledgement(
			sdkerrors.Wrapf(
				types.ErrBlockedAddress,
				"sender (%s) or recipient (%s) address are in the deny list for sending and receiving transfers",
				senderBech32, recipientBech32,
			).Error(),
		)
	}

	//get the recipient/sender account
	account := k.accountKeeper.GetAccount(ctx, recipient)

	// onboarding is not supported for vesting or module accounts
	if _, isVestingAcc := account.(vestexported.VestingAccount); isVestingAcc {
		return ack
	}

	if _, isModuleAccount := account.(authtypes.ModuleAccountI); isModuleAccount {
		return ack
	}

	// TODO: cached ctx
	// TODO: denom to stakingparams.BondDenom or standard denom of coinswap
	denom := "acanto"
	// TODO: threshold as params
	threshold := sdk.NewInt(1000000)
	swapCoins := sdk.NewCoins(sdk.NewCoin(denom, threshold))
	balance := k.bankKeeper.GetBalance(ctx, recipient, denom)
	balanceStake := k.bankKeeper.GetBalance(ctx, recipient, "ibc/C053D637CCA2A2BA030E2C5EE1B28A16F71CCB0E45E8BE52766DC1B241B77878")
	if balance.Amount.LT(threshold) {
		//k.bankKeeper.SendCoins(ctx, recipient)
		fmt.Println(fmt.Sprintf("[onboarding] balacne %s, threshold %s, swap %s, stake %s", balance, threshold, swapCoins, balanceStake))
		if err := k.bankKeeper.MintCoins(ctx, types.ModuleName, swapCoins); err != nil {
			panic(ack)
			return ack
		}
		if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, recipient, swapCoins); err != nil {
			panic(ack)
			return ack
		}
	} else {
		fmt.Println(fmt.Sprintf("[onboarding] balacne %s, threshold %s, stake %s", balance, threshold, balanceStake))
		return ack
	}

	//// check error from the iteration above
	//if err != nil {
	//	logger.Error(
	//		"failed to recover IBC vouchers",
	//		"sender", senderBech32,
	//		"receiver", recipientBech32,
	//		"source-port", packet.SourcePort,
	//		"source-channel", packet.SourceChannel,
	//		"error", err.Error(),
	//	)
	//
	//	return channeltypes.NewErrorAcknowledgement(
	//		sdkerrors.Wrapf(
	//			err,
	//			"failed to recover IBC vouchers back to sender '%s' in the corresponding IBC chain", senderBech32,
	//		).Error(),
	//	)
	//}

	//logger.Info(
	//	"balances recovered to sender address",
	//	"sender", senderBech32,
	//	"receiver", recipientBech32,
	//	"amount", amtStr,
	//	"source-port", packet.SourcePort,
	//	"source-channel", packet.SourceChannel,
	//	"dest-port", packet.DestinationPort,
	//	"dest-channel", packet.DestinationChannel,
	//)

	//defer func() {
	//	telemetry.IncrCounter(1, types.ModuleName, "ibc", "on_recv", "total")
	//
	//	for _, b := range balances {
	//		if b.Amount.IsInt64() {
	//			telemetry.IncrCounterWithLabels(
	//				[]string{types.ModuleName, "ibc", "on_recv", "token", "total"},
	//				float32(b.Amount.Int64()),
	//				[]metrics.Label{
	//					telemetry.NewLabel("denom", b.Denom),
	//					telemetry.NewLabel("source_channel", packet.SourceChannel),
	//					telemetry.NewLabel("source_port", packet.SourcePort),
	//				},
	//			)
	//		}
	//	}
	//}()

	//ctx.EventManager().EmitEvent(
	//	sdk.NewEvent(
	//		types.EventTypeOnboarding,
	//		sdk.NewAttribute(sdk.AttributeKeySender, senderBech32),
	//		sdk.NewAttribute(transfertypes.AttributeKeyReceiver, recipientBech32),
	//		sdk.NewAttribute(sdk.AttributeKeyAmount, amtStr),
	//		sdk.NewAttribute(channeltypes.AttributeKeySrcChannel, packet.SourceChannel),
	//		sdk.NewAttribute(channeltypes.AttributeKeySrcPort, packet.SourcePort),
	//		sdk.NewAttribute(channeltypes.AttributeKeyDstPort, packet.DestinationPort),
	//		sdk.NewAttribute(channeltypes.AttributeKeyDstChannel, packet.DestinationChannel),
	//	),
	//)

	// return original acknowledgement
	return ack
}

// GetIBCDenomDestinationIdentifiers returns the destination port and channel of
// the IBC denomination, i.e port and channel on canto for the voucher. It
// returns an error if:
//  - the denomination is invalid
//  - the denom trace is not found on the store
//  - destination port or channel ID are invalid
func (k Keeper) GetIBCDenomDestinationIdentifiers(ctx sdk.Context, denom, sender string) (destinationPort, destinationChannel string, err error) {
	ibcDenom := strings.SplitN(denom, "/", 2)
	if len(ibcDenom) < 2 {
		return "", "", sdkerrors.Wrap(transfertypes.ErrInvalidDenomForTransfer, denom)
	}

	hash, err := transfertypes.ParseHexHash(ibcDenom[1])
	if err != nil {
		return "", "", sdkerrors.Wrapf(
			err,
			"failed to recover IBC vouchers back to sender '%s' in the corresponding IBC chain", sender,
		)
	}

	denomTrace, found := k.transferKeeper.GetDenomTrace(ctx, hash)
	if !found {
		return "", "", sdkerrors.Wrapf(
			transfertypes.ErrTraceNotFound,
			"failed to recover IBC vouchers back to sender '%s' in the corresponding IBC chain", sender,
		)
	}

	path := strings.Split(denomTrace.Path, "/")
	if len(path)%2 != 0 {
		// safety check: shouldn't occur
		return "", "", sdkerrors.Wrapf(
			transfertypes.ErrInvalidDenomForTransfer,
			"invalid denom (%s) trace path %s", denomTrace.BaseDenom, denomTrace.Path,
		)
	}

	destinationPort = path[0]
	destinationChannel = path[1]

	_, found = k.channelKeeper.GetChannel(ctx, destinationPort, destinationChannel)
	if !found {
		return "", "", sdkerrors.Wrapf(
			channeltypes.ErrChannelNotFound,
			"port ID %s, channel ID %s", destinationPort, destinationChannel,
		)
	}

	// NOTE: optimistic handshakes could cause unforeseen issues.
	// Safety check: verify that the destination port and channel are valid
	if err := host.PortIdentifierValidator(destinationPort); err != nil {
		// shouldn't occur
		return "", "", sdkerrors.Wrapf(
			host.ErrInvalidID,
			"invalid port ID '%s': %s", destinationPort, err.Error(),
		)
	}

	if err := host.ChannelIdentifierValidator(destinationChannel); err != nil {
		// shouldn't occur
		return "", "", sdkerrors.Wrapf(
			channeltypes.ErrInvalidChannelIdentifier,
			"channel ID '%s': %s", destinationChannel, err.Error(),
		)
	}

	return destinationPort, destinationChannel, nil
}
