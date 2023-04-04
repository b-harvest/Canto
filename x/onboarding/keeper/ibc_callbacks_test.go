package keeper_test

import (
	"fmt"
	"github.com/Canto-Network/Canto/v6/contracts"
	coinswaptypes "github.com/b-harvest/coinswap/modules/coinswap/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/mock"
	"golang.org/x/exp/slices"
	"time"

	"github.com/Canto-Network/Canto/v6/testutil"

	transfertypes "github.com/cosmos/ibc-go/v3/modules/apps/transfer/types"
	clienttypes "github.com/cosmos/ibc-go/v3/modules/core/02-client/types"
	channeltypes "github.com/cosmos/ibc-go/v3/modules/core/04-channel/types"
	ibcgotesting "github.com/cosmos/ibc-go/v3/testing"
	ibcmock "github.com/cosmos/ibc-go/v3/testing/mock"

	erc20types "github.com/Canto-Network/Canto/v6/x/erc20/types"
	"github.com/Canto-Network/Canto/v6/x/onboarding/keeper"
	"github.com/Canto-Network/Canto/v6/x/onboarding/types"
)

var (
	metadataIbc = banktypes.Metadata{
		Description: "USDC IBC voucher (channel 0)",
		Base:        ibcUsdcDenom,
		// NOTE: Denom units MUST be increasing
		DenomUnits: []*banktypes.DenomUnit{
			{
				Denom:    ibcUsdcDenom,
				Exponent: 0,
			},
		},
		Name:    "USDC channel-0",
		Symbol:  "ibcUSDC-0",
		Display: ibcUsdcDenom,
	}
)

func (suite *KeeperTestSuite) setupRegisterCoin(metadata banktypes.Metadata) *erc20types.TokenPair {
	pair, err := suite.app.Erc20Keeper.RegisterCoin(suite.ctx, metadata)
	suite.Require().NoError(err)
	return pair
}

func (s *KeeperTestSuite) FindEvent(events []sdk.Event, name string) sdk.Event {
	index := slices.IndexFunc(events, func(e sdk.Event) bool { return e.Type == name })
	if index == -1 {
		return sdk.Event{}
	}
	return events[index]
}

func (s *KeeperTestSuite) ExtractAttributes(event sdk.Event) map[string]string {
	attrs := make(map[string]string)
	if event.Attributes == nil {
		return attrs
	}
	for _, a := range event.Attributes {
		attrs[string(a.Key)] = string(a.Value)
	}
	return attrs
}

