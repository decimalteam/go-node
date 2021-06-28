package cli

import (
	"bitbucket.org/decimalteam/go-node/x/validator/exported"
	"fmt"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"strings"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"

	"bitbucket.org/decimalteam/go-node/x/validator/types"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd(queryRoute string, cdc *codec.LegacyAmino) *cobra.Command {
	// Group validator queries under a subcommand
	validatorQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the validator module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	validatorQueryCmd.AddCommand(
		GetCmdQueryDelegation(queryRoute, cdc),
		GetCmdQueryDelegations(queryRoute, cdc),
		GetCmdQueryUnbondingDelegation(queryRoute, cdc),
		GetCmdQueryUnbondingDelegations(queryRoute, cdc),
		GetCmdQueryValidator(queryRoute, cdc),
		GetCmdQueryValidators(queryRoute, cdc),
		GetCmdQueryValidatorDelegations(queryRoute, cdc),
		GetCmdQueryValidatorUnbondingDelegations(queryRoute, cdc),
		GetCmdQueryParams(queryRoute, cdc),
		GetCmdQueryPool(queryRoute, cdc))

	return validatorQueryCmd

}

// GetCmdQueryValidator implements the validator query command.
func GetCmdQueryValidator(storeName string, cdc *codec.LegacyAmino) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validator [validator-addr]",
		Short: "Query a validator",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query details about an individual validator.

Example:
$ %s query validator validator dxvaloper1gghjut3ccd8ay0zduzj64hwre2fxs9ldmqhffj
`,
				version.AppName,
			),
		),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			addr, err := sdk.ValAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			res, _, err := clientCtx.QueryStore(types.GetValidatorKey(addr), storeName)
			if err != nil {
				return err
			}

			if len(res) == 0 {
				return fmt.Errorf("No validator found with address %s ", addr)
			}

			validator, err := types.UnmarshalValidator(cdc, res)
			if err != nil {
				return err
			}
			return clientCtx.PrintObjectLegacy(validator)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetCmdQueryValidators implements the query all validators command.
func GetCmdQueryValidators(storeName string, cdc *codec.LegacyAmino) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validators",
		Short: "Query for all validators",
		Args:  cobra.NoArgs,
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query details about all validators on a network.

Example:
$ %s query validator validators
`,
				version.AppName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			resKVs, _, err := clientCtx.QueryStore([]byte{types.ValidatorsKey}, storeName)
			if err != nil {
				return err
			}

			var validators types.Validators
			for _, kv := range resKVs {
				validator, err := types.UnmarshalValidator(cdc, kv)
				if err != nil {
					return err
				}
				validators = append(validators, validator)
			}

			return clientCtx.PrintObjectLegacy(validators)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetCmdQueryValidatorUnbondingDelegations implements the query all unbonding delegatations from a validator command.
func GetCmdQueryValidatorUnbondingDelegations(queryRoute string, cdc *codec.LegacyAmino) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "unbonding-delegations-from [validator-addr]",
		Short: "Query all unbonding delegatations from a validator",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query delegations that are unbonding _from_ a validator.

Example:
$ %s query validator unbonding-delegations-from cosmosvaloper1gghjut3ccd8ay0zduzj64hwre2fxs9ldmqhffj
`,
				version.AppName,
			),
		),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			valAddr, err := sdk.ValAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			bz, err := cdc.MarshalJSON(types.NewQueryValidatorParams(valAddr))
			if err != nil {
				return err
			}

			route := fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryValidatorUnbondingDelegations)
			res, _, err := clientCtx.QueryWithData(route, bz)
			if err != nil {
				return err
			}

			var ubds types.UnbondingDelegations
			cdc.MustUnmarshalJSON(res, &ubds)
			return clientCtx.PrintObjectLegacy(ubds)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetCmdQueryDelegation the query delegation command.
func GetCmdQueryDelegation(queryRoute string, cdc *codec.LegacyAmino) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delegation [delegator-addr] [validator-addr] [coin]",
		Short: "Query a delegation based on address and validator address",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query delegations for an individual delegator on an individual validator.

Example:
$ %s query validator delegation cosmos1gghjut3ccd8ay0zduzj64hwre2fxs9ld75ru9p cosmosvaloper1gghjut3ccd8ay0zduzj64hwre2fxs9ldmqhffj
`,
				version.AppName,
			),
		),
		Args: cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			delAddr, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			valAddr, err := sdk.ValAddressFromBech32(args[1])
			if err != nil {
				return err
			}

			coin := args[2]

			bz, err := cdc.MarshalJSON(types.NewQueryBondsParams(delAddr, valAddr, coin))
			if err != nil {
				return err
			}

			route := fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryDelegation)
			res, _, err := clientCtx.QueryWithData(route, bz)
			if err != nil {
				return err
			}

			var resp exported.DelegationI
			if err := cdc.UnmarshalJSON(res, &resp); err != nil {
				return err
			}

			return clientCtx.PrintObjectLegacy(resp)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetCmdQueryDelegations implements the command to query all the delegations
// made from one delegator.
func GetCmdQueryDelegations(queryRoute string, cdc *codec.LegacyAmino) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delegations [delegator-addr]",
		Short: "Query all delegations made by one delegator",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query delegations for an individual delegator on all validators.

