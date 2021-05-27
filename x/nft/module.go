package nft

// DONTCOVER

import (
	"bitbucket.org/decimalteam/go-node/x/nft/client/cli"
	"bitbucket.org/decimalteam/go-node/x/nft/client/rest"
	"bitbucket.org/decimalteam/go-node/x/nft/internal/types"
	"bitbucket.org/decimalteam/go-node/x/nft/simulation"
	"encoding/json"
	"math/rand"

	"github.com/gorilla/mux"
	"github.com/spf13/cobra"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	sim "github.com/cosmos/cosmos-sdk/x/simulation"
)

var (
	_ module.AppModule           = AppModule{}
	_ module.AppModuleBasic      = AppModuleBasic{}
	_ module.AppModuleSimulation = AppModule{}
)

// AppModuleBasic app moduleKey basics object
type AppModuleBasic struct{}

var _ module.AppModuleBasic = AppModuleBasic{}

// Name defines moduleKey name
func (AppModuleBasic) Name() string {
	return ModuleName
}

// RegisterCodec registers moduleKey codec
func (AppModuleBasic) RegisterCodec(cdc *codec.Codec) {
	RegisterCodec(cdc)
}

// DefaultGenesis default genesis state
func (AppModuleBasic) DefaultGenesis() json.RawMessage {
	return ModuleCdc.MustMarshalJSON(DefaultGenesisState())
}

// ValidateGenesis moduleKey validate genesis
func (AppModuleBasic) ValidateGenesis(bz json.RawMessage) error {
	var data GenesisState
	err := ModuleCdc.UnmarshalJSON(bz, &data)
	if err != nil {
		return err
	}
	return ValidateGenesis(data)
}

// RegisterRESTRoutes registers rest routes
func (AppModuleBasic) RegisterRESTRoutes(ctx context.CLIContext, rtr *mux.Router) {
	rest.RegisterRoutes(ctx, rtr, ModuleCdc, RouterKey)
}

// GetTxCmd gets the root tx command of this moduleKey
func (AppModuleBasic) GetTxCmd(cdc *codec.Codec) *cobra.Command {
	return cli.GetTxCmd(StoreKey, cdc)
}

// GetQueryCmd gets the root query command of this moduleKey
func (AppModuleBasic) GetQueryCmd(cdc *codec.Codec) *cobra.Command {
	return cli.GetQueryCmd(RouterKey, cdc)
}

//____________________________________________________________________________

// AppModule supply app moduleKey
type AppModule struct {
	AppModuleBasic

	keeper Keeper

	// Account keeper is used for testing purposes only
	accountKeeper types.AccountKeeper
}

// NewAppModule creates a new AppModule object
func NewAppModule(keeper Keeper, accountKeeper types.AccountKeeper) AppModule {
	return AppModule{
		AppModuleBasic: AppModuleBasic{},

		keeper:        keeper,
		accountKeeper: accountKeeper,
	}
}

// Name defines moduleKey name
func (AppModule) Name() string {
	return ModuleName
}

// RegisterInvariants registers the nft moduleKey invariants
func (am AppModule) RegisterInvariants(ir sdk.InvariantRegistry) {
	RegisterInvariants(ir, am.keeper)
}

// Route moduleKey message route name
func (AppModule) Route() string {
	return RouterKey
}

// NewHandler moduleKey handler
func (am AppModule) NewHandler() sdk.Handler {
	return GenericHandler(am.keeper)
}

// QuerierRoute moduleKey querier route name
func (AppModule) QuerierRoute() string {
	return QuerierRoute
}

// NewQuerierHandler moduleKey querier
func (am AppModule) NewQuerierHandler() sdk.Querier {
	return NewQuerier(am.keeper)
}

// InitGenesis moduleKey init-genesis
func (am AppModule) InitGenesis(ctx sdk.Context, data json.RawMessage) []abci.ValidatorUpdate {
	var genesisState GenesisState
	ModuleCdc.MustUnmarshalJSON(data, &genesisState)
	InitGenesis(ctx, am.keeper, genesisState)
	return []abci.ValidatorUpdate{}
}

// ExportGenesis moduleKey export genesis
func (am AppModule) ExportGenesis(ctx sdk.Context) json.RawMessage {
	gs := ExportGenesis(ctx, am.keeper)
	return ModuleCdc.MustMarshalJSON(gs)
}

// BeginBlock moduleKey begin-block
func (AppModule) BeginBlock(_ sdk.Context, _ abci.RequestBeginBlock) {}

// EndBlock moduleKey end-block
func (am AppModule) EndBlock(ctx sdk.Context, _ abci.RequestEndBlock) []abci.ValidatorUpdate {
	return []abci.ValidatorUpdate{}
}

// RegisterStoreDecoder registers a decoder for nft moduleKey's types
func (AppModule) RegisterStoreDecoder(sdr sdk.StoreDecoderRegistry) {
	sdr[StoreKey] = simulation.DecodeStore
}

// ProposalContents doesn't return any content functions for governance proposals.
func (AppModule) ProposalContents(_ module.SimulationState) []sim.WeightedProposalContent { return nil }

// GenerateGenesisState creates a randomized GenState of the nft moduleKey.
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	simulation.RandomizedGenState(simState)
}

// RandomizedParams doesn't create randomized nft param changes for the simulator.
func (AppModule) RandomizedParams(_ *rand.Rand) []sim.ParamChange { return nil }

// WeightedOperations doesn't return any operation for the nft moduleKey.
func (am AppModule) WeightedOperations(simState module.SimulationState) []sim.WeightedOperation {
	return nil
	//return simulation.WeightedOperations(simState.AppParams, simState.Cdc, am.accountKeeper, am.keeper)
}
