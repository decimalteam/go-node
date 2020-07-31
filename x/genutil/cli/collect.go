package cli

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/pkg/errors"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	cfg "github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/libs/cli"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authexported "github.com/cosmos/cosmos-sdk/x/auth/exported"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/genutil/types"

	"bitbucket.org/decimalteam/go-node/x/genutil"
	"bitbucket.org/decimalteam/go-node/x/validator"
)

const flagGenTxDir = "gentx-dir"

// CollectGenTxsCmd - return the cobra command to collect genesis transactions
func CollectGenTxsCmd(ctx *server.Context, cdc *codec.Codec,
	genAccIterator types.GenesisAccountsIterator, defaultNodeHome string) *cobra.Command {

	cmd := &cobra.Command{
		Use:   "collect-gentxs",
		Short: "Collect genesis txs and output a genesis.json file",
		RunE: func(_ *cobra.Command, _ []string) error {
			config := ctx.Config
			config.SetRoot(viper.GetString(cli.HomeFlag))

			config.Mempool.CacheSize = 100000
			config.Mempool.Recheck = false
			config.Mempool.Size = 10000

			config.P2P.RecvRate = 15360000 // 15 mB/s
			config.P2P.SendRate = 15360000 // 15 mB/s
			config.P2P.FlushThrottleTimeout = 10 * time.Millisecond

			if viper.GetString("network") != "" {
				switch viper.GetString("network") {
				case "mainnet":
					config.P2P.Seeds = "0906b583daebe8951226e56cf75e1d2175f19671@decimal-node-1.mainnet.decimalchain.com:26656,1e9a5adb32f39a62849c94dbec95f251f5ebd728@decimal-node-2.mainnet.decimalchain.com:26656"
				case "testnet":
					config.P2P.Seeds = "bf7a6b366e3c451a3c12b3a6c01af7230fb92fc7@decimal-node-1.testnet.decimalchain.com:26656,76b81a4b817b39d63a3afe1f3a294f2a8f5c55b0@decimal-node-2.testnet.decimalchain.com:26656"
				case "devnet":
					config.P2P.Seeds = "8a2cc38f5264e9699abb8db91c9b4a4a061f000d@decimal-node-1.devnet.decimalchain.com:26656,27fcfef145b3717c5d639ec72fb12f9c43da98f0@decimal-node-2.devnet.decimalchain.com:26656"
				default:
					return fmt.Errorf("invalid network")
				}
			}

			name := viper.GetString(flags.FlagName)
			nodeID, valPubKey, err := genutil.InitializeNodeValidatorFiles(config)
			if err != nil {
				return err
			}

			genDoc, err := tmtypes.GenesisDocFromFile(config.GenesisFile())
			if err != nil {
				return err
			}

			genTxsDir := viper.GetString(flagGenTxDir)
			if genTxsDir == "" {
				genTxsDir = filepath.Join(config.RootDir, "config", "gentx")
			}

			toPrint := newPrintInfo(config.Moniker, genDoc.ChainID, nodeID, genTxsDir, json.RawMessage(""))
			initCfg := genutil.NewInitConfig(genDoc.ChainID, genTxsDir, name, nodeID, valPubKey)

			appMessage, err := GenAppStateFromConfig(cdc, config, initCfg, *genDoc, genAccIterator)
			if err != nil {
				return err
			}

			toPrint.AppMessage = appMessage

			// print out some key information
			return displayInfo(cdc, toPrint)
		},
	}

	cmd.Flags().String(cli.HomeFlag, defaultNodeHome, "node's home directory")
	cmd.Flags().String(flagGenTxDir, "",
		"override default \"gentx\" directory from which collect and execute "+
			"genesis transactions; default [--home]/config/gentx/")
	return cmd
}

type printInfo struct {
	Moniker    string          `json:"moniker" yaml:"moniker"`
	ChainID    string          `json:"chain_id" yaml:"chain_id"`
	NodeID     string          `json:"node_id" yaml:"node_id"`
	GenTxsDir  string          `json:"gentxs_dir" yaml:"gentxs_dir"`
	AppMessage json.RawMessage `json:"app_message" yaml:"app_message"`
}

func newPrintInfo(moniker, chainID, nodeID, genTxsDir string,
	appMessage json.RawMessage) printInfo {

	return printInfo{
		Moniker:    moniker,
		ChainID:    chainID,
		NodeID:     nodeID,
		GenTxsDir:  genTxsDir,
		AppMessage: appMessage,
	}
}

func displayInfo(cdc *codec.Codec, info printInfo) error {
	out, err := codec.MarshalJSONIndent(cdc, info)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(os.Stderr, "%s\n", string(sdk.MustSortJSON(out)))
	return err
}

