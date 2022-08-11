package keeper

import (
	"bitbucket.org/decimalteam/go-node/x/swap/internal/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) HasChain(ctx sdk.Context, chainNumber int) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.GetChainKey(chainNumber))
}

func (k Keeper) SetChain(ctx sdk.Context, chainNumber int, chain types.Chain) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(types.NewChain(chain.Name, chain.Active))
	store.Set(types.GetChainKey(chainNumber), bz)
}

func (k Keeper) GetChain(ctx sdk.Context, chainNumber int) (types.Chain, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetChainKey(chainNumber))
	if bz == nil {
		return types.Chain{}, false
	}

	var chain types.Chain
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &chain)
	return chain, true
}

func (k Keeper) GetAllChains(ctx sdk.Context) []types.Chain {
	var chains []types.Chain
	iterator := k.GetChainsIterator(ctx)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var chain types.Chain
		err := k.cdc.UnmarshalBinaryLengthPrefixed(iterator.Value(), &chain)
		if err != nil {
			panic(err)
		}
		chains = append(chains, chain)
	}
	return chains
}

// GetChainsIterator gets an iterator over all chains
func (k Keeper) GetChainsIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, types.ChainKey)
}
