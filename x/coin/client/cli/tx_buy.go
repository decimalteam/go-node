package cli

import (
	"bitbucket.org/decimalteam/go-node/utils/formulas"
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

func GetCmdBuyCoin(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "buy [coinToBuy] [coinToSell] [amountToSell]",
		Short: "Buy coin",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			var coinToBuySymbol = args[0]

			var coinToSellSymbol = args[1]
			var amountToSell, _ = sdk.NewIntFromString(args[2])

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
			// TODO: Validate limits and check if sufficient balance (formulas)

			valueSell := formulas.CalculatePurchaseReturn(coinToBuy.Volume, coinToBuy.Reserve, coinToBuy.ConstantReserveRatio, amountToSell)

			// Do basic validating
			msg := types.NewMsgBuyCoin(cliCtx.GetFromAddress(), coinToBuySymbol, coinToSellSymbol, valueSell, amountToSell)
			err := msg.ValidateBasic()
			if err != nil {
				return err
			}

			// Get account balance
			acc, _ := cliUtils.GetAccount(cliCtx, cliCtx.GetFromAddress())
			balance := acc.GetCoins()
			if balance.AmountOf(strings.ToLower(coinToSellSymbol)).LT(amountToSell) {
				return sdk.NewError(types.DefaultCodespace, types.InsufficientCoinToSell, "Not enough coin to sell")
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}
