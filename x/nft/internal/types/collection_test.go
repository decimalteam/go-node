package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"testing"

	"github.com/stretchr/testify/require"
)

// ---------------------------------------- Collection ---------------------------------------------------

func TestNewCollection(t *testing.T) {
	reserve := sdk.NewInt(100)
	var subTokenIDs []int64
	testNFT := NewBaseNFT(ID1, Addrs[0], Addrs[0], TokenURI1, reserve, subTokenIDs, true)
	testNFT2 := NewBaseNFT(ID2, Addrs[0], Addrs[0], TokenURI1, reserve, subTokenIDs, true)
	nfts := NewNFTs(testNFT, testNFT2)

	collection := NewCollection(fmt.Sprintf("      %s      ", Denom1), nfts)
	require.Equal(t, collection.Denom, Denom1)
	require.Equal(t, len(collection.NFTs), 2)
}

func TestEmptyCollection(t *testing.T) {
	collection := EmptyCollection()
	require.Equal(t, collection.Denom, "")
	require.Equal(t, len(collection.NFTs), 0)
}

func TestCollectionGetNFTMethod(t *testing.T) {
	reserve := sdk.NewInt(100)
	testNFT := NewBaseNFT(ID1, Addrs[0], Addrs[0], TokenURI1, reserve, []int64{}, true)
	nfts := NewNFTs(testNFT)
	collection := NewCollection(Denom1, nfts)

	returnedNFT, err := collection.GetNFT(ID1)
	require.NoError(t, err)
	require.Equal(t, testNFT.String(), returnedNFT.String())

	returnedNFT, err = collection.GetNFT(ID2)
	require.Error(t, err)
	require.Nil(t, returnedNFT)
}

func TestCollectionContainsNFTMethod(t *testing.T) {
	reserve := sdk.NewInt(100)
	var subTokenIDs []int64
	testNFT := NewBaseNFT(ID1, Addrs[0], Addrs[0], TokenURI1, reserve, subTokenIDs, true)
	nfts := NewNFTs(testNFT)
	collection := NewCollection(Denom1, nfts)

	contains := collection.ContainsNFT(ID1)
	require.True(t, contains)

	contains = collection.ContainsNFT(ID2)
	require.False(t, contains)
}

func TestCollectionAddNFTMethod(t *testing.T) {
	reserve := sdk.NewInt(100)
	var subTokenIDs []int64
	testNFT := NewBaseNFT(ID1, Addrs[0], Addrs[0], TokenURI1, reserve, subTokenIDs, true)
	testNFT2 := NewBaseNFT(ID2, Addrs[0], Addrs[0], TokenURI1, reserve, subTokenIDs, true)
	nfts := NewNFTs(testNFT)
	collection := NewCollection(Denom1, nfts)

	newCollection, err := collection.AddNFT(testNFT)
	require.NoError(t, err)
	require.Equal(t, collection.String(), newCollection.String())

	newCollection, err = collection.AddNFT(testNFT2)
	require.NoError(t, err)
	require.NotEqual(t, collection.String(), newCollection.String())
	require.Equal(t, len(newCollection.NFTs), 2)
}

func TestCollectionUpdateNFTMethod(t *testing.T) {
	reserve := sdk.NewInt(100)
	var subTokenIDs []int64
	testNFT := NewBaseNFT(ID1, Addrs[0], Addrs[0], TokenURI1, reserve, subTokenIDs, true)
	testNFT2 := NewBaseNFT(ID2, Addrs[1], Addrs[1], TokenURI2, reserve, subTokenIDs, true)
	testNFT3 := NewBaseNFT(ID1, Addrs[1], Addrs[1], TokenURI2, reserve, subTokenIDs, true)
	nfts := NewNFTs(testNFT)
	collection := NewCollection(Denom1, nfts)

	// should fail when nft does not exist in the collection
	newCollection, err := collection.UpdateNFT(testNFT2)
	require.Error(t, err)
	require.Equal(t, collection.String(), newCollection.String())

	collection, err = collection.UpdateNFT(testNFT3)
	require.NoError(t, err)

	returnedNFT, err := collection.GetNFT(ID1)
	require.NoError(t, err)

	require.Equal(t, returnedNFT.GetOwners().GetOwners()[0].GetAddress(), Addrs[1])
	require.Equal(t, returnedNFT.GetTokenURI(), TokenURI2)
}

