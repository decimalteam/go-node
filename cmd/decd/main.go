package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/server/cmd"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/cosmos/cosmos-sdk/store/types"
	types3 "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	bankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	"github.com/spf13/cast"

	"bitbucket.org/decimalteam/go-node/utils/keys"
	"github.com/pkg/errors"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/tendermint/tendermint/libs/cli"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdkCfg "github.com/cosmos/cosmos-sdk/client/config"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	"bitbucket.org/decimalteam/go-node/app"
	"bitbucket.org/decimalteam/go-node/config"
	"bitbucket.org/decimalteam/go-node/x/genutil"
	genutilcli "bitbucket.org/decimalteam/go-node/x/genutil/cli"
	"bitbucket.org/decimalteam/go-node/x/validator"
)

const flagInvCheckPeriod = "inv-check-period"

var invCheckPeriod uint

func main() {
	encodingConfig := app.MakeEncodingConfig()
	initClientCtx := client.Context{}.
		WithJSONMarshaler(encodingConfig.Codec).
		WithInterfaceRegistry(encodingConfig.InterfaceRegistry).
		WithTxConfig(encodingConfig.TxConfig).
		WithLegacyAmino(encodingConfig.Amino).
		WithInput(os.Stdin).
		WithAccountRetriever(authtypes.AccountRetriever{}).
		WithHomeDir(app.DefaultNodeHome).
		WithViper("")

	_config := sdk.GetConfig()
	_config.SetCoinType(60)
	_config.SetFullFundraiserPath("44'/60'/0'/0/0")
	_config.SetBech32PrefixForAccount(config.DecimalPrefixAccAddr, config.DecimalPrefixAccPub)
	_config.SetBech32PrefixForValidator(config.DecimalPrefixValAddr, config.DecimalPrefixValPub)
	_config.SetBech32PrefixForConsensusNode(config.DecimalPrefixConsAddr, config.DecimalPrefixConsPub)
	_config.Seal()

	ctx := server.NewDefaultContext()
	cobra.EnableCommandSorting = false
	rootCmd := &cobra.Command{
		Use:   "decd",
		Short: "Decimal Go Node",
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			// set the default command outputs
			cmd.SetOut(cmd.OutOrStdout())
			cmd.SetErr(cmd.ErrOrStderr())

			initClientCtx = client.ReadHomeFlag(initClientCtx, cmd)

			initClientCtx, err := sdkCfg.ReadFromClientConfig(initClientCtx)
			if err != nil {
				return err
			}

			if err := client.SetCmdClientContextHandler(initClientCtx, cmd); err != nil {
				return err
			}

			// LABEL-TEST: added line for test
			serverCtx := server.GetServerContextFromCmd(cmd)
			config := serverCtx.Config
			config.SetRoot(initClientCtx.HomeDir)

			return server.InterceptConfigsPreRunHandler(cmd)
		},
	}

	rootCmd.AddCommand(
		genutilcli.InitCmd(ctx, app.ModuleBasics, app.DefaultNodeHome),
		genutilcli.CollectGenTxsCmd(ctx, bankTypes.GenesisBalancesIterator{}, app.DefaultNodeHome),
		genutilcli.GenTxCmd(
			ctx, encodingConfig.TxConfig, app.ModuleBasics, validator.AppModuleBasic{},
			bankTypes.GenesisBalancesIterator{}, app.DefaultNodeHome, app.DefaultCLIHome,
		),
		app.MigrateGenesisCmd(),
		genutilcli.GenDeclareCandidateTxCmd(
			ctx, app.ModuleBasics, validator.AppModuleBasic{}, app.DefaultNodeHome, app.DefaultCLIHome,
		),
		genutilcli.ValidateGenesisCmd(ctx, app.ModuleBasics),
		// AddGenesisAccountCmd allows users to add accounts to the genesis file
		addGenesisAccountCmd(ctx, app.DefaultNodeHome, app.DefaultCLIHome),
		fixAppHashError(ctx, app.DefaultNodeHome),
	)
	server.AddCommands(rootCmd, app.DefaultNodeHome, newApp, exportAppStateAndTMValidators, func(cmd *cobra.Command) {})

	// prepare and add flags
	//rootCmd.PrsistentFlags().UintVar(&invCheckPeriod, flagInvCheckPeriod,
	//	0, "Assert registered invariants every N blocks")

	if err := cmd.Execute(rootCmd, app.DefaultNodeHome); err != nil {
		switch e := err.(type) {
		case server.ErrorCode:
			os.Exit(e.Code)

		default:
			os.Exit(1)
		}
	}
}

func newApp(logger log.Logger, db dbm.DB, traceStore io.Writer, options servertypes.AppOptions) servertypes.Application {
	skipUpgradesHeight := map[int64]bool{}

	for _, h := range cast.ToIntSlice(options.Get(server.FlagUnsafeSkipUpgrades)) {
		skipUpgradesHeight[int64(h)] = true
	}

	return app.NewInitApp(
		logger, db, traceStore, true, skipUpgradesHeight, invCheckPeriod,
		baseapp.SetPruning(types.NewPruningOptionsFromString(viper.GetString("pruning"))),
		baseapp.SetMinGasPrices(viper.GetString(server.FlagMinGasPrices)),
		baseapp.SetHaltHeight(uint64(viper.GetInt(server.FlagHaltHeight))),
	)
}

