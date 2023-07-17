---
Title: Liquidstaking
Description: A high-level overview of how the command-line interfaces (CLI) works for the liquidstaking module.
---

# Liquidstaking Module

## Synopsis

This document provides a high-level overview of how the command line (CLI) interface works for the `liquidstaking` module. 
To set up a local testing environment, you should run [init_testnet.sh](https://github.com/b-harvest/Canto/blob/liquidstaking-module/init_testnet.sh) 

Note that [jq](https://stedolan.github.io/jq/) is recommended to be installed as it is used to process JSON throughout the document.

## Command Line Interfaces

- [Transaction](#Transaction)
  - [ProvideInsurance](#ProvideInsurance)
  - [CancelProvideInsurance](#CancelProvideInsurance)
  - [LiquidStake](#LiquidStake)
  - [LiquidUnstake](#LiquidUnstake)
  - [DepositInsurance](#DepositInsurance)
  - [WithdrawInsurance](#WithdrawInsurance)
  - [WithdrawInsuranceCommission](#WithdrawInsuranceCommission)
  - [ClaimDiscountedReward](#ClaimDiscountedReward)
- [Query](#Query)
  - [Params](#Params)
  - [Epoch](#Epoch)
  - [Chunks](#Chunks)
  - [Chunk](#Chunk)
  - [Insurances](#Insurances)
  - [Insurance](#Insurance)
  - [WithdrawInsuranceRequests](#WithdrawInsuranceRequests)
  - [WithdrawInsuranceRequest](#WithdrawInsuranceRequest)
  - [UnpairingForUnstakingChunkInfos](#UnpairingForUnstakingChunkInfos)
  - [UnpairingForUnstakingChunkInfo](#UnpairingForUnstakingChunkInfo)
  - [RedelegationInfos](#RedelegationInfos)
  - [RedelegationInfo](#RedelegationInfo)
  - [ChunkSize](#ChunkSize)
  - [MinimumCollateral](#MinimumCollateral)
  - [States](#States)

# Transaction

## ProvideInsurance

Provide insurance.

Usage

```bash
provide-insurance [validator-address] [amount] [fee-rate]
```

| **Argument**      | **Description**                                                                                                          |
|:------------------|:-------------------------------------------------------------------------------------------------------------------------|
| validator-address | the validator address that the insurance provider wants to cover                                                         |
| amount            | amount of collalteral; it must be acanto and amount must be bigger than 7% of ChunkSize(=250K) tokens(9% is recommended) |
| fee-rate          | how much commission will you receive for providing insurance? (fee-rate x chunk's delegation reward) will be commission. |

Example

```bash
# Provide insurance with 9% of ChunkSize collateral and 10% as fee-rate.
cantod tx liquidstaking provide-insurance <validator-address> 22500000000000000000000acanto 0.1 --from key1 --fees 200000acanto  \
--from key1 \
--keyring-backend test \
--fees 200000acanto \
--output json | jq

#
# Tips
# 
# Query validators first you want to cover and copy operator_address of the validator.
# And use that address at <validator-address>
cantod q staking validators
#
# Query chunks
# You can see newly created insurances (initial status of insurance is "Pairing")
cantod q liquidstaking insurances
```

## CancelProvideInsurance

Provide insurance.

Usage

```bash
cancel-provide-insurance [insurance-id]
```

| **Argument** | **Description**                                |
|:-------------|:-----------------------------------------------|
| insurance-id | the id of pairing insurance you want to cancel |

Example

```bash
cantod tx liquidstaking cancel-provide-insurance 3
--from key1 \
--keyring-backend test \
--fees 200000acanto \
--output json | jq

#
# Tips
#
# Query insurances
# If it is succeeded, then you cannot see the insurance with the id in result.
cantod q liquidstaking insurances
```

## LiquidStake

Liquid stake coin.

Usage

```bash
liquid-stake [amount]
```

| **Argument**  | **Description**                                                                                          |
| :------------ |:---------------------------------------------------------------------------------------------------------|
| amount        | amount of coin to stake; it must be acanto and amount must be multiple of ChunkSize(=250K) tokens |

Example

```bash
# Liquid stake 1 chunk (250K tokens)
cantod tx liquidstaking liquid-stake 250000000000000000000000acanto \
--from key1 \
--keyring-backend test \
--fees 3000000acanto \
--gas 3000000 \
--output json | jq

#
# Tips
#
# Query account balances
# If liquid stake succeeded, you can see the newly minted lsToken
cantod q bank balances <address> -o json | jq

# Query chunks
# And you can see newly created chunk with new id
cantod q liquidstaking chunks
```

## LiquidUnstake

Liquid stake coin.

Usage

```bash
liquid-unstake [amount]
```

| **Argument**  | **Description**                                                                                      |
| :------------ |:-----------------------------------------------------------------------------------------------------|
| amount        | amount of coin to un-stake; it must be acanto and amount must be multiple of ChunkSize(=250K) tokens |

Example

```bash
# Liquid unstake 1 chunk (250K tokens)
cantod tx liquidstaking liquid-unstake 250000000000000000000000acanto \
--from key1 \
--keyring-backend test \
--fees 3000000acanto
--gas 3000000 \
--output json | jq

#
# Tips
#
# Query account balances
# If liquid unstake request is accepted, you can see lsToken corresponding msg.Amount is escrowed(=decreased).
# When the actual unstaking process is finished, then you can see unstaked token in your account.
# Notice the newly minted lsToken
cantod q bank balances <address> -o json | jq

# Query your unstaking request
# If your unstake request is accepted, then you can query your unstaking request.
cantod q liquidstaking unpairing-for-unstaking-chunk-infos --queued="true"
```


# Query

## Params


Query the current liquidstaking parameters information.

Usage

```bash
params
```

Example

```bash
cantod query liquidstaking params -o json | jq
```

## LiquidValidators

Query all liquid validators.

Usage

```bash
liquid-validators
```

Example

```bash
cantod query liquidstaking liquid-validators -o json | jq
```
## States

Query net amount state.

Usage

```bash
states
```

Example

```bash
cantod query liquidstaking states -o json | jq
```

## VotingPower

Query the voter’s staking and liquid staking voting power.

Usage

```bash
voting-power [voter]
```

| **Argument** |  **Description**      |
| :----------- | :-------------------- |
| voter        | voter account address |

Example

```bash
cantod query liquidstaking voting-power cre1mzgucqnfr2l8cj5apvdpllhzt4zeuh2c5l33n3 -o json | jq
```
