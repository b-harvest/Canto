package cli

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/version"

	"github.com/Canto-Network/Canto/v6/x/liquidstaking/types"
)

// GetQueryCmd returns the cli query commands for the CSR module
func GetQueryCmd(queryRoute string) *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		CmdQueryParams(),
		CmdQueryEpoch(),
		CmdQueryChunks(),
		CmdQueryChunk(),
		CmdQueryInsurances(),
		CmdQueryInsurance(),
		CmdQueryWithdrawInsuranceRequests(),
		CmdQueryWithdrawInsuranceRequest(),
		CmdQueryUnpairingForUnstakingChunkInfosRequests(),
		CmdQueryUnpairingForUnstakingChunkInfosRequest(),
		CmdQueryChunkSizeRequest(),
		CmdQueryMinimumCollateral(),
		CmdQueryStates(),
	)

	return cmd
}

// CmdQueryParams implements a command that will return the current parameters of the
// liquidstaking module.
func CmdQueryParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "params",
		Short: fmt.Sprintf("Query the current parameters of %s module", types.ModuleName),
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			request := &types.QueryParamsRequest{}

			// Query store
			response, err := queryClient.Params(context.Background(), request)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(&response.Params)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// CmdQueryEpoch implements a command that will return the Epoch from the Epoch store
func CmdQueryEpoch() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "epoch",
		Short: fmt.Sprintf("Query the epoch of %s module", types.ModuleName),
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			request := &types.QueryEpochRequest{}

			// Query store
			response, err := queryClient.Epoch(context.Background(), request)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(response)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// CmdQueryChunk implements a command that will return a Chunk given a chunk id
func CmdQueryChunks() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "chunks [optional flags]",
		Args:  cobra.ExactArgs(0),
		Short: "Query all chunks",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query all chunks on a network.
