package nft

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"bitbucket.org/decimalteam/go-node/x/nft/internal/types"
)

const (
	moduleKey            = "module"
	denom                = "denom"
	nftID                = "nft_id"
	sender               = "sender"
	recipient            = "recipient"
	tokenURI             = "token_uri"
	amount               = "amount"
	subTokenIdStartRange = "sub_token_id_start_range"
)

func TestInvalidMsg(t *testing.T) {
	ctx, _, NFTKeeper := createTestApp(t, false)
	h := GenericHandler(NFTKeeper)
	_, err := h(ctx, sdk.NewTestMsg())

	require.Error(t, err)
}

func TestTransferNFTMsg(t *testing.T) {
	ctx, _, NFTKeeper := createTestApp(t, false)
	h := GenericHandler(NFTKeeper)

	// An NFT to be transferred
	reserve := sdk.NewInt(100)
	basenft := types.NewBaseNFT(ID1, Addrs[0], Addrs[0], TokenURI1, reserve, []int64{}, true)

	// Define MsgTransferNft
	transferNftMsg := types.NewMsgTransferNFT(Addrs[0], Addrs[1], Denom1, ID1, []int64{})

	// handle should fail trying to transfer NFT that doesn't exist
	res, err := h(ctx, transferNftMsg)
	require.Error(t, err)

	// Create token (collection and owner)
	_, err = NFTKeeper.MintNFT(ctx, Denom1, basenft.GetID(), basenft.GetReserve(), sdk.NewInt(1), basenft.GetCreator(), Addrs[0], basenft.GetTokenURI(), basenft.GetAllowMint())
	require.Nil(t, err)
	require.True(t, CheckInvariants(NFTKeeper, ctx))

	// handle should succeed when nft exists and is transferred by owner
	res, err = h(ctx, transferNftMsg)
	require.NoError(t, err)
	require.True(t, CheckInvariants(NFTKeeper, ctx))

	// event events should be emitted correctly
	for _, event := range res.Events {
		for _, attribute := range event.Attributes {
			value := string(attribute.Value)
			switch key := string(attribute.Key); key {
			case moduleKey:
				require.Equal(t, value, types.ModuleName)
			case denom:
				require.Equal(t, value, Denom1)
			case nftID:
				require.Equal(t, value, ID1)
			case tokenURI:
				require.Equal(t, value, TokenURI1)
			case sender:
				require.Equal(t, value, Addrs[0].String())
			case recipient:
				// require.Equal(t, value, Addrs[0].String())
			case amount:
			default:
				require.Fail(t, fmt.Sprintf("unrecognized event %s", key))
			}
		}
	}

	// nft should have been transferred as a result of the message
	nftAfterwards, err := NFTKeeper.GetNFT(ctx, Denom1, ID1)
	require.NoError(t, err)
	require.Equal(t, nftAfterwards.GetOwners().GetOwners()[1].GetAddress().String(), Addrs[1].String())

	transferNftMsg = types.NewMsgTransferNFT(Addrs[1], Addrs[2], Denom1, ID1, []int64{})

	// handle should succeed when nft exists and is transferred by owner
	res, err = h(ctx, transferNftMsg)
	require.NoError(t, err)
	require.True(t, CheckInvariants(NFTKeeper, ctx))

	// Create token (collection and owner)
	_, err = NFTKeeper.MintNFT(ctx,
		Denom2, basenft.GetID(),
		basenft.GetReserve(), sdk.NewInt(100),
		basenft.GetCreator(),
		Addrs[1],
		basenft.GetTokenURI(), basenft.GetAllowMint(),
	)
	require.Nil(t, err)
	require.True(t, CheckInvariants(NFTKeeper, ctx))

	transferNftMsg = types.NewMsgTransferNFT(Addrs[1], Addrs[2], Denom2, ID1, []int64{})

	// handle should succeed when nft exists and is transferred by owner
	res, err = h(ctx, transferNftMsg)
	require.NoError(t, err)
	require.True(t, CheckInvariants(NFTKeeper, ctx))
}

