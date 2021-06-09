package keeper

import (
	"bitbucket.org/decimalteam/go-node/x/nft/internal/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMintNFT(t *testing.T) {
	ctx, _, NFTKeeper := createTestApp(t, false)

	// MintNFT shouldn't fail when collection does not exist
	nft := types.NewBaseNFT(
		ID1,
		Addrs[0],
		Addrs[0],
		TokenURI1,
		sdk.NewInt(100),
		[]int64{},
		true,
	)
	_, err := NFTKeeper.MintNFT(ctx, Denom1, nft.GetID(), nft.GetReserve(), sdk.NewInt(1), nft.GetCreator(), Addrs[0], nft.GetTokenURI(), nft.GetAllowMint())

	require.NoError(t, err)

	// MintNFT shouldn't fail when collection exists
	nft2 := types.NewBaseNFT(
		ID2,
		Addrs[0],
		Addrs[0],
		TokenURI1,
		sdk.NewInt(100),
		[]int64{},
		true,
	)
	_, err = NFTKeeper.MintNFT(ctx, Denom1, nft2.GetID(), nft2.GetReserve(), sdk.NewInt(1), nft2.GetCreator(), Addrs[0], nft2.GetTokenURI(), nft2.GetAllowMint())
	require.NoError(t, err)

	// MintNFT should correctly calculate lastSubTokenId
	nft3 := types.NewBaseNFT(
		ID3,
		Addrs[0],
		Addrs[0],
		TokenURI1,
		sdk.NewInt(100),
		[]int64{3, 2, 4},
		true,
	)
	quantity := sdk.NewInt(50)
	lastSubTokenID, err := NFTKeeper.MintNFT(
		ctx, Denom1, nft3.GetID(),
		nft3.GetReserve(), quantity, nft3.GetCreator(),
		Addrs[0], nft3.GetTokenURI(), nft3.GetAllowMint(),
	)
	require.NoError(t, err)
	require.Equal(t, quantity.AddRaw(1).Int64(), lastSubTokenID)
}

func TestGetNFT(t *testing.T) {
	ctx, _, NFTKeeper := createTestApp(t, false)

	// MintNFT shouldn't fail when collection does not exist
	nft := types.NewBaseNFT(
		ID1,
		Addrs[0],
		Addrs[0],
		TokenURI1,
		sdk.NewInt(100),
		[]int64{},
		true,
	)
	_, err := NFTKeeper.MintNFT(ctx, Denom1, nft.GetID(), nft.GetReserve(), sdk.NewInt(1), nft.GetCreator(), Addrs[0], nft.GetTokenURI(), nft.GetAllowMint())

	require.NoError(t, err)

	// GetNFT should get the NFT
	receivedNFT, err := NFTKeeper.GetNFT(ctx, Denom1, ID1)
	require.NoError(t, err)
	require.Equal(t, receivedNFT.GetID(), ID1)
	require.True(t, receivedNFT.GetCreator().Equals(Addrs[0]))
	require.Equal(t, receivedNFT.GetTokenURI(), TokenURI1)

	// MintNFT shouldn't fail when collection exists
	nft2 := types.NewBaseNFT(
		ID2,
		Addrs[0],
		Addrs[0],
		TokenURI1,
		sdk.NewInt(100),
		[]int64{},
		true,
	)
	_, err = NFTKeeper.MintNFT(ctx, Denom1, nft2.GetID(), nft2.GetReserve(), sdk.NewInt(1), nft2.GetCreator(), Addrs[0], nft2.GetTokenURI(), nft2.GetAllowMint())
	require.NoError(t, err)

	// GetNFT should get the NFT when collection exists
	receivedNFT2, err := NFTKeeper.GetNFT(ctx, Denom1, ID2)
	require.NoError(t, err)
	require.Equal(t, receivedNFT2.GetID(), ID2)
	require.True(t, receivedNFT2.GetCreator().Equals(Addrs[0]))
	require.Equal(t, receivedNFT2.GetTokenURI(), TokenURI1)

	msg, fail := SupplyInvariant(NFTKeeper)(ctx)
	require.False(t, fail, msg)
}

