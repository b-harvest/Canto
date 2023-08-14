package keeper

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"os"
	"path/filepath"
	"time"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/Canto-Network/Canto/v7/x/liquidstaking/types"
)

// Keeper of the inflation store
type Keeper struct {
	storeKey   sdk.StoreKey
	cdc        codec.BinaryCodec
	paramstore paramtypes.Subspace

	accountKeeper      types.AccountKeeper
	bankKeeper         types.BankKeeper
	distributionKeeper types.DistributionKeeper
	stakingKeeper      types.StakingKeeper
	slashingKeeper     types.SlashingKeeper
	evidenceKeeper     types.EvidenceKeeper

	fileLogger log.Logger
}

// NewKeeper creates a new mint Keeper instance
func NewKeeper(
	storeKey sdk.StoreKey,
	cdc codec.BinaryCodec,
	subspace paramtypes.Subspace,
	accountKeeper types.AccountKeeper,
	bankKeeper types.BankKeeper,
	distributionKeeper types.DistributionKeeper,
	stakingKeeper types.StakingKeeper,
	slashingKeeper types.SlashingKeeper,
	evidenceKeeper types.EvidenceKeeper,
) Keeper {
	// set KeyTable if it has not already been set
	if !subspace.HasKeyTable() {
		subspace = subspace.WithKeyTable(types.ParamKeyTable())
	}
	file, err := os.OpenFile(filepath.Join(os.Getenv("HOME"), "logs", fmt.Sprintf("liquidstaking.log-%s", time.Now().Format(time.RFC3339))), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	fileLogger := log.NewTMLogger(log.NewSyncWriter(file)).With("module", "x/liquidstaking")
	return Keeper{
		storeKey:           storeKey,
		cdc:                cdc,
		paramstore:         subspace,
		accountKeeper:      accountKeeper,
		bankKeeper:         bankKeeper,
		distributionKeeper: distributionKeeper,
		stakingKeeper:      stakingKeeper,
		slashingKeeper:     slashingKeeper,
		evidenceKeeper:     evidenceKeeper,

		fileLogger: fileLogger,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", "x/"+types.ModuleName)
}
