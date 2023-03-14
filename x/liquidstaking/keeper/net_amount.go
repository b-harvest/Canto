package keeper

import (
	"github.com/Canto-Network/Canto/v6/x/liquidstaking/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) GetNetAmountState(ctx sdk.Context) (nas types.NetAmountState) {
	// TODO: Fill calculate logic
	nas = types.NetAmountState{
		MintRate:               sdk.Dec{},
		LsTokensTotalSupply:    sdk.Int{},
		NetAmount:              sdk.Dec{},
		TotalDelShares:         sdk.Dec{},
		TotalRemainingRewards:  sdk.Dec{},
		TotalChunksBalance:     sdk.Int{},
		TotalLiquidTokens:      sdk.Int{},
		TotalInsuranceTokens:   sdk.Int{},
		TotalUnbondingBalance:  sdk.Int{},
		RewardModuleAccBalance: sdk.Int{},
	}
	return
}
