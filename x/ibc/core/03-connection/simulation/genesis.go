package simulation

import (
	"math/rand"

	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"bitbucket.org/decimalteam/go-node/x/ibc/core/03-connection/types"
)

// GenConnectionGenesis returns the default connection genesis state.
func GenConnectionGenesis(_ *rand.Rand, _ []simtypes.Account) types.GenesisState {
	return types.DefaultGenesisState()
}