func TestEditNFTMetadataMsg(t *testing.T) {
	ctx, _, NFTKeeper := createTestApp(t, false)
	h := GenericHandler(NFTKeeper)

	reserve := sdk.NewInt(101)

	// An NFT to be edited
	basenft := types.NewBaseNFT(ID1, Addrs[0], Addrs[0], TokenURI1, reserve, []int64{}, true)

	// Create token (collection and address)
	_, err := NFTKeeper.MintNFT(ctx, Denom1, basenft.GetID(), basenft.GetReserve(), sdk.NewInt(1), basenft.GetCreator(), Addrs[0], basenft.GetTokenURI(), basenft.GetAllowMint())

	require.Nil(t, err)

	// Define MsgTransferNft
	failingEditNFTMetadata := types.NewMsgEditNFTMetadata(Addrs[0], ID1, Denom2, TokenURI2)

	res, err := h(ctx, failingEditNFTMetadata)
	require.Error(t, err)

	// Define MsgTransferNft
	editNFTMetadata := types.NewMsgEditNFTMetadata(Addrs[0], ID1, Denom1, TokenURI2)

	res, err = h(ctx, editNFTMetadata)
	require.NoError(t, err)

	// event events should be emitted correctly
	for _, event := range res.Events {
		for _, attribute := range event.Attributes {
			value := string(attribute.Value)
			switch key := string(attribute.Key); key {
			case moduleKey:
				require.Equal(t, value, types.ModuleName)
			case denom:
				require.Equal(t, value, Denom1)
			case nftID:
				require.Equal(t, value, ID1)
			case sender:
				require.Equal(t, value, Addrs[0].String())
			case tokenURI:
				require.Equal(t, value, TokenURI2)
			case recipient:
				// require.Equal(t, value, Addrs[0].String())
			case amount:
				// require.Equal(t, value, reserve)
			default:
				require.Fail(t, fmt.Sprintf("unrecognized event %s", key))
			}
		}
	}

	nftAfterwards, err := NFTKeeper.GetNFT(ctx, Denom1, ID1)
	require.NoError(t, err)
	require.Equal(t, TokenURI2, nftAfterwards.GetTokenURI())
}

func TestMintNFTMsg(t *testing.T) {
	ctx, _, NFTKeeper := createTestApp(t, false)
	h := GenericHandler(NFTKeeper)

	// Define MsgMintNFT
	reserve := sdk.NewInt(101)
	mintNFT := types.NewMsgMintNFT(Addrs[0], Addrs[0], ID1, Denom1, TokenURI1, sdk.NewInt(1), reserve, false)

	// minting a token should succeed
	res, err := h(ctx, mintNFT)
	require.NoError(t, err)

	// event events should be emitted correctly
	for _, event := range res.Events {
		for _, attribute := range event.Attributes {
			value := string(attribute.Value)
			switch key := string(attribute.Key); key {
			case moduleKey:
				require.Equal(t, value, types.ModuleName)
			case denom:
				require.Equal(t, value, Denom1)
			case nftID:
				require.Equal(t, value, ID1)
			case sender:
				require.Equal(t, value, Addrs[0].String())
			case tokenURI:
				require.Equal(t, value, TokenURI1)
			case subTokenIdStartRange:
				require.Equal(t, value, ID1)
			case recipient:
				// require.Equal(t, value, Addrs[0].String())
			case amount:
				// require.Equal(t, value, reserve)
			default:
				require.Fail(t, fmt.Sprintf("unrecognized event %s", key))
			}
		}
	}

	nftAfterwards, err := NFTKeeper.GetNFT(ctx, Denom1, ID1)

	require.NoError(t, err)
	require.Equal(t, TokenURI1, nftAfterwards.GetTokenURI())

	// minting the same token should fail if allowMint=false
	res, err = h(ctx, mintNFT)
	require.Error(t, err)

	require.True(t, CheckInvariants(NFTKeeper, ctx))
}