Example:
$ %s query %s chunks
$ %s query %s chunks --status [CHUNK_STATUS_PAIRING | CHUNK_STATUS_PAIRED | CHUNK_STATUS_UNPAIRING | CHUNK_STATUS_UNPAIRING_FOR_UNSTAKING]
`,
				version.AppName, types.ModuleName,
				version.AppName, types.ModuleName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			pageRequest, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			request := &types.QueryChunksRequest{
				Pagination: pageRequest,
			}
			chunkStatusStr, _ := cmd.Flags().GetString(FlagChunkStatus)
			if chunkStatusStr != "" {
				status := types.ChunkStatus_value[chunkStatusStr]
				if status == 0 {
					return sdkerrors.Wrap(
						sdkerrors.ErrInvalidRequest,
						fmt.Sprintf("chunk status must be either %s, %s, %s, or %s",
							types.ChunkStatus_name[1],
							types.ChunkStatus_name[2],
							types.ChunkStatus_name[3],
							types.ChunkStatus_name[4]),
					)
				}
				request.Status = types.ChunkStatus(status)
			}

			queryClient := types.NewQueryClient(clientCtx)

			// Query store
			response, err := queryClient.Chunks(context.Background(), request)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(response)
		},
	}
	cmd.Flags().AddFlagSet(flagSetChunks())
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// CmdQueryChunk implements a command that will return a Chunk given a chunk id
func CmdQueryChunk() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "chunk [chunkId]",
		Args:    cobra.ExactArgs(1),
		Short:   "Query the Chunk associated with a given chunk id",
		Example: fmt.Sprintf("%s query %s chunk 1", version.AppName, types.ModuleName),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			// arg must be converted to a uint
			chunkId, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			request := &types.QueryChunkRequest{Id: chunkId}
			// Query store
			response, err := queryClient.Chunk(context.Background(), request)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(response)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// CmdQueryInsurances implements a command that will return insurances in liquidstaking module
func CmdQueryInsurances() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "insurances [optional flags]",
		Args:  cobra.ExactArgs(0),
		Short: "Query all insurances",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query all insurancces on a network.
Example:
$ %s query %s insurances --validator-address cantovaloper1gxl6usug4cz60yhpsjj7vw7vzysrz772yxjzsf
$ %s query %s insurances --provider-address canto1czxcryk6qw30erz3dc6ucjcvl5kp88uk3k4cj8 
$ %s query %s insurances --status [INSURANCE_STATUS_PAIRING | INSURANCE_STATUS_PAIRED | INSURANCE_STATUS_UNPAIRING | INSURANCE_STATUS_UNPAIRING_FOR_WITHDRAWAL, INSURANCE_STATUS_UNPAIRED]
$ %s query %s insurances --validator-address cantovaloper1gxl6usug4cz60yhpsjj7vw7vzysrz772yxjzsf --provider-address canto1czxcryk6qw30erz3dc6ucjcvl5kp88uk3k4cj8 
$ %s query %s insurances --validator-address cantovaloper1gxl6usug4cz60yhpsjj7vw7vzysrz772yxjzsf --provider-address canto1czxcryk6qw30erz3dc6ucjcvl5kp88uk3k4cj8 --status [INSURANCE_STATUS_PAIRING | INSURANCE_STATUS_PAIRED | INSURANCE_STATUS_UNPAIRING | INSURANCE_STATUS_UNPAIRING_FOR_WITHDRAWAL, INSURANCE_STATUS_UNPAIRED]
`,
				version.AppName, types.ModuleName,
				version.AppName, types.ModuleName,
				version.AppName, types.ModuleName,
				version.AppName, types.ModuleName,
				version.AppName, types.ModuleName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			pageRequest, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			insuranceStatusStr, _ := cmd.Flags().GetString(FlagInsuranceStatus)
			validatorAddress, _ := cmd.Flags().GetString(FlagValidatorAddress)
			providerAddress, _ := cmd.Flags().GetString(FlagProviderAddress)

			request := &types.QueryInsurancesRequest{
				Status:           types.InsuranceStatus(types.InsuranceStatus_value[insuranceStatusStr]),
				ValidatorAddress: validatorAddress,
				ProviderAddress:  providerAddress,
				Pagination:       pageRequest,
			}

			queryClient := types.NewQueryClient(clientCtx)

			// Query store
			response, err := queryClient.Insurances(context.Background(), request)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(response)
		},
	}
	cmd.Flags().AddFlagSet(flagSetInsurances())
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// CmdQueryInsurance implements a command that will return a Chunk given an insurance id
func CmdQueryInsurance() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "insurance [insuranceId]",
		Args:    cobra.ExactArgs(1),
		Short:   "Query the Insurance associated with a given insurance id",
		Example: fmt.Sprintf("%s query liquidstaking insurance 1", version.AppName),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			// arg must be converted to a uint
			insuranceId, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			request := &types.QueryInsuranceRequest{Id: insuranceId}
			// Query store
			response, err := queryClient.Insurance(context.Background(), request)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(response)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// CmdQueryWithdrawInsuranceRequests CmdQueryWithdrawRequests implements a command that will return withdraw requests in liquidstaking module.
func CmdQueryWithdrawInsuranceRequests() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "withdraw-insurance-requests [optional flags]",
		Args:  cobra.ExactArgs(0),
		Short: "Query all withdraw requests",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query details about all withdraw requests on a network.
Example:
$ %s query %s withdraw-insurance-requests 
$ %s query %s withdraw-insurance-requests --provider-address canto1czxcryk6qw30erz3dc6ucjcvl5kp88uk3k4cj8
`,
				version.AppName, types.ModuleName,
				version.AppName, types.ModuleName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			pageRequest, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			request := &types.QueryWithdrawInsuranceRequestsRequest{
				Pagination: pageRequest,
			}

			providerAddress, _ := cmd.Flags().GetString(FlagProviderAddress)
			if providerAddress != "" {
				_, err = sdk.AccAddressFromBech32(providerAddress)
				if err != nil {
					return err
				}
				request.ProviderAddress = providerAddress
			}

			queryClient := types.NewQueryClient(clientCtx)

			// Query store
			response, err := queryClient.WithdrawInsuranceRequests(context.Background(), request)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(response)
		},
	}
	cmd.Flags().AddFlagSet(flagSetWithdrawInsuranceRequests())
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// CmdQueryWithdrawInsuranceRequest CmdQueryWithdrawRequest implements a command that will return a withdraw request given an insurance id.
func CmdQueryWithdrawInsuranceRequest() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "withdraw-insurance-request [insurance-id]",
		Args:  cobra.ExactArgs(1),
		Short: "Query the withdraw request associated with a given insurance id",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query details about a withdraw request on a network.	
Example:
$ %s query %s withdraw-insurance-request 1
`,
				version.AppName, types.ModuleName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			// arg must be converted to a uint
			insuranceId, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			request := &types.QueryWithdrawInsuranceRequestRequest{Id: insuranceId}
			// Query store
			response, err := queryClient.WithdrawInsuranceRequest(context.Background(), request)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(response)
		},
	}
	return cmd
}

