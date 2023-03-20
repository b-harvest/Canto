package keeper_test

import (
	"github.com/Canto-Network/Canto/v6/x/liquidstaking/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/simulation"
	"math/rand"
)

func (suite *KeeperTestSuite) TestDoLiquidStake() {
	// TODO: Write fail test cases first
}

func (suite *KeeperTestSuite) TestDoLiquidStakeFailCases() {
	// init rand with seed
	s := rand.NewSource(0)
	r := rand.New(s)

	valAddrs := suite.CreateValidators([]int64{10, 10, 10})
	valNum := len(valAddrs)

	bondDenom := suite.app.StakingKeeper.BondDenom(suite.ctx)
	minimumRequirement := sdk.NewCoin(bondDenom, sdk.NewInt(types.ChunkSize))

	addrs := suite.AddTestAddrs(10, sdk.NewInt(minimumRequirement.Amount.Int64()))

	// TC: There are no pairing insurances yet. Insurances must be proivded to liquid stake
	acc1 := addrs[0]
	_, _, err := suite.app.LiquidStakingKeeper.DoLiquidStake(
		suite.ctx,
		acc1,
		minimumRequirement.AddAmount(
			sdk.NewInt(types.ChunkSize),
		),
	)
	suite.ErrorIs(err, types.ErrNoPairingInsurance)

	// TODO: Add tc for max paired chunk size exceeded

	fraction := sdk.NewDecWithPrec(types.SlashFractionInt, types.SlashFractionPrec)
	minimumCoverage := sdk.NewCoin(bondDenom, sdk.NewInt(minimumRequirement.Amount.ToDec().Mul(fraction).TruncateInt().Int64()))
	providers := suite.AddTestAddrs(10, minimumCoverage.Amount)

	// Provide insurances
	for i, provider := range providers {
		_, err := suite.app.LiquidStakingKeeper.DoInsuranceProvide(
			suite.ctx,
			provider,
			valAddrs[i%(valNum)], // can point same validators
			sdk.NewDecWithPrec(
				int64(simulation.RandIntBetween(r, 1, 10)),
				2,
			), // 1 ~ 10% insurance fee
			minimumCoverage,
		)
		suite.NoError(err)
	}

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
