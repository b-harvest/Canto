<!-- order: 3 -->

# State Transitions

## AliveChunks / UnbondingChunks

State transitions of alive chunks are performed on every `EndBlock` of last block in epoch to keep in track of any changes in alive chunk set

**Bond relocation**

- When chunk bond/unbond request is accepted
- When insurance unbond request is accepted
- When fee ranking is changed

**Slashing**

## InsuranceBids

- When exisitng `InsuranceBid` is bonded with chunk
  - after bond relocation process, an insurance bid bonded to a new alive chunk is removed from the state
- When new `InsuranceBid` is added
  - when a user submits a new insurance bid, it is added to the state
- When exisiting `InsuranceBid` is canceled
  - when a user cancel the exisiting bid through `MsgCancelInsuranceBid`, it is removed from the state in that block
- When chunk-insurance bond is broken unintentionally
  - chunk unbond request, pushed out of the valid fee rank
  - if the chunk that was bonded to the insurance is unbonded or redelegated, the insurance is added to the `InsuranceBids` after the unbonding period has passed

## InsuranceUnbondRequest

- When new insurance unbond request is added
  - when a user submits a new insurance unbond request, it is added to the state
- When existing insurance unbond request is accepted or canceled
  - when a user cancel the exisiting request through `MsgCancelInsuranceUnbond`, it is removed from the state in that block
  - after bond relocation process, accepted request is removed from the state

## ChunkBondRequest

- When new chunk bond request(liquid staking request) is added
  - when a user submits a new chunk bond request, it is added to the state
- When existing chunk bond request is accepted or canceled
  - when a user cancel the exisiting request through `MsgCancelLiquidStaking`, it is removed from the state in that block
  - after bond relocation process, accepted request is removed from the state
- When chunk-insurance bond is broken unintentionally
  - insurance unbond request, pushed out of the valid fee rank
  - if the chunk that was bonded to the insurance is unbonded, the chunk is added to the `ChunkBondRequest` after the unbonding period has passed

## ChunkUnbondRequest

- When new chunk unbond request(liquid unstaking request) is added
  - when a user submits a new chunk unbond request, it is added to the state
- When existing chunk unbond request is accepted or canceled
  - when a user cancel the exisiting request through `MsgCancelLiquidUnstaking`, it is removed from the state in that block
  - after bond relocation process, accepted request is removed from the state

## Reward distribution

- In the last block of the epoch, the staking reward accumulated in the epoch is claimed and distributed so that all alive chunks have the same amount of tokens
