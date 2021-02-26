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
func (k Keeper) GetNFT(ctx sdk.Context, denom, id string) (nft exported.NFT, err error) {
	collection, found := k.GetCollection(ctx, denom)
	if !found {
		return nil, sdkerrors.Wrap(types.ErrUnknownCollection, fmt.Sprintf("collection of %s doesn't exist", denom))
	}
	nft, err = collection.GetNFT(id)

	if err != nil {
		return nil, err
	}
	return nft, err
}

// MintNFT mints an NFT and manages that NFTs existence within Collections and Owners
func (k Keeper) MintNFT(ctx sdk.Context, denom string, nft exported.NFT) (err error) {
	collection, found := k.GetCollection(ctx, denom)
	if found {
		collection, err = collection.AddNFT(nft)
		if err != nil {
			return err
		}
	} else {
		collection = types.NewCollection(denom, types.NewNFTs(nft))
	}
	k.SetCollection(ctx, denom, collection)

	ownerIDCollection, _ := k.GetOwnerByDenom(ctx, nft.GetCreator(), denom)
	ownerIDCollection = ownerIDCollection.AddID(nft.GetID())
	k.SetOwnerByDenom(ctx, nft.GetCreator(), denom, ownerIDCollection.IDs)
	return
}

// DeleteNFT deletes an existing NFT from store
func (k Keeper) DeleteNFT(ctx sdk.Context, denom, id string, quantity sdk.Int) (err error) {
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

	nft.GetOwners().GetOwner(nft.GetCreator()).SetQuantity(nft.GetOwners().GetOwner(nft.GetCreator()).GetQuantity().Sub(quantity))

	k.SetCollection(ctx, denom, collection)

	return
}
