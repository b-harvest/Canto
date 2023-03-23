package keeper_test

import (
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/Canto-Network/Canto/v6/app"
	"github.com/Canto-Network/Canto/v6/x/liquidstaking/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/evmos/ethermint/crypto/ethsecp256k1"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/tmhash"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmversion "github.com/tendermint/tendermint/proto/tendermint/version"
	"github.com/tendermint/tendermint/version"
)

type KeeperTestSuite struct {
	suite.Suite
	// use keeper for tests
	ctx         sdk.Context
	app         *app.Canto
	queryClient types.QueryClient
	consAddress sdk.ConsAddress
	address     common.Address
	validator   stakingtypes.Validator

	denom string
}

var s *KeeperTestSuite

func TestKeeperTestSuite(t *testing.T) {
	s = new(KeeperTestSuite)

	suite.Run(t, s)

	// Run Ginkgo integration tests
	RegisterFailHandler(Fail)
	RunSpecs(t, "Keeper Suite")
}

func (suite *KeeperTestSuite) SetupTest() {
	// instantiate app
	suite.app = app.Setup(false, nil)
	// initialize ctx for tests
	suite.SetupApp()
}

func (suite *KeeperTestSuite) SetupApp() {
	t := suite.T()

	priv, err := ethsecp256k1.GenerateKey()
	require.NoError(t, err)

	suite.address = common.BytesToAddress(priv.PubKey().Address().Bytes())
	suite.denom = "acanto"

	// consensus key
	privCons, err := ethsecp256k1.GenerateKey()
	require.NoError(t, err)
	suite.consAddress = sdk.ConsAddress(privCons.PubKey().Address())
	suite.ctx = suite.app.BaseApp.NewContext(false, tmproto.Header{
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

	queryHelper := baseapp.NewQueryServerTestHelper(suite.ctx, suite.app.InterfaceRegistry())
	types.RegisterQueryServer(queryHelper, suite.app.LiquidStakingKeeper)
	suite.queryClient = types.NewQueryClient(queryHelper)

	stakingParams := suite.app.StakingKeeper.GetParams(suite.ctx)
	stakingParams.BondDenom = suite.denom
	suite.app.StakingKeeper.SetParams(suite.ctx, stakingParams)

	// Set Validator
	valAddr := sdk.ValAddress(suite.address.Bytes())
	validator, err := stakingtypes.NewValidator(valAddr, privCons.PubKey(), stakingtypes.Description{})
	require.NoError(t, err)

	validator = stakingkeeper.TestingUpdateValidator(suite.app.StakingKeeper, suite.ctx, validator, true)
	suite.app.StakingKeeper.AfterValidatorCreated(suite.ctx, validator.GetOperator())
	err = suite.app.StakingKeeper.SetValidatorByConsAddr(suite.ctx, validator)
	require.NoError(t, err)

	validators := s.app.StakingKeeper.GetValidators(s.ctx, 1)
	suite.validator = validators[0]
}

// Commit commits and starts a new block with an updated context.
func (suite *KeeperTestSuite) Commit() {
	suite.CommitAfter(time.Second * 0)
}

// Commit commits a block at a given time.
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
}

func (suite *KeeperTestSuite) CreateValidators(powers []int64) (valAddrs []sdk.ValAddress) {
	notBondedPool := suite.app.StakingKeeper.GetNotBondedPool(suite.ctx)

	for _, power := range powers {
		priv := ed25519.GenPrivKey()
		valAddr := sdk.ValAddress(priv.PubKey().Address().Bytes())
		validator, err := stakingtypes.NewValidator(valAddr, priv.PubKey(), stakingtypes.Description{})
		suite.NoError(err)

		tokens := suite.app.StakingKeeper.TokensFromConsensusPower(suite.ctx, power)

		// Mint tokens for not bonded pool
		err = suite.app.BankKeeper.MintCoins(suite.ctx, types.ModuleName, sdk.NewCoins(sdk.NewCoin(suite.denom, tokens)))
		suite.NoError(err)
		err = suite.app.BankKeeper.SendCoinsFromModuleToModule(suite.ctx, types.ModuleName, notBondedPool.GetName(), sdk.NewCoins(sdk.NewCoin(suite.denom, tokens)))
		suite.NoError(err)

		validator, _ = validator.AddTokensFromDel(tokens)
		_, err = validator.SetInitialCommission(stakingtypes.NewCommission(sdk.NewDecWithPrec(10, 2), sdk.NewDecWithPrec(10, 2), sdk.NewDecWithPrec(10, 2)))
		if err != nil {
			return
		}
		validator = stakingkeeper.TestingUpdateValidator(suite.app.StakingKeeper, suite.ctx, validator, true)
		suite.app.StakingKeeper.AfterValidatorCreated(suite.ctx, validator.GetOperator())
		suite.app.StakingKeeper.SetValidator(suite.ctx, validator)
		suite.NoError(suite.app.StakingKeeper.SetValidatorByConsAddr(suite.ctx, validator))
		valAddrs = append(valAddrs, valAddr)
	}
	suite.app.EndBlocker(suite.ctx, abci.RequestEndBlock{})
	return
}

// Add test addresses with funds
func (suite *KeeperTestSuite) AddTestAddrs(accNum int, amount sdk.Int) ([]sdk.AccAddress, []sdk.Coin) {
	addrs := make([]sdk.AccAddress, 0, accNum)
	balances := make([]sdk.Coin, 0, accNum)
	for i := 0; i < accNum; i++ {
		addr := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
		addrs = append(addrs, addr)
		balances = append(balances, sdk.NewCoin(suite.denom, amount))

		// fund each account
		err := suite.app.BankKeeper.MintCoins(suite.ctx, types.ModuleName, sdk.NewCoins(sdk.NewCoin(suite.denom, amount)))
		suite.NoError(err)
		err = suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, types.ModuleName, addr, sdk.NewCoins(sdk.NewCoin(suite.denom, amount)))
		suite.NoError(err)
	}
	return addrs, balances
}

// TODO: impl
//func (suite *KeeperTestSuite) advanceHeight(height int) {
//	feeCollector := suite.app.AccountKeeper.GetModuleAccount(suite.ctx, authtypes.FeeCollectorName)
//	for i := 0; i < height; i++ {
//		s.ctx = s.ctx.WithBlockHeight(s.ctx.BlockHeight() + 1).WithBlockTime(s.ctx.BlockTime().Add(time.Second))
//	}
//}
