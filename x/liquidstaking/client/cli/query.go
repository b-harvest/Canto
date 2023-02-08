package cli

import (
	"context"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/Canto-Network/Canto-Testnet-v2/v1/x/liquidstaking/types"
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
		GetCmdQueryAliveChunk(),
		GetCmdQueryAliveChunks(),
		GetCmdQueryAliveChunksByInsuranceProvider(),
		GetCmdQueryAliveChunksByValidator(),
		GetCmdQueryUnbondingChunks(),
		GetCmdQueryChunkBondRequest(),
		GetCmdQueryChunkBondRequests(),
		GetCmdQueryChunkBondRequestsByDelegator(),
		GetCmdQueryChunkUnbondRequest(),
		GetCmdQueryChunkUnbondRequests(),
		GetCmdQueryChunkUnbondRequestsByUndelegator(),
		GetCmdQueryInsuranceBid(),
		GetCmdQueryInsuranceBids(),
		GetCmdQueryInsuranceBidsByInsuranceProvider(),
		GetCmdQueryInsuranceBidsByValidator(),
		GetCmdQueryInsuranceUnbondRequest(),
		GetCmdQueryInsuranceUnbondRequests(),
		GetCmdQueryInsuranceUnbondRequestsByInsuranceProvider(),
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

func GetCmdQueryAliveChunk() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "alive-chunks",
		Short: "Query the current chunks",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			request := &types.QueryAliveChunkRequest{
				Id: id,
			}

			res, err := queryClient.AliveChunk(context.Background(), request)
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

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}
			request := &types.QueryAliveChunksRequest{
				Pagination: pageReq,
			}

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

func GetCmdQueryAliveChunksByInsuranceProvider() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "alive-chunks",
		Short: "Query the current chunks",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			addr, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}
			request := &types.QueryAliveChunksByInsuranceProviderRequest{
				InsuranceProviderAddr: addr.String(),
				Pagination:            pageReq,
			}

			res, err := queryClient.AliveChunksByInsuranceProvider(context.Background(), request)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func GetCmdQueryAliveChunksByValidator() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "alive-chunks",
		Short: "Query the current chunks",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			addr, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}
			request := &types.QueryAliveChunksByValidatorRequest{
				ValidatorAddr: addr.String(),
				Pagination:    pageReq,
			}

			res, err := queryClient.AliveChunksByValidator(context.Background(), request)
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

func GetCmdQueryChunkBondRequest() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "chunk-bond-requests",
		Short: "Query the current chunks",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			request := &types.QueryChunkBondRequestRequest{
				Id: id,
			}

			res, err := queryClient.ChunkBondRequest(context.Background(), request)
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

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}
			request := &types.QueryChunkBondRequestsRequest{
				Pagination: pageReq,
			}

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

func GetCmdQueryChunkBondRequestsByDelegator() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "chunk-bond-requests",
		Short: "Query the current chunks",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			delegator, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			request := &types.QueryChunkBondRequestsByDelegatorRequest{
				DelegatorAddr: delegator.String(),
				Pagination:    pageReq,
			}

			res, err := queryClient.ChunkBondRequestsByDelegator(context.Background(), request)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func GetCmdQueryChunkUnbondRequest() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "chunk-unbond-requests",
		Short: "Query the current chunks",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}
			request := &types.QueryChunkUnbondRequestRequest{
				Id: id,
			}

			res, err := queryClient.ChunkUnbondRequest(context.Background(), request)
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

			request := &types.QueryChunkUnbondRequestsRequest{}

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

func GetCmdQueryChunkUnbondRequestsByUndelegator() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "chunk-unbond-requests",
		Short: "Query the current chunks",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			delegator, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			request := &types.QueryChunkUnbondRequestsByUndelegatorRequest{
				DelegatorAddr: delegator.String(),
				Pagination:    pageReq,
			}

			res, err := queryClient.ChunkUnbondRequestsByUndelegator(context.Background(), request)
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
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)
			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}
			request := &types.QueryInsuranceBidRequest{
				Id: id,
			}

			res, err := queryClient.InsuranceBid(context.Background(), request)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func GetCmdQueryInsuranceBids() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "insurance-bids",
		Short: "Query the current insurances",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}
			request := &types.QueryInsuranceBidsRequest{
				Pagination: pageReq,
			}

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

func GetCmdQueryInsuranceBidsByInsuranceProvider() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "insurance-bids",
		Short: "Query the current insurances",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			addr, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}
			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			request := &types.QueryInsuranceBidsByInsuranceProviderRequest{
				InsuranceProviderAddr: addr.String(),
				Pagination:            pageReq,
			}

			res, err := queryClient.InsuranceBidsByInsuranceProvider(context.Background(), request)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func GetCmdQueryInsuranceBidsByValidator() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "insurance-bids",
		Short: "Query the current insurances",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			addr, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}
			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			request := &types.QueryInsuranceBidsByValidatorRequest{
				ValidatorAddr: addr.String(),
				Pagination:    pageReq,
			}

			res, err := queryClient.InsuranceBidsByValidator(context.Background(), request)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func GetCmdQueryInsuranceUnbondRequest() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "insurance-bid",
		Short: "Query the current insurances",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)
			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}
			request := &types.QueryInsuranceUnbondRequestRequest{
				Id: id,
			}

			res, err := queryClient.InsuranceUnbondRequest(context.Background(), request)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func GetCmdQueryInsuranceUnbondRequests() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "insurance-bids",
		Short: "Query the current insurances",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}
			request := &types.QueryInsuranceUnbondRequestsRequest{
				Pagination: pageReq,
			}

			res, err := queryClient.InsuranceUnbondRequests(context.Background(), request)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func GetCmdQueryInsuranceUnbondRequestsByInsuranceProvider() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "insurance-bids",
		Short: "Query the current insurances",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			addr, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}
			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}
			request := &types.QueryInsuranceUnbondRequestsByInsuranceProviderRequest{
				InsuranceProviderAddr: addr.String(),
				Pagination:            pageReq,
			}

			res, err := queryClient.InsuranceUnbondRequestsByInsuranceProvider(context.Background(), request)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
