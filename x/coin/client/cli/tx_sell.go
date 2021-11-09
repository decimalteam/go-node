
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

func GetCmdSellCoin() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sell [coinToSell] [amountToSell] [coinToBuy] [minAmountToBuy] [from]",
		Short: "Sell coin",
		Args:  cobra.ExactArgs(5),
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.Flags().Set(flags.FlagFrom, args[4])

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			var coinToSellSymbol = args[0]
			var amountToSell, _ = sdk.NewIntFromString(args[1])

			var coinToBuySymbol = args[2]
			var amountToBuy, _ = sdk.NewIntFromString(args[3])

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
			// TODO: Calculate amounts and check limits
			// Do basic validating
			msg := types2.NewMsgSellCoin(clientCtx.GetFromAddress(), sdk.NewCoin(coinToSellSymbol, amountToSell), sdk.NewCoin(coinToBuySymbol, amountToBuy))
			validationErr := msg.ValidateBasic()
			if validationErr != nil {
				return validationErr
			}

			// Get account balance
			balance, err := cliUtils.GetAccountCoins(clientCtx, clientCtx.GetFromAddress())

			if err != nil {
				return err
			}

			if balance.AmountOf(strings.ToLower(coinToSellSymbol)).LT(amountToSell) {
				return types2.ErrInsufficientFunds(amountToSell.String(), balance.AmountOf(strings.ToLower(coinToSellSymbol)).String())
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}



/*package cli

import (
	types2 "bitbucket.org/decimalteam/go-node/x/coin/types"
	"strings"

	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	cliUtils "bitbucket.org/decimalteam/go-node/x/coin/client/utils"
)

func GetCmdSellCoin(cdc *codec.LegacyAmino) *cobra.Command {
	return &cobra.Command{
		Use:   "sell [coinToSell] [amountToSell] [coinToBuy] [minAmountToBuy]",
		Short: "Sell coin",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd).WithLegacyAmino(cdc)

			var coinToSellSymbol = args[0]
			var amountToSell, _ = sdk.NewIntFromString(args[1])

			var coinToBuySymbol = args[2]
			var amountToBuy, _ = sdk.NewIntFromString(args[3])

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
			// TODO: Calculate amounts and check limits
			// Do basic validating
			msg := types2.NewMsgSellCoin(clientCtx.GetFromAddress(), sdk.NewCoin(coinToSellSymbol, amountToSell), sdk.NewCoin(coinToBuySymbol, amountToBuy))
			validationErr := msg.ValidateBasic()
			if validationErr != nil {
				return validationErr
			}

			// Get account balance
			balance, err := cliUtils.GetAccountCoins(clientCtx, clientCtx.GetFromAddress())

			if err != nil {
				return err
			}

			if balance.AmountOf(strings.ToLower(coinToSellSymbol)).LT(amountToSell) {
				return types2.ErrInsufficientFunds(amountToSell.String(), balance.AmountOf(strings.ToLower(coinToSellSymbol)).String())
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}
}
*/