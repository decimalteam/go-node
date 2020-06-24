package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"

	cfg "github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/crypto"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"

	"bitbucket.org/decimalteam/go-node/x/validator/internal/types"
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

	validatorTxCmd.AddCommand(flags.PostCommands(
		GetCmdDeclareCandidate(cdc),
		GetDelegate(cdc),
		GetSetOnline(cdc),
		GetSetOffline(cdc),
		GetUnbond(cdc),
		GetEditCandidate(cdc),
	)...)

	return validatorTxCmd
}

// GetCmdDeclareCandidate is the CLI command for doing CreateValidator
func GetCmdDeclareCandidate(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Short: "Declare candidate",
		Use:   "declare [pub_key] [coin] [commission] --from name/address",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := auth.NewTxBuilderFromCLI(cliCtx.Input).WithTxEncoder(utils.GetTxEncoder(cdc))

			valAddress := cliCtx.GetFromAddress()

			stake, err := sdk.ParseCoin(args[1])
			if err != nil {
				return err
			}
			commission, err := sdk.NewDecFromStr(args[2])
			if err != nil {
				return err
			}

			pubKey, err := sdk.GetPubKeyFromBech32(sdk.Bech32PubKeyTypeConsPub, args[0])
			if err != nil {
				return err
			}

			rewardAddressStr := viper.GetString(FlagRewardAddress)
			rewardAddress := sdk.AccAddress{}
			if rewardAddressStr != "" {
				rewardAddress, err = sdk.AccAddressFromBech32(rewardAddressStr)
				if err != nil {
					return err
				}
			} else {
				rewardAddress = valAddress
			}

			msg := types.NewMsgDeclareCandidate(sdk.ValAddress(valAddress), pubKey, commission, stake, types.Description{}, rewardAddress)
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd.Flags().AddFlagSet(FsPk)
	cmd.Flags().AddFlagSet(FsAmount)
	cmd.Flags().AddFlagSet(fsDescriptionCreate)
	cmd.Flags().AddFlagSet(FsCommissionCreate)

	cmd.Flags().String(FlagIP, "", fmt.Sprintf("The node's public IP. It takes effect only when used in combination with --%s", flags.FlagGenerateOnly))
	cmd.Flags().String(FlagNodeID, "", "The node's ID")

	cmd.MarkFlagRequired(flags.FlagFrom)
	cmd.MarkFlagRequired(FlagAmount)
	cmd.MarkFlagRequired(FlagPubKey)
	cmd.MarkFlagRequired(FlagMoniker)

	return cmd
}

// CreateValidatorMsgHelpers returns the flagset, particular flags and a description of defaults
// this is anticipated to be used with the gen-tx.
func CreateValidatorMsgHelpers(ipDefault string) (fs *flag.FlagSet, nodeIDFlag, pubkeyFlag, amountFlag, defaultsDesc string) {

	fsCreateValidator := flag.NewFlagSet("", flag.ContinueOnError)
	fsCreateValidator.String(FlagIP, ipDefault, "The node's public IP")
	fsCreateValidator.String(FlagNodeID, "", "The node's NodeID")
	fsCreateValidator.String(FlagWebsite, "", "The validator's (optional) website")
	fsCreateValidator.String(FlagDetails, "", "The validator's (optional) details")
	fsCreateValidator.String(FlagIdentity, "", "The (optional) identity signature (ex. UPort or Keybase)")
	fsCreateValidator.AddFlagSet(FsCommissionCreate)
	fsCreateValidator.AddFlagSet(FsAmount)
	fsCreateValidator.AddFlagSet(FsPk)

	return fsCreateValidator, FlagNodeID, FlagPubKey, FlagAmount, defaultsDesc
}

// PrepareFlagsForTxCreateValidator prepare flags in config.
func PrepareFlagsForTxCreateValidator(config *cfg.Config, nodeID, chainID string, valPubKey crypto.PubKey) {

	ip := viper.GetString(FlagIP)
	if ip == "" {
		fmt.Println("couldn't retrieve an external IP; the tx's memo field will be unset")
	}

	website := viper.GetString(FlagWebsite)
	securityContact := viper.GetString(FlagSecurityContact)
	details := viper.GetString(FlagDetails)
	identity := viper.GetString(FlagIdentity)

	viper.Set(flags.FlagChainID, chainID)
	viper.Set(flags.FlagFrom, viper.GetString(flags.FlagName))
	viper.Set(FlagNodeID, nodeID)
	viper.Set(FlagIP, ip)
	viper.Set(FlagPubKey, sdk.MustBech32ifyPubKey(sdk.Bech32PubKeyTypeConsPub, valPubKey))
	viper.Set(FlagMoniker, config.Moniker)
	viper.Set(FlagWebsite, website)
	viper.Set(FlagSecurityContact, securityContact)
	viper.Set(FlagDetails, details)
	viper.Set(FlagIdentity, identity)

	if config.Moniker == "" {
		viper.Set(FlagMoniker, viper.GetString(flags.FlagName))
	}
	if viper.GetString(FlagAmount) == "" {
		viper.Set(FlagAmount, types.TokensFromConsensusPower(100).String()+types.DefaultBondDenom)
	}
	if viper.GetString(FlagCommissionRate) == "" {
		viper.Set(FlagCommissionRate, "0.1")
	}
}

