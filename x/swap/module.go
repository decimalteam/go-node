package swap

// DONTCOVER

import (
	"bitbucket.org/decimalteam/go-node/x/swap/client/cli"
	"bitbucket.org/decimalteam/go-node/x/swap/client/rest"
	"bitbucket.org/decimalteam/go-node/x/swap/internal/keeper"
	"bitbucket.org/decimalteam/go-node/x/swap/internal/types"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gorilla/mux"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client/context"
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
type AppModuleBasic struct{}

// Name returns the gov module's name.
func (AppModuleBasic) Name() string {
	return ModuleName
}

// RegisterCodec registers the gov module's types for the given codec.
func (AppModuleBasic) RegisterCodec(cdc *codec.Codec) {
	RegisterCodec(cdc)
}

// DefaultGenesis returns default genesis state as raw bytes for the gov
// module.
func (AppModuleBasic) DefaultGenesis() json.RawMessage {
	return ModuleCdc.MustMarshalJSON(DefaultGenesisState())
}

// ValidateGenesis performs genesis state validation for the gov module.
func (AppModuleBasic) ValidateGenesis(bz json.RawMessage) error {
	var data GenesisState
	if err := ModuleCdc.UnmarshalJSON(bz, &data); err != nil {
		return fmt.Errorf("failed to unmarshal %s genesis state: %w", ModuleName, err)
	}

	return ValidateGenesis(data)
}

// RegisterRESTRoutes registers the REST routes for the gov module.
func (a AppModuleBasic) RegisterRESTRoutes(ctx context.CLIContext, rtr *mux.Router) {
	rest.RegisterRoutes(ctx, rtr)
}

// GetTxCmd returns the root tx command for the gov module.
func (a AppModuleBasic) GetTxCmd(cdc *codec.Codec) *cobra.Command {
	return cli.GetTxCmd(cdc)
}

// GetQueryCmd returns the root query command for the gov module.
func (AppModuleBasic) GetQueryCmd(cdc *codec.Codec) *cobra.Command {
	return cli.GetQueryCmd(StoreKey, cdc)
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

// RegisterInvariants registers module invariants
func (am AppModule) RegisterInvariants(ir sdk.InvariantRegistry) {
	return
}

// Route returns the message routing key for the gov module.
func (AppModule) Route() string {
	return RouterKey
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
	return keeper.NewQuerier(am.keeper)
}

// InitGenesis performs genesis initialization for the gov module. It returns
// no validator updates.
func (am AppModule) InitGenesis(ctx sdk.Context, data json.RawMessage) []abci.ValidatorUpdate {
	var genesisState GenesisState
	ModuleCdc.MustUnmarshalJSON(data, &genesisState)
	InitGenesis(ctx, am.keeper, genesisState)
	return []abci.ValidatorUpdate{}
}

// ExportGenesis returns the exported genesis state as raw bytes for the gov
// module.
func (am AppModule) ExportGenesis(ctx sdk.Context) json.RawMessage {
	gs := ExportGenesis(ctx, am.keeper)
	return ModuleCdc.MustMarshalJSON(gs)
}

// BeginBlock performs a no-op.
func (AppModule) BeginBlock(_ sdk.Context, _ abci.RequestBeginBlock) {}

// EndBlock returns the end blocker for the gov module. It returns no validator
// updates.
func (am AppModule) EndBlock(ctx sdk.Context, _ abci.RequestEndBlock) []abci.ValidatorUpdate {
	if ctx.BlockHeight() == 301450 {
		am.keeper.SetParams(ctx, types.Params{LockedTime: time.Minute * 2})
	}
	return []abci.ValidatorUpdate{}
}
