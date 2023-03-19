package keeper

import (
	"github.com/Canto-Network/Canto/v6/x/liquidstaking/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

func (k Keeper) DoLiquidStake(ctx sdk.Context, delAddr sdk.AccAddress, amount sdk.Coin) (newShares sdk.Dec, lsTokenMintAmount sdk.Int, err error) {
	// Check if max paired chunk size is exceeded or not
	currenPairedChunkSize := 0
	err = k.IterateAllChunks(ctx, func(chunk types.Chunk) (bool, error) {
		if chunk.Status == types.CHUNK_STATUS_PAIRED {
			currenPairedChunkSize++
		}
		return false, nil
	})
	if err != nil {
		return
	}
	availableChunks := types.MaxPairedChunks - currenPairedChunkSize
	if availableChunks <= 0 {
		err = sdkerrors.Wrapf(types.ErrMaxPairedChunkSizeExceeded, "current paired chunk size: %d", currenPairedChunkSize)
		return
	}

	bondDenom := k.stakingKeeper.BondDenom(ctx)
	minimumRequirement := sdk.NewCoin(bondDenom, sdk.NewInt(types.ChunkSize))
	// amount must be greater than or equal to one chunk size
	if !amount.IsGTE(minimumRequirement) {
		err = sdkerrors.Wrapf(types.ErrInvalidAmount, "amount must be greater than or equal to %s", minimumRequirement.String())
		return
	}

	// Check if there are any pairing insurances
	var pairingInsurances []types.Insurance
	validatorMap := make(map[string]stakingtypes.Validator)
	err = k.IterateAllInsurances(ctx, func(insurance types.Insurance) (bool, error) {
		if insurance.Status == types.INSURANCE_STATUS_PAIRING {
			// Store validator of insurance to map
			if _, ok := validatorMap[insurance.ValidatorAddress]; !ok {
				// If validator is not in map, get validator from staking keeper
				validator, found := k.stakingKeeper.GetValidator(ctx, sdk.ValAddress(insurance.ValidatorAddress))
				if found && !validator.IsJailed() {
					validatorMap[insurance.ValidatorAddress] = validator
				} else {
					return false, nil
				}
			} else {
				pairingInsurances = append(pairingInsurances, insurance)
			}
		}
		return false, nil
	})
	if err != nil {
		return
	}
	if len(pairingInsurances) == 0 {
		err = types.ErrNoPairingInsurance
		return
	}

	// Liquid stakers can send amount of tokens corresponding multiple of chunk size and create multiple chunks
	// Check the liquid staker's balance
	n := amount.Amount.Quo(minimumRequirement.Amount).Int64()
	amount = sdk.NewCoin(bondDenom, sdk.NewInt(n*types.ChunkSize))
	if !k.bankKeeper.HasBalance(ctx, delAddr, amount) {
		err = sdkerrors.Wrapf(sdkerrors.ErrInsufficientFunds, "insufficient funds to liquid stake")
		return
	}

	if n > int64(availableChunks) {
		n = int64(availableChunks)
		amount = sdk.NewCoin(bondDenom, sdk.NewInt(n*types.ChunkSize))
	}

	types.SortInsurances(validatorMap, pairingInsurances)
	totalNewShares := sdk.Dec{}
	totalLsTokenMintAmount := sdk.Int{}
	for i := int64(0); i < n; i++ {
		// We can create paired chunk only with available pairing insurances
		if len(pairingInsurances) == 0 {
			break
		}
		cheapestInsurance := pairingInsurances[0]
		pairingInsurances = pairingInsurances[1:]

		// Now we have the cheapest pairing insurance and valid msg liquid stake! Let's create a chunk
		// Create a chunk
		chunkId := k.getNextChunkIdWithUpdate(ctx)
		chunk := types.NewChunk(chunkId)

		// Escrow liquid staker's coins
		if err = k.bankKeeper.SendCoinsFromAccountToModule(
			ctx,
			delAddr,
			chunk.DerivedAddress().String(),
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
		totalNewShares.Add(newShares)

		nas := k.GetNetAmountState(ctx)
		// TODO: bond denom must be set at Genesis
		liquidBondDenom := k.GetParams(ctx).LiquidBondDenom
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
		totalLsTokenMintAmount.Add(lsTokenMintAmount)
		if err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, delAddr, sdk.NewCoins(mintedCoin)); err != nil {
			return
		}
		chunk.Status = types.CHUNK_STATUS_PAIRED
		cheapestInsurance.Status = types.INSURANCE_STATUS_PAIRED
		k.SetChunk(ctx, chunk)
		k.SetInsurance(ctx, cheapestInsurance)
	}

	return totalNewShares, totalLsTokenMintAmount, err
}
