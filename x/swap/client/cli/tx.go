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
	"strconv"
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
		GetSwapInitialize(cdc),
		GetRedeemV2(cdc),
		GetChainActivate(cdc),
		GetChainDeactivate(cdc),
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

func GetSwapInitialize(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init [recipient] [amount] [token_symbol] [tx_number] [from_chain] [dest_chain] --from",
		Short: "Swap initialize",
		Args:  cobra.ExactArgs(7),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := auth.NewTxBuilderFromCLI(cliCtx.Input).WithTxEncoder(utils.GetTxEncoder(cdc))

			from := cliCtx.GetFromAddress()

			recipient := args[0]
			amount, ok := sdk.NewIntFromString(args[1])
			if !ok {
				return fmt.Errorf("invalid amount")
			}
			symbol := args[2]
			txNumber := args[3]
			fromChain, err := strconv.Atoi(args[4])
			if err != nil {
				return err
			}
			destChain, err := strconv.Atoi(args[5])
			if err != nil {
				return err
			}

			msg := types.NewMsgSwapInitialize(from, recipient, amount, symbol, txNumber, fromChain, destChain)
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
	return cmd
}

func GetRedeemV2(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "redeemV2 [from] [recipient] [amount] [token_symbol] [tx_number] [from_chain] [dest_chain] [v] [r] [s] --from",
		Short: "Swap initialize",
		Args:  cobra.ExactArgs(11),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := auth.NewTxBuilderFromCLI(cliCtx.Input).WithTxEncoder(utils.GetTxEncoder(cdc))

			sender := cliCtx.GetFromAddress()

			from := args[0]
			recipient, err := sdk.AccAddressFromBech32(args[1])
			if err != nil {
				return err
			}
			amount, ok := sdk.NewIntFromString(args[2])
			if !ok {
				return fmt.Errorf("invalid amount")
			}
			symbol := args[3]
			txNumber := args[4]
			fromChain, err := strconv.Atoi(args[5])
			if err != nil {
				return err
			}
			destChain, err := strconv.Atoi(args[6])
			if err != nil {
				return err
			}

			v, err := strconv.Atoi(args[7])
			if err != nil {
				return err
			}

			r, err := hex.DecodeString(args[8])
			s, err := hex.DecodeString(args[9])

			var _r [32]byte
			copy(_r[:], r)

			var _s [32]byte
			copy(_s[:], s)

			msg := types.NewMsgRedeemV2(
				sender, recipient, from, amount, symbol, txNumber, fromChain, destChain, uint8(v), _r, _s)
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
	return cmd
}

func GetChainActivate(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "chain-activate [number] [name] --from",
		Short: "Activate chain",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := auth.NewTxBuilderFromCLI(cliCtx.Input).WithTxEncoder(utils.GetTxEncoder(cdc))

			from := cliCtx.GetFromAddress()

			number, err := strconv.Atoi(args[0])
			if err != nil {
				return err
			}

			name := args[1]

			msg := types.NewMsgChainActivate(from, number, name)
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
	return cmd
}

func GetChainDeactivate(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "chain-deactivate [number] --from",
		Short: "Deactivate chain",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := auth.NewTxBuilderFromCLI(cliCtx.Input).WithTxEncoder(utils.GetTxEncoder(cdc))

			from := cliCtx.GetFromAddress()

			number, err := strconv.Atoi(args[0])
			if err != nil {
				return err
			}

			msg := types.NewMsgChainDeactivate(from, number)
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
	return cmd
}
