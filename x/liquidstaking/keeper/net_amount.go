package keeper

import (
	"github.com/Canto-Network/Canto/v6/x/liquidstaking/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// TODO: Include Pairing Chunk Balances
func (k Keeper) GetNetAmountState(ctx sdk.Context) (nas types.NetAmountState) {
	liquidBondDenom := k.GetLiquidBondDenom(ctx)
	bondDenom := k.stakingKeeper.BondDenom(ctx)
	totalDelShares := sdk.ZeroDec()
	totalChunksBalance := sdk.NewDec(0)
	totalRemainingRewards := sdk.ZeroDec()
	totalRemainingInsuranceCommissions := sdk.ZeroDec()
	totalLiquidTokens := sdk.ZeroInt()
	totalInsuranceTokens := sdk.ZeroInt()
	totalPairedInsuranceTokens := sdk.ZeroInt()
	totalUnpairingInsuranceTokens := sdk.ZeroInt()
	totalUnbondingChunksBalance := sdk.ZeroDec()
	numPairedChunks := sdk.ZeroInt()

	moduleFeeRate, utilizationRatio := k.CalcDynamicFeeRate(ctx)
	discountRate := k.CalcDiscountRate(ctx)
	k.IterateAllChunks(ctx, func(chunk types.Chunk) (stop bool) {
		balance := k.bankKeeper.GetBalance(ctx, chunk.DerivedAddress(), k.stakingKeeper.BondDenom(ctx))
		totalChunksBalance = totalChunksBalance.Add(balance.Amount.ToDec())

		switch chunk.Status {
		case types.CHUNK_STATUS_PAIRED:
			numPairedChunks = numPairedChunks.Add(sdk.OneInt())
			// chunk is paired which means have delegation
			pairedIns, _ := k.GetInsurance(ctx, chunk.PairedInsuranceId)
			valAddr, err := sdk.ValAddressFromBech32(pairedIns.ValidatorAddress)
			if err != nil {
				panic(err)
			}
			validator := k.stakingKeeper.Validator(ctx, valAddr)
			delegation, found := k.stakingKeeper.GetDelegation(ctx, chunk.DerivedAddress(), valAddr)
			if !found {
				return false
			}
			totalDelShares = totalDelShares.Add(delegation.GetShares())
			tokenValue := validator.TokensFromSharesTruncated(delegation.GetShares()).TruncateInt()
			penaltyAmt := types.ChunkSize.Sub(tokenValue)
			// If penaltyAmt > 0 and paired insurance can cover it, then token value is same with ChunkSize
			if penaltyAmt.IsPositive() {
				pairedInsBal := k.bankKeeper.GetBalance(ctx, pairedIns.DerivedAddress(), liquidBondDenom)
				if pairedInsBal.Amount.LT(penaltyAmt) {
					penaltyAmt = penaltyAmt.Sub(pairedInsBal.Amount)
				} else {
					penaltyAmt = sdk.ZeroInt()
				}
				tokenValue = types.ChunkSize.Sub(penaltyAmt)
			}
			totalLiquidTokens = totalLiquidTokens.Add(tokenValue)
			cachedCtx, _ := ctx.CacheContext()
			endingPeriod := k.distributionKeeper.IncrementValidatorPeriod(cachedCtx, validator)
			delRewards := k.distributionKeeper.CalculateDelegationRewards(cachedCtx, validator, delegation, endingPeriod)
			// chunk's remaining reward is calculated by
			// 1. rest = del_reward - insurance_commission
			// 2. remaining = rest x (1 - module_fee_rate)
			// 3. result = remaining x (1 - discount_rate)
			delReward := delRewards.AmountOf(bondDenom)
			insuranceCommission := delReward.Mul(pairedIns.FeeRate)
			restReward := delReward.Sub(insuranceCommission)
			remainingReward := restReward.Mul(sdk.OneDec().Sub(moduleFeeRate))
			totalRemainingRewards = totalRemainingRewards.Add(
				remainingReward.Mul(
					sdk.OneDec().Sub(discountRate),
				),
			)
		default:
			k.stakingKeeper.IterateDelegatorUnbondingDelegations(ctx, chunk.DerivedAddress(), func(ubd stakingtypes.UnbondingDelegation) (stop bool) {
				for _, entry := range ubd.Entries {
					totalUnbondingChunksBalance = totalUnbondingChunksBalance.Add(entry.Balance.ToDec())
				}
				return false
			})
		}
		return false
	})

	// Iterate all paired insurances to get total insurance tokens
	k.IterateAllInsurances(ctx, func(insurance types.Insurance) (stop bool) {
		insuranceBalance := k.bankKeeper.GetBalance(ctx, insurance.DerivedAddress(), bondDenom)
		commission := k.bankKeeper.GetBalance(ctx, insurance.FeePoolAddress(), bondDenom)
		switch insurance.Status {
		case types.INSURANCE_STATUS_PAIRED:
			totalPairedInsuranceTokens = totalPairedInsuranceTokens.Add(insuranceBalance.Amount)
		case types.INSURANCE_STATUS_UNPAIRING:
			totalUnpairingInsuranceTokens = totalUnpairingInsuranceTokens.Add(insuranceBalance.Amount)
		case types.INSURANCE_STATUS_UNPAIRED:
		}
		totalInsuranceTokens = totalInsuranceTokens.Add(insuranceBalance.Amount)
		totalRemainingInsuranceCommissions = totalRemainingInsuranceCommissions.Add(commission.Amount.ToDec())
		return false
	})

	nas = types.NetAmountState{
		LsTokensTotalSupply:                k.bankKeeper.GetSupply(ctx, liquidBondDenom).Amount,
		TotalLiquidTokens:                  totalLiquidTokens,
		TotalChunksBalance:                 totalChunksBalance.TruncateInt(),
		TotalDelShares:                     totalDelShares,
		TotalRemainingRewards:              totalRemainingRewards,
		TotalUnbondingChunksBalance:        totalUnbondingChunksBalance.TruncateInt(),
		NumPairedChunks:                    numPairedChunks,
		ChunkSize:                          types.ChunkSize,
		TotalRemainingInsuranceCommissions: totalRemainingInsuranceCommissions,
		TotalInsuranceTokens:               totalInsuranceTokens,
		TotalPairedInsuranceTokens:         totalPairedInsuranceTokens,
		TotalUnpairingInsuranceTokens:      totalUnpairingInsuranceTokens,
	}
	nas.NetAmount = nas.CalcNetAmount(k.bankKeeper.GetBalance(ctx, types.RewardPool, bondDenom).Amount)
	nas.MintRate = nas.CalcMintRate()
	nas.RewardModuleAccBalance = k.bankKeeper.GetBalance(ctx, types.RewardPool, bondDenom).Amount
	nas.FeeRate, nas.UtilizationRatio = moduleFeeRate, utilizationRatio
	nas.RemainingChunkSlots = k.GetAvailableChunkSlots(ctx)
	nas.DiscountRate = discountRate
	return
}
