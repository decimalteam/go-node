package keeper

import (
	"bitbucket.org/decimalteam/go-node/x/nft/internal/types"
	"fmt"
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

//func TestUpdateNFT(t *testing.T) {
//	ctx, _, NFTKeeper := createTestApp(t, false)
//
//	nft := types.NewBaseNFT(ID2, Addrs[0, Addrs[0, TokenURI1, sdk.NewInt(1), sdk.NewInt(101), true)
//
//	// UpdateNFT should fail when NFT doesn't exists
//	err := NFTKeeper.MintNFT(ctx, Denom1, &nft)
//	require.Error(t, err)
//
//	// MintNFT shouldn't fail when collection does not exist
//	err = NFTKeeper.MintNFT(ctx, Denom1, &nft)
//	require.NoError(t, err)
//
//	nonnft := types.NewBaseNFT(ID2, Addrs[0, TokenURI1)
//	// UpdateNFT should fail when NFT doesn't exists
//	err = NFTKeeper.UpdateNFT(ctx, Denom1, &nonnft)
//	require.Error(t, err)
//
//	// UpdateNFT shouldn't fail when NFT exists
//	nft2 := types.NewBaseNFT(ID1, Addrs[0, TokenURI2)
//	err = NFTKeeper.UpdateNFT(ctx, Denom1, &nft2)
//	require.NoError(t, err)
//
//	// UpdateNFT shouldn't fail when NFT exists
//	nft2 = types.NewBaseNFT(ID1, address2, TokenURI2)
//	err = NFTKeeper.UpdateNFT(ctx, Denom1, &nft2)
//	require.NoError(t, err)
//
//	// GetNFT should get the NFT with new TokenURI1
//	receivedNFT, err := NFTKeeper.GetNFT(ctx, Denom1, ID1)
//	require.NoError(t, err)
//	require.Equal(t, receivedNFT.GetTokenURI(), TokenURI2)
//}

func TestDeleteNFT(t *testing.T) {
	ctx, _, NFTKeeper := createTestApp(t, false)

	// DeleteNFT should fail when NFT doesn't exist and collection doesn't exist
	err := NFTKeeper.DeleteNFT(ctx, Denom1, ID1, []int64{})
	require.Error(t, err)

	// MintNFT should not fail when collection does not exist
	nft := types.NewBaseNFT(
		ID1,
		Addrs[0],
		Addrs[0],
		TokenURI1,
		sdk.NewInt(100),
		[]int64{},
		true,
	)
	_, err = NFTKeeper.MintNFT(ctx, Denom1, nft.GetID(), nft.GetReserve(), sdk.NewInt(1), nft.GetCreator(), Addrs[0], nft.GetTokenURI(), nft.GetAllowMint())

	require.NoError(t, err)

	// DeleteNFT should fail when NFT doesn't exist but collection does exist
	err = NFTKeeper.DeleteNFT(ctx, Denom1, ID2, []int64{})
	require.Error(t, err)

	// DeleteNFT should fail when delete quantity is more than the reserved one
	//err = NFTKeeper.DeleteNFT(ctx, Denom1, ID1, []int64{})
	//require.Error(t, err)

	// DeleteNFT should not fail when NFT and collection exist
	err = NFTKeeper.DeleteNFT(ctx, Denom1, ID1, []int64{})
	require.NoError(t, err)

	// NFT should no longer exist
	isNFT := NFTKeeper.IsNFT(ctx, Denom1, ID1)
	require.False(t, isNFT)

	owner := NFTKeeper.GetOwner(ctx, Addrs[0])
	s := owner.String()
	fmt.Println(s)
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
