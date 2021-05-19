package cli

import (
	"strings"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"

	cliUtils "bitbucket.org/decimalteam/go-node/x/coin/client/utils"
	"bitbucket.org/decimalteam/go-node/x/coin/internal/types"
)

func GetCmdSellCoin(cdc *codec.LegacyAmino) *cobra.Command {
	return &cobra.Command{
		Use:   "sell [coinToSell] [amountToSell] [coinToBuy] [minAmountToBuy]",
		Short: "Sell coin",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := auth.NewTxBuilderFromCLI(cliCtx.Input).WithTxEncoder(utils.GetTxEncoder(cdc))

			var coinToSellSymbol = args[0]
			var amountToSell, _ = sdk.NewIntFromString(args[1])

			var coinToBuySymbol = args[2]
			var amountToBuy, _ = sdk.NewIntFromString(args[3])

			// Check if coin to buy exists
			coinToBuy, _ := cliUtils.GetCoin(cliCtx, coinToBuySymbol)
			if coinToBuy.Symbol != coinToBuySymbol {
				return types.ErrCoinDoesNotExist(coinToBuySymbol)
			}
			// Check if coin to sell exists
			coinToSell, _ := cliUtils.GetCoin(cliCtx, coinToSellSymbol)
			if coinToSell.Symbol != coinToSellSymbol {
				return types.ErrCoinDoesNotExist(coinToSellSymbol)
			}
			// TODO: Calculate amounts and check limits
			// Do basic validating
			msg := types.NewMsgSellCoin(cliCtx.GetFromAddress(), sdk.NewCoin(coinToSellSymbol, amountToSell), sdk.NewCoin(coinToBuySymbol, amountToBuy))
			validationErr := msg.ValidateBasic()
			if validationErr != nil {
				return validationErr
			}

			// Get account balance
			acc, _ := cliUtils.GetAccount(cliCtx, cliCtx.GetFromAddress())
			balance := acc.GetCoins()
			if balance.AmountOf(strings.ToLower(coinToSellSymbol)).LT(amountToSell) {
				return types.ErrInsufficientFunds(amountToSell.String(), balance.AmountOf(strings.ToLower(coinToSellSymbol)).String())
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}
