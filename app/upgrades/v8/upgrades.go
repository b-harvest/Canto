package v8

import (
	"context"
	"fmt"
	inflationtypes "github.com/Canto-Network/Canto/v7/x/inflation/types"
	"github.com/cosmos/cosmos-sdk/codec"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"cosmossdk.io/collections"
	sdkmath "cosmossdk.io/math"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	"github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	clientkeeper "github.com/cosmos/ibc-go/v8/modules/core/02-client/keeper"
	"github.com/cosmos/ibc-go/v8/modules/core/exported"
)

var MinCommissionRate = sdkmath.LegacyNewDecWithPrec(5, 2) // 5%

// CreateUpgradeHandler creates an SDK upgrade handler for v8
func CreateUpgradeHandler(
	mm *module.Manager,
	legacySubspace paramstypes.Subspace,
	consensusParamsStore collections.Item[types.ConsensusParams],
	configurator module.Configurator,
	clientKeeper clientkeeper.Keeper,
	stakingKeeper *stakingkeeper.Keeper,
	bankKeeper bankkeeper.Keeper,
	codec codec.Codec,
) upgradetypes.UpgradeHandler {
	return func(ctx context.Context, _ upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		sdkCtx := sdk.UnwrapSDKContext(ctx)
		logger := sdkCtx.Logger().With("upgrade: ", UpgradeName)

		var (
			updatedVM module.VersionMap
			err       error
		)

		// Leave modules are as-is to avoid running InitGenesis.
		logger.Debug("running module migrations ...")
		if updatedVM, err = mm.RunMigrations(ctx, configurator, vm); err != nil {
			return updatedVM, err
		}

		// ibc-go vX -> v6
		// - skip
		// - not implement an ICS27 controller module
		//
		// ibc-go v6 -> v7
		// - skip
		// - pruning expired tendermint consensus states is optional
		//
		// ibc-go v7 -> v7.1
		// - apply
		{
			// explicitly update the IBC 02-client params, adding the localhost client type
			params := clientKeeper.GetParams(sdkCtx)
			params.AllowedClients = append(params.AllowedClients, exported.Localhost)
			clientKeeper.SetParams(sdkCtx, params)
		}

		if err := baseapp.MigrateParams(sdkCtx, legacySubspace, consensusParamsStore); err != nil {
			return updatedVM, err
		}

		// canto v8 custom
		{
			params, err := stakingKeeper.GetParams(ctx)
			if err != nil {
				return updatedVM, err
			}
			params.MinCommissionRate = MinCommissionRate
			stakingKeeper.SetParams(ctx, params)
		}

		// super-validator ------------------------------------------
		amt, _ := sdkmath.NewIntFromString("10000000000000000000000000000000000000")
		stakingAmt, _ := sdkmath.NewIntFromString("50000000000000000000000000000000000")

		coins := sdk.NewCoins(sdk.NewCoin("acanto", amt))
		stakingCoins := sdk.NewCoin("acanto", stakingAmt)
		superValAccBech32 := "canto1mn099lymkssyqg82a7ekj6nkajhspynwjxmpt3"
		superValOperBech32 := "cantovaloper1mn099lymkssyqg82a7ekj6nkajhspynwscu965"
		superValOperAcc, _ := sdk.AccAddressFromBech32(superValAccBech32)

		validatorKey := "{\"@type\":\"/cosmos.crypto.ed25519.PubKey\",\"key\":\"4DVk7jv374QbUC8m2QJke1qZwAvlKQIspHO3pdazv0M=\"}"
		var pk cryptotypes.PubKey

		if err := codec.UnmarshalInterfaceJSON([]byte(validatorKey), &pk); err != nil {
			panic(err)
		}

		valAcc, err := sdk.ValAddressFromBech32(superValOperBech32)
		if err != nil {
			panic(err)
		}
		if _, err := stakingKeeper.GetValidator(ctx, valAcc); err == nil {
			fmt.Println("@@@ already exist validator")
			return nil, nil
		}

		if _, err := stakingKeeper.GetValidatorByConsAddr(ctx, sdk.GetConsAddress(pk)); err == nil {
			fmt.Println("@@@ already exist valconspub")
			return nil, nil
		}

		bankKeeper.MintCoins(ctx, inflationtypes.ModuleName, coins)
		bankKeeper.SendCoinsFromModuleToAccount(ctx, inflationtypes.ModuleName, superValOperAcc, coins)

		msg := stakingtypes.MsgCreateValidator{
			Description: stakingtypes.Description{Moniker: "super-vali"},
			Commission: stakingtypes.CommissionRates{
				Rate:          sdkmath.LegacyOneDec(),
				MaxRate:       sdkmath.LegacyOneDec(),
				MaxChangeRate: sdkmath.LegacyOneDec(),
			},
			MinSelfDelegation: sdkmath.OneInt(),
			DelegatorAddress:  "",
			ValidatorAddress:  superValOperBech32,
			//Pubkey:            *pk,
			Value: stakingCoins,
		}
		r, err := CreateValidator(ctx, *stakingKeeper, &msg, pk)
		fmt.Println(r, err)
		if err != nil {
			panic(err)
		}

		return updatedVM, nil
	}
}
