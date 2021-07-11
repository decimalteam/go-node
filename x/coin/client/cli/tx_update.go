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

func GetCmdUpdateCoin(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "update [symbol] [limitVolume] [identity]",
		Short: "Update custom coin",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			txBldr := auth.NewTxBuilderFromCLI(cliCtx.Input).WithTxEncoder(utils.GetTxEncoder(cdc))
			// Parsing parameters to variables
			var symbol = args[0]
			var limitVolume, ok = sdk.NewIntFromString(args[1])
			if !ok {
				return types.ErrInvalidLimitVolume
			}
			var identity = args[2]

			msg := types.NewMsgUpdateCoin(cliCtx.GetFromAddress(), symbol, limitVolume, identity)
			// Check if coin does not exist yet
			coinExists, err := cliUtils.ExistsCoin(cliCtx, symbol)
			if err != nil {
				return err
			}
			if !coinExists {
				return types.ErrCoinDoesNotExist(symbol)
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}
