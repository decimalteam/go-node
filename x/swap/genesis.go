package swap

import (
	"bitbucket.org/decimalteam/go-node/x/swap/internal/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func ValidateGenesis(data types.GenesisState) error {
	err := data.Params.Validate()
	if err != nil {
		return err
	}
	return nil
}

func InitGenesis(ctx sdk.Context, k Keeper, data GenesisState) {
	k.SetParams(ctx, data.Params)

	for _, swap := range data.Swaps {
		k.SetSwap(ctx, swap)
	}
}

func ExportGenesis(ctx sdk.Context, k Keeper) types.GenesisState {
	params := k.GetParams(ctx)
	swaps := k.GetAllSwaps(ctx)
	return types.GenesisState{
		Swaps:  swaps,
		Params: params,
	}
}
