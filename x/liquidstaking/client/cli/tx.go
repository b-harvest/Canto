package cli

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Canto-Network/Canto-Testnet-v2/v1/x/liquidstaking/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/spf13/cobra"
)

// NewTxCmd returns a root CLI command handler for certain modules/liquidstaking transaction commands.
func NewTxCmd() *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "liquidstaking subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	txCmd.AddCommand(
		NewLiquidStakeCmd(),
		NewLiquidUnstakeCmd(),
		NewRegisterInsuranceCmd(),
		NewUnregisterInsuranceCmd(),
	)
	return txCmd
}

// NewLiquidStakeCmd returns a CLI command handler for liquid staking chunks.
func NewLiquidStakeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "liquid-stake num_chunk",
		Args:  cobra.ExactArgs(1),
		Short: "Liquid-stake num_chunk",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Liquid-stake num_chunk.

Example:
$ %s tx %s liquid-stake 2 --from mykey
`,
				version.AppName, types.ModuleName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			liquidStaker := clientCtx.GetFromAddress()
			numChunks, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return err
			}

			msg := &types.MsgLiquidStake{
				DelegatorAddress: liquidStaker.String(),
				NumChunks:        numChunks,
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// NewLiquidStakeCmd returns a CLI command handler for liquid unstaking chunks.
func NewLiquidUnstakeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "liquid-unstake num_chunk",
		Args:  cobra.ExactArgs(1),
		Short: "Liquid-unstake num_chunk",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Liquid-unstake num_chunk.

Example:
$ %s tx %s liquid-unstake 3 --from mykey
`,
				version.AppName, types.ModuleName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			liquidStaker := clientCtx.GetFromAddress()
			numChunks, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("parse pair id: %w", err)
			}

			msg := &types.MsgLiquidUnstake{DelegatorAddress: liquidStaker.String(), NumChunks: numChunks}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// NewRegisterInsuranceCmd returns a CLI command handler for registering insurance.
func NewRegisterInsuranceCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "register-insurance amount",
		Args:  cobra.ExactArgs(1),
		Short: "Register-insurance amount",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Register-insurance amount.

Example:
$ %s tx %s register-insurance 500stake --from mykey
`,
				version.AppName, types.ModuleName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			insurer := clientCtx.GetFromAddress()

			amount, err := sdk.ParseCoinNormalized(args[0])
			if err != nil {
				return err
			}

			msg := &types.MsgRegisterInsurance{InsurerAddress: insurer.String(), Amount: amount}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// NewUnregisterInsuranceCmd returns a CLI command handler for unregistering given insurance.
func NewUnregisterInsuranceCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "unregister-insurance insurance_id...",
		Args:  cobra.ExactArgs(1),
		Short: "Unregister-insurance insurance_id...",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Unregister-insurance insurance_id.

Example:
$ %s tx %s register-insurance 3,4 --from mykey
`,
				version.AppName, types.ModuleName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			insurer := clientCtx.GetFromAddress()

			var ids []uint64
			for _, idStr := range strings.Split(args[0], ",") {
				id, err := strconv.ParseUint(idStr, 10, 64)
				if err != nil {
					return fmt.Errorf("parse pair id: %w", err)
				}
				ids = append(ids, id)
			}

			msg := &types.MsgUnregisterInsurance{InsurerAddress: insurer.String(), InsuranceIds: ids}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
