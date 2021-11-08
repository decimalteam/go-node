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
		Use:   "create [from] [title] [symbol] [crr] [initReserve] [initVolume] [limitVolume] [identity]",
		Short: "Creates new coin",
		Args:  cobra.ExactArgs(8),
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.Flags().Set(flags.FlagFrom, args[0])

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			title := args[1]
			symbol := args[2]
			crr, err := strconv.ParseUint(args[3], 10, 8)
			// If error when convert crr
			if err != nil {
				return types.ErrInvalidCRR()
			}
			initReserve, _ := sdk.NewIntFromString(args[4])
			initVolume, _ := sdk.NewIntFromString(args[5])
			limitVolume, _ := sdk.NewIntFromString(args[6])
			identity := args[7]

			msg := types.NewMsgCreateCoin(clientCtx.GetFromAddress(), title, symbol, uint(crr), initVolume, initReserve, limitVolume, identity)
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}
			/*acc, err := cliUtils.GetAccount(clientCtx, clientCtx.GetFromAddress())
			if err != nil {
				return err
			}

			balance, _ := cliUtils.GetAccountCoins(clientCtx, acc.GetAddress())
			if balance.AmountOf(cliUtils.GetBaseCoin()).LT(initReserve) {
				return types.ErrInsufficientCoinReserve()
			}*/
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
