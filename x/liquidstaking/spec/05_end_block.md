<!-- order: 5 -->

# End-Block

## Before-End-Block

These operations occur before the end-block operations for the liquidity module.

**Store/Remove requests from messages**

- After successful message verification and coin escrow, the incoming `MsgCancelInsuranceBid`, `MsgCancelInsuranceUnbond`, `MsgCancelLiquidStaking` and `MsgCancelLiquidUnstaking` messages are used to remove existing requests.
- And `MsgBidInsurance`, `MsgUnbondInsuranced`, `MsgLiquidStaking` and `MsgLiquidUnstaking` messages are converted to requests and stored.

## Bond Relocation

If there are `{*action}Request` and `InsuranceBid` that have not yet matched, the bond relocation is executed. It means relocating the bond between validator-chunk-insurance in the most efficient way when there is a change in list of `AliveChunk`, taking into account `{*action}Request`, `InsuranceBid` and new fee rankings.

**Bond relocation process**

1. Create new fee ranking by integrating multiple requests
    1. Alive chunks with insurance unbond requests are excluded from candidates
    2. A new candidate is created by combining the insurance fee rate of the insurance bid and the validator fee rate of the designated validator
    3. All candidates are ranked in descending order based on their total fee
2. Determine new alive chunks
    1. The maximum size of new alive chunks is affected by the chunk bond and unbond requests, and it must not exceed the `MaxiAliveChunk`
3. Bond/Unbond/Redelegation process
    1. If the number of alive chunks remains unchanged before and after the new ranking decision, bond relocation can only be done through redelegation
    2. Bonding or unbonding is necessary only when there is a difference in the number of alive chunks, and the remainder can be adjusted by redelegation
    3. Unbonding or redelegating chunks and insurance are managed by `UnbondingChunk` until the unbonding period elapses, and slashing risk is also insured
