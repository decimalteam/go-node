package cli

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	client2 "github.com/cosmos/cosmos-sdk/x/auth/client"
	"github.com/pkg/errors"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"bitbucket.org/decimalteam/go-node/x/genutil"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
)

// GenDeclareCandidateTxCmd builds and prints transaction to declare validator candidate.
// nolint: errcheck
func GenDeclareCandidateTxCmd(ctx *server.Context, mbm module.BasicManager, smbh StakingMsgBuildingHelpers, defaultNodeHome, defaultCLIHome string) *cobra.Command {

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
			clientCtx := client.GetClientContextFromCmd(cmd)

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
				_, err := sdk.GetFromBech32(sdk.Bech32PrefixConsPub, valPubKeyString)
				if err != nil {
					return err
				}
			}

			// Read chain ID
			chainID := viper.GetString(flags.FlagChainID)
			if len(chainID) == 0 {
				return fmt.Errorf("chain ID must be specified")
			}

			_, err = keyring.New(sdk.KeyringServiceName(), viper.GetString(flags.FlagKeyringBackend), viper.GetString(flagClientHome), cmd.InOrStdin())
			if err != nil {
				return err
			}

			name := viper.GetString(flags.FlagName)
			key, err := clientCtx.Keyring.Key(name)
			if err != nil {
				return err
			}

			// Set flags for creating declare validator candidate tx
			viper.Set(flags.FlagHome, viper.GetString(flagClientHome))
			smbh.PrepareFlagsForTxCreateValidator(config, nodeID, chainID, valPubKey)

			// Fetch the amount of coins staked
			amount := viper.GetString(flagAmount)
			_, err = sdk.ParseCoinNormalized(amount)
			if err != nil {
				return err
			}
			inBuf := bufio.NewReader(cmd.InOrStdin())

			createValCfg, err := smbh.PrepareFlagsForTxCreateValidator(config, nodeID, chainID, valPubKey)
			if err != nil {
				return err
			}

			txFactory := tx.NewFactoryCLI(clientCtx, cmd.Flags())
			clientCtx = clientCtx.WithInput(inBuf).WithFromAddress(key.GetAddress())

			viper.Set(flags.FlagGenerateOnly, true)

			// create a 'create-validator' message
			txBldr, msg, err := smbh.BuildCreateValidatorMsg(clientCtx, createValCfg, txFactory)
			if err != nil {
				return err
			}

			if key.GetType() == keyring.TypeOffline || key.GetType() == keyring.TypeMulti {
				fmt.Println("Offline key passed in. Use `tx sign` command to sign:")
				return client2.PrintUnsignedStdTx(txBldr, clientCtx, []sdk.Msg{msg})
			}

			// write the unsigned transaction to the buffer
			w := bytes.NewBuffer([]byte{})
			clientCtx = clientCtx.WithOutput(w)

			if err = client2.PrintUnsignedStdTx(txFactory, clientCtx, []sdk.Msg{msg}); err != nil {
				return err
			}

			// read the transaction
			stdTx, err := readUnsignedGenTxFile(clientCtx, w)
			if err != nil {
				return err
			}

			// sign the transaction and write it to the output file
			// fixme (tx.Sign)
			txBuilder, err := clientCtx.TxConfig.WrapTxBuilder(stdTx)
			if err != nil {
				return err
			}

			err = client2.SignTx(txFactory, clientCtx, name, txBuilder, true, true)
			if err != nil {
				return errors.Wrap(err, "failed to sign std tx")
			}

			//signedTx, err := utils.SignStdTx(txBldr, clientCtx, name, stdTx, false, true)
			//if err != nil {
			//	return err
			//}

			txJSON, err := json.MarshalIndent(stdTx, "", "  ")
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
