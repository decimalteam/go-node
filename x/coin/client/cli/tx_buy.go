package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"

	cliUtils "bitbucket.org/decimalteam/go-node/x/coin/client/utils"
	"bitbucket.org/decimalteam/go-node/x/coin/internal/types"
)

func GetCmdBuyCoin(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "buy [coinToBuy] [amountToBuy] [coinToSell] [maxAmountToSell]",
		Short: "Buy coin",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := auth.NewTxBuilderFromCLI(cliCtx.Input).WithTxEncoder(utils.GetTxEncoder(cdc))

			var coinToBuySymbol = args[0]
			var coinToSellSymbol = args[2]
			var amountToBuy, _ = sdk.NewIntFromString(args[1])
			var maxAmountToSell, _ = sdk.NewIntFromString(args[3])

			// Check if coin to buy exists
			coinToBuy, _ := cliUtils.GetCoin(cliCtx, coinToBuySymbol)
			if coinToBuy.Symbol != coinToBuySymbol {
				return sdkerrors.New(types.DefaultCodespace, types.CoinToBuyNotExists, fmt.Sprintf("Coin to buy with symbol %s does not exist", coinToBuySymbol))
			}
			// Check if coin to sell exists
			coinToSell, _ := cliUtils.GetCoin(cliCtx, coinToSellSymbol)
			if coinToSell.Symbol != coinToSellSymbol {
				return sdkerrors.New(types.DefaultCodespace, types.CoinToSellNotExists, fmt.Sprintf("Coin to sell with symbol %s does not exist", coinToSellSymbol))
			}
			// Calculate amounts and check limits
			amountBuy, amountSell, calcErr := cliUtils.BuyCoinCalculateAmounts(coinToBuy, coinToSell, amountToBuy, maxAmountToSell)
			if calcErr != nil {
				return calcErr
			}
			// Do basic validating
			msg := types.NewMsgBuyCoin(cliCtx.GetFromAddress(), coinToBuySymbol, coinToSellSymbol, amountBuy, amountSell)
			err := msg.ValidateBasic()
			if err != nil {
				return err
			}

			// Get account balance
			acc, _ := cliUtils.GetAccount(cliCtx, cliCtx.GetFromAddress())
			balance := acc.GetCoins()
			if balance.AmountOf(strings.ToLower(coinToSellSymbol)).LT(amountSell) {
				return sdkerrors.New(types.DefaultCodespace, types.InsufficientCoinToSell, "Not enough coin to sell")
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}
