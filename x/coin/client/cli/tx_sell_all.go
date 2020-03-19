package cli

import (
	cliUtils "bitbucket.org/decimalteam/go-node/x/coin/client/utils"
	"bitbucket.org/decimalteam/go-node/x/coin/internal/types"
	"fmt"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/spf13/cobra"
	"strings"
)

func GetCmdSellAllCoin(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "sell_all [coinToSell] [coinToBuy] [minAmountToBuy]",
		Short: "Sell all coin",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			coinToSellSymbol := args[0]
			coinToBuySymbol := args[1]
			minAmountToBuy, _ := sdk.NewIntFromString(args[2])

			// Check if coin to buy exists
			coinToBuy, _ := cliUtils.GetCoin(cliCtx, coinToBuySymbol)
			if coinToBuy.Symbol != coinToBuySymbol {
				return sdk.NewError(types.DefaultCodespace, types.CoinToBuyNotExists, fmt.Sprintf("Coin to buy with symbol %s does not exist", coinToBuySymbol))
			}
			// Check if coin to sell exists
			coinToSell, _ := cliUtils.GetCoin(cliCtx, coinToSellSymbol)
			if coinToSell.Symbol != coinToSellSymbol {
				return sdk.NewError(types.DefaultCodespace, types.CoinToSellNotExists, fmt.Sprintf("Coin to sell with symbol %s does not exist", coinToSellSymbol))
			}

			acc, _ := cliUtils.GetAccount(cliCtx, cliCtx.GetFromAddress())
			balance := acc.GetCoins()
			amountToSell := balance.AmountOf(strings.ToLower(coinToSellSymbol))

			_, _, err := cliUtils.SellCoinCalculateAmounts(coinToBuy, coinToSell, minAmountToBuy, amountToSell)
			if err != nil {
				return err
			}
			// Do basic validating
			msg := types.NewMsgSellAllCoin(cliCtx.GetFromAddress(), coinToBuySymbol, coinToSellSymbol, minAmountToBuy)
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}
