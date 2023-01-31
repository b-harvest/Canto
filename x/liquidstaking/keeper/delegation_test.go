package keeper_test

import (
	"strconv"

	"github.com/Canto-Network/Canto-Testnet-v2/v1/x/liquidstaking/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
	tmrand "github.com/tendermint/tendermint/libs/rand"
)

func (suite *KeeperTestSuite) TestKeeperDelegateRedelegateUndelegate() {
	seed := tmrand.Int63n(10000)
	tmrand.Seed(seed)

	bondDenom := suite.app.StakingKeeper.BondDenom(suite.ctx)
	tokenAmount := generateRandomTokenAmount()
	coin := sdk.NewCoin(bondDenom, tokenAmount)
	depositCoinsIntoModule(suite, coin)

	// prepare delegator
	del := generateRandomAccount()

	// prepare validator1
	valOperator := generateRandomValidatorAccount(suite)
	suite.Require().NotNil(valOperator)
	suite.app.EndBlocker(suite.ctx, abci.RequestEndBlock{})
	val := suite.app.StakingKeeper.Validator(suite.ctx, valOperator)
	suite.Require().NotNil(val)
	suite.Require().Equal(val.GetOperator(), valOperator)

	// prepare validator2
	newValOperator := generateRandomValidatorAccount(suite)
	suite.Require().NotNil(newValOperator)
	suite.app.EndBlocker(suite.ctx, abci.RequestEndBlock{})
	newVal := suite.app.StakingKeeper.Validator(suite.ctx, newValOperator)
	suite.Require().NotNil(newVal)
	suite.Require().Equal(newVal.GetOperator(), newValOperator)

	suite.Run("delegate random amount, seed: "+strconv.FormatInt(seed, 10), func() {
		err := suite.keeper.DelegateTokenAmount(suite.ctx, valOperator.String(), tokenAmount)
		suite.Require().NoError(err)

		coin := suite.app.BankKeeper.GetBalance(suite.ctx, types.LiquidStakingModuleAccount, bondDenom)
		suite.True(sdk.ZeroInt().Equal(coin.Amount))

		delegation, found := suite.app.StakingKeeper.GetDelegation(suite.ctx, types.LiquidStakingModuleAccount, valOperator)
		suite.Require().True(found)
		suite.Equal(types.LiquidStakingModuleAccount.String(), delegation.DelegatorAddress)
		suite.Equal(valOperator.String(), delegation.ValidatorAddress)
		suite.True(delegation.Shares.GT(sdk.ZeroDec()))
	})
	suite.Run("redelegate", func() {
		suite.keeper.RedelegateTokenAmount(suite.ctx, valOperator.String(), newValOperator.String(), tokenAmount)

		_, found := suite.app.StakingKeeper.GetDelegation(suite.ctx, types.LiquidStakingModuleAccount, valOperator)
		suite.False(found)

		delegation, found := suite.app.StakingKeeper.GetDelegation(suite.ctx, types.LiquidStakingModuleAccount, newValOperator)
		suite.Require().True(found)
		suite.Equal(types.LiquidStakingModuleAccount.String(), delegation.DelegatorAddress)
		suite.Equal(newValOperator.String(), delegation.ValidatorAddress)
		suite.True(delegation.Shares.GT(sdk.ZeroDec()))
	})
	suite.Run("undelegate", func() {
		suite.keeper.UndelegateTokenAmount(suite.ctx, del.String(), newValOperator.String(), tokenAmount)
		_, found := suite.app.StakingKeeper.GetDelegation(suite.ctx, types.LiquidStakingModuleAccount, newValOperator)
		suite.False(found)
	})
}
