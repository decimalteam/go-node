package app

import (
	"encoding/json"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	"io"
	"os"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	tos "github.com/tendermint/tendermint/libs/os"
	dbm "github.com/tendermint/tm-db"

	"github.com/cosmos/cosmos-sdk/codec/types"
	codec2 "github.com/cosmos/cosmos-sdk/crypto/codec"

	bam "github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/capability"
	"github.com/cosmos/cosmos-sdk/x/params"
	ibc "github.com/cosmos/ibc-go/modules/core"

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

	stakingKeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	ibckeeper "github.com/cosmos/ibc-go/modules/core/keeper"

	porttypes "github.com/cosmos/ibc-go/modules/core/05-port/types"
	ibchost "github.com/cosmos/ibc-go/modules/core/24-host"

	appParams "bitbucket.org/decimalteam/go-node/app/params"

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
		ibc.AppModuleBasic{},
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

// MakeAminoCodec generates necessary codec for Amino
func MakeAminoCodec() *codec.LegacyAmino {
	var cdc = codec.NewLegacyAmino()
	ModuleBasics.RegisterLegacyAminoCodec(cdc)
	sdk.RegisterLegacyAminoCodec(cdc)
	codec2.RegisterCrypto(cdc)
	return cdc
}

type versionSetter struct {}

func (vs *versionSetter) SetProtocolVersion(version uint64) {}

type newApp struct {
	*bam.BaseApp
	cdc               *codec.LegacyAmino
	appCodec          codec.Codec
	interfaceRegistry types.InterfaceRegistry

	// keys to access the substores
	keys    map[string]*sdk.KVStoreKey
	tkeys   map[string]*sdk.TransientStoreKey
	memKeys map[string]*sdk.MemoryStoreKey

	// Keepers
	accountKeeper    authKeeper.AccountKeeper
	ibcKeeper        *ibckeeper.Keeper
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
func NewInitApp(logger log.Logger, db dbm.DB, traceStore io.Writer, loadLatest bool,
	invCheckPeriod uint, baseAppOptions ...func(*bam.BaseApp)) *newApp {

	encodingConfig := appParams.NewEncodingConfig()
	encodingConfig.Amino = MakeAminoCodec()

	// BaseApp handles interactions with Tendermint through the ABCI protocol
	bApp := bam.NewBaseApp(appName, logger, db, encodingConfig.TxConfig.TxDecoder(), baseAppOptions...)

	bApp.SetVersion(config.DecimalVersion)
	bApp.SetInterfaceRegistry(encodingConfig.InterfaceRegistry)

	// TODO: Add the keys that module requires
	keys := sdk.NewKVStoreKeys(
		bam.Paramspace,
		authTypes.StoreKey,
		paramsTypes.StoreKey,
		coin.StoreKey,
		multisig.StoreKey,
		validator.StoreKey,
		gov.StoreKey,
		swap.StoreKey,
		ibchost.StoreKey,
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

	// Set specific subspaces
	authSubspace := app.paramsKeeper.Subspace(authTypes.ModuleName)
	stakingSubspace := app.paramsKeeper.Subspace(stakingTypes.ModuleName)
	bankSupspace := app.paramsKeeper.Subspace(bankTypes.ModuleName)
	coinSubspace := app.paramsKeeper.Subspace(coin.ModuleName)
	multisigSubspace := app.paramsKeeper.Subspace(multisig.ModuleName)
	validatorSubspace := app.paramsKeeper.Subspace(validator.ModuleName)
	govSubspace := app.paramsKeeper.Subspace(gov.ModuleName).WithKeyTable(gov.ParamKeyTable())
	swapSubspace := app.paramsKeeper.Subspace(swap.ModuleName)
	_ = app.paramsKeeper.Subspace(ibchost.ModuleName)

	// The AccountKeeper handles address -> account lookups
	app.accountKeeper = authKeeper.NewAccountKeeper(
		app.appCodec,
		keys[authTypes.StoreKey],
		authSubspace,
		authTypes.ProtoBaseAccount,
		maccPerms,
	)

	upgradesMap := map[int64]bool{}

	binaryCdc := codec.NewProtoCodec(encodingConfig.InterfaceRegistry)

	// func NewKeeper(skipUpgradeHeights map[int64]bool, storeKey sdk.StoreKey, cdc codec.BinaryCodec, homePath string, vs xp.ProtocolVersionSetter) Keeper {
	app.upgradeKeeper = upgradeKeeper.NewKeeper(upgradesMap, keys[upgradeTypes.StoreKey], binaryCdc , "/upgrades", &versionSetter{})

	app.capabilityKeeper = capabilityKeeper.NewKeeper(app.appCodec, keys[capabilityTypes.StoreKey], memKeys[capabilityTypes.MemStoreKey])
	scopedIBCKeeper := app.capabilityKeeper.ScopeToModule(ibchost.StoreKey)

	// The BankKeeper allows you perform sdk.Coins interactions
	app.bankKeeper = bankKeeper.NewBaseKeeper(
		app.appCodec,
		keys[bankTypes.StoreKey],
		app.accountKeeper,
		bankSupspace,
		app.ModuleAccountAddrs(),
	)

	stakingKeeper := stakingKeeper.NewKeeper(app.appCodec, keys[stakingTypes.StoreKey], app.accountKeeper, app.bankKeeper, stakingSubspace)

	// Create ibc keeper
	app.ibcKeeper = ibckeeper.NewKeeper(
		app.appCodec, keys[ibchost.StoreKey], app.paramsKeeper.Subspace(ibchost.ModuleName), stakingKeeper, app.upgradeKeeper, scopedIBCKeeper,
	)

	app.coinKeeper = coin.NewKeeper(
		app.cdc,
		keys[coin.StoreKey],
		coinSubspace,
		app.accountKeeper,
		app.bankKeeper,
		app.ibcKeeper.ChannelKeeper,
		&app.ibcKeeper.PortKeeper,
		scopedIBCKeeper,
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
	app.govKeeper = gov.NewKeeper(
		app.cdc,
		keys[gov.StoreKey],
		govSubspace,
		app.accountKeeper,
		&app.validatorKeeper,
		govRouter,
	)

	app.swapKeeper = swap.NewKeeper(
		app.cdc,
		keys[swap.StoreKey],
		swapSubspace,
		app.coinKeeper,
		app.accountKeeper,
		app.bankKeeper,
	)

	//type RandomGenesisAccountsFn func(simState *module.SimulationState) GenesisAccounts
	app.mm = module.NewManager(
		genutil.NewAppModule(app.accountKeeper, app.validatorKeeper, app.BaseApp.DeliverTx, encodingConfig.TxConfig),
		auth.NewAppModule(app.appCodec, app.accountKeeper, func(state *module.SimulationState) authTypes.GenesisAccounts {
			return authTypes.GenesisAccounts{}
		}),
		bank.NewAppModule(app.appCodec, app.bankKeeper, app.accountKeeper),
		coin.NewAppModule(app.coinKeeper, app.accountKeeper),
		capability.NewAppModule(app.appCodec, *app.capabilityKeeper),
		multisig.NewAppModule(app.multisigKeeper, app.accountKeeper, app.bankKeeper),
		validator.NewAppModule(app.validatorKeeper, app.accountKeeper, app.bankKeeper, app.coinKeeper),
		gov.NewAppModule(app.govKeeper, app.accountKeeper),
		swap.NewAppModule(app.swapKeeper),
		nft.NewAppModule(app.nftKeeper, app.accountKeeper),
		ibc.NewAppModule(app.ibcKeeper),
	)

	coinModule := coin.NewAppModule(app.coinKeeper, app.accountKeeper)

	ibcRouter := porttypes.NewRouter()
	ibcRouter.AddRoute(coin.ModuleName, coinModule)
	app.ibcKeeper.SetRouter(ibcRouter)

	//app.mm.SetOrderBeginBlockers(distr.ModuleName, /*slashing.ModuleName*/)
	app.mm.SetOrderEndBlockers(validator.ModuleName, gov.ModuleName, ibchost.ModuleName)

	// Sets the order of Genesis - Order matters, genutil is to always come last
	// NOTE: The genutils moodule must occur after staking so that pools are
	// properly initialized with tokens from genesis accounts.
	app.mm.SetOrderInitGenesis(
		capabilityTypes.ModuleName,
		validator.ModuleName,
		authTypes.ModuleName,
		bankTypes.ModuleName,
		coin.ModuleName,
		multisig.ModuleName,
		genutil.ModuleName,
		gov.ModuleName,
		swap.ModuleName,
		nft.ModuleName,
		ibchost.ModuleName,
	)

	// register all module routes and module queriers
	app.mm.RegisterRoutes(app.Router(), app.QueryRouter(), app.cdc)

	// The initChainer handles translating the genesis.json file into initial state for the network
	app.SetInitChainer(app.InitChainer)
	app.SetBeginBlocker(app.BeginBlocker)
	app.SetEndBlocker(app.EndBlocker)

	// The AnteHandler handles signature verification and transaction pre-processing
	app.SetAnteHandler(
		utils.NewAnteHandler(
			app.accountKeeper,
			app.bankKeeper,
			app.validatorKeeper,
			app.coinKeeper,
			ante.DefaultSigVerificationGasConsumer,
		),
	)

	// initialize stores
	app.MountKVStores(keys)
	app.MountTransientStores(tkeys)

	err := app.LoadLatestVersion()
	if err != nil {
		tos.Exit(err.Error())
	}

	return app
}

// GenesisState represents chain state at the start of the chain. Any initial state (account balances) are stored here.
type GenesisState map[string]json.RawMessage

func (app *newApp) NewDefaultGenesisState() GenesisState {
	return ModuleBasics.DefaultGenesis(app.appCodec)
}

func (app *newApp) InitChainer(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
	var genesisState GenesisState

	err := app.cdc.UnmarshalJSON(req.AppStateBytes, &genesisState)
	if err != nil {
		panic(err)
	}

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
