package onboarding_test

import (
	"fmt"
	coinswaptypes "github.com/b-harvest/coinswap/modules/coinswap/types"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"

	ibctesting "github.com/Canto-Network/Canto/v6/testing"
	"github.com/cosmos/ibc-go/v3/modules/apps/transfer/types"
	clienttypes "github.com/cosmos/ibc-go/v3/modules/core/02-client/types"
)

type TransferTestSuite struct {
	suite.Suite

	coordinator *ibctesting.Coordinator

	// testing chains used for convenience and readability
	chainA *ibctesting.TestChain
	chainB *ibctesting.TestChain
	chainC *ibctesting.TestChain
}

func (suite *TransferTestSuite) SetupTest() {
	suite.coordinator = ibctesting.NewCoordinator(suite.T(), 3)
	suite.chainA = suite.coordinator.GetChain(ibctesting.GetChainID(1))
	suite.chainB = suite.coordinator.GetChain(ibctesting.GetChainID(2))
	suite.chainC = suite.coordinator.GetChain(ibctesting.GetChainID(3))
}

func NewTransferPath(chainA, chainB *ibctesting.TestChain) *ibctesting.Path {
	path := ibctesting.NewPath(chainA, chainB)
	path.EndpointA.ChannelConfig.PortID = ibctesting.TransferPort
	path.EndpointB.ChannelConfig.PortID = ibctesting.TransferPort
	path.EndpointA.ChannelConfig.Version = types.Version
	path.EndpointB.ChannelConfig.Version = types.Version

	return path
}

