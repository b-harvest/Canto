package keeper

import (
	"time"

	"github.com/Canto-Network/Canto-Testnet-v2/v1/x/liquidstaking/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// TODO: rename func
func (k *Keeper) delegateTokenAmount(ctx *sdk.Context, valAddress string, tokenAmount sdk.Int) error {
	val, err := sdk.ValAddressFromBech32(valAddress)
	if err != nil {
		return err
	}
	validator, found := k.stk.GetValidator(*ctx, val)
	if !found {
		panic("validator must exist")
	}

	if _, err := k.stk.Delegate(*ctx,
		types.LiquidStakingModuleAccount,
		tokenAmount,
		stakingtypes.Unbonded,
		validator,
		true); err != nil {
		return err
	}
	return nil
}

func (k Keeper) DelegateTokenAmount(ctx sdk.Context, valAddress string, tokenAmount sdk.Int) error {
	bondDenom := k.stk.BondDenom(ctx)
	if err := k.bk.SendCoinsFromModuleToAccount(
		ctx,
		types.ModuleName,
		types.LiquidStakingModuleAccount,
		sdk.NewCoins(sdk.NewCoin(bondDenom, tokenAmount))); err != nil {
		return err
	}
	return k.delegateTokenAmount(&ctx, valAddress, tokenAmount)
}

func (k Keeper) RedelegateTokenAmount(ctx sdk.Context, srcValAddress, dstValAddress string, tokenAmount sdk.Int) error {
	srcVal, err := sdk.ValAddressFromBech32(srcValAddress)
	if err != nil {
		return err
	}
	dstVal, err := sdk.ValAddressFromBech32(dstValAddress)
	if err != nil {
		return err
	}

	srcValStr := srcVal.String()
	dstValStr := dstVal.String()
	_ = srcValStr
	_ = dstValStr

	// check the source validator already has receiving transitive redelegation
	hasReceiving := k.stk.HasReceivingRedelegation(ctx, types.LiquidStakingModuleAccount, srcVal)
	if hasReceiving {
		return stakingtypes.ErrTransitiveRedelegation
	}

	// calculate delShares from tokens with validation
	shares, err := k.stk.ValidateUnbondAmount(ctx, types.LiquidStakingModuleAccount, srcVal, tokenAmount)
	if err != nil {
		return err
	}

	// TODO: check lastRedelegation

	cachedCtx, writeCache := ctx.CacheContext()
	_, err = k.stk.BeginRedelegation(cachedCtx, types.LiquidStakingModuleAccount, srcVal, dstVal, shares)
	if err != nil {
		return err
	}
	writeCache()
	return nil
}

func (k Keeper) UndelegateTokenAmount(
	ctx sdk.Context, del, val string, tokenAmount sdk.Int,
) (time.Time, error) {
	delAddr, err := sdk.AccAddressFromBech32(del)
	if err != nil {
		return time.Time{}, err
	}
	valAddr, err := sdk.ValAddressFromBech32(val)
	if err != nil {
		return time.Time{}, err
	}
	validator, found := k.stk.GetValidator(ctx, valAddr)
	if !found {
		return time.Time{}, stakingtypes.ErrNoDelegatorForAddress
	}

	// calculate delShares from tokens with validation
	share, err := k.stk.ValidateUnbondAmount(ctx,
		types.LiquidStakingModuleAccount,
		valAddr,
		tokenAmount)
	if err != nil {
		return time.Time{}, err
	}
	if !share.IsPositive() {
		return time.Time{}, types.ErrInvalidTokenAmount // TODO: define error
	}

	if k.stk.HasMaxUnbondingDelegationEntries(ctx, delAddr, valAddr) {
		return time.Time{}, stakingtypes.ErrMaxUnbondingDelegationEntries
	}

	returnAmount, err := k.stk.Unbond(ctx, types.LiquidStakingModuleAccount, valAddr, share)
	if err != nil {
		return time.Time{}, err
	}

	// transfer the validator tokens to the not bonded pool
	if validator.IsBonded() {
		coins := sdk.NewCoins(sdk.NewCoin(k.stk.BondDenom(ctx), returnAmount))
		if err = k.bk.SendCoinsFromModuleToModule(ctx, stakingtypes.BondedPoolName, stakingtypes.NotBondedPoolName, coins); err != nil {
			// panic(err)
			return time.Time{}, err
		}
	}

	completionTime := ctx.BlockHeader().Time.Add(k.stk.UnbondingTime(ctx))
	ubd := k.stk.SetUnbondingDelegationEntry(ctx, delAddr, valAddr, ctx.BlockHeight(), completionTime, returnAmount)
	k.stk.InsertUBDQueue(ctx, ubd, completionTime)

	return completionTime, nil
}

func (k *Keeper) getModuleAccountBalance(ctx *sdk.Context) sdk.Coin {
	bondDenom := k.stk.BondDenom(*ctx)
	return sdk.NewCoin(bondDenom, k.bk.SpendableCoins(*ctx, types.LiquidStakingModuleAccount).AmountOf(bondDenom))
}

func (k *Keeper) getDelegationReward(ctx *sdk.Context) sdk.Dec {
	bondDenom := k.stk.BondDenom(*ctx)
	totalReward := sdk.ZeroDec()

	// Cache ctx for calculate rewards
	cachedCtx, _ := ctx.CacheContext()
	k.stk.IterateDelegations(
		cachedCtx, types.LiquidStakingModuleAccount,
		func(_ int64, delegation stakingtypes.DelegationI) (stop bool) {
			valAddr := delegation.GetValidatorAddr()
			val := k.stk.Validator(cachedCtx, valAddr)
			endingPeriod := k.dk.IncrementValidatorPeriod(cachedCtx, val)
			reward := k.dk.CalculateDelegationRewards(cachedCtx, val, delegation, endingPeriod).AmountOf(bondDenom)

			if reward.IsPositive() {
				totalReward = totalReward.Add(reward)
			}
			return false
		},
	)
	return totalReward
}

func (k *Keeper) withdrawDelegationReward(ctx *sdk.Context) sdk.Int {
	totalRewards := sdk.ZeroInt()
	bondDenom := k.stk.BondDenom(*ctx)
	k.stk.IterateDelegations(
		*ctx, types.LiquidStakingModuleAccount,
		func(_ int64, delegation stakingtypes.DelegationI) (stop bool) {
			valAddr := delegation.GetValidatorAddr()
			reward, err := k.dk.WithdrawDelegationRewards(*ctx, types.LiquidStakingModuleAccount, valAddr)
			if err != nil {
				panic(err)
			}
			totalRewards = totalRewards.Add(reward.AmountOf(bondDenom))
			return false
		},
	)
	return totalRewards
}

func (k Keeper) DistributeReward(ctx sdk.Context) error {
	// TODO: check threshold
	// delegationReward := k.getDelegationReward(&ctx)

	reward := k.withdrawDelegationReward(&ctx)
	balance := k.getModuleAccountBalance(&ctx)
	if !reward.Equal(balance.Amount) {
		// NOTE: just for debug why?
		reward = balance.Amount
	}
	aliveChunks := k.GetAllAliveChunks(ctx)
	numAliveChunks := len(aliveChunks)
	rewardPerAliveChunk := reward.QuoRaw(int64(numAliveChunks))
	if rewardPerAliveChunk.IsZero() {
		return nil
	}
	for _, aliveChunk := range aliveChunks {
		aliveChunk.TokenAmount = aliveChunk.TokenAmount.Add(rewardPerAliveChunk)
		if err := k.DelegateTokenAmount(ctx, aliveChunk.ValidatorAddress, rewardPerAliveChunk); err != nil {
			return err
		}
		k.SetAliveChunk(ctx, aliveChunk)
	}

	return nil
}
