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

**☆ Actual Bond/Unbond/Redelegation determination**

1. Create two lists of chunk-validator pairs before and after bond relocation. In other words, the two lists mean the list of existing chunk-validator bond and the proposed chunk-validator bond considering all requests submitted up to the current block
2. By comparing the two lists, we can know exactly how many chunks each validator needs to receive or give (if both inputs and outputs of chunk need to occur in one validator, the amount will be offset)
3. Iterate the matching until all validators’ chunks equal to their target (proposed state)

    matching logic : the validator with the most to give matches the validator with the most to receive

    moving amount of chunk : min(surplus of validator with the most to give, deficit of validator with the most to receive)

4. If the last remaining validator is the receiver after the above process, the quantity in the bond chunk request is newly staked to the validator. If the last remaining validator is the giver, the amount is unstaked
