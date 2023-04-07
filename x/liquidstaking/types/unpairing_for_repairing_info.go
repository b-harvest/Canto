package types

import "time"

func NewUnpairingForRepairingInfo(
	chunkId, insuranceId uint64,
	completionTime time.Time,
) UnpairingForRepairingInfo {
	return UnpairingForRepairingInfo{
		ChunkId:        chunkId,
		InsuranceId:    insuranceId,
		CompletionTime: completionTime,
	}
}
