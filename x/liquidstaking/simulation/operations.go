package simulation

import (
	"github.com/Canto-Network/Canto/v6/app"
	"github.com/Canto-Network/Canto/v6/x/liquidstaking/keeper"
	"github.com/Canto-Network/Canto/v6/x/liquidstaking/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"
	"math/rand"
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
func WeightedOperations(appParams simtypes.AppParams, cdc codec.JSONCodec, ak types.AccountKeeper, bk types.BankKeeper, k keeper.Keeper) simulation.WeightedOperations {
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
			SimulateMsgLiquidStake(ak, bk, k),
		),
	}
}

func SimulateMsgLiquidStake(ak types.AccountKeeper, bk types.BankKeeper, k keeper.Keeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		simAccount, _ := simtypes.RandomAcc(r, accs)
		account := ak.GetAccount(ctx, simAccount.Address)
		delegator := account.GetAddress()
		spendable := bk.SpendableCoins(ctx, delegator)

		amount := msg.Amount

		if err = k.ShouldBeBondDenom(ctx, amount.Denom); err != nil {
			return
		}
		// Liquid stakers can send amount of tokens corresponding multiple of chunk size and create multiple chunks
		if err = k.ShouldBeMultipleOfChunkSize(amount.Amount); err != nil {
			return
		}
		chunksToCreate := amount.Amount.Quo(types.ChunkSize).Int64()

		availableChunkSlots := k.GetAvailableChunkSlots(ctx).Int64()
		if (availableChunkSlots - chunksToCreate) < 0 {
			err = sdkerrors.Wrapf(
				types.ErrExceedAvailableChunks,
				"requested chunks to create: %d, available chunks: %d",
				chunksToCreate,
				availableChunkSlots,
			)
			return
		}

		pairingInsurances, validatorMap := k.getPairingInsurances(ctx)
		if chunksToCreate > int64(len(pairingInsurances)) {
			err = types.ErrNoPairingInsurance
			return
		}

		nas := k.GetNetAmountState(ctx)
		types.SortInsurances(validatorMap, pairingInsurances, false)
		totalNewShares := sdk.ZeroDec()
		totalLsTokenMintAmount := sdk.ZeroInt()
		for i := int64(0); i < chunksToCreate; i++ {
			cheapestInsurance := pairingInsurances[0]
			pairingInsurances = pairingInsurances[1:]

			// Now we have the cheapest pairing insurance and valid msg liquid stake! Let's create a chunk
			// Create a chunk
			chunkId := k.getNextChunkIdWithUpdate(ctx)
			chunk := types.NewChunk(chunkId)

			// Escrow liquid staker's coins
			if err = k.bankKeeper.SendCoins(
				ctx,
				delAddr,
				chunk.DerivedAddress(),
				sdk.NewCoins(sdk.NewCoin(amount.Denom, types.ChunkSize)),
			); err != nil {
				return
			}
			validator := validatorMap[cheapestInsurance.ValidatorAddress]

			// Delegate to the validator
			// Delegator: DerivedAddress(chunk.Id)
			// Validator: insurance.ValidatorAddress
			// Amount: msg.Amount
			chunk, cheapestInsurance, newShares, err = k.pairChunkAndInsurance(
				ctx,
				chunk,
				cheapestInsurance,
				validator,
			)
			if err != nil {
				return
			}
			totalNewShares = totalNewShares.Add(newShares)

			liquidBondDenom := k.GetLiquidBondDenom(ctx)
			// Mint the liquid staking token
			lsTokenMintAmount = amount.Amount
			if nas.LsTokensTotalSupply.IsPositive() {
				lsTokenMintAmount = types.NativeTokenToLiquidStakeToken(amount.Amount, nas.LsTokensTotalSupply, nas.NetAmount)
			}
			if !lsTokenMintAmount.IsPositive() {
				err = sdkerrors.Wrapf(types.ErrInvalidAmount, "amount must be greater than or equal to %s", amount.String())
				return
			}
			mintedCoin := sdk.NewCoin(liquidBondDenom, lsTokenMintAmount)
			if err = k.bankKeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(mintedCoin)); err != nil {
				return
			}
			totalLsTokenMintAmount = totalLsTokenMintAmount.Add(lsTokenMintAmount)
			if err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, delAddr, sdk.NewCoins(mintedCoin)); err != nil {
				return
			}
			chunks = append(chunks, chunk)
		}
		return
	}
}
