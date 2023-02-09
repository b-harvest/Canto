<!-- order: 1 -->

# Concepts

Protocols using the PoS(Proof-of-Stake) consensus mechanism usually require token owners to stake their tokens on the network in order to participate in the governance process. At this time, the user's token is locked and loses its potential utility. Liquid staking is a staking method that can avoid this loss of capital efficiency. In other words, liquid staking allows holders to earn staking rewards while still being able to trade or use their assets as needed, whereas traditional staking typically requires locking up assets for a fixed period of time.
Basically, in liquid staking, a new token(lsToken) is minted as proof of staking the native token, and the lsToken is circulated in the market instead of native token. In order for lsToken to be fully fungible, the relationship between the staking status of the native token used in liquid staking and the minted lsToken must be fungible. In other words, regardless of which validator the user chooses for liquid staking, the reward accumulated in lsToken and the risk of lsToken must be the same.
However, as is well known, each validator inevitably has differences in node operating ability, security level, and required fee rate, so the reward or risk of staking varies depending on which validator the user chooses. To solve this problem, we would like to propose our own unique liquid staking that has features such as insurance, fee-rate competition and reward distribution.

## Insurance

Insurance protects against the potential loss of funds from the slashing of staked tokens. In simpler terms, the risk of loss from slashing is transferred to the insurance provider, ensuring that the initially staked tokens through the liquid staking module are always protected. As previously mentioned, each validator carries a varying level of risk for slashing. By transferring this risk to the insurance provider, the user’s choice of validator becomes no longer important. This means that the minted lsToken is independent of the user's choice and is completely fungible in terms of risk for slashing.

## Fee-rate competition

Insurance providers charge a fee for the protection of staked tokens. Tokens can only be staked through the liquid staking module if the corresponding insurance is in place. This requirement for insurance incentivizes insurance providers to charge high fees for their services. However, an increase in insurance costs decreases the return of liquid staking, which in turn reduces the motivation for users to use it. Therefore, it is necessary to prevent insurance providers from raising fees arbitrarily, we achieved this through fee rate competition.

Fee rate competition allows liquid staking only for slots whose fee rate fall within a certain rank as determined by the governance. The fee rate here is the total of the insurance fee rate required by the insurance provider and the commission fee rate set by the validator selected by the insurance provider. In this case, both the validator and the insurance provider are discouraged from raising the fee rate excessively, because if the fee rate is set too high and the amount to be staked is not allocated, no profit will be made.

## Reward distribution

Service providers(validators and insurance providers) who win the fee rate competition receive a commission for their services, in the form of a fee corresponding to the fee rate they have set. Therefore, since the fee rates set by service providers vary, the staking rewards earned for each chunk(unit amount of tokens in liquid staking) also differ. Different rewards for each chunk means that the chunk has different status depending on the validator it is delegated to. Since this causes lsToken to become non-fungible, it is important to ensure that rewards are distributed equally among different chunks. Therefore, the liquid staking module preserves the fungibility of lsToken by collecting all the rewards through the account managing each chunk and then distributing them fairly.
