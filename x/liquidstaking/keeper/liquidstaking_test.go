package keeper_test

import (
	"github.com/Canto-Network/Canto/v6/x/liquidstaking/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/simulation"
	"math/rand"
)

// Provide insurance with random fee (1 ~ 10%)
func (suite *KeeperTestSuite) provideInsurances(providers []sdk.AccAddress, valAddrs []sdk.ValAddress, amounts []sdk.Coin) []types.Insurance {
	s := rand.NewSource(0)
	r := rand.New(s)

	valNum := len(valAddrs)
	var providedInsurances []types.Insurance
	for i, provider := range providers {
		insurance, err := suite.app.LiquidStakingKeeper.DoInsuranceProvide(
			suite.ctx,
			provider,
			valAddrs[i%valNum], // can point same validators
			sdk.NewDecWithPrec(
				int64(simulation.RandIntBetween(r, 1, 10)),
				2,
			), // 1 ~ 10% insurance fee
			amounts[i],
		)
		suite.NoError(err)
		providedInsurances = append(providedInsurances, insurance)
	}
	return providedInsurances
}

func (suite *KeeperTestSuite) getMinimumRequirements() (oneChunkAmount, slashingCoverage sdk.Coin) {
	bondDenom := suite.app.StakingKeeper.BondDenom(suite.ctx)
	oneChunkAmount = sdk.NewCoin(bondDenom, sdk.NewInt(types.ChunkSize))
	fraction := sdk.MustNewDecFromStr(types.SlashFraction)
	slashingCoverage = sdk.NewCoin(bondDenom, sdk.NewInt(oneChunkAmount.Amount.ToDec().Mul(fraction).TruncateInt().Int64()))
	return
}

func (suite *KeeperTestSuite) TestDoInsuranceProvide() {
	valAddrs := suite.CreateValidators([]int64{10, 10, 10})
	_, minimumCoverage := suite.getMinimumRequirements()
	providers, amounts := suite.AddTestAddrs(10, minimumCoverage.Amount)

	providedInsurances := suite.provideInsurances(providers, valAddrs, amounts)
	var storedInsurances []types.Insurance
	suite.NoError(suite.app.LiquidStakingKeeper.IterateAllInsurances(suite.ctx, func(insurance types.Insurance) (bool, error) {
		storedInsurances = append(storedInsurances, insurance)
		return false, nil
	}))

	suite.Equal(len(providedInsurances), len(storedInsurances))
	for i, insurance := range providedInsurances {
		suite.True(insurance.Equal(storedInsurances[i]))
	}
}

func (suite *KeeperTestSuite) TestDoLiquidStake() {
	valAddrs := suite.CreateValidators([]int64{10, 10, 10})
	minimumRequirement, minimumCoverage := suite.getMinimumRequirements()
	providers, balances := suite.AddTestAddrs(10, minimumCoverage.Amount)
	suite.provideInsurances(providers, valAddrs, balances)

	delegators, balances := suite.AddTestAddrs(10, minimumRequirement.Amount)
	nas := suite.app.LiquidStakingKeeper.GetNetAmountState(suite.ctx)

	// First try
	del1 := delegators[0]
	amt1 := balances[0]
	newShares, lsTokenMintAmount, err := suite.app.LiquidStakingKeeper.DoLiquidStake(suite.ctx, del1, amt1)
	suite.NoError(err)
	suite.True(amt1.Amount.Equal(newShares.TruncateInt()))
	suite.True(amt1.Amount.Equal(lsTokenMintAmount))

	afterNas := suite.app.LiquidStakingKeeper.GetNetAmountState(suite.ctx)
	suite.True(nas.TotalLiquidTokens.Add(amt1.Amount).Equal(afterNas.TotalLiquidTokens))
	suite.True(nas.NetAmount.Add(amt1.Amount.ToDec()).Equal(afterNas.NetAmount))
}

func (suite *KeeperTestSuite) TestDoLiquidStakeFailCases() {
	valAddrs := suite.CreateValidators([]int64{10, 10, 10})
	minimumRequirement, minimumCoverage := suite.getMinimumRequirements()

	addrs, _ := suite.AddTestAddrs(10, sdk.NewInt(minimumRequirement.Amount.Int64()))

	// TC: There are no pairing insurances yet. Insurances must be provided to liquid stake
	acc1 := addrs[0]
	_, _, err := suite.app.LiquidStakingKeeper.DoLiquidStake(
		suite.ctx,
		acc1,
		minimumRequirement,
	)
	suite.ErrorIs(err, types.ErrNoPairingInsurance)

	// TODO: Add tc for max paired chunk size exceeded

	providers, balances := suite.AddTestAddrs(10, minimumCoverage.Amount)
	suite.provideInsurances(providers, valAddrs, balances)

	// TC: Not enough amount to liquid stake
	// acc1 tries to liquid stake 2 * ChunkSize tokens, but he has only ChunkSize tokens
	_, _, err = suite.app.LiquidStakingKeeper.DoLiquidStake(
		suite.ctx,
		acc1,
		minimumRequirement.AddAmount(
			sdk.NewInt(types.ChunkSize),
		),
	)
	suite.ErrorIs(err, sdkerrors.ErrInsufficientFunds)
}
