package nft

import (
	"bitbucket.org/decimalteam/go-node/utils/updates"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func BeginBlocker(ctx sdk.Context, k Keeper) {
	if ctx.BlockHeight() >= updates.Update10Block {
		k.SetBaseDenom()
	}
	if ctx.BlockHeight() == updates.Update11Block {
		collections := k.GetCollections(ctx)
		for _, collection := range collections {
			for _, nft := range collection.NFTs {
				k.SetTokenIDIndex(ctx, nft.GetID())
			}
		}
	}
}
