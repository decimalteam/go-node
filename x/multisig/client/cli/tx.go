package cli

import (
	types2 "bitbucket.org/decimalteam/go-node/x/multisig/types"
	"fmt"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd() *cobra.Command {
	multisigTxCmd := &cobra.Command{
		Use:                        types2.ModuleName,
		Short:                      fmt.Sprintf("%s transactions subcommands", types2.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	multisigTxCmd.AddCommand(
		getCmdCreateWallet(),
		getCmdCreateTransaction(),
		getCmdSignTransaction(),
	)

	return multisigTxCmd
}

// getCmdCreateWallet is the CLI command for sending a CreateWallet transaction.
func getCmdCreateWallet() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-wallet [owners] [weights] [threshold]",
		Short: "create a new multi-signature wallet",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error
			clientCtx := client.GetClientContextFromCmd(cmd)

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

			weights := make([]uint64, len(weightsStrings))
			for i, weightString := range weightsStrings {
				weight, err := strconv.ParseUint(weightString, 10, 64)
				if err != nil {
					return err
				}
				weights[i] = weight
			}

			threshold, err := strconv.ParseUint(thresholdString, 10, 64)
			if err != nil {
				return err
			}

			msg := types2.NewMsgCreateWallet(clientCtx.GetFromAddress(), owners, weights, threshold)
			if err != nil {
				return err
			}

			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// getCmdCreateTransaction is the CLI command for sending a CreateTransaction transaction.
func getCmdCreateTransaction() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-transaction [wallet] [receiver] [coins]",
		Short: "create a new multi-signature transaction",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error
			clientCtx := client.GetClientContextFromCmd(cmd)

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

			coins, err := sdk.ParseCoinsNormalized(coinsString)
			if err != nil {
				return err
			}

			// TODO: Check coins exist?
			// for _, c := range coins {
			// 	coin, err := cliUtils.GetCoin(clientCtx, c.Denom)
			// 	if err != nil {
			// 		return err
			// 	}
			// 	if coin.Symbol != c.Denom {
			// 		return sdkerrors.New(types.DefaultCodespace, types.InvalidCoinToSend, fmt.Sprintf("Coin to send with symbol %s does not exist", c.Denom))
			// 	}
			// }

			msg := types2.NewMsgCreateTransaction(clientCtx.GetFromAddress(), wallet, receiver, coins)
			if err != nil {
				return err
			}
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			fmt.Printf("%+v\n", msg)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// getCmdSignTransaction is the CLI command for saving a transaction signature.
func getCmdSignTransaction() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sign-transaction [tx-id]",
		Short: "Save a signature generated for a specific transaction",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error
			clientCtx := client.GetClientContextFromCmd(cmd)

			txID := args[0]

			msg := types2.NewMsgSignTransaction(clientCtx.GetFromAddress(), txID)
			if err != nil {
				return err
			}
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
