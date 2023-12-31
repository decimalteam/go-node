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

func GetCmdSellAllCoin(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "sell_all [coinToSell] [coinToBuy] [minAmountToBuy]",
		Short: "Sell all coin",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := auth.NewTxBuilderFromCLI(cliCtx.Input).WithTxEncoder(utils.GetTxEncoder(cdc))

			coinToSellSymbol := args[0]
			coinToBuySymbol := args[1]
			minAmountToBuy, _ := sdk.NewIntFromString(args[2])

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

			// Get account balance
			acc, err := cliUtils.GetAccount(cliCtx, cliCtx.GetFromAddress())
			if err != nil {
				return err
			}
			balance := acc.GetCoins()
			if balance.AmountOf(strings.ToLower(coinToSellSymbol)).Sign() <= 0 {
				return types.ErrInsufficientFundsToSellAll()
			}

			// TODO: Calculate amounts and check limits
			// Do basic validating
			msg := types.NewMsgSellAllCoin(cliCtx.GetFromAddress(), sdk.NewCoin(coinToSellSymbol, sdk.NewInt(0)), sdk.NewCoin(coinToBuySymbol, minAmountToBuy))
			validationErr := msg.ValidateBasic()
			if validationErr != nil {
				return validationErr
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}
