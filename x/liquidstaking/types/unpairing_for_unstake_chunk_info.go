package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func NewUnpairingForUnstakeChunkInfo(
	chunkId uint64,
	delegatorAddress string,
	escrowedLsTokens sdk.Coin,
) UnpairingForUnstakeChunkInfo {
	return UnpairingForUnstakeChunkInfo{
		ChunkId:          chunkId,
		DelegatorAddress: delegatorAddress,
		EscrowedLstokens: escrowedLsTokens,
	}
}

func (info *UnpairingForUnstakeChunkInfo) GetDelegator() sdk.AccAddress {
	return sdk.MustAccAddressFromBech32(info.DelegatorAddress)
}

func (info *UnpairingForUnstakeChunkInfo) Validate(chunkMap map[uint64]Chunk) error {
	chunk, ok := chunkMap[info.ChunkId]
	if !ok {
		return sdkerrors.Wrapf(
			ErrNotFoundUnpairingForUnstakeChunkInfoChunkId,
			"chunk id: %d",
			info.ChunkId,
		)
	}
	if chunk.Status != CHUNK_STATUS_UNPAIRING_FOR_UNSTAKE {
		return ErrInvalidChunkStatus
	}
	return nil
}
