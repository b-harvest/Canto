package keeper_test

import (
	"math"
	"testing"
	"time"

	"github.com/Canto-Network/Canto-Testnet-v2/v1/app"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/tmhash"
	tmrand "github.com/tendermint/tendermint/libs/rand"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmversion "github.com/tendermint/tendermint/proto/tendermint/version"
	"github.com/tendermint/tendermint/version"

	evm "github.com/Canto-Network/ethermint-v2/x/evm/types"

	epochstypes "github.com/Canto-Network/Canto-Testnet-v2/v1/x/epochs/types"
	"github.com/Canto-Network/Canto-Testnet-v2/v1/x/liquidstaking/keeper"
	"github.com/Canto-Network/Canto-Testnet-v2/v1/x/liquidstaking/types"
)

type KeeperTestSuite struct {
	suite.Suite

	ctx            sdk.Context
	app            *app.Canto
	keeper         keeper.Keeper
	queryClientEvm evm.QueryClient
	queryClient    types.QueryClient
	consAddress    sdk.ConsAddress
}

func TestKeeperTestSuite(t *testing.T) {
	s := new(KeeperTestSuite)
	suite.Run(t, s)
}

func (suite *KeeperTestSuite) SetupTest() {
	suite.DoSetupTest(suite.T())
}

// Test helpers
func (suite *KeeperTestSuite) DoSetupTest(t require.TestingT) {
	checkTx := false

	// init app
	suite.app = app.Setup(checkTx, nil)

	// setup context
	suite.ctx = suite.app.BaseApp.NewContext(checkTx, tmproto.Header{
		Height:          1,
		ChainID:         "canto_9001-1",
		Time:            time.Now().UTC(),
		ProposerAddress: suite.consAddress.Bytes(),

		Version: tmversion.Consensus{
			Block: version.BlockProtocol,
		},
		LastBlockId: tmproto.BlockID{
			Hash: tmhash.Sum([]byte("block_id")),
			PartSetHeader: tmproto.PartSetHeader{
				Total: 11,
				Hash:  tmhash.Sum([]byte("partset_header")),
			},
		},
		AppHash:            tmhash.Sum([]byte("app")),
		DataHash:           tmhash.Sum([]byte("data")),
		EvidenceHash:       tmhash.Sum([]byte("evidence")),
		ValidatorsHash:     tmhash.Sum([]byte("validators")),
		NextValidatorsHash: tmhash.Sum([]byte("next_validators")),
		ConsensusHash:      tmhash.Sum([]byte("consensus")),
		LastResultsHash:    tmhash.Sum([]byte("last_result")),
	})
	suite.keeper = suite.app.LiquidStakingKeeper

	// setup query helpers
	queryHelper := baseapp.NewQueryServerTestHelper(suite.ctx, suite.app.InterfaceRegistry())
	types.RegisterQueryServer(queryHelper, keeper.NewQuerier(suite.app.LiquidStakingKeeper))
	suite.queryClient = types.NewQueryClient(queryHelper)

	// Set epoch start time and height for all epoch identifiers
	identifiers := []string{epochstypes.WeekEpochID, epochstypes.DayEpochID}
	for _, identifier := range identifiers {
		epoch, found := suite.app.EpochsKeeper.GetEpochInfo(suite.ctx, identifier)
		suite.Require().True(found)
		epoch.StartTime = suite.ctx.BlockTime()
		epoch.CurrentEpochStartHeight = suite.ctx.BlockHeight()
		suite.app.EpochsKeeper.SetEpochInfo(suite.ctx, epoch)
	}
}

func (suite *KeeperTestSuite) Commit() {
	suite.CommitAfter(time.Nanosecond)
}

func (suite *KeeperTestSuite) CommitAfter(t time.Duration) {
	_ = suite.app.Commit()
	header := suite.ctx.BlockHeader()
	header.Height += 1
	header.Time = header.Time.Add(t)
	suite.app.BeginBlock(abci.RequestBeginBlock{
		Header: header,
	})

	// update ctx
	suite.ctx = suite.app.BaseApp.NewContext(false, header)

	queryHelper := baseapp.NewQueryServerTestHelper(suite.ctx, suite.app.InterfaceRegistry())
	evm.RegisterQueryServer(queryHelper, suite.app.EvmKeeper)
	suite.queryClientEvm = evm.NewQueryClient(queryHelper)
}

