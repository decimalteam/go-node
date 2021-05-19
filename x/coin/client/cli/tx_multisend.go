package cli

import (
	cliUtils "bitbucket.org/decimalteam/go-node/x/coin/client/utils"
	"bitbucket.org/decimalteam/go-node/x/coin/internal/types"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/spf13/cobra"
)

func GetCmdMultiSendCoin(cdc *codec.LegacyAmino) *cobra.Command {
	return &cobra.Command{
		Use:   "multisend [coin receiver] [coin receiver] ...",
		Short: "Multisend coin",
		Args:  cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := auth.NewTxBuilderFromCLI(cliCtx.Input).WithTxEncoder(utils.GetTxEncoder(cdc))

			sends := make([]types.Send, len(args)/2)
			coins := make([]sdk.Coin, len(args)/2)

			for i, value := range args {
				if i%2 == 0 {
					coin, err := sdk.ParseCoin(value)
					if err != nil {
						return err
					}
					sends[i/2].Coin = coin
					coins[i/2] = coin
				} else {
					receiver, err := sdk.AccAddressFromBech32(value)
					if err != nil {
						return err
					}
					sends[i/2].Receiver = receiver
				}
			}

			msg := types.NewMsgMultiSendCoin(cliCtx.GetFromAddress(), sends)

			// Check if enough balance
			acc, err := cliUtils.GetAccount(cliCtx, cliCtx.GetFromAddress())
			if err != nil {
				return err
			}
			balance := acc.GetCoins()
			if !balance.IsAllGTE(coins) {
				var wantFunds string
				for _, send := range sends {
					wantFunds += send.Coin.String() + ", "
				}
				wantFunds = wantFunds[:len(wantFunds)-2]
				return types.ErrInsufficientFunds(wantFunds, balance.String())
			}
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}
