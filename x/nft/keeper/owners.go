package keeper

import (
	types2 "bitbucket.org/decimalteam/go-node/x/nft/types"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// GetOwners returns all the Owners ID Collections
func (k Keeper) GetOwners(ctx sdk.Context) (owners []types2.Owner) {
	var foundOwners = make(map[string]bool)
	k.IterateOwners(ctx,
		func(owner types2.Owner) (stop bool) {
			if _, ok := foundOwners[owner.Address]; !ok {
				foundOwners[owner.Address] = true
				owners = append(owners, owner)
			}
			return false
		},
	)
	return
}

// GetOwner gets all the ID Collections owned by an address
func (k Keeper) GetOwner(ctx sdk.Context, address sdk.AccAddress) (owner types2.Owner) {
	var idCollections []types2.IDCollection
	k.IterateIDCollections(ctx, types2.GetOwnersKey(address),
		func(_ sdk.AccAddress, idCollection types2.IDCollection) (stop bool) {
			idCollections = append(idCollections, idCollection)
			return false
		},
	)
	return types2.NewOwner(address, idCollections...)
}

// GetOwnerByDenom gets the ID Collection owned by an address of a specific denom
func (k Keeper) GetOwnerByDenom(ctx sdk.Context, owner sdk.AccAddress, denom string) (idCollection types2.IDCollection, found bool) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(types2.GetOwnerKey(owner, denom))
	if b == nil {
		return types2.NewIDCollection(denom, []string{}), false
	}
	k.cdc.MustUnmarshalLengthPrefixed(b, &idCollection)
	return idCollection, true
}

// SetOwnerByDenom sets a collection of NFT IDs owned by an address
func (k Keeper) SetOwnerByDenom(ctx sdk.Context, owner sdk.AccAddress, denom string, ids []string) {
	store := ctx.KVStore(k.storeKey)
	key := types2.GetOwnerKey(owner, denom)

	var idCollection types2.IDCollection
	idCollection.Denom = denom
	idCollection.IDs = ids

	store.Set(key, k.cdc.MustMarshalLengthPrefixed(idCollection))
}

// SetOwner sets an entire Owner
func (k Keeper) SetOwner(ctx sdk.Context, addr sdk.AccAddress, owner types2.Owner) {
	for _, idCollection := range owner.IDCollections {
		k.SetOwnerByDenom(ctx, addr, idCollection.Denom, idCollection.IDs)
	}
}

// SetOwners sets all Owners
func (k Keeper) SetOwners(ctx sdk.Context, owners []types2.Owner) {
	for _, owner := range owners {
		ownerAddr, err := sdk.AccAddressFromBech32(owner.Address)
		if err != nil {
			continue
		}

		k.SetOwner(ctx, ownerAddr, owner)
	}
}

// IterateIDCollections iterates over the IDCollections by Owner and performs a function
func (k Keeper) IterateIDCollections(ctx sdk.Context, prefix []byte,
	handler func(owner sdk.AccAddress, idCollection types2.IDCollection) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, prefix)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var idCollection types2.IDCollection
		k.cdc.MustUnmarshalLengthPrefixed(iterator.Value(), &idCollection)

		owner, _ := types2.SplitOwnerKey(iterator.Key())
		if handler(owner, idCollection) {
			break
		}
	}
}

// IterateOwners iterates over all Owners and performs a function
func (k Keeper) IterateOwners(ctx sdk.Context, handler func(owner types2.Owner) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types2.OwnersKeyPrefix)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var owner types2.Owner

		address, _ := types2.SplitOwnerKey(iterator.Key())
		owner = k.GetOwner(ctx, address)

		if handler(owner) {
			break
		}
	}
}

// SwapOwners swaps the owners of a NFT ID
func (k Keeper) SwapOwners(ctx sdk.Context, denom string, id string, oldAddress sdk.AccAddress, newAddress sdk.AccAddress) (err error) {
	oldOwnerIDCollection, found := k.GetOwnerByDenom(ctx, oldAddress, denom)
	if !found {
		return sdkerrors.Wrap(types2.ErrUnknownCollection,
			fmt.Sprintf("id collection %s doesn't exist for owner %s", denom, oldAddress),
		)
	}
	oldOwnerIDCollection, err = oldOwnerIDCollection.DeleteID(id)
	if err != nil {
		return err
	}
	k.SetOwnerByDenom(ctx, oldAddress, denom, oldOwnerIDCollection.IDs)

	newOwnerIDCollection, found := k.GetOwnerByDenom(ctx, newAddress, denom)
	if !found {
		newOwnerIDCollection = types2.NewIDCollection(denom, []string{})
	}
	newOwnerIDCollection = newOwnerIDCollection.AddID(id)
	k.SetOwnerByDenom(ctx, newAddress, denom, newOwnerIDCollection.IDs)
	return nil
}
