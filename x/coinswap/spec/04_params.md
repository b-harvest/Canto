<!--
order: 4
-->

# Parameters

The coinswap module contains the following parameters:

| Key                    | Type         | Example                                |
|:-----------------------|:-------------|:---------------------------------------|
| Fee                    | string (dec) | "0.003"                                |
| PoolCreationFee        | sdk.Coin     | "5000acanto"                           |
| TaxRate                | string (dec) | "0.4"                                  |
| MaxStandardCoinPerPool | string (int) | "100000000"                            |
| MaxSwapAmount          | sdk.Coins    | [{"denom":"stake","amount":"1000000"}] |

### Fee
Swap fee rate for swap. In this version, swap fees aren't paid upon swap orders directly. Instead, pool just adjust pool's quoting prices to reflect the swap fees.

### PoolCreationFee
Fee paid for to create a pool. This fee prevents spamming and is collected in the fee collector.

### TaxRate
Community tax rate for pool creation fee. This tax is collected in the fee collector.

### MaxStandardCoinPerPool
Maximum amount of standard coin per pool. This parameter is used to prevent pool from being too large.

### MaxSwapAmount
Maximum amount of swap amount. This parameter is used to prevent swap from being too large. It is also used as whitelist for pool creation.