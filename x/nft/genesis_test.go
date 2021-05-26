package nft_test

/*
import (
	"testing"

	"bitbucket.org/decimalteam/go-node/x/nft"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestInitGenesis(t *testing.T) {
	ctx, _, NFTKeeper := createTestApp(t, false)

	genesisState := nft.DefaultGenesisState()
	require.Equal(t, 0, len(genesisState.Owners))
	require.Equal(t, 0, len(genesisState.Collections))

	ids := []string{id, id2, id3}
	idCollection := nft.NewIDCollection(denom, ids)
	idCollection2 := nft.NewIDCollection(denom2, ids)
	owner := nft.NewOwner(address, idCollection)

	owner2 := nft.NewOwner(address2, idCollection2)

	owners := []nft.Owner{owner, owner2}

	nft1 := nft.NewBaseNFT(id, address, address, tokenURI1, sdk.NewInt(1), sdk.NewInt(101), true)
	nft2 := nft.NewBaseNFT(id2, address, address, tokenURI1, sdk.NewInt(1), sdk.NewInt(101), true)
	nft3 := nft.NewBaseNFT(id3, address, address, tokenURI1, sdk.NewInt(1), sdk.NewInt(101), true)
	nfts := nft.NewNFTs(nft1, nft2, nft3)
	collection := nft.NewCollection(denom, nfts)

	nftx := nft.NewBaseNFT(id, address2, address2, tokenURI1, sdk.NewInt(1), sdk.NewInt(101), true)
	nft2x := nft.NewBaseNFT(id2, address2, address2, tokenURI1, sdk.NewInt(1), sdk.NewInt(101), true)
	nft3x := nft.NewBaseNFT(id3, address2, address2, tokenURI1, sdk.NewInt(1), sdk.NewInt(101), true)
	nftsx := nft.NewNFTs(nftx, nft2x, nft3x)
	collection2 := nft.NewCollection(denom2, nftsx)

	collections := nft.NewCollections(collection, collection2)

	genesisState = nft.NewGenesisState(owners, collections)

	nft.InitGenesis(ctx, NFTKeeper, genesisState)

	returnedOwners := NFTKeeper.GetOwners(ctx)
	require.Equal(t, 2, len(owners))
	require.Equal(t, returnedOwners[0].String(), owners[0].String())
	require.Equal(t, returnedOwners[1].String(), owners[1].String())

	returnedCollections := NFTKeeper.GetCollections(ctx)
	require.Equal(t, 2, len(returnedCollections))
	require.Equal(t, returnedCollections[0].String(), collections[0].String())
	require.Equal(t, returnedCollections[1].String(), collections[1].String())

	exportedGenesisState := nft.ExportGenesis(ctx, NFTKeeper)
	require.Equal(t, len(genesisState.Owners), len(exportedGenesisState.Owners))
	require.Equal(t, genesisState.Owners[0].String(), exportedGenesisState.Owners[0].String())
	require.Equal(t, genesisState.Owners[1].String(), exportedGenesisState.Owners[1].String())

	require.Equal(t, len(genesisState.Collections), len(exportedGenesisState.Collections))
	require.Equal(t, genesisState.Collections[0].String(), exportedGenesisState.Collections[0].String())
	require.Equal(t, genesisState.Collections[1].String(), exportedGenesisState.Collections[1].String())
}*/
