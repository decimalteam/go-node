package app

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"bitbucket.org/decimalteam/go-node/utils/updates"
	genutilcli "bitbucket.org/decimalteam/go-node/x/genutil/cli"
	"github.com/cosmos/cosmos-sdk/x/supply/exported"
	tmtypes "github.com/tendermint/tendermint/types"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	tos "github.com/tendermint/tendermint/libs/os"
	dbm "github.com/tendermint/tm-db"

	bam "github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/supply"

	"bitbucket.org/decimalteam/go-node/config"
	"bitbucket.org/decimalteam/go-node/utils"
	"bitbucket.org/decimalteam/go-node/x/coin"
	"bitbucket.org/decimalteam/go-node/x/genutil"
	"bitbucket.org/decimalteam/go-node/x/gov"
	"bitbucket.org/decimalteam/go-node/x/multisig"
	"bitbucket.org/decimalteam/go-node/x/nft"
	"bitbucket.org/decimalteam/go-node/x/swap"
	"bitbucket.org/decimalteam/go-node/x/validator"
)

const appName = "decimal"

var (
	// default home directories for the application CLI
	DefaultCLIHome = os.ExpandEnv("$HOME/.decimal/cli")

	// DefaultNodeHome sets the folder where the applcation data and configuration will be stored
	DefaultNodeHome = os.ExpandEnv("$HOME/.decimal/daemon")

	// NewBasicManager is in charge of setting up basic module elements
	ModuleBasics = module.NewBasicManager(
		genutil.AppModuleBasic{},
		auth.AppModuleBasic{},
		bank.AppModuleBasic{},
		params.AppModuleBasic{},
		supply.AppModuleBasic{},
		coin.AppModuleBasic{},
		multisig.AppModuleBasic{},
		validator.AppModuleBasic{},
		gov.AppModuleBasic{},
		swap.AppModuleBasic{},
		nft.AppModuleBasic{},
	)
	// account permissions
	maccPerms = map[string][]string{
		auth.FeeCollectorName:       {supply.Burner, supply.Minter},
		validator.BondedPoolName:    {supply.Burner, supply.Staking},
		validator.NotBondedPoolName: {supply.Burner, supply.Staking},
		swap.PoolName:               {supply.Minter, supply.Burner},
		nft.ReservedPool:            {supply.Burner},
	}
)

// MakeCodec generates the necessary codecs for Amino
func MakeCodec() *codec.Codec {
	var cdc = codec.New()
	ModuleBasics.RegisterCodec(cdc)
	sdk.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)
	return cdc
}

type newApp struct {
	*bam.BaseApp
	cdc *codec.Codec

	// keys to access the substores
	keys  map[string]*sdk.KVStoreKey
	tkeys map[string]*sdk.TransientStoreKey

	// Keepers
	accountKeeper   auth.AccountKeeper
	bankKeeper      bank.Keeper
	supplyKeeper    supply.Keeper
	paramsKeeper    params.Keeper
	coinKeeper      coin.Keeper
	multisigKeeper  multisig.Keeper
	validatorKeeper validator.Keeper
	govKeeper       gov.Keeper
	swapKeeper      swap.Keeper
	nftKeeper       nft.Keeper

	// Module Manager
	mm *module.Manager

	updated   bool
	initChain bool
}

var cfg = &config.Config{}

