package app

import (
	"encoding/json"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/grpc/tmservice"
	"github.com/cosmos/cosmos-sdk/server/api"
	config2 "github.com/cosmos/cosmos-sdk/server/config"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	"github.com/cosmos/cosmos-sdk/x/auth/simulation"
	authTx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	"io"
	"os"

	appparams "bitbucket.org/decimalteam/go-node/app/params"
	abci "github.com/tendermint/tendermint/abci/types"
	tmjson "github.com/tendermint/tendermint/libs/json"
	"github.com/tendermint/tendermint/libs/log"
	tos "github.com/tendermint/tendermint/libs/os"
	dbm "github.com/tendermint/tm-db"

	bam "github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/capability"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/upgrade"

	authKeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authTypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	bankKeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	bankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	capabilityKeeper "github.com/cosmos/cosmos-sdk/x/capability/keeper"
	capabilityTypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	paramsKeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	paramsTypes "github.com/cosmos/cosmos-sdk/x/params/types"

	upgradeKeeper "github.com/cosmos/cosmos-sdk/x/upgrade/keeper"
	upgradeTypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"

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
		coin.AppModuleBasic{},
		multisig.AppModuleBasic{},
		validator.AppModuleBasic{},
		capability.AppModuleBasic{},
		gov.AppModuleBasic{},
		swap.AppModuleBasic{},
		nft.AppModuleBasic{},
	)
	// account permissions
	maccPerms = map[string][]string{
		authTypes.FeeCollectorName:  {authTypes.Burner, authTypes.Minter},
		validator.BondedPoolName:    {authTypes.Burner, authTypes.Staking},
		validator.NotBondedPoolName: {authTypes.Burner, authTypes.Staking},
		swap.PoolName:               {authTypes.Minter, authTypes.Burner},
		nft.ReservedPool:            {authTypes.Burner},
	}
)

type newApp struct {
	*bam.BaseApp
	cdc               *codec.LegacyAmino
	appCodec          codec.Marshaler
	interfaceRegistry types.InterfaceRegistry

	// keys to access the substores
	keys    map[string]*sdk.KVStoreKey
	tkeys   map[string]*sdk.TransientStoreKey
	memKeys map[string]*sdk.MemoryStoreKey

	// Keepers
	accountKeeper    authKeeper.AccountKeeper
	capabilityKeeper *capabilityKeeper.Keeper
	bankKeeper       bankKeeper.BaseKeeper
	paramsKeeper     paramsKeeper.Keeper
	coinKeeper       coin.Keeper
	multisigKeeper   multisig.Keeper
	validatorKeeper  validator.Keeper
	govKeeper        gov.Keeper
	swapKeeper       swap.Keeper
	nftKeeper        nft.Keeper
	upgradeKeeper    upgradeKeeper.Keeper

	// Module Manager
	mm *module.Manager
}

