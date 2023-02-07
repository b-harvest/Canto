<!-- order: 4 -->

# Messages

## MsgBidInsurance

```go
type MsgBidInsurance struct {
  ValidatorAddress      string
  InsuranceAmount       sdk.Int
  InsuranceFeeRate      sdk.Dec
}
```

**Validity checks**

The transaction that is triggered with `MsgBidInsurance`fails if:

- The designated validator does not exist
- The balance of msg sender does not have enough amount of coins for `InsuranceAmount`
- `InsuranceAmount` is less than `MinInsurancePercentage` * `ChunkSize` / `MintRate`
- condition for `InsuranceFeeRate`?

## MsgCancelInsuranceBid

```go
type MsgCancelInsuranceBid struct {
  BidId   uint64
}
```

**Validity checks**

The transaction that is triggered with `MsgCancelInsuranceBid`fails if:

- The `BidId` does not exist
- The address of msg sender does not match with `InsuranceBid.InsuranceProviderAddress`

## MsgUnbondInsurance

```go
type MsgUnbondInsurance struct {
  AliveChunkId   uint64
}
```

**Validity checks**

The transaction that is triggered with `MsgUnbondInsurance`fails if:

- The `AliveChunkId` does not exist
- The address of msg sender does not match with `AliveChunk.InsuranceProviderAddress`

## MsgCancelInsuranceUnbond

```go
type MsgCancelInsuranceUnbond struct {
  AliveChunkId   uint64
}
```

**Validity checks**

The transaction that is triggered with `MsgCancelInsuranceUnbond`fails if:

- The `AliveChunkId` does not exist
- The address of msg sender does not match with `AliveChunk.InsuranceProviderAddress`

## MsgLiquidStaking

```go
type MsgLiquidStaking struct {
  TokenAmount      sdk.Int
}
```

**Validity checks**

The transaction that is triggered with `MsgLiquidStaking`fails if:

- The balance of msg sender does not have enough amount of coins for `TokenAmount`
- The `TokenAmount` is less than `ChunkSize` / `MintRate`

## MsgCancelLiquidStaking

```go
type MsgCancelLiquidStaking struct {
  ChunkBondRequestId  uint64
}
```

**Validity checks**

The transaction that is triggered with `MsgCancelLiquidStaking`fails if:

- The `ChunkBondRequestId` does not exist
- The address of msg sender does not match with `ChunkBondRequest.RequestorAddress`

## MsgLiquidUnstaking

```go
type MsgLiquidUnstaking struct {
  NumChunkUnstake  uint64
}
```

**Validity checks**

The transaction that is triggered with `MsgLiquidUnstaking`fails if:

- The balance of msg sender does not have enough amount of coins for `NumChunkUnstake` * `ChunkSize`

## MsgCancelLiquidUnstaking

```go
type MsgCancelLiquidUnstaking struct {
  ChunkUnbondRequestId  uint64
}
```

**Validity checks**

The transaction that is triggered with `MsgCancelLiquidUntaking`fails if:

- The `ChunkUnbondRequestId` does not exist
- The address of msg sender does not match with `ChunkUnbondRequest.RequestorAddress`
