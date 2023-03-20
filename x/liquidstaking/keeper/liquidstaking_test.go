package keeper_test

import (
	"github.com/Canto-Network/Canto/v6/x/liquidstaking/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (suite *KeeperTestSuite) TestDoLiquidStake() {
	// TODO: Write fail test cases first
}

func (suite *KeeperTestSuite) TestDoLiquidStakeFailCases() {
	_ = suite.CreateValidators([]int64{10, 10, 10})
	// TODO: Create pairing insurances with created validators

	bondDenom := suite.app.StakingKeeper.BondDenom(suite.ctx)
	minimumRequirement := sdk.NewCoin(bondDenom, sdk.NewInt(types.ChunkSize))
	addrs := suite.AddTestAddrs(10, sdk.NewInt(minimumRequirement.Amount.Int64()))

	// CaseN: Not enough amount to liquid stake
	// acc1 tries to liquid stake 2 * ChunkSize tokens, but he has only ChunkSize tokens
	acc1 := addrs[0]
	_, _, err := suite.app.LiquidStakingKeeper.DoLiquidStake(
		suite.ctx,
		acc1,
		minimumRequirement.AddAmount(
			sdk.NewInt(types.ChunkSize),
		),
	)
	suite.ErrorIs(err, sdkerrors.ErrInsufficientFunds)
}