func TestCollectionDeleteNFTMethod(t *testing.T) {
	reserve := sdk.NewInt(100)
	var subTokenIDs []int64
	testNFT := NewBaseNFT(ID1, Addrs[0], Addrs[0], TokenURI1, reserve, subTokenIDs, true)
	testNFT2 := NewBaseNFT(ID2, Addrs[1], Addrs[1], TokenURI2, reserve, subTokenIDs, true)
	testNFT3 := NewBaseNFT(ID3, Addrs[0], Addrs[0], TokenURI1, reserve, subTokenIDs, true)
	nfts := NewNFTs(testNFT, testNFT2)
	collection := NewCollection(Denom1, nfts)

	newCollection, err := collection.DeleteNFT(testNFT3)
	require.Error(t, err)
	require.Equal(t, collection.String(), newCollection.String())

	collection, err = collection.DeleteNFT(testNFT2)
	require.NoError(t, err)
	require.Equal(t, len(collection.NFTs), 1)

	returnedNFT, err := collection.GetNFT(ID2)
	require.Nil(t, returnedNFT)
	require.Error(t, err)
}

func TestCollectionSupplyMethod(t *testing.T) {
	empty := EmptyCollection()
	require.Equal(t, empty.Supply(), 0)

	reserve := sdk.NewInt(100)
	var subTokenIDs []int64
	testNFT := NewBaseNFT(ID1, Addrs[0], Addrs[0], TokenURI1, reserve, subTokenIDs, true)
	testNFT2 := NewBaseNFT(ID2, Addrs[1], Addrs[1], TokenURI2, reserve, subTokenIDs, true)
	nfts := NewNFTs(testNFT, testNFT2)
	collection := NewCollection(Denom1, nfts)

	require.Equal(t, collection.Supply(), 2)

	collection, err := collection.DeleteNFT(testNFT)
	require.Nil(t, err)
	require.Equal(t, collection.Supply(), 1)

	collection, err = collection.DeleteNFT(testNFT2)
	require.Nil(t, err)
	require.Equal(t, collection.Supply(), 0)

	collection, err = collection.AddNFT(testNFT)
	require.Nil(t, err)
	require.Equal(t, collection.Supply(), 1)
}

func TestCollectionStringMethod(t *testing.T) {
	reserve := sdk.NewInt(100)
	var subTokenIDs []int64
	testNFT := NewBaseNFT(ID1, Addrs[0], Addrs[0], TokenURI1, reserve, subTokenIDs, true)
	testNFT2 := NewBaseNFT(ID2, Addrs[1], Addrs[1], TokenURI2, reserve, subTokenIDs, true)
	nfts := NewNFTs(testNFT, testNFT2)
	collection := NewCollection(Denom1, nfts)
	require.Equal(t, collection.String(),
		fmt.Sprintf(`Denom: 				%s
NFTs:

ID:				%s
Owners:			%s 
TokenURI:		%s
ID:				%s
Owners:			%s 
TokenURI:		%s`, Denom1, ID1, Addrs[0].String(), TokenURI1,
			ID2, Addrs[1].String(), TokenURI2))
}

// ---------------------------------------- Collections ---------------------------------------------------

