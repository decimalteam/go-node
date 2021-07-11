package cli

import (
	"bitbucket.org/decimalteam/go-node/x/swap/internal/types"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/libs/rand"
)

func GetTxCmd(cdc *codec.Codec) *cobra.Command {
	swapTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%s transactions subcommands", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	swapTxCmd.AddCommand(flags.PostCommands(
		GetCmdHTLT(cdc),
		GetCmdRedeem(cdc),
		GetRefund(cdc),
	)...)

	return swapTxCmd
}

func GetCmdHTLT(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "htlt [in | out] [recipient] [amount] [--hash] --from",
		Short: "Create swap",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := auth.NewTxBuilderFromCLI(cliCtx.Input).WithTxEncoder(utils.GetTxEncoder(cdc))

			from := cliCtx.GetFromAddress()
			recipient := args[1]
			amount, err := sdk.ParseCoins(args[2])
			if err != nil {
				return err
			}

			transferType, err := types.TransferTypeFromString(args[0])
			if err != nil {
				return err
			}

			var hash [32]byte
			var secret [32]byte
			hashStr := viper.GetString(FlagHash)
			if hashStr == "" {
				copy(secret[:], rand.Bytes(32))
				hash = sha256.Sum256(secret[:])
				fmt.Println("Secret = ", hex.EncodeToString(secret[:]))
			} else {
				h, err := hex.DecodeString(hashStr)
				if err != nil {
					return err
				}
				copy(hash[:], h)
			}

			msg := types.NewMsgHTLT(
				transferType,
				from,
				recipient,
				hash,
				amount,
			)
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd.Flags().AddFlagSet(FsHash)

	return cmd
}

func GetCmdRedeem(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "redeem [secret] --from",
		Short: "Redeem swap",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := auth.NewTxBuilderFromCLI(cliCtx.Input).WithTxEncoder(utils.GetTxEncoder(cdc))

			from := cliCtx.GetFromAddress()

			secretStr := args[0]
			secret, err := hex.DecodeString(secretStr)
			if err != nil {
				return err
			}

			msg := types.NewMsgRedeem(from, secret)
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
	return cmd
}

func GetRefund(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "refund [hash] --from",
		Short: "Refund locked coins",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := auth.NewTxBuilderFromCLI(cliCtx.Input).WithTxEncoder(utils.GetTxEncoder(cdc))

			from := cliCtx.GetFromAddress()

			var hash [32]byte
			hashStr := args[0]
			h, err := hex.DecodeString(hashStr)
			if err != nil {
				return err
			}
			copy(hash[:], h)

			msg := types.NewMsgRefund(from, hash)
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	return cmd
}
