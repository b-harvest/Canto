<!-- order: 8 -->

# Parameters

The `liquidstaking` module contains the following parameters:

| Param          | Type             | Default                |
|----------------|------------------|------------------------|  
| DynamicFeeRate | string (sdk.Dec) | "0.000000000000000000" |

| Param      | Type             | Default                |
|------------|------------------|------------------------|  
| R0         | string (sdk.Dec) | "0.000000000000000000" |
| USoftCap   | string (sdk.Dec) | "0.050000000000000000" |
| UHardCap   | string (sdk.Dec) | "0.100000000000000000" |
| UOptimal   | string (sdk.Dec) | "0.090000000000000000" |
| Slope1     | string (sdk.Dec) | "0.100000000000000000" |
| Slope2     | string (sdk.Dec) | "0.400000000000000000" |
| MaxFeeRate | string (sdk.Dec) | "0.500000000000000000" |

## R0

Minimum fee rate.

## USoftCap

SoftCap for utilization ratio. If U is below softcap, fee rate is R0.

## UHardCap

HardCap for utilization ratio. U cannot bigger than hardcap.

## UOptimal

Optimal utilization ratio.

## Slope1

If the current utilization ratio is below optimal, the fee rate increases at a slow pace.

## Slope2

If the current utilization ratio is above optimal, the fee rate increases at a faster pace.

## MaxFeeRate

Maximum fee rate. Fee rate cannot exceed this value.
