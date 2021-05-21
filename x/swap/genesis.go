package swap

import (
	types2 "bitbucket.org/decimalteam/go-node/x/swap/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func ValidateGenesis(data types2.GenesisState) error {
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

func ExportGenesis(ctx sdk.Context, k Keeper) types2.GenesisState {
	params := k.GetParams(ctx)
	swaps := k.GetAllSwaps(ctx)
	return types2.GenesisState{
		Swaps:  swaps,
		Params: params,
	}
}