func (suite *KeeperTestSuite) TestOnRecvPacket() {
	// secp256k1 account
	secpPk := secp256k1.GenPrivKey()
	secpAddr := sdk.AccAddress(secpPk.PubKey().Address())
	secpAddrcanto := secpAddr.String()
	secpAddrCosmos := sdk.MustBech32ifyAddressBytes(sdk.Bech32MainPrefix, secpAddr)

	// Setup Cosmos <=> canto IBC relayer
	denom := "uUSDC"
	ibcDenom := ibcUsdcDenom
	transferAmount := sdk.NewIntWithDecimal(25, 6)
	sourceChannel := "channel-0"
	cantoChannel := "channel-0"
	path := fmt.Sprintf("%s/%s", transfertypes.PortID, cantoChannel)

	timeoutHeight := clienttypes.NewHeight(0, 100)
	disabledTimeoutTimestamp := uint64(0)
	mockPacket := channeltypes.NewPacket(ibcgotesting.MockPacketData, 1, transfertypes.PortID, sourceChannel, transfertypes.PortID, cantoChannel, timeoutHeight, disabledTimeoutTimestamp)
	packet := mockPacket
	expAck := ibcmock.MockAcknowledgement

	voucher := sdk.NewCoins(
		sdk.NewCoin(ibcUsdcDenom, transferAmount),
	)

	testCases := []struct {
		name                 string
		malleate             func()
		ackSuccess           bool
		expOnboarding        bool
		enableConvert        bool
		expSwapAmount        sdk.Int
		expConvertAmount     sdk.Int
		receiverAcantoAmount sdk.Coin
		expVoucher           sdk.Coins
		expErc20Balance      int64
	}{
		{
			"continue - params disabled",
			func() {
				params := suite.app.OnboardingKeeper.GetParams(suite.ctx)
				params.EnableOnboarding = false
				suite.app.OnboardingKeeper.SetParams(suite.ctx, params)
			},
			true,
			false,
			true,
			sdk.ZeroInt(),
			sdk.ZeroInt(),
			sdk.NewCoin("acanto", sdk.ZeroInt()),
			voucher,
			0,
		},
		{
			"fail - no liquidity pool exists",
			func() {

				denom = "uUSDT"
				ibcDenom = ibcUsdtDenom

				transfer := transfertypes.NewFungibleTokenPacketData(denom, transferAmount.String(), secpAddrCosmos, secpAddrcanto)
				bz := transfertypes.ModuleCdc.MustMarshalJSON(&transfer)
				packet = channeltypes.NewPacket(bz, 100, transfertypes.PortID, sourceChannel, transfertypes.PortID, cantoChannel, timeoutHeight, 0)
			},
			true,
			false,
			true,
			sdk.ZeroInt(),
			sdk.ZeroInt(),
			sdk.NewCoin("acanto", sdk.ZeroInt()),
			sdk.NewCoins(sdk.NewCoin(ibcUsdtDenom, transferAmount)),
			0,
		},
		{
			"continue - no liquidity pool exists but acanto balance is already bigger than threshold",
			func() {

				denom = "uUSDT"
				ibcDenom = ibcUsdtDenom

				transfer := transfertypes.NewFungibleTokenPacketData(denom, transferAmount.String(), secpAddrCosmos, secpAddrcanto)
				bz := transfertypes.ModuleCdc.MustMarshalJSON(&transfer)
				packet = channeltypes.NewPacket(bz, 100, transfertypes.PortID, sourceChannel, transfertypes.PortID, cantoChannel, timeoutHeight, 0)
			},
			true,
			false,
			true,
			sdk.ZeroInt(),
			sdk.ZeroInt(),
			sdk.NewCoin("acanto", sdk.NewIntWithDecimal(4, 18)),
			sdk.NewCoins(sdk.NewCoin(ibcUsdtDenom, transferAmount)),
			0,
		},
		{
			"fail - not enough ibc coin to swap threshold",
			func() {
				denom = "uUSDC"
				ibcDenom = ibcUsdcDenom

				transferAmount = sdk.NewIntWithDecimal(1, 6)
				transfer := transfertypes.NewFungibleTokenPacketData(denom, transferAmount.String(), secpAddrCosmos, secpAddrcanto)
				bz := transfertypes.ModuleCdc.MustMarshalJSON(&transfer)
				packet = channeltypes.NewPacket(bz, 100, transfertypes.PortID, sourceChannel, transfertypes.PortID, cantoChannel, timeoutHeight, 0)
			},
			true,
			false,
			true,
			sdk.ZeroInt(),
			sdk.ZeroInt(),
			sdk.NewCoin("acanto", sdk.ZeroInt()),
			sdk.NewCoins(sdk.NewCoin(ibcUsdcDenom, sdk.NewIntWithDecimal(1, 6))),
			0,
		},
		{
			"continue - acanto balance is already bigger than threshold",
			func() {
				transferAmount = sdk.NewIntWithDecimal(25, 6)
				transfer := transfertypes.NewFungibleTokenPacketData(denom, transferAmount.String(), secpAddrCosmos, secpAddrcanto)
				bz := transfertypes.ModuleCdc.MustMarshalJSON(&transfer)
				packet = channeltypes.NewPacket(bz, 100, transfertypes.PortID, sourceChannel, transfertypes.PortID, cantoChannel, timeoutHeight, 0)
			},
			true,
			false,
			true,
			sdk.ZeroInt(),
			sdk.ZeroInt(),
			sdk.NewCoin("acanto", sdk.NewIntWithDecimal(4, 18)),
			voucher,
			0,
		},
		{
			"fail - unauthorized  channel",
			func() {
				cantoChannel = "channel-100"
				transferAmount = sdk.NewIntWithDecimal(25, 6)
				transfer := transfertypes.NewFungibleTokenPacketData(denom, transferAmount.String(), secpAddrCosmos, secpAddrcanto)
				bz := transfertypes.ModuleCdc.MustMarshalJSON(&transfer)
				packet = channeltypes.NewPacket(bz, 100, transfertypes.PortID, sourceChannel, transfertypes.PortID, cantoChannel, timeoutHeight, 0)
			},
			true,
			false,
			true,
			sdk.ZeroInt(),
			sdk.ZeroInt(),
			sdk.NewCoin("acanto", sdk.ZeroInt()),
			voucher,
			0,
		},
		{
			"success - swap and erc20 conversion are successful",
			func() {
				cantoChannel = "channel-0"
				transferAmount = sdk.NewIntWithDecimal(25, 6)
				transfer := transfertypes.NewFungibleTokenPacketData(denom, transferAmount.String(), secpAddrCosmos, secpAddrcanto)
				bz := transfertypes.ModuleCdc.MustMarshalJSON(&transfer)
				packet = channeltypes.NewPacket(bz, 100, transfertypes.PortID, sourceChannel, transfertypes.PortID, cantoChannel, timeoutHeight, 0)
			},
			true,
			true,
			true,
			sdk.NewInt(4001601),
			sdk.NewInt(20998399),
			sdk.NewCoin("acanto", sdk.ZeroInt()),
			sdk.NewCoins(sdk.NewCoin(ibcUsdcDenom, sdk.ZeroInt())),
			20998399,
		},
		{
			"success - swap is successful but erc20 conversion is not done",
			func() {
				sourceChannel = "channel-0"
				transferAmount = sdk.NewIntWithDecimal(25, 6)
				transfer := transfertypes.NewFungibleTokenPacketData(denom, transferAmount.String(), secpAddrCosmos, secpAddrcanto)
				bz := transfertypes.ModuleCdc.MustMarshalJSON(&transfer)
				packet = channeltypes.NewPacket(bz, 100, transfertypes.PortID, sourceChannel, transfertypes.PortID, cantoChannel, timeoutHeight, 0)
			},
			true,
			true,
			false,
			sdk.NewInt(4001601),
			sdk.ZeroInt(),
			sdk.NewCoin("acanto", sdk.ZeroInt()),
			sdk.NewCoins(sdk.NewCoin(ibcUsdcDenom, sdk.NewInt(20998399))),
			0,
		},
		{
			"success - swap and erc20 conversion are successful (acanto balance is positive but less than threshold)",
			func() {
				transferAmount = sdk.NewIntWithDecimal(25, 6)
				transfer := transfertypes.NewFungibleTokenPacketData(denom, transferAmount.String(), secpAddrCosmos, secpAddrcanto)
				bz := transfertypes.ModuleCdc.MustMarshalJSON(&transfer)
				packet = channeltypes.NewPacket(bz, 100, transfertypes.PortID, sourceChannel, transfertypes.PortID, cantoChannel, timeoutHeight, 0)
			},
			true,
			true,
			true,
			sdk.NewInt(4001601),
			sdk.NewInt(20998399),
			sdk.NewCoin("acanto", sdk.NewIntWithDecimal(3, 18)),
			sdk.NewCoins(sdk.NewCoin(ibcUsdcDenom, sdk.ZeroInt())),
			20998399,
		},
		{
			"success - swap is successful but erc20 conversion is not done (acanto balance is positive but less than threshold)",
			func() {
				transferAmount = sdk.NewIntWithDecimal(25, 6)
				transfer := transfertypes.NewFungibleTokenPacketData(denom, transferAmount.String(), secpAddrCosmos, secpAddrcanto)
				bz := transfertypes.ModuleCdc.MustMarshalJSON(&transfer)
				packet = channeltypes.NewPacket(bz, 100, transfertypes.PortID, sourceChannel, transfertypes.PortID, cantoChannel, timeoutHeight, 0)
			},
			true,
			true,
			false,
			sdk.NewInt(4001601),
			sdk.ZeroInt(),
			sdk.NewCoin("acanto", sdk.NewIntWithDecimal(3, 18)),
			sdk.NewCoins(sdk.NewCoin(ibcUsdcDenom, sdk.NewInt(20998399))),
			0,
		},

		{
			"fail - required ibc token to swap exceeds max swap amount limit",
			func() {
				coinswapParams := suite.app.CoinswapKeeper.GetParams(suite.ctx)
				coinswapParams.MaxSwapAmount = sdk.NewCoins(sdk.NewCoin(ibcUsdcDenom, sdk.NewIntWithDecimal(5, 5)))
				suite.app.CoinswapKeeper.SetParams(suite.ctx, coinswapParams)

				transferAmount = sdk.NewIntWithDecimal(25, 6)
				transfer := transfertypes.NewFungibleTokenPacketData(denom, transferAmount.String(), secpAddrCosmos, secpAddrcanto)
				bz := transfertypes.ModuleCdc.MustMarshalJSON(&transfer)
				packet = channeltypes.NewPacket(bz, 100, transfertypes.PortID, sourceChannel, transfertypes.PortID, cantoChannel, timeoutHeight, 0)
			},
			true,
			false,
			true,
			sdk.ZeroInt(),
			sdk.ZeroInt(),
			sdk.NewCoin("acanto", sdk.NewIntWithDecimal(3, 18)),
			voucher,
			0,
		},
	}
	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			suite.SetupTest() // reset

			// Setup liquidity pool (acanto, uUSDC)
			coinswapParams := suite.app.CoinswapKeeper.GetParams(suite.ctx)
			coinswapParams.MaxSwapAmount = sdk.NewCoins(sdk.NewCoin(ibcUsdcDenom, sdk.NewIntWithDecimal(10, 6)))
			suite.app.CoinswapKeeper.SetParams(suite.ctx, coinswapParams)

			testutil.FundAccount(suite.app.BankKeeper, suite.ctx, secpAddr, sdk.NewCoins(sdk.NewCoin("acanto", sdk.NewIntWithDecimal(10000, 18)), sdk.NewCoin(ibcUsdcDenom, sdk.NewIntWithDecimal(10000, 6))))

			msgAddLiquidity := coinswaptypes.MsgAddLiquidity{
				MaxToken:         sdk.NewCoin(ibcUsdcDenom, sdk.NewIntWithDecimal(10000, 6)),
				ExactStandardAmt: sdk.NewIntWithDecimal(10000, 18),
				MinLiquidity:     sdk.NewInt(1),
				Deadline:         time.Now().Add(time.Minute * 10).Unix(),
				Sender:           secpAddr.String(),
			}

			suite.app.CoinswapKeeper.AddLiquidity(suite.ctx, &msgAddLiquidity)

			// Enable Onboarding
			params := suite.app.OnboardingKeeper.GetParams(suite.ctx)
			params.EnableOnboarding = true
			params.WhitelistedChannels = []string{"channel-0"}
			suite.app.OnboardingKeeper.SetParams(suite.ctx, params)

			tc.malleate()

			// Set Denom Trace
			denomTrace := transfertypes.DenomTrace{
				Path:      path,
				BaseDenom: denom,
			}
			suite.app.TransferKeeper.SetDenomTrace(suite.ctx, denomTrace)

			// Set Cosmos Channel
			channel := channeltypes.Channel{
				State:          channeltypes.INIT,
				Ordering:       channeltypes.UNORDERED,
				Counterparty:   channeltypes.NewCounterparty(transfertypes.PortID, sourceChannel),
				ConnectionHops: []string{sourceChannel},
			}
			suite.app.IBCKeeper.ChannelKeeper.SetChannel(suite.ctx, transfertypes.PortID, cantoChannel, channel)

			// Set Next Sequence Send
			suite.app.IBCKeeper.ChannelKeeper.SetNextSequenceSend(suite.ctx, transfertypes.PortID, cantoChannel, 1)

			// Mock the Transferkeeper to always return nil on SendTransfer(), as this
			// method requires a successfull handshake with the counterparty chain.
			// This, however, exceeds the requirements of the unit tests.
			mockTransferKeeper := &MockTransferKeeper{
				Keeper: suite.app.BankKeeper,
			}

			mockTransferKeeper.On("GetDenomTrace", mock.Anything, mock.Anything).Return(denomTrace, true)
			mockTransferKeeper.On("SendTransfer", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

			sp, found := suite.app.ParamsKeeper.GetSubspace(types.ModuleName)
			suite.Require().True(found)
			suite.app.OnboardingKeeper = keeper.NewKeeper(sp, suite.app.AccountKeeper, suite.app.BankKeeper, suite.app.IBCKeeper.ChannelKeeper, mockTransferKeeper, suite.app.CoinswapKeeper, suite.app.Erc20Keeper)

			// Fund receiver account with canto, ERC20 coins and IBC vouchers
			testutil.FundAccount(suite.app.BankKeeper, suite.ctx, secpAddr, sdk.NewCoins(sdk.NewCoin(ibcDenom, transferAmount), tc.receiverAcantoAmount))

			pair := suite.setupRegisterCoin(metadataIbc)
			suite.Require().NotNil(pair)

			if !tc.enableConvert {
				pair.Enabled = false
				suite.app.Erc20Keeper.SetTokenPair(suite.ctx, *pair)
			}

			// Perform IBC callback
			ack := suite.app.OnboardingKeeper.OnRecvPacket(suite.ctx, packet, expAck)

			// Check acknowledgement
			if tc.ackSuccess {
				suite.Require().True(ack.Success(), string(ack.Acknowledgement()))
				suite.Require().Equal(expAck, ack)
			} else {
				suite.Require().False(ack.Success(), string(ack.Acknowledgement()))
			}

			// Check onboarding
			cantoBalance := suite.app.BankKeeper.GetBalance(suite.ctx, secpAddr, "acanto")
			voucherBalance := suite.app.BankKeeper.GetBalance(suite.ctx, secpAddr, ibcDenom)
			// Check ERC20 balances
			erc20balance := suite.app.Erc20Keeper.BalanceOf(suite.ctx, contracts.ERC20MinterBurnerDecimalsContract.ABI, pair.GetERC20Contract(), common.BytesToAddress(secpAddr.Bytes()))

			if tc.expOnboarding {
				suite.Require().True(cantoBalance.Equal(tc.receiverAcantoAmount.Add(sdk.NewCoin("acanto", params.AutoSwapThreshold))))
			} else {
				suite.Require().Equal(tc.expVoucher, sdk.NewCoins(voucherBalance))
			}
			suite.Require().Equal(tc.expVoucher, sdk.NewCoins(voucherBalance))
			suite.Require().Equal(tc.expErc20Balance, erc20balance.Int64())

			events := suite.ctx.EventManager().Events()

			attrs := suite.ExtractAttributes(suite.FindEvent(events, "swap"))
			if tc.expSwapAmount.IsPositive() {
				suite.Require().Equal(tc.expSwapAmount.String(), attrs["amount"])
			} else {
				suite.Require().Equal(0, len(attrs))
			}

			attrs = suite.ExtractAttributes(suite.FindEvent(events, "convert_coin"))
			if tc.enableConvert && tc.expConvertAmount.IsPositive() {
				suite.Require().Equal(tc.expConvertAmount.String(), attrs["amount"])
			} else {
				suite.Require().Equal(0, len(attrs))
			}
		})
	}
}
