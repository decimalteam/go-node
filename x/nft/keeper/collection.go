package keeper

import (
	types2 "bitbucket.org/decimalteam/go-node/x/nft/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// IterateCollections iterates over collections and performs a function
func (k Keeper) IterateCollections(ctx sdk.Context, handler func(collection types2.Collection) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types2.CollectionsKeyPrefix)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var collection types2.Collection
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &collection)
		if handler(collection) {
			break
		}
	}
}

// SetCollection sets the entire collection of a single denom
func (k Keeper) SetCollection(ctx sdk.Context, denom string, collection types2.Collection) {
	store := ctx.KVStore(k.storeKey)
	collectionKey := types2.GetCollectionKey(denom)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(collection)
	store.Set(collectionKey, bz)
}

// GetCollection returns a collection of NFTs
func (k Keeper) GetCollection(ctx sdk.Context, denom string) (collection types2.Collection, found bool) {
	store := ctx.KVStore(k.storeKey)
	collectionKey := types2.GetCollectionKey(denom)
	bz := store.Get(collectionKey)
	if bz == nil {
		return
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &collection)
	return collection, true
}

// GetCollections returns all the NFTs collections
func (k Keeper) GetCollections(ctx sdk.Context) (collections []types2.Collection) {
	k.IterateCollections(ctx,
		func(collection types2.Collection) (stop bool) {
			collections = append(collections, collection)
			return false
		},
	)
	return
}

// GetDenoms returns all the NFT denoms
func (k Keeper) GetDenoms(ctx sdk.Context) (denoms []string) {
	k.IterateCollections(ctx,
		func(collection types2.Collection) (stop bool) {
			denoms = append(denoms, collection.Denom)
			return false
		},
	)
	return
}
