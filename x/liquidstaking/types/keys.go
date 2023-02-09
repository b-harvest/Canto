package types

import (
	"bytes"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
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
	prefixChunkBondRequestIndex
	prefixChunkUnbondRequest
	prefixChunkUnbondRequestIndex
	prefixInsuranceBid
	prefixInsuranceBidIndexByInsuranceProvider
	prefixInsuranceBidIndexByValidator
	prefixInsuranceUnbondRequest
	prefixInsuranceUnbondRequestIndexByInsuranceProvider
	prefixAliveChunk
	prefixAliveChunkIndexByInsuranceProvider
	prefixAliveChunkIndexByValidator
	prefixUnbondingChunk
)

// KVStore key prefixes
var (
	KeyPrefixChunkBondRequestId   = []byte{prefixChunkBondRequestId}
	KeyPrefixChunkUnbondRequestId = []byte{prefixChunkUnbondRequestId}
	KeyPrefixInsuranceBidId       = []byte{prefixInsuranceBidId}
	KeyPrefixAliveChunkId         = []byte{prefixAliveChunkId}
	KeyPrefixUnbondingChunkId     = []byte{prefixUnbondingChunkId}

	KeyPrefixChunkBondRequest                               = []byte{prefixChunkBondRequest}
	KeyPrefixChunkBondRequestIndex                          = []byte{prefixChunkBondRequestIndex}
	KeyPrefixChunkUnbondRequest                             = []byte{prefixChunkUnbondRequest}
	KeyPrefixChunkUnbondRequestIndex                        = []byte{prefixChunkUnbondRequestIndex}
	KeyPrefixInsuranceBid                                   = []byte{prefixInsuranceBid}
	KeyPrefixInsuranceBidIndexByInsuranceProvider           = []byte{prefixInsuranceBidIndexByInsuranceProvider}
	KeyPrefixInsuranceBidIndexByValidator                   = []byte{prefixInsuranceBidIndexByValidator}
	KeyPrefixInsuranceUnbondRequest                         = []byte{prefixInsuranceUnbondRequest}
	KeyPrefixInsuranceUnbondRequestIndexByInsuranceProvider = []byte{prefixInsuranceUnbondRequestIndexByInsuranceProvider}
	KeyPrefixAliveChunk                                     = []byte{prefixAliveChunk}
	KeyPrefixAliveChunkIndexByInsuranceProvider             = []byte{prefixAliveChunkIndexByInsuranceProvider}
	KeyPrefixAliveChunkIndexByValidator                     = []byte{prefixAliveChunkIndexByValidator}
	KeyPrefixUnbondingChunk                                 = []byte{prefixUnbondingChunk}
)

func GetChunkBondRequestKey(id ChunkBondRequestId) []byte {
	return append(KeyPrefixChunkBondRequest, sdk.Uint64ToBigEndian(id)...)
}

func GetChunkBondRequestIndexPrefixKey(requesterAddr sdk.AccAddress) []byte {
	return append(KeyPrefixChunkBondRequestIndex, address.MustLengthPrefix(requesterAddr)...)
}

func GetChunkBondRequestIndexKey(requesterAddr sdk.AccAddress, id ChunkBondRequestId) []byte {
	return append(GetChunkBondRequestIndexPrefixKey(requesterAddr), sdk.Uint64ToBigEndian(id)...)
}

func GetChunkUnbondRequestKey(id ChunkUnbondRequestId) []byte {
	return append(KeyPrefixChunkUnbondRequest, sdk.Uint64ToBigEndian(id)...)
}

func GetChunkUnbondRequestIndexPrefixKey(requesterAddr sdk.AccAddress) []byte {
	return append(KeyPrefixChunkUnbondRequestIndex, address.MustLengthPrefix(requesterAddr)...)
}

func GetChunkUnbondRequestIndexKey(requesterAddr sdk.AccAddress, id ChunkUnbondRequestId) []byte {
	return append(GetChunkUnbondRequestIndexPrefixKey(requesterAddr), sdk.Uint64ToBigEndian(id)...)
}

func GetInsuranceBidKey(id InsuranceBidId) []byte {
	return append(KeyPrefixInsuranceBid, sdk.Uint64ToBigEndian(id)...)
}

func GetInsuranceBidIndexByInsuranceProviderPrefixKey(insuranceProvider sdk.AccAddress) []byte {
	return append(KeyPrefixInsuranceBidIndexByInsuranceProvider, address.MustLengthPrefix(insuranceProvider)...)
}

func GetInsuranceBidIndexByInsuranceProviderKey(insuranceProvider sdk.AccAddress, id InsuranceBidId) []byte {
	return append(GetInsuranceBidIndexByInsuranceProviderPrefixKey(insuranceProvider), sdk.Uint64ToBigEndian(id)...)
}

func GetInsuranceBidIndexByValidatorPrefixKey(validator sdk.ValAddress) []byte {
	return append(KeyPrefixInsuranceBidIndexByValidator, address.MustLengthPrefix(validator)...)
}

