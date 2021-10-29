package main

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/server"
	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"
	"os"
	"path"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/tendermint/tendermint/libs/cli"

	"github.com/cosmos/cosmos-sdk/codec"

	cliconfig "github.com/cosmos/cosmos-sdk/client/config"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/client/rpc"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	authcmd "github.com/cosmos/cosmos-sdk/x/auth/client/cli"
	bankcmd "github.com/cosmos/cosmos-sdk/x/bank/client/cli"

	"bitbucket.org/decimalteam/go-node/app"
	"bitbucket.org/decimalteam/go-node/config"
	coinCmd "bitbucket.org/decimalteam/go-node/x/coin/client/cli"
	multisigCmd "bitbucket.org/decimalteam/go-node/x/multisig/client/cli"
)

func main() {
	cobra.EnableCommandSorting = false

	encodingConfig := app.MakeEncodingConfig()
	/*initClientCtx := client.Context{}.
		WithJSONMarshaler(encodingConfig.Codec).
		WithInterfaceRegistry(encodingConfig.InterfaceRegistry).
		WithTxConfig(encodingConfig.TxConfig).
		WithLegacyAmino(encodingConfig.Amino).
		WithInput(os.Stdin).
		WithAccountRetriever(authtypes.AccountRetriever{}).
		WithHomeDir(app.DefaultNodeHome).
		WithViper("AU")*/

	cdc := encodingConfig.Amino

	// Read in the configuration file for the sdk
	_config := sdk.GetConfig()
	_config.SetCoinType(60)
	_config.SetFullFundraiserPath("44'/60'/0'/0/0")
	_config.SetBech32PrefixForAccount(config.DecimalPrefixAccAddr, config.DecimalPrefixAccPub)
	_config.SetBech32PrefixForValidator(config.DecimalPrefixValAddr, config.DecimalPrefixValPub)
	_config.SetBech32PrefixForConsensusNode(config.DecimalPrefixConsAddr, config.DecimalPrefixConsPub)
	_config.Seal()

	rootCmd := &cobra.Command{
		Use:   "deccli",
		Short: "Decimal Client Console",
	}

	// Add --chain-id to persistent flags and mark it required
	rootCmd.PersistentFlags().String(flags.FlagChainID, "", "Chain ID of decimal node")
	rootCmd.Flags().String(cli.HomeFlag, app.DefaultCLIHome, "node's home directory")

	rootCmd.PersistentPreRunE = func(cmd *cobra.Command, _ []string) error {
		// set the default command outputs
		/*cmd.SetOut(cmd.OutOrStdout())
		cmd.SetErr(cmd.ErrOrStderr())*/

		home, err := cmd.PersistentFlags().GetString(cli.HomeFlag)
		if err != nil {
			return err
		}

		cfgFile := path.Join(home, "config", "config.toml")
		if _, err := os.Stat(cfgFile); err == nil {
			viper.SetConfigFile(cfgFile)

			if err := viper.ReadInConfig(); err != nil {
				return err
			}
		}
		if err := viper.BindPFlag(flags.FlagChainID, cmd.PersistentFlags().Lookup(flags.FlagChainID)); err != nil {
			return err
		}
		if err := viper.BindPFlag(cli.EncodingFlag, cmd.PersistentFlags().Lookup(cli.EncodingFlag)); err != nil {
			return err
		}
		return viper.BindPFlag(cli.OutputFlag, cmd.PersistentFlags().Lookup(cli.OutputFlag))

	/*	initClientCtx = client.ReadHomeFlag(initClientCtx, cmd)

		initClientCtx, err := cliconfig.ReadFromClientConfig(initClientCtx)
		if err != nil {
			return err
		}

		if err := client.SetCmdClientContextHandler(initClientCtx, cmd); err != nil {
			return err
		}

		return server.InterceptConfigsPreRunHandler(cmd)*/
	}

	// Construct Root Command
	rootCmd.AddCommand(
		rpc.StatusCommand(),
		cliconfig.Cmd(),
		queryCmd(cdc),
		txCmd(cdc),
		flags.LineBreak,
		flags.LineBreak,
		flags.LineBreak,
		version.NewVersionCommand(),
		cli.NewCompletionCmd(rootCmd, true),
		keys.Commands(app.DefaultNodeHome),
	)

	if err := svrcmd.Execute(rootCmd, app.DefaultNodeHome); err != nil {
		switch e := err.(type) {
		case server.ErrorCode:
			os.Exit(e.Code)

		default:
			os.Exit(1)
		}
	}
}

func queryCmd(cdc *codec.LegacyAmino) *cobra.Command {
	queryCmd := &cobra.Command{
		Use:     "query",
		Aliases: []string{"q"},
		Short:   "Querying subcommands",
		RunE:    client.ValidateCmd,
	}

	queryCmd.AddCommand(
		authcmd.GetAccountCmd(),
		flags.LineBreak,
		rpc.ValidatorCommand(),
		rpc.BlockCommand(),
		authcmd.QueryTxsByEventsCmd(),
		authcmd.QueryTxCmd(),
		flags.LineBreak,
	)

	// add modules' query commands
	app.ModuleBasics.AddQueryCommands(coinCmd.GetQueryCmd("coin", cdc))
	app.ModuleBasics.AddQueryCommands(multisigCmd.GetQueryCmd("multisig", cdc))
	app.ModuleBasics.AddQueryCommands(queryCmd)

	queryCmd.Flags().String(flags.FlagChainID, "", "The network chain ID")

	return queryCmd
}

func txCmd(cdc *codec.LegacyAmino) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tx",
		Short: "Transactions subcommands",
	}

	cmd.AddCommand(
		bankcmd.NewSendTxCmd(),
		flags.LineBreak,
		authcmd.GetSignCommand(),
		authcmd.GetMultiSignCommand(),
		flags.LineBreak,
		authcmd.GetBroadcastCommand(),
		authcmd.GetEncodeCommand(),
		flags.LineBreak,
	)

	// add modules' tx commands
	app.ModuleBasics.AddTxCommands(cmd)

	cmd.PersistentFlags().String(flags.FlagChainID, "", "The network chain ID")

	return cmd
}

func initConfig(cmd *cobra.Command) error {
	home, err := cmd.PersistentFlags().GetString(cli.HomeFlag)
	if err != nil {
		return err
	}

	cfgFile := path.Join(home, "config", "config.toml")
	if _, err := os.Stat(cfgFile); err == nil {
		viper.SetConfigFile(cfgFile)

		if err := viper.ReadInConfig(); err != nil {
			return err
		}
	}
	return nil

	if err := viper.BindPFlag(flags.FlagChainID, cmd.PersistentFlags().Lookup(flags.FlagChainID)); err != nil {
		return err
	}
	if err := viper.BindPFlag(cli.EncodingFlag, cmd.PersistentFlags().Lookup(cli.EncodingFlag)); err != nil {
		return err
	}
	return viper.BindPFlag(cli.OutputFlag, cmd.PersistentFlags().Lookup(cli.OutputFlag))
}
