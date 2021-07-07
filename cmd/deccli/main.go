package main

import (
	"fmt"
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

	cdc := app.MakeAminoCodec()

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
	rootCmd.PersistentPreRunE = func(_ *cobra.Command, _ []string) error {
		return initConfig(rootCmd)
	}

	// Construct Root Command
	rootCmd.AddCommand(
		rpc.StatusCommand(),
		cliconfig.Cmd(),
		queryCmd(cdc),
		txCmd(cdc),
		flags.LineBreak,
		//lcd.ServeCommand(cdc, registerRoutes),
		flags.LineBreak,
		keys.Commands(app.DefaultNodeHome),
		flags.LineBreak,
		version.NewVersionCommand(),
		cli.NewCompletionCmd(rootCmd, true),
	)

	executor := cli.PrepareMainCmd(rootCmd, "AU", app.DefaultCLIHome)
	err := executor.Execute()
	if err != nil {
		fmt.Printf("Failed executing CLI command: %s, exiting...\n", err)
		os.Exit(1)
	}
}

//func registerRoutes(rs *lcd.RestServer) {
//	client.RegisterRoutes(rs.CliCtx, rs.Mux)
//	app.ModuleBasics.RegisterRESTRoutes(rs.CliCtx, rs.Mux)
//	authrest.RegisterTxRoutes(rs.CliCtx, rs.Mux)
//}

func queryCmd(cdc *codec.LegacyAmino) *cobra.Command {
	queryCmd := &cobra.Command{
		Use:     "query",
		Aliases: []string{"q"},
		Short:   "Querying subcommands",
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

	return queryCmd
}

func txCmd(cdc *codec.LegacyAmino) *cobra.Command {
	txCmd := &cobra.Command{
		Use:   "tx",
		Short: "Transactions subcommands",
	}

	txCmd.AddCommand(
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
	app.ModuleBasics.AddTxCommands(coinCmd.GetTxCmd(cdc))
	app.ModuleBasics.AddTxCommands(txCmd)

	return txCmd
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
	if err := viper.BindPFlag(flags.FlagChainID, cmd.PersistentFlags().Lookup(flags.FlagChainID)); err != nil {
		return err
	}
	if err := viper.BindPFlag(cli.EncodingFlag, cmd.PersistentFlags().Lookup(cli.EncodingFlag)); err != nil {
		return err
	}
	return viper.BindPFlag(cli.OutputFlag, cmd.PersistentFlags().Lookup(cli.OutputFlag))
}
