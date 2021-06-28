package cli

import (
	types2 "bitbucket.org/decimalteam/go-node/x/coin/types"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"strings"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"

	cliUtils "bitbucket.org/decimalteam/go-node/x/coin/client/utils"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func GetCmdSendCoin(cdc *codec.LegacyAmino) *cobra.Command {
	return &cobra.Command{
		Use:   "send [coin] [amount] [receiver]",
		Short: "Send coin",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd).WithLegacyAmino(cdc)

			coin := args[0]
			amount, _ := sdk.NewIntFromString(args[1])
			receiver, err := sdk.AccAddressFromBech32(args[2])
			if err != nil {
				return err
			}
			msg := types2.NewMsgSendCoin(clientCtx.GetFromAddress(), sdk.NewCoin(coin, amount), receiver)
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			// Check if coin exists
			//existsCoin, _ := cliUtils.ExistsCoin(clientCtx, coin)
			//if !existsCoin {
			//	return types.ErrCoinDoesNotExist(coin)
			//}

			// Check if enough balance

			balance, err := cliUtils.GetAccountCoins(clientCtx, clientCtx.GetFromAddress())

			if err != nil {
				return err
			}
			if balance.AmountOf(strings.ToLower(coin)).LT(amount) {
				return types2.ErrInsufficientFunds(amount.String(), balance.AmountOf(strings.ToLower(coin)).String())
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}
}
