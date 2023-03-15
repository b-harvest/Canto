package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewInsurance(id uint64) Insurance {
	return Insurance{
		Id:      id,
		ChunkId: 0, // Not yet assigned
		Status:  INSURANCE_STATUS_PAIRING,
	}
}

func (i *Insurance) DerivedAddress() sdk.AccAddress {
	return DeriveAddress(ModuleName, fmt.Sprint("insurance%d", i.Id))
}

func (i *Insurance) FeePoolAddress() sdk.AccAddress {
	return DeriveAddress(ModuleName, fmt.Sprint("insurancefee%d", i.Id))
}
