package cli

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/spf13/cobra"
	"strconv"

	"bitbucket.org/decimalteam/go-node/x/coin/internal/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd(cdc *codec.Codec) *cobra.Command {
	coinTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%s transactions subcommands", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	coinTxCmd.AddCommand(client.PostCommands(
		GetCmdCreateCoin(cdc),
	)...)

	return coinTxCmd
}

func GetCmdCreateCoin(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "create [title] [symbol] [crr] [initReserve] [initAmount] [limitAmount]",
		Short: "Creates new coin",
		Args:  cobra.ExactArgs(6),
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Check if coin does not exist yet
			// At this moment transaction pass CheckTx,
			// But on DeliverTx it fails and become invalid
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			var title = args[0]
			var symbol = args[1]
			var crr, err = strconv.ParseUint(args[2], 10, 8)
			// If error when convert crr
			if err != nil {
				return sdk.NewError(types.DefaultCodespace, types.DecodeError, "Failed to convert CRR to uint")
			}
			var initReserve, _ = sdk.NewIntFromString(args[3])
			var initAmount, _ = sdk.NewIntFromString(args[4])
			var limitAmount, _ = sdk.NewIntFromString(args[5])

			msg := types.NewMsgCreateCoin(title, uint(crr), symbol, initAmount, initReserve, limitAmount, cliCtx.GetFromAddress())
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}