<!-- order: 2 -->

# State

## AliveChunk

`AliveChunk` stores information on staked chunks through the liquid staking module.

```go
type AliveChunk struct {
  Id                       uint64     // id of alive chunk
  ValidatorAddress         string     // address of the validator that the chunk is staked to
  InsuranceProviderAddress string     // address of the insurance provider for the chunk
  TokenAmount              sdk.Int    // amount of the native token in the chunk
  InsuranceAmount          sdk.Int    // amount of the native tokne for insurance
  InsuranceFeeRate         sdk.Dec    // insurance fee rate of the chunk
}
```

## UnbondingChunk

```go
type UnbondingChunk struct {
  ValidatorAddress         string     // address of the validator that the chunk is staked to
  InsuranceProviderAddress string     // address of the insurance provider for the chunk
  TokenAmount              sdk.Int    // amount of the native token in the chunk
  InsuranceAmount          sdk.Int    // amount of the native tokne for insurance
  InsuranceFeeRate         sdk.Dec    // insurance fee rate of the chunk
}
```

## LiquidStakingInfo (for querying?)

`LiquidStakingInfo` stores information necessary for minting and burning. As a ledger, it collects and manages elements for determining the fair value of lsToken.

```go
type LiquidStakingInfo struct {
  MintRate                 sdk.Dec    // LSTokenTotalSupply / NetAmount
  LStokenTotalSupply       sdk.Int    // total supply of LStoken
  NetAmount                sdk.Int    // sum of all proxy account's native token balance
  TokenAmount              sdk.Int    // sum of all token amount worth of delegation shares + native token balance of each proxy accounts
  TotalRemainingRewards    sdk.Dec    // sum of remaining rewards of all proxy accounts
  TotalUnbondingBalance    sdk.Int    // sum of unbonding balance of all proxy accounts
}
```

## InsuranceBid

`InsuranceBid` contains information on a user's(potential insurance provider) specified validator, insurance amount, and fee rate for providing insurance

```go
type InsuranceBid struct {
  Id                       uint64
  ValidatorAddress         string
  InsuranceProviderAddress string
  InsuranceAmount          sdk.Int
  InsuranceFeeRate         sdk.Dec
}
```

## InsuranceUnbondRequest

`InsuranceUnbondRequest` contains id for insurance unbond request only. By specifying the alive chunk id, all necessary information in the `AliveChunk` can be accessed.

```go
type InsuranceUnbondRequest struct {
  AliveChunkId             uint64
}
```

## ChunkBondRequest

`ChunkBondRequest` stores information about chunk bond requests.

```go
type ChunkBondRequest struct {
  Id                       uint64
  RequestorAddress         string
  TokenAmount              sdk.Int
}
```

## ChunkUnbondRequest

`ChunkUnbondRequest` stores information about chunk unbond requests.

```go
type ChunkUnbondRequest struct {
  Id                       uint64
  RequestorAddress         string
  NumChunkUnbond           uint64
}
```
