package types

import sdk "github.com/cosmos/cosmos-sdk/types"

func NewRedelegationInfo(id uint64) RedelegationInfo {
	return RedelegationInfo{
		ChunkId: id,
		Covered: sdk.ZeroInt(),
	}
}
