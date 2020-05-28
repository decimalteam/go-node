package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"

	"bitbucket.org/decimalteam/go-node/x/coin/internal/types"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd(cdc *codec.Codec) *cobra.Command {
	coinTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%s transactions subcommands", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	coinTxCmd.AddCommand(setFlagFeeCoin(flags.PostCommands(
		GetCmdCreateCoin(cdc),
		GetCmdBuyCoin(cdc),
		GetCmdSellCoin(cdc),
		GetCmdSendCoin(cdc),
		GetCmdSellAllCoin(cdc),
		GetCmdIssueCheck(cdc),
		GetCmdRedeemCheck(cdc),
	))...)

	return coinTxCmd
}

func setFlagFeeCoin(cmds []*cobra.Command) []*cobra.Command {
	for _, cmd := range cmds {
		cmd.Flags().String("fee-coin", "tDEL", "Coin for paying fee")
	}
	return cmds
}
