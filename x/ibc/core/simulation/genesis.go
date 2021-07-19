package simulation

// DONTCOVER

import (
	"encoding/json"
	"fmt"
	"math/rand"

	clientsims "bitbucket.org/decimalteam/go-node/x/ibc/core/02-client/simulation"
	clienttypes "bitbucket.org/decimalteam/go-node/x/ibc/core/02-client/types"
	connectionsims "bitbucket.org/decimalteam/go-node/x/ibc/core/03-connection/simulation"
	connectiontypes "bitbucket.org/decimalteam/go-node/x/ibc/core/03-connection/types"
	channelsims "bitbucket.org/decimalteam/go-node/x/ibc/core/04-channel/simulation"
	channeltypes "bitbucket.org/decimalteam/go-node/x/ibc/core/04-channel/types"
	host "bitbucket.org/decimalteam/go-node/x/ibc/core/24-host"
	"bitbucket.org/decimalteam/go-node/x/ibc/core/types"
	"github.com/cosmos/cosmos-sdk/types/module"
)

// Simulation parameter constants
const (
	clientGenesis     = "client_genesis"
	connectionGenesis = "connection_genesis"
	channelGenesis    = "channel_genesis"
)

// RandomizedGenState generates a random GenesisState for evidence
func RandomizedGenState(simState *module.SimulationState) {
	var (
		clientGenesisState     clienttypes.GenesisState
		connectionGenesisState connectiontypes.GenesisState
		channelGenesisState    channeltypes.GenesisState
	)

	simState.AppParams.GetOrGenerate(
		simState.Cdc, clientGenesis, &clientGenesisState, simState.Rand,
		func(r *rand.Rand) { clientGenesisState = clientsims.GenClientGenesis(r, simState.Accounts) },
	)

	simState.AppParams.GetOrGenerate(
		simState.Cdc, connectionGenesis, &connectionGenesisState, simState.Rand,
		func(r *rand.Rand) { connectionGenesisState = connectionsims.GenConnectionGenesis(r, simState.Accounts) },
	)

	simState.AppParams.GetOrGenerate(
		simState.Cdc, channelGenesis, &channelGenesisState, simState.Rand,
		func(r *rand.Rand) { channelGenesisState = channelsims.GenChannelGenesis(r, simState.Accounts) },
	)

	ibcGenesis := types.GenesisState{
		ClientGenesis:     clientGenesisState,
		ConnectionGenesis: connectionGenesisState,
		ChannelGenesis:    channelGenesisState,
	}

	bz, err := json.MarshalIndent(&ibcGenesis, "", " ")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Selected randomly generated %s parameters:\n%s\n", host.ModuleName, bz)
	simState.GenState[host.ModuleName] = simState.Cdc.MustMarshalJSON(&ibcGenesis)
}
