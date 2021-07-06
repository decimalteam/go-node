package cli

import (
	"strconv"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"

	cliUtils "bitbucket.org/decimalteam/go-node/x/coin/client/utils"
	"bitbucket.org/decimalteam/go-node/x/coin/internal/types"
)

func GetCmdCreateCoin(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "create [title] [symbol] [crr] [initReserve] [initVolume] [limitVolume] [identity]",
		Short: "Creates new coin",
		Args:  cobra.ExactArgs(7),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := auth.NewTxBuilderFromCLI(cliCtx.Input).WithTxEncoder(utils.GetTxEncoder(cdc))
			// Parsing parameters to variables
			var title = args[0]
			var symbol = args[1]
			var crr, err = strconv.ParseUint(args[2], 10, 8)
			// If error when convert crr
			if err != nil {
				return types.ErrInvalidCRR(args[2])
			}
			var initReserve, _ = sdk.NewIntFromString(args[3])
			var initVolume, _ = sdk.NewIntFromString(args[4])
			var limitVolume, _ = sdk.NewIntFromString(args[5])
			var identity = args[6]

			msg := types.NewMsgCreateCoin(cliCtx.GetFromAddress(), title, symbol, uint(crr), initVolume, initReserve, limitVolume, identity)
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}
			// Check if coin does not exist yet
			coinExists, _ := cliUtils.ExistsCoin(cliCtx, symbol)
			if coinExists {
				return types.ErrCoinAlreadyExist(symbol)
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}
