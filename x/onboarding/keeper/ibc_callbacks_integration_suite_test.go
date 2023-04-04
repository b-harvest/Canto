package keeper_test

import (
	"bytes"
	"fmt"
	"strconv"
	"testing"
	"time"

	erc20types "github.com/Canto-Network/Canto/v6/x/erc20/types"
	inflationtypes "github.com/Canto-Network/Canto/v6/x/inflation/types"
	coinswaptypes "github.com/b-harvest/coinswap/modules/coinswap/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/stretchr/testify/suite"

	sdk "github.com/cosmos/cosmos-sdk/types"
	transfertypes "github.com/cosmos/ibc-go/v3/modules/apps/transfer/types"
	clienttypes "github.com/cosmos/ibc-go/v3/modules/core/02-client/types"
	channeltypes "github.com/cosmos/ibc-go/v3/modules/core/04-channel/types"
	ibcgotesting "github.com/cosmos/ibc-go/v3/testing"

	ibctesting "github.com/Canto-Network/Canto/v6/ibc/testing"

	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"

	"github.com/Canto-Network/Canto/v6/app"
	"github.com/Canto-Network/Canto/v6/x/onboarding/types"
)

type IBCTestingSuite struct {
	suite.Suite
	coordinator *ibcgotesting.Coordinator

	// testing chains used for convenience and readability
	cantoChain      *ibcgotesting.TestChain
	IBCGravityChain *ibcgotesting.TestChain
	IBCCosmosChain  *ibcgotesting.TestChain

	pathGravitycanto  *ibcgotesting.Path
	pathCosmoscanto   *ibcgotesting.Path
	pathGravityCosmos *ibcgotesting.Path
}

var s *IBCTestingSuite

func TestIBCTestingSuite(t *testing.T) {
	s = new(IBCTestingSuite)
	suite.Run(t, s)

	// Run Ginkgo integration tests
	RegisterFailHandler(Fail)
	RunSpecs(t, "Keeper Suite")
}

func (suite *IBCTestingSuite) SetupTest() {
	// initializes 3 test chains
	suite.coordinator = ibctesting.NewCoordinator(suite.T(), 1, 2)
	suite.cantoChain = suite.coordinator.GetChain(ibcgotesting.GetChainID(1))
	suite.IBCGravityChain = suite.coordinator.GetChain(ibcgotesting.GetChainID(2))
	suite.IBCCosmosChain = suite.coordinator.GetChain(ibcgotesting.GetChainID(3))
	suite.coordinator.CommitNBlocks(suite.cantoChain, 2)
	suite.coordinator.CommitNBlocks(suite.IBCGravityChain, 2)
	suite.coordinator.CommitNBlocks(suite.IBCCosmosChain, 2)

	// Mint coins on the gravity side which we'll use to unlock our acanto
	coinUsdc := sdk.NewCoin("uUSDC", sdk.NewIntWithDecimal(10000, 6))
	coinUsdt := sdk.NewCoin("uUSDT", sdk.NewIntWithDecimal(10000, 6))
	coins := sdk.NewCoins(coinUsdc, coinUsdt)
	err := suite.IBCGravityChain.GetSimApp().BankKeeper.MintCoins(suite.IBCGravityChain.GetContext(), minttypes.ModuleName, coins)
	suite.Require().NoError(err)
	err = suite.IBCGravityChain.GetSimApp().BankKeeper.SendCoinsFromModuleToAccount(suite.IBCGravityChain.GetContext(), minttypes.ModuleName, suite.IBCGravityChain.SenderAccount.GetAddress(), coins)
	suite.Require().NoError(err)

	// Mint coins on the cosmos side which we'll use to unlock our acanto
	coinAtom := sdk.NewCoin("uatom", sdk.NewIntWithDecimal(10000, 6))
	coins = sdk.NewCoins(coinAtom)
	err = suite.IBCCosmosChain.GetSimApp().BankKeeper.MintCoins(suite.IBCCosmosChain.GetContext(), minttypes.ModuleName, coins)
	suite.Require().NoError(err)
	err = suite.IBCCosmosChain.GetSimApp().BankKeeper.SendCoinsFromModuleToAccount(suite.IBCCosmosChain.GetContext(), minttypes.ModuleName, suite.IBCCosmosChain.SenderAccount.GetAddress(), coins)
	suite.Require().NoError(err)

	params := types.DefaultParams()
	params.EnableOnboarding = true
	suite.cantoChain.App.(*app.Canto).OnboardingKeeper.SetParams(suite.cantoChain.GetContext(), params)

	suite.pathGravitycanto = ibctesting.NewTransferPath(suite.IBCGravityChain, suite.cantoChain) // clientID, connectionID, channelID empty
	suite.pathCosmoscanto = ibctesting.NewTransferPath(suite.IBCCosmosChain, suite.cantoChain)
	suite.pathGravityCosmos = ibctesting.NewTransferPath(suite.IBCCosmosChain, suite.IBCGravityChain)
	suite.coordinator.Setup(suite.pathGravitycanto) // clientID, connectionID, channelID filled
	suite.coordinator.Setup(suite.pathCosmoscanto)
	suite.coordinator.Setup(suite.pathGravityCosmos)
	suite.Require().Equal("07-tendermint-0", suite.pathGravitycanto.EndpointA.ClientID)
	suite.Require().Equal("connection-0", suite.pathGravitycanto.EndpointA.ConnectionID)
	suite.Require().Equal("channel-0", suite.pathGravitycanto.EndpointA.ChannelID)

	// Set the proposer address for the current header
	// It because EVMKeeper.GetCoinbaseAddress requires ProposerAddress in block header
	suite.cantoChain.CurrentHeader.ProposerAddress = suite.cantoChain.LastHeader.ValidatorSet.Proposer.Address
	suite.IBCGravityChain.CurrentHeader.ProposerAddress = suite.IBCGravityChain.LastHeader.ValidatorSet.Proposer.Address
	suite.IBCCosmosChain.CurrentHeader.ProposerAddress = suite.IBCCosmosChain.LastHeader.ValidatorSet.Proposer.Address
}

