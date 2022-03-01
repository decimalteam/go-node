package main

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/cosmos/cosmos-sdk/store/types"

	"github.com/pkg/errors"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/cli"
	"github.com/tendermint/tendermint/libs/log"
	tmtypes "github.com/tendermint/tendermint/types"
	dbm "github.com/tendermint/tm-db"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	authexported "github.com/cosmos/cosmos-sdk/x/auth/exported"
	authvesting "github.com/cosmos/cosmos-sdk/x/auth/vesting"

	"bitbucket.org/decimalteam/go-node/app"
	"bitbucket.org/decimalteam/go-node/config"
	"bitbucket.org/decimalteam/go-node/x/genutil"
	genutilcli "bitbucket.org/decimalteam/go-node/x/genutil/cli"
	"bitbucket.org/decimalteam/go-node/x/validator"
)

const flagInvCheckPeriod = "inv-check-period"

var invCheckPeriod uint

func main() {
	cdc := app.MakeCodec()

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
		Use:               "decd",
		Short:             "Decimal Go Node",
		PersistentPreRunE: server.PersistentPreRunEFn(ctx),
	}

	rootCmd.AddCommand(flags.PostCommands(
		genutilcli.InitCmd(ctx, cdc, app.ModuleBasics, app.DefaultNodeHome),
		genutilcli.CollectGenTxsCmd(ctx, cdc, auth.GenesisAccountIterator{}, app.DefaultNodeHome),
		genutilcli.GenTxCmd(
			ctx, cdc, app.ModuleBasics, validator.AppModuleBasic{},
			auth.GenesisAccountIterator{}, app.DefaultNodeHome, app.DefaultCLIHome,
		),
		genutilcli.GenDeclareCandidateTxCmd(
			ctx, cdc, app.ModuleBasics, validator.AppModuleBasic{},
			auth.GenesisAccountIterator{}, app.DefaultNodeHome, app.DefaultCLIHome,
		),
		genutilcli.ValidateGenesisCmd(ctx, cdc, app.ModuleBasics),
		// AddGenesisAccountCmd allows users to add accounts to the genesis file
		addGenesisAccountCmd(ctx, cdc, app.DefaultNodeHome, app.DefaultCLIHome),
		fixAppHashError(ctx, app.DefaultNodeHome),
	)...)

	server.AddCommands(ctx, cdc, rootCmd, newApp, exportAppStateAndTMValidators)

	// prepare and add flags
	executor := cli.PrepareBaseCmd(rootCmd, "AU", app.DefaultNodeHome)
	rootCmd.PersistentFlags().UintVar(&invCheckPeriod, flagInvCheckPeriod,
		0, "Assert registered invariants every N blocks")
	err := executor.Execute()
	if err != nil {
		panic(err)
	}
}

func newApp(logger log.Logger, db dbm.DB, traceStore io.Writer) abci.Application {
	return app.NewInitApp(
		logger, db, traceStore, true, invCheckPeriod,
		baseapp.SetPruning(types.NewPruningOptionsFromString(viper.GetString("pruning"))),
		baseapp.SetMinGasPrices(viper.GetString(server.FlagMinGasPrices)),
		baseapp.SetHaltHeight(uint64(viper.GetInt(server.FlagHaltHeight))),
	)
}

func exportAppStateAndTMValidators(
	logger log.Logger, db dbm.DB, traceStore io.Writer, height int64, forZeroHeight bool, jailWhiteList []string,
) (json.RawMessage, []tmtypes.GenesisValidator, error) {

	if height != -1 {
		aApp := app.NewInitApp(logger, db, traceStore, true, uint(1))
		err := aApp.LoadHeight(height)
		if err != nil {
			return nil, nil, err
		}
		return aApp.ExportAppStateAndValidators(forZeroHeight, jailWhiteList)
	}

	aApp := app.NewInitApp(logger, db, traceStore, true, uint(1))

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
func addGenesisAccountCmd(ctx *server.Context, cdc *codec.Codec,
	defaultNodeHome, defaultClientHome string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-genesis-account [address_or_key_name] [coin][,[coin]]",
		Short: "Add genesis account to genesis.json",
		Args:  cobra.ExactArgs(2),
		RunE: func(_ *cobra.Command, args []string) error {
			cfg := ctx.Config
			cfg.SetRoot(viper.GetString(cli.HomeFlag))

			addr, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				kb, err := keys.NewKeyBaseFromDir(viper.GetString(flagClientHome))
				if err != nil {
					return err
				}

				info, err := kb.Get(args[0])
				if err != nil {
					return err
				}

				addr = info.GetAddress()
			}

			coins, err := sdk.ParseCoins(args[1])
			if err != nil {
				return err
			}

			vestingStart := viper.GetInt64(flagVestingStart)
			vestingEnd := viper.GetInt64(flagVestingEnd)
			vestingAmt, err := sdk.ParseCoins(viper.GetString(flagVestingAmt))
			if err != nil {
				return err
			}

			// create concrete account type based on input parameters
			var genAccount authexported.GenesisAccount

			baseAccount := auth.NewBaseAccount(addr, coins.Sort(), nil, 0, 0)
			if !vestingAmt.IsZero() {
				baseVestingAccount, err := authvesting.NewBaseVestingAccount(
					baseAccount, vestingAmt.Sort(), vestingEnd,
				)
				if err != nil {
					return err
				}

				switch {
				case vestingStart != 0 && vestingEnd != 0:
					genAccount = authvesting.NewContinuousVestingAccountRaw(baseVestingAccount, vestingStart)

				case vestingEnd != 0:
					genAccount = authvesting.NewDelayedVestingAccountRaw(baseVestingAccount)

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
			genFile := cfg.GenesisFile()
			appState, genDoc, err := genutil.GenesisStateFromGenFile(cdc, genFile)
			if err != nil {
				return err
			}

			authGenState := auth.GetGenesisStateFromAppState(cdc, appState)

			if authGenState.Accounts.Contains(addr) {
				return fmt.Errorf("cannot add account at existing address %s", addr)
			}

			// Add the new account to the set of genesis accounts and sanitize the
			// accounts afterwards.
			authGenState.Accounts = append(authGenState.Accounts, genAccount)
			authGenState.Accounts = auth.SanitizeGenesisAccounts(authGenState.Accounts)

			authGenStateBz, err := cdc.MarshalJSON(authGenState)
			if err != nil {
				return fmt.Errorf("failed to marshal auth genesis state: %w", err)
			}
			appState[auth.ModuleName] = authGenStateBz

			appStateJSON, err := cdc.MarshalJSON(appState)
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