Example:
$ %s query validator delegations cosmos1gghjut3ccd8ay0zduzj64hwre2fxs9ld75ru9p
`,
				version.AppName,
			),
		),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			delAddr, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			bz, err := cdc.MarshalJSON(types.NewQueryDelegatorParams(delAddr))
			if err != nil {
				return err
			}

			route := fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryDelegatorDelegations)
			res, _, err := clientCtx.QueryWithData(route, bz)
			if err != nil {
				return err
			}

			var resp []exported.DelegationI
			if err := cdc.UnmarshalJSON(res, &resp); err != nil {
				return err
			}

			return clientCtx.PrintObjectLegacy(resp)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetCmdQueryValidatorDelegations implements the command to query all the
// delegations to a specific validator.
func GetCmdQueryValidatorDelegations(queryRoute string, cdc *codec.LegacyAmino) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delegations-to [validator-addr]",
		Short: "Query all delegations made to one validator",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query delegations on an individual validator.

Example:
$ %s query validator delegations-to cosmosvaloper1gghjut3ccd8ay0zduzj64hwre2fxs9ldmqhffj
`,
				version.AppName,
			),
		),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			valAddr, err := sdk.ValAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			bz, err := cdc.MarshalJSON(types.NewQueryValidatorParams(valAddr))
			if err != nil {
				return err
			}

			route := fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryValidatorDelegations)
			res, _, err := clientCtx.QueryWithData(route, bz)
			if err != nil {
				return err
			}

			var resp types.DelegationResponse
			if err := cdc.UnmarshalJSON(res, &resp); err != nil {
				return err
			}

			return clientCtx.PrintObjectLegacy(resp)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetCmdQueryUnbondingDelegation implements the command to query a single
// unbonding-delegation record.
func GetCmdQueryUnbondingDelegation(queryRoute string, cdc *codec.LegacyAmino) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "unbonding-delegation [delegator-addr] [validator-addr]",
		Short: "Query an unbonding-delegation record based on delegator and validator address",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query unbonding delegations for an individual delegator on an individual validator.

Example:
$ %s query validator unbonding-delegation cosmos1gghjut3ccd8ay0zduzj64hwre2fxs9ld75ru9p cosmosvaloper1gghjut3ccd8ay0zduzj64hwre2fxs9ldmqhffj
`,
				version.AppName,
			),
		),
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			valAddr, err := sdk.ValAddressFromBech32(args[1])
			if err != nil {
				return err
			}

			delAddr, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			bz, err := cdc.MarshalJSON(types.NewQueryBondsParams(delAddr, valAddr, ""))
			if err != nil {
				return err
			}

			route := fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryUnbondingDelegation)
			res, _, err := clientCtx.QueryWithData(route, bz)
			if err != nil {
				return err
			}

			return clientCtx.PrintObjectLegacy(types.MustUnmarshalUBD(cdc, res))
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetCmdQueryUnbondingDelegations implements the command to query all the
// unbonding-delegation records for a delegator.
func GetCmdQueryUnbondingDelegations(queryRoute string, cdc *codec.LegacyAmino) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "unbonding-delegations [delegator-addr]",
		Short: "Query all unbonding-delegations records for one delegator",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query unbonding delegations for an individual delegator.

Example:
$ %s query validator unbonding-delegation cosmos1gghjut3ccd8ay0zduzj64hwre2fxs9ld75ru9p
`,
				version.AppName,
			),
		),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			delegatorAddr, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			bz, err := cdc.MarshalJSON(types.NewQueryDelegatorParams(delegatorAddr))
			if err != nil {
				return err
			}

			route := fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryDelegatorUnbondingDelegations)
			res, _, err := clientCtx.QueryWithData(route, bz)
			if err != nil {
				return err
			}

			var ubds types.UnbondingDelegations
			if err = cdc.UnmarshalJSON(res, &ubds); err != nil {
				return err
			}

			return clientCtx.PrintObjectLegacy(ubds)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetCmdQueryPool implements the pool query command.
func GetCmdQueryPool(storeName string, cdc *codec.LegacyAmino) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pool",
		Args:  cobra.NoArgs,
		Short: "Query the current validator pool values",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query values for amounts stored in the validator pool.

Example:
$ %s query validator pool
`,
				version.AppName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			bz, _, err := clientCtx.QueryWithData(fmt.Sprintf("custom/%s/pool", storeName), nil)
			if err != nil {
				return err
			}

			var pool types.Pool
			if err := cdc.UnmarshalJSON(bz, &pool); err != nil {
				return err
			}

			return clientCtx.PrintObjectLegacy(pool)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetCmdQueryParams implements the params query command.
func GetCmdQueryParams(storeName string, cdc *codec.LegacyAmino) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "params",
		Args:  cobra.NoArgs,
		Short: "Query the current validator parameters information",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query values set as validator parameters.

Example:
$ %s query validator params
`,
				version.AppName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			route := fmt.Sprintf("custom/%s/%s", storeName, types.QueryParameters)
			bz, _, err := clientCtx.QueryWithData(route, nil)
			if err != nil {
				return err
			}

			var params types.Params
			cdc.MustUnmarshalJSON(bz, &params)
			return clientCtx.PrintObjectLegacy(params)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
