package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"bitbucket.org/decimalteam/go-node/x/swap/internal/types"
)

func (k Keeper) SetSwap(ctx sdk.Context, swap types.Swap) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(swap)
	store.Set(types.GetSwapKey(swap.HashedSecret), bz)
}

func (k Keeper) HasSwap(ctx sdk.Context, hash types.Hash) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.GetSwapKey(hash))
}

func (k Keeper) GetSwap(ctx sdk.Context, hash types.Hash) (types.Swap, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetSwapKey(hash))
	if bz == nil {
		return types.Swap{}, false
	}

	var swap types.Swap
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &swap)
	return swap, true
}

func (k Keeper) GetAllSwaps(ctx sdk.Context) types.Swaps {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.SwapKey)
	defer iterator.Close()

	var swaps types.Swaps

	for ; iterator.Valid(); iterator.Next() {
		var swap types.Swap
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &swap)
		swaps = append(swaps, swap)
	}

	return swaps
}
