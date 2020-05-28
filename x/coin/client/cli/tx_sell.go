package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	cliUtils "bitbucket.org/decimalteam/go-node/x/coin/client/utils"
	"bitbucket.org/decimalteam/go-node/x/coin/internal/types"
	"bitbucket.org/decimalteam/go-node/x/validator/client/utils"
	vtypes "bitbucket.org/decimalteam/go-node/x/validator/utils"
)

func GetCmdSellCoin(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "sell [coinToSell] [amountToSell] [coinToBuy] [minAmountToBuy]",
		Short: "Sell coin",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := vtypes.NewTxBuilderFromCLI(cliCtx.Input).WithTxEncoder(utils.GetTxEncoder(cdc))

			var coinToSellSymbol = args[0]
			var amountToSell, _ = sdk.NewIntFromString(args[1])

			var coinToBuySymbol = args[2]
			var amountToBuy, _ = sdk.NewIntFromString(args[3])

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
			// TODO: Calculate amounts and check limits
			// Do basic validating
			msg := types.NewMsgSellCoin(cliCtx.GetFromAddress(), coinToBuySymbol, coinToSellSymbol, amountToSell, amountToBuy)
			validationErr := msg.ValidateBasic()
			if validationErr != nil {
				return validationErr
			}

			// Get account balance
			acc, _ := cliUtils.GetAccount(cliCtx, cliCtx.GetFromAddress())
			balance := acc.GetCoins()
			if balance.AmountOf(strings.ToLower(coinToSellSymbol)).LT(amountToSell) {
				return sdkerrors.New(types.DefaultCodespace, types.InsufficientCoinToSell, "Not enough coin to sell")
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}
