package keeper

import (
	"github.com/Canto-Network/Canto/v6/x/liquidstaking/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

func (k Keeper) coverUnpairingForUnstakeCase(ctx sdk.Context) error {
	var unbondingDelegationInfos []types.LiquidUnstakeUnbondingDelegationInfo
	err := k.IterateAllLiquidUnstakeUnbondingDelegationInfos(ctx, func(liquidUnstakeUnbondingDelegationInfo types.LiquidUnstakeUnbondingDelegationInfo) (bool, error) {
		unbondingDelegationInfos = append(unbondingDelegationInfos, liquidUnstakeUnbondingDelegationInfo)
		return false, nil
	})
	if err != nil {
		panic(err)
	}

	// For all completed unboding delegation infos
	for _, unbondingDelegationInfo := range unbondingDelegationInfos {
		chunk, found := k.GetChunk(ctx, unbondingDelegationInfo.ChunkId)
		if !found {
			panic(types.ErrNotFoundChunk.Error())
		}
		if chunk.Status != types.CHUNK_STATUS_UNPAIRING_FOR_UNSTAKE {
			panic(types.ErrInvalidChunkStatus.Error())
		}

		// get insurance from chunk.InsuranceId
		insurance, found := k.GetInsurance(ctx, chunk.InsuranceId)
		if !found {
			panic(types.ErrNotFoundInsurance.Error())
		}

		// TODO: retreive using unstaker address, not chunk derived address
		// get unbonding delegation using staking keeper
		unbondingDelegation, found := k.stakingKeeper.GetUnbondingDelegation(
			ctx, chunk.DerivedAddress(),
			insurance.GetValidator(),
		)
		if !found {
			panic(types.ErrNotFoundUnbondingDelegation.Error())
		}

		// check if chunk got damaged during unbonding
		// for all entries of unbondingDelegation
		for _, entry := range unbondingDelegation.Entries {
			if entry.CompletionTime.Equal(unbondingDelegationInfo.CompletionTime) &&
				entry.InitialBalance.Equal(types.ChunkSize) {
				diff := entry.InitialBalance.Sub(entry.Balance)
				if diff.IsPositive() {
					// chunk got damaged, insurance should cover it
				}
			}

		}
		// if entry is not completed
		// if entry.IsMature(ctx.BlockTime()) {
		// then chunk got damaged
		// }

		// Check if chunk got damaged or not

		k.DeleteChunk(ctx, chunk.Id)
		return nil
	}
	return nil
}

func (k Keeper) DoLiquidStake(ctx sdk.Context, msg *types.MsgLiquidStake) (chunks []types.Chunk, newShares sdk.Dec, lsTokenMintAmount sdk.Int, err error) {
	delAddr := msg.GetDelegator()
	amount := msg.Amount

	// Check if max paired chunk size is exceeded or not
	currenPairedChunks := 0
	err = k.IterateAllChunks(ctx, func(chunk types.Chunk) (bool, error) {
		if chunk.Status == types.CHUNK_STATUS_PAIRED {
			currenPairedChunks++
		}
		return false, nil
	})
	if err != nil {
		return
	}
	availableChunks := types.MaxPairedChunks - currenPairedChunks
	if availableChunks <= 0 {
		err = sdkerrors.Wrapf(types.ErrMaxPairedChunkSizeExceeded, "current paired chunk size: %d", currenPairedChunks)
		return
	}

	if err = k.ShouldBeBondDenom(ctx, amount.Denom); err != nil {
		return
	}
	// Liquid stakers can send amount of tokens corresponding multiple of chunk size and create multiple chunks
	if err = k.ShouldBeMultipleOfChunkSize(amount.Amount); err != nil {
		return
	}
	chunksToCreate := amount.Amount.Quo(types.ChunkSize).Int64()
	if chunksToCreate > int64(availableChunks) {
		err = sdkerrors.Wrapf(
			types.ErrExceedAvailableChunks,
			"requested chunks to create: %d, available chunks: %d",
			chunksToCreate,
			availableChunks,
		)
		return
	}

	// Check if there are any pairing insurances
	var pairingInsurances []types.Insurance
	validatorMap := make(map[string]stakingtypes.Validator)
	err = k.IteratePairingInsurances(ctx, func(insurance types.Insurance) (bool, error) {
		if _, ok := validatorMap[insurance.ValidatorAddress]; !ok {
			validator, found := k.stakingKeeper.GetValidator(ctx, insurance.GetValidator())
			valid, err := k.isValidValidator(ctx, validator, found)
			if err != nil {
				return false, nil
			}
			if valid {
				validatorMap[insurance.ValidatorAddress] = validator
			} else {
				return false, nil
			}
		}
		pairingInsurances = append(pairingInsurances, insurance)
		return false, nil
	})
	if err != nil {
		return
	}
	if chunksToCreate > int64(len(pairingInsurances)) {
		err = types.ErrNoPairingInsurance
		return
	}

	// TODO: Must be proved that this is safe, we must call this before sending
	nas := k.GetNetAmountState(ctx)
	types.SortInsurances(validatorMap, pairingInsurances, false)
	totalNewShares := sdk.ZeroDec()
	totalLsTokenMintAmount := sdk.ZeroInt()
	for i := int64(0); i < chunksToCreate; i++ {
		cheapestInsurance := pairingInsurances[0]
		pairingInsurances = pairingInsurances[1:]

		// Now we have the cheapest pairing insurance and valid msg liquid stake! Let's create a chunk
		// Create a chunk
		chunkId := k.getNextChunkIdWithUpdate(ctx)
		chunk := types.NewChunk(chunkId)

		// Escrow liquid staker's coins
		if err = k.bankKeeper.SendCoins(
			ctx,
			delAddr,
			chunk.DerivedAddress(),
			sdk.NewCoins(amount),
		); err != nil {
			return
		}
		chunk.InsuranceId = cheapestInsurance.Id
		validator := validatorMap[cheapestInsurance.ValidatorAddress]

		// Delegate to the validator
		// Delegator: DerivedAddress(chunk.Id)
		// Validator: insurance.ValidatorAddress
		// Amount: msg.Amount
		newShares, err = k.stakingKeeper.Delegate(ctx, chunk.DerivedAddress(), amount.Amount, stakingtypes.Unbonded, validator, true)
		if err != nil {
			return
		}
		totalNewShares = totalNewShares.Add(newShares)

		liquidBondDenom := k.GetLiquidBondDenom(ctx)
		// Mint the liquid staking token
		lsTokenMintAmount = amount.Amount
		if nas.LsTokensTotalSupply.IsPositive() {
			lsTokenMintAmount = types.NativeTokenToLiquidStakeToken(amount.Amount, nas.LsTokensTotalSupply, nas.NetAmount)
		}
		if !lsTokenMintAmount.IsPositive() {
			err = sdkerrors.Wrapf(types.ErrInvalidAmount, "amount must be greater than or equal to %s", amount.String())
			return
		}
		mintedCoin := sdk.NewCoin(liquidBondDenom, lsTokenMintAmount)
		if err = k.bankKeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(mintedCoin)); err != nil {
			return
		}
		totalLsTokenMintAmount = totalLsTokenMintAmount.Add(lsTokenMintAmount)
		if err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, delAddr, sdk.NewCoins(mintedCoin)); err != nil {
			return
		}
		chunk.SetStatus(types.CHUNK_STATUS_PAIRED)
		cheapestInsurance.SetStatus(types.INSURANCE_STATUS_PAIRED)
		k.SetChunk(ctx, chunk)
		k.SetInsurance(ctx, cheapestInsurance)
		k.DeletePairingInsuranceIndex(ctx, cheapestInsurance)
		chunks = append(chunks, chunk)
	}
	return
}

