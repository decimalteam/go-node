package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"

	"bitbucket.org/decimalteam/go-node/x/validator/internal/types"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd(queryRoute string, cdc *codec.Codec) *cobra.Command {
	// Group validator queries under a subcommand
	validatorQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	validatorQueryCmd.AddCommand(
		flags.GetCommands(
			GetCmdQueryValidators(queryRoute, cdc),
			GetCmdQueryValidatorDelegations(queryRoute, cdc),
			GetCmdQueryPool(queryRoute, cdc),
		)...,
	)

	return validatorQueryCmd

}

// GetCmdQueryValidators implements the query all validators command.
func GetCmdQueryValidators(storeName string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "validators",
		Short: "Query for all validators",
		Args:  cobra.NoArgs,
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query details about all validators on a network.

Example:
$ %s query validator validators
`,
				version.ClientName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			resKVs, _, err := cliCtx.QuerySubspace([]byte{types.ValidatorsKey}, storeName)
			if err != nil {
				return err
			}

			var validators types.Validators
			for _, kv := range resKVs {
				validator, err := types.UnmarshalValidator(cdc, kv.Value)
				if err != nil {
					return err
				}
				validators = append(validators, validator)
			}

			return cliCtx.PrintOutput(validators)
		},
	}
}

// GetCmdQueryValidatorDelegations implements the command to query all the
// delegations to a specific validator.
func GetCmdQueryValidatorDelegations(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "delegations-to [validator-addr]",
		Short: "Query all delegations made to one validator",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query delegations on an individual validator.

Example:
$ %s query staking delegations-to dxvaloper1gghjut3ccd8ay0zduzj64hwre2fxs9ldmqhffj
`,
				version.ClientName,
			),
		),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			valAddr, err := sdk.ValAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			bz, err := cdc.MarshalJSON(types.NewQueryValidatorParams(valAddr))
			if err != nil {
				return err
			}

			route := fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryValidatorDelegations)
			res, _, err := cliCtx.QueryWithData(route, bz)
			if err != nil {
				return err
			}

			var resp types.DelegationResponses
			if err := cdc.UnmarshalJSON(res, &resp); err != nil {
				return err
			}

			return cliCtx.PrintOutput(resp)
		},
	}
}

// GetCmdQueryPool implements the pool query command.
func GetCmdQueryPool(storeName string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "pool",
		Args:  cobra.NoArgs,
		Short: "Query the current staking pool values",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query values for amounts stored in the staking pool.

Example:
$ %s query staking pool
`,
				version.ClientName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			bz, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/pool", storeName), nil)
			if err != nil {
				return err
			}

			var pool types.Pool
			if err := cdc.UnmarshalJSON(bz, &pool); err != nil {
				return err
			}

			return cliCtx.PrintOutput(pool)
		},
	}
}
