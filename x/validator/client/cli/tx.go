package cli

import (
	"fmt"
	tx "github.com/cosmos/cosmos-sdk/client/tx"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	cli2 "github.com/cosmos/cosmos-sdk/x/staking/client/cli"

	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"

	cfg "github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/crypto"

	"bitbucket.org/decimalteam/go-node/x/validator/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd(cdc *codec.LegacyAmino) *cobra.Command {
	validatorTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%s transactions subcommands", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	validatorTxCmd.AddCommand(
		GetCmdDeclareCandidate(cdc),
		GetDelegate(cdc),
		GetDelegateNFT(cdc),
		GetUnbondNFT(cdc),
		GetSetOnline(cdc),
		GetSetOffline(cdc),
		GetUnbond(cdc),
		GetEditCandidate(cdc),
	)

	return validatorTxCmd
}

// GetCmdDeclareCandidate is the CLI command for doing CreateValidator
func GetCmdDeclareCandidate(cdc *codec.LegacyAmino) *cobra.Command {
	cmd := &cobra.Command{
		Short: "Declare candidate",
		Use:   "declare [pub_key] [coin] [commission] --from name/address",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd).WithLegacyAmino(cdc)

			valAddress := clientCtx.GetFromAddress()

			stake, err := sdk.ParseCoinNormalized(args[1])
			if err != nil {
				return err
			}
			commission, err := sdk.NewDecFromStr(args[2])
			if err != nil {
				return err
			}

			var pk cryptotypes.PubKey
			if err := clientCtx.Codec.UnmarshalInterfaceJSON([]byte(args[0]), &pk); err != nil {
				return err
			}
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

			msg := types.NewMsgDeclareCandidate(sdk.ValAddress(valAddress), pk, commission, stake, types.Description{}, rewardAddress)
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}

	cmd.Flags().AddFlagSet(FsPk)
	cmd.Flags().AddFlagSet(FsAmount)
	cmd.Flags().AddFlagSet(FsDescriptionCreate)
	cmd.Flags().AddFlagSet(FsCommissionCreate)

	cmd.Flags().String(FlagIP, "", fmt.Sprintf("The node's public IP. It takes effect only when used in combination with --%s", flags.FlagGenerateOnly))
	cmd.Flags().String(FlagNodeID, "", "The node's ID")

	cmd.MarkFlagRequired(flags.FlagFrom)
	cmd.MarkFlagRequired(FlagAmount)
	cmd.MarkFlagRequired(FlagPubKey)
	cmd.MarkFlagRequired(FlagMoniker)

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// CreateValidatorMsgHelpers returns the flagset, particular flags and a description of defaults
// this is anticipated to be used with the gen-tx.
func CreateValidatorMsgHelpers(ipDefault string) (fs *flag.FlagSet, nodeIDFlag, pubkeyFlag, amountFlag, defaultsDesc string) {

	fsCreateValidator := flag.NewFlagSet("", flag.ContinueOnError)
	fsCreateValidator.String(FlagIP, ipDefault, "The node's public IP")
	fsCreateValidator.String(FlagNodeID, "", "The node's NodeID")
	fsCreateValidator.String(FlagMoniker, "", "The validator's (optional) moniker")
	fsCreateValidator.String(FlagWebsite, "", "The validator's (optional) website")
	fsCreateValidator.String(FlagDetails, "", "The validator's (optional) details")
	fsCreateValidator.String(FlagSecurityContact, "", "The (optional) security contract")
	fsCreateValidator.String(FlagIdentity, "", "The (optional) identity signature (ex. UPort or Keybase)")
	fsCreateValidator.String(FlagRewardAddress, "", "Address of account receiving validator's rewards (optional)")

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
	rewardAddr := viper.GetString(FlagRewardAddress)

	viper.Set(flags.FlagChainID, chainID)
	viper.Set(flags.FlagFrom, viper.GetString(flags.FlagName))
	viper.Set(FlagNodeID, nodeID)
	viper.Set(FlagIP, ip)
	viper.Set(FlagPubKey, sdk.MustBech32ifyAddressBytes(sdk.Bech32PrefixConsPub, valPubKey.Bytes()))
	viper.Set(FlagMoniker, config.Moniker)
	viper.Set(FlagWebsite, website)
	viper.Set(FlagSecurityContact, securityContact)
	viper.Set(FlagDetails, details)
	viper.Set(FlagIdentity, identity)
	viper.Set(FlagRewardAddress, rewardAddr)

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
func BuildCreateValidatorMsg(clientCtx client.Context, config cli2.TxCreateValidatorConfig, txBldr tx.Factory) (tx.Factory, sdk.Msg, error) {
	amounstStr := viper.GetString(FlagAmount)
	amount, err := sdk.ParseCoinNormalized(amounstStr)
	if err != nil {
		return txBldr, nil, err
	}

	valAddr := clientCtx.GetFromAddress()
	rewardAddr := valAddr
	rewardAddrStr := viper.GetString(FlagRewardAddress)
	if len(rewardAddrStr) > 0 {
		rewardAddr, err = sdk.AccAddressFromBech32(rewardAddrStr)
		if err != nil {
			return txBldr, nil, err
		}
	}
	pkStr := viper.GetString(FlagPubKey)

	var pk cryptotypes.PubKey
	if err := clientCtx.JSONCodec.UnmarshalInterfaceJSON([]byte(pkStr), &pk); err != nil {
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

	msg := types.NewMsgDeclareCandidate(sdk.ValAddress(valAddr), pk, commission, amount, description, rewardAddr)

	// NOTE: No need to show public IP of the node
	// ip := viper.GetString(FlagIP)
	// nodeID := viper.GetString(FlagNodeID)
	// if nodeID != "" && ip != "" {
	// 	txBldr = txBldr.WithMemo(fmt.Sprintf("%s@%s:26656", nodeID, ip))
	// }

	return txBldr, &msg, nil
}

// GetDelegate .
func GetDelegate(cdc *codec.LegacyAmino) *cobra.Command {
	cmd := &cobra.Command{
		Short: "Delegate coins",
		Use:   "delegate [validator_address] [coin] --from name/address",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd).WithLegacyAmino(cdc)

			delAddress := clientCtx.GetFromAddress()

			coin, err := sdk.ParseCoinNormalized(args[1])
			if err != nil {
				return err
			}

			valAddress, err := sdk.ValAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			msg := types.NewMsgDelegate(valAddress, delAddress, coin)
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func GetDelegateNFT(cdc *codec.LegacyAmino) *cobra.Command {
	cmd := &cobra.Command{
		Short: "Delegate NFT",
		Use:   "delegate-nft [validator_address] [tokenID] [denom] [quantity] --from name/address",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd).WithLegacyAmino(cdc)

			delAddress := clientCtx.GetFromAddress()

			quantity, ok := sdk.NewIntFromString(args[3])
			if !ok {
				return fmt.Errorf("invalid quantity")
			}

			valAddress, err := sdk.ValAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			tokenID := args[1]
			denom := args[2]

			msg := types.NewMsgDelegateNFT(valAddress, delAddress, tokenID, denom, quantity)
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// GetSetOnline .
func GetSetOnline(cdc *codec.LegacyAmino) *cobra.Command {
	return &cobra.Command{
		Short: "Set online validator",
		Use:   "set-online --from name/address",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd).WithLegacyAmino(cdc)

			valAddress := clientCtx.GetFromAddress()

			msg := types.NewMsgSetOnline(sdk.ValAddress(valAddress))
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}
}

// GetSetOffline .
func GetSetOffline(cdc *codec.LegacyAmino) *cobra.Command {
	return &cobra.Command{
		Short: "Set offline validator",
		Use:   "set-offline --from name/address",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd).WithLegacyAmino(cdc)

			valAddress := clientCtx.GetFromAddress()

			msg := types.NewMsgSetOffline(sdk.ValAddress(valAddress))
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}
}

// GetUnbond .
func GetUnbond(cdc *codec.LegacyAmino) *cobra.Command {
	cmd := &cobra.Command{
		Short: "Unbond delegation",
		Use:   "unbond [validator-address] [coin] --from name/address",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd).WithLegacyAmino(cdc)

			valAddress, err := sdk.ValAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			delAddress := clientCtx.GetFromAddress()

			coin, err := sdk.ParseCoinNormalized(args[1])
			if err != nil {
				return err
			}

			msg := types.NewMsgUnbond(valAddress, delAddress, coin)
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// GetUnbondNFT .
func GetUnbondNFT(cdc *codec.LegacyAmino) *cobra.Command {
	cmd := &cobra.Command{
		Short: "Unbond-nft delegation",
		Use:   "unbond-nft [validator-address] [tokenID] [denom] [quantity] --from name/address",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd).WithLegacyAmino(cdc)

			valAddress, err := sdk.ValAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			delAddress := clientCtx.GetFromAddress()

			quantity, ok := sdk.NewIntFromString(args[3])
			if !ok {
				return fmt.Errorf("invalid quantity")
			}

			tokenID := args[1]
			denom := args[2]

			msg := types.NewMsgUnbondNFT(valAddress, delAddress, tokenID, denom, quantity)
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// GetEditCandidate .
func GetEditCandidate(cdc *codec.LegacyAmino) *cobra.Command {
	cmd := &cobra.Command{
		Short: "Edit candidate",
		Use:   "edit-candidate [validator-address] [reward-address]",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd).WithLegacyAmino(cdc)

			valAddress, err := sdk.ValAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			rewardAddress, err := sdk.AccAddressFromBech32(args[1])
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
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}

	cmd.Flags().AddFlagSet(FsDescriptionEdit)
	cmd.Flags().AddFlagSet(FsCommissionUpdate)
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