// Newgo-nodeApp is a constructor function for go-nodeApp
func NewInitApp(logger log.Logger, db dbm.DB, traceStore io.Writer, loadLatest bool,
	invCheckPeriod uint, baseAppOptions ...func(*bam.BaseApp)) *newApp {
	//fmt.Printf("decd version: %s\n", config.DecimalVersion)

	// First define the top level codec that will be shared by the different modules
	cdc := MakeCodec()

	// BaseApp handles interactions with Tendermint through the ABCI protocol
	bApp := bam.NewBaseApp(appName, logger, db, auth.DefaultTxDecoder(cdc), baseAppOptions...)

	// Load file with updates info: last_update and all_updates
	err := config.UpdatesInfo.Load()
	if err != nil {
		panic(fmt.Sprintf("error: read permissions '%s'", err.Error()))
	}
	config.UpdatesInfo.PushNewPlanHeight(config.UpdatesInfo.LastBlock)
	err = config.UpdatesInfo.Save()
	if err != nil {
		panic(fmt.Sprintf("error: write permissions '%s'", err.Error()))
	}

	bApp.SetAppVersion(config.DecimalVersion)

	keys := sdk.NewKVStoreKeys(
		bam.MainStoreKey,
		auth.StoreKey,
		supply.StoreKey,
		params.StoreKey,
		coin.StoreKey,
		multisig.StoreKey,
		validator.StoreKey,
	)

	tkeys := sdk.NewTransientStoreKeys(params.TStoreKey)

	cfg = config.GetDefaultConfig(config.ChainID)

	// Here you initialize your application with the store keys it requires
	var app = &newApp{
		BaseApp: bApp,
		cdc:     cdc,
		keys:    keys,
		tkeys:   tkeys,
	}

	// The ParamsKeeper handles parameter storage for the application
	app.paramsKeeper = params.NewKeeper(app.cdc, keys[params.StoreKey], tkeys[params.TStoreKey])
	// Set specific subspaces
	authSubspace := app.paramsKeeper.Subspace(auth.DefaultParamspace)
	bankSupspace := app.paramsKeeper.Subspace(bank.DefaultParamspace)
	coinSubspace := app.paramsKeeper.Subspace(coin.DefaultParamspace)
	multisigSubspace := app.paramsKeeper.Subspace(multisig.DefaultParamspace)
	validatorSubspace := app.paramsKeeper.Subspace(validator.DefaultParamSpace)
	govSubspace := app.paramsKeeper.Subspace(gov.DefaultParamspace).WithKeyTable(gov.ParamKeyTable())
	swapSubspace := app.paramsKeeper.Subspace(swap.DefaultParamspace)

	// The AccountKeeper handles address -> account lookups
	app.accountKeeper = auth.NewAccountKeeper(
		app.cdc,
		keys[auth.StoreKey],
		authSubspace,
		auth.ProtoBaseAccount,
	)

	// The BankKeeper allows you perform sdk.Coins interactions
	app.bankKeeper = bank.NewBaseKeeper(
		app.accountKeeper,
		bankSupspace,
		app.ModuleAccountAddrs(),
	)

	// The SupplyKeeper collects transaction fees and renders them to the fee distribution module
	app.supplyKeeper = supply.NewKeeper(
		app.cdc,
		keys[supply.StoreKey],
		app.accountKeeper,
		app.bankKeeper,
		maccPerms,
	)

	app.coinKeeper = coin.NewKeeper(
		app.cdc,
		keys[coin.StoreKey],
		coinSubspace,
		app.accountKeeper,
		app.bankKeeper,
		app.supplyKeeper,
		cfg,
	)

	app.multisigKeeper = multisig.NewKeeper(
		app.cdc,
		keys[multisig.StoreKey],
		multisigSubspace,
		app.accountKeeper,
		app.coinKeeper,
		app.bankKeeper,
	)

	app.nftKeeper = nft.NewKeeper(app.cdc, keys[nft.StoreKey], app.supplyKeeper, validator.DefaultBondDenom)

	app.validatorKeeper = validator.NewKeeper(
		app.cdc,
		keys[validator.StoreKey],
		validatorSubspace,
		app.coinKeeper,
		app.accountKeeper,
		app.supplyKeeper,
		app.multisigKeeper,
		app.nftKeeper,
		auth.FeeCollectorName,
	)

	govRouter := gov.NewRouter()
	app.govKeeper = gov.NewKeeper(
		app.cdc,
		app.keys[gov.StoreKey],
		govSubspace,
		app.supplyKeeper,
		&app.validatorKeeper,
		govRouter,
	)

	app.swapKeeper = swap.NewKeeper(
		app.cdc,
		app.keys[swap.StoreKey],
		swapSubspace,
		app.coinKeeper,
		app.accountKeeper,
		app.supplyKeeper,
	)

	app.mm = module.NewManager(
		genutil.NewAppModule(app.accountKeeper, app.validatorKeeper, app.BaseApp.DeliverTx),
		auth.NewAppModule(app.accountKeeper),
		bank.NewAppModule(app.bankKeeper, app.accountKeeper),
		supply.NewAppModule(app.supplyKeeper, app.accountKeeper),
		coin.NewAppModule(app.coinKeeper, app.accountKeeper),
		multisig.NewAppModule(app.multisigKeeper, app.accountKeeper, app.bankKeeper),
		validator.NewAppModule(app.validatorKeeper, app.supplyKeeper, app.coinKeeper),
		swap.NewAppModule(app.swapKeeper),
		gov.NewAppModule(app.govKeeper, app.accountKeeper, app.supplyKeeper),
		nft.NewAppModule(app.nftKeeper, app.accountKeeper),
	)

	// app.mm.SetOrderBeginBlockers(coin.ModuleName, validator.ModuleName)
	app.mm.SetOrderEndBlockers(validator.ModuleName, gov.ModuleName)

	// Sets the order of Genesis - Order matters, genutil is to always come last
	// NOTE: The genutils moodule must occur after staking so that pools are
	// properly initialized with tokens from genesis accounts.
	app.mm.SetOrderInitGenesis(
		validator.ModuleName,
		auth.ModuleName,
		bank.ModuleName,
		coin.ModuleName,
		supply.ModuleName,
		multisig.ModuleName,
		genutil.ModuleName,
		swap.ModuleName,
		gov.ModuleName,
		nft.ModuleName,
	)

	// register all module routes and module queriers
	app.mm.RegisterRoutes(app.Router(), app.QueryRouter())

	// The initChainer handles translating the genesis.json file into initial state for the network
	app.SetInitChainer(app.InitChainer)
	app.SetBeginBlocker(app.BeginBlocker)
	app.SetEndBlocker(app.EndBlocker)

	// The AnteHandler handles signature verification and transaction pre-processing
	app.SetAnteHandler(
		utils.NewAnteHandler(
			app.accountKeeper,
			app.validatorKeeper,
			app.coinKeeper,
			app.supplyKeeper,
			auth.DefaultSigVerificationGasConsumer,
		),
	)

	// initialize stores
	app.MountKVStores(keys)
	app.MountTransientStores(tkeys)

	err = app.LoadLatestVersion(app.keys[bam.MainStoreKey])
	if err != nil {
		tos.Exit(err.Error())
	}

	return app
}