// BuildCreateValidatorMsg makes a new MsgCreateValidator.
func BuildCreateValidatorMsg(cliCtx context.CLIContext, txBldr auth.TxBuilder) (auth.TxBuilder, sdk.Msg, error) {
	amounstStr := viper.GetString(FlagAmount)
	amount, err := sdk.ParseCoin(amounstStr)
	if err != nil {
		return txBldr, nil, err
	}

	valAddr := cliCtx.GetFromAddress()
	pkStr := viper.GetString(FlagPubKey)

	pk, err := sdk.GetPubKeyFromBech32(sdk.Bech32PubKeyTypeConsPub, pkStr)
	if err != nil {
		return txBldr, nil, err
	}

	description := types.NewDescription(
		viper.GetString(FlagMoniker),
		viper.GetString(FlagIdentity),
		viper.GetString(FlagWebsite),
		viper.GetString(FlagSecurityContact),
		viper.GetString(FlagDetails),
	)

	// get the initial validator commission parameters
	rateStr := viper.GetString(FlagCommissionRate)
	commission, err := sdk.NewDecFromStr(rateStr)
	if err != nil {
		return txBldr, nil, err
	}

	msg := types.NewMsgDeclareCandidate(sdk.ValAddress(valAddr), pk, commission, amount, description, valAddr)

	ip := viper.GetString(FlagIP)
	nodeID := viper.GetString(FlagNodeID)
	if nodeID != "" && ip != "" {
		txBldr = txBldr.WithMemo(fmt.Sprintf("%s@%s:26656", nodeID, ip))
	}

	return txBldr, msg, nil
}

// GetDelegate .
func GetDelegate(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Short: "Delegate coins",
		Use:   "delegate [validator_address] [coin] --from name/address",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := auth.NewTxBuilderFromCLI(cliCtx.Input).WithTxEncoder(utils.GetTxEncoder(cdc))

			delAddress := cliCtx.GetFromAddress()

			coin, err := sdk.ParseCoin(args[1])
			if err != nil {
				return err
			}

			valAddress, err := sdk.ValAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			msg := types.NewMsgDelegate(valAddress, delAddress, coin)
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

// GetSetOnline .
func GetSetOnline(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Short: "Set online validator",
		Use:   "set-online --from name/address",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := auth.NewTxBuilderFromCLI(cliCtx.Input).WithTxEncoder(utils.GetTxEncoder(cdc))

			valAddress := cliCtx.GetFromAddress()

			msg := types.NewMsgSetOnline(sdk.ValAddress(valAddress))
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

// GetSetOffline .
func GetSetOffline(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Short: "Set offline validator",
		Use:   "set-offline --from name/address",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := auth.NewTxBuilderFromCLI(cliCtx.Input).WithTxEncoder(utils.GetTxEncoder(cdc))

			valAddress := cliCtx.GetFromAddress()

			msg := types.NewMsgSetOffline(sdk.ValAddress(valAddress))
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

// GetUnbond .
func GetUnbond(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Short: "Unbond delegation",
		Use:   "unbond [validator-address] [coin] --from name/address",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := auth.NewTxBuilderFromCLI(cliCtx.Input).WithTxEncoder(utils.GetTxEncoder(cdc))

			valAddress, err := sdk.ValAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			delAddress := cliCtx.GetFromAddress()

			coin, err := sdk.ParseCoin(args[1])
			if err != nil {
				return err
			}

			msg := types.NewMsgUnbond(valAddress, delAddress, coin)
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

// GetEditCandidate .
func GetEditCandidate(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Short: "Edit candidate",
		Use:   "edit-candidate [pub_key] [validator-address] [reward-address]",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := auth.NewTxBuilderFromCLI(cliCtx.Input).WithTxEncoder(utils.GetTxEncoder(cdc))

			valAddress, err := sdk.ValAddressFromBech32(args[1])
			if err != nil {
				return err
			}

			rewardAddress, err := sdk.AccAddressFromBech32(args[2])
			if err != nil {
				return err
			}

			description := types.NewDescription(
				viper.GetString(FlagMoniker),
				viper.GetString(FlagIdentity),
				viper.GetString(FlagWebsite),
				viper.GetString(FlagSecurityContact),
				viper.GetString(FlagDetails),
			)

			msg := types.NewMsgEditCandidate(valAddress, rewardAddress, description)
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd.Flags().AddFlagSet(fsDescriptionEdit)
	cmd.Flags().AddFlagSet(fsCommissionUpdate)

	return cmd
}
