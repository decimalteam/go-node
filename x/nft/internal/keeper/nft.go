package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"bitbucket.org/decimalteam/go-node/x/nft/exported"
	"bitbucket.org/decimalteam/go-node/x/nft/internal/types"
)

// IsNFT returns whether an NFT exists
func (k Keeper) IsNFT(ctx sdk.Context, denom, id string) (exists bool) {
	_, err := k.GetNFT(ctx, denom, id)
	return err == nil
}

// GetNFT gets the entire NFT metadata struct for a uint64
func (k Keeper) GetNFT(ctx sdk.Context, denom, id string) (exported.NFT, error) {
	collection, found := k.GetCollection(ctx, denom)
	if !found {
		return nil, sdkerrors.Wrap(types.ErrUnknownCollection, fmt.Sprintf("collection of %s doesn't exist", denom))
	}
	nft, err := collection.GetNFT(id)

	if err != nil {
		return nil, err
	}
	return nft, err
}

func (k Keeper) GetSubToken(ctx sdk.Context, denom, id string, subTokenID sdk.Int) sdk.Int {
	store := ctx.KVStore(k.storeKey)
	subTokenKey := types.GetSubTokenKey(denom, id, subTokenID)
	bz := store.Get(subTokenKey)
	if bz == nil {
		return sdk.Int{}
	}

	subToken := sdk.ZeroInt()

	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &subToken)
	return subToken
}

func (k Keeper) SetSubToken(ctx sdk.Context, denom, id string, subTokenID sdk.Int, reserve sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	subTokenKey := types.GetSubTokenKey(denom, id, subTokenID)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(reserve)
	store.Set(subTokenKey, bz)
}

func (k Keeper) GetLastSubTokenID(ctx sdk.Context, denom, id string) sdk.Int {
	store := ctx.KVStore(k.storeKey)
	lastSubTokenIDKey := types.GetLastSubTokenIDKey(denom, id)
	bz := store.Get(lastSubTokenIDKey)
	if bz == nil {
		return sdk.ZeroInt()
	}

	lastTokenID := sdk.ZeroInt()

	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &lastTokenID)
	return lastTokenID
}

func (k Keeper) SetLastSubTokenID(ctx sdk.Context, denom, id string, lastSubTokenID sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	lastSubTokenIDKey := types.GetLastSubTokenIDKey(denom, id)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(lastSubTokenID.String())
	store.Set(lastSubTokenIDKey, bz)
}

// MintNFT mints an NFT and manages that NFTs existence within Collections and Owners
func (k Keeper) MintNFT(ctx sdk.Context, denom string, nft exported.NFT) (sdk.Int, error) {
	quantity := nft.GetOwners().GetOwners()[0].GetQuantity()

	err := k.ReserveTokens(ctx,
		sdk.NewCoins(
			sdk.NewCoin(
				k.baseDenom,
				nft.GetReserve().Mul(quantity), // reserve * quantity
			)),
		nft.GetCreator())
	if err != nil {
		return sdk.Int{}, err
	}

	collection, found := k.GetCollection(ctx, denom)
	if found {
		collection, err = collection.AddNFT(nft)
		if err != nil {
			return sdk.Int{}, err
		}
	} else {
		collection = types.NewCollection(denom, types.NewNFTs(nft))
	}
	k.SetCollection(ctx, denom, collection)

	lastSubTokenID := k.GetLastSubTokenID(ctx, denom, nft.GetID())

	newLastSubTokenID := lastSubTokenID.Add(quantity)

	for i := lastSubTokenID; i.LT(newLastSubTokenID); i = i.AddRaw(1) {
		k.SetSubToken(ctx, denom, nft.GetID(), i, nft.GetReserve())
	}

	k.SetLastSubTokenID(ctx, denom, nft.GetID(), newLastSubTokenID)

	ownerIDCollection, _ := k.GetOwnerByDenom(ctx, nft.GetCreator(), denom)
	ownerIDCollection = ownerIDCollection.AddID(nft.GetID())
	k.SetOwnerByDenom(ctx, nft.GetCreator(), denom, ownerIDCollection.IDs)
	return newLastSubTokenID, err
}

// DeleteNFT deletes an existing NFT from store
func (k Keeper) DeleteNFT(ctx sdk.Context, denom, id string, quantity sdk.Int) error {
	collection, found := k.GetCollection(ctx, denom)
	if !found {
		return sdkerrors.Wrap(types.ErrUnknownCollection, fmt.Sprintf("collection of %s doesn't exist", denom))
	}
	nft, err := collection.GetNFT(id)
	if err != nil {
		return err
	}
	ownerIDCollection, found := k.GetOwnerByDenom(ctx, nft.GetCreator(), denom)
	if !found {
		return sdkerrors.Wrap(types.ErrUnknownCollection,
			fmt.Sprintf("id collection #%s doesn't exist for owner %s", denom, nft.GetCreator()),
		)
	}

	if quantity.GT(nft.GetOwners().GetOwner(nft.GetCreator()).GetQuantity()) {
		return sdkerrors.Wrap(types.ErrNotAllowedBurn,
			fmt.Sprintf("owner %s has only %s tokens", nft.GetCreator(), nft.GetOwners().GetOwner(nft.GetCreator()).GetQuantity().String()))
	}

	if quantity.Equal(nft.GetOwners().GetOwner(nft.GetCreator()).GetQuantity()) {
		ownerIDCollection, err = ownerIDCollection.DeleteID(nft.GetID())
		if err != nil {
			return err
		}
		k.SetOwnerByDenom(ctx, nft.GetCreator(), denom, ownerIDCollection.IDs)
	}

	owner := nft.GetOwners().
		GetOwner(nft.GetCreator())

	nft = nft.SetOwners(nft.
		GetOwners().
		SetOwner(
			owner.
				SetQuantity(owner.GetQuantity().Sub(quantity))))

	collection, err = collection.DeleteNFT(nft)
	if err != nil {
		return err
	}

	k.SetCollection(ctx, denom, collection)

	err = k.BurnTokens(ctx, sdk.NewCoins(
		sdk.NewCoin(k.baseDenom, nft.GetReserve().Mul(quantity))))
	if err != nil {
		return err
	}

	return nil
}
