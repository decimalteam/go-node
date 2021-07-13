package nft

// DONTCOVER

import (
	"bitbucket.org/decimalteam/go-node/x/nft/client/cli"
	"bitbucket.org/decimalteam/go-node/x/nft/client/rest"
	"bitbucket.org/decimalteam/go-node/x/nft/simulation"
	"bitbucket.org/decimalteam/go-node/x/nft/types"
	authKeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	channeltypes "github.com/cosmos/ibc-go/modules/core/04-channel/types"
	ibcexported "github.com/cosmos/ibc-go/modules/core/exported"

	//"bitbucket.org/decimalteam/go-node/x/nft/simulation"
	"encoding/json"
	"github.com/cosmos/cosmos-sdk/client"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"math/rand"

	"github.com/gorilla/mux"
	"github.com/spf13/cobra"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	simulationtypes "github.com/cosmos/cosmos-sdk/types/simulation"
)

var (
	_ module.AppModule           = AppModule{}
	_ module.AppModuleBasic      = AppModuleBasic{}
	_ module.AppModuleSimulation = AppModule{}
)

// AppModuleBasic app module basics object
type AppModuleBasic struct {
	cdc codec.LegacyAmino
}

// var _ module.AppModuleBasic = AppModuleBasic{}

// Name defines module name
func (AppModuleBasic) Name() string {
	return ModuleName
}

// RegisterCodec registers module codec
func (AppModuleBasic) RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	RegisterCodec(cdc)
}

// DefaultGenesis default genesis state
func (AppModuleBasic) DefaultGenesis(_ codec.JSONCodec) json.RawMessage {
	return ModuleCdc.MustMarshalJSON(DefaultGenesisState())
}

// ValidateGenesis module validate genesis
func (AppModuleBasic) ValidateGenesis(marshaler codec.JSONCodec, config client.TxEncodingConfig, bz json.RawMessage) error {
	var data GenesisState
	err := ModuleCdc.UnmarshalJSON(bz, &data)
	if err != nil {
		return err
	}
	return ValidateGenesis(data)
}

// RegisterInterfaces implements InterfaceModule.RegisterInterfaces
func (AppModuleBasic) RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	types.RegisterInterfaces(registry)
}

// RegisterGRPCGatewayRoutes registers the gRPC Gateway routes for the nft module.
func (AppModuleBasic) RegisterGRPCGatewayRoutes(_ client.Context, _ *runtime.ServeMux) {}

// RegisterRESTRoutes registers rest routes
func (AppModuleBasic) RegisterRESTRoutes(ctx client.Context, rtr *mux.Router) {
	rest.RegisterRoutes(ctx, rtr, ModuleCdc, RouterKey)
}

// GetTxCmd gets the root tx command of this module
func (a AppModuleBasic) GetTxCmd() *cobra.Command {
	return cli.GetTxCmd(StoreKey, &a.cdc)
}

// GetQueryCmd gets the root query command of this module
func (a AppModuleBasic) GetQueryCmd() *cobra.Command {
	return cli.GetQueryCmd(RouterKey, &a.cdc)
}

//____________________________________________________________________________

// AppModule supply app module
type AppModule struct {
	AppModuleBasic

	keeper Keeper

	// Account keeper is used for testing purposes only
	accountKeeper authKeeper.AccountKeeper
}

// NewAppModule creates a new AppModule object
func NewAppModule(keeper Keeper, accountKeeper authKeeper.AccountKeeper) AppModule {
	return AppModule{
		AppModuleBasic: AppModuleBasic{},

		keeper:        keeper,
		accountKeeper: accountKeeper,
	}
}

// Name defines module name
func (AppModule) Name() string {
	return ModuleName
}

// RegisterInvariants registers the nft module invariants
func (am AppModule) RegisterInvariants(ir sdk.InvariantRegistry) {
	RegisterInvariants(ir, am.keeper)
}

// Route module message route name
func (AppModule) Route() sdk.Route {
	return sdk.NewRoute(RouterKey, func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		return &sdk.Result{}, nil
	})
}

func (AppModule) RegisterServices(module.Configurator) {}

// NewHandler module handler
func (am AppModule) NewHandler() sdk.Handler {
	return GenericHandler(am.keeper)
}

