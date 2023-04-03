package keeper

import (
	errorsmod "cosmossdk.io/errors"
	"fmt"
	"github.com/Canto-Network/Canto/v6/contracts"
	erc20types "github.com/Canto-Network/Canto/v6/x/erc20/types"
	coinswaptypes "github.com/b-harvest/coinswap/modules/coinswap/types"
	"github.com/ethereum/go-ethereum/common"
	"strings"
	"time"

	"github.com/Canto-Network/Canto/v6/ibc"
	"github.com/Canto-Network/Canto/v6/x/onboarding/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	vestexported "github.com/cosmos/cosmos-sdk/x/auth/vesting/exported"
	transfertypes "github.com/cosmos/ibc-go/v3/modules/apps/transfer/types"
	channeltypes "github.com/cosmos/ibc-go/v3/modules/core/04-channel/types"
	host "github.com/cosmos/ibc-go/v3/modules/core/24-host"
	"github.com/cosmos/ibc-go/v3/modules/core/exported"
)

// OnRecvPacket performs an IBC receive callback. It returns the tokens that
// users transferred to their Cosmos secp256k1 address instead of the Ethereum
// ethsecp256k1 address. The expected behavior is as follows:
//
// First transfer from authorized source chain:
//   - sends back IBC tokens which originated from the source chain
//   - sends over all canto native tokens
//
// Second transfer from a different authorized source chain:
//   - only sends back IBC tokens which originated from the source chain
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

	// check source channel is in the whitelist channels
	var found bool
	for _, s := range params.WhitelistedChannels {
		if s == packet.DestinationChannel {
			found = true
		}
	}

	if !found {
		return ack
	}

	fmt.Println(fmt.Sprintf("[onboarding] source channel: %s", packet.SourceChannel))
	fmt.Println(fmt.Sprintf("[onboarding] destination channel: %s", packet.DestinationChannel))

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
	standardDenom := k.coinswapKeeper.GetStandardDenom(ctx)
	fmt.Println(fmt.Sprintf("[onboarding] denom %s", standardDenom))

	var data transfertypes.FungibleTokenPacketData
	if err = transfertypes.ModuleCdc.UnmarshalJSON(packet.GetData(), &data); err != nil {
		// NOTE: shouldn't happen as the packet has already
		// been decoded on ICS20 transfer logic
		err = errorsmod.Wrapf(types.ErrInvalidType, "cannot unmarshal ICS-20 transfer packet data")
		return channeltypes.NewErrorAcknowledgement(err.Error())
	}
	// parse the transferred denom
	transferredCoin := ibc.GetReceivedCoin(
		packet.SourcePort, packet.SourceChannel,
		packet.DestinationPort, packet.DestinationChannel,
		data.Denom, data.Amount,
	)
	fmt.Println(fmt.Sprintf("[onboarding] coin %s", transferredCoin))

	threshold := k.GetParams(ctx).AutoSwapThreshold
	swapCoins := sdk.NewCoin(standardDenom, threshold)
	standardCoinBalance := k.bankKeeper.GetBalance(ctx, recipient, standardDenom)
	transferredCoinBalance := k.bankKeeper.GetBalance(ctx, recipient, transferredCoin.Denom)

	if standardCoinBalance.Amount.LT(threshold) {
		fmt.Println(fmt.Sprintf("[onboarding] before swap balacne %s, threshold %s, swap %s, stake %s", standardCoinBalance, threshold, swapCoins, transferredCoinBalance))

		swapDuration := k.GetParams(ctx).AutoSwapDuration
		msg := coinswaptypes.MsgSwapOrder{
			coinswaptypes.Input{Coin: transferredCoin, Address: recipient.String()},
			coinswaptypes.Output{Coin: swapCoins, Address: recipient.String()},
			time.Now().Add(swapDuration).Unix(),
			true,
		}

		if err = k.coinswapKeeper.Swap(ctx, &msg); err != nil {
			fmt.Println(fmt.Sprintf("[onboarding] swap error %s", err))
			return ack
		}

		transferredCoinBalance = k.bankKeeper.GetBalance(ctx, recipient, transferredCoin.Denom)

		//convert coins to ERC20 token
		pairID := k.erc20Keeper.GetTokenPairID(ctx, transferredCoin.Denom)
		if len(pairID) == 0 {
			// short-circuit: if the denom is not registered, conversion will fail
			// so we can continue with the rest of the stack
			return ack
		}

		pair, _ := k.erc20Keeper.GetTokenPair(ctx, pairID)
		if !pair.Enabled {
			// no-op: continue with the rest of the stack without conversion
			return ack
		}

		// Build MsgConvertCoin, from recipient to recipient since IBC transfer already occurred
		convertMsg := erc20types.NewMsgConvertCoin(transferredCoinBalance, common.BytesToAddress(recipient.Bytes()), recipient)

		// NOTE: we don't use ValidateBasic the msg since we've already validated
		// the ICS20 packet data

		// Use MsgConvertCoin to convert the Cosmos Coin to an ERC20
		if _, err = k.erc20Keeper.ConvertCoin(sdk.WrapSDKContext(ctx), convertMsg); err != nil {
			return ack
		}

		abi := contracts.ERC20BurnableContract.ABI
		ercBalance := k.erc20Keeper.BalanceOf(ctx, abi, pair.GetERC20Contract(), common.BytesToAddress(recipient.Bytes()))
		res, err := k.erc20Keeper.CallEVM(ctx, abi, common.BytesToAddress(recipient.Bytes()), pair.GetERC20Contract(), false, "symbol")
		if err != nil {
			return ack
		}
		var symbolRes erc20types.ERC20StringResponse
		if err := abi.UnpackIntoInterface(&symbolRes, "symbol", res.Ret); err != nil {
			return ack
		}

		standardCoinBalance = k.bankKeeper.GetBalance(ctx, recipient, standardDenom)
		transferredCoinBalance = k.bankKeeper.GetBalance(ctx, recipient, transferredCoin.Denom)
		fmt.Println(fmt.Sprintf("[onboarding] after swap balacne %s, threshold %s, swap %s, stake %s", standardCoinBalance, threshold, swapCoins, transferredCoinBalance))
		fmt.Println(fmt.Sprintf("[onboarding] erc20 token balance %s %s", ercBalance, symbolRes.Value))

	} else {
		fmt.Println(fmt.Sprintf("[onboarding] balacne %s, threshold %s, stake %s", standardCoinBalance, threshold, transferredCoinBalance))
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
//   - the denomination is invalid
//   - the denom trace is not found on the store
//   - destination port or channel ID are invalid
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
