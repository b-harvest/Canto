package types

import (
	"bytes"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
)

// constants
const (
	// ModuleName is the name of the module
	ModuleName = "liquidstaking"

	// StoreKey is the string store representation
	StoreKey = ModuleName

	// RouterKey is the msg router key for the liquidstaking module
	RouterKey = ModuleName
)

// prefix bytes for the liquidstaking persistent store
const (
	prefixLastChunkId = iota + 1
	prefixLastInsuranceId
	prefixChunk
	prefixInsurance
	prefixInsurancesByProviderIndex
	prefixDelegatorIndex
	prefixWithdrawingInsurance
	prefixPreviousInsuranceIndex
	prefixLiquidUnstakeUnbondingDelegation
	prefixEpoch
)

// KVStore key prefixes
var (
	KeyPrefixLastChunkId                      = []byte{prefixLastChunkId}
	KeyPrefixLastInsuranceId                  = []byte{prefixLastInsuranceId}
	KeyPrefixChunk                            = []byte{prefixChunk}
	KeyPrefixInsurance                        = []byte{prefixInsurance}
	KeyPrefixInsurancesByProviderIndex        = []byte{prefixInsurancesByProviderIndex}
	KeyPrefixDelegatorIndex                   = []byte{prefixDelegatorIndex}
	KeyPrefixWithdrawingInsurance             = []byte{prefixWithdrawingInsurance}
	KeyPrefixPreviousInsuranceIndex           = []byte{prefixPreviousInsuranceIndex}
	KeyPrefixLiquidUnstakeUnbondingDelegation = []byte{prefixLiquidUnstakeUnbondingDelegation}
	KeyPrefixEpoch                            = []byte{prefixEpoch}
)

func GetChunkKey(chunkId uint64) []byte {
	return append(KeyPrefixChunk, sdk.Uint64ToBigEndian(chunkId)...)
}

func GetInsuranceKey(insuranceId uint64) []byte {
	return append(KeyPrefixInsurance, sdk.Uint64ToBigEndian(insuranceId)...)
}

func GetInsurancesByProviderIndexKey(providerAddress sdk.AccAddress, insuranceId uint64) []byte {
	return append(append(KeyPrefixInsurancesByProviderIndex, address.MustLengthPrefix(providerAddress)...), sdk.Uint64ToBigEndian(insuranceId)...)
}

func ParseInsurancesByProviderIndexKey(key []byte) (providerAddress sdk.AccAddress, insuranceId uint64) {
	if !bytes.HasPrefix(key, KeyPrefixInsurancesByProviderIndex) {
		panic("invalid insurances by provider index key")
	}

	providerAddressLength := key[1]
	providerAddress = sdk.AccAddress(key[2 : 2+providerAddressLength])
	insuranceId = sdk.BigEndianToUint64(key[2+providerAddressLength:])
	return
}
