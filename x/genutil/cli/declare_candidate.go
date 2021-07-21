package cli

import (
	"bitbucket.org/decimalteam/go-node/x/validator/client/cli"
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	authclient "github.com/cosmos/cosmos-sdk/x/auth/client"
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
	fsCreateValidator, defaultsDesc := smbh.CreateValidatorMsgHelpers(ipDefault)

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
			serverCtx := server.GetServerContextFromCmd(cmd)

			config := serverCtx.Config
			config.SetRoot(clientCtx.HomeDir)

			nodeID, valPubKey, err := genutil.InitializeNodeValidatorFiles(ctx.Config)
			if err != nil {
				return err
			}

			// Read --nodeID, if empty take it from priv_validator.json
			if nodeIDString, _ := cmd.Flags().GetString(cli.FlagNodeID); nodeIDString != "" {
				nodeID = nodeIDString
			}
			// Read --pubkey, if empty take it from priv_validator.json
			if valPubKeyString, _ := cmd.Flags().GetString(cli.FlagPubKey); valPubKeyString != "" {
				valPubKey, err = sdk.GetPubKeyFromBech32(sdk.Bech32PubKeyTypeConsPub, valPubKeyString)
				if err != nil {
					return err
				}
			}

			// Read chain ID
			chainID, _ := cmd.Flags().GetString(flags.FlagChainID)
			if len(chainID) == 0 {
				return fmt.Errorf("chain ID must be specified")
			}

			name, _ := cmd.Flags().GetString(flags.FlagName)
			key, err := clientCtx.Keyring.Key(name)
			if err != nil {
				return err
			}

			// Set flags for creating declare validator candidate tx
			cliHome, _ := cmd.Flags().GetString(flagClientHome)
			cmd.Flags().Set(flags.FlagHome, cliHome)

			moniker := config.Moniker
			if m, _ := cmd.Flags().GetString(cli.FlagMoniker); m != "" {
				moniker = m
			}

			createValCfg, _ := smbh.PrepareFlagsForTxCreateValidator(cmd.Flags(), moniker, nodeID, chainID, valPubKey)

			// Fetch the amount of coins staked
			amount, _ := cmd.Flags().GetString(cli.FlagAmount)
			_, err = sdk.ParseCoinNormalized(amount)
			if err != nil {
				return err
			}
			inBuf := bufio.NewReader(cmd.InOrStdin())

			txFactory := tx.NewFactoryCLI(clientCtx, cmd.Flags()).WithTxConfig(clientCtx.TxConfig).WithAccountRetriever(clientCtx.AccountRetriever)
			clientCtx = clientCtx.WithInput(inBuf).WithFromAddress(key.GetAddress())

			viper.Set(flags.FlagGenerateOnly, true)

			// create a 'create-validator' message
			txBldr, msg, err := smbh.BuildCreateValidatorMsg(clientCtx, createValCfg, txFactory, cmd.Flags(), true)
			if err != nil {
				return err
			}

			if key.GetType() == keyring.TypeOffline || key.GetType() == keyring.TypeMulti {
				fmt.Println("Offline key passed in. Use `tx sign` command to sign:")
				return authclient.PrintUnsignedStdTx(txBldr, clientCtx, []sdk.Msg{msg})
			}

			// write the unsigned transaction to the buffer
			w := bytes.NewBuffer([]byte{})
			clientCtx = clientCtx.WithOutput(w)

			if err = authclient.PrintUnsignedStdTx(txFactory, clientCtx, []sdk.Msg{msg}); err != nil {
				return err
			}

			// read the transaction
			stdTx, err := readUnsignedGenTxFile(clientCtx, w)
			if err != nil {
				return err
			}

			// sign the transaction and write it to the output file
			fmt.Printf("1\n")
			txBuilder, err := clientCtx.TxConfig.WrapTxBuilder(stdTx)
			if err != nil {
				return err
			}

			fmt.Printf("2\n")
			err = authclient.SignTx(txFactory, clientCtx, name, txBuilder, true, true)
			if err != nil {
				return errors.Wrap(err, "failed to sign std tx")
			}

			fmt.Printf("3\n")
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
