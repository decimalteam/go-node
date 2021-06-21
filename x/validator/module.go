package validator

import (
	"encoding/json"
	"github.com/cosmos/cosmos-sdk/client"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	authKeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	bankKeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"

	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"

	abci "github.com/tendermint/tendermint/abci/types"
	cfg "github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/crypto"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"

	"bitbucket.org/decimalteam/go-node/x/coin"
	"bitbucket.org/decimalteam/go-node/x/validator/client/cli"
	"bitbucket.org/decimalteam/go-node/x/validator/client/rest"
)

// Type check to ensure the interface is properly implemented
var (
	_ module.AppModule      = AppModule{}
	_ module.AppModuleBasic = AppModuleBasic{}
)

// AppModuleBasic defines the basic application module used by the validator module.
type AppModuleBasic struct {
	cdc codec.LegacyAmino
}

// Name returns the validator module's name.
func (AppModuleBasic) Name() string {
	return ModuleName
}

// RegisterLegacyAminoCodec registers the validator module's types for the given codec.
func (AppModuleBasic) RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	RegisterCodec(cdc)
}

// DefaultGenesis returns default genesis state as raw bytes for the validator
// module.
func (AppModuleBasic) DefaultGenesis(_ codec.JSONMarshaler) json.RawMessage {
	return ModuleCdc.MustMarshalJSON(DefaultGenesisState())
}

// ValidateGenesis performs genesis state validation for the validator module.
func (AppModuleBasic) ValidateGenesis(marshaler codec.JSONMarshaler, config client.TxEncodingConfig, bz json.RawMessage) error {
	var data GenesisState
	err := ModuleCdc.UnmarshalJSON(bz, &data)
	if err != nil {
		return err
	}
	return ValidateGenesis(data)
}

// RegisterRESTRoutes registers the REST routes for the validator module.
func (AppModuleBasic) RegisterRESTRoutes(ctx client.Context, rtr *mux.Router) {
	rest.RegisterRoutes(ctx, rtr)
}

// RegisterInterfaces implements InterfaceModule.RegisterInterfaces
func (AppModuleBasic) RegisterInterfaces(registry cdctypes.InterfaceRegistry) {}

// RegisterGRPCGatewayRoutes registers the gRPC Gateway routes for the validator module.
func (AppModuleBasic) RegisterGRPCGatewayRoutes(clientCtx client.Context, mux *runtime.ServeMux) {}

// GetTxCmd returns the root tx command for the validator module.
func (amb AppModuleBasic) GetTxCmd() *cobra.Command {
	return cli.GetTxCmd(&amb.cdc)
}

// GetQueryCmd returns no root query command for the validator module.
func (amb AppModuleBasic) GetQueryCmd() *cobra.Command {
	return cli.GetQueryCmd(StoreKey, &amb.cdc)
}

//_____________________________________
// extra helpers

// CreateValidatorMsgHelpers - used for gen-tx
func (AppModuleBasic) CreateValidatorMsgHelpers(ipDefault string) (
	fs *flag.FlagSet, nodeIDFlag, pubkeyFlag, amountFlag, defaultsDesc string) {
	return cli.CreateValidatorMsgHelpers(ipDefault)
}

// PrepareFlagsForTxCreateValidator - used for gen-tx
func (AppModuleBasic) PrepareFlagsForTxCreateValidator(config *cfg.Config, nodeID,
	chainID string, valPubKey crypto.PubKey) {
	cli.PrepareFlagsForTxCreateValidator(config, nodeID, chainID, valPubKey)
}

// BuildCreateValidatorMsg - used for gen-tx
func (AppModuleBasic) BuildCreateValidatorMsg(
	cliCtx client.Context,
	txBldr client.TxBuilder,
) (client.TxBuilder, sdk.Msg, error) {
	return cli.BuildCreateValidatorMsg(cliCtx, txBldr)
}

// AppModule implements an application module for the validator module.
type AppModule struct {
	AppModuleBasic

	keeper     Keeper
	accKeeper  authKeeper.AccountKeeper
	coinKeeper coin.Keeper
	bankKeeper bankKeeper.BaseKeeper
}

// NewAppModule creates a new AppModule object
func NewAppModule(k Keeper, accKeeper authKeeper.AccountKeeper, bankKeeper bankKeeper.BaseKeeper, coinKeeper coin.Keeper) AppModule {

	return AppModule{
		AppModuleBasic: AppModuleBasic{},
		keeper:         k,
		accKeeper:      accKeeper,
		bankKeeper:     bankKeeper,
		coinKeeper:     coinKeeper,
	}
}

// Name returns the validator module's name.
func (AppModule) Name() string {
	return ModuleName
}

// Route returns the message routing key for the validator module.
func (AppModule) Route() sdk.Route {
	return sdk.NewRoute(RouterKey, func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		return &sdk.Result{}, nil
	})
}

// RegisterInvariants registers the validator module invariants.
func (AppModule) RegisterInvariants(_ sdk.InvariantRegistry) {}

// RegisterServices registers module services.
func (AppModule) RegisterServices(_ module.Configurator) {}

// LegacyQuerierHandler returns no sdk.Querier.
func (am AppModule) LegacyQuerierHandler(amino *codec.LegacyAmino) sdk.Querier {
	return NewQuerier(am.keeper)
}

// NewHandler returns an sdk.Handler for the validator module.
func (am AppModule) NewHandler() sdk.Handler {
	return NewHandler(am.keeper)
}

// QuerierRoute returns the validator module's querier route name.
func (AppModule) QuerierRoute() string {
	return QuerierRoute
}

// NewQuerierHandler returns the validator module sdk.Querier.
func (am AppModule) NewQuerierHandler() sdk.Querier {
	return NewQuerier(am.keeper)
}

// InitGenesis performs genesis initialization for the validator module. It returns
// no validator updates.
func (am AppModule) InitGenesis(ctx sdk.Context, _ codec.JSONMarshaler, data json.RawMessage) []abci.ValidatorUpdate {
	var genesisState GenesisState
	ModuleCdc.MustUnmarshalJSON(data, &genesisState)
	return InitGenesis(ctx, am.accKeeper, am.keeper, am.bankKeeper, genesisState)
}

// ExportGenesis returns the exported genesis state as raw bytes for the validator
// module.
func (am AppModule) ExportGenesis(ctx sdk.Context, marshaler codec.JSONMarshaler) json.RawMessage {
	gs := ExportGenesis(ctx, am.keeper)
	return ModuleCdc.MustMarshalJSON(gs)
}

// BeginBlock returns the begin blocker for the validator module.
func (am AppModule) BeginBlock(ctx sdk.Context, req abci.RequestBeginBlock) {
	BeginBlocker(ctx, req, am.keeper)
}

// EndBlock returns the end blocker for the validator module. It returns no validator
// updates.
func (am AppModule) EndBlock(ctx sdk.Context, _ abci.RequestEndBlock) []abci.ValidatorUpdate {
	return EndBlocker(ctx, am.keeper, am.coinKeeper, am.accKeeper, am.bankKeeper, true)
}

func (am AppModule) ConsensusVersion() uint64 {
	return 1
}
