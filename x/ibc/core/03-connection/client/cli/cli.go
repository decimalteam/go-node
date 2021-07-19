package cli

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/spf13/cobra"

	types2 "bitbucket.org/decimalteam/go-node/x/ibc/core/02-client/types"
)

// GetQueryCmd returns the query commands for IBC connections
func GetQueryCmd() *cobra.Command {
	queryCmd := &cobra.Command{
		Use:                        types2.SubModuleName,
		Short:                      "IBC connection query subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
	}

	queryCmd.AddCommand(
		GetCmdQueryConnections(),
		GetCmdQueryConnection(),
		GetCmdQueryClientConnections(),
	)

	return queryCmd
}

// NewTxCmd returns a CLI command handler for all x/ibc connection transaction commands.
func NewTxCmd() *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        types2.SubModuleName,
		Short:                      "IBC connection transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	txCmd.AddCommand(
		NewConnectionOpenInitCmd(),
		NewConnectionOpenTryCmd(),
		NewConnectionOpenAckCmd(),
		NewConnectionOpenConfirmCmd(),
	)

	return txCmd
}
