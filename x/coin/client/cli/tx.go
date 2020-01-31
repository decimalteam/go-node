package cli

import (
	"bitbucket.org/decimalteam/go-node/x/coin/internal/types"
	utils2 "bitbucket.org/decimalteam/go-node/x/coin/utils"
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/spf13/cobra"
	"strconv"
	"strings"
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
		GetCmdBuyCoin(cdc),
	)...)

	return coinTxCmd
}

func GetCmdCreateCoin(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "create [title] [symbol] [crr] [initReserve] [initVolume] [limitVolume]",
		Short: "Creates new coin",
		Args:  cobra.ExactArgs(6),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			// Parsing parameters to variables
			var title = args[0]
			var symbol = args[1]
			var crr, err = strconv.ParseUint(args[2], 10, 8)
			// If error when convert crr
			if err != nil {
				return sdk.NewError(types.DefaultCodespace, types.DecodeError, "Failed to convert CRR to uint")
			}
			var initReserve, _ = sdk.NewIntFromString(args[3])
			var initVolume, _ = sdk.NewIntFromString(args[4])
			var limitVolume, _ = sdk.NewIntFromString(args[5])

			msg := types.NewMsgCreateCoin(title, uint(crr), symbol, initVolume, initReserve, limitVolume, cliCtx.GetFromAddress())
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}
			// Check if coin does not exist yet
			coinExists, _ := utils2.ExistsCoin(cliCtx, symbol)
			if coinExists {
				return sdk.NewError(types.DefaultCodespace, types.CoinAlreadyExists, fmt.Sprintf("Coin with symbol %s already exists", symbol))
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

func GetCmdBuyCoin(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "buy [coinToBuy] [amountToBuy] [coinToSell] [maxAmountToSell]",
		Short: "Buy coin",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			var coinToBuySymbol = args[0]
			var amountToBuy, _ = sdk.NewIntFromString(args[1])

			var coinToSellSymbol = args[2]
			var maxAmountToSell, _ = sdk.NewIntFromString(args[3])

			// Do basic validating
			msg := types.NewMsgBuyCoin(cliCtx.GetFromAddress(), coinToBuySymbol, coinToSellSymbol, amountToBuy, maxAmountToSell)
			err := msg.ValidateBasic()
			if err != nil {
				return err
			}

			// Check if coin to buy exists
			coinToBuy, _ := utils2.GetCoin(cliCtx, coinToBuySymbol)
			if coinToBuy.Symbol != coinToBuySymbol {
				return sdk.NewError(types.DefaultCodespace, types.CoinToBuyNotExists, fmt.Sprintf("Coin to buy with symbol %s does not exist", coinToBuySymbol))
			}
			// Check if coin to sell exists
			coinToSell, _ := utils2.GetCoin(cliCtx, coinToSellSymbol)
			if coinToSell.Symbol != coinToSellSymbol {
				return sdk.NewError(types.DefaultCodespace, types.CoinToSellNotExists, fmt.Sprintf("Coin to sell with symbol %s does not exist", coinToSellSymbol))
			}
			// TODO: Validate limits and check if sufficient balance (formulas)
			// Get account balance
			acc, _ := utils2.GetAccount(cliCtx, cliCtx.GetFromAddress())
			balance := acc.GetCoins()
			if balance.AmountOf(strings.ToLower(coinToSellSymbol)).LT(amountToBuy) {
				return sdk.NewError(types.DefaultCodespace, types.InsufficientCoinToSell, "Not enough coin to sell")
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}