// CmdQueryUnpairingForUnstakingChunkInfosRequests implements a command that will return unpairing for unstaking chunk infos requests in liquidstaking module.
func CmdQueryUnpairingForUnstakingChunkInfosRequests() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "unpairing-for-unstaking-chunk-infos [optional flags]",
		Args:  cobra.ExactArgs(0),
		Short: "Query all unpairing for unstaking chunk infos",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query details about all unpairing for unstaking chunk infos on a network.
Example:
$ %s query %s unpairing-for-unstaking-chunk-infos
`,
				version.AppName, types.ModuleName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			pageRequest, err := client.ReadPageRequest(cmd.Flags())

			if err != nil {
				return err
			}
			request := &types.QueryUnpairingForUnstakingChunkInfosRequest{
				Pagination: pageRequest,
			}
			delegatorAddress, _ := cmd.Flags().GetString(FlagDelegatorAddress)
			if delegatorAddress != "" {
				_, err = sdk.AccAddressFromBech32(delegatorAddress)
				if err != nil {
					return err
				}
				request.DelegatorAddress = delegatorAddress
			}
			queryClient := types.NewQueryClient(clientCtx)
			// Query store
			response, err := queryClient.UnpairingForUnstakingChunkInfos(context.Background(), request)
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(response)
		},
	}
	cmd.Flags().AddFlagSet(flagSetUnstakingChunkInfoRequests())
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// CmdQueryUnpairingForUnstakingChunkInfosRequest implements a command that will return unpairing for unstaking chunk info in liquidstaking module.
func CmdQueryUnpairingForUnstakingChunkInfosRequest() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "unpairing-for-unstaking-chunk-info [chunk-id]",
		Args:  cobra.ExactArgs(1),
		Short: "Query the unpairing for unstaking chunk info associated with a given chunk id",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query details about a unpairing for unstaking chunk info on a network.
Example:
$ %s query %s unpairing-for-unstaking-chunk-info 1
`,
				version.AppName, types.ModuleName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			// arg must be converted to a uint
			chunkId, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)
			request := &types.QueryUnpairingForUnstakingChunkInfoRequest{Id: chunkId}
			// Query store
			response, err := queryClient.UnpairingForUnstakingChunkInfo(context.Background(), request)
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(response)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// CmdQueryChunkSizeRequest implements a command that will return chunk size in liquidstaking module.
func CmdQueryChunkSizeRequest() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "chunk-size",
		Args:  cobra.ExactArgs(0),
		Short: "Query the chunk size tokens(=how many tokens are needed to create a chunk)",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query the chunk size on a network.
Example:
$ %s query %s chunk-size
`,
				version.AppName, types.ModuleName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			request := &types.QueryChunkSizeRequest{}
			// Query store
			response, err := queryClient.ChunkSize(context.Background(), request)
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(response)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func CmdQueryMinimumCollateral() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "minimum-collateral",
		Args:  cobra.ExactArgs(0),
		Short: "Query the minimum collateral",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query the minimum collateral on a network.
Example:
$ %s query %s minimum-collateral
`,
				version.AppName, types.ModuleName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)
			request := &types.QueryMinimumCollateralRequest{}
			// Query store
			response, err := queryClient.MinimumCollateral(context.Background(), request)
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(response)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// CmdQueryStates implements a command that will return states in liquidstaking module.
func CmdQueryStates() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "states",
		Args:  cobra.ExactArgs(0),
		Short: "Query the states of liquid staking module",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query the states of liquid staking module on a network.
Example:
$ %s query %s states
`,
				version.AppName, types.ModuleName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			request := &types.QueryStatesRequest{}
			// Query store
			response, err := queryClient.States(context.Background(), request)
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(response)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}