func (k Keeper) DoLiquidUnstake(ctx sdk.Context, msg *types.MsgLiquidUnstake) (
	unstakedChunks []types.Chunk,
	unstakeUnbondingDelegationInfos []types.LiquidUnstakeUnbondingDelegationInfo,
	err error,
) {
	delAddr := msg.GetDelegator()
	amount := msg.Amount

	if err = k.ShouldBeMultipleOfChunkSize(amount.Amount); err != nil {
		return
	}
	if err = k.ShouldBeBondDenom(ctx, amount.Denom); err != nil {
		return
	}
	chunksToLiquidUnstake := amount.Amount.Quo(types.ChunkSize).Int64()

	chunksWithInsuranceId := make(map[uint64]types.Chunk)
	var insurances []types.Insurance
	validatorMap := make(map[string]stakingtypes.Validator)
	// Get a chunk which have most expensive insurance
	err = k.IterateAllChunks(ctx, func(chunk types.Chunk) (stop bool, err error) {
		if chunk.Status != types.CHUNK_STATUS_PAIRED {
			return false, nil
		}
		insurance, found := k.GetInsurance(ctx, chunk.InsuranceId)
		if found == false {
			return false, types.ErrNotFoundInsurance
		}

		if _, ok := validatorMap[insurance.ValidatorAddress]; !ok {
			// If validator is not in map, get validator from staking keeper
			validator, found := k.stakingKeeper.GetValidator(ctx, insurance.GetValidator())
			valid, err := k.isValidValidator(ctx, validator, found)
			if err != nil {
				return false, nil
			}
			if valid {
				validatorMap[insurance.ValidatorAddress] = validator
			} else {
				return false, nil
			}
		}
		insurances = append(insurances, insurance)
		chunksWithInsuranceId[chunk.InsuranceId] = chunk
		return false, nil
	})
	if err != nil {
		return
	}

	pairedChunks := int64(len(chunksWithInsuranceId))
	if pairedChunks == 0 {
		err = types.ErrNoPairedChunk
		return
	}
	if chunksToLiquidUnstake > pairedChunks {
		err = sdkerrors.Wrapf(
			types.ErrExceedAvailableChunks,
			"requested chunks to liquid unstake: %d, paired chunks: %d",
			chunksToLiquidUnstake,
			pairedChunks,
		)
		return
	}
	// Sort insurances by descend order
	types.SortInsurances(validatorMap, insurances, true)

	// How much ls tokens must be burned
	nas := k.GetNetAmountState(ctx)
	lsTokenBurnAmount := amount.Amount
	if nas.LsTokensTotalSupply.IsPositive() {
		lsTokenBurnAmount = amount.Amount.ToDec().Mul(nas.MintRate).TruncateInt()
	}
	liquidBondDenom := k.GetLiquidBondDenom(ctx)
	lsTokensToBurn := sdk.NewCoin(liquidBondDenom, lsTokenBurnAmount)
	// escrow
	if err = k.bankKeeper.SendCoins(
		ctx, delAddr, types.LsTokenEscrowAcc, sdk.NewCoins(lsTokensToBurn),
	); err != nil {
		return
	}
	completionTime := ctx.BlockHeader().Time.Add(k.stakingKeeper.UnbondingTime(ctx))
	for i := int64(0); i < chunksToLiquidUnstake; i++ {
		mostExpensiveInsurance := insurances[i]
		chunkToBeUndelegated := chunksWithInsuranceId[mostExpensiveInsurance.Id]
		chunkToBeUndelegated.SetStatus(types.CHUNK_STATUS_UNPAIRING_FOR_UNSTAKE)
		mostExpensiveInsurance.SetStatus(types.INSURANCE_STATUS_UNPAIRING_FOR_REPAIRING)

		var shares sdk.Dec
		shares, err = k.stakingKeeper.ValidateUnbondAmount(ctx, chunkToBeUndelegated.DerivedAddress(), mostExpensiveInsurance.GetValidator(), types.ChunkSize)
		if err != nil {
			return
		}

		if k.stakingKeeper.HasMaxUnbondingDelegationEntries(ctx, chunkToBeUndelegated.DerivedAddress(), mostExpensiveInsurance.GetValidator()) {
			err = stakingtypes.ErrMaxUnbondingDelegationEntries
			return
		}
		var unbondedNativeToken sdk.Int
		unbondedNativeToken, err = k.stakingKeeper.Unbond(
			ctx,
			chunkToBeUndelegated.DerivedAddress(),
			mostExpensiveInsurance.GetValidator(),
			shares,
		)
		if err != nil {
			return
		}
		if !validatorMap[mostExpensiveInsurance.ValidatorAddress].IsBonded() {
			coins := sdk.NewCoins(sdk.NewCoin(k.stakingKeeper.BondDenom(ctx), unbondedNativeToken))
			if err = k.bankKeeper.SendCoinsFromModuleToModule(ctx, stakingtypes.NotBondedPoolName, stakingtypes.BondedPoolName, coins); err != nil {
				return
			}
		}
		ubd := k.stakingKeeper.SetUnbondingDelegationEntry(
			ctx,
			delAddr,
			mostExpensiveInsurance.GetValidator(),
			ctx.BlockHeight(),
			completionTime,
			unbondedNativeToken,
		)
		k.stakingKeeper.InsertUBDQueue(ctx, ubd, completionTime)
		liquidUnstakeUnbondingDelegationInfo := types.NewLiquidUnstakeUnbondingDelegationInfo(
			chunkToBeUndelegated.Id,
			delAddr.String(),
			mostExpensiveInsurance.ValidatorAddress,
			lsTokensToBurn,
			completionTime,
		)
		k.SetLiquidUnstakeUnbondingDelegationInfo(ctx, liquidUnstakeUnbondingDelegationInfo)
		k.SetChunk(ctx, chunkToBeUndelegated)
		k.SetInsurance(ctx, mostExpensiveInsurance)
		unstakedChunks = append(unstakedChunks, chunkToBeUndelegated)
		unstakeUnbondingDelegationInfos = append(unstakeUnbondingDelegationInfos, liquidUnstakeUnbondingDelegationInfo)
	}
	return
}

