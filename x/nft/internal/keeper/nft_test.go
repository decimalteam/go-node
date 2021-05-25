package keeper_test

import (
	"testing"

	"bitbucket.org/decimalteam/go-node/x/nft/internal/keeper"
	"bitbucket.org/decimalteam/go-node/x/nft/internal/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/stretchr/testify/require"
)

func TestMintNFT(t *testing.T) {
	ctx, _, NFTKeeper := createTestApp(t, false)

	// MintNFT shouldn't fail when collection does not exist
	nft := types.NewBaseNFT(id, address, address, tokenURI, sdk.NewInt(1), sdk.NewInt(101), true)
	err := NFTKeeper.MintNFT(ctx, denom, nft)
	require.NoError(t, err)

	// MintNFT shouldn't fail when collection exists
	nft2 := types.NewBaseNFT(id2, address, address, tokenURI, sdk.NewInt(1), sdk.NewInt(101), true)
	err = NFTKeeper.MintNFT(ctx, denom, nft2)
	require.NoError(t, err)
}

func TestGetNFT(t *testing.T) {
	ctx, _, NFTKeeper := createTestApp(t, false)

	// MintNFT shouldn't fail when collection does not exist
	nft := types.NewBaseNFT(id, address, address, tokenURI, sdk.NewInt(1), sdk.NewInt(101), true)
	err := NFTKeeper.MintNFT(ctx, denom, nft)
	require.NoError(t, err)

	// GetNFT should get the NFT
	receivedNFT, err := NFTKeeper.GetNFT(ctx, denom, id)
	require.NoError(t, err)
	require.Equal(t, receivedNFT.GetID(), id)
	require.True(t, receivedNFT.GetCreator().Equals(address))
	require.Equal(t, receivedNFT.GetTokenURI(), tokenURI)

	// MintNFT shouldn't fail when collection exists
	nft2 := types.NewBaseNFT(id2, address, address, tokenURI, sdk.NewInt(1), sdk.NewInt(101), true)
	err = NFTKeeper.MintNFT(ctx, denom, nft2)
	require.NoError(t, err)

	// GetNFT should get the NFT when collection exists
	receivedNFT2, err := NFTKeeper.GetNFT(ctx, denom, id2)
	require.NoError(t, err)
	require.Equal(t, receivedNFT2.GetID(), id2)
	require.True(t, receivedNFT2.GetCreator().Equals(address))
	require.Equal(t, receivedNFT2.GetTokenURI(), tokenURI)

	msg, fail := keeper.SupplyInvariant(NFTKeeper)(ctx)
	require.False(t, fail, msg)
}

//func TestUpdateNFT(t *testing.T) {
//	ctx, _, NFTKeeper := createTestApp(t, false)
//
//	nft := types.NewBaseNFT(id2, address, address, tokenURI, sdk.NewInt(1), sdk.NewInt(101), true)
//
//	// UpdateNFT should fail when NFT doesn't exists
//	err := NFTKeeper.MintNFT(ctx, denom, &nft)
//	require.Error(t, err)
//
//	// MintNFT shouldn't fail when collection does not exist
//	err = NFTKeeper.MintNFT(ctx, denom, &nft)
//	require.NoError(t, err)
//
//	nonnft := types.NewBaseNFT(id2, address, tokenURI)
//	// UpdateNFT should fail when NFT doesn't exists
//	err = NFTKeeper.UpdateNFT(ctx, denom, &nonnft)
//	require.Error(t, err)
//
//	// UpdateNFT shouldn't fail when NFT exists
//	nft2 := types.NewBaseNFT(id, address, tokenURI2)
//	err = NFTKeeper.UpdateNFT(ctx, denom, &nft2)
//	require.NoError(t, err)
//
//	// UpdateNFT shouldn't fail when NFT exists
//	nft2 = types.NewBaseNFT(id, address2, tokenURI2)
//	err = NFTKeeper.UpdateNFT(ctx, denom, &nft2)
//	require.NoError(t, err)
//
//	// GetNFT should get the NFT with new tokenURI
//	receivedNFT, err := NFTKeeper.GetNFT(ctx, denom, id)
//	require.NoError(t, err)
//	require.Equal(t, receivedNFT.GetTokenURI(), tokenURI2)
//}

func TestDeleteNFT(t *testing.T) {
	ctx, _, NFTKeeper := createTestApp(t, false)

	// DeleteNFT should fail when NFT doesn't exist and collection doesn't exist
	err := NFTKeeper.DeleteNFT(ctx, denom, id, sdk.NewInt(1))
	require.Error(t, err)

	// MintNFT should not fail when collection does not exist
	nft := types.NewBaseNFT(id, address, address, tokenURI, sdk.NewInt(1), sdk.NewInt(101), true)
	err = NFTKeeper.MintNFT(ctx, denom, nft)
	require.NoError(t, err)

	// DeleteNFT should fail when NFT doesn't exist but collection does exist
	err = NFTKeeper.DeleteNFT(ctx, denom, id2, sdk.NewInt(1))
	require.Error(t, err)

	// DeleteNFT should fail when delete quantity is more than the reserved one
	err = NFTKeeper.DeleteNFT(ctx, denom, id, sdk.NewInt(4))
	require.Error(t, err)

	// DeleteNFT should not fail when NFT and collection exist
	err = NFTKeeper.DeleteNFT(ctx, denom, id, sdk.NewInt(1))
	require.NoError(t, err)

	// NFT should no longer exist
	isNFT := NFTKeeper.IsNFT(ctx, denom, id)
	require.False(t, isNFT)

	owner := NFTKeeper.GetOwner(ctx, address)
	require.Equal(t, 0, owner.Supply())

	msg, fail := keeper.SupplyInvariant(NFTKeeper)(ctx)
	require.False(t, fail, msg)
}

func TestIsNFT(t *testing.T) {
	ctx, _, NFTKeeper := createTestApp(t, false)

	// IsNFT should return false
	isNFT := NFTKeeper.IsNFT(ctx, denom, id)
	require.False(t, isNFT)

	// MintNFT shouldn't fail when collection does not exist
	nft := types.NewBaseNFT(id, address, address, tokenURI, sdk.NewInt(1), sdk.NewInt(101), true)
	err := NFTKeeper.MintNFT(ctx, denom, nft)
	require.NoError(t, err)

	// IsNFT should return true
	isNFT = NFTKeeper.IsNFT(ctx, denom, id)
	require.True(t, isNFT)
}