func generateRandomId() uint64 {
	return uint64(tmrand.Int63n(10000))
}

func generateRandomAccount() sdk.AccAddress {
	pk := ed25519.GenPrivKey().PubKey()
	return sdk.AccAddress(pk.Address())
}

func generateRandomTokenAmount() sdk.Int {
	return sdk.NewInt(tmrand.Int63n(math.MaxInt64)).AddRaw(1)
}

func generateRandomValidatorAccount(suite *KeeperTestSuite) (ret sdk.ValAddress) {
	pk := ed25519.GenPrivKey().PubKey()
	ret = sdk.ValAddress(pk.Address())
	val, err := stakingtypes.NewValidator(ret, pk, stakingtypes.Description{})
	rate := sdk.NewDec(tmrand.Int63n(100)).QuoInt64(int64(100))
	val.SetInitialCommission(stakingtypes.NewCommission(rate, sdk.OneDec(), sdk.OneDec()))
	suite.Require().NoError(err)

	suite.app.StakingKeeper.SetValidator(suite.ctx, val)
	err = suite.app.StakingKeeper.SetValidatorByConsAddr(suite.ctx, val)
	suite.Require().NoError(err)

	suite.app.StakingKeeper.AfterValidatorCreated(suite.ctx, val.GetOperator())

	return val.GetOperator()
}

func depositCoinsIntoModule(suite *KeeperTestSuite, coin sdk.Coin) {
	err := suite.app.BankKeeper.MintCoins(suite.ctx, types.ModuleName, sdk.NewCoins(coin))
	suite.Require().NoError(err)
}

func generateRandomAliveChunk(suite *KeeperTestSuite, id types.AliveChunkId) types.AliveChunk {
	return types.NewAliveChunk(
		id,
		generateRandomChunkBondRequest(generateRandomId()),
		generateRandomInsuranceBid(suite, generateRandomId()),
	)
}

func generateRandomInsuranceBid(suite *KeeperTestSuite, id types.InsuranceBidId) types.InsuranceBid {
	return types.InsuranceBid{
		Id:                       id,
		ValidatorAddress:         generateRandomValidatorAccount(suite).String(),
		InsuranceProviderAddress: generateRandomAccount().String(),
		InsuranceAmount:          generateRandomTokenAmount(),
		InsuranceFeeRate:         sdk.NewDec(tmrand.Int63n(math.MaxInt64)),
	}
}

func generateRandomUnbondingChunk(suite *KeeperTestSuite, id types.UnbondingChunkId) types.UnbondingChunk {
	return types.UnbondingChunk{
		Id:                       id,
		ValidatorAddress:         generateRandomValidatorAccount(suite).String(),
		InsuranceProviderAddress: generateRandomAccount().String(),
		TokenAmount:              generateRandomTokenAmount(),
		InsuranceAmount:          generateRandomTokenAmount(),
	}
}

func generateRandomChunkBondRequest(id types.ChunkBondRequestId) types.ChunkBondRequest {
	return types.ChunkBondRequest{
		Id:          id,
		Address:     generateRandomAccount().String(),
		TokenAmount: generateRandomTokenAmount(),
	}
}

func generateRandomChunkBondRequestWithTokenAmount(id types.ChunkBondRequestId, tokenAmount sdk.Int) types.ChunkBondRequest {
	return types.ChunkBondRequest{
		Id:          id,
		Address:     generateRandomAccount().String(),
		TokenAmount: tokenAmount,
	}
}

func generateRandomChunkUnbondRequest(id types.ChunkUnbondRequestId) types.ChunkUnbondRequest {
	return types.ChunkUnbondRequest{
		Id:             id,
		Address:        generateRandomAccount().String(),
		NumChunkUnbond: uint64(tmrand.Int63n(math.MaxInt64)),
	}
}

func generateRandomInsuranceUnbondRequest(id types.AliveChunkId) types.InsuranceUnbondRequest {
	return types.InsuranceUnbondRequest{
		AliveChunkId: id,
	}
}

func accountWithCoins(addr sdk.AccAddress, bankKeeper bankkeeper.Keeper, ctx sdk.Context, coins sdk.Coins) sdk.AccAddress {
	// TODO: check this logic
	err := bankKeeper.MintCoins(ctx, types.ModuleName, coins)
	if err != nil {
		panic(err)
	}
	err = bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, addr, coins)
	if err != nil {
		panic(err)
	}
	return addr
}