func (suite *IBCTestingSuite) FundCantoChain(coins sdk.Coins) {
	err := suite.cantoChain.App.(*app.Canto).BankKeeper.MintCoins(suite.cantoChain.GetContext(), inflationtypes.ModuleName, coins)
	suite.Require().NoError(err)
	err = suite.cantoChain.App.(*app.Canto).BankKeeper.SendCoinsFromModuleToAccount(suite.cantoChain.GetContext(), inflationtypes.ModuleName, suite.cantoChain.SenderAccount.GetAddress(), coins)
	suite.Require().NoError(err)
}

func (suite *IBCTestingSuite) setupRegisterCoin(metadata banktypes.Metadata) *erc20types.TokenPair {
	err := suite.cantoChain.App.(*app.Canto).BankKeeper.MintCoins(suite.cantoChain.GetContext(), inflationtypes.ModuleName, sdk.Coins{sdk.NewInt64Coin(metadata.Base, 1)})
	suite.Require().NoError(err)

	pair, err := suite.cantoChain.App.(*app.Canto).Erc20Keeper.RegisterCoin(suite.cantoChain.GetContext(), metadata)
	suite.Require().NoError(err)
	return pair
}

func (suite *IBCTestingSuite) CreatePool(denom string) {

	coincanto := sdk.NewCoin("acanto", sdk.NewIntWithDecimal(10000, 18))
	coinIBC := sdk.NewCoin(denom, sdk.NewIntWithDecimal(10000, 6))
	coins := sdk.NewCoins(coincanto, coinIBC)
	suite.FundCantoChain(coins)

	// create ibc/uUSDC, acanto pool
	coinswapParams := suite.cantoChain.App.(*app.Canto).CoinswapKeeper.GetParams(suite.cantoChain.GetContext())
	coinswapParams.MaxSwapAmount = sdk.NewCoins(sdk.NewCoin(denom, sdk.NewIntWithDecimal(10, 6)))
	suite.cantoChain.App.(*app.Canto).CoinswapKeeper.SetParams(suite.cantoChain.GetContext(), coinswapParams)
	msgAddLiquidity := coinswaptypes.MsgAddLiquidity{
		MaxToken:         sdk.NewCoin(denom, sdk.NewIntWithDecimal(10000, 6)),
		ExactStandardAmt: sdk.NewIntWithDecimal(10000, 18),
		MinLiquidity:     sdk.NewInt(1),
		Deadline:         time.Now().Add(time.Minute * 10).Unix(),
		Sender:           suite.cantoChain.SenderAccount.GetAddress().String(),
	}
	suite.cantoChain.App.(*app.Canto).CoinswapKeeper.AddLiquidity(suite.cantoChain.GetContext(), &msgAddLiquidity)
}

var (
	timeoutHeight   = clienttypes.NewHeight(1000, 1000)
	uusdcDenomtrace = transfertypes.DenomTrace{
		Path:      "transfer/channel-0",
		BaseDenom: "uUSDC",
	}
	uusdcIbcdenom = uusdcDenomtrace.IBCDenom()

	uusdtDenomtrace = transfertypes.DenomTrace{
		Path:      "transfer/channel-0",
		BaseDenom: "uUSDT",
	}
	uusdtIbcdenom = uusdtDenomtrace.IBCDenom()

	uosmoDenomtrace = transfertypes.DenomTrace{
		Path:      "transfer/channel-0",
		BaseDenom: "uosmo",
	}

	uosmoIbcdenom = uosmoDenomtrace.IBCDenom()

	uatomDenomtrace = transfertypes.DenomTrace{
		Path:      "transfer/channel-1",
		BaseDenom: "uatom",
	}
	uatomIbcdenom = uatomDenomtrace.IBCDenom()

	acantoDenomtrace = transfertypes.DenomTrace{
		Path:      "transfer/channel-0",
		BaseDenom: "acanto",
	}
	acantoIbcdenom = acantoDenomtrace.IBCDenom()

	uatomOsmoDenomtrace = transfertypes.DenomTrace{
		Path:      "transfer/channel-0/transfer/channel-1",
		BaseDenom: "uatom",
	}
	uatomOsmoIbcdenom = uatomOsmoDenomtrace.IBCDenom()
)

