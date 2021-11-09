package cli

import (
	cliUtils "bitbucket.org/decimalteam/go-node/x/coin/client/utils"
	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"

	types2"bitbucket.org/decimalteam/go-node/x/coin/types"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
)

func GetCmdBuyCoin() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "buy [coinToBuy] [amountToBuy] [coinToSell] [maxAmountToSell] [from]",
		Short: "Buy coin",
		Args:  cobra.ExactArgs(5),
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.Flags().Set(flags.FlagFrom, args[4])

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			var coinToBuySymbol = args[0]
			var coinToSellSymbol = args[2]
			var amountToBuy, _ = sdk.NewIntFromString(args[1])
			var maxAmountToSell, _ = sdk.NewIntFromString(args[3])
			// Check if coin to buy exists
			coinToBuy, _ := cliUtils.GetCoin(clientCtx, coinToBuySymbol)
			if coinToBuy.Symbol != coinToBuySymbol {
				return types2.ErrCoinDoesNotExist(coinToBuySymbol)
			}
			// Check if coin to sell exists
			coinToSell, _ := cliUtils.GetCoin(clientCtx, coinToSellSymbol)
			if coinToSell.Symbol != coinToSellSymbol {
				return types2.ErrCoinDoesNotExist(coinToSellSymbol)
			}
			// TODO: Calculate amounts and check limits
			// Do basic validating
			msg := types2.NewMsgBuyCoin(clientCtx.GetFromAddress(), sdk.NewCoin(coinToBuySymbol, amountToBuy), sdk.NewCoin(coinToSellSymbol, maxAmountToSell))
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			// TODO: Check account balance
			// acc, _ := cliUtils.GetAccount(clientCtx, clientCtx.GetFromAddress())
			// balance := acc.GetCoins()
			// if balance.AmountOf(strings.ToLower(coinToSellSymbol)).LT(maxAmountToSell) {
			// 	return sdkerrors.New(types.DefaultCodespace, types.InsufficientCoinToSell, "Not enough coin to sell")
			// }

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

/*package cli

import (
	types2 "bitbucket.org/decimalteam/go-node/x/coin/types"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"

	cliUtils "bitbucket.org/decimalteam/go-node/x/coin/client/utils"
)

func GetCmdBuyCoin(cdc *codec.LegacyAmino) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "buy [coinToBuy] [amountToBuy] [coinToSell] [maxAmountToSell]",
		Short: "Buy coin",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd).WithLegacyAmino(cdc)

			var coinToBuySymbol = args[0]
			var coinToSellSymbol = args[2]
			var amountToBuy, _ = sdk.NewIntFromString(args[1])
			var maxAmountToSell, _ = sdk.NewIntFromString(args[3])

			// Check if coin to buy exists
			coinToBuy, _ := cliUtils.GetCoin(clientCtx, coinToBuySymbol)
			if coinToBuy.Symbol != coinToBuySymbol {
				return types2.ErrCoinDoesNotExist(coinToBuySymbol)
			}
			// Check if coin to sell exists
			coinToSell, _ := cliUtils.GetCoin(clientCtx, coinToSellSymbol)
			if coinToSell.Symbol != coinToSellSymbol {
				return types2.ErrCoinDoesNotExist(coinToSellSymbol)
			}
			// TODO: Calculate amounts and check limits
			// Do basic validating
			msg := types2.NewMsgBuyCoin(clientCtx.GetFromAddress(), sdk.NewCoin(coinToBuySymbol, amountToBuy), sdk.NewCoin(coinToSellSymbol, maxAmountToSell))
			err := msg.ValidateBasic()
			if err != nil {
				return err
			}

			// TODO: Check account balance
			// acc, _ := cliUtils.GetAccount(clientCtx, clientCtx.GetFromAddress())
			// balance := acc.GetCoins()
			// if balance.AmountOf(strings.ToLower(coinToSellSymbol)).LT(maxAmountToSell) {
			// 	return sdkerrors.New(types.DefaultCodespace, types.InsufficientCoinToSell, "Not enough coin to sell")
			// }

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
*/