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

func GetCmdSendCoin(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "send [coin] [amount] [receiver]",
		Short: "Send coin",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			coin := args[0]
			amount, _ := sdk.NewIntFromString(args[1])
			receiver, err := sdk.AccAddressFromBech32(args[2])
			print(err)
			msg := types.NewMsgSendCoin(cliCtx.GetFromAddress(), coin, amount, receiver)
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			// Check if coin exists
			existsCoin, _ := cliUtils.ExistsCoin(cliCtx, coin)
			print(err)
			if !existsCoin {
				return sdk.NewError(types.DefaultCodespace, types.CoinToBuyNotExists, fmt.Sprintf("Coin to sent with symbol %s does not exist", coin))
			}

			// Check if enough balance
			acc, _ := cliUtils.GetAccount(cliCtx, cliCtx.GetFromAddress())
			balance := acc.GetCoins()
			if balance.AmountOf(strings.ToLower(coin)).LT(amount) {
				return sdk.NewError(types.DefaultCodespace, types.InsufficientCoinToSell, "Not enough coin to send")
			}
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}
