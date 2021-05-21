package capability

import (
	"bitbucket.org/decimalteam/go-node/x/capability/internal/types"
	keeper2 "bitbucket.org/decimalteam/go-node/x/capability/keeper"
	types2 "bitbucket.org/decimalteam/go-node/x/capability/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// InitGenesis initializes the capability module's state from a provided genesis
// state.
func InitGenesis(ctx sdk.Context, k keeper2.Keeper, genState types2.GenesisState) {
	if err := k.InitializeIndex(ctx, genState.Index); err != nil {
		panic(err)
	}

	// set owners for each index and initialize capability
	for _, genOwner := range genState.Owners {
		k.SetOwners(ctx, genOwner.Index, genOwner.IndexOwners)
		k.InitializeCapability(ctx, genOwner.Index, genOwner.IndexOwners)
	}
}

// ExportGenesis returns the capability module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper2.Keeper) *types2.GenesisState {
	index := k.GetLatestIndex(ctx)
	owners := []types.GenesisOwners{}

	for i := uint64(1); i < index; i++ {
		capabilityOwners, ok := k.GetOwners(ctx, i)
		if !ok || len(capabilityOwners.Owners) == 0 {
			continue
		}

		genOwner := types.GenesisOwners{
			Index:       i,
			IndexOwners: capabilityOwners,
		}
		owners = append(owners, genOwner)
	}

	return &types2.GenesisState{
		Index:  index,
		Owners: owners,
	}
}
