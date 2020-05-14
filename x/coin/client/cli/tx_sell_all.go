package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	decsdk "bitbucket.org/decimalteam/go-node/utils/types"
	"bitbucket.org/decimalteam/go-node/x/auth"
	"bitbucket.org/decimalteam/go-node/x/auth/client/utils"
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
				return sdkerrors.New(types.DefaultCodespace, types.CoinToBuyNotExists, fmt.Sprintf("Coin to buy with symbol %s does not exist", coinToBuySymbol))
			}
			// Check if coin to sell exists
			coinToSell, _ := cliUtils.GetCoin(cliCtx, coinToSellSymbol)
			if coinToSell.Symbol != coinToSellSymbol {
				return sdkerrors.New(types.DefaultCodespace, types.CoinToSellNotExists, fmt.Sprintf("Coin to sell with symbol %s does not exist", coinToSellSymbol))
			}

			// Get account balance
			acc, _ := cliUtils.GetAccount(cliCtx, cliCtx.GetFromAddress())
			balance := acc.GetCoins()
			if balance.AmountOf(strings.ToLower(coinToSellSymbol)).Sign() <= 0 {
				return sdkerrors.New(types.DefaultCodespace, types.InsufficientCoinToSell, "No coins to sell")
			}

			// TODO: Calculate amounts and check limits
			// Do basic validating
			msg := types.NewMsgSellAllCoin(decsdk.AccAddress(cliCtx.GetFromAddress()), coinToBuySymbol, coinToSellSymbol, minAmountToBuy)
			validationErr := msg.ValidateBasic()
			if validationErr != nil {
				return validationErr
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}
