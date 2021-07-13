package cli

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/types/errors"
	client2 "github.com/cosmos/cosmos-sdk/x/auth/client"
	"github.com/cosmos/cosmos-sdk/x/staking/client/cli"

	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"

	cfg "github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/crypto"
	tos "github.com/tendermint/tendermint/libs/os"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/genutil/types"

	"bitbucket.org/decimalteam/go-node/x/genutil"
)

const (
	flagOverwrite  = "overwrite"
	flagClientHome = "home-client"
)

// StakingMsgBuildingHelpers helpers for message building gen-tx command
type StakingMsgBuildingHelpers interface {
	CreateValidatorMsgHelpers(ipDefault string) (fs *flag.FlagSet, nodeIDFlag, pubkeyFlag, amountFlag, defaultsDesc string)
	PrepareFlagsForTxCreateValidator(config *cfg.Config, nodeID, chainID string, valPubKey crypto.PubKey) (cli.TxCreateValidatorConfig, error)
	BuildCreateValidatorMsg(cliCtx client.Context, config cli.TxCreateValidatorConfig, txBldr tx.Factory) (tx.Factory, sdk.Msg, error)
}

// GenTxCmd builds the application's gentx command.
// nolint: errcheck
func GenTxCmd(ctx *server.Context, txEncodingConfig client.TxConfig, mbm module.BasicManager, smbh StakingMsgBuildingHelpers,
	genBalIterator types.GenesisBalancesIterator, defaultNodeHome, defaultCLIHome string) *cobra.Command {

	ipDefault, _ := server.ExternalIP()
	fsCreateValidator, flagNodeID, flagPubKey, flagAmount, defaultsDesc := smbh.CreateValidatorMsgHelpers(ipDefault)

	cmd := &cobra.Command{
		Use:   "gentx",
		Short: "Generate a genesis tx carrying a self delegation",
		Args:  cobra.NoArgs,
		Long: fmt.Sprintf(`This command is an alias of the 'tx create-validator' command'.

		It creates a genesis transaction to create a validator. 
		The following default parameters are included: 
		    %s`, defaultsDesc),

		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			serverCtx := server.GetServerContextFromCmd(cmd)

			cdc := clientCtx.Codec
			amino := clientCtx.LegacyAmino

			config := serverCtx.Config

			config.SetRoot(clientCtx.HomeDir)

			if err != nil {
				return err
			}

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

			genDoc, err := tmtypes.GenesisDocFromFile(config.GenesisFile())
			if err != nil {
				return err
			}

			var genesisState map[string]json.RawMessage
			if err = json.Unmarshal(genDoc.AppState, &genesisState); err != nil {
				return err
			}

			// LABEL-TEST: check this part of code
			if err = mbm.ValidateGenesis(cdc, txEncodingConfig, genesisState); err != nil {
				return err
			}

			inBuf := bufio.NewReader(cmd.InOrStdin())

			//kb, err := keyring.New(sdk.KeyringServiceName(), viper.GetString(flags.FlagKeyringBackend), viper.GetString(flagClientHome), cmd.InOrStdin())
			//if err != nil {
			//	return err
			//}

			name := viper.GetString(flags.FlagName)
			key, err := clientCtx.Keyring.Key(name)
			if err != nil {
				return err
			}

			// Set flags for creating gentx
			viper.Set(flags.FlagHome, viper.GetString(flagClientHome))

			createValCfg, err := smbh.PrepareFlagsForTxCreateValidator(config, nodeID, genDoc.ChainID, valPubKey)
			if err != nil {
				return err
			}

			// Fetch the amount of coins staked
			amount := viper.GetString(flagAmount)
			coins, err := sdk.ParseCoinsNormalized(amount)
			if err != nil {
				return err
			}

			err = genutil.ValidateAccountInGenesis(genesisState, genBalIterator, key.GetAddress(), coins, amino, cdc)
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
			log.Println(msg)

			//keys.NewInMemoryKeyBase()
			//info, err := txBldr.Keybase().Get(name)
			//if err != nil {
			//	return err
			//}

			if key.GetType() == keyring.TypeOffline || key.GetType() == keyring.TypeMulti {
				fmt.Println("Offline key passed in. Use `tx sign` command to sign:")
				return client2.PrintUnsignedStdTx(txBldr, clientCtx, []sdk.Msg{msg})
			}

			// write the unsigned transaction to the buffer
			w := bytes.NewBuffer([]byte{})
			clientCtx = clientCtx.WithOutput(w)

			if err = client2.PrintUnsignedStdTx(txBldr, clientCtx, []sdk.Msg{msg}); err != nil {
				return err
			}

			// read the transaction
			stdTx, err := readUnsignedGenTxFile(clientCtx, w)
			if err != nil {
				return err
			}

			// sign the transaction and write it to the output file
			txBuilder, err := clientCtx.TxConfig.WrapTxBuilder(stdTx)
			if err != nil {
				return err
			}

			err = client2.SignTx(txFactory, clientCtx, name, txBuilder, true, true)
			if err != nil {
				return errors.Wrap(err, "failed to sign std tx")
			}

			// Fetch output file name
			outputDocument := viper.GetString(flags.FlagOutputDocument)
			if outputDocument == "" {
				outputDocument, err = makeOutputFilepath(config.RootDir, nodeID)
				if err != nil {
					return err
				}
			}

			if err := writeSignedGenTx(clientCtx, outputDocument, stdTx); err != nil {
				return err
			}

			fmt.Println("Genesis transaction written to ", outputDocument)
			return nil

		},
	}

	cmd.Flags().String(flags.FlagHome, defaultNodeHome, "node's home directory")
	cmd.Flags().String(flagClientHome, defaultCLIHome, "client's home directory")
	cmd.Flags().String(flags.FlagName, "", "name of private key with which to sign the gentx")
	cmd.Flags().String(flags.FlagOutputDocument, "",
		"write the genesis transaction JSON document to the given file instead of the default location")
	cmd.Flags().AddFlagSet(fsCreateValidator)

	cmd.MarkFlagRequired(flags.FlagName)
	return cmd
}

// func printJSON(data interface{}) {
// 	x, err := json.MarshalIndent(data, "", "\t")
// 	if err != nil {
// 		panic(err)
// 	}
// 	fmt.Println(string(x))
// }

func makeOutputFilepath(rootDir, nodeID string) (string, error) {
	writePath := filepath.Join(rootDir, "config", "gentx")
	if err := tos.EnsureDir(writePath, 0700); err != nil {
		return "", err
	}
	return filepath.Join(writePath, fmt.Sprintf("gentx-%v.json", nodeID)), nil
}

func readUnsignedGenTxFile(clientCtx client.Context, r io.Reader) (sdk.Tx, error) {
	bytes, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	aTx, err := clientCtx.TxConfig.TxJSONDecoder()(bytes)
	if err != nil {
		return nil, err
	}

	return aTx, err
}

func writeSignedGenTx(clientCtx client.Context, outputDocument string, tx sdk.Tx) error {
	outputFile, err := os.OpenFile(outputDocument, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer outputFile.Close()
	json, err := clientCtx.TxConfig.TxJSONEncoder()(tx)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(outputFile, "%s\n", json)

	return err
}
