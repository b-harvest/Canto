package simulation

import (
	"bytes"
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/kv"
)

// NewDecodeStore returns a decoder function closure that unmarshals the KVPair's
// Value to the corresponding distribution type.
func NewDecodeStore(cdc codec.Codec) func(kvA, kvB kv.Pair) string {
	return func(kvA, kvB kv.Pair) string {
		switch {
		case bytes.Equal(kvA.Key[:1], types.KeyPrefixLast):
			var feePoolA, feePoolB types.
			cdc.MustUnmarshal(kvA.Value, &feePoolA)
			cdc.MustUnmarshal(kvB.Value, &feePoolB)
			return fmt.Sprintf("%v\n%v", feePoolA, feePoolB)

		default:
			panic(fmt.Sprintf("invalid liquidstaking key prefix %X", kvA.Key[:1]))
		}
	}
}
