package cli

import (
	types2 "bitbucket.org/decimalteam/go-node/x/coin/types"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	cliUtils "bitbucket.org/decimalteam/go-node/x/coin/client/utils"
)

func GetCmdCreateCoin(cdc *codec.LegacyAmino) *cobra.Command {
	return &cobra.Command{
		Use:   "create [title] [symbol] [crr] [initReserve] [initVolume] [limitVolume] [identity] [from]",
		Short: "Creates new coin",
		Args:  cobra.ExactArgs(8),
		RunE: func(cmd *cobra.Command, args []string) error {

			clientCtx := client.GetClientContextFromCmd(cmd)
			sender, err := sdk.AccAddressFromBech32(args[7])
			if err != nil {
				return err
			}
			// Parsing parameters to variables
			title := args[0]
			symbol := args[1]
			crr, err := strconv.ParseUint(args[2], 10, 8)
			// If error when convert crr
			if err != nil {
				return types2.ErrInvalidCRR()
			}
			initReserve, _ := sdk.NewIntFromString(args[3])
			initVolume, _ := sdk.NewIntFromString(args[4])
			limitVolume, _ := sdk.NewIntFromString(args[5])
			identity := args[6]

			msg := types2.NewMsgCreateCoin(sender, title, symbol, uint(crr), initVolume, initReserve, limitVolume, identity)
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}
			acc, err := cliUtils.GetAccount(clientCtx, sender)
			if err != nil {
				return err
			}
			balance, _ := cliUtils.GetAccountCoins(clientCtx, acc.GetAddress())
			if balance.AmountOf(cliUtils.GetBaseCoin()).LT(initReserve) {
				return types2.ErrInsufficientCoinReserve()
			}
			// Check if coin does not exist yet
			coinExists, _ := cliUtils.ExistsCoin(clientCtx, symbol)
			if coinExists {
				return types2.ErrCoinAlreadyExist(symbol)
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}
}
