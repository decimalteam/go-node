package cli

import (
	"errors"
	"fmt"
	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"

	"github.com/spf13/cobra"

	"bitbucket.org/decimalteam/go-node/x/validator/internal/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd(cdc *codec.Codec) *cobra.Command {
	validatorTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%s transactions subcommands", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	validatorTxCmd.AddCommand(client.PostCommands(
		GetCmdCreateValidator(cdc),
	)...)

	return validatorTxCmd
}

// GetCmdCreateValidator is the CLI command for doing CreateValidator
func GetCmdCreateValidator(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "Create validator",
		Short: "create [val_address] [pub_key] [amount] [coin] [commission]",
		Args:  cobra.ExactArgs(5),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			valAddress, err := sdk.ValAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			amount, ok := sdk.NewIntFromString(args[2])
			if !ok {
				return errors.New("Invalid amount ")
			}

			stake := sdk.NewCoin(args[3], amount)
			commission, err := sdk.NewDecFromStr(args[4])
			if err != nil {
				return err
			}

			pubKeyStr := args[1]

			pubKey, err := sdk.GetConsPubKeyBech32(pubKeyStr)
			if err != nil {
				return err
			}

			msg := types.NewMsgCreateValidator(valAddress, pubKey, types.Commission{Rate: commission}, stake)
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}
