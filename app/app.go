package app

import (
	"encoding/json"
	"io"
	"os"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	tos "github.com/tendermint/tendermint/libs/os"
	dbm "github.com/tendermint/tm-db"

	bam "github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/supply"

	"bitbucket.org/decimalteam/go-node/config"
	"bitbucket.org/decimalteam/go-node/x/check"
	"bitbucket.org/decimalteam/go-node/x/coin"
	"bitbucket.org/decimalteam/go-node/x/genutil"
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
		//staking.AppModuleBasic{},
		//distr.AppModuleBasic{},
		params.AppModuleBasic{},
		//slashing.AppModuleBasic{},
		supply.AppModuleBasic{},
		coin.AppModuleBasic{},
		check.AppModuleBasic{},
		validator.AppModuleBasic{},
	)
	// account permissions
	maccPerms = map[string][]string{
		auth.FeeCollectorName: {supply.Burner, supply.Minter},
		//distr.ModuleName:            nil,
		validator.BondedPoolName:    {supply.Burner, supply.Staking},
		validator.NotBondedPoolName: {supply.Burner, supply.Staking},
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
	accountKeeper auth.AccountKeeper
	bankKeeper    bank.Keeper
	//stakingKeeper   staking.Keeper
	//slashingKeeper  slashing.Keeper
	//distrKeeper     distr.Keeper
	supplyKeeper    supply.Keeper
	paramsKeeper    params.Keeper
	coinKeeper      coin.Keeper
	validatorKeeper validator.Keeper
	checkKeeper     check.Keeper

	// Module Manager
	mm *module.Manager
}

