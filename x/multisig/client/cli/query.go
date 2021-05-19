package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"

	"bitbucket.org/decimalteam/go-node/x/multisig/internal/types"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd(queryRoute string, cdc *codec.LegacyAmino) *cobra.Command {
	// Group multisig queries under a subcommand
	multisigQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	multisigQueryCmd.AddCommand(
		flags.GetCommands(
			listWalletsCommand(queryRoute, cdc),
			getWalletCommand(queryRoute, cdc),
			listTransactionsCommand(queryRoute, cdc),
			getTransactionCommand(queryRoute, cdc),
		)...,
	)

	return multisigQueryCmd
}

// listWalletsCommand queries a list of wallets containing owner with specified address.
func listWalletsCommand(queryRoute string, cdc *codec.LegacyAmino) *cobra.Command {
	return &cobra.Command{
		Use:   "list-wallets [owner]",
		Short: "List all multi-signature wallets by owner address",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			owner := args[0]

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/listWallets/%s", queryRoute, owner), nil)
			if err != nil {
				fmt.Printf("could not list multi-signature wallets\n")
				return nil
			}

			var out types.QueryWallets
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}
}

// getWalletCommand queries a wallet with specified address.
func getWalletCommand(queryRoute string, cdc *codec.LegacyAmino) *cobra.Command {
	return &cobra.Command{
		Use:   "get-wallet [address]",
		Short: "Get multi-signature wallet by address",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			address := args[0]

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/getWallet/%s", queryRoute, address), nil)
			if err != nil {
				fmt.Printf("could not resolve multi-signature wallet %s\n", address)
				return nil
			}

			var out types.Wallet
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}
}

// listTransactionsCommand queries a list of transactions for the wallet with specified address.
func listTransactionsCommand(queryRoute string, cdc *codec.LegacyAmino) *cobra.Command {
	return &cobra.Command{
		Use:   "list-transactions [wallet]",
		Short: "List all multi-signature transactions by wallet address",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			wallet := args[0]

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/listTransactions/%s", queryRoute, wallet), nil)
			if err != nil {
				fmt.Printf("could not list multi-signature transactions\n")
				return nil
			}

			var out types.QueryTransactions
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}
}

// getTransactionCommand queries a transaction with specified ID.
func getTransactionCommand(queryRoute string, cdc *codec.LegacyAmino) *cobra.Command {
	return &cobra.Command{
		Use:   "get-transaction [tx_id]",
		Short: "Get multi-signature transaction by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			txID := args[0]

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/getTransaction/%s", queryRoute, txID), nil)
			if err != nil {
				fmt.Printf("could not resolve multi-signature transaction %s\n", txID)
				return nil
			}

			var out types.Transaction
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}
}
