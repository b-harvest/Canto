package keeper_test

import (
	"github.com/Canto-Network/Canto/v6/x/liquidstaking/keeper"
	"github.com/Canto-Network/Canto/v6/x/liquidstaking/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (suite *KeeperTestSuite) TestChunksInvariant() {
	env := suite.setupLiquidStakeTestingEnv(
		testingEnvOptions{
			"TestDoWithdrawInsuranceCommission",
			3,
			TenPercentFeeRate,
			nil,
			onePower,
			nil,
			3,
			TenPercentFeeRate,
			nil,
			3,
			types.ChunkSize.MulRaw(500),
		},
	)
	_, broken := keeper.ChunksInvariant(suite.app.LiquidStakingKeeper)(suite.ctx)
	suite.False(broken, "completely normal")

	// 1: PAIRED CHUNK
	var origin, mutated types.Chunk = env.pairedChunks[0], env.pairedChunks[0]
	// forcefully change status of chunk as invalid
	{
		mutated.PairedInsuranceId = types.Empty
		suite.app.LiquidStakingKeeper.SetChunk(suite.ctx, mutated)
		_, broken = keeper.ChunksInvariant(suite.app.LiquidStakingKeeper)(suite.ctx)
		suite.True(broken, "paired chunk must have valid paired insurance id")
		// recover
		suite.app.LiquidStakingKeeper.SetChunk(suite.ctx, origin)
	}

	originIns, _ := suite.app.LiquidStakingKeeper.GetInsurance(suite.ctx, origin.PairedInsuranceId)
	// delete paired insurance
	{
		suite.app.LiquidStakingKeeper.DeleteInsurance(suite.ctx, originIns.Id)
		_, broken = keeper.ChunksInvariant(suite.app.LiquidStakingKeeper)(suite.ctx)
		suite.True(broken, "paired insurance must exist in store")
		// recover
		suite.app.LiquidStakingKeeper.SetInsurance(suite.ctx, originIns)
		suite.mustPassInvariants()
	}

	// forcefully change status of insurance as invalid
	{
		mutatedIns := originIns
		mutatedIns.Status = types.INSURANCE_STATUS_UNSPECIFIED
		suite.app.LiquidStakingKeeper.SetInsurance(suite.ctx, mutatedIns)
		_, broken = keeper.ChunksInvariant(suite.app.LiquidStakingKeeper)(suite.ctx)
		suite.True(broken, "insurance must have valid status")
		// recover
		suite.app.LiquidStakingKeeper.SetInsurance(suite.ctx, originIns)
		suite.mustPassInvariants()
	}

	originDel, _ := suite.app.StakingKeeper.GetDelegation(suite.ctx, origin.DerivedAddress(), originIns.GetValidator())
	// forcefully delete delegation obj of paired chunk
	{
		suite.app.StakingKeeper.RemoveDelegation(suite.ctx, originDel)
		_, broken = keeper.ChunksInvariant(suite.app.LiquidStakingKeeper)(suite.ctx)
		suite.True(broken, "delegation must exist in store")
		// recover
		suite.app.StakingKeeper.SetDelegation(suite.ctx, originDel)
		suite.mustPassInvariants()
	}

	// forcefully delegation shares as invalid
	{
		mutatedDel := originDel
		mutatedDel.Shares = mutatedDel.Shares.Sub(sdk.OneDec())
		suite.app.StakingKeeper.SetDelegation(suite.ctx, mutatedDel)
		_, broken = keeper.ChunksInvariant(suite.app.LiquidStakingKeeper)(suite.ctx)
		suite.True(broken, "delegation must have valid shares")
		// recover
		suite.app.StakingKeeper.SetDelegation(suite.ctx, originDel)
		suite.mustPassInvariants()
	}
	suite.ctx = suite.advanceEpoch(suite.ctx)
	suite.ctx = suite.advanceHeight(suite.ctx, 1, "")

	// 2: UNPAIRING CHUNK
	// first, create unpairing chunk
	insToBeWithdrawn, _, err := suite.app.LiquidStakingKeeper.DoWithdrawInsurance(
		suite.ctx,
		types.NewMsgWithdrawInsurance(env.insurances[2].ProviderAddress, env.insurances[2].Id),
	)
	suite.NoError(err)
	suite.ctx = suite.advanceEpoch(suite.ctx)
	suite.ctx = suite.advanceHeight(suite.ctx, 1, "start withdrawing insurance")

	origin, _ = suite.app.LiquidStakingKeeper.GetChunk(suite.ctx, insToBeWithdrawn.ChunkId)
	suite.checkUnpairingAndUnpairingForUnstakingChunks(suite.ctx, origin)

	// 3: PAIRING
	suite.ctx = suite.advanceEpoch(suite.ctx)
	suite.ctx = suite.advanceHeight(suite.ctx, 1, "unpairing finished")
	origin, _ = suite.app.LiquidStakingKeeper.GetChunk(suite.ctx, origin.Id)
	suite.Equal(
		types.CHUNK_STATUS_PAIRING, origin.Status,
		"after unpairing finished, chunk's status must be pairing",
	)
	// forcefully change paired insurance id of pairing chunk
	{
		mutated := origin
		mutated.PairedInsuranceId = 5
		suite.app.LiquidStakingKeeper.SetChunk(suite.ctx, mutated)
		_, broken = keeper.ChunksInvariant(suite.app.LiquidStakingKeeper)(suite.ctx)
		suite.True(broken, "pairing chunk must not have paired insurance id")
		// recover
		suite.app.LiquidStakingKeeper.SetChunk(suite.ctx, origin)
		suite.mustPassInvariants()
	}

	chunkBal := suite.app.BankKeeper.GetBalance(suite.ctx, origin.DerivedAddress(), suite.denom)
	suite.True(chunkBal.Amount.GTE(types.ChunkSize))
	// forcefully change chunk's balance
	{
		oneToken := sdk.NewCoins(sdk.NewCoin(suite.denom, sdk.OneInt()))
		suite.app.BankKeeper.SendCoins(
			suite.ctx,
			origin.DerivedAddress(),
			sdk.AccAddress(env.valAddrs[0]),
			oneToken,
		)
		_, broken = keeper.ChunksInvariant(suite.app.LiquidStakingKeeper)(suite.ctx)
		suite.True(broken, "chunk must have valid balance")
		// recover
		suite.app.BankKeeper.SendCoins(
			suite.ctx,
			sdk.AccAddress(env.valAddrs[0]),
			origin.DerivedAddress(),
			oneToken,
		)
		suite.mustPassInvariants()
	}

	// 4: UNPAIRING FOR UNSTAKING CHUNK
	// first, create unpairing for unstaking chunk
	toBeUnstakedChunks, _, err := suite.app.LiquidStakingKeeper.QueueLiquidUnstake(
		suite.ctx,
		types.NewMsgLiquidUnstake(
			env.delegators[0].String(),
			sdk.NewCoin(suite.denom, types.ChunkSize),
		),
	)
	suite.NoError(err)
	suite.ctx = suite.advanceEpoch(suite.ctx)
	suite.ctx = suite.advanceHeight(suite.ctx, 1, "unstaking chunk started")

	origin, _ = suite.app.LiquidStakingKeeper.GetChunk(suite.ctx, toBeUnstakedChunks[0].Id)
	suite.checkUnpairingAndUnpairingForUnstakingChunks(suite.ctx, origin)
}

