package cli

import (
	"encoding/json"
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/tendermint/tendermint/types"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	cfg "github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/libs/cli"
	tos "github.com/tendermint/tendermint/libs/os"
	trand "github.com/tendermint/tendermint/libs/rand"
	tmprototypes "github.com/tendermint/tendermint/proto/tendermint/types"
)

// InitCmd returns a command that initializes all files needed for Tendermint
// and the respective application
func InitCmd(ctx *server.Context, mbm module.BasicManager,
	defaultNodeHome string) *cobra.Command { // nolint: golint

	cmd := &cobra.Command{
		Use:   "init [moniker] --network mainnet|testnet|devnet",
		Short: "Initialize private validator, p2p, genesis, and application configuration files",
		Long:  `Initialize validators's and node's configuration files.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			cdc := clientCtx.JSONCodec
			config := ctx.Config

			config.SetRoot(viper.GetString(cli.HomeFlag))

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
					appState = []byte(testNetGenesis)
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
				appState, err = json.MarshalIndent(mbm.DefaultGenesis(cdc), "", " ")
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
				genDoc.ConsensusParams = &tmprototypes.ConsensusParams{
					Block: tmprototypes.BlockParams{
						MaxBytes:   10000000,
						MaxGas:     100000,
						TimeIotaMs: 1000,
					},
					Evidence: tmprototypes.EvidenceParams{
						MaxAgeNumBlocks: 100000,
						MaxAgeDuration:  86400000000000,
					},
					Validator: types.DefaultValidatorParams(),
				}
				if err = genutil.ExportGenesisFile(genDoc, genFile); err != nil {
					return err
				}
			}

			toPrint := newPrintInfo(config.Moniker, chainID, nodeID, "", appState)

			cfg.WriteConfigFile(filepath.Join(config.RootDir, "config", "config.toml"), config)
			return displayInfo(toPrint)
		},
	}

	cmd.Flags().String(cli.HomeFlag, defaultNodeHome, "node's home directory")
	cmd.Flags().BoolP(flagOverwrite, "o", false, "overwrite the genesis.json file")
	cmd.Flags().String(flags.FlagChainID, "", "genesis file chain-id, if left blank will be randomly created")
	cmd.Flags().String("network", "", "mainnet, testnet or devnet")

	return cmd
}
