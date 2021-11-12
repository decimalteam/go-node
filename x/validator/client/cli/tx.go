package cli

import (
	"fmt"
	tx "github.com/cosmos/cosmos-sdk/client/tx"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	cli2 "github.com/cosmos/cosmos-sdk/x/staking/client/cli"
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"

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
		GetCmdDeclareCandidate(),
		GetDelegate(),
		GetDelegateNFT(),
		GetUnbondNFT(),
		GetSetOnline(),
		GetSetOffline(),
		GetUnbond(),
		GetEditCandidate(),
	)

	return validatorTxCmd
}

// GetCmdDeclareCandidate is the CLI command for doing CreateValidator
func GetCmdDeclareCandidate() *cobra.Command {
	cmd := &cobra.Command{
		Short: "Declare candidate",
		Use:   "declare [pub_key] [coin] [commission] [from]",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.Flags().Set(flags.FlagFrom, args[3])

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

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
			if err := clientCtx.JSONMarshaler.UnmarshalInterfaceJSON([]byte(args[0]), &pk); err != nil {
				return err
			}

			rewardAddressStr, _ := cmd.Flags().GetString(FlagRewardAddress)
			rewardAddress := sdk.AccAddress{}
			if rewardAddressStr != "" {
				rewardAddress, err = sdk.AccAddressFromBech32(rewardAddressStr)
				if err != nil {
					return err
				}
			} else {
				rewardAddress = valAddress
			}

			msg, err := types.NewMsgDeclareCandidate(sdk.ValAddress(valAddress), pk, commission, stake, types.Description{}, rewardAddress)
			if err != nil {
				return err
			}

			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().AddFlagSet(FsPk)
	cmd.Flags().AddFlagSet(FsAmount)
	cmd.Flags().AddFlagSet(FsDescriptionCreate)
	cmd.Flags().AddFlagSet(FsCommissionCreate)

	cmd.Flags().String(FlagIP, "", fmt.Sprintf("The node's public IP. It takes effect only when used in combination with --%s", flags.FlagGenerateOnly))
	cmd.Flags().String(FlagNodeID, "", "The node's ID")

	//cmd.MarkFlagRequired(flags.FlagFrom)
	cmd.MarkFlagRequired(FlagAmount)
	cmd.MarkFlagRequired(FlagPubKey)
	cmd.MarkFlagRequired(FlagMoniker)

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// GetCmdDeclareCandidate is the CLI command for doing CreateValidator
/*func GetCmdDeclareCandidate(cdc *codec.LegacyAmino) *cobra.Command {
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
			if err := clientCtx.JSONMarshaler.UnmarshalInterfaceJSON([]byte(args[0]), &pk); err != nil {
				return err
			}

			rewardAddressStr, _ := cmd.Flags().GetString(FlagRewardAddress)
			rewardAddress := sdk.AccAddress{}
			if rewardAddressStr != "" {
				rewardAddress, err = sdk.AccAddressFromBech32(rewardAddressStr)
				if err != nil {
					return err
				}
			} else {
				rewardAddress = valAddress
			}

			msg, err := types.NewMsgDeclareCandidate(sdk.ValAddress(valAddress), pk, commission, stake, types.Description{}, rewardAddress)
			if err != nil {
				return err
			}

			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
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
}*/

// CreateValidatorMsgHelpers returns the flagset, particular flags and a description of defaults
// this is anticipated to be used with the gen-tx.
func CreateValidatorMsgHelpers(ipDefault string) (fs *flag.FlagSet, defaultsDesc string) {

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

	return fsCreateValidator, defaultsDesc
}

// PrepareFlagsForTxCreateValidator prepare flags in config.
func PrepareFlagsForTxCreateValidator(flagSet *flag.FlagSet, moniker, nodeID, chainID string, valPubKey cryptotypes.PubKey) (cli2.TxCreateValidatorConfig, error) {
	c := cli2.TxCreateValidatorConfig{}

	ip, _ := flagSet.GetString(FlagIP)
	if ip == "" {
		fmt.Println("couldn't retrieve an external IP; the tx's memo field will be unset")
	}
	c.IP = ip

	website, err := flagSet.GetString(FlagWebsite)
	if err != nil {
		return c, err
	}
	c.Website = website

	securityContact, err := flagSet.GetString(FlagSecurityContact)
	if err != nil {
		return c, err
	}
	c.SecurityContact = securityContact

	details, err := flagSet.GetString(FlagDetails)
	if err != nil {
		return c, err
	}
	c.SecurityContact = details

	identity, err := flagSet.GetString(FlagIdentity)
	if err != nil {
		return c, err
	}
	c.Identity = identity

	c.Amount, err = flagSet.GetString(FlagAmount)
	if err != nil {
		return c, err
	}

	c.CommissionRate, err = flagSet.GetString(FlagCommissionRate)
	if err != nil {
		return c, err
	}

	c.NodeID = nodeID
	c.PubKey = sdk.MustBech32ifyPubKey(sdk.Bech32PubKeyTypeConsPub, valPubKey)
	c.Website = website
	c.SecurityContact = securityContact
	c.Details = details
	c.Identity = identity
	c.ChainID = chainID
	c.Moniker = moniker

	if c.Amount == "" {
		c.Amount = types.TokensFromConsensusPower(100).String() + types.DefaultBondDenom
	}

	if c.CommissionRate == "" {
		c.CommissionRate = "0.1"
	}

	if c.Moniker == "" {
		c.Moniker, _ = flagSet.GetString(flags.FlagName)
	}

	viper.Set(flags.FlagFrom, viper.GetString(flags.FlagName))
	viper.Set(FlagRewardAddress, viper.Get(FlagRewardAddress))

	return c, nil
}

// BuildCreateValidatorMsg makes a new MsgCreateValidator.
func BuildCreateValidatorMsg(clientCtx client.Context, config cli2.TxCreateValidatorConfig, txBldr tx.Factory, fs *flag.FlagSet, generateOnly bool) (tx.Factory, sdk.Msg, error) {
	amounstStr := config.Amount
	amount, err := sdk.ParseCoinNormalized(amounstStr)
	if err != nil {
		return txBldr, nil, err
	}

	valAddr := clientCtx.GetFromAddress()
	rewardAddr := valAddr
	rewardAddrStr, _ := fs.GetString(FlagRewardAddress)
	pkStr := config.PubKey

	if len(rewardAddrStr) > 0 {
		rewardAddr, err = sdk.AccAddressFromBech32(rewardAddrStr)
		if err != nil {
			return txBldr, nil, err
		}
	}

	pk, err := sdk.GetPubKeyFromBech32(sdk.Bech32PubKeyTypeConsPub, pkStr)
	if err != nil {
		return txBldr, nil, err
	}

	description := types.NewDescription(
		config.Moniker,
		config.Identity,
		config.Website,
		config.SecurityContact,
		config.Details,
	)

	// get the initial validator commission parameters
	rateStr := config.CommissionRate
	commission, err := sdk.NewDecFromStr(rateStr)
	if err != nil {
		return txBldr, nil, err
	}

	msg, err := types.NewMsgDeclareCandidate(sdk.ValAddress(valAddr), pk, commission, amount, description, rewardAddr)

	if err != nil {
		return txBldr, nil, err
	}

	// NOTE: No need to show public IP of the node
	if generateOnly {
		ip := config.IP
		nodeID := config.NodeID

		if nodeID != "" && ip != "" {
			txBldr = txBldr.WithMemo(fmt.Sprintf("%s@%s:26656", nodeID, ip))
		}
	}

	return txBldr, msg, nil
}

// GetDelegate .
func GetDelegate() *cobra.Command {
	cmd := &cobra.Command{
		Short: "Delegate coins",
		Use:   "delegate [validator_address] [coin] [from]",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.Flags().Set(flags.FlagFrom, args[2])

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

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


/*func GetDelegate(cdc *codec.LegacyAmino) *cobra.Command {
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
}*/


func GetDelegateNFT() *cobra.Command {
	cmd := &cobra.Command{
		Short: "Delegate NFT",
		Use:   "delegate-nft [validator_address] [tokenID] [denom] [quantity] [from]",
		Args:  cobra.ExactArgs(5),
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.Flags().Set(flags.FlagFrom, args[4])

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

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

/*func GetDelegateNFT(cdc *codec.LegacyAmino) *cobra.Command {
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
}*/

// GetSetOnline .

func GetSetOnline() *cobra.Command {
	return &cobra.Command{
		Short: "Set online validator",
		Use:   "set-online [from]",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.Flags().Set(flags.FlagFrom, args[0])

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			valAddress := clientCtx.GetFromAddress()

			msg := types.NewMsgSetOnline(sdk.ValAddress(valAddress))
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}

}

/*func GetSetOnline(cdc *codec.LegacyAmino) *cobra.Command {
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
}*/

// GetSetOffline .
func GetSetOffline() *cobra.Command {
	return &cobra.Command{
		Short: "Set offline validator",
		Use:   "set-offline [from]",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.Flags().Set(flags.FlagFrom, args[0])

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			valAddress := clientCtx.GetFromAddress()

			msg := types.NewMsgSetOffline(sdk.ValAddress(valAddress))
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}

}

/*func GetSetOffline(cdc *codec.LegacyAmino) *cobra.Command {
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
}*/

// GetUnbond .
func GetUnbond() *cobra.Command {
	cmd := &cobra.Command{
		Short: "Unbond delegation",
		Use:   "unbond [validator-address] [coin] [from]",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.Flags().Set(flags.FlagFrom, args[2])

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

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
/*
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
}*/

// GetUnbondNFT .
func GetUnbondNFT() *cobra.Command {
	cmd := &cobra.Command{
		Short: "Unbond-nft delegation",
		Use:   "unbond-nft [validator-address] [tokenID] [denom] [quantity] [from]",
		Args:  cobra.ExactArgs(5),
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.Flags().Set(flags.FlagFrom, args[4])

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

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
/*
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
}*/

// GetEditCandidate .
func GetEditCandidate() *cobra.Command {
	cmd := &cobra.Command{
		Short: "Edit candidate",
		Use:   "edit-candidate [validator-address] [reward-address]",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.Flags().Set(flags.FlagFrom, args[4])

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

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
/*
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
}*/
