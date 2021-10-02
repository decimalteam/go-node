package keeper

import (
	"bitbucket.org/decimalteam/go-node/x/swap/internal/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) SetSwapV2(ctx sdk.Context, hash types.Hash) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetSwapV2Key(hash), []byte{})
}

func (k Keeper) HasSwapV2(ctx sdk.Context, hash types.Hash) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.GetSwapV2Key(hash))
}
