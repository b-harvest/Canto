package types

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
	prefixChunkByInsurance
	prefixInsurance
	prefixInsurancesByProvider
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
	KeyPrefixChunkByInsurance                 = []byte{prefixChunkByInsurance}
	KeyPrefixInsurance                        = []byte{prefixInsurance}
	KeyPrefixInsurancesByProvider             = []byte{prefixInsurancesByProvider}
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