func (suite *KeeperTestSuite) checkUnpairingAndUnpairingForUnstakingChunks(
	ctx sdk.Context,
	origin types.Chunk,
) {
	// forcefully change status of chunk as invalid
	{
		mutated := origin
		mutated.UnpairingInsuranceId = types.Empty
		suite.app.LiquidStakingKeeper.SetChunk(suite.ctx, mutated)
		_, broken := keeper.ChunksInvariant(suite.app.LiquidStakingKeeper)(suite.ctx)
		suite.True(broken, "unpairing chunk must have valid unpairing insurance id")
		// recover
		suite.app.LiquidStakingKeeper.SetChunk(suite.ctx, origin)
		suite.mustPassInvariants()
	}

	originIns, _ := suite.app.LiquidStakingKeeper.GetInsurance(ctx, origin.UnpairingInsuranceId)
	// forcefully delete unpairing insurance
	{
		suite.app.LiquidStakingKeeper.DeleteInsurance(ctx, originIns.Id)
		_, broken := keeper.ChunksInvariant(suite.app.LiquidStakingKeeper)(ctx)
		suite.True(broken, "unpairing insurance must exist in store")
		// recover
		suite.app.LiquidStakingKeeper.SetInsurance(ctx, originIns)
		suite.mustPassInvariants()
	}

	ubd, _ := suite.app.StakingKeeper.GetUnbondingDelegation(ctx, origin.DerivedAddress(), originIns.GetValidator())
	// forcefully delete unbonding delegation obj of unpairing chunk
	{
		suite.app.StakingKeeper.RemoveUnbondingDelegation(ctx, ubd)
		_, broken := keeper.ChunksInvariant(suite.app.LiquidStakingKeeper)(ctx)
		suite.True(broken, "unbonding delegation must exist in store")
		// recover
		suite.app.StakingKeeper.SetUnbondingDelegation(ctx, ubd)
		suite.mustPassInvariants()
	}

	// forcefully add unbonding entry
	{
		ubd.Entries = append(ubd.Entries, ubd.Entries[0])
		suite.app.StakingKeeper.SetUnbondingDelegation(ctx, ubd)
		_, broken := keeper.ChunksInvariant(suite.app.LiquidStakingKeeper)(ctx)
		suite.True(broken, "chunk's unbonding delegation must have one entry")
		// recover
		ubd.Entries = ubd.Entries[:len(ubd.Entries)-1]
		suite.app.StakingKeeper.SetUnbondingDelegation(ctx, ubd)
		suite.mustPassInvariants()
	}

	// forcefully change initial balance of unbonding entry
	{
		ubd.Entries[0].InitialBalance = ubd.Entries[0].InitialBalance.Sub(sdk.OneInt())
		suite.app.StakingKeeper.SetUnbondingDelegation(ctx, ubd)
		_, broken := keeper.ChunksInvariant(suite.app.LiquidStakingKeeper)(ctx)
		suite.True(broken, "chunk's unbonding delegation's entry must have valid initial balance")
		// recover
		ubd.Entries[0].InitialBalance = ubd.Entries[0].InitialBalance.Add(sdk.OneInt())
		suite.app.StakingKeeper.SetUnbondingDelegation(ctx, ubd)
		suite.mustPassInvariants()
	}
}