// GenesisState represents chain state at the start of the chain. Any initial state (account balances) are stored here.
type GenesisState map[string]json.RawMessage

func NewDefaultGenesisState() GenesisState {
	return ModuleBasics.DefaultGenesis()
}

func (app *newApp) InitChainer(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
	if app.initChain {
		return abci.ResponseInitChain{}
	}
	var genesisState GenesisState

	err := app.cdc.UnmarshalJSON(req.AppStateBytes, &genesisState)
	if err != nil {
		panic(err)
	}

	return app.mm.InitGenesis(ctx, genesisState)
}

func (app *newApp) BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock) abci.ResponseBeginBlock {
	if !cfg.Initialized {
		config.ChainID = ctx.ChainID()
		if strings.HasPrefix(config.ChainID, "decimal-testnet") {
			cfg.TitleBaseCoin = config.TitleTestBaseCoin
			cfg.SymbolBaseCoin = config.SymbolTestBaseCoin
			cfg.InitialVolumeBaseCoin = config.InitialVolumeTestBaseCoin
		} else if strings.HasPrefix(config.ChainID, "decimal") {
			cfg.TitleBaseCoin = config.TitleBaseCoin
			cfg.SymbolBaseCoin = config.SymbolBaseCoin
			cfg.InitialVolumeBaseCoin = config.InitialVolumeBaseCoin
		}
		cfg.Initialized = true
	}

	if ctx.BlockHeight() == updates.Update2Block {
		govAppModule := app.mm.Modules[gov.ModuleName].(gov.AppModule)
		swapAppModule := app.mm.Modules[swap.ModuleName].(swap.AppModule)

		genesis := genutilcli.MainNetGenesis
		var err error

		var genDoc *tmtypes.GenesisDoc
		if genDoc, err = tmtypes.GenesisDocFromJSON([]byte(genesis)); err != nil {
			panic(err)
		}

		var genState map[string]json.RawMessage
		if err = app.cdc.UnmarshalJSON(genDoc.AppState, &genState); err != nil {
			panic(err)
		}

		govAppModule.InitGenesis(ctx, genState[gov.ModuleName])
		gov.InitGenesis(ctx, app.govKeeper, gov.InitialGenesisState)

		swapAppModule.InitGenesis(ctx, genState[swap.ModuleName])
		swap.InitGenesis(ctx, app.swapKeeper, app.supplyKeeper, swap.InitialGenesisState)

		moduleAddress := app.supplyKeeper.GetModuleAddress(swap.PoolName)
		moduleAccount := app.accountKeeper.GetAccount(ctx, moduleAddress)

		if moduleAccount != nil {
			if _, ok := moduleAccount.(exported.ModuleAccountI); !ok {
				app.accountKeeper.RemoveAccount(ctx, moduleAccount)
			}
		}
	}

	return app.mm.BeginBlock(ctx, req)
}
func (app *newApp) EndBlocker(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
	return app.mm.EndBlock(ctx, req)
}
func (app *newApp) LoadHeight(height int64) error {
	return app.LoadVersion(height, app.keys[bam.MainStoreKey])
}

// ModuleAccountAddrs returns all the app's module account addresses.
func (app *newApp) ModuleAccountAddrs() map[string]bool {
	modAccAddrs := make(map[string]bool)
	for acc := range maccPerms {
		modAccAddrs[supply.NewModuleAddress(acc).String()] = true
	}
	return modAccAddrs
}
