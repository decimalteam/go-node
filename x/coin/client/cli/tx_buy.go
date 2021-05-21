package cli

import (
	types2 "bitbucket.org/decimalteam/go-node/x/coin/types"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"

	cliUtils "bitbucket.org/decimalteam/go-node/x/coin/client/utils"
)

func GetCmdBuyCoin(cdc *codec.LegacyAmino) *cobra.Command {
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
				return types2.ErrCoinDoesNotExist(coinToBuySymbol)
			}
			// Check if coin to sell exists
			coinToSell, _ := cliUtils.GetCoin(cliCtx, coinToSellSymbol)
			if coinToSell.Symbol != coinToSellSymbol {
				return types2.ErrCoinDoesNotExist(coinToSellSymbol)
			}
			// TODO: Calculate amounts and check limits
			// Do basic validating
			msg := types2.NewMsgBuyCoin(cliCtx.GetFromAddress(), sdk.NewCoin(coinToBuySymbol, amountToBuy), sdk.NewCoin(coinToSellSymbol, maxAmountToSell))
			err := msg.ValidateBasic()
			if err != nil {
				return err
			}

			// TODO: Check account balance
			// acc, _ := cliUtils.GetAccount(cliCtx, cliCtx.GetFromAddress())
			// balance := acc.GetCoins()
			// if balance.AmountOf(strings.ToLower(coinToSellSymbol)).LT(maxAmountToSell) {
			// 	return sdkerrors.New(types.DefaultCodespace, types.InsufficientCoinToSell, "Not enough coin to sell")
			// }

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}
