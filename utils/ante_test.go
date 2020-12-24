package utils

import (
	"bitbucket.org/decimalteam/go-node/config"
	"bitbucket.org/decimalteam/go-node/utils/updates"
	"bitbucket.org/decimalteam/go-node/x/coin"
	"bitbucket.org/decimalteam/go-node/x/genutil"
	"bitbucket.org/decimalteam/go-node/x/multisig"
	"bitbucket.org/decimalteam/go-node/x/swap"
	"bitbucket.org/decimalteam/go-node/x/validator"
	"encoding/json"
	bam "github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/auth/vesting"
	"github.com/cosmos/cosmos-sdk/x/bank"
	tmos "github.com/tendermint/tendermint/libs/os"

	"github.com/cosmos/cosmos-sdk/x/gov"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/supply"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"
	"io"
	"testing"
)

var (
	ModuleBasics = module.NewBasicManager(
		auth.AppModuleBasic{},
		supply.AppModuleBasic{},
		genutil.AppModuleBasic{},
		bank.AppModuleBasic{},
		gov.AppModuleBasic{},
		params.AppModuleBasic{},
		coin.AppModuleBasic{},
		swap.AppModuleBasic{},
		validator.AppModuleBasic{},
	)

	// module account permissions
	maccPerms = map[string][]string{
		auth.FeeCollectorName:       {supply.Burner, supply.Minter},
		validator.BondedPoolName:    {supply.Burner, supply.Staking},
		validator.NotBondedPoolName: {supply.Burner, supply.Staking},
		swap.PoolName:               {supply.Minter, supply.Burner},
	}
)

// SimApp extends an ABCI application, but with most of its parameters exported.
// They are exported for convenience in creating helper functions, as object
// capabilities aren't needed for testing.
type SimApp struct {
	*bam.BaseApp
	cdc *codec.Codec

	invCheckPeriod uint

	// keys to access the substores
	keys  map[string]*sdk.KVStoreKey
	tkeys map[string]*sdk.TransientStoreKey

	// subspaces
	subspaces map[string]params.Subspace

	// keepers
	AccountKeeper   auth.AccountKeeper
	BankKeeper      bank.Keeper
	SupplyKeeper    supply.Keeper
	ValidatorKeeper validator.Keeper
	GovKeeper       gov.Keeper
	SwapKeeper      swap.Keeper
	ParamsKeeper    params.Keeper
	CoinKeeper      coin.Keeper
	MultisigKeeper  multisig.Keeper

	// the module manager
	mm *module.Manager

	// simulation manager
	sm *module.SimulationManager
}

// MakeCodec - custom tx codec
func MakeCodec() *codec.Codec {
	var cdc = codec.New()
	ModuleBasics.RegisterCodec(cdc)
	vesting.RegisterCodec(cdc)
	sdk.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)
	return cdc
}

