---
Title: Liquidstaking
Description: A high-level overview of what gRPC-gateway REST routes are supported in the liquidstaking module.
---

# Liquidstaking Module

## Synopsis

This document provides a high-level overview of what gRPC-gateway REST routes are supported in the liquidstaking module.
To set up a local testing environment, you should run [init_testnet.sh](https://github.com/b-harvest/Canto/blob/liquidstaking-module/init_testnet.sh)

## gRPC-gateway REST Routes


++https://github.com/crescent-network/crescent/blob/main/proto/crescent/farming/v1beta1/query.proto
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

## Params

Query the current liquidstaking parameters information.

Example Request

```bash
http://localhost:1317/canto/liquidstaking/v1/params
```

Example Response

```json
{
  "params": {
    "dynamic_fee_rate": {
      "r0": "0.000000000000000000",
      "u_soft_cap": "0.050000000000000000",
      "u_hard_cap": "0.100000000000000000",
      "u_optimal": "0.090000000000000000",
      "slope1": "0.100000000000000000",
      "slope2": "0.400000000000000000",
      "max_fee_rate": "0.500000000000000000"
    }
  }
}
```

## Epoch

Query the epoch information.

Example Request

```bash
http://localhost:1317/canto/liquidstaking/v1/epoch
```

Example Response

```json
{
  "epoch": {
    "current_number": "648",
    "start_time": "2060-10-01T01:34:14.723955Z",
    "duration": "1814400s",
    "start_height": "3235"
  }
}
```

## Chunks

Query chunks.

Usage

```bash
http://localhost:1317/canto/liquidstaking/v1/chunks
```

Example Response

```json
{
  "chunks": [
    {
      "chunk": {
        "id": "1",
        "paired_insurance_id": "4",
        "unpairing_insurance_id": "0",
        "status": "CHUNK_STATUS_PAIRED"
      },
      "derived_address": "canto14zq9dj3mde6kwl7302zxcf2nv83m3k3qj9cq3k"
    },
    {
      "chunk": {
        "id": "2",
        "paired_insurance_id": "7",
        "unpairing_insurance_id": "0",
        "status": "CHUNK_STATUS_PAIRED"
      },
      "derived_address": "canto15r7jycu6dsljrrngnuez8ytpk8sey3awyleeht"
    }
  ],
  "pagination": {
    "next_key": null,
    "total": "2"
  }
}
```

## Chunk

Query a chunk by id.

Example Request

```bash
http://localhost:1317/canto/liquidstaking/v1/chunks/1
```

Example Response

```json
{
  "chunk": {
    "id": "1",
    "paired_insurance_id": "4",
    "unpairing_insurance_id": "0",
    "status": "CHUNK_STATUS_PAIRED"
  },
  "derived_address": "canto14zq9dj3mde6kwl7302zxcf2nv83m3k3qj9cq3k"
}
```

## Insurances

Query insurances.

Example Request

```bash
http://localhost:1317/canto/liquidstaking/v1/insurances
```

Example Response

```json
{
  "insurances": [
    {
      "insurance": {
        "id": "1",
        "validator_address": "cantovaloper1xjlslz2vl7v6gu807fmfw8ae7726q9pf84kzqs",
        "provider_address": "canto1xjlslz2vl7v6gu807fmfw8ae7726q9pf9t3x34",
        "fee_rate": "0.100000000000000000",
        "chunk_id": "0",
        "status": "INSURANCE_STATUS_UNPAIRED"
      },
      "derived_address": "canto1p6qg4xu665ld3l8nr72z0vpsujf0s9ekhfjhuv",
      "fee_pool_address": "canto1fy0mcah0tcedpyqyz423mefdxh7zqz4g2lu8jf"
    },
    {
      "insurance": {
        "id": "2",
        "validator_address": "cantovaloper1xjlslz2vl7v6gu807fmfw8ae7726q9pf84kzqs",
        "provider_address": "canto1xjlslz2vl7v6gu807fmfw8ae7726q9pf9t3x34",
        "fee_rate": "0.100000000000000000",
        "chunk_id": "0",
        "status": "INSURANCE_STATUS_UNPAIRED"
      },
      "derived_address": "canto1hk5wgk3js5uqymxppawk87tv0j0fnc3pefcex4",
      "fee_pool_address": "canto1a3f65vrngauvsj066067qsjh068hgxezpdr6rg"
    },
    {
      "insurance": {
        "id": "3",
        "validator_address": "cantovaloper1xjlslz2vl7v6gu807fmfw8ae7726q9pf84kzqs",
        "provider_address": "canto1xjlslz2vl7v6gu807fmfw8ae7726q9pf9t3x34",
        "fee_rate": "0.100000000000000000",
        "chunk_id": "0",
        "status": "INSURANCE_STATUS_UNPAIRED"
      },
      "derived_address": "canto1yqg5xesskfhmzwdn3gaas6faz6d0yjwd34rrct",
      "fee_pool_address": "canto10m6pl6am95swkaa480789y2wqlmpyqang9uxeh"
    },
    {
      "insurance": {
        "id": "4",
        "validator_address": "cantovaloper1xjlslz2vl7v6gu807fmfw8ae7726q9pf84kzqs",
        "provider_address": "canto1xjlslz2vl7v6gu807fmfw8ae7726q9pf9t3x34",
        "fee_rate": "0.100000000000000000",
        "chunk_id": "1",
        "status": "INSURANCE_STATUS_PAIRED"
      },
      "derived_address": "canto1my633g6sqx9fr4szzxuj70zutmsd78zymhv5kf",
      "fee_pool_address": "canto1sdl4z9y8x59979qjx8ut9zyndsux9sld0s6kcv"
    },
    {
      "insurance": {
        "id": "5",
        "validator_address": "cantovaloper1xjlslz2vl7v6gu807fmfw8ae7726q9pf84kzqs",
        "provider_address": "canto1xjlslz2vl7v6gu807fmfw8ae7726q9pf9t3x34",
        "fee_rate": "0.100000000000000000",
        "chunk_id": "0",
        "status": "INSURANCE_STATUS_PAIRING"
      },
      "derived_address": "canto10fhsthfzmwcfjhqulyy8w0r2fqa9q5s6a88sz4",
      "fee_pool_address": "canto1mfkktq2nj4dcc9ypkfengcpjufntk9eak0atc3"
    },
    {
      "insurance": {
        "id": "6",
        "validator_address": "cantovaloper1xjlslz2vl7v6gu807fmfw8ae7726q9pf84kzqs",
        "provider_address": "canto1xjlslz2vl7v6gu807fmfw8ae7726q9pf9t3x34",
        "fee_rate": "0.100000000000000000",
        "chunk_id": "0",
        "status": "INSURANCE_STATUS_PAIRING"
      },
      "derived_address": "canto1nwhkqs073xj7sx3tq7yp7sxkqa6n4fawcy5fxf",
      "fee_pool_address": "canto1tjqa0tc9dxv7f23szzj3m9rn88uzegf5d0ks9n"
    },
    {
      "insurance": {
        "id": "7",
        "validator_address": "cantovaloper1xjlslz2vl7v6gu807fmfw8ae7726q9pf84kzqs",
        "provider_address": "canto1xjlslz2vl7v6gu807fmfw8ae7726q9pf9t3x34",
        "fee_rate": "0.010000000000000000",
        "chunk_id": "2",
        "status": "INSURANCE_STATUS_PAIRED"
      },
      "derived_address": "canto1dcd7gu8s5xez9hadt8sum0lnvfz93ntfgtq04q",
      "fee_pool_address": "canto1v9cv3cxst0sxrz92mptx6fj8utfsg8c0qr6v74"
    }
  ],
  "pagination": {
    "next_key": null,
    "total": "7"
  }
}
```


## Insurance

Query an Insurance by id.

Example Request

```bash
http://localhost:1317/canto/liquidstaking/v1/insurances/4
```

Example Response

```json
{
  "insurance": {
    "id": "4",
    "validator_address": "cantovaloper1xjlslz2vl7v6gu807fmfw8ae7726q9pf84kzqs",
    "provider_address": "canto1xjlslz2vl7v6gu807fmfw8ae7726q9pf9t3x34",
    "fee_rate": "0.100000000000000000",
    "chunk_id": "1",
    "status": "INSURANCE_STATUS_PAIRED"
  },
  "derived_address": "canto1my633g6sqx9fr4szzxuj70zutmsd78zymhv5kf",
  "fee_pool_address": "canto1sdl4z9y8x59979qjx8ut9zyndsux9sld0s6kcv"
}
```

## WithdrawInsuranceRequests

Query WithdrawInsuranceRequests.

Example Request

```bash
http://localhost:1317/canto/liquidstaking/v1/withdraw_insurance_requests
```

Example Response

```json
{
  "withdraw_insurance_requests": [
    {
      "insurance_id": "7"
    }
  ],
  "pagination": {
    "next_key": null,
    "total": "1"
  }
}
```

## WithdrawInsuranceRequest

Query a WithdrawInsuranceRequest by insurance id.

Example Request

```bash
http://localhost:1317/canto/liquidstaking/v1/insurances/5/withdraw_insurance_requests
```

Example Response

```json
{
  "withdraw_insurance_request": {
    "insurance_id": "5"
  }
}
```

## UnpairingForUnstakingChunkInfos

Query UnpairingForUnstakingChunkInfos.

Example Request

```bash
http://localhost:1317/canto/liquidstaking/v1/unpairing_for_unstaking_chunk_infos 
```

Example Response

```json
{
  "unpairing_for_unstaking_chunk_info": {
    "chunk_id": "2",
    "delegator_address": "canto1xjlslz2vl7v6gu807fmfw8ae7726q9pf9t3x34",
    "escrowed_lstokens": {
      "denom": "lscanto",
      "amount": "240214408039107442750000"
    }
  }
}
```

## UnpairingForUnstakingChunkInfo

Query an UnpairingForUnstakingChunkInfo by chunk id.

Example Request

```bash
http://localhost:1317/canto/liquidstaking/v1/chunks/2/unpairing_for_unstaking_chunk_infos
```

Example Response

```json
{
  "unpairing_for_unstaking_chunk_info": {
    "chunk_id": "2",
    "delegator_address": "canto1xjlslz2vl7v6gu807fmfw8ae7726q9pf9t3x34",
    "escrowed_lstokens": {
      "denom": "lscanto",
      "amount": "240214408039107442750000"
    }
  }
}
```

## RedelegationInfos

Query RedelegationInfos.

Example Request

```bash
http://localhost:1317/canto/liquidstaking/v1/redelegation_infos
```

Example Response

```bash
# Query redelegation-infos by chunk id
cantod query liquidstaking redelegation-infos -o json | jq
```

## RedelegationInfo

Query RedelegationInfo.

Usage

```bash
redelegation-info [chunk-id]
```

Example

```bash
# Query redelegation-info by chunk id
cantod query liquidstaking redelegation-info 1 -o json | jq
```

## ChunkSize

Query ChunkSize.

Usage

```bash
chunk-size
```

Example

```bash
# Query chunk size
cantod query liquidstaking chunk-size -o json | jq
```

## MinimumCollateral

Query MinimumCollateral.

Usage

```bash
minimum-collateral
```

Example

```bash
# Query minimum collateral  
cantod query liquidstaking minimum-collateral -o json | jq
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