// GenAppStateFromConfig gets the genesis app state from the config
func GenAppStateFromConfig(cdc *codec.Codec, config *cfg.Config,
	initCfg genutil.InitConfig, genDoc tmtypes.GenesisDoc,
	genAccIterator types.GenesisAccountsIterator,
) (appState json.RawMessage, err error) {

	// process genesis transactions, else create default genesis.json
	appGenTxs, err := CollectStdTxs(
		cdc, config.Moniker, initCfg.GenTxsDir, genDoc, genAccIterator)
	if err != nil {
		return appState, err
	}

	// if there are no gen txs to be processed, return the default empty state
	if len(appGenTxs) == 0 {
		return appState, errors.New("there must be at least one genesis tx")
	}

	// create the app state
	appGenesisState, err := genutil.GenesisStateFromGenDoc(cdc, genDoc)
	if err != nil {
		return appState, err
	}

	appGenesisState, err = genutil.SetGenTxsInAppGenesisState(cdc, appGenesisState, appGenTxs)
	if err != nil {
		return appState, err
	}
	appState, err = codec.MarshalJSONIndent(cdc, appGenesisState)
	if err != nil {
		return appState, err
	}

	genDoc.AppState = appState
	err = ExportGenesisFile(&genDoc, config.GenesisFile())
	return appState, err
}

// CollectStdTxs processes and validates application's genesis StdTxs and returns
// the list of appGenTxs, and persistent peers required to generate genesis.json.
func CollectStdTxs(cdc *codec.Codec, moniker, genTxsDir string,
	genDoc tmtypes.GenesisDoc, genAccIterator types.GenesisAccountsIterator,
) (appGenTxs []authtypes.StdTx, err error) {

	var fos []os.FileInfo
	fos, err = ioutil.ReadDir(genTxsDir)
	if err != nil {
		return appGenTxs, err
	}

	// prepare a map of all accounts in genesis state to then validate
	// against the validators addresses
	var appState map[string]json.RawMessage
	if err := cdc.UnmarshalJSON(genDoc.AppState, &appState); err != nil {
		return appGenTxs, err
	}

	addrMap := make(map[string]authexported.Account)
	genAccIterator.IterateGenesisAccounts(cdc, appState,
		func(acc authexported.Account) (stop bool) {
			addrMap[acc.GetAddress().String()] = acc
			return false
		},
	)

	for _, fo := range fos {
		filename := filepath.Join(genTxsDir, fo.Name())
		if !fo.IsDir() && (filepath.Ext(filename) != ".json") {
			continue
		}

		// get the genStdTx
		var jsonRawTx []byte
		if jsonRawTx, err = ioutil.ReadFile(filename); err != nil {
			return appGenTxs, err
		}
		var genStdTx authtypes.StdTx
		if err = cdc.UnmarshalJSON(jsonRawTx, &genStdTx); err != nil {
			return appGenTxs, err
		}
		appGenTxs = append(appGenTxs, genStdTx)

		// genesis transactions must be single-message
		msgs := genStdTx.GetMsgs()
		if len(msgs) != 1 {
			return appGenTxs, errors.New(
				"each genesis transaction must provide a single genesis message")
		}

		msg := msgs[0].(validator.MsgDeclareCandidate)
		// validate delegator and validator addresses and funds against the accounts in the state
		delAddr := sdk.AccAddress(msg.ValidatorAddr).String()
		valAddr := sdk.AccAddress(msg.ValidatorAddr).String()

		delAcc, delOk := addrMap[delAddr]
		if !delOk {
			return appGenTxs, fmt.Errorf(
				"account %v not in genesis.json: %+v", delAddr, addrMap)
		}

		_, valOk := addrMap[valAddr]
		if !valOk {
			return appGenTxs, fmt.Errorf(
				"account %v not in genesis.json: %+v", valAddr, addrMap)
		}

		if delAcc.GetCoins().AmountOf(msg.Stake.Denom).LT(msg.Stake.Amount) {
			return appGenTxs, fmt.Errorf(
				"insufficient fund for delegation %v: %v < %v",
				delAcc.GetAddress(), delAcc.GetCoins().AmountOf(msg.Stake.Denom), msg.Stake.Amount,
			)
		}
	}

	return appGenTxs, nil
}

// ExportGenesisFile creates and writes the genesis configuration to disk. An
// error is returned if building or writing the configuration to file fails.
func ExportGenesisFile(genDoc *tmtypes.GenesisDoc, genFile string) error {
	if err := genDoc.ValidateAndComplete(); err != nil {
		return err
	}

	return genDoc.SaveAs(genFile)
}
