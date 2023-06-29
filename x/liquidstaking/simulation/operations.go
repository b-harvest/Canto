package simulation

import (
	"math/rand"

	"github.com/Canto-Network/Canto/v6/app"
	"github.com/Canto-Network/Canto/v6/x/liquidstaking/keeper"
	"github.com/Canto-Network/Canto/v6/x/liquidstaking/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"
)

// Simulation operation weights constants.
const (
	OpWeightMsgLiquidStake                 = "op_weight_msg_liquid_stake"
	OpWeightMsgLiquidUnstake               = "op_weight_msg_liquid_unstake"
	OpWeightMsgProvideInsurance            = "op_weight_msg_provide_insurance"
	OpWeightMsgCancelProvideInsurance      = "op_weight_msg_cancel_provide_insurance"
	OpWeightMsgDepositInsurance            = "op_weight_msg_deposit_insurance"
	OpWeightMsgWithdrawInsurance           = "op_weight_msg_withdraw_insurance"
	OpWeightMsgWithdrawInsuranceCommission = "op_weight_msg_withdraw_insurance_commission"
	OpWeightMsgClaimDiscountedReward       = "op_weight_msg_claim_discounted_reward"
)

var (
	Gas  = uint64(20000000)
	Fees = sdk.Coins{
		{
			Denom:  sdk.DefaultBondDenom,
			Amount: sdk.NewInt(0),
		},
	}
)

// WeightedOperations returns all the operations from the module with their respective weights
func WeightedOperations(appParams simtypes.AppParams, cdc codec.JSONCodec, ak types.AccountKeeper, bk types.BankKeeper, sk types.StakingKeeper, k keeper.Keeper) simulation.WeightedOperations {
	var weightMsgLiquidStake int
	appParams.GetOrGenerate(cdc, OpWeightMsgLiquidStake, &weightMsgLiquidStake, nil, func(_ *rand.Rand) {
		weightMsgLiquidStake = app.DefaultWeightMsgLiquidStake
	})

	var weightMsgLiquidUnstake int
	appParams.GetOrGenerate(cdc, OpWeightMsgLiquidUnstake, &weightMsgLiquidUnstake, nil, func(_ *rand.Rand) {
		weightMsgLiquidUnstake = app.DefaultWeightMsgLiquidUnstake
	})

	var weightMsgProvideInsurance int
	appParams.GetOrGenerate(cdc, OpWeightMsgProvideInsurance, &weightMsgProvideInsurance, nil, func(_ *rand.Rand) {
		weightMsgProvideInsurance = app.DefaultWeightMsgProvideInsurance
	})

	var weightMsgCancelProvideInsurance int
	appParams.GetOrGenerate(cdc, OpWeightMsgCancelProvideInsurance, &weightMsgCancelProvideInsurance, nil, func(_ *rand.Rand) {
		weightMsgCancelProvideInsurance = app.DefaultWeightMsgCancelProvideInsurance
	})

	var weightMsgDepositInsurance int
	appParams.GetOrGenerate(cdc, OpWeightMsgDepositInsurance, &weightMsgDepositInsurance, nil, func(_ *rand.Rand) {
		weightMsgDepositInsurance = app.DefaultWeightMsgDepositInsurance
	})

	var weightMsgWithdrawInsurance int
	appParams.GetOrGenerate(cdc, OpWeightMsgWithdrawInsurance, &weightMsgWithdrawInsurance, nil, func(_ *rand.Rand) {
		weightMsgWithdrawInsurance = app.DefaultWeightMsgWithdrawInsurance
	})

	var weightMsgWithdrawInsuranceCommission int
	appParams.GetOrGenerate(cdc, OpWeightMsgWithdrawInsuranceCommission, &weightMsgWithdrawInsuranceCommission, nil, func(_ *rand.Rand) {
		weightMsgWithdrawInsuranceCommission = app.DefaultWeightMsgWithdrawInsuranceCommission
	})

	var weightMsgClaimDiscountedReward int
	appParams.GetOrGenerate(cdc, OpWeightMsgClaimDiscountedReward, &weightMsgClaimDiscountedReward, nil, func(_ *rand.Rand) {
		weightMsgClaimDiscountedReward = app.DefaultWeightMsgClaimDiscountedReward
	})

	return simulation.WeightedOperations{
		simulation.NewWeightedOperation(
			weightMsgLiquidStake,
			SimulateMsgLiquidStake(ak, bk),
		),
		simulation.NewWeightedOperation(
			weightMsgLiquidUnstake,
			SimulateMsgLiquidUnstake(ak, bk, k),
		),
		simulation.NewWeightedOperation(
			weightMsgProvideInsurance,
			SimulateMsgProvideInsurance(ak, bk, sk),
		),
	}
}