func (suite *IBCTestingSuite) SendAndReceiveMessage(path *ibcgotesting.Path, origin *ibcgotesting.TestChain, coin string, amount int64, sender string, receiver string, seq uint64) *sdk.Result {
	// Send coin from A to B
	transferMsg := transfertypes.NewMsgTransfer(path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID, sdk.NewCoin(coin, sdk.NewInt(amount)), sender, receiver, timeoutHeight, 0)
	_, err := origin.SendMsgs(transferMsg)
	suite.Require().NoError(err) // message committed
	// Recreate the packet that was sent
	transfer := transfertypes.NewFungibleTokenPacketData(coin, strconv.Itoa(int(amount)), sender, receiver)
	packet := channeltypes.NewPacket(transfer.GetBytes(), seq, path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID, path.EndpointB.ChannelConfig.PortID, path.EndpointB.ChannelID, timeoutHeight, 0)
	// Receive message on the counterparty side, and send ack

	// original call
	// err = path.RelayPacket(packet)

	// patched RelayPacket call to get res
	res, err := RelayPacket(path, packet)

	// ---------- Temporary Print for Debugging
	//for _, ev := range res.GetEvents() {
	//	fmt.Println(string(ev.Type))
	//	for _, e := range ev.Attributes {
	//		fmt.Println("\t", string(e.Key), string(e.Value))
	//	}
	//}
	// ---------- Temporary Print for Debugging

	suite.Require().NoError(err)
	return res
}

// RelayPacket attempts to relay the packet first on EndpointA and then on EndpointB
// if EndpointA does not contain a packet commitment for that packet. An error is returned
// if a relay step fails or the packet commitment does not exist on either endpoint.
func RelayPacket(path *ibcgotesting.Path, packet channeltypes.Packet) (*sdk.Result, error) {
	pc := path.EndpointA.Chain.App.GetIBCKeeper().ChannelKeeper.GetPacketCommitment(path.EndpointA.Chain.GetContext(), packet.GetSourcePort(), packet.GetSourceChannel(), packet.GetSequence())
	if bytes.Equal(pc, channeltypes.CommitPacket(path.EndpointA.Chain.App.AppCodec(), packet)) {

		// packet found, relay from A to B
		if err := path.EndpointB.UpdateClient(); err != nil {
			return nil, err
		}

		res, err := path.EndpointB.RecvPacketWithResult(packet)
		if err != nil {
			return nil, err
		}

		ack, err := ibcgotesting.ParseAckFromEvents(res.GetEvents())
		if err != nil {
			return nil, err
		}

		if err := path.EndpointA.AcknowledgePacket(packet, ack); err != nil {
			return nil, err
		}

		return res, nil
	}

	pc = path.EndpointB.Chain.App.GetIBCKeeper().ChannelKeeper.GetPacketCommitment(path.EndpointB.Chain.GetContext(), packet.GetSourcePort(), packet.GetSourceChannel(), packet.GetSequence())
	if bytes.Equal(pc, channeltypes.CommitPacket(path.EndpointB.Chain.App.AppCodec(), packet)) {

		// packet found, relay B to A
		if err := path.EndpointA.UpdateClient(); err != nil {
			return nil, err
		}

		res, err := path.EndpointA.RecvPacketWithResult(packet)
		ack, err := ibcgotesting.ParseAckFromEvents(res.GetEvents())
		if err != nil {
			return nil, err
		}

		if err := path.EndpointB.AcknowledgePacket(packet, ack); err != nil {
			return nil, err
		}
		return res, nil
	}

	return nil, fmt.Errorf("packet commitment does not exist on either endpoint for provided packet")
}

func CreatePacket(amount, denom, sender, receiver, srcPort, srcChannel, dstPort, dstChannel string, seq, timeout uint64) channeltypes.Packet {
	transfer := transfertypes.FungibleTokenPacketData{
		Amount:   amount,
		Denom:    denom,
		Receiver: sender,
		Sender:   receiver,
	}
	return channeltypes.NewPacket(
		transfer.GetBytes(),
		seq,
		srcPort,
		srcChannel,
		dstPort,
		dstChannel,
		clienttypes.ZeroHeight(), // timeout height disabled
		timeout,
	)
}