func TestUpdateNFT(t *testing.T) {
	ctx, _, NFTKeeper := createTestApp(t, false)

	subTokenIDs := []int64{}
	reserve := sdk.NewInt(100)
	quantity := sdk.NewInt(1)

	nft := types.NewBaseNFT(
		ID1,
		Addrs[0],
		Addrs[0],
		TokenURI1,
		reserve,
		subTokenIDs,
		true,
	)

	// UpdateNFT should fail when nft doesn't exist
	//_, err := NFTKeeper.MintNFT(ctx, Denom1, ID2, reserve, quantity, Addrs[0], Addrs[0], TokenURI1, true)
	//require.Error(t, err)

	// MintNFT shouldn't fail when collection does not exist
	_, err := NFTKeeper.MintNFT(ctx, Denom1, ID1, nft.GetReserve(), quantity, nft.GetCreator(), Addrs[0], nft.GetTokenURI(), nft.GetAllowMint())
	require.NoError(t, err)

	// UpdateNFT shouldn't fail when NFT exists
	nft2 := types.NewBaseNFT(
		ID1,
		Addrs[1],
		Addrs[1],
		TokenURI2,
		reserve,
		subTokenIDs,
		true,
	)
	err = NFTKeeper.UpdateNFT(ctx, Denom1, nft2)
	require.NoError(t, err)

	// GetNFT should get the NFT with new TokenURI1
	receivedNFT, err := NFTKeeper.GetNFT(ctx, Denom1, ID1)
	require.NoError(t, err)
	require.Equal(t, receivedNFT.GetTokenURI(), TokenURI2)
}

func TestDeleteNFT(t *testing.T) {
	ctx, _, NFTKeeper := createTestApp(t, false)

	// DeleteNFT should fail when NFT doesn't exist and collection doesn't exist
	subTokenIDs := []int64{}
	err := NFTKeeper.DeleteNFT(ctx, Denom1, ID1, subTokenIDs)
	require.Error(t, err)

	// MintNFT should not fail when collection does not exist
	nft := types.NewBaseNFT(
		ID1,
		Addrs[0],
		Addrs[0],
		TokenURI1,
		sdk.NewInt(100),
		subTokenIDs,
		true,
	)
	_, err = NFTKeeper.MintNFT(ctx, Denom1, nft.GetID(), nft.GetReserve(), sdk.NewInt(1), nft.GetCreator(), Addrs[0], nft.GetTokenURI(), nft.GetAllowMint())

	require.NoError(t, err)

	// DeleteNFT should fail when NFT doesn't exist but collection does exist
	err = NFTKeeper.DeleteNFT(ctx, Denom1, ID2, subTokenIDs)
	require.Error(t, err)

	// DeleteNFT should fail when at least of nft's subtokenIds is not in the owner's subTokenIDs
	err = NFTKeeper.DeleteNFT(ctx, Denom1, ID1, []int64{3})
	require.Error(t, err)

	// DeleteNFT should not fail when NFT and collection exist
	err = NFTKeeper.DeleteNFT(ctx, Denom1, ID1, subTokenIDs)
	require.NoError(t, err)

	// NFT should no longer exist
	isNFT := NFTKeeper.IsNFT(ctx, Denom1, ID1)
	require.True(t, isNFT)

	owner := NFTKeeper.GetOwner(ctx, Addrs[0])
	require.Equal(t, 0, owner.Supply())
}

func TestIsNFT(t *testing.T) {
	ctx, _, NFTKeeper := createTestApp(t, false)

	// IsNFT should return false
	isNFT := NFTKeeper.IsNFT(ctx, Denom1, ID1)
	require.False(t, isNFT)

	// MintNFT shouldn't fail when collection does not exist
	nft := types.NewBaseNFT(
		ID1,
		Addrs[0],
		Addrs[0],
		TokenURI1,
		sdk.NewInt(100),
		[]int64{},
		true,
	)
	_, err := NFTKeeper.MintNFT(ctx, Denom1, nft.GetID(), nft.GetReserve(), sdk.NewInt(1), nft.GetCreator(), Addrs[0], nft.GetTokenURI(), nft.GetAllowMint())

	require.NoError(t, err)

	// IsNFT should return true
	isNFT = NFTKeeper.IsNFT(ctx, Denom1, ID1)
	require.True(t, isNFT)
}
