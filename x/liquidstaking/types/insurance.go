package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"sort"
)

const (
	// 5%
	SlashFractionInt  = 5
	SlashFractionPrec = 2
)

func NewInsurance(id uint64, providerAddress, validatorAddress string, feeRate sdk.Dec) Insurance {
	return Insurance{
		Id:               id,
		ChunkId:          0, // Not yet assigned
		Status:           INSURANCE_STATUS_PAIRING,
		ProviderAddress:  providerAddress,
		ValidatorAddress: validatorAddress,
	}
}

func (i *Insurance) DerivedAddress() sdk.AccAddress {
	return DeriveAddress(ModuleName, fmt.Sprintf("insurance%d", i.Id))
}

func (i *Insurance) FeePoolAddress() sdk.AccAddress {
	return DeriveAddress(ModuleName, fmt.Sprintf("insurancefee%d", i.Id))
}

func SortInsurances(validatorMap map[string]stakingtypes.Validator, insurances []Insurance) {
	sort.Slice(insurances, func(i, j int) bool {
		iInsurance := insurances[i]
		jInsurance := insurances[j]

		iValidator := validatorMap[iInsurance.ValidatorAddress]
		jValidator := validatorMap[jInsurance.ValidatorAddress]

		iFee := iValidator.Commission.Rate.Add(iInsurance.FeeRate)
		jFee := jValidator.Commission.Rate.Add(jInsurance.FeeRate)

		if !iFee.Equal(jFee) {
			return iFee.LT(jFee)
		}
		return iInsurance.Id < jInsurance.Id
	})
}
