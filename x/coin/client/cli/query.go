package cli

import (
	types2 "bitbucket.org/decimalteam/go-node/x/coin/types"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
)

// GetQueryCmd returns the CLI query commands for this module.
func GetQueryCmd(queryRoute string, cdc *codec.LegacyAmino) *cobra.Command {
	// Group coin queries under a subcommand
	coinQueryCmd := &cobra.Command{
		Use:                        types2.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types2.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	coinQueryCmd.AddCommand(
		flags.GetCommands(
			listCoinsCommand(queryRoute, cdc),
			getCoinCommand(queryRoute, cdc),
		)...,
	)

	return coinQueryCmd

}

func listCoinsCommand(queryRoute string, cdc *codec.LegacyAmino) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all existing coins",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := client.Context{}.WithLegacyAmino(cdc)

			path := fmt.Sprintf("custom/%s/%s", queryRoute, types2.QueryListCoins)
			res, _, err := ctx.QueryWithData(path, nil)
			if err != nil {
				fmt.Printf("could not get coins\n%s\n", err.Error())
				return nil
			}

			var out types2.QueryResCoins
			cdc.MustUnmarshalJSON(res, &out)
			return ctx.PrintOutput(out)
		},
	}
}

func getCoinCommand(queryRoute string, cdc *codec.LegacyAmino) *cobra.Command {
	return &cobra.Command{
		Use:   "get [symbol]",
		Short: "Returns coin information by symbol",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := client.Context{}.WithLegacyAmino(cdc)
			symbol := args[0]

			path := fmt.Sprintf("custom/%s/%s/%s", queryRoute, types2.QueryGetCoin, symbol)
			res, _, err := ctx.QueryWithData(path, nil)
			if err != nil {
				fmt.Printf("could not resolve coin %s\n%s\n", symbol, err.Error())

				return nil
			}

			var out types2.Coin
			cdc.MustUnmarshalJSON(res, &out)
			return ctx.PrintOutput(out)
		},
	}
}