func TestNewCollections(t *testing.T) {
	emptyCollections := NewCollections()
	require.Empty(t, emptyCollections)

	reserve := sdk.NewInt(100)
	var subTokenIDs []int64

	testNFT := NewBaseNFT(ID1, Addrs[0], Addrs[0], TokenURI1, reserve, subTokenIDs, true)
	nfts := NewNFTs(testNFT)
	collection := NewCollection(Denom1, nfts)

	testNFT2 := NewBaseNFT(ID2, Addrs[1], Addrs[1], TokenURI2, reserve, subTokenIDs, true)
	nfts2 := NewNFTs(testNFT2)
	collection2 := NewCollection(Denom2, nfts2)

	collections := NewCollections(collection, collection2)
	require.Equal(t, len(collections), 2)
}
func TestCollectionsAppendMethod(t *testing.T) {
	reserve := sdk.NewInt(100)
	var subTokenIDs []int64
	testNFT := NewBaseNFT(ID1, Addrs[0], Addrs[0], TokenURI1, reserve, subTokenIDs, true)
	nfts := NewNFTs(testNFT)
	collection := NewCollection(Denom1, nfts)

	collections := NewCollections(collection)

	testNFT2 := NewBaseNFT(ID2, Addrs[1], Addrs[1], TokenURI2, reserve, subTokenIDs, true)
	nfts2 := NewNFTs(testNFT2)
	collection2 := NewCollection(Denom2, nfts2)
	collections2 := NewCollections(collection2)

	collections = collections.Append(collections2...)
	require.Equal(t, len(collections), 2)
}
func TestCollectionsFindMethod(t *testing.T) {
	reserve := sdk.NewInt(100)
	var subTokenIDs []int64
	testNFT := NewBaseNFT(ID1, Addrs[0], Addrs[0], TokenURI1, reserve, subTokenIDs, true)
	nfts := NewNFTs(testNFT)
	collection := NewCollection(Denom1, nfts)

	testNFT2 := NewBaseNFT(ID2, Addrs[1], Addrs[1], TokenURI2, reserve, subTokenIDs, true)
	nfts2 := NewNFTs(testNFT2)
	collection2 := NewCollection(Denom2, nfts2)

	collections := NewCollections(collection)

	foundCollection, found := collections.Find(Denom2)
	require.False(t, found)
	require.Empty(t, foundCollection)

	collections = NewCollections(collection, collection2)

	foundCollection, found = collections.Find(Denom2)
	require.True(t, found)
	require.Equal(t, foundCollection.String(), collection2.String())

	collection3 := NewCollection(Denom3, nfts)
	collections = NewCollections(collection, collection2, collection3)

	_, found = collections.Find(Denom1)
	require.True(t, found)

	_, found = collections.Find(Denom2)
	require.True(t, found)

	_, found = collections.Find(Denom3)
	require.True(t, found)
}

func TestCollectionsRemoveMethod(t *testing.T) {
	reserve := sdk.NewInt(100)
	var subTokenIDs []int64
	testNFT := NewBaseNFT(ID1, Addrs[0], Addrs[0], TokenURI1, reserve, subTokenIDs, true)
	nfts := NewNFTs(testNFT)
	collection := NewCollection(Denom1, nfts)

	collections := NewCollections(collection)

	returnedCollections, removed := collections.Remove(Denom2)
	require.False(t, removed)
	require.Equal(t, returnedCollections.String(), collections.String())

	testNFT2 := NewBaseNFT(ID2, Addrs[1], Addrs[1], TokenURI2, reserve, subTokenIDs, true)
	nfts2 := NewNFTs(testNFT2)
	collection2 := NewCollection(Denom2, nfts2)

	collections = NewCollections(collection, collection2)

	returnedCollections, removed = collections.Remove(Denom2)
	require.True(t, removed)
	require.NotEqual(t, returnedCollections.String(), collections.String())
	require.Equal(t, 1, len(returnedCollections))

	foundCollection, found := returnedCollections.Find(Denom2)
	require.False(t, found)
	require.Empty(t, foundCollection)
}

func TestCollectionsStringMethod(t *testing.T) {
	collections := NewCollections()
	require.Equal(t, collections.String(), "")

	reserve := sdk.NewInt(100)
	var subTokenIDs []int64

	testNFT := NewBaseNFT(ID1, Addrs[0], Addrs[0], TokenURI1, reserve, subTokenIDs, true)
	nfts := NewNFTs(testNFT)
	collection := NewCollection(Denom1, nfts)

	testNFT2 := NewBaseNFT(ID2, Addrs[1], Addrs[1], TokenURI2, reserve, subTokenIDs, true)
	nfts2 := NewNFTs(testNFT2)
	collection2 := NewCollection(Denom2, nfts2)

	collections = NewCollections(collection, collection2)
	require.Equal(t, fmt.Sprintf(`Denom: 				%s
NFTs:

ID:				%s
Owners:			%s 
TokenURI:		%s
Denom: 				%s
NFTs:

ID:				%s
Owners:			%s 
TokenURI:		%s`, Denom1, ID1, Addrs[0].String(), TokenURI1,
		Denom2, ID2, Addrs[1].String(), TokenURI2), collections.String())
}

