package simulation

import (
	"math/rand"

	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"bitbucket.org/decimalteam/go-node/x/ibc/core/04-channel/types"
)

// GenChannelGenesis returns the default channel genesis state.
func GenChannelGenesis(_ *rand.Rand, _ []simtypes.Account) types.GenesisState {
	return types.DefaultGenesisState()
}