// Newgo-nodeApp is a constructor function for go-nodeApp
func NewInitApp(logger log.Logger, db dbm.DB, traceStore io.Writer, loadLatest bool,
	invCheckPeriod uint, baseAppOptions ...func(*bam.BaseApp)) *newApp {

	// First define the top level codec that will be shared by the different modules
	cdc := MakeCodec()

	// BaseApp handles interactions with Tendermint through the ABCI protocol
	bApp := bam.NewBaseApp(appName, logger, db, auth.DefaultTxDecoder(cdc), baseAppOptions...)

	bApp.SetAppVersion(version.Version)

	// TODO: Add the keys that module requires
	keys := sdk.NewKVStoreKeys(bam.MainStoreKey, auth.StoreKey, /*staking.StoreKey,*/
		supply.StoreKey /*, distr.StoreKey*/ /*slashing.StoreKey,*/, params.StoreKey, coin.StoreKey, validator.StoreKey)

	tkeys := sdk.NewTransientStoreKeys(params.TStoreKey)

	config := config.GetDefaultConfig(config.ChainID)

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
	//stakingSubspace := app.paramsKeeper.Subspace(staking.DefaultParamspace)
	//distrSubspace := app.paramsKeeper.Subspace(distr.DefaultParamspace)
	//slashingSubspace := app.paramsKeeper.Subspace(slashing.DefaultParamspace)
	coinSubspace := app.paramsKeeper.Subspace(coin.DefaultParamspace)
	checkSubspace := app.paramsKeeper.Subspace(check.DefaultParamspace)
	validatorSubspace := app.paramsKeeper.Subspace(validator.DefaultParamSpace)
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

	// The staking keeper
	//stakingKeeper := staking.NewKeeper(
	//	app.cdc,
	//	keys[staking.StoreKey],
	//	tkeys[staking.TStoreKey],
	//	app.supplyKeeper,
	//	stakingSubspace,
	//)

	//app.distrKeeper = distr.NewKeeper(
	//	app.cdc,
	//	keys[distr.StoreKey],
	//	distrSubspace,
	//	&stakingKeeper,
	//	app.supplyKeeper,
	//	auth.FeeCollectorName,
	//	app.ModuleAccountAddrs(),
	//)

	//app.slashingKeeper = slashing.NewKeeper(
	//	app.cdc,
	//	keys[slashing.StoreKey],
	//	&stakingKeeper,
	//	slashingSubspace,
	//)

	// register the staking hooks
	// NOTE: stakingKeeper above is passed by reference, so that it will contain these hooks
	//app.stakingKeeper = *stakingKeeper.SetHooks(
	//	staking.NewMultiStakingHooks(
	//		app.distrKeeper.Hooks(),
	//		app.slashingKeeper.Hooks()),
	//)

	app.coinKeeper = coin.NewKeeper(
		app.cdc,
		keys[coin.StoreKey],
		coinSubspace,
		app.accountKeeper,
		app.bankKeeper,
		config,
	)

	app.checkKeeper = check.NewKeeper(
		app.cdc,
		keys[check.StoreKey],
		checkSubspace,
		app.coinKeeper,
		app.accountKeeper,
	)

	app.validatorKeeper = validator.NewKeeper(
		app.cdc,
		keys[validator.StoreKey],
		validatorSubspace,
		app.coinKeeper,
		app.supplyKeeper,
		auth.FeeCollectorName,
	)

	app.mm = module.NewManager(
		genutil.NewAppModule(app.accountKeeper, app.validatorKeeper, app.BaseApp.DeliverTx),
		auth.NewAppModule(app.accountKeeper),
		bank.NewAppModule(app.bankKeeper, app.accountKeeper),
		supply.NewAppModule(app.supplyKeeper, app.accountKeeper),
		//distr.NewAppModule(app.distrKeeper, app.supplyKeeper),
		coin.NewAppModule(app.coinKeeper, app.accountKeeper),
		validator.NewAppModule(app.validatorKeeper, app.supplyKeeper, app.coinKeeper),
		check.NewAppModule(app.checkKeeper, app.coinKeeper, app.accountKeeper),
		//slashing.NewAppModule(app.slashingKeeper, app.validatorKeeper),
		//staking.NewAppModule(app.stakingKeeper, app.distrKeeper, app.accountKeeper, app.supplyKeeper),
	)

	//app.mm.SetOrderBeginBlockers(distr.ModuleName, /*slashing.ModuleName*/)
	app.mm.SetOrderEndBlockers(validator.ModuleName)

	// Sets the order of Genesis - Order matters, genutil is to always come last
	// NOTE: The genutils moodule must occur after staking so that pools are
	// properly initialized with tokens from genesis accounts.
	app.mm.SetOrderInitGenesis(
		//distr.ModuleName,
		//staking.ModuleName,
		validator.ModuleName,
		auth.ModuleName,
		bank.ModuleName,
		//slashing.ModuleName,
		coin.ModuleName,
		supply.ModuleName,
		check.ModuleName,
		genutil.ModuleName,
	)

	// register all module routes and module queriers
	app.mm.RegisterRoutes(app.Router(), app.QueryRouter())

	// The initChainer handles translating the genesis.json file into initial state for the network
	app.SetInitChainer(app.InitChainer)
	app.SetBeginBlocker(app.BeginBlocker)
	app.SetEndBlocker(app.EndBlocker)

	// The AnteHandler handles signature verification and transaction pre-processing
	app.SetAnteHandler(
		validator.NewAnteHandler(
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

	err := app.LoadLatestVersion(app.keys[bam.MainStoreKey])
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
	var genesisState GenesisState

	err := app.cdc.UnmarshalJSON(req.AppStateBytes, &genesisState)
	if err != nil {
		panic(err)
	}

	return app.mm.InitGenesis(ctx, genesisState)
}

func (app *newApp) BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock) abci.ResponseBeginBlock {
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

	// Получаем адрес и ключи формата Ed25519 (Ethereum)
	//privateKey, err := crypto.GenerateKey()
	//if err != nil {
	//	panic(err)
	//}
	//privateKeyBytes := crypto.FromECDSA(privateKey)
	//publicKey := privateKey.Public()
	//publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	//if !ok {
	//	panic(err)
	//}
	//publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)
	//address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
	//hash := sha3.NewLegacyKeccak256()
	//hash.Write(publicKeyBytes[1:])

	//fmt.Printf("address: (%v) \n", address)
	//fmt.Printf("publicKey: (%v) \n", hexutil.Encode(hash.Sum(nil)[12:]))
	//fmt.Printf("publicKeyBytes: (%v) \n", hexutil.Encode(publicKeyBytes)[4:])
	//fmt.Printf("privateKeyBytes: (%v) \n", hexutil.Encode(privateKeyBytes)[2:])

	modAccAddrs := make(map[string]bool)
	for acc := range maccPerms {
		modAccAddrs[supply.NewModuleAddress(acc).String()] = true
	}
	return modAccAddrs
}