// Newgo-nodeApp is a constructor function for go-nodeApp
func NewInitApp(logger log.Logger, db dbm.DB, traceStore io.Writer, loadLatest bool, skipUpgradeHeights map[int64]bool,
	homePath string, invCheckPeriod uint, encodingConfig appparams.EncodingConfig, baseAppOptions ...func(*bam.BaseApp)) *newApp {
	// BaseApp handles interactions with Tendermint through the ABCI protocol
	bApp := bam.NewBaseApp(appName, logger, db, encodingConfig.TxConfig.TxDecoder(), baseAppOptions...)

	bApp.SetInterfaceRegistry(encodingConfig.InterfaceRegistry)
	bApp.SetCommitMultiStoreTracer(traceStore)

	// TODO: Add the keys that module requires
	keys := sdk.NewKVStoreKeys(
		bam.Paramspace,
		authTypes.StoreKey,
		paramsTypes.StoreKey,
		bankTypes.StoreKey,
		coin.StoreKey,
		multisig.StoreKey,
		validator.StoreKey,
		gov.StoreKey,
		upgradeTypes.ModuleName,
		swap.StoreKey,
		capabilityTypes.StoreKey,
	)

	tkeys := sdk.NewTransientStoreKeys(paramsTypes.TStoreKey)
	memKeys := sdk.NewMemoryStoreKeys(capabilityTypes.MemStoreKey)

	config := config.GetDefaultConfig(config.ChainID)

	// Here you initialize your application with the store keys it requires
	var app = &newApp{
		BaseApp:           bApp,
		cdc:               encodingConfig.Amino,
		appCodec:          encodingConfig.Codec,
		interfaceRegistry: encodingConfig.InterfaceRegistry,
		keys:              keys,
		tkeys:             tkeys,
		memKeys:           memKeys,
	}

	// The ParamsKeeper handles parameter storage for the application
	app.paramsKeeper = paramsKeeper.NewKeeper(app.appCodec, app.cdc, keys[paramsTypes.StoreKey], tkeys[paramsTypes.TStoreKey])

	// set the BaseApp's parameter store
	bApp.SetParamStore(app.paramsKeeper.Subspace(bam.Paramspace).WithKeyTable(paramsKeeper.ConsensusParamsKeyTable()))

	// Set specific subspaces
	authSubspace := app.paramsKeeper.Subspace(authTypes.ModuleName)
	bankSupspace := app.paramsKeeper.Subspace(bankTypes.ModuleName)
	coinSubspace := app.paramsKeeper.Subspace(coin.ModuleName)
	multisigSubspace := app.paramsKeeper.Subspace(multisig.ModuleName)
	validatorSubspace := app.paramsKeeper.Subspace(validator.ModuleName)
	govSubspace := app.paramsKeeper.Subspace(gov.ModuleName).WithKeyTable(gov.ParamKeyTable())
	swapSubspace := app.paramsKeeper.Subspace(swap.ModuleName)

	// The AccountKeeper handles address -> account lookups
	app.accountKeeper = authKeeper.NewAccountKeeper(
		app.appCodec,
		keys[authTypes.StoreKey],
		authSubspace,
		authTypes.ProtoBaseAccount,
		maccPerms,
	)

	app.capabilityKeeper = capabilityKeeper.NewKeeper(app.appCodec, keys[capabilityTypes.StoreKey], memKeys[capabilityTypes.MemStoreKey])

	// The BankKeeper allows you perform sdk.Coins interactions
	app.bankKeeper = bankKeeper.NewBaseKeeper(
		app.appCodec,
		keys[bankTypes.StoreKey],
		app.accountKeeper,
		bankSupspace,
		app.ModuleAccountAddrs(),
	)

	app.coinKeeper = coin.NewKeeper(
		app.cdc,
		keys[coin.StoreKey],
		coinSubspace,
		app.accountKeeper,
		app.bankKeeper,
		config,
	)

	app.multisigKeeper = multisig.NewKeeper(
		app.cdc,
		keys[multisig.StoreKey],
		multisigSubspace,
		app.accountKeeper,
		app.coinKeeper,
		app.bankKeeper,
	)

	app.nftKeeper = nft.NewKeeper(
		app.cdc,
		keys[nft.StoreKey],
		app.bankKeeper,
		app.accountKeeper,
		validator.DefaultBondDenom,
	)

	app.validatorKeeper = validator.NewKeeper(
		app.cdc,
		keys[validator.StoreKey],
		validatorSubspace,
		app.coinKeeper,
		app.accountKeeper,
		app.bankKeeper,
		app.multisigKeeper,
		app.nftKeeper,
		authTypes.FeeCollectorName,
	)

	// register the proposal types
	govRouter := gov.NewRouter()
	//govRouter.AddRoute(
	//	upgradeTypes.RouterKey, upgrade.NewSoftwareUpgradeProposalHandler(app.upgradeKeeper),
	//)
	app.govKeeper = gov.NewKeeper(
		app.cdc,
		keys[gov.StoreKey],
		govSubspace,
		app.accountKeeper,
		&app.validatorKeeper,
		govRouter,
	)

	app.swapKeeper = swap.NewKeeper(
		*app.cdc,
		keys[swap.StoreKey],
		swapSubspace,
		app.coinKeeper,
		app.accountKeeper,
		app.bankKeeper,
	)

	app.upgradeKeeper = upgradeKeeper.NewKeeper(skipUpgradeHeights, keys[upgradeTypes.StoreKey], encodingConfig.Codec, homePath)

	app.mm = module.NewManager(
		genutil.NewAppModule(app.accountKeeper, app.validatorKeeper, app.BaseApp.DeliverTx, encodingConfig.TxConfig),
		auth.NewAppModule(app.appCodec, app.accountKeeper, simulation.RandomGenesisAccounts),
		bank.NewAppModule(app.appCodec, app.bankKeeper, app.accountKeeper),
		coin.NewAppModule(app.coinKeeper, app.accountKeeper),
		capability.NewAppModule(app.appCodec, *app.capabilityKeeper),
		multisig.NewAppModule(app.multisigKeeper, app.accountKeeper, app.bankKeeper),
		validator.NewAppModule(app.validatorKeeper, app.accountKeeper, app.bankKeeper, app.coinKeeper),
		gov.NewAppModule(app.govKeeper, app.accountKeeper),
		swap.NewAppModule(app.swapKeeper),
		nft.NewAppModule(app.nftKeeper, app.accountKeeper),
		upgrade.NewAppModule(app.upgradeKeeper),
	)

	app.mm.SetOrderBeginBlockers(upgradeTypes.ModuleName)
	app.mm.SetOrderEndBlockers(validator.ModuleName, gov.ModuleName)

	// Sets the order of Genesis - Order matters, genutil is to always come last
	// NOTE: The genutils moodule must occur after staking so that pools are
	// properly initialized with tokens from genesis accounts.
	app.mm.SetOrderInitGenesis(
		capabilityTypes.ModuleName,
		authTypes.ModuleName,
		bankTypes.ModuleName,
		validator.ModuleName,
		coin.ModuleName,
		multisig.ModuleName,
		genutil.ModuleName,
		gov.ModuleName,
		swap.ModuleName,
		nft.ModuleName,
	)

	// register all module routes and module queriers
	app.mm.RegisterRoutes(app.Router(), app.QueryRouter(), app.cdc)

	// initialize stores
	app.MountKVStores(keys)
	app.MountTransientStores(tkeys)
	app.MountMemoryStores(memKeys)

	// The initChainer handles translating the genesis.json file into initial state for the network
	app.SetInitChainer(app.InitChainer)
	app.SetBeginBlocker(app.BeginBlocker)

	// The AnteHandler handles signature verification and transaction pre-processing
	app.SetAnteHandler(
		utils.NewAnteHandler(
			app.accountKeeper,
			app.bankKeeper,
			bankKeeper.NewBaseViewKeeper(app.appCodec, keys[bankTypes.StoreKey], app.accountKeeper),
			app.validatorKeeper,
			app.coinKeeper,
			ante.DefaultSigVerificationGasConsumer,
		),
	)

	app.SetEndBlocker(app.EndBlocker)

	if loadLatest {
		err := app.LoadLatestVersion()
		if err != nil {
			tos.Exit(err.Error())
		}

		ctx := app.BaseApp.NewUncachedContext(true, tmproto.Header{})
		app.capabilityKeeper.InitializeAndSeal(ctx)
	}

	return app
}

