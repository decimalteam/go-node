package keeper

import (
	types2 "bitbucket.org/decimalteam/go-node/x/swap/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) SetSwap(ctx sdk.Context, swap types2.Swap) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalLengthPrefixed(swap)
	store.Set(types2.GetSwapKey(*swap.HashedSecret), bz)
}

func (k Keeper) HasSwap(ctx sdk.Context, hash types2.Hash) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types2.GetSwapKey(hash))
}

func (k Keeper) GetSwap(ctx sdk.Context, hash types2.Hash) (types2.Swap, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types2.GetSwapKey(hash))
	if bz == nil {
		return types2.Swap{}, false
	}

	var swap types2.Swap
	k.cdc.MustUnmarshalLengthPrefixed(bz, &swap)
	return swap, true
}

func (k Keeper) GetAllSwaps(ctx sdk.Context) types2.Swaps {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types2.SwapKey)
	defer iterator.Close()

	var swaps types2.Swaps

	for ; iterator.Valid(); iterator.Next() {
		var swap types2.Swap
		k.cdc.MustUnmarshalLengthPrefixed(iterator.Value(), &swap)
		swaps = append(swaps, swap)
	}

	return swaps
}
