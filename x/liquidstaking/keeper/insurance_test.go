package keeper_test

import (
	"github.com/Canto-Network/Canto/v6/x/liquidstaking/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/evmos/ethermint/tests"
)

// Sets a bunch of insurances in the store and then get and ensure that each of them
// match up with what is stored on stack vs keeper
func (suite *KeeperTestSuite) TestInsuranceSetGet() {
	numberInsurances := 10
	insurances := GenerateInsurances(numberInsurances, false)
	for _, insurance := range insurances {
		suite.app.LiquidStakingKeeper.SetInsurance(suite.ctx, insurance)
	}

	for _, insurance := range insurances {
		id := insurance.Id
		status := insurance.Status
		chunkId := insurance.ChunkId
		// Get insurance from the store
		result, found := suite.app.LiquidStakingKeeper.GetInsurance(suite.ctx, id)

		// Validation
		suite.Require().True(found)
		suite.Require().Equal(result.Id, id)
		suite.Require().Equal(result.Status, status)
		suite.Require().Equal(result.ChunkId, chunkId)
	}
}

func (suite *KeeperTestSuite) TestDeleteInsurance() {
	numberInsurances := 10
	insurances := GenerateInsurances(numberInsurances, false)
	for _, insurance := range insurances {
		suite.app.LiquidStakingKeeper.SetInsurance(suite.ctx, insurance)
	}

	for _, insurance := range insurances {
		id := insurance.Id
		// Get insurance from the store
		result, found := suite.app.LiquidStakingKeeper.GetInsurance(suite.ctx, id)

		// Validation
		suite.Require().True(found)
		suite.Require().Equal(result.Id, id)

		// Delete insurance from the store
		suite.app.LiquidStakingKeeper.DeleteInsurance(suite.ctx, id)

		// Get insurance from the store
		result, found = suite.app.LiquidStakingKeeper.GetInsurance(suite.ctx, id)

		// Validation
		suite.Require().False(found)
		suite.Require().Equal(result.Id, uint64(0))
	}
}

func (suite *KeeperTestSuite) TestLastInsuranceIdSetGet() {
	// Set LastInsuranceId and retrieve it
	id := uint64(10)
	suite.app.LiquidStakingKeeper.SetLastInsuranceId(suite.ctx, id)

	result := suite.app.LiquidStakingKeeper.GetLastInsuranceId(suite.ctx)
	suite.Require().Equal(result, id)
}

// Creates a bunch of insurances
func GenerateInsurances(number int, sameAddress bool) []types.Insurance {
	insurances := make([]types.Insurance, number)
	for i := 0; i < number; i++ {
		var addr string
		if sameAddress {
			addr = authtypes.NewModuleAddress("test").String()
		} else {
			addr = sdk.AccAddress(tests.GenerateAddress().Bytes()).String()
		}

		insurances[i] = types.NewInsurance(uint64(i), addr, "", sdk.NewDec(0))
		insurances[i].ProviderAddress = addr
	}
	return insurances
}
