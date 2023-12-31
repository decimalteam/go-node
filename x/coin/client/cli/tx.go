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

	coinTxCmd.AddCommand(flags.PostCommands(
		GetCmdCreateCoin(cdc),
		GetCmdUpdateCoin(cdc),
		GetCmdBuyCoin(cdc),
		GetCmdSellCoin(cdc),
		GetCmdSendCoin(cdc),
		GetCmdMultiSendCoin(cdc),
		GetCmdSellAllCoin(cdc),
		GetCmdIssueCheck(cdc),
		GetCmdRedeemCheck(cdc),
	)...)

	return coinTxCmd
}