// QuerierRoute module querier route name
func (AppModule) QuerierRoute() string {
	return QuerierRoute
}

// NewQuerierHandler module querier
func (am AppModule) NewQuerierHandler() sdk.Querier {
	return NewQuerier(am.keeper)
}

func (am AppModule) LegacyQuerierHandler(*codec.LegacyAmino) sdk.Querier {
	return NewQuerier(am.keeper)
}

// InitGenesis module init-genesis
func (am AppModule) InitGenesis(ctx sdk.Context, _ codec.JSONCodec, data json.RawMessage) []abci.ValidatorUpdate {
	var genesisState GenesisState
	ModuleCdc.MustUnmarshalJSON(data, &genesisState)
	InitGenesis(ctx, am.keeper, genesisState)
	return []abci.ValidatorUpdate{}
}

// ExportGenesis module export genesis
func (am AppModule) ExportGenesis(ctx sdk.Context, _ codec.JSONCodec) json.RawMessage {
	gs := ExportGenesis(ctx, am.keeper)
	return ModuleCdc.MustMarshalJSON(gs)
}

// BeginBlock module begin-block
func (AppModule) BeginBlock(_ sdk.Context, _ abci.RequestBeginBlock) {}

// EndBlock module end-block
func (am AppModule) EndBlock(ctx sdk.Context, _ abci.RequestEndBlock) []abci.ValidatorUpdate {
	return []abci.ValidatorUpdate{}
}

// RegisterStoreDecoder registers a decoder for nft module's types
func (AppModule) RegisterStoreDecoder(sdr sdk.StoreDecoderRegistry) {
	//todo change it for correct decode store
	//sdr[StoreKey] = sim.NewDecodeStore(am.cdc)
}

// ProposalContents doesn't return any content functions for governance proposals.
func (AppModule) ProposalContents(_ module.SimulationState) []simulationtypes.WeightedProposalContent {
	return nil
}

// GenerateGenesisState creates a randomized GenState of the nft module.
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	simulation.RandomizedGenState(simState)
}

// RandomizedParams doesn't create randomized nft param changes for the simulator.
func (AppModule) RandomizedParams(_ *rand.Rand) []simulationtypes.ParamChange { return nil }

// WeightedOperations doesn't return any operation for the nft module.
func (am AppModule) WeightedOperations(_ module.SimulationState) []simulationtypes.WeightedOperation {
	return nil
	//return simulation.WeightedOperations(simState.AppParams, simState.Cdc, am.accountKeeper, am.keeper)
}

func (am AppModule) OnChanOpenInit(ctx sdk.Context, order channeltypes.Order, connectionHops []string, portID string, channelID string, channelCap *capabilitytypes.Capability, counterparty channeltypes.Counterparty, version string) error {
	return nil
}

func (am AppModule) OnChanOpenTry(ctx sdk.Context, order channeltypes.Order, connectionHops []string, portID, channelID string, channelCap *capabilitytypes.Capability, counterparty channeltypes.Counterparty, version, counterpartyVersion string) error {
	return nil
}

func (am AppModule) OnChanOpenAck(ctx sdk.Context, portID, channelID string, counterpartyVersion string) error {
	return nil
}

func (am AppModule) OnChanOpenConfirm(ctx sdk.Context, portID, channelID string) error {
	return nil
}

func (am AppModule) OnChanCloseInit(ctx sdk.Context, portID, channelID string) error {
	return nil
}

func (am AppModule) OnChanCloseConfirm(ctx sdk.Context, portID, channelID string) error {
	return nil
}

func (am AppModule) OnRecvPacket(ctx sdk.Context, packet channeltypes.Packet, relayer sdk.AccAddress) ibcexported.Acknowledgement {
	return nil
}

func (am AppModule) OnAcknowledgementPacket(ctx sdk.Context, packet channeltypes.Packet, acknowledgement []byte, relayer sdk.AccAddress) (*sdk.Result, error) {
	return &sdk.Result{}, nil
}

func (am AppModule) OnTimeoutPacket(ctx sdk.Context, packet channeltypes.Packet, relayer sdk.AccAddress) (*sdk.Result, error) {
	return &sdk.Result{}, nil
}

func (am AppModule) ConsensusVersion() uint64 {
	return 1
}
