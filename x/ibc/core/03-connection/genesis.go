package connection

import (
	"bitbucket.org/decimalteam/go-node/x/ibc/core/03-connection/keeper"
	"bitbucket.org/decimalteam/go-node/x/ibc/core/03-connection/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// InitGenesis initializes the ibc connection submodule's state from a provided genesis
// state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, gs types.GenesisState) {
	for _, connection := range gs.Connections {
		conn := types.NewConnectionEnd(connection.State, connection.ClientId, connection.Counterparty, connection.Versions, connection.DelayPeriod)
		k.SetConnection(ctx, connection.Id, conn)
	}
	for _, connPaths := range gs.ClientConnectionPaths {
		k.SetClientConnectionPaths(ctx, connPaths.ClientId, connPaths.Paths)
	}
	k.SetNextConnectionSequence(ctx, gs.NextConnectionSequence)
}

// ExportGenesis returns the ibc connection submodule's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) types.GenesisState {
	return types.GenesisState{
		Connections:           k.GetAllConnections(ctx),
		ClientConnectionPaths: k.GetAllClientConnectionPaths(ctx),
	}
}
