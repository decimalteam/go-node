package cli

import (
	cliUtils "bitbucket.org/decimalteam/go-node/x/coin/client/utils"
	types2 "bitbucket.org/decimalteam/go-node/x/coin/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
)

/*func GetCmdUpdateCoin(cdc *codec.LegacyAmino) *cobra.Command {
	return &cobra.Command{
		Use:   "update [symbol] [limitVolume] [identity]",
		Short: "Update custom coin",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.Context{}.WithLegacyAmino(cdc)

			// Parsing parameters to variables
			var symbol = args[0]
			var limitVolume, ok = sdk.NewIntFromString(args[1])
			if !ok {
				return types2.ErrInvalidLimitVolume
			}
			var identity = args[2]

			msg := types2.NewMsgUpdateCoin(clientCtx.GetFromAddress(), symbol, limitVolume, identity)
			// Check if coin does not exist yet
			coinExists, err := cliUtils.ExistsCoin(clientCtx, symbol)
			if err != nil {
				return err
			}
			if !coinExists {
				return types2.ErrCoinDoesNotExist(symbol)
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}
}*/

func GetCmdUpdateCoin() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update [symbol] [limitVolume] [identity] [from]",
		Short: "Update custom coin",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.Flags().Set(flags.FlagFrom, args[3])

			clientCtx, err := client.GetClientTxContext(cmd)
			var symbol = args[0]
			var limitVolume, ok = sdk.NewIntFromString(args[1])
			if !ok {
				return types2.ErrInvalidLimitVolume
			}
			var identity = args[2]

			msg := types2.NewMsgUpdateCoin(clientCtx.GetFromAddress(), symbol, limitVolume, identity)
			// Check if coin does not exist yet
			coinExists, err := cliUtils.ExistsCoin(clientCtx, symbol)
			if err != nil {
				return err
			}
			if !coinExists {
				return types2.ErrCoinDoesNotExist(symbol)
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}