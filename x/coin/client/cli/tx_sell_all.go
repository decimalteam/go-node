

package cli

import (
	types2 "bitbucket.org/decimalteam/go-node/x/coin/types"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"strings"

	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"

	cliUtils "bitbucket.org/decimalteam/go-node/x/coin/client/utils"
)

func GetCmdSellAllCoin() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sell_all [coinToSell] [coinToBuy] [minAmountToBuy] [from]",
		Short: "Sell all coin",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.Flags().Set(flags.FlagFrom, args[3])

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			coinToSellSymbol := args[0]
			coinToBuySymbol := args[1]
			minAmountToBuy, _ := sdk.NewIntFromString(args[2])

			// Check if coin to buy exists
			coinToBuy, _ := cliUtils.GetCoin(clientCtx, coinToBuySymbol)
			if coinToBuy.Symbol != coinToBuySymbol {
				return types2.ErrCoinDoesNotExist(coinToBuySymbol)
			}
			// Check if coin to sell exists
			coinToSell, _ := cliUtils.GetCoin(clientCtx, coinToSellSymbol)
			if coinToSell.Symbol != coinToSellSymbol {
				return types2.ErrCoinDoesNotExist(coinToSellSymbol)
			}

			// Get account balance
			balance, err := cliUtils.GetAccountCoins(clientCtx, clientCtx.GetFromAddress())

			if err != nil {
				return err
			}

			if balance.AmountOf(strings.ToLower(coinToSellSymbol)).Sign() <= 0 {
				return types2.ErrInsufficientFundsToSellAll()
			}

			// TODO: Calculate amounts and check limits
			// Do basic validating
			msg := types2.NewMsgSellAllCoin(clientCtx.GetFromAddress(), sdk.NewCoin(coinToSellSymbol, sdk.NewInt(0)), sdk.NewCoin(coinToBuySymbol, minAmountToBuy))
			validationErr := msg.ValidateBasic()
			if validationErr != nil {
				return validationErr
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}


