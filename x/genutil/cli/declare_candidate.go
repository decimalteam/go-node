package cli

import (
	"bytes"
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	client2 "github.com/cosmos/cosmos-sdk/x/auth/client"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	kbkeys "github.com/cosmos/cosmos-sdk/crypto"
	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/genutil/types"

	"bitbucket.org/decimalteam/go-node/x/genutil"
)

// GenDeclareCandidateTxCmd builds and prints transaction to declare validator candidate.
// nolint: errcheck
func GenDeclareCandidateTxCmd(ctx *server.Context, cdc *codec.LegacyAmino, mbm module.BasicManager, smbh StakingMsgBuildingHelpers,
	genAccIterator types.GenesisAccountsIterator, defaultNodeHome, defaultCLIHome string) *cobra.Command {

	ipDefault, _ := server.ExternalIP()
	fsCreateValidator, flagNodeID, flagPubKey, flagAmount, defaultsDesc := smbh.CreateValidatorMsgHelpers(ipDefault)

	cmd := &cobra.Command{
		Use:   "gen-declare-candidate-tx",
		Short: "Generate a genesis declare validator candidate tx carrying a self delegation",
		Args:  cobra.NoArgs,
		Long: fmt.Sprintf(`This command is an alias of the 'tx create-validator' command'.

		It creates a genesis transaction to create a validator. 
		The following default parameters are included: 
		    %s`, defaultsDesc),

		RunE: func(cmd *cobra.Command, args []string) error {
			config := ctx.Config
			config.SetRoot(viper.GetString(flags.FlagHome))
			nodeID, valPubKey, err := genutil.InitializeNodeValidatorFiles(ctx.Config)
			if err != nil {
				return err
			}

			// Read --nodeID, if empty take it from priv_validator.json
			if nodeIDString := viper.GetString(flagNodeID); nodeIDString != "" {
				nodeID = nodeIDString
			}
			// Read --pubkey, if empty take it from priv_validator.json
			if valPubKeyString := viper.GetString(flagPubKey); valPubKeyString != "" {
				valPubKey, err = sdk.GetFromBech32(sdk.Bech32PrefixConsPub, valPubKeyString)
				if err != nil {
					return err
				}
			}

			// Read chain ID
			chainID := viper.GetString(flags.FlagChainID)
			if len(chainID) == 0 {
				return fmt.Errorf("chain ID must be specified")
			}

			_, err = kbkeys.NewKeyring(sdk.KeyringServiceName(), viper.GetString(flags.FlagKeyringBackend), viper.GetString(flagClientHome), cmd.InOrStdin())
			if err != nil {
				return err
			}

			name := viper.GetString(flags.FlagName)

			// Set flags for creating declare validator candidate tx
			viper.Set(flags.FlagHome, viper.GetString(flagClientHome))
			smbh.PrepareFlagsForTxCreateValidator(config, nodeID, chainID, valPubKey)

			// Fetch the amount of coins staked
			amount := viper.GetString(flagAmount)
			_, err = sdk.ParseCoinNormalized(amount)
			if err != nil {
				return err
			}

			clientCtx := client.GetClientContextFromCmd(cmd)

			viper.Set(flags.FlagGenerateOnly, true)

			// create a 'create-validator' message
			txBldr, msg, err := smbh.BuildCreateValidatorMsg(clientCtx, txBldr)
			if err != nil {
				return err
			}

			info, err := txBldr.Keybase().Get(name)
			if err != nil {
				return err
			}

			if info.GetType() == kbkeys.TypeOffline || info.GetType() == kbkeys.TypeMulti {
				fmt.Println("Offline key passed in. Use `tx sign` command to sign:")
				return utils.PrintUnsignedStdTx(txBldr, clientCtx, []sdk.Msg{msg})
			}

			// write the unsigned transaction to the buffer
			w := bytes.NewBuffer([]byte{})
			clientCtx = clientCtx.WithOutput(w)
			txFactory := tx.Factory{}

			if err = client2.PrintUnsignedStdTx(txFactory, clientCtx, []sdk.Msg{msg}); err != nil {
				return err
			}

			// read the transaction
			stdTx, err := readUnsignedGenTxFile(cdc, w)
			if err != nil {
				return err
			}

			// sign the transaction and write it to the output file
			signedTx, err := utils.SignStdTx(txBldr, clientCtx, name, stdTx, false, true)
			if err != nil {
				return err
			}

			txJSON, err := cdc.MarshalJSONIndent(signedTx, "", "  ")
			if err != nil {
				return err
			}
			fmt.Println(string(txJSON))
			return nil

		},
	}

	cmd.Flags().String(flags.FlagHome, defaultNodeHome, "node's home directory")
	cmd.Flags().String(flagClientHome, defaultCLIHome, "client's home directory")
	cmd.Flags().String(flags.FlagName, "", "name of private key with which to sign the gentx")
	cmd.Flags().String(flags.FlagChainID, "", "Chain ID for which tx should be signed")
	cmd.Flags().AddFlagSet(fsCreateValidator)

	return cmd
}
