package types

import (
	"bitbucket.org/decimalteam/go-node/x/nft/exported"
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

// ---------------------------------------- BaseNFT ---------------------------------------------------

func TestBaseNFTGetMethods(t *testing.T) {
	testNFT := NewBaseNFT(ID1, Addrs[0], Addrs[0], TokenURI1, sdk.NewInt(1), []int64{}, true)

	require.Equal(t, ID1, testNFT.GetID())
	require.Equal(t, Addrs[0], testNFT.GetOwners().GetOwners()[0].GetAddress())
	require.Equal(t, TokenURI1, testNFT.GetTokenURI())
}

func TestBaseNFTSetMethods(t *testing.T) {
	testNFT := NewBaseNFT(ID1, Addrs[0], Addrs[0], TokenURI1, sdk.NewInt(1), []int64{1}, true)

	testNFT = testNFT.SetOwners(
		&TokenOwners{
			Owners: []exported.TokenOwner{&TokenOwner{
				Address:     Addrs[1],
				SubTokenIDs: []int64{},
			}}})
	require.Equal(t, Addrs[1], testNFT.GetOwners().GetOwners()[0].GetAddress())

	testNFT = testNFT.EditMetadata(TokenURI2)
	require.Equal(t, TokenURI2, testNFT.GetTokenURI())
}

func TestBaseNFTStringFormat(t *testing.T) {
	reserve := sdk.NewInt(1)
	testNFT := NewBaseNFT(ID1, Addrs[0], Addrs[0], TokenURI1, reserve, []int64{}, true)
	expected := fmt.Sprintf(`ID:				%s
Owners:			%s 
TokenURI:		%s`,
		ID1, Addrs[0].String(), TokenURI1)
	require.Equal(t, expected, testNFT.String())
}

// ---------------------------------------- NFTs ---------------------------------------------------

func TestNewNFTs(t *testing.T) {
	emptyNFTs := NewNFTs()
	require.Equal(t, len(emptyNFTs), 0)

	testNFT := NewBaseNFT(ID1, Addrs[0], Addrs[0], TokenURI1, sdk.NewInt(1), []int64{}, true)
	oneNFTs := NewNFTs(testNFT)
	require.Equal(t, len(oneNFTs), 1)

	testNFT2 := NewBaseNFT(ID2, Addrs[0], Addrs[0], TokenURI1, sdk.NewInt(1), []int64{}, true)
	twoNFTs := NewNFTs(testNFT, testNFT2)
	require.Equal(t, len(twoNFTs), 2)
}

func TestNFTsAppendMethod(t *testing.T) {
	testNFT := NewBaseNFT(ID1, Addrs[0], Addrs[0], TokenURI1, sdk.NewInt(1), []int64{}, true)
	nfts := NewNFTs(testNFT)
	require.Equal(t, len(nfts), 1)

	testNFT2 := NewBaseNFT(ID2, Addrs[0], Addrs[0], TokenURI1, sdk.NewInt(1), []int64{}, true)
	nfts2 := NewNFTs(testNFT2)

	nfts = nfts.Append(nfts2...)
	require.Equal(t, len(nfts), 2)

	var id3 = string('3')
	var id4 = string('4')
	var id5 = string('5')
	testNFT3 := NewBaseNFT(id3, Addrs[0], Addrs[0], TokenURI1, sdk.NewInt(1), []int64{}, true)
	testNFT4 := NewBaseNFT(id4, Addrs[0], Addrs[0], TokenURI1, sdk.NewInt(1), []int64{}, true)
	testNFT5 := NewBaseNFT(id5, Addrs[0], Addrs[0], TokenURI1, sdk.NewInt(1), []int64{}, true)

	nfts3 := NewNFTs(testNFT5, testNFT3, testNFT4)
	nfts = nfts.Append(nfts3...)
	require.Equal(t, len(nfts), 5)

	nft, found := nfts.Find(ID2)
	require.True(t, found)
	require.Equal(t, nft.String(), testNFT2.String())

	nft, found = nfts.Find(id5)
	require.True(t, found)
	require.Equal(t, nft.String(), testNFT5.String())

	nft, found = nfts.Find(id3)
	require.True(t, found)
	require.Equal(t, nft.String(), testNFT3.String())
}

func TestNFTsFindMethod(t *testing.T) {
	testNFT := NewBaseNFT(ID1, Addrs[0], Addrs[0], TokenURI1, sdk.NewInt(1), []int64{}, true)
	testNFT2 := NewBaseNFT(ID2, Addrs[0], Addrs[0], TokenURI1, sdk.NewInt(1), []int64{}, true)

	var id3 = string('3')
	var id4 = string('4')
	var id5 = string('5')
	testNFT3 := NewBaseNFT(id3, Addrs[0], Addrs[0], TokenURI1, sdk.NewInt(1), []int64{}, true)
	testNFT4 := NewBaseNFT(id4, Addrs[0], Addrs[0], TokenURI1, sdk.NewInt(1), []int64{}, true)
	testNFT5 := NewBaseNFT(id5, Addrs[0], Addrs[0], TokenURI1, sdk.NewInt(1), []int64{}, true)

	nfts := NewNFTs(testNFT, testNFT3, testNFT4, testNFT5, testNFT2)
	nft, found := nfts.Find(ID1)
	require.True(t, found)
	require.Equal(t, nft.String(), testNFT.String())

	nft, found = nfts.Find(ID2)
	require.True(t, found)
	require.Equal(t, nft.String(), testNFT2.String())

	nft, found = nfts.Find(id3)
	require.True(t, found)
	require.Equal(t, nft.String(), testNFT3.String())

	nft, found = nfts.Find(id4)
	require.True(t, found)
	require.Equal(t, nft.String(), testNFT4.String())

	nft, found = nfts.Find(id5)
	require.True(t, found)
	require.Equal(t, nft.String(), testNFT5.String())

	var id6 = string('6')
	nft, found = nfts.Find(id6)
	require.False(t, found)
	require.Nil(t, nft)
}

func TestNFTsUpdateMethod(t *testing.T) {
	testNFT := NewBaseNFT(ID1, Addrs[0], Addrs[0], TokenURI1, sdk.NewInt(1), []int64{}, true)
	testNFT2 := NewBaseNFT(ID2, Addrs[0], Addrs[0], TokenURI1, sdk.NewInt(1), []int64{}, true)
	nfts := NewNFTs(testNFT)
	var success bool
	nfts, success = nfts.Update(ID1, testNFT2)
	require.True(t, success)

	nft, found := nfts.Find(ID2)
	require.True(t, found)
	require.Equal(t, nft.String(), testNFT2.String())

	nft, found = nfts.Find(ID1)
	require.False(t, found)
	require.Nil(t, nft)

	var returnedNFTs NFTs
	returnedNFTs, success = nfts.Update(ID1, testNFT2)
	require.False(t, success)
	require.Equal(t, returnedNFTs.String(), nfts.String())
}

func TestNFTsRemoveMethod(t *testing.T) {
	testNFT := NewBaseNFT(ID1, Addrs[0], Addrs[0], TokenURI1, sdk.NewInt(1), []int64{}, true)
	testNFT2 := NewBaseNFT(ID2, Addrs[0], Addrs[0], TokenURI1, sdk.NewInt(1), []int64{}, true)
	nfts := NewNFTs(testNFT, testNFT2)

	var success bool
	nfts, success = nfts.Remove(ID1)
	require.True(t, success)
	require.Equal(t, len(nfts), 1)

	nfts, success = nfts.Remove(ID2)
	require.True(t, success)
	require.Equal(t, len(nfts), 0)

	var returnedNFTs NFTs
	returnedNFTs, success = nfts.Remove(ID2)
	require.False(t, success)
	require.Equal(t, nfts.String(), returnedNFTs.String())
}

func TestNFTsStringMethod(t *testing.T) {
	quantity := sdk.NewInt(1)
	testNFT := NewBaseNFT(ID1, Addrs[0], Addrs[0], TokenURI1, quantity, []int64{}, true)
	nfts := NewNFTs(testNFT)
	require.Equal(t, nfts.String(), fmt.Sprintf(`ID:				%s
Owners:			%s 
TokenURI:		%s`, ID1, Addrs[0].String(), TokenURI1))
}

func TestNFTsEmptyMethod(t *testing.T) {
	nfts := NewNFTs()
	require.True(t, nfts.Empty())
	testNFT := NewBaseNFT(ID1, Addrs[0], Addrs[0], TokenURI1, sdk.NewInt(1), []int64{}, true)
	nfts = NewNFTs(testNFT)
	require.False(t, nfts.Empty())
}

func TestNFTsMarshalUnmarshalJSON(t *testing.T) {
	subTokenIDs := []int64{}
	testNFT := NewBaseNFT(ID1, Addrs[0], Addrs[0], TokenURI1, sdk.NewInt(1), subTokenIDs, true)
	nfts := NewNFTs(testNFT)
	bz, err := nfts.MarshalJSON()
	require.NoError(t, err)
	require.Equal(t, string(bz),
		fmt.Sprintf(`{"%s":{"id":"%s","owners":{"owners":[{"address":"%s","sub_token_ids":%v}]},"creator":"%s","token_uri":"%s","reserve":"1","allow_mint":true}}`,
			ID1, ID1, Addrs[0].String(), subTokenIDs, Addrs[0].String(), TokenURI1))

	var unmarshalledNFTs NFTs
	err = unmarshalledNFTs.UnmarshalJSON(bz)
	require.NoError(t, err)
	require.Equal(t, unmarshalledNFTs.String(), nfts.String())

	bz = []byte{}
	err = unmarshalledNFTs.UnmarshalJSON(bz)
	require.Error(t, err)
}

func TestNFTsSortInterface(t *testing.T) {
	testNFT := NewBaseNFT(ID1, Addrs[0], Addrs[0], TokenURI1, sdk.NewInt(1), []int64{}, true)
	testNFT2 := NewBaseNFT(ID2, Addrs[0], Addrs[0], TokenURI1, sdk.NewInt(1), []int64{}, true)

	nfts := NewNFTs(testNFT)
	require.Equal(t, nfts.Len(), 1)

	nfts = NewNFTs(testNFT, testNFT2)
	require.Equal(t, nfts.Len(), 2)

	require.True(t, nfts.Less(0, 1))
	require.False(t, nfts.Less(1, 0))

	nfts.Swap(0, 1)
	require.False(t, nfts.Less(0, 1))
	require.True(t, nfts.Less(1, 0))

	nfts.Sort()
	require.True(t, nfts.Less(0, 1))
	require.False(t, nfts.Less(1, 0))
}