// constructs a send from chainA to chainB on the established channel/connection
// and sends the same coin back from chainB to chainA.
func (suite *TransferTestSuite) TestHandleMsgTransfer() {
	// setup between chainA and chainB
	path := NewTransferPath(suite.chainA, suite.chainB)
	suite.coordinator.Setup(path)
	coinswapKeeper := suite.chainB.App.GetCoinswapKeeper()
	params := coinswapKeeper.GetParams(suite.chainB.GetContext())

	// register ibc denoms (set params)
	params.MaxSwapAmount = sdk.NewCoins(sdk.NewCoin("ibc/C053D637CCA2A2BA030E2C5EE1B28A16F71CCB0E45E8BE52766DC1B241B77878", sdk.NewInt(10000000)))
	coinswapKeeper.SetParams(suite.chainB.GetContext(), params)

	middlewareParams := suite.chainB.App.GetOnboardingKeeper().GetParams(suite.chainB.GetContext())
	middlewareParams.AutoSwapThreshold = sdk.NewInt(4000000)
	suite.chainB.App.GetOnboardingKeeper().SetParams(suite.chainB.GetContext(), middlewareParams)

	// Pool creation
	msgAddLiquidity := coinswaptypes.MsgAddLiquidity{
		MaxToken:         sdk.NewCoin("ibc/C053D637CCA2A2BA030E2C5EE1B28A16F71CCB0E45E8BE52766DC1B241B77878", sdk.NewInt(10000000000)),
		ExactStandardAmt: sdk.NewInt(10000000000),
		MinLiquidity:     sdk.NewInt(1),
		Deadline:         time.Now().Add(time.Minute * 10).Unix(),
		Sender:           suite.chainB.SenderAccount.GetAddress().String(),
	}

	_, err := coinswapKeeper.AddLiquidity(suite.chainB.GetContext(), &msgAddLiquidity)
	if err != nil {
		fmt.Println(err)
	}

	//	originalBalance := suite.chainA.GetSimApp().BankKeeper.GetBalance(suite.chainA.GetContext(), suite.chainA.SenderAccount.GetAddress(), sdk.DefaultBondDenom)
	timeoutHeight := clienttypes.NewHeight(10, 100)

	amount, ok := sdk.NewIntFromString("9223372036854775808") // 2^63 (one above int64)
	suite.Require().True(ok)
	coinToSendToB := sdk.NewCoin(sdk.DefaultBondDenom, amount)

	// send from chainA to chainB
	msg := types.NewMsgTransfer(path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID, coinToSendToB, suite.chainA.SenderAccount.GetAddress().String(), suite.chainB.SenderAccount.GetAddress().String(), timeoutHeight, 0)
	res, err := suite.chainA.SendMsgs(msg)
	suite.Require().NoError(err) // message committed

	packet, err := ibctesting.ParsePacketFromEvents(res.GetEvents())
	suite.Require().NoError(err)

	voucherDenomTrace := types.ParseDenomTrace(types.GetPrefixedDenom(packet.GetDestPort(), packet.GetDestChannel(), sdk.DefaultBondDenom))

	balanceBefore := suite.chainB.GetSimApp().BankKeeper.GetBalance(suite.chainB.GetContext(), suite.chainB.SenderAccount.GetAddress(), voucherDenomTrace.IBCDenom())
	balanceCantoBefore := suite.chainB.GetSimApp().BankKeeper.GetBalance(suite.chainB.GetContext(), suite.chainB.SenderAccount.GetAddress(), "acanto")

	// relay send
	err = path.RelayPacket(packet)
	suite.Require().NoError(err) // relay committed

	// check that voucher exists on chain B
	balance := suite.chainB.GetSimApp().BankKeeper.GetBalance(suite.chainB.GetContext(), suite.chainB.SenderAccount.GetAddress(), voucherDenomTrace.IBCDenom())
	balanceCanto := suite.chainB.GetSimApp().BankKeeper.GetBalance(suite.chainB.GetContext(), suite.chainB.SenderAccount.GetAddress(), "acanto")
	fmt.Println(balanceBefore, balance, balanceCanto)

	coinSentFromAToB := types.GetTransferCoin(path.EndpointB.ChannelConfig.PortID, path.EndpointB.ChannelID, sdk.DefaultBondDenom, amount)
	suite.Require().True(coinSentFromAToB.Add(balanceBefore).Amount.GTE(balance.Amount))
	// check whether the canto is swapped and the amount is greater than the threshold
	suite.Require().True(balanceCanto.Amount.GTE(balanceCantoBefore.Amount))
	suite.Require().True(balanceCanto.Amount.GTE(middlewareParams.AutoSwapThreshold))

	// Send again from chainA to chainB
	coinToSendToB = suite.chainA.GetSimApp().BankKeeper.GetBalance(suite.chainA.GetContext(), suite.chainA.SenderAccount.GetAddress(), sdk.DefaultBondDenom)
	balanceBefore = suite.chainB.GetSimApp().BankKeeper.GetBalance(suite.chainB.GetContext(), suite.chainB.SenderAccount.GetAddress(), voucherDenomTrace.IBCDenom())

	msg = types.NewMsgTransfer(path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID, coinToSendToB, suite.chainA.SenderAccount.GetAddress().String(), suite.chainB.SenderAccount.GetAddress().String(), timeoutHeight, 0)

	res, err = suite.chainA.SendMsgs(msg)
	suite.Require().NoError(err) // message committed

	packet, err = ibctesting.ParsePacketFromEvents(res.GetEvents())
	suite.Require().NoError(err)

	// relay send
	err = path.RelayPacket(packet)
	suite.Require().NoError(err) // relay committed

	balance = suite.chainB.GetSimApp().BankKeeper.GetBalance(suite.chainB.GetContext(), suite.chainB.SenderAccount.GetAddress(), voucherDenomTrace.IBCDenom())
	coinSentFromAToB = types.GetTransferCoin(path.EndpointB.ChannelConfig.PortID, path.EndpointB.ChannelID, sdk.DefaultBondDenom, coinToSendToB.Amount)
	balancCantoAfter := suite.chainB.GetSimApp().BankKeeper.GetBalance(suite.chainB.GetContext(), suite.chainB.SenderAccount.GetAddress(), "acanto")
	suite.Require().Equal(balance, balanceBefore.Add(coinSentFromAToB))
	suite.Require().Equal(balancCantoAfter.Amount, balanceCanto.Amount)

	/*



		// setup between chainB to chainC
		// NOTE:
		// pathBtoC.EndpointA = endpoint on chainB
		// pathBtoC.EndpointB = endpoint on chainC
		pathBtoC := NewTransferPath(suite.chainB, suite.chainC)
		suite.coordinator.Setup(pathBtoC)

		// send from chainB to chainC
		msg = types.NewMsgTransfer(pathBtoC.EndpointA.ChannelConfig.PortID, pathBtoC.EndpointA.ChannelID, coinSentFromAToB, suite.chainB.SenderAccount.GetAddress().String(), suite.chainC.SenderAccount.GetAddress().String(), timeoutHeight, 0)
		res, err = suite.chainB.SendMsgs(msg)
		suite.Require().NoError(err) // message committed

		packet, err = ibctesting.ParsePacketFromEvents(res.GetEvents())
		suite.Require().NoError(err)

		err = pathBtoC.RelayPacket(packet)
		suite.Require().NoError(err) // relay committed

		// NOTE: fungible token is prefixed with the full trace in order to verify the packet commitment
		fullDenomPath := types.GetPrefixedDenom(pathBtoC.EndpointB.ChannelConfig.PortID, pathBtoC.EndpointB.ChannelID, voucherDenomTrace.GetFullDenomPath())

		coinSentFromBToC := sdk.NewCoin(types.ParseDenomTrace(fullDenomPath).IBCDenom(), amount)
		balance = suite.chainC.GetSimApp().BankKeeper.GetBalance(suite.chainC.GetContext(), suite.chainC.SenderAccount.GetAddress(), coinSentFromBToC.Denom)
		balanceCanto = suite.chainC.GetSimApp().BankKeeper.GetBalance(suite.chainC.GetContext(), suite.chainC.SenderAccount.GetAddress(), "acanto")
		fmt.Println(balanceCanto)

		// check that the balance is updated on chainC
		suite.Require().Equal(coinSentFromBToC, balance)

		// check that balance on chain B is empty
		balance = suite.chainB.GetSimApp().BankKeeper.GetBalance(suite.chainB.GetContext(), suite.chainB.SenderAccount.GetAddress(), coinSentFromBToC.Denom)
		suite.Require().Zero(balance.Amount.Int64())

		// send from chainC back to chainB
		msg = types.NewMsgTransfer(pathBtoC.EndpointB.ChannelConfig.PortID, pathBtoC.EndpointB.ChannelID, coinSentFromBToC, suite.chainC.SenderAccount.GetAddress().String(), suite.chainB.SenderAccount.GetAddress().String(), timeoutHeight, 0)
		res, err = suite.chainC.SendMsgs(msg)
		suite.Require().NoError(err) // message committed

		packet, err = ibctesting.ParsePacketFromEvents(res.GetEvents())
		suite.Require().NoError(err)

		err = pathBtoC.RelayPacket(packet)
		suite.Require().NoError(err) // relay committed

		balance = suite.chainB.GetSimApp().BankKeeper.GetBalance(suite.chainB.GetContext(), suite.chainB.SenderAccount.GetAddress(), coinSentFromAToB.Denom)

		// check that the balance on chainA returned back to the original state
		suite.Require().Equal(coinSentFromAToB, balance)

		// check that module account escrow address is empty
		escrowAddress := types.GetEscrowAddress(packet.GetDestPort(), packet.GetDestChannel())
		balance = suite.chainB.GetSimApp().BankKeeper.GetBalance(suite.chainB.GetContext(), escrowAddress, sdk.DefaultBondDenom)
		suite.Require().Equal(sdk.NewCoin(sdk.DefaultBondDenom, sdk.ZeroInt()), balance)

		// check that balance on chain B is empty
		balance = suite.chainC.GetSimApp().BankKeeper.GetBalance(suite.chainC.GetContext(), suite.chainC.SenderAccount.GetAddress(), voucherDenomTrace.IBCDenom())
		balanceCanto = suite.chainC.GetSimApp().BankKeeper.GetBalance(suite.chainC.GetContext(), suite.chainC.SenderAccount.GetAddress(), "acanto")
		fmt.Println(balanceCanto)

		suite.Require().Zero(balance.Amount.Int64())

	*/
}

func TestTransferTestSuite(t *testing.T) {
	suite.Run(t, new(TransferTestSuite))
}
