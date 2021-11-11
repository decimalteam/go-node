package cli

import (
	cliUtils "bitbucket.org/decimalteam/go-node/x/coin/client/utils"
	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
	"strconv"

	"bitbucket.org/decimalteam/go-node/x/coin/types"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
)

func GetCmdCreateCoin() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create [title] [symbol] [crr] [initReserve] [initVolume] [limitVolume] [identity] [from]",
		Short: "Creates new coin",
		Args:  cobra.ExactArgs(8),
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.Flags().Set(flags.FlagFrom, args[7])

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			title := args[0]
			symbol := args[1]
			crr, err := strconv.ParseUint(args[2], 10, 8)
			// If error when convert crr
			if err != nil {
				return types.ErrInvalidCRR()
			}
			initReserve, _ := sdk.NewIntFromString(args[3])
			initVolume, _ := sdk.NewIntFromString(args[4])
			limitVolume, _ := sdk.NewIntFromString(args[5])
			identity := args[6]

			msg := types.NewMsgCreateCoin(clientCtx.GetFromAddress(), title, symbol, uint(crr), initVolume, initReserve, limitVolume, identity)
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}
			acc, err := cliUtils.GetAccount(clientCtx, clientCtx.GetFromAddress())
			if err != nil {
				return err
			}
			acc.GetAddress()

			balance, _ := cliUtils.GetAccountCoins(clientCtx, clientCtx.GetFromAddress())
			if balance.AmountOf(cliUtils.GetBaseCoin()).LT(initReserve) {
				return types.ErrInsufficientCoinReserve()
			}
			// Check if coin does not exist yet
			coinExists, _ := cliUtils.ExistsCoin(clientCtx, symbol)
			if coinExists {
				return types.ErrCoinAlreadyExist(symbol)
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