func exportAppStateAndTMValidators(
	logger log.Logger, db dbm.DB, traceStore io.Writer, height int64, forZeroHeight bool, jailWhiteList []string, options servertypes.AppOptions,
) (servertypes.ExportedApp, error) {

	if height != -1 {
		aApp := app.NewInitApp(logger, db, traceStore, true, map[int64]bool{}, uint(1))
		err := aApp.LoadHeight(height)
		if err != nil {
			return servertypes.ExportedApp{}, err
		}

		return aApp.ExportAppStateAndValidators(forZeroHeight, jailWhiteList)
	}

	aApp := app.NewInitApp(logger, db, traceStore, true, map[int64]bool{}, uint(1))

	return aApp.ExportAppStateAndValidators(forZeroHeight, jailWhiteList)
}

// NOTE: Following part of the code was copied from file:
// github.com/cosmos/cosmos-sdk@v0.37.4/x/genaccounts/client/cli/genesis_accts.go
// since it was removed in the latest cosmos-sdk release.

const (
	flagClientHome   = "home-client"
	flagVestingStart = "vesting-start-time"
	flagVestingEnd   = "vesting-end-time"
	flagVestingAmt   = "vesting-amount"
)

// addGenesisAccountCmd returns add-genesis-account cobra Command.
func addGenesisAccountCmd(ctx *server.Context,
	defaultNodeHome, defaultClientHome string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-genesis-account [address_or_key_name] [coin][,[coin]]",
		Short: "Add genesis account to genesis.json",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			depCdc := clientCtx.JSONMarshaler
			cdc := depCdc.(codec.Marshaler)

			serverCtx := server.GetServerContextFromCmd(cmd)
			config := serverCtx.Config

			config.SetRoot(clientCtx.HomeDir)

			addr, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				kb, err := keys.NewKeyBaseFromDir(viper.GetString(flagClientHome))
				if err != nil {
					return err
				}

				// todo
				_, err = kb.Export(args[0])
				if err != nil {
					return err
				}

				//addr = info.GetAddress()
			}

			coins, err := sdk.ParseCoinsNormalized(args[1])
			if err != nil {
				return err
			}

			vestingStart := viper.GetInt64(flagVestingStart)
			vestingEnd := viper.GetInt64(flagVestingEnd)
			vestingAmt, err := sdk.ParseCoinsNormalized(viper.GetString(flagVestingAmt))
			if err != nil {
				return err
			}

			// create concrete account type based on input parameters
			var genAccount authtypes.GenesisAccount

			balances := bankTypes.Balance{
				Address: addr.String(),
				Coins:   coins.Sort(),
			}
			baseAccount := authtypes.NewBaseAccount(addr, nil, 0, 0)

			if !vestingAmt.IsZero() {
				baseVestingAccount := types3.NewBaseVestingAccount(
					baseAccount, vestingAmt.Sort(), vestingEnd,
				)
				if err != nil {
					return err
				}

				switch {
				case vestingStart != 0 && vestingEnd != 0:
					genAccount = types3.NewContinuousVestingAccountRaw(baseVestingAccount, vestingStart)

				case vestingEnd != 0:
					genAccount = types3.NewDelayedVestingAccountRaw(baseVestingAccount)

				default:
					return errors.New("invalid vesting parameters; must supply start and end time or end time")
				}
			} else {
				genAccount = baseAccount
			}

			if err := genAccount.Validate(); err != nil {
				return fmt.Errorf("failed to validate new genesis account: %w", err)
			}

			// retrieve the app state
			genFile := config.GenesisFile()
			appState, genDoc, err := genutil.GenesisStateFromGenFile(cdc, genFile)
			if err != nil {
				return err
			}

			authGenState := authtypes.GetGenesisStateFromAppState(cdc, appState)

			accs, err := authtypes.UnpackAccounts(authGenState.Accounts)
			if err != nil {
				return fmt.Errorf("failed to get accounts from any: %w", err)
			}

			if accs.Contains(addr) {
				return fmt.Errorf("cannot add account at existing address %s", addr)
			}

			// Add the new account to the set of genesis accounts and sanitize the
			// accounts afterwards.
			accs = append(accs, genAccount)

			accs = authtypes.SanitizeGenesisAccounts(accs)

			genAccs, err := authtypes.PackAccounts(accs)
			if err != nil {
				return fmt.Errorf("failed to convert accounts into any's: %w", err)
			}
			authGenState.Accounts = genAccs

			authGenStateBz, err := cdc.MarshalJSON(&authGenState)
			if err != nil {
				return fmt.Errorf("failed to marshal auth genesis state: %w", err)
			}
			appState[authtypes.ModuleName] = authGenStateBz

			bankGetState := bankTypes.GetGenesisStateFromAppState(depCdc, appState)
			bankGetState.Balances = append(bankGetState.Balances, balances)
			bankGetState.Balances = bankTypes.SanitizeGenesisBalances(bankGetState.Balances)
			bankGetState.Supply = bankGetState.Supply.Add(balances.Coins...)

			bankGenStateBz, err := cdc.MarshalJSON(bankGetState)
			if err != nil {
				return fmt.Errorf("falied to marshal bank genesis state: %w", err)
			}

			appState[bankTypes.ModuleName] = bankGenStateBz

			appStateJSON, err := json.Marshal(appState)
			if err != nil {
				return err
			}

			// export app state
			genDoc.AppState = appStateJSON

			return genutilcli.ExportGenesisFile(genDoc, genFile)
		},
	}

	cmd.Flags().String(cli.HomeFlag, defaultNodeHome, "node's home directory")
	cmd.Flags().String(flagClientHome, defaultClientHome, "client's home directory")
	cmd.Flags().String(flagVestingAmt, "", "amount of coins for vesting accounts")
	cmd.Flags().Uint64(flagVestingStart, 0, "schedule start time (unix epoch) for vesting accounts")
	cmd.Flags().Uint64(flagVestingEnd, 0, "schedule end time (unix epoch) for vesting accounts")
	return cmd
}
