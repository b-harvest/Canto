package v8

import (
	"context"

	"cosmossdk.io/collections"
	sdkmath "cosmossdk.io/math"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	csrkeeper "github.com/Canto-Network/Canto/v7/x/csr/keeper"
	govshuttlekeeper "github.com/Canto-Network/Canto/v7/x/govshuttle/keeper"
	"github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	clientkeeper "github.com/cosmos/ibc-go/v8/modules/core/02-client/keeper"
	"github.com/cosmos/ibc-go/v8/modules/core/exported"
	"github.com/ethereum/go-ethereum/common"
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
	csrKeeper csrkeeper.Keeper,
	govshuttleKeeper govshuttlekeeper.Keeper,
) upgradetypes.UpgradeHandler {
	return func(ctx context.Context, _ upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		sdkCtx := sdk.UnwrapSDKContext(ctx)
		logger := sdkCtx.Logger().With("upgrade: ", UpgradeName)

		// Leave modules are as-is to avoid running InitGenesis.
		logger.Debug("running module migrations ...")
		if vm, err := mm.RunMigrations(ctx, configurator, vm); err != nil {
			return vm, err
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
			return vm, err
		}

		// canto v8 custom
		{
			params, err := stakingKeeper.GetParams(ctx)
			if err != nil {
				return vm, err
			}
			params.MinCommissionRate = MinCommissionRate
			stakingKeeper.SetParams(ctx, params)

			// hardcode for missing states
			csrs := csrKeeper.GetAllCSRs(sdkCtx)
			if len(csrs) == 0 {
				csrs = GetCSRState(logger)
			}

			for _, csr := range csrs {
				csrKeeper.SetCSR(sdkCtx, csr)
			}

			_, exist := csrKeeper.GetTurnstile(sdkCtx)
			if !exist {
				csrKeeper.SetTurnstile(sdkCtx, common.HexToAddress(TunstileState))
			}

			_, exist = govshuttleKeeper.GetPort(sdkCtx)
			if !exist {
				govshuttleKeeper.SetPort(sdkCtx, common.HexToAddress(PortState))
			}

		}

		return vm, nil
	}
}
