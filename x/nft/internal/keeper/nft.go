package keeper

import (
	"encoding/binary"
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

func (k Keeper) GetSubToken(ctx sdk.Context, denom, id string, subTokenID int64) (sdk.Int, bool) {
	store := ctx.KVStore(k.storeKey)
	subTokenKey := types.GetSubTokenKey(denom, id, subTokenID)
	bz := store.Get(subTokenKey)
	if bz == nil {
		return sdk.Int{}, false
	}

	reserve := sdk.ZeroInt()

	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &reserve)
	return reserve, true
}

func (k Keeper) SetSubToken(ctx sdk.Context, denom, id string, subTokenID int64, reserve sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	subTokenKey := types.GetSubTokenKey(denom, id, subTokenID)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(reserve)
	store.Set(subTokenKey, bz)
}

func (k Keeper) RemoveSubToken(ctx sdk.Context, denom, id string, subTokenID int64) {
	store := ctx.KVStore(k.storeKey)
	subTokenKey := types.GetSubTokenKey(denom, id, subTokenID)
	store.Delete(subTokenKey)
}

func (k Keeper) GetLastSubTokenID(ctx sdk.Context, denom, id string) int64 {
	store := ctx.KVStore(k.storeKey)
	lastSubTokenIDKey := types.GetLastSubTokenIDKey(denom, id)
	bz := store.Get(lastSubTokenIDKey)
	if bz == nil {
		return 0
	}

	b := make([]byte, 8)
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &b)
	return int64(binary.LittleEndian.Uint64(b))
}

func (k Keeper) SetLastSubTokenID(ctx sdk.Context, denom, id string, lastSubTokenID int64) {
	store := ctx.KVStore(k.storeKey)
	lastSubTokenIDKey := types.GetLastSubTokenIDKey(denom, id)
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(lastSubTokenID))
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(b)
	store.Set(lastSubTokenIDKey, bz)
}

// MintNFT mints an NFT and manages that NFTs existence within Collections and Owners
func (k Keeper) MintNFT(ctx sdk.Context, denom, id string, reserve, quantity sdk.Int,
	creator, owner sdk.AccAddress, tokenURI string, allowMint bool) (int64, error) {

	nft, err := k.GetNFT(ctx, denom, id)
	if err == nil {
		reserve = nft.GetReserve()
	}

	err = k.ReserveTokens(ctx,
		sdk.NewCoins(
			sdk.NewCoin(
				k.baseDenom,
				reserve.Mul(quantity), // reserve * quantity
			)),
		creator)
	if err != nil {
		return 0, err
	}

	lastSubTokenID := k.GetLastSubTokenID(ctx, denom, id)

	if lastSubTokenID == 0 {
		lastSubTokenID = 1
	}

	tempSubTokenID := lastSubTokenID
	subTokenIDs := make([]int64, quantity.Int64())
	for i := int64(0); i < quantity.Int64(); i++ {
		subTokenIDs[i] = tempSubTokenID
		tempSubTokenID++
	}

	nft = types.NewBaseNFT(id, creator, owner, tokenURI, reserve, subTokenIDs, allowMint)
	collection, found := k.GetCollection(ctx, denom)
	if found {
		collection, err = collection.AddNFT(nft)
		if err != nil {
			return 0, err
		}
	} else {
		collection = types.NewCollection(denom, types.NewNFTs(nft))
	}
	k.SetCollection(ctx, denom, collection)

	newLastSubTokenID := lastSubTokenID + quantity.Int64()

	for i := lastSubTokenID; i < newLastSubTokenID; i++ {
		k.SetSubToken(ctx, denom, nft.GetID(), i, nft.GetReserve())
	}

	k.SetLastSubTokenID(ctx, denom, nft.GetID(), newLastSubTokenID)

	ownerIDCollection, _ := k.GetOwnerByDenom(ctx, nft.GetCreator(), denom)
	ownerIDCollection = ownerIDCollection.AddID(nft.GetID())
	k.SetOwnerByDenom(ctx, nft.GetCreator(), denom, ownerIDCollection.IDs)
	return newLastSubTokenID, err
}

// UpdateNFT updates an already existing NFTs
func (k Keeper) UpdateNFT(ctx sdk.Context, denom string, nft exported.NFT) (err error) {
	collection, found := k.GetCollection(ctx, denom)
	if !found {
		return sdkerrors.Wrap(types.ErrUnknownCollection, fmt.Sprintf("collection #%s doesn't exist", denom))
	}

	oldNFT, err := collection.GetNFT(nft.GetID())

	if err != nil {
		return err
	}

	collection.NFTs, _ = collection.NFTs.Update(oldNFT.GetID(), nft)

	k.SetCollection(ctx, denom, collection)
	return nil
}

// DeleteNFT deletes an existing NFT from store
func (k Keeper) DeleteNFT(ctx sdk.Context, denom, id string, subTokenIDs []int64) error {
	collection, found := k.GetCollection(ctx, denom)
	if !found {
		return sdkerrors.Wrap(types.ErrUnknownCollection, fmt.Sprintf("collection of %s doesn't exist", denom))
	}
	nft, err := collection.GetNFT(id)
	if err != nil {
		return err
	}

	owner := nft.GetOwners().GetOwner(nft.GetCreator())
	ownerSubTokenIDs := types.SortedIntArray(owner.GetSubTokenIDs())
	for _, subTokenID := range subTokenIDs {
		if ownerSubTokenIDs.Find(subTokenID) == -1 {
			return sdkerrors.Wrap(types.ErrNotAllowedBurn,
				fmt.Sprintf("owner %s has only %s tokens", nft.GetCreator(),
					types.SortedIntArray(nft.GetOwners().GetOwner(nft.GetCreator()).GetSubTokenIDs()).String()))
		}
		owner = owner.RemoveSubTokenID(subTokenID)
		k.RemoveSubToken(ctx, denom, id, subTokenID)
	}

	nft = nft.SetOwners(nft.
		GetOwners().
		SetOwner(owner))

	nftOwner, _ := k.GetOwner(ctx, nft.GetCreator()).DeleteID(denom, nft.GetID())
	k.SetOwner(ctx, nftOwner)

	collection, err = collection.DeleteNFT(nft)
	if err != nil {
		return err
	}

	k.SetCollection(ctx, denom, collection)

	err = k.BurnTokens(ctx, sdk.NewCoins(
		sdk.NewCoin(k.baseDenom, nft.GetReserve().MulRaw(int64(len(subTokenIDs))))))
	if err != nil {
		return err
	}

	return nil
}
