package keys

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/libs/cli"

	"github.com/cosmos/cosmos-sdk/client/flags"
)

// Commands registers a sub-tree of commands to interact with
// local private key storage.
func Commands(defaultNodeHome string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "keys",
		Short: "Add or view local private keys",
		Long: `Keys allows you to manage your local keystore for tendermint.

    These keys may be in any format supported by go-crypto and can be
    used by light-clients, full nodes, or any other application that
    needs to sign with a private key.`,
	}
	cmd.AddCommand(
		MnemonicKeyCommand(),
		AddKeyCommand(),
		ExportKeyCommand(),
		ImportKeyCommand(),
		ListKeysCmd(),
		ShowKeysCmd(),
		flags.LineBreak,
		DeleteKeyCommand(),
		//UpdateKeyCommand(),
		ParseKeyStringCommand(),
		MigrateCommand(),
	)
	cmd.PersistentFlags().String(flags.FlagKeyringBackend, flags.DefaultKeyringBackend, "Select keyring's backend (os|file|test)")
	cmd.PersistentFlags().String(flags.FlagHome, defaultNodeHome, "The application home directory")
	cmd.PersistentFlags().String(flags.FlagKeyringDir, "", "The client Keyring directory; if omitted, the default 'home' directory will be used")
	cmd.PersistentFlags().String(cli.OutputFlag, "text", "Output format (text|json)")

	viper.BindPFlag(flags.FlagKeyringBackend, cmd.Flags().Lookup(flags.FlagKeyringBackend))
	return cmd
}