// NewSimApp returns a reference to an initialized SimApp.
func NewSimApp(
	logger log.Logger, db dbm.DB, traceStore io.Writer, loadLatest bool, skipUpgradeHeights map[int64]bool,
	invCheckPeriod uint, baseAppOptions ...func(*bam.BaseApp),
) *SimApp {

	cdc := MakeCodec()

	bApp := bam.NewBaseApp("SimApp", logger, db, auth.DefaultTxDecoder(cdc), baseAppOptions...)
	bApp.SetCommitMultiStoreTracer(traceStore)
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

	app := &SimApp{
		BaseApp:        bApp,
		cdc:            cdc,
		invCheckPeriod: invCheckPeriod,
		keys:           keys,
		tkeys:          tkeys,
		subspaces:      make(map[string]params.Subspace),
	}

	// init params keeper and subspaces
	app.ParamsKeeper = params.NewKeeper(app.cdc, keys[params.StoreKey], tkeys[params.TStoreKey])
	app.subspaces[auth.ModuleName] = app.ParamsKeeper.Subspace(auth.DefaultParamspace)
	app.subspaces[bank.ModuleName] = app.ParamsKeeper.Subspace(bank.DefaultParamspace)
	app.subspaces[swap.ModuleName] = app.ParamsKeeper.Subspace(swap.DefaultParamspace)
	app.subspaces[multisig.ModuleName] = app.ParamsKeeper.Subspace(multisig.DefaultParamspace)
	app.subspaces[validator.ModuleName] = app.ParamsKeeper.Subspace(validator.DefaultParamSpace)
	app.subspaces[gov.ModuleName] = app.ParamsKeeper.Subspace(gov.DefaultParamspace).WithKeyTable(gov.ParamKeyTable())

	// add keepers
	app.AccountKeeper = auth.NewAccountKeeper(
		app.cdc, keys[auth.StoreKey], app.subspaces[auth.ModuleName], auth.ProtoBaseAccount,
	)
	app.BankKeeper = bank.NewBaseKeeper(
		app.AccountKeeper, app.subspaces[bank.ModuleName], app.BlacklistedAccAddrs(),
	)
	app.SupplyKeeper = supply.NewKeeper(
		app.cdc, keys[supply.StoreKey], app.AccountKeeper, app.BankKeeper, maccPerms,
	)
	app.CoinKeeper = coin.NewKeeper(
		app.cdc, keys[coin.StoreKey], app.subspaces[coin.ModuleName], app.AccountKeeper, app.BankKeeper, config.GetDefaultConfig(config.ChainID),
	)
	app.ValidatorKeeper = validator.NewKeeper(
		app.cdc, keys[validator.StoreKey], app.subspaces[validator.ModuleName], app.CoinKeeper, app.AccountKeeper, app.SupplyKeeper, app.MultisigKeeper, auth.FeeCollectorName,
	)

	// NOTE: Any module instantiated in the module manager that is later modified
	// must be passed by reference here.
	app.mm = module.NewManager(
		genutil.NewAppModule(app.AccountKeeper, app.ValidatorKeeper, app.BaseApp.DeliverTx),
		auth.NewAppModule(app.AccountKeeper),
		bank.NewAppModule(app.BankKeeper, app.AccountKeeper),
		supply.NewAppModule(app.SupplyKeeper, app.AccountKeeper),
		gov.NewAppModule(app.GovKeeper, app.AccountKeeper, app.SupplyKeeper),
		validator.NewAppModule(app.ValidatorKeeper, app.SupplyKeeper, app.CoinKeeper),
	)

	// During begin block slashing happens after distr.BeginBlocker so that
	// there is nothing left over in the validator fee pool, so as to keep the
	// CanWithdrawInvariant invariant.
	app.mm.SetOrderBeginBlockers(validator.ModuleName)
	app.mm.SetOrderEndBlockers(validator.ModuleName)

	// NOTE: The genutils moodule must occur after staking so that pools are
	// properly initialized with tokens from genesis accounts.
	app.mm.SetOrderInitGenesis(
		auth.ModuleName, bank.ModuleName,
		gov.ModuleName, supply.ModuleName,
		genutil.ModuleName, validator.ModuleName,
	)

	app.mm.RegisterRoutes(app.Router(), app.QueryRouter())

	// create the simulation manager and define the order of the modules for deterministic simulations
	//
	// NOTE: this is not required apps that don't use the simulator for fuzz testing
	// transactions
	app.sm = module.NewSimulationManager(
		auth.NewAppModule(app.AccountKeeper),
		bank.NewAppModule(app.BankKeeper, app.AccountKeeper),
		supply.NewAppModule(app.SupplyKeeper, app.AccountKeeper),
		gov.NewAppModule(app.GovKeeper, app.AccountKeeper, app.SupplyKeeper),
		params.NewAppModule(), // NOTE: only used for simulation to generate randomized param change proposals
	)

	app.sm.RegisterStoreDecoders()

	// initialize stores
	app.MountKVStores(keys)
	app.MountTransientStores(tkeys)

	// initialize BaseApp
	app.SetInitChainer(app.InitChainer)
	app.SetBeginBlocker(app.BeginBlocker)
	app.SetAnteHandler(ante.NewAnteHandler(app.AccountKeeper, app.SupplyKeeper, auth.DefaultSigVerificationGasConsumer))
	app.SetEndBlocker(app.EndBlocker)

	if loadLatest {
		err := app.LoadLatestVersion(app.keys[bam.MainStoreKey])
		if err != nil {
			tmos.Exit(err.Error())
		}
	}

	return app
}

// BlacklistedAccAddrs returns all the app's module account addresses black listed for receiving tokens.
func (app *SimApp) BlacklistedAccAddrs() map[string]bool {
	blacklistedAddrs := make(map[string]bool)

	return blacklistedAddrs
}

// Name returns the name of the App
func (app *SimApp) Name() string { return app.BaseApp.Name() }

// BeginBlocker application updates every begin block
func (app *SimApp) BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock) abci.ResponseBeginBlock {
	return app.mm.BeginBlock(ctx, req)
}

// EndBlocker application updates every end block
func (app *SimApp) EndBlocker(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
	return app.mm.EndBlock(ctx, req)
}

