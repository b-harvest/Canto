package keeper_test

import (
	"fmt"
	coinswaptypes "github.com/b-harvest/coinswap/modules/coinswap/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/evmos/ethermint/tests"
	"github.com/stretchr/testify/mock"
	"time"

	"github.com/Canto-Network/Canto/v6/testutil"

	transfertypes "github.com/cosmos/ibc-go/v3/modules/apps/transfer/types"
	clienttypes "github.com/cosmos/ibc-go/v3/modules/core/02-client/types"
	channeltypes "github.com/cosmos/ibc-go/v3/modules/core/04-channel/types"
	ibcgotesting "github.com/cosmos/ibc-go/v3/testing"
	ibcmock "github.com/cosmos/ibc-go/v3/testing/mock"

	"github.com/Canto-Network/Canto/v6/x/onboarding/keeper"
	"github.com/Canto-Network/Canto/v6/x/onboarding/types"
)

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
	cantoChannel := "channel-3"
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
		receiverAcantoAmount sdk.Coin
		expVoucher           sdk.Coins
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
			sdk.NewCoin("acanto", sdk.ZeroInt()),
			voucher,
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
			sdk.NewCoin("acanto", sdk.ZeroInt()),
			sdk.NewCoins(sdk.NewCoin(ibcUsdtDenom, transferAmount)),
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
			sdk.NewCoin("acanto", sdk.NewIntWithDecimal(4, 18)),
			sdk.NewCoins(sdk.NewCoin(ibcUsdtDenom, transferAmount)),
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
			sdk.NewCoin("acanto", sdk.ZeroInt()),
			sdk.NewCoins(sdk.NewCoin(ibcUsdcDenom, sdk.NewIntWithDecimal(1, 6))),
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
			sdk.NewCoin("acanto", sdk.NewIntWithDecimal(4, 18)),
			voucher,
		},
		{
			"fail - unauthorized source channel",
			func() {
				sourceChannel = "channel-100"
				transferAmount = sdk.NewIntWithDecimal(25, 6)
				transfer := transfertypes.NewFungibleTokenPacketData(denom, transferAmount.String(), secpAddrCosmos, secpAddrcanto)
				bz := transfertypes.ModuleCdc.MustMarshalJSON(&transfer)
				packet = channeltypes.NewPacket(bz, 100, transfertypes.PortID, sourceChannel, transfertypes.PortID, cantoChannel, timeoutHeight, 0)
			},
			true,
			false,
			sdk.NewCoin("acanto", sdk.ZeroInt()),
			voucher,
		},
		{
			"success - swap is successful",
			func() {
				sourceChannel = "channel-0"
				transferAmount = sdk.NewIntWithDecimal(25, 6)
				transfer := transfertypes.NewFungibleTokenPacketData(denom, transferAmount.String(), secpAddrCosmos, secpAddrcanto)
				bz := transfertypes.ModuleCdc.MustMarshalJSON(&transfer)
				packet = channeltypes.NewPacket(bz, 100, transfertypes.PortID, sourceChannel, transfertypes.PortID, cantoChannel, timeoutHeight, 0)
			},
			true,
			true,
			sdk.NewCoin("acanto", sdk.ZeroInt()),
			voucher,
		},

		{
			"success - swap is successful (acanto balance is positive but less than threshold)",
			func() {
				transferAmount = sdk.NewIntWithDecimal(25, 6)
				transfer := transfertypes.NewFungibleTokenPacketData(denom, transferAmount.String(), secpAddrCosmos, secpAddrcanto)
				bz := transfertypes.ModuleCdc.MustMarshalJSON(&transfer)
				packet = channeltypes.NewPacket(bz, 100, transfertypes.PortID, sourceChannel, transfertypes.PortID, cantoChannel, timeoutHeight, 0)
			},
			true,
			true,
			sdk.NewCoin("acanto", sdk.NewIntWithDecimal(3, 18)),
			voucher,
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
			sdk.NewCoin("acanto", sdk.NewIntWithDecimal(3, 18)),
			voucher,
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
			if tc.expOnboarding {
				suite.Require().True(cantoBalance.Equal(tc.receiverAcantoAmount.Add(sdk.NewCoin("acanto", params.AutoSwapThreshold))))
				suite.Require().True(voucherBalance.Amount.LT(transferAmount))
			} else {
				suite.Require().Equal(tc.expVoucher, sdk.NewCoins(voucherBalance))
				suite.Require().Equal(tc.receiverAcantoAmount, cantoBalance)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestGetIBCDenomDestinationIdentifiers() {
	address := sdk.AccAddress(tests.GenerateAddress().Bytes()).String()

	testCases := []struct {
		name                                      string
		denom                                     string
		malleate                                  func()
		expError                                  bool
		expDestinationPort, expDestinationChannel string
	}{
		{
			"invalid native denom",
			"acanto",
			func() {},
			true,
			"", "",
		},
		{
			"invalid IBC denom hash",
			"ibc/acanto",
			func() {},
			true,
			"", "",
		},
		{
			"denom trace not found",
			ibcAtomDenom,
			func() {},
			true,
			"", "",
		},
		{
			"channel not found",
			ibcAtomDenom,
			func() {
				denomTrace := transfertypes.DenomTrace{
					Path:      "transfer/channel-3",
					BaseDenom: "uatom",
				}
				suite.app.TransferKeeper.SetDenomTrace(suite.ctx, denomTrace)
			},
			true,
			"", "",
		},
		{
			"invalid destination port - insufficient length",
			"ibc/B9A49AA0AB0EB977D4EC627D7D9F747AF11BB1D74F430DE759CA37B22ECACF30", // denomTrace.Hash()
			func() {
				denomTrace := transfertypes.DenomTrace{
					Path:      "t/channel-3",
					BaseDenom: "uatom",
				}
				suite.app.TransferKeeper.SetDenomTrace(suite.ctx, denomTrace)

				channel := channeltypes.Channel{
					Counterparty: channeltypes.NewCounterparty("t", "channel-292"),
				}
				suite.app.IBCKeeper.ChannelKeeper.SetChannel(suite.ctx, "t", "channel-3", channel)
			},
			true,
			"", "",
		},
		{
			"invalid channel identifier - insufficient length",
			"ibc/5E3E083402F07599C795A7B75058EC3F13A8E666A8FEA2E51B6F3D93C755DFBC", // denomTrace.Hash()
			func() {
				denomTrace := transfertypes.DenomTrace{
					Path:      "transfer/c-3",
					BaseDenom: "uatom",
				}
				suite.app.TransferKeeper.SetDenomTrace(suite.ctx, denomTrace)

				channel := channeltypes.Channel{
					Counterparty: channeltypes.NewCounterparty("transfer", "channel-292"),
				}
				suite.app.IBCKeeper.ChannelKeeper.SetChannel(suite.ctx, "transfer", "c-3", channel)
			},
			true,
			"", "",
		},
		{
			"success - ATOM",
			ibcAtomDenom,
			func() {
				denomTrace := transfertypes.DenomTrace{
					Path:      "transfer/channel-3",
					BaseDenom: "uatom",
				}
				suite.app.TransferKeeper.SetDenomTrace(suite.ctx, denomTrace)

				channel := channeltypes.Channel{
					Counterparty: channeltypes.NewCounterparty("transfer", "channel-292"),
				}
				suite.app.IBCKeeper.ChannelKeeper.SetChannel(suite.ctx, "transfer", "channel-3", channel)
			},
			false,
			"transfer", "channel-3",
		},
		{
			"success - OSMO",
			ibcOsmoDenom,
			func() {
				denomTrace := transfertypes.DenomTrace{
					Path:      "transfer/channel-0",
					BaseDenom: "uosmo",
				}
				suite.app.TransferKeeper.SetDenomTrace(suite.ctx, denomTrace)

				channel := channeltypes.Channel{
					Counterparty: channeltypes.NewCounterparty("transfer", "channel-204"),
				}
				suite.app.IBCKeeper.ChannelKeeper.SetChannel(suite.ctx, "transfer", "channel-0", channel)
			},
			false,
			"transfer", "channel-0",
		},
		{
			"success - ibcATOM (via Osmosis)",
			"ibc/6CDD4663F2F09CD62285E2D45891FC149A3568E316CE3EBBE201A71A78A69388",
			func() {
				denomTrace := transfertypes.DenomTrace{
					Path:      "transfer/channel-0/transfer/channel-0",
					BaseDenom: "uatom",
				}

				suite.app.TransferKeeper.SetDenomTrace(suite.ctx, denomTrace)

				channel := channeltypes.Channel{
					Counterparty: channeltypes.NewCounterparty("transfer", "channel-204"),
				}
				suite.app.IBCKeeper.ChannelKeeper.SetChannel(suite.ctx, "transfer", "channel-0", channel)
			},
			false,
			"transfer", "channel-0",
		},
	}
	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			suite.SetupTest() // reset

			tc.malleate()

			destinationPort, destinationChannel, err := suite.app.OnboardingKeeper.GetIBCDenomDestinationIdentifiers(suite.ctx, tc.denom, address)
			if tc.expError {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
				suite.Require().Equal(tc.expDestinationPort, destinationPort)
				suite.Require().Equal(tc.expDestinationChannel, destinationChannel)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestOnRecvPacketFailTransfer() {
	// secp256k1 account
	secpPk := secp256k1.GenPrivKey()
	secpAddr := sdk.AccAddress(secpPk.PubKey().Address())
	secpAddrcanto := secpAddr.String()
	secpAddrCosmos := sdk.MustBech32ifyAddressBytes(sdk.Bech32MainPrefix, secpAddr)

	// Setup Cosmos <=> canto IBC relayer
	denom := "uatom"
	sourceChannel := "channel-292"
	// cantoChannel := claimstypes.DefaultAuthorizedChannels[1]

	cantoChannel := "channel-3"
	path := fmt.Sprintf("%s/%s", transfertypes.PortID, cantoChannel)

	var mockTransferKeeper *MockTransferKeeper
	expAck := ibcmock.MockAcknowledgement
	testCases := []struct {
		name     string
		malleate func()
	}{
		{
			"Fail to retrieve ibc denom trace",
			func() {
				mockTransferKeeper.On("GetDenomTrace", mock.Anything, mock.Anything).Return(transfertypes.DenomTrace{}, false)
				mockTransferKeeper.On("SendTransfer", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
			},
		},
		{
			"invalid ibc denom trace",
			func() {
				// Set Denom Trace
				denomTrace := transfertypes.DenomTrace{
					Path:      "badpath",
					BaseDenom: denom,
				}
				suite.app.TransferKeeper.SetDenomTrace(suite.ctx, denomTrace)
				mockTransferKeeper.On("GetDenomTrace", mock.Anything, mock.Anything).Return(denomTrace, true)
				mockTransferKeeper.On("SendTransfer", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
			},
		},

		{
			"Fail to send transfer",
			func() {
				// Set Denom Trace
				denomTrace := transfertypes.DenomTrace{
					Path:      path,
					BaseDenom: denom,
				}
				suite.app.TransferKeeper.SetDenomTrace(suite.ctx, denomTrace)
				mockTransferKeeper.On("GetDenomTrace", mock.Anything, mock.Anything).Return(denomTrace, true)
				mockTransferKeeper.On("SendTransfer", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(fmt.Errorf("Fail to transfer"))
			},
		},
	}
	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			suite.SetupTest() // reset

			// Enable Onboarding
			params := suite.app.OnboardingKeeper.GetParams(suite.ctx)
			params.EnableOnboarding = true
			suite.app.OnboardingKeeper.SetParams(suite.ctx, params)

			transfer := transfertypes.NewFungibleTokenPacketData(denom, "100", secpAddrCosmos, secpAddrcanto)
			packet := channeltypes.NewPacket(transfer.GetBytes(), 100, transfertypes.PortID, sourceChannel, transfertypes.PortID, cantoChannel, timeoutHeight, 0)

			mockTransferKeeper = &MockTransferKeeper{
				Keeper: suite.app.BankKeeper,
			}

			tc.malleate()

			sp, found := suite.app.ParamsKeeper.GetSubspace(types.ModuleName)
			suite.Require().True(found)
			suite.app.OnboardingKeeper = keeper.NewKeeper(sp, suite.app.AccountKeeper, suite.app.BankKeeper, suite.app.IBCKeeper.ChannelKeeper, mockTransferKeeper, suite.app.CoinswapKeeper, suite.app.Erc20Keeper)

			// Fund receiver account with canto
			coins := sdk.NewCoins(
				sdk.NewCoin("acanto", sdk.NewInt(1000)),
				sdk.NewCoin(ibcAtomDenom, sdk.NewInt(1000)),
			)
			testutil.FundAccount(suite.app.BankKeeper, suite.ctx, secpAddr, coins)

			// Perform IBC callback
			ack := suite.app.OnboardingKeeper.OnRecvPacket(suite.ctx, packet, expAck)
			// Onboarding should Fail
			suite.Require().Equal(expAck, ack)
		})
	}
}
