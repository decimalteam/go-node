package cli

import (
	"fmt"
	"github.com/tendermint/tendermint/types"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/server"
	cfgApp "github.com/cosmos/cosmos-sdk/server/config"
	"github.com/cosmos/cosmos-sdk/store"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	cfg "github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/libs/cli"
	tos "github.com/tendermint/tendermint/libs/os"
	trand "github.com/tendermint/tendermint/libs/rand"
)

// InitCmd returns a command that initializes all files needed for Tendermint
// and the respective application
func InitCmd(ctx *server.Context, cdc *codec.Codec, mbm module.BasicManager,
	defaultNodeHome string) *cobra.Command { // nolint: golint
	cmd := &cobra.Command{
		Use:   "init [moniker] --network mainnet|testnet|devnet",
		Short: "Initialize private validator, p2p, genesis, and application configuration files",
		Long:  `Initialize validators's and node's configuration files.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			config := ctx.Config

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

			chainID := viper.GetString(flags.FlagChainID)
			if chainID == "" {
				chainID = fmt.Sprintf("test-chain-%v", trand.Str(6))
			}

			nodeID, _, err := genutil.InitializeNodeValidatorFiles(config)
			if err != nil {
				return err
			}

			config.Moniker = args[0]

			genFile := config.GenesisFile()
			if !viper.GetBool(flagOverwrite) && tos.FileExists(genFile) {
				return fmt.Errorf("genesis.json file already exists: %v", genFile)
			}

			var appState []byte
			if viper.GetString("network") != "" {
				switch viper.GetString("network") {
				case "mainnet":
					appState = []byte(mainNetGenesis)
				case "testnet":
					appState = []byte(TestNetGenesis)
				case "devnet":
					appState = []byte(devNetGenesis)
				default:
					return fmt.Errorf("invalid network")
				}

				genDoc, err := types.GenesisDocFromJSON(appState)
				if err != nil {
					return err
				}
				if err = genutil.ExportGenesisFile(genDoc, genFile); err != nil {
					return err
				}

			} else {
				appState, err = codec.MarshalJSONIndent(cdc, mbm.DefaultGenesis())
				if err != nil {
					return err
				}

				genDoc := &types.GenesisDoc{}
				if _, err := os.Stat(genFile); err != nil {
					if !os.IsNotExist(err) {
						return err
					}
				} else {
					genDoc, err = types.GenesisDocFromFile(genFile)
					if err != nil {
						return err
					}
				}

				genDoc.ChainID = chainID
				genDoc.Validators = nil
				genDoc.AppState = appState
				genDoc.ConsensusParams = &types.ConsensusParams{
					Block: types.BlockParams{
						MaxBytes:   10000000,
						MaxGas:     100000,
						TimeIotaMs: 1000,
					},
					Evidence: types.EvidenceParams{
						MaxAgeNumBlocks: 100000,
						MaxAgeDuration:  86400000000000,
					},
					Validator: types.DefaultValidatorParams(),
				}
				if err = genutil.ExportGenesisFile(genDoc, genFile); err != nil {
					return err
				}
			}

			// Set pruning from 'syncable' to 'nothing'
			appConfigFilePath := filepath.Join(config.RootDir, "config", "app.toml")
			appConf, _ := cfgApp.ParseConfig()
			appConf.Pruning = store.PruningStrategyNothing
			cfgApp.WriteConfigFile(appConfigFilePath, appConf)

			toPrint := newPrintInfo(config.Moniker, chainID, nodeID, "", appState)
			cfg.WriteConfigFile(filepath.Join(config.RootDir, "config", "config.toml"), config)

			return displayInfo(cdc, toPrint)
		},
	}

	cmd.Flags().String(cli.HomeFlag, defaultNodeHome, "node's home directory")
	cmd.Flags().BoolP(flagOverwrite, "o", false, "overwrite the genesis.json file")
	cmd.Flags().String(flags.FlagChainID, "", "genesis file chain-id, if left blank will be randomly created")
	cmd.Flags().String("network", "", "mainnet, testnet or devnet")

	return cmd
}