func (k Keeper) DoInsuranceProvide(ctx sdk.Context, msg *types.MsgInsuranceProvide) (insurance types.Insurance, err error) {
	providerAddr := msg.GetProvider()
	valAddr := msg.GetValidator()
	feeRate := msg.FeeRate
	amount := msg.Amount

	// Check if the amount is greater than minimum coverage
	_, minimumCoverage := k.GetMinimumRequirements(ctx)
	if amount.Amount.LT(minimumCoverage.Amount) {
		err = sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "amount must be greater than minimum coverage: %s", minimumCoverage.String())
		return
	}

	// Check if the validator is valid
	validator, found := k.stakingKeeper.GetValidator(ctx, valAddr)
	_, err = k.isValidValidator(ctx, validator, found)
	if err != nil {
		return
	}

	// Create an insurance
	insuranceId := k.getNextInsuranceIdWithUpdate(ctx)
	insurance = types.NewInsurance(insuranceId, providerAddr.String(), valAddr.String(), feeRate)

	// Escrow provider's balance
	if err = k.bankKeeper.SendCoins(
		ctx,
		providerAddr,
		insurance.DerivedAddress(),
		sdk.NewCoins(amount),
	); err != nil {
		return
	}
	k.SetInsurance(ctx, insurance)
	k.SetPairingInsuranceIndex(ctx, insurance)
	k.SetInsurancesByProviderIndex(ctx, insurance)

	return
}

