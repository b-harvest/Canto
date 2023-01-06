package cli

import (
	"context"
	"fmt"

	"github.com/Canto-Network/Canto-Testnet-v2/v1/x/liquidstaking/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd() *cobra.Command {
	feesQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	feesQueryCmd.AddCommand(
		GetCmdQueryParams(),
		GetCmdQueryLiquidValidators(),
		GetCmdQueryLiquidStakingState(),
		GetCmdQueryAliveChunks(),
		GetCmdQueryUnbondingChunks(),
		GetCmdQueryChunkBondRequests(),
		GetCmdQueryChunkUnbondRequests(),
		GetCmdQueryInsuranceBid(),
	)

	return feesQueryCmd
}

// GetCmdQueryParams implements a command to return the current liquidstaking
// parameters.
func GetCmdQueryParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "params",
		Short: "Query the current liquidstaking parameters",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			request := &types.QueryParamsRequest{}

			res, err := queryClient.Params(context.Background(), request)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(&res.Params)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetCmdQueryParams implements a command to return the current
// liquid validators.
func GetCmdQueryLiquidValidators() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validators",
		Short: "Query the current liquid validators",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			request := &types.QueryLiquidValidatorsRequest{}

			res, err := queryClient.LiquidValidators(context.Background(), request)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetCmdQueryParams implements a command to return the current
// liquidstaking state.
func GetCmdQueryLiquidStakingState() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "state",
		Short: "Query the current liquidstaking state",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			request := &types.QueryLiquidStakingStateRequest{}

			res, err := queryClient.LiquidStakingState(context.Background(), request)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetCmdQueryParams implements a command to return the current
// Chunks.
func GetCmdQueryAliveChunks() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "alive-chunks",
		Short: "Query the current chunks",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			request := &types.QueryAliveChunksRequest{}

			res, err := queryClient.AliveChunks(context.Background(), request)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func GetCmdQueryUnbondingChunks() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "unbonding-chunks",
		Short: "Query the current chunks",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			request := &types.QueryUnbondingChunksRequest{}

			res, err := queryClient.UnbondingChunks(context.Background(), request)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func GetCmdQueryChunkBondRequests() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "chunk-bond-requests",
		Short: "Query the current chunks",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			request := &types.QueryChunkBondRequests{}

			res, err := queryClient.ChunkBondRequests(context.Background(), request)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func GetCmdQueryChunkUnbondRequests() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "chunk-unbond-requests",
		Short: "Query the current chunks",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			request := &types.QueryChunkUnbondRequests{}

			res, err := queryClient.ChunkUnbondRequests(context.Background(), request)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetCmdQueryParams implements a command to return the current
// Insurances.
func GetCmdQueryInsuranceBid() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "insurance-bid",
		Short: "Query the current insurances",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			request := &types.QueryInsuranceBidRequest{}

			res, err := queryClient.InsuranceBids(context.Background(), request)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
