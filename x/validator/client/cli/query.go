package cli

import (
	"fmt"
	"github.com/spf13/cobra"

	"bitbucket.org/decimalteam/go-node/x/validator/internal/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
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
		client.GetCommands(
		// TODO: Add query Cmds
		)...,
	)

	return validatorQueryCmd

}

// TODO: Add Query Commands
