package cli

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"

	"bitbucket.org/decimalteam/go-node/x/multisig/internal/types"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd(cdc *codec.LegacyAmino) *cobra.Command {
	multisigTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%s transactions subcommands", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	multisigTxCmd.AddCommand(flags.PostCommands(
		getCmdCreateWallet(cdc),
		getCmdCreateTransaction(cdc),
		getCmdSignTransaction(cdc),
	)...)

	return multisigTxCmd
}

// getCmdCreateWallet is the CLI command for sending a CreateWallet transaction.
func getCmdCreateWallet(cdc *codec.LegacyAmino) *cobra.Command {
	return &cobra.Command{
		Use:   "create-wallet [owners] [weights] [threshold]",
		Short: "create a new multi-signature wallet",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := auth.NewTxBuilderFromCLI(cliCtx.Input).WithTxEncoder(utils.GetTxEncoder(cdc))

			ownersStrings := strings.Split(args[0], ",")
			weightsStrings := strings.Split(args[1], ",")
			thresholdString := args[2]

			owners := make([]sdk.AccAddress, len(ownersStrings))
			for i, address := range ownersStrings {
				owners[i], err = sdk.AccAddressFromBech32(address)
				if err != nil {
					return err
				}
			}

			weights := make([]uint, len(weightsStrings))
			for i, weightString := range weightsStrings {
				weight, err := strconv.ParseUint(weightString, 10, 64)
				if err != nil {
					return err
				}
				weights[i] = uint(weight)
			}

			threshold, err := strconv.ParseUint(thresholdString, 10, 64)
			if err != nil {
				return err
			}

			msg := types.NewMsgCreateWallet(cliCtx.GetFromAddress(), owners, weights, uint(threshold))
			if err != nil {
				return err
			}

			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

// getCmdCreateTransaction is the CLI command for sending a CreateTransaction transaction.
func getCmdCreateTransaction(cdc *codec.LegacyAmino) *cobra.Command {
	return &cobra.Command{
		Use:   "create-transaction [wallet] [receiver] [coins]",
		Short: "create a new multi-signature transaction",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := auth.NewTxBuilderFromCLI(cliCtx.Input).WithTxEncoder(utils.GetTxEncoder(cdc))

			walletString := args[0]
			receiverString := args[1]
			coinsString := args[2]

			wallet, err := sdk.AccAddressFromBech32(walletString)
			if err != nil {
				return err
			}

			receiver, err := sdk.AccAddressFromBech32(receiverString)
			if err != nil {
				return err
			}

			coins, err := sdk.ParseCoins(coinsString)
			if err != nil {
				return err
			}

			// TODO: Check coins exist?
			// for _, c := range coins {
			// 	coin, err := cliUtils.GetCoin(cliCtx, c.Denom)
			// 	if err != nil {
			// 		return err
			// 	}
			// 	if coin.Symbol != c.Denom {
			// 		return sdkerrors.New(types.DefaultCodespace, types.InvalidCoinToSend, fmt.Sprintf("Coin to send with symbol %s does not exist", c.Denom))
			// 	}
			// }

			msg := types.NewMsgCreateTransaction(cliCtx.GetFromAddress(), wallet, receiver, coins)
			if err != nil {
				return err
			}
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			fmt.Printf("%+v\n", msg)

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

// getCmdSignTransaction is the CLI command for saving a transaction signature.
func getCmdSignTransaction(cdc *codec.LegacyAmino) *cobra.Command {
	return &cobra.Command{
		Use:   "sign-transaction [tx-id]",
		Short: "Save a signature generated for a specific transaction",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := auth.NewTxBuilderFromCLI(cliCtx.Input).WithTxEncoder(utils.GetTxEncoder(cdc))

			txID := args[0]

			msg := types.NewMsgSignTransaction(cliCtx.GetFromAddress(), txID)
			if err != nil {
				return err
			}
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}
