package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (c *Chunk) DerivedAddress() sdk.AccAddress {
	return DeriveAddress(ModuleName, fmt.Sprint("chunk%d", c.Id))
}
