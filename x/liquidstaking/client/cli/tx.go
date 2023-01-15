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
		NewLiquidStakingCmd(),
		NewCancelLiquidStakingCmd(),
		NewLiquidUnstakingCmd(),
		NewCancelLiquidUnstakingCmd(),
		NewBidInsuranceCmd(),
		NewCancelInsuranceBidCmd(),
		NewUnbondInsuranceCmd(),
		NewCancelInsuranceUnbondCmd(),
	)
	return txCmd
}

// NewLiquidStakingCmd returns a CLI command handler for liquid staking chunks.
func NewLiquidStakingCmd() *cobra.Command {
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
			tokenAmount, ok := sdk.NewIntFromString(args[1])
			if !ok {
				return fmt.Errorf("invalid amount %s", args[1])
			}

			msg := &types.MsgLiquidStaking{
				RequesterAddress: liquidStaker.String(),
				TokenAmount:      tokenAmount,
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func NewCancelLiquidStakingCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cancel-liquid-stake num_chunk",
		Args:  cobra.ExactArgs(1),
		Short: "Cancel_liquid-stake num_chunk",
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
			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			msg := &types.MsgCancelLiquidStaking{
				RequesterAddress:   liquidStaker.String(),
				ChunkBondRequestId: id,
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// NewLiquidStakeCmd returns a CLI command handler for liquid unstaking chunks.
func NewLiquidUnstakingCmd() *cobra.Command {
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

			msg := &types.MsgLiquidUnstaking{RequesterAddress: liquidStaker.String(), NumChunkUnstake: numChunks}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func NewCancelLiquidUnstakingCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cancel-liquid-unstaking num_chunk",
		Args:  cobra.ExactArgs(1),
		Short: "Cancel_liquid-unstaking num_chunk",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Cancel-Liquid-staking num_chunk.

Example:
$ %s tx %s cancel-liquid-staking 2 --from mykey
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
			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			msg := &types.MsgCancelLiquidUnstaking{
				// TODO:
				RequesterAddress:     liquidStaker.String(),
				ChunkUnbondRequestId: id,
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// NewBidInsuranceCmd returns a CLI command handler for registering insurance.
func NewBidInsuranceCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "register-insurance amount",
		Args:  cobra.ExactArgs(4),
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
			validator, err := sdk.ValAddressFromBech32(args[1])
			if err != nil {
				return err
			}
			insuranceFeeRate, err := sdk.NewDecFromStr(args[2])
			if err != nil {
				return err
			}
			amount, err := sdk.ParseCoinNormalized(args[3])
			if err != nil {
				return err
			}
			if err != nil {
				return err
			}

			msg := &types.MsgBidInsurance{
				RequesterAddress: insurer.String(),
				ValidatorAddress: validator.String(),
				InsuranceFeeRate: insuranceFeeRate,
				Amount:           amount,
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func NewCancelInsuranceBidCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cancel-insurance-bid num_chunk",
		Args:  cobra.ExactArgs(1),
		Short: "Cancel_insurance-bid num_chunk",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Cancel-insurance-bid num_chunk.

Example:
$ %s tx %s cancel-insurance-bid 2 --from mykey
`,
				version.AppName, types.ModuleName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			msg := &types.MsgCancelInsuranceBid{
				BidId: id,
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// NewUnbondInsuranceCmd returns a CLI command handler for unregistering given insurance.
func NewUnbondInsuranceCmd() *cobra.Command {
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
			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("parse pair id: %w", err)
			}

			msg := &types.MsgUnbondInsurance{RequesterAddress: insurer.String(), AliveChunkId: id}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func NewCancelInsuranceUnbondCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cancel-insurance-unbond num_chunk",
		Args:  cobra.ExactArgs(1),
		Short: "Cancel_insurance-unbond num_chunk",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Cancel-insurance-unbond num_chunk.

Example:
$ %s tx %s cancel-insurance-bid 2 --from mykey
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
			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			msg := &types.MsgCancelInsuranceUnbond{
				RequesterAddress: insurer.String(),
				UnbondRequestId:  id,
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
