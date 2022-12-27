package cli

import (
	"github.com/Canto-Network/Canto-Testnet-v2/v1/x/liquidstaking/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/spf13/cobra"
)

// NewTxCmd returns a root CLI command handler for certain modules/erc20 transaction commands.
func NewTxCmd() *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "liquidstaking subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	txCmd.AddCommand()
	return txCmd
}