func TestCollectionsEmptyMethod(t *testing.T) {
	collections := NewCollections()
	require.True(t, collections.Empty())

	reserve := sdk.NewInt(100)
	var subTokenIDs []int64

	testNFT := NewBaseNFT(ID1, Addrs[0], Addrs[0], TokenURI1, reserve, subTokenIDs, true)
	nfts := NewNFTs(testNFT)
	collection := NewCollection(Denom1, nfts)

	collections = NewCollections(collection)
	require.False(t, collections.Empty())
}

func TestCollectionsSortInterface(t *testing.T) {
	reserve := sdk.NewInt(100)
	var subTokenIDs []int64
	testNFT := NewBaseNFT(ID1, Addrs[0], Addrs[0], TokenURI1, reserve, subTokenIDs, true)
	nfts := NewNFTs(testNFT)
	collection := NewCollection(Denom1, nfts)

	testNFT2 := NewBaseNFT(ID2, Addrs[1], Addrs[1], TokenURI2, reserve, subTokenIDs, true)
	nfts2 := NewNFTs(testNFT2)
	collection2 := NewCollection(Denom2, nfts2)

	collections := NewCollections(collection, collection2)
	require.Equal(t, 2, collections.Len())

	require.True(t, collections.Less(0, 1))
	require.False(t, collections.Less(1, 0))

	collections.Swap(0, 1)
	require.False(t, collections.Less(0, 1))
	require.True(t, collections.Less(1, 0))

	collections.Sort()
	require.True(t, collections.Less(0, 1))
	require.False(t, collections.Less(1, 0))
}

func TestCollectionMarshalAndUnmarshalJSON(t *testing.T) {
	reserve := sdk.NewInt(100)
	subTokenIDs := []int64{}

	testNFT := NewBaseNFT(ID1, Addrs[0], Addrs[0], TokenURI1, reserve, subTokenIDs, true)
	nfts := NewNFTs(testNFT)
	collection := NewCollection(Denom1, nfts)

	testNFT2 := NewBaseNFT(ID2, Addrs[1], Addrs[1], TokenURI2, reserve, subTokenIDs, true)
	nfts2 := NewNFTs(testNFT2)
	collection2 := NewCollection(Denom2, nfts2)

	collections := NewCollections(collection, collection2)

	bz, err := collections.MarshalJSON()
	require.NoError(t, err)
	require.Equal(t, string(bz), fmt.Sprintf(`{"%s":{"denom":"%s","nfts":{"%s":{"id":"%s","owners":{"owners":[{"address":"%s","sub_token_ids":%v}]},"creator":"%s","token_uri":"%s","reserve":"%s","allow_mint":%t}}},"%s":{"denom":"%s","nfts":{"%s":{"id":"%s","owners":{"owners":[{"address":"%s","sub_token_ids":%v}]},"creator":"%s","token_uri":"%s","reserve":"%s","allow_mint":%t}}}}`,
		Denom1, Denom1, ID1, ID1, Addrs[0].String(), subTokenIDs, Addrs[0].String(), TokenURI1, reserve.String(), testNFT.GetAllowMint(),
		Denom2, Denom2, ID2, ID2, Addrs[1].String(), subTokenIDs, Addrs[1].String(), TokenURI2, reserve.String(), testNFT2.GetAllowMint(),
	))

	var newCollections Collections
	err = newCollections.UnmarshalJSON(bz)
	require.NoError(t, err)

	err = newCollections.UnmarshalJSON([]byte{})
	require.Error(t, err)
}