// GenesisState represents chain state at the start of the chain. Any initial state (account balances) are stored here.
type GenesisState map[string]json.RawMessage

func (app *newApp) NewDefaultGenesisState() GenesisState {
	return ModuleBasics.DefaultGenesis(app.appCodec)
}

func (app *newApp) RegisterAPIRoutes(server *api.Server, _ config2.APIConfig) {
	clientCtx := server.ClientCtx

	authTx.RegisterGRPCGatewayRoutes(clientCtx, server.GRPCGatewayRouter)

	tmservice.RegisterGRPCGatewayRoutes(clientCtx, server.GRPCGatewayRouter)

	ModuleBasics.RegisterRESTRoutes(clientCtx, server.Router)
	ModuleBasics.RegisterGRPCGatewayRoutes(clientCtx, server.GRPCGatewayRouter)
}

func (app *newApp) RegisterTxService(clientCtx client.Context) {
	authTx.RegisterTxService(app.BaseApp.GRPCQueryRouter(), clientCtx, app.BaseApp.Simulate, app.interfaceRegistry)
}

func (app *newApp) RegisterTendermintService(clientCtx client.Context) {
	tmservice.RegisterTendermintService(app.BaseApp.GRPCQueryRouter(), clientCtx, app.interfaceRegistry)
}

func (app *newApp) InitChainer(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
	var genesisState GenesisState

	err := tmjson.Unmarshal(req.AppStateBytes, &genesisState)
	if err != nil {
		panic(err)
	}
	//app.upgradeKeeper.SetModuleVersionMap(ctx, app.mm.GetVersionMap())
	return app.mm.InitGenesis(ctx, app.appCodec, genesisState)
}

func (app *newApp) BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock) abci.ResponseBeginBlock {
	return app.BeginBlock(req)
}
func (app *newApp) EndBlocker(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
	return app.EndBlock(req)
}
func (app *newApp) LoadHeight(height int64) error {
	return app.LoadVersion(height)
}

// ModuleAccountAddrs returns all the app's module account addresses.
func (app *newApp) ModuleAccountAddrs() map[string]bool {
	modAccAddrs := make(map[string]bool)
	for acc := range maccPerms {
		modAccAddrs[authTypes.NewModuleAddress(acc).String()] = true
	}
	return modAccAddrs
}