// SimulateMsgLiquidStake generates a MsgLiquidStake with random values.
func SimulateMsgLiquidStake(ak types.AccountKeeper, bk types.BankKeeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		simAccount, _ := simtypes.RandomAcc(r, accs)
		account := ak.GetAccount(ctx, simAccount.Address)
		delegator := account.GetAddress()
		spendable := bk.SpendableCoins(ctx, delegator)

		chunksToLiquidStake := int64(simtypes.RandIntBetween(r, 1, 3))
		stakingCoins := sdk.NewCoins(
			sdk.NewCoin(
				sdk.DefaultBondDenom,
				types.ChunkSize.MulRaw(chunksToLiquidStake),
			),
		)
		if !spendable.AmountOf(sdk.DefaultBondDenom).GTE(stakingCoins[0].Amount) {
			if err := bk.MintCoins(ctx, types.ModuleName, stakingCoins); err != nil {
				panic(err)
			}
			if err := bk.SendCoinsFromModuleToAccount(ctx, types.ModuleName, delegator, stakingCoins); err != nil {
				panic(err)
			}
			spendable = bk.SpendableCoins(ctx, delegator)
		}

		msg := types.NewMsgLiquidStake(delegator.String(), stakingCoins[0])
		txCtx := simulation.OperationInput{
			R:               r,
			App:             app,
			TxGen:           simapp.MakeTestEncodingConfig().TxConfig,
			Cdc:             nil,
			Msg:             msg,
			MsgType:         msg.Type(),
			CoinsSpentInMsg: spendable,
			Context:         ctx,
			SimAccount:      simAccount,
			AccountKeeper:   ak,
			Bankkeeper:      bk,
			ModuleName:      types.ModuleName,
		}
		return types.GenAndDeliverTxWithFees(txCtx, Gas, Fees)
	}
}

// SimulateMsgLiquidUnstake generates a MsgLiquidUnstake with random values.
func SimulateMsgLiquidUnstake(ak types.AccountKeeper, bk types.BankKeeper, k keeper.Keeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		simAccount, _ := simtypes.RandomAcc(r, accs)
		account := ak.GetAccount(ctx, simAccount.Address)
		delegator := account.GetAddress()
		spendable := bk.SpendableCoins(ctx, delegator)

		nas := k.GetNetAmountState(ctx)
		chunksToLiquidStake := int64(simtypes.RandIntBetween(r, 1, 3))
		unstakingCoin := sdk.NewCoin(
			sdk.DefaultBondDenom,
			types.ChunkSize.MulRaw(chunksToLiquidStake),
		)
		// mustHaveCoin means ls tokens to successfully liquid unstake the coin.
		mustHaveCoins := sdk.NewCoins(
			sdk.NewCoin(
				types.DefaultLiquidBondDenom,
				unstakingCoin.Amount.ToDec().Mul(nas.MintRate).Ceil().TruncateInt(),
			),
		)
		if !spendable.AmountOf(types.DefaultLiquidBondDenom).GTE(mustHaveCoins[0].Amount) {
			if err := bk.MintCoins(ctx, types.ModuleName, mustHaveCoins); err != nil {
				panic(err)
			}
			if err := bk.SendCoinsFromModuleToAccount(ctx, types.ModuleName, delegator, mustHaveCoins); err != nil {
				panic(err)
			}
			spendable = bk.SpendableCoins(ctx, delegator)
		}

		msg := types.NewMsgLiquidUnstake(delegator.String(), unstakingCoin)
		txCtx := simulation.OperationInput{
			R:               r,
			App:             app,
			TxGen:           simapp.MakeTestEncodingConfig().TxConfig,
			Cdc:             nil,
			Msg:             msg,
			MsgType:         msg.Type(),
			CoinsSpentInMsg: spendable,
			Context:         ctx,
			SimAccount:      simAccount,
			AccountKeeper:   ak,
			Bankkeeper:      bk,
			ModuleName:      types.ModuleName,
		}
		return types.GenAndDeliverTxWithFees(txCtx, Gas, Fees)
	}
}

func SimulateMsgProvideInsurance(ak types.AccountKeeper, bk types.BankKeeper, sk types.StakingKeeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		simAccount, _ := simtypes.RandomAcc(r, accs)
		account := ak.GetAccount(ctx, simAccount.Address)
		provider := account.GetAddress()
		spendable := bk.SpendableCoins(ctx, provider)

		upperThanMinimumCollateral := simtypes.RandomDecAmount(r, sdk.MustNewDecFromStr("0.03"))
		minCollateral := sdk.MustNewDecFromStr(types.MinimumCollateral)
		minCollateral = minCollateral.Add(upperThanMinimumCollateral)
		collaterals := sdk.NewCoins(
			sdk.NewCoin(
				sdk.DefaultBondDenom,
				minCollateral.Ceil().TruncateInt(),
			),
		)
		feeRate := simtypes.RandomDecAmount(r, sdk.MustNewDecFromStr("0.15"))

		if !spendable.AmountOf(sdk.DefaultBondDenom).GTE(collaterals[0].Amount) {
			if err := bk.MintCoins(ctx, types.ModuleName, collaterals); err != nil {
				panic(err)
			}
			if err := bk.SendCoinsFromModuleToAccount(ctx, types.ModuleName, provider, collaterals); err != nil {
				panic(err)
			}
			spendable = bk.SpendableCoins(ctx, provider)
		}

		validators := sk.GetAllValidators(ctx)
		// select one validator randomly
		validator := validators[r.Intn(len(validators))]

		msg := types.NewMsgProvideInsurance(provider.String(), validator.GetOperator().String(), collaterals[0], feeRate)
		txCtx := simulation.OperationInput{
			R:               r,
			App:             app,
			TxGen:           simapp.MakeTestEncodingConfig().TxConfig,
			Cdc:             nil,
			Msg:             msg,
			MsgType:         msg.Type(),
			CoinsSpentInMsg: spendable,
			Context:         ctx,
			SimAccount:      simAccount,
			AccountKeeper:   ak,
			Bankkeeper:      bk,
			ModuleName:      types.ModuleName,
		}
		return types.GenAndDeliverTxWithFees(txCtx, Gas, Fees)
	}
}
