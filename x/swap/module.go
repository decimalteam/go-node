package swap

// DONTCOVER

import (
	"bitbucket.org/decimalteam/go-node/x/swap/client/cli"
	"bitbucket.org/decimalteam/go-node/x/swap/client/rest"
	keeper2 "bitbucket.org/decimalteam/go-node/x/swap/keeper"
	"bitbucket.org/decimalteam/go-node/x/swap/types"
	"encoding/json"
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	abci "github.com/tendermint/tendermint/abci/types"
)

var (
	_ module.AppModule      = AppModule{}
	_ module.AppModuleBasic = AppModuleBasic{}
)

// AppModuleBasic defines the basic application module used by the gov module.
type AppModuleBasic struct {
	cdc *codec.LegacyAmino
}

// Name returns the gov module's name.
func (AppModuleBasic) Name() string {
	return ModuleName
}

// RegisterCodec registers the gov module's types for the given codec.
func (AppModuleBasic) RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	RegisterCodec(cdc)
}

// RegisterInterfaces implements InterfaceModule.RegisterInterfaces
func (a AppModuleBasic) RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	types.RegisterInterfaces(registry)
}

// DefaultGenesis returns default genesis state as raw bytes for the gov
// module.
func (AppModuleBasic) DefaultGenesis(cdc codec.JSONMarshaler) json.RawMessage {
	return cdc.MustMarshalJSON(DefaultGenesisState())
}

// ValidateGenesis performs genesis state validation for the gov module.
func (AppModuleBasic) ValidateGenesis(cdc codec.JSONMarshaler, _ client.TxEncodingConfig, rawdata json.RawMessage) error {
	var data GenesisState
	if err := cdc.UnmarshalJSON(rawdata, &data); err != nil {
		return fmt.Errorf("failed to unmarshal %s genesis state: %w", ModuleName, err)
	}

	return ValidateGenesis(data)
}

// RegisterRESTRoutes registers the REST routes for the swap module.
func (a AppModuleBasic) RegisterRESTRoutes(ctx client.Context, rtr *mux.Router) {
	rest.RegisterRoutes(ctx, rtr)
}

func (a AppModuleBasic) RegisterGRPCGatewayRoutes(client.Context, *runtime.ServeMux) {}

// GetTxCmd returns the root tx command for the gov module.
func (a AppModuleBasic) GetTxCmd() *cobra.Command {
	return cli.GetTxCmd()
}

// GetQueryCmd returns the root query command for the gov module.
func (a AppModuleBasic) GetQueryCmd() *cobra.Command {
	return cli.GetQueryCmd(StoreKey, a.cdc)
}

//____________________________________________________________________________

// AppModule implements an application module for the gov module.
type AppModule struct {
	AppModuleBasic

	keeper Keeper
}

// NewAppModule creates a new AppModule object
func NewAppModule(keeper Keeper) AppModule {
	return AppModule{
		AppModuleBasic: AppModuleBasic{},
		keeper:         keeper,
	}
}

// Name returns the gov module's name.
func (AppModule) Name() string {
	return ModuleName
}

// Route returns the message routing key for the gov module.
func (am AppModule) Route() sdk.Route {
	return sdk.NewRoute(RouterKey, NewHandler(am.keeper))
}

// RegisterInvariants registers module invariants
func (am AppModule) RegisterInvariants(sdk.InvariantRegistry) {}

// RegisterServices registers module services.
func (am AppModule) RegisterServices(module.Configurator) {
}

// NewHandler returns an sdk.Handler for the gov module.
func (am AppModule) NewHandler() sdk.Handler {
	return NewHandler(am.keeper)
}

// QuerierRoute returns the gov module's querier route name.
func (AppModule) QuerierRoute() string {
	return QuerierRoute
}

// NewQuerierHandler returns no sdk.Querier.
func (am AppModule) NewQuerierHandler() sdk.Querier {
	return nil
}

// LegacyQuerierHandler returns no sdk.Querier.
func (am AppModule) LegacyQuerierHandler(*codec.LegacyAmino) sdk.Querier {
	return keeper2.NewQuerier(am.keeper)
}

// InitGenesis performs genesis initialization for the gov module. It returns
// no validator updates.
func (am AppModule) InitGenesis(ctx sdk.Context, cdc codec.JSONMarshaler, data json.RawMessage) []abci.ValidatorUpdate {
	var genesisState GenesisState
	cdc.MustUnmarshalJSON(data, &genesisState)
	InitGenesis(ctx, am.keeper, genesisState)
	return []abci.ValidatorUpdate{}
}

// ExportGenesis returns the exported genesis state as raw bytes for the gov
// module.
func (am AppModule) ExportGenesis(ctx sdk.Context, cdc codec.JSONMarshaler) json.RawMessage {
	gs := ExportGenesis(ctx, am.keeper)
	return cdc.MustMarshalJSON(gs)
}

// BeginBlock performs a no-op.
func (AppModule) BeginBlock(_ sdk.Context, _ abci.RequestBeginBlock) {}

// EndBlock returns the end blocker for the gov module. It returns no validator
// updates.
func (am AppModule) EndBlock(_ sdk.Context, _ abci.RequestEndBlock) []abci.ValidatorUpdate {
	return []abci.ValidatorUpdate{}
}

func (am AppModule) ConsensusVersion() uint64 {
	return 1
}