func GetInsuranceBidIndexByValidatorKey(validator sdk.ValAddress, id InsuranceBidId) []byte {
	return append(GetInsuranceBidIndexByValidatorPrefixKey(validator), sdk.Uint64ToBigEndian(id)...)
}

func GetInsuranceUnbondRequestKey(id AliveChunkId) []byte {
	return append(KeyPrefixInsuranceUnbondRequest, sdk.Uint64ToBigEndian(id)...)
}

func GetInsuranceUnbondRequestIndexByInsuranceProviderPrefixKey(insuranceProvider sdk.AccAddress) []byte {
	return append(KeyPrefixInsuranceUnbondRequestIndexByInsuranceProvider,
		address.MustLengthPrefix(insuranceProvider)...)
}

func GetInsuranceUnbondRequestIndexByInsuranceProviderKey(insuranceProvider sdk.AccAddress, id InsuranceBidId) []byte {
	return append(GetInsuranceUnbondRequestIndexByInsuranceProviderPrefixKey(insuranceProvider),
		sdk.Uint64ToBigEndian(id)...)
}

func GetAliveChunkKey(id AliveChunkId) []byte {
	return append(KeyPrefixAliveChunk, sdk.Uint64ToBigEndian(id)...)
}

func GetAliveChunkIndexByInsuranceProviderPrefixKey(insuranceProvider sdk.AccAddress) []byte {
	return append(KeyPrefixAliveChunkIndexByInsuranceProvider, address.MustLengthPrefix(insuranceProvider)...)
}

func GetAliveChunkIndexByInsuranceProviderKey(insuranceProvider sdk.AccAddress, id AliveChunkId) []byte {
	return append(GetAliveChunkIndexByInsuranceProviderPrefixKey(insuranceProvider), sdk.Uint64ToBigEndian(id)...)
}

func GetAliveChunkIndexByValidatorPrefixKey(validator sdk.ValAddress) []byte {
	return append(KeyPrefixAliveChunkIndexByValidator, address.MustLengthPrefix(validator)...)
}

func GetAliveChunkIndexByValidatorKey(validator sdk.ValAddress, id AliveChunkId) []byte {
	return append(GetAliveChunkIndexByValidatorPrefixKey(validator), sdk.Uint64ToBigEndian(id)...)
}

func GetUnbondingChunkKey(id UnbondingChunkId) []byte {
	return append(KeyPrefixUnbondingChunk, sdk.Uint64ToBigEndian(id)...)
}

func parseIndexByAccAddressKey(key, keyPrefix []byte) (accAddr sdk.AccAddress, id uint64) {
	if !bytes.HasPrefix(key, keyPrefix) {
		panic("key does not have proper prefix")
	}

	addrLen := key[1]
	accAddr = key[2 : 2+addrLen]
	id = sdk.BigEndianToUint64(key[2+addrLen:])
	return
}

func parseIndexByValAddressKey(key, keyPrefix []byte) (valAddr sdk.ValAddress, id uint64) {
	if !bytes.HasPrefix(key, keyPrefix) {
		panic("key does not have proper prefix")
	}

	addrLen := key[1]
	valAddr = key[2 : 2+addrLen]
	id = sdk.BigEndianToUint64(key[2+addrLen:])
	return
}

func ParseChunkBondRequestIndex(key []byte) (requesterAddr sdk.AccAddress, id uint64) {
	return parseIndexByAccAddressKey(key, KeyPrefixChunkBondRequestIndex)
}

func ParseChunkUnbondRequestIndex(key []byte) (requesterAddr sdk.AccAddress, id uint64) {
	return parseIndexByAccAddressKey(key, KeyPrefixChunkUnbondRequestIndex)
}

func ParseAliveChunkIndexByInsuranceProviderKey(key []byte) (insuranceProviderAddr sdk.AccAddress, id uint64) {
	return parseIndexByAccAddressKey(key, KeyPrefixAliveChunkIndexByInsuranceProvider)
}

func ParseAliveChunkIndexByValidatorKey(key []byte) (validatorAddr sdk.ValAddress, id uint64) {
	return parseIndexByValAddressKey(key, KeyPrefixAliveChunkIndexByValidator)
}

func ParseInsuranceBidIndexByInsuranceProviderKey(key []byte) (insuranceProviderAddr sdk.AccAddress, id uint64) {
	return parseIndexByAccAddressKey(key, KeyPrefixInsuranceBidIndexByInsuranceProvider)
}

func ParseInsuranceBidIndexByValidatorKey(key []byte) (validatorAddr sdk.ValAddress, id uint64) {
	return parseIndexByValAddressKey(key, KeyPrefixInsuranceBidIndexByValidator)
}

func ParseInsuranceUnbondRequestIndexByInsuranceProviderKey(key []byte) (insuranceProviderAddr sdk.AccAddress, id uint64) {
	return parseIndexByAccAddressKey(key, KeyPrefixInsuranceUnbondRequestIndexByInsuranceProvider)
}
