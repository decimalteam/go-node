package nft

import (
	"bitbucket.org/decimalteam/go-node/utils/updates"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func BeginBlocker(ctx sdk.Context, k Keeper) {
	if ctx.BlockHeight() >= updates.Update10Block {
		k.SetBaseDenom()
	}
}