func (k Keeper) DoCancelInsuranceProvide(ctx sdk.Context, msg *types.MsgCancelInsuranceProvide) (insurance types.Insurance, err error) {
	providerAddr := msg.GetProvider()
	insuranceId := msg.Id

	// Check if the insurance exists
	insurance, found := k.GetPairingInsurance(ctx, insuranceId)
	if !found {
		err = sdkerrors.Wrapf(types.ErrPairingInsuranceNotFound, "insurance id: %d", insuranceId)
		return
	}

	// Check if the provider is the same
	if insurance.ProviderAddress != providerAddr.String() {
		err = sdkerrors.Wrapf(types.ErrNotProviderOfInsurance, "insurance id: %d", insuranceId)
		return
	}

	// Unescrow provider's balance
	escrowed := k.bankKeeper.GetBalance(ctx, insurance.DerivedAddress(), k.stakingKeeper.BondDenom(ctx))
	if err = k.bankKeeper.SendCoins(
		ctx,
		insurance.DerivedAddress(),
		providerAddr,
		sdk.NewCoins(escrowed),
	); err != nil {
		return
	}
	k.DeleteInsurance(ctx, insuranceId)
	k.DeleteInsurancesByProviderIndex(ctx, insurance)
	k.DeletePairingInsuranceIndex(ctx, insurance)
	return
}

func (k Keeper) SetLiquidBondDenom(ctx sdk.Context, denom string) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.KeyLiquidBondDenom, []byte(denom))
}

func (k Keeper) GetLiquidBondDenom(ctx sdk.Context) string {
	store := ctx.KVStore(k.storeKey)
	return string(store.Get(types.KeyLiquidBondDenom))
}

func (k Keeper) isValidValidator(ctx sdk.Context, validator stakingtypes.Validator, found bool) (bool, error) {
	if !found {
		return false, types.ErrValidatorNotFound
	}
	pubKey, err := validator.ConsPubKey()
	if err != nil {
		return false, err
	}
	if k.slashingKeeper.IsTombstoned(ctx, sdk.ConsAddress(pubKey.Address())) {
		return false, types.ErrTombstonedValidator
	}
	return true, nil
}

// Get minimum requirements for liquid staking
// Liquid staker must provide at least one chunk amount
// Insurance provider must provide at least slashing coverage
func (k Keeper) GetMinimumRequirements(ctx sdk.Context) (oneChunkAmount, slashingCoverage sdk.Coin) {
	bondDenom := k.stakingKeeper.BondDenom(ctx)
	oneChunkAmount = sdk.NewCoin(bondDenom, types.ChunkSize)
	fraction := sdk.MustNewDecFromStr(types.SlashFraction)
	slashingCoverage = sdk.NewCoin(bondDenom, oneChunkAmount.Amount.ToDec().Mul(fraction).TruncateInt())
	return
}

// ShouldBeMultipleOfChunkSize returns error if amount is not a multiple of chunk size
func (k Keeper) ShouldBeMultipleOfChunkSize(amount sdk.Int) error {
	if !amount.IsPositive() || !amount.Mod(types.ChunkSize).IsZero() {
		return sdkerrors.Wrapf(types.ErrInvalidAmount, "got: %s", amount.String())
	}
	return nil
}

// ShouldBeBondDenom returns error if denom is not the same as the bond denom
func (k Keeper) ShouldBeBondDenom(ctx sdk.Context, denom string) error {
	if denom == k.stakingKeeper.BondDenom(ctx) {
		return nil
	}
	return sdkerrors.Wrapf(types.ErrInvalidBondDenom, "expected: %s, got: %s", k.stakingKeeper.BondDenom(ctx), denom)
}
