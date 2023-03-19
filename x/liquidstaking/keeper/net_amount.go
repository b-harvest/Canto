package keeper

import (
	"github.com/Canto-Network/Canto/v6/x/liquidstaking/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// TODO: 테스트할 때 실질적으로 어떤 수치들을 가지고 테스트해야할지는 태영님과 논의
func (k Keeper) GetNetAmountState(ctx sdk.Context) (nas types.NetAmountState) {
	bondDenom := k.stakingKeeper.BondDenom(ctx)
	totalDelShares := sdk.ZeroDec()
	totalChunksBalance := sdk.NewDec(0)
	totalRewards := sdk.ZeroDec()
	totalLiquidTokens := sdk.ZeroInt()
	totalInsuranceTokens := sdk.ZeroInt()
	totalUnbondingBalance := sdk.ZeroDec()

	err := k.IterateAllChunks(ctx, func(chunk types.Chunk) (stop bool, err error) {
		balance := k.bankKeeper.GetBalance(ctx, chunk.DerivedAddress(), k.stakingKeeper.BondDenom(ctx))
		totalChunksBalance = totalChunksBalance.Add(balance.Amount.ToDec())

		insurance, _ := k.GetInsurance(ctx, chunk.InsuranceId)
		validator := k.stakingKeeper.Validator(ctx, sdk.ValAddress(insurance.ValidatorAddress))
		delegation, found := k.stakingKeeper.GetDelegation(ctx, chunk.DerivedAddress(), sdk.ValAddress(insurance.ValidatorAddress))
		if !found {
			return false, nil
		}
		totalDelShares.Add(delegation.GetShares())
		tokens := validator.TokensFromSharesTruncated(delegation.GetShares()).TruncateInt()
		totalLiquidTokens = totalLiquidTokens.Add(tokens)
		cachedCtx, _ := ctx.CacheContext()
		endingPeriod := k.distributionKeeper.IncrementValidatorPeriod(cachedCtx, validator)
		delReward := k.distributionKeeper.CalculateDelegationRewards(cachedCtx, validator, delegation, endingPeriod)
		totalRewards = totalRewards.Add(delReward.AmountOf(bondDenom))

		ubds := k.stakingKeeper.GetAllUnbondingDelegations(ctx, chunk.DerivedAddress())
		for _, ubd := range ubds {
			for _, entry := range ubd.Entries {
				totalUnbondingBalance = totalUnbondingBalance.Add(entry.Balance.ToDec())
			}
		}
		return false, nil
	})
	if err != nil {
		panic(err)
	}

	// Iterate all paired insurances to get total insurance tokens
	err = k.IterateAllInsurances(ctx, func(insurance types.Insurance) (stop bool, err error) {
		if insurance.Status == types.INSURANCE_STATUS_PAIRED {
			insuranceBalance := k.bankKeeper.GetBalance(ctx, insurance.DerivedAddress(), k.stakingKeeper.BondDenom(ctx))
			// TODO: Should we add insurance fee to total insurance tokens?
			totalInsuranceTokens = totalInsuranceTokens.Add(insuranceBalance.Amount)
		}
		return false, nil
	})
	if err != nil {
		panic(err)
	}

	nas = types.NetAmountState{
		LsTokensTotalSupply:   k.bankKeeper.GetSupply(ctx, bondDenom).Amount,
		TotalChunksBalance:    totalChunksBalance.TruncateInt(),
		TotalDelShares:        totalDelShares,
		TotalRemainingRewards: totalRewards,
		TotalLiquidTokens:     totalLiquidTokens,
		TotalInsuranceTokens:  totalInsuranceTokens,
	}

	nas.NetAmount = nas.CalcNetAmount()
	nas.MintRate = nas.CalcMintRate()
	return
}
