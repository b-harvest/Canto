package v6

import (
	onboardingkeeper "github.com/Canto-Network/Canto/v6/x/onboarding/keeper"
	coinswapkeeper "github.com/b-harvest/coinswap/modules/coinswap/keeper"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
)

// CreateUpgradeHandler creates an SDK upgrade handler for v2
func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	onboardingKeeper onboardingkeeper.Keeper,
	coinswapKeeper coinswapkeeper.Keeper,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		logger := ctx.Logger().With("upgrading to v6.0.0", UpgradeName)

		params := onboardingKeeper.GetParams(ctx)
		params.WhitelistedChannels = []string{"channel-0"}
		params.AutoSwapThreshold = sdk.NewIntWithDecimal(4, 18)
		onboardingKeeper.SetParams(ctx, params)

		coinswapParams := coinswapKeeper.GetParams(ctx)
		coinswapParams.PoolCreationFee = sdk.NewCoin("uatom", sdk.ZeroInt())
		coinswapParams.MaxSwapAmount = sdk.NewCoins(
			sdk.NewCoin(usdcIBCDenom, sdk.NewIntWithDecimal(10, 6)),
			sdk.NewCoin(usdtIBCDenom, sdk.NewIntWithDecimal(10, 6)),
			sdk.NewCoin(ethIBCDenom, sdk.NewIntWithDecimal(1, 17)),
		)
		coinswapKeeper.SetParams(ctx, coinswapParams)
		coinswapKeeper.SetStandardDenom(ctx, "acanto")

		// Leave modules are as-is to avoid running InitGenesis.
		logger.Debug("running module migrations ...")
		return mm.RunMigrations(ctx, configurator, vm)
	}
}
