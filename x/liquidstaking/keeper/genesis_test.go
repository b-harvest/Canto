package keeper_test

import (
	"github.com/Canto-Network/Canto/v6/x/liquidstaking/keeper"
	"github.com/Canto-Network/Canto/v6/x/liquidstaking/types"
)

func (suite *KeeperTestSuite) TestDefaultGenesis() {
	genState := types.DefaultGenesisState()

	keeper.InitGenesis(suite.ctx, suite.app.LiquidStakingKeeper, *genState)
	got := keeper.ExportGenesis(suite.ctx, suite.app.LiquidStakingKeeper)
	suite.Require().Equal(genState, got)
}
