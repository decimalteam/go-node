package genutil

import (
	genutiltypes "bitbucket.org/decimalteam/go-node/x/genutil/types"
	"bitbucket.org/decimalteam/go-node/x/validator"
	"encoding/json"
	"github.com/cosmos/cosmos-sdk/client"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"

	"github.com/gorilla/mux"
	"github.com/spf13/cobra"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/genutil/types"
)

var (
	_ module.AppModuleGenesis = AppModule{}
	_ module.AppModuleBasic   = AppModuleBasic{}
)

// AppModuleBasic defines the basic application module used by the genutil module.
type AppModuleBasic struct {
	cdc codec.JSONMarshaler
}

// Name returns the genutil module's name.
func (AppModuleBasic) Name() string {
	return ModuleName
}

// RegisterLegacyAminoCodec registers the genutil module's types for the given codec.
func (AppModuleBasic) RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
}

// DefaultGenesis returns default genesis state as raw bytes for the genutil
// module.
func (AppModuleBasic) DefaultGenesis(cdc codec.JSONMarshaler) json.RawMessage {
	return cdc.MustMarshalJSON(&types.GenesisState{})
}

// ValidateGenesis performs genesis state validation for the genutil module.
func (AppModuleBasic) ValidateGenesis(cdc codec.JSONMarshaler, config client.TxEncodingConfig, bz json.RawMessage) error {
	var data genutiltypes.GenesisState
	err := cdc.UnmarshalJSON(bz, &data)
	if err != nil {
		return err
	}
	return ValidateGenesis(&data, config.TxJSONDecoder())
}

// RegisterInterfaces implements InterfaceModule.RegisterInterfaces
func (AppModuleBasic) RegisterInterfaces(registry cdctypes.InterfaceRegistry) {}

// RegisterGRPCGatewayRoutes registers the gRPC Gateway routes for the genutil module.
func (AppModuleBasic) RegisterGRPCGatewayRoutes(clientCtx client.Context, mux *runtime.ServeMux) {}

// RegisterRESTRoutes registers the REST routes for the genutil module.
func (amb AppModuleBasic) RegisterRESTRoutes(context client.Context, router *mux.Router) {
}

// GetTxCmd returns the root tx command for the genutil module.
func (amb AppModuleBasic) GetTxCmd() *cobra.Command {
	return nil
}

// GetQueryCmd returns no root query command for the genutil module.
func (amb AppModuleBasic) GetQueryCmd() *cobra.Command {
	return nil
}

//___________________________
// app module
type AppModule struct {
	AppModuleBasic
	accountKeeper   types.AccountKeeper
	validatorKeeper validator.Keeper
	deliverTx       deliverTxfn
	txConfig        client.TxEncodingConfig
}

// NewAppModule creates a new AppModule object
func NewAppModule(accountKeeper types.AccountKeeper,
	validatorKeeper validator.Keeper, deliverTx deliverTxfn, txConfig client.TxEncodingConfig) module.AppModule {

	return module.NewGenesisOnlyAppModule(AppModule{
		AppModuleBasic:  AppModuleBasic{},
		accountKeeper:   accountKeeper,
		validatorKeeper: validatorKeeper,
		deliverTx:       deliverTx,
		txConfig:        txConfig,
	})
}

// Name returns the genutil module's name.
func (AppModule) Name() string {
	return ModuleName
}

// Route returns the message routing key for the genutil module.
func (AppModule) Route() sdk.Route {
	return sdk.NewRoute(ModuleName, func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		return &sdk.Result{}, nil
	})
}

// RegisterInvariants registers the genutil module invariants.
func (AppModule) RegisterInvariants(_ sdk.InvariantRegistry) {}

// RegisterServices registers module services.
func (AppModule) RegisterServices(_ module.Configurator) {}

// LegacyQuerierHandler returns no sdk.Querier.
func (am AppModule) LegacyQuerierHandler(amino *codec.LegacyAmino) sdk.Querier {
	return nil
}

// NewHandler returns an sdk.Handler for the genutil module.
func (am AppModule) NewHandler() sdk.Handler {
	return nil
}

// QuerierRoute returns the genutil module's querier route name.
func (AppModule) QuerierRoute() string {
	return ModuleName
}

// NewQuerierHandler returns the genutil module sdk.Querier.
func (am AppModule) NewQuerierHandler() sdk.Querier {
	return nil
}

// InitGenesis performs genesis initialization for the genutil module. It returns
// no genutil updates.
func (am AppModule) InitGenesis(ctx sdk.Context, cdc codec.JSONMarshaler, data json.RawMessage) []abci.ValidatorUpdate {
	var genesisState genutiltypes.GenesisState
	cdc.MustUnmarshalJSON(data, &genesisState)
	return InitGenesis(ctx, cdc, am.validatorKeeper, am.deliverTx, genesisState, am.txConfig)
}

// module export genesis
func (am AppModule) ExportGenesis(sdk.Context, codec.JSONMarshaler) json.RawMessage {
	return nil
}
func (am AppModule) ConsensusVersion() uint64 {
	return 1
}
