package nft

import (
	"bitbucket.org/decimalteam/go-node/x/nft/internal/types"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestInitGenesis(t *testing.T) {
	ctx, _, NFTKeeper := createTestApp(t, false)

	genesisState := DefaultGenesisState()
	require.Equal(t, 0, len(genesisState.Owners))
	require.Equal(t, 0, len(genesisState.Collections))

	ids := []string{ID1, ID2, ID3}
	idCollection := types.NewIDCollection(Denom1, ids)
	idCollection2 := types.NewIDCollection(Denom2, ids)
	owner := types.NewOwner(Addrs[0], idCollection)

	owner2 := types.NewOwner(Addrs[1], idCollection2)

	owners := []Owner{owner, owner2}

	reserve := sdk.NewInt(100)
	subTokenIds := []int64{}

	nft1 := types.NewBaseNFT(ID1, Addrs[0], Addrs[0], TokenURI1, reserve, subTokenIds, true)
	nft2 := types.NewBaseNFT(ID2, Addrs[0], Addrs[0], TokenURI1, reserve, subTokenIds, true)
	nft3 := types.NewBaseNFT(ID3, Addrs[0], Addrs[0], TokenURI1, reserve, subTokenIds, true)
	nfts := types.NewNFTs(nft1, nft2, nft3)
	collection := NewCollection(Denom1, nfts)

	nftx := types.NewBaseNFT(ID1, Addrs[1], Addrs[1], TokenURI1, reserve, subTokenIds, true)
	nft2x := types.NewBaseNFT(ID2, Addrs[1], Addrs[1], TokenURI1, reserve, subTokenIds, true)
	nft3x := types.NewBaseNFT(ID3, Addrs[1], Addrs[1], TokenURI1, reserve, subTokenIds, true)
	nftsx := types.NewNFTs(nftx, nft2x, nft3x)
	collection2 := NewCollection(Denom2, nftsx)

	collections := NewCollections(collection, collection2)

	genesisState = NewGenesisState(owners, collections)

	InitGenesis(ctx, NFTKeeper, genesisState)

	returnedOwners := NFTKeeper.GetOwners(ctx)
	require.Equal(t, 2, len(owners))
	require.Equal(t, returnedOwners[0].String(), owners[0].String())
	require.Equal(t, returnedOwners[1].String(), owners[1].String())

	returnedCollections := NFTKeeper.GetCollections(ctx)
	require.Equal(t, 2, len(returnedCollections))
	require.Equal(t, returnedCollections[0].String(), collections[0].String())
	require.Equal(t, returnedCollections[1].String(), collections[1].String())

	exportedGenesisState := ExportGenesis(ctx, NFTKeeper)
	require.Equal(t, len(genesisState.Owners), len(exportedGenesisState.Owners))
	require.Equal(t, genesisState.Owners[0].String(), exportedGenesisState.Owners[0].String())
	require.Equal(t, genesisState.Owners[1].String(), exportedGenesisState.Owners[1].String())

	require.Equal(t, len(genesisState.Collections), len(exportedGenesisState.Collections))
	require.Equal(t, genesisState.Collections[0].String(), exportedGenesisState.Collections[0].String())
	require.Equal(t, genesisState.Collections[1].String(), exportedGenesisState.Collections[1].String())
}
