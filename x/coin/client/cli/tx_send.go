package cli

import (
	"fmt"
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

func GetCmdSendCoin(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "send [coin] [amount] [receiver]",
		Short: "Send coin",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := auth.NewTxBuilderFromCLI(cliCtx.Input).WithTxEncoder(utils.GetTxEncoder(cdc))

			coin := args[0]
			amount, ok := sdk.NewIntFromString(args[1])
			if !ok {
				return fmt.Errorf("bitLen(int) > maxBitLen(=255)")
			}
			receiver, err := sdk.AccAddressFromBech32(args[2])
			if err != nil {
				return err
			}
			msg := types.NewMsgSendCoin(cliCtx.GetFromAddress(), sdk.NewCoin(coin, amount), receiver)
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			// Check if coin exists
			//existsCoin, _ := cliUtils.ExistsCoin(cliCtx, coin)
			//if !existsCoin {
			//	return types.ErrCoinDoesNotExist(coin)
			//}

			// Check if enough balance
			acc, err := cliUtils.GetAccount(cliCtx, cliCtx.GetFromAddress())
			if err != nil {
				return err
			}
			balance := acc.GetCoins()
			if balance.AmountOf(strings.ToLower(coin)).LT(amount) {
				return types.ErrInsufficientFunds(amount.String(), balance.AmountOf(strings.ToLower(coin)).String())
			}
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}
