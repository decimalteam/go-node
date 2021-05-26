package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/stretchr/testify/require"

	"bitbucket.org/decimalteam/go-node/x/nft/internal/keeper"
	"bitbucket.org/decimalteam/go-node/x/nft/internal/types"
)

func TestGetOwners(t *testing.T) {
	ctx, _, NFTKeeper := createTestApp(t, false)

	nft := types.NewBaseNFT(id, address, address, tokenURI, sdk.NewInt(1), sdk.NewInt(101), true)
	_, err := NFTKeeper.MintNFT(ctx, denom, nft)
	require.NoError(t, err)

	nft2 := types.NewBaseNFT(id2, address, address, tokenURI, sdk.NewInt(1), sdk.NewInt(101), true)
	_, err = NFTKeeper.MintNFT(ctx, denom, nft)
	require.NoError(t, err)

	nft3 := types.NewBaseNFT(id3, address3, address3, tokenURI, sdk.NewInt(1), sdk.NewInt(101), true)
	_, err = NFTKeeper.MintNFT(ctx, denom, nft3)
	require.NoError(t, err)

	owners := NFTKeeper.GetOwners(ctx)
	require.Equal(t, 3, len(owners))

	nft = types.NewBaseNFT(id, address, address, tokenURI, sdk.NewInt(1), sdk.NewInt(101), true)
	_, err = NFTKeeper.MintNFT(ctx, denom2, nft)
	require.NoError(t, err)

	nft2 = types.NewBaseNFT(id2, address2, address, tokenURI, sdk.NewInt(1), sdk.NewInt(101), true)
	_, err = NFTKeeper.MintNFT(ctx, denom2, nft2)
	require.NoError(t, err)

	nft3 = types.NewBaseNFT(id3, address3, address3, tokenURI, sdk.NewInt(1), sdk.NewInt(101), true)
	_, err = NFTKeeper.MintNFT(ctx, denom2, nft3)
	require.NoError(t, err)

	owners = NFTKeeper.GetOwners(ctx)
	require.Equal(t, 3, len(owners))

	msg, fail := keeper.SupplyInvariant(NFTKeeper)(ctx)
	require.False(t, fail, msg)
}

func TestSetOwner(t *testing.T) {
	ctx, _, NFTKeeper := createTestApp(t, false)

	nft := types.NewBaseNFT(id, address, address, tokenURI, sdk.NewInt(1), sdk.NewInt(101), true)
	_, err := NFTKeeper.MintNFT(ctx, denom, nft)
	require.NoError(t, err)

	idCollection := types.NewIDCollection(denom, []string{id, id2, id3})
	owner := types.NewOwner(address, idCollection)

	oldOwner := NFTKeeper.GetOwner(ctx, address)

	NFTKeeper.SetOwner(ctx, owner)

	newOwner := NFTKeeper.GetOwner(ctx, address)
	require.NotEqual(t, oldOwner.String(), newOwner.String())
	require.Equal(t, owner.String(), newOwner.String())

	// for invariant sanity
	NFTKeeper.SetOwner(ctx, oldOwner)

	msg, fail := keeper.SupplyInvariant(NFTKeeper)(ctx)
	require.False(t, fail, msg)
}

func TestSetOwners(t *testing.T) {
	ctx, _, NFTKeeper := createTestApp(t, false)

	// create NFT where id = "id" with owner = "address"
	nft := types.NewBaseNFT(id, address, address, tokenURI, sdk.NewInt(1), sdk.NewInt(101), true)
	_, err := NFTKeeper.MintNFT(ctx, denom, nft)
	require.NoError(t, err)

	// create NFT where id = "id2" with owner = "address2"
	nft = types.NewBaseNFT(id2, address2, address2, tokenURI, sdk.NewInt(1), sdk.NewInt(101), true)
	_, err = NFTKeeper.MintNFT(ctx, denom, nft)
	require.NoError(t, err)

	// create two owners (address and address2) with the same id collections of "id", "id2"  "id3"
	idCollection := types.NewIDCollection(denom, []string{id, id2, id3})
	owner := types.NewOwner(address, idCollection)
	owner2 := types.NewOwner(address2, idCollection)

	// get both owners that were created during the NFT mint process
	oldOwner := NFTKeeper.GetOwner(ctx, address)
	oldOwner2 := NFTKeeper.GetOwner(ctx, address2)

	// replace previous old owners with updated versions (that have multiple ids)
	NFTKeeper.SetOwners(ctx, []types.Owner{owner, owner2})

	newOwner := NFTKeeper.GetOwner(ctx, address)
	require.NotEqual(t, oldOwner.String(), newOwner.String())
	require.Equal(t, owner.String(), newOwner.String())

	newOwner2 := NFTKeeper.GetOwner(ctx, address2)
	require.NotEqual(t, oldOwner2.String(), newOwner2.String())
	require.Equal(t, owner2.String(), newOwner2.String())

	// replace old owners for invariance sanity
	NFTKeeper.SetOwners(ctx, []types.Owner{oldOwner, oldOwner2})

	msg, fail := keeper.SupplyInvariant(NFTKeeper)(ctx)
	require.False(t, fail, msg)
}

func TestSwapOwners(t *testing.T) {
	ctx, _, NFTKeeper := createTestApp(t, false)

	nft := types.NewBaseNFT(id, address, address, tokenURI, sdk.NewInt(1), sdk.NewInt(101), true)
	_, err := NFTKeeper.MintNFT(ctx, denom, nft)
	require.NoError(t, err)

	err = NFTKeeper.SwapOwners(ctx, denom, id, address, address2)
	require.NoError(t, err)

	err = NFTKeeper.SwapOwners(ctx, denom, id, address, address2)
	require.Error(t, err)

	err = NFTKeeper.SwapOwners(ctx, denom2, id, address, address2)
	require.Error(t, err)

	msg, fail := keeper.SupplyInvariant(NFTKeeper)(ctx)
	require.False(t, fail, msg)
}
