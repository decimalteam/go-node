package swap

import (
	swaptypes "bitbucket.org/decimalteam/go-node/x/swap/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func ValidateGenesis(data swaptypes.GenesisState) error {
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

func ExportGenesis(ctx sdk.Context, k Keeper) *swaptypes.GenesisState {
	params := k.GetParams(ctx)
	swaps := k.GetAllSwaps(ctx)
	return &swaptypes.GenesisState{
		Swaps:  swaps,
		Params: params,
	}
}
