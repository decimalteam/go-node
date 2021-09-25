package cli

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"

	"github.com/spf13/cobra"

	neturl "net/url"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/tendermint/tendermint/libs/cli"

	ncfg "bitbucket.org/decimalteam/go-node/config"
)

func UpdaterCmd(ctx *server.Context, cdc *codec.Codec, mbm module.BasicManager,
	defaultNodeHome string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "updater [url] [bech32_address]",
		Short: "Initialize updater file with URL and bech32",
		Long:  `Initialize updater file with 'http://example.com' and 'dx...'`,
		Args:  cobra.ExactArgs(2),
		RunE: func(_ *cobra.Command, args []string) error {
			url := args[0]
			addr := args[1]

			_, err := neturl.ParseRequestURI(url)
			if err != nil {
				return err
			}

			_, err = sdk.AccAddressFromBech32(addr)
			if err != nil {
				return err
			}

			return NewUpdateCFG(filepath.Join(ncfg.ConfigPath, ncfg.UpdaterName)).Push(url, addr)
		},
	}

	cmd.Flags().String(cli.HomeFlag, defaultNodeHome, "node's home directory")
	cmd.Flags().String("network", "", "mainnet, testnet or devnet")

	return cmd
}

type UpdateCFG struct {
	filename string
	URL      string `json:"url"`
	Address  string `json:"address"`
}

func NewUpdateCFG(filename string) *UpdateCFG {
	return &UpdateCFG{
		filename: filename,
	}
}

func (cfg *UpdateCFG) Push(url, address string) error {
	res, err := json.MarshalIndent(&UpdateCFG{
		filename: cfg.filename,
		URL:      url,
		Address:  address,
	}, "", "\t")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(cfg.filename, res, 0600)
	if err != nil {
		return err
	}
	return nil
}

func (cfg *UpdateCFG) Load() *UpdateCFG {
	res, err := ioutil.ReadFile(cfg.filename)
	if err != nil {
		return nil
	}
	err = json.Unmarshal(res, cfg)
	if err != nil {
		return nil
	}
	return cfg
}
