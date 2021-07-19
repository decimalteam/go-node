package client

import (
	"bitbucket.org/decimalteam/go-node/x/ibc/core/02-client/keeper"
	"bitbucket.org/decimalteam/go-node/x/ibc/core/exported"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// BeginBlocker updates an existing localhost client with the latest block height.
func BeginBlocker(ctx sdk.Context, k keeper.Keeper) {
	_, found := k.GetClientState(ctx, exported.Localhost)
	if !found {
		return
	}

	// update the localhost client with the latest block height
	if err := k.UpdateClient(ctx, exported.Localhost, nil); err != nil {
		panic(err)
	}
}
