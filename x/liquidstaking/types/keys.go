package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// constants
const (
	// module name
	ModuleName = "liquidstaking"

	// StoreKey to be used when creating the KVStore
	StoreKey = ModuleName

	// RouterKey to be used for message routing
	RouterKey = ModuleName
)

// prefix bytes for the inflation persistent store
const (
	prefixChunkBondRequestId = iota + 1
	prefixChunkUnbondRequestId
	prefixInsuranceBidId
	prefixAliveChunkId
	prefixUnbondingChunkId

	prefixChunkBondRequest
	prefixChunkUnbondRequest
	prefixInsuranceBid
	prefixInsuranceUnbondRequest
	prefixAliveChunk
	prefixUnbondingChunk
)

// KVStore key prefixes
var (
	KeyPrefixChunkBondRequestId   = []byte{prefixChunkBondRequestId}
	KeyPrefixChunkUnbondRequestId = []byte{prefixChunkUnbondRequestId}
	KeyPrefixInsuranceBidId       = []byte{prefixInsuranceBidId}
	KeyPrefixAliveChunkId         = []byte{prefixAliveChunkId}
	KeyPrefixUnbondingChunkId     = []byte{prefixUnbondingChunkId}

	KeyPrefixChunkBondRequest       = []byte{prefixChunkBondRequest}
	KeyPrefixChunkUnbondRequest     = []byte{prefixChunkUnbondRequest}
	KeyPrefixInsuranceBid           = []byte{prefixInsuranceBid}
	KeyPrefixInsuranceUnbondRequest = []byte{prefixInsuranceUnbondRequest}
	KeyPrefixAliveChunk             = []byte{prefixAliveChunk}
	KeyPrefixUnbondingChunk         = []byte{prefixUnbondingChunk}
)

func GetChunkBondRequestKey(id ChunkBondRequestId) []byte {
	return append(KeyPrefixChunkBondRequest, sdk.Uint64ToBigEndian(id)...)
}

func GetChunkUnbondRequestKey(id ChunkUnbondRequestId) []byte {
	return append(KeyPrefixChunkUnbondRequest, sdk.Uint64ToBigEndian(id)...)
}

func GetInsuranceBidKey(id InsuranceBidId) []byte {
	return append(KeyPrefixInsuranceBid, sdk.Uint64ToBigEndian(id)...)
}

func GetInsuranceUnbondRequestKey(id AliveChunkId) []byte {
	return append(KeyPrefixInsuranceUnbondRequest, sdk.Uint64ToBigEndian(id)...)
}

func GetAliveChunkKey(id AliveChunkId) []byte {
	return append(KeyPrefixAliveChunk, sdk.Uint64ToBigEndian(id)...)
}

func GetUnbondingChunkKey(id UnbondingChunkId) []byte {
	return append(KeyPrefixUnbondingChunk, sdk.Uint64ToBigEndian(id)...)
}
