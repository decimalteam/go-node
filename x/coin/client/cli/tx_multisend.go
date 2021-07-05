package cli

import (
	cliUtils "bitbucket.org/decimalteam/go-node/x/coin/client/utils"
	types2 "bitbucket.org/decimalteam/go-node/x/coin/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
)

func GetCmdMultiSendCoin(cdc *codec.LegacyAmino) *cobra.Command {
	return &cobra.Command{
		Use:   "multisend [coin receiver] [coin receiver] ...",
		Short: "Multisend coin",
		Args:  cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd).WithLegacyAmino(cdc)

			sends := make([]types2.Send, len(args)/2)
			coins := make([]sdk.Coin, len(args)/2)

			for i, value := range args {
				if i%2 == 0 {
					coin, err := sdk.ParseCoinNormalized(value)
					if err != nil {
						return err
					}
					sends[i/2].Coin = coin
					coins[i/2] = coin
				} else {
					sends[i/2].Receiver = value
				}
			}

			msg := types2.NewMsgMultiSendCoin(clientCtx.GetFromAddress(), sends)

			// Check if enough balance
			balance, err := cliUtils.GetAccountCoins(clientCtx, clientCtx.GetFromAddress())
			if err != nil {
				return err
			}

			if !balance.IsAllGTE(coins) {
				var wantFunds string
				for _, send := range sends {
					wantFunds += send.Coin.String() + ", "
				}
				wantFunds = wantFunds[:len(wantFunds)-2]
				return types2.ErrInsufficientFunds(wantFunds, balance.String())
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}
}
