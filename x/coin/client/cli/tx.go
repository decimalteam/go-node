package cli

import (
	types2 "bitbucket.org/decimalteam/go-node/x/coin/types"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd(cdc *codec.LegacyAmino) *cobra.Command {
	coinTxCmd := &cobra.Command{
		Use:                        types2.ModuleName,
		Short:                      fmt.Sprintf("%s transactions subcommands", types2.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	coinTxCmd.AddCommand(
		GetCmdCreateCoin(),
		GetCmdUpdateCoin(),
		GetCmdBuyCoin(),
		GetCmdSellCoin(),
		GetCmdSendCoin(),
		GetCmdMultiSendCoin(),
		GetCmdSellAllCoin(),
		GetCmdIssueCheck(),
		GetCmdRedeemCheck(),
	)

	return coinTxCmd
}
