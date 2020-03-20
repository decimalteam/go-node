package cli

import (
	"bitbucket.org/decimalteam/go-node/x/validator/internal/types"
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
	cfg "github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/crypto"
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
		GetCmdDeclareCandidate(cdc),
	)...)

	return validatorTxCmd
}

// GetCmdDeclareCandidate is the CLI command for doing CreateValidator
func GetCmdDeclareCandidate(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Short: "Declare candidate",
		Use:   "declare [pub_key] [coin] [commission]",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			valAddress := cliCtx.GetFromAddress()

			stake, err := sdk.ParseCoin(args[1])
			if err != nil {
				return err
			}
			commission, err := sdk.NewDecFromStr(args[2])
			if err != nil {
				return err
			}

			pubKeyStr := args[0]

			pubKey, err := sdk.GetConsPubKeyBech32(pubKeyStr)
			if err != nil {
				return err
			}

			msg := types.NewMsgDeclareCandidate(valAddress, pubKey, types.Commission{Rate: commission}, stake, types.Description{})
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

//__________________________________________________________

var (
	defaultTokens         = sdk.TokensFromConsensusPower(100)
	defaultAmount         = defaultTokens.String() + sdk.DefaultBondDenom
	defaultCommissionRate = "0.1"
)

// Return the flagset, particular flags, and a description of defaults
// this is anticipated to be used with the gen-tx
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

// prepare flags in config
func PrepareFlagsForTxCreateValidator(
	config *cfg.Config, nodeID, chainID string, valPubKey crypto.PubKey,
) {

	ip := viper.GetString(FlagIP)
	if ip == "" {
		fmt.Println("couldn't retrieve an external IP; the tx's memo field will be unset")
	}

	website := viper.GetString(FlagWebsite)
	details := viper.GetString(FlagDetails)
	identity := viper.GetString(FlagIdentity)

	viper.Set(client.FlagChainID, chainID)
	viper.Set(client.FlagFrom, viper.GetString(client.FlagName))
	viper.Set(FlagNodeID, nodeID)
	viper.Set(FlagIP, ip)
	viper.Set(FlagPubKey, sdk.MustBech32ifyConsPub(valPubKey))
	viper.Set(FlagMoniker, config.Moniker)
	viper.Set(FlagWebsite, website)
	viper.Set(FlagDetails, details)
	viper.Set(FlagIdentity, identity)

	if config.Moniker == "" {
		viper.Set(FlagMoniker, viper.GetString(client.FlagName))
	}
	if viper.GetString(FlagAmount) == "" {
		viper.Set(FlagAmount, defaultAmount)
	}
	if viper.GetString(FlagCommissionRate) == "" {
		viper.Set(FlagCommissionRate, defaultCommissionRate)
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

	pk, err := sdk.GetConsPubKeyBech32(pkStr)
	if err != nil {
		return txBldr, nil, err
	}

	description := staking.Description{
		Moniker:  viper.GetString(FlagMoniker),
		Identity: viper.GetString(FlagIdentity),
		Website:  viper.GetString(FlagWebsite),
		Details:  viper.GetString(FlagDetails),
	}

	// get the initial validator commission parameters
	rateStr := viper.GetString(FlagCommissionRate)
	commission, err := sdk.NewDecFromStr(rateStr)
	if err != nil {
		return txBldr, nil, err
	}
	commissionRates := types.Commission{Rate: commission}

	msg := types.NewMsgDeclareCandidate(valAddr, pk, commissionRates, amount, types.Description(description))

	ip := viper.GetString(FlagIP)
	nodeID := viper.GetString(FlagNodeID)
	if nodeID != "" && ip != "" {
		txBldr = txBldr.WithMemo(fmt.Sprintf("%s@%s:26656", nodeID, ip))
	}

	return txBldr, msg, nil
}