// InitChainer application update at chain initialization
func (app *SimApp) InitChainer(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
	var genesisState GenesisState
	app.cdc.MustUnmarshalJSON(req.AppStateBytes, &genesisState)
	return app.mm.InitGenesis(ctx, genesisState)
}

func (app *SimApp) Codec() *codec.Codec {
	return app.cdc
}

type GenesisState map[string]json.RawMessage

// NewDefaultGenesisState generates the default state for the application.
func NewDefaultGenesisState() GenesisState {
	return ModuleBasics.DefaultGenesis()
}

// LoadHeight loads a particular height
func (app *SimApp) LoadHeight(height int64) error {
	return app.LoadVersion(height, app.keys[bam.MainStoreKey])
}

// ModuleAccountAddrs returns all the app's module account addresses.
func (app *SimApp) ModuleAccountAddrs() map[string]bool {
	modAccAddrs := make(map[string]bool)
	for acc := range maccPerms {
		modAccAddrs[supply.NewModuleAddress(acc).String()] = true
	}

	return modAccAddrs
}

// returns context and app with params set on account keeper
func createTestApp(isCheckTx bool) (*SimApp, sdk.Context) {
	_config := sdk.GetConfig()
	_config.SetCoinType(60)
	_config.SetFullFundraiserPath("44'/60'/0'/0/0")
	_config.SetBech32PrefixForAccount(config.DecimalPrefixAccAddr, config.DecimalPrefixAccPub)
	_config.SetBech32PrefixForValidator(config.DecimalPrefixValAddr, config.DecimalPrefixValPub)
	_config.SetBech32PrefixForConsensusNode(config.DecimalPrefixConsAddr, config.DecimalPrefixConsPub)

	db := dbm.NewMemDB()
	app := NewSimApp(log.NewNopLogger(), db, nil, true, map[int64]bool{}, 0)
	if !isCheckTx {
		// init chain must be called to stop deliverState from being nil
		genesisState := NewDefaultGenesisState()
		stateBytes, err := codec.MarshalJSONIndent(app.Codec(), genesisState)
		if err != nil {
			panic(err)
		}

		// Initialize the chain
		app.InitChain(
			abci.RequestInitChain{
				Validators:    []abci.ValidatorUpdate{},
				AppStateBytes: stateBytes,
			},
		)
	}
	ctx := app.BaseApp.NewContext(isCheckTx, abci.Header{})
	app.AccountKeeper.SetParams(ctx, types.DefaultParams())

	return app, ctx
}

// Test that simulate transaction accurately estimates gas cost
func TestSimulateGasCost(t *testing.T) {
	// setup
	app, ctx := createTestApp(true)
	ctx = ctx.WithBlockHeight(updates.Update3Block)
	anteHandler := NewAnteHandler(app.AccountKeeper, app.ValidatorKeeper, app.CoinKeeper, app.SupplyKeeper, ante.DefaultSigVerificationGasConsumer)

	// keys and addresses
	priv1, _, addr1 := types.KeyTestPubAddr()

	swap.ServiceAddress = addr1.String()

	// set the accounts
	acc1 := app.AccountKeeper.NewAccountWithAddress(ctx, addr1)
	acc1.SetCoins(sdk.NewCoins(sdk.NewCoin(validator.DefaultBondDenom, sdk.NewInt(1))))
	require.NoError(t, acc1.SetAccountNumber(0))
	app.AccountKeeper.SetAccount(ctx, acc1)

	// set up msgs and fee
	var tx sdk.Tx
	msg1 := swap.MsgHTLT{
		TransferType: 1,
		From:         acc1.GetAddress(),
		Recipient:    "",
		HashedSecret: [32]byte{1, 2, 3, 4, 5, 6, 7, 8},
		Amount:       sdk.NewCoins(sdk.NewCoin(validator.DefaultBondDenom, sdk.NewInt(100))),
	}
	msgs := []sdk.Msg{msg1}
	fee := types.NewTestStdFee()

	// signers in order. accnums are all 0 because it is in genesis block
	privs, accnums, seqs := []crypto.PrivKey{priv1}, []uint64{0}, []uint64{0}
	tx = types.NewTestTx(ctx, msgs, privs, accnums, seqs, fee)

	cc, _ := ctx.CacheContext()
	newCtx, err := anteHandler(cc, tx, true)
	require.Nil(t, err, "transaction failed on simulate mode")

	simulatedGas := newCtx.GasMeter().GasConsumed()
	fee.Gas = simulatedGas

	// update tx with simulated gas estimate
	tx = types.NewTestTx(ctx, msgs, privs, accnums, seqs, fee)
	_, err = anteHandler(ctx, tx, false)

	require.Nil(t, err, "transaction failed with gas estimate")
}