func TestBurnNFTMsg(t *testing.T) {
	ctx, _, NFTKeeper := createTestApp(t, false)
	h := GenericHandler(NFTKeeper)

	reserve := sdk.NewInt(100)
	// An NFT to be burned
	basenft := types.NewBaseNFT(ID1, Addrs[0], Addrs[0], TokenURI1, reserve, []int64{1, 2, 3}, true)

	// Create token (collection and address)
	_, err := NFTKeeper.MintNFT(ctx, Denom1, basenft.GetID(), basenft.GetReserve(), sdk.NewInt(3), basenft.GetCreator(), Addrs[0], basenft.GetTokenURI(), basenft.GetAllowMint())
	require.Nil(t, err)

	exists := NFTKeeper.IsNFT(ctx, Denom1, ID1)
	require.True(t, exists)

	// burning a non-existent NFT should fail
	failBurnNFT := types.NewMsgBurnNFT(Addrs[0], ID2, Denom1, []int64{4})
	res, err := h(ctx, failBurnNFT)
	require.Error(t, err)

	// NFT should still exist
	exists = NFTKeeper.IsNFT(ctx, Denom1, ID1)
	require.True(t, exists)

	// burning the NFt should succeed
	burnNFT := types.NewMsgBurnNFT(Addrs[0], ID1, Denom1, []int64{2})

	res, err = h(ctx, burnNFT)
	require.NoError(t, err)

	// event events should be emitted correctly
	for _, event := range res.Events {
		for _, attribute := range event.Attributes {
			if event.Type != sdk.EventTypeMessage || event.Type != types.EventTypeBurnNFT {
				continue
			}
			value := string(attribute.Value)
			switch key := string(attribute.Key); key {
			case moduleKey:
				require.Equal(t, value, types.ModuleName)
			case denom:
				require.Equal(t, value, Denom1)
			case nftID:
				require.Equal(t, value, ID1)
			case sender:
				require.Equal(t, value, Addrs[0].String())
			case recipient:
				// require.Equal(t, value, Addrs[0].String())
			case amount:
				// require.Equal(t, value, reserve)
			default:
				require.Fail(t, fmt.Sprintf("unrecognized event %s", key))
			}
		}
	}

	nft, err := NFTKeeper.GetNFT(ctx, Denom1, ID1)
	require.NoError(t, err)
	require.Equal(t, []int64{1, 3}, nft.GetOwners().GetOwners()[0].GetSubTokenIDs())

	// the NFT should not exist after burn
	exists = NFTKeeper.IsNFT(ctx, Denom1, ID1)
	require.True(t, exists)

	ownerReturned := NFTKeeper.GetOwner(ctx, Addrs[0])
	require.Equal(t, 1, ownerReturned.Supply())

	//require.True(t, CheckInvariants(NFTKeeper, ctx))
}

func TestUniqueTokenURI(t *testing.T) {
	ctx, _, nftKeeper := createTestApp(t, false)

	reserve := sdk.NewInt(100)

	const tokenURI1 = "tokenURI1"
	const tokenURI2 = "tokenURI2"

	msg := types.NewMsgMintNFT(Addrs[0], Addrs[0], "token1", "denom1", tokenURI1, sdk.NewInt(1), reserve, true)
	_, err := HandleMsgMintNFT(ctx, msg, nftKeeper)
	require.NoError(t, err)

	msg = types.NewMsgMintNFT(Addrs[0], Addrs[0], "token1", "denom1", tokenURI1, sdk.NewInt(1), reserve, true)
	_, err = HandleMsgMintNFT(ctx, msg, nftKeeper)
	require.NoError(t, err)

	msg = types.NewMsgMintNFT(Addrs[0], Addrs[0], "token2", "denom1", tokenURI2, sdk.NewInt(1), reserve, true)
	_, err = HandleMsgMintNFT(ctx, msg, nftKeeper)
	require.NoError(t, err)

	msg = types.NewMsgMintNFT(Addrs[0], Addrs[0], "token3", "denom1", tokenURI1, sdk.NewInt(1), reserve, true)
	_, err = HandleMsgMintNFT(ctx, msg, nftKeeper)
	require.Error(t, types.ErrNotUniqueTokenURI(), err)
}
