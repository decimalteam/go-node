package nft_test

/*
import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"bitbucket.org/decimalteam/go-node/x/nft"
	"bitbucket.org/decimalteam/go-node/x/nft/internal/types"
)

const (
	module    = "module"
	denom     = "denom"
	nftID     = "nft-id"
	sender    = "sender"
	recipient = "recipient"
	tokenURI  = "token-uri"
)

func TestInvalidMsg(t *testing.T) {
	ctx, _, NFTKeeper := createTestApp(t, false)
	h := nft.GenericHandler(NFTKeeper)
	_, err := h(ctx, sdk.NewTestMsg())

	require.Error(t, err)
}

func TestTransferNFTMsg(t *testing.T) {
	ctx, _, NFTKeeper := createTestApp(t, false)
	h := nft.GenericHandler(NFTKeeper)

	// An NFT to be transferred
	basenft := types.NewBaseNFT(id, address, address, tokenURI, sdk.NewInt(1), sdk.NewInt(101), true)

	// Define MsgTransferNft
	transferNftMsg := types.NewMsgTransferNFT(address, address2, denom, id, []sdk.Int{})

	// handle should fail trying to transfer NFT that doesn't exist
	res, err := h(ctx, transferNftMsg)
	require.Error(t, err)

	// Create token (collection and owner)
	err = NFTKeeper.MintNFT(ctx, denom, basenft)
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
			case module:
				require.Equal(t, value, types.ModuleName)
			case denom:
				require.Equal(t, value, denom)
			case nftID:
				require.Equal(t, value, id)
			case sender:
				require.Equal(t, value, address.String())
			case recipient:
				require.Equal(t, value, address2.String())
			default:
				require.Fail(t, fmt.Sprintf("unrecognized event %s", key))
			}
		}
	}

	// nft should have been transferred as a result of the message
	nftAfterwards, err := NFTKeeper.GetNFT(ctx, denom, id)
	require.NoError(t, err)
	require.True(t, nftAfterwards.GetOwners().GetOwners()[0].GetAddress().Equals(address2))

	transferNftMsg = types.NewMsgTransferNFT(address2, address3, denom, id, []sdk.Int{})

	// handle should succeed when nft exists and is transferred by owner
	res, err = h(ctx, transferNftMsg)
	require.NoError(t, err)
	require.True(t, CheckInvariants(NFTKeeper, ctx))

	// Create token (collection and owner)
	err = NFTKeeper.MintNFT(ctx, denom2, basenft)
	require.Nil(t, err)
	require.True(t, CheckInvariants(NFTKeeper, ctx))

	transferNftMsg = types.NewMsgTransferNFT(address2, address3, denom2, id, []sdk.Int{})

	// handle should succeed when nft exists and is transferred by owner
	res, err = h(ctx, transferNftMsg)
	require.NoError(t, err)
	require.True(t, CheckInvariants(NFTKeeper, ctx))
}

func TestEditNFTMetadataMsg(t *testing.T) {
	ctx, _, NFTKeeper := createTestApp(t, false)
	h := nft.GenericHandler(NFTKeeper)

	// An NFT to be edited
	basenft := types.NewBaseNFT(id, address, address, tokenURI, sdk.NewInt(1), sdk.NewInt(101), true)

	// Create token (collection and address)
	err := NFTKeeper.MintNFT(ctx, denom, basenft)
	require.Nil(t, err)

	// Define MsgTransferNft
	failingEditNFTMetadata := types.NewMsgEditNFTMetadata(address, id, denom2, tokenURI2)

	res, err := h(ctx, failingEditNFTMetadata)
	require.Error(t, err)

	// Define MsgTransferNft
	editNFTMetadata := types.NewMsgEditNFTMetadata(address, id, denom, tokenURI2)

	res, err = h(ctx, editNFTMetadata)
	require.NoError(t, err)

	// event events should be emitted correctly
	for _, event := range res.Events {
		for _, attribute := range event.Attributes {
			value := string(attribute.Value)
			switch key := string(attribute.Key); key {
			case module:
				require.Equal(t, value, types.ModuleName)
			case denom:
				require.Equal(t, value, denom)
			case nftID:
				require.Equal(t, value, id)
			case sender:
				require.Equal(t, value, address.String())
			case tokenURI:
				require.Equal(t, value, tokenURI2)
			default:
				require.Fail(t, fmt.Sprintf("unrecognized event %s", key))
			}
		}
	}

	nftAfterwards, err := NFTKeeper.GetNFT(ctx, denom, id)
	require.NoError(t, err)
	require.Equal(t, tokenURI2, nftAfterwards.GetTokenURI())
}

func TestMintNFTMsg(t *testing.T) {
	ctx, _, NFTKeeper := createTestApp(t, false)
	h := nft.GenericHandler(NFTKeeper)

	// Define MsgMintNFT
	mintNFT := types.NewMsgMintNFT(address, address, id, denom, tokenURI, sdk.NewInt(1), sdk.NewInt(101), true)

	// minting a token should succeed
	res, err := h(ctx, mintNFT)
	require.NoError(t, err)

	// event events should be emitted correctly
	for _, event := range res.Events {
		for _, attribute := range event.Attributes {
			value := string(attribute.Value)
			switch key := string(attribute.Key); key {
			case module:
				require.Equal(t, value, types.ModuleName)
			case denom:
				require.Equal(t, value, denom)
			case nftID:
				require.Equal(t, value, id)
			case sender:
				require.Equal(t, value, address.String())
			case recipient:
				require.Equal(t, value, address.String())
			case tokenURI:
				require.Equal(t, value, tokenURI)
			default:
				require.Fail(t, fmt.Sprintf("unrecognized event %s", key))
			}
		}
	}

	nftAfterwards, err := NFTKeeper.GetNFT(ctx, denom, id)

	require.NoError(t, err)
	require.Equal(t, tokenURI, nftAfterwards.GetTokenURI())

	// minting the same token should fail
	res, err = h(ctx, mintNFT)
	require.Error(t, err)

	require.True(t, CheckInvariants(NFTKeeper, ctx))
}

func TestBurnNFTMsg(t *testing.T) {
	ctx, _, NFTKeeper := createTestApp(t, false)
	h := nft.GenericHandler(NFTKeeper)

	// An NFT to be burned
	basenft := types.NewBaseNFT(id, address, address, tokenURI, sdk.NewInt(1), sdk.NewInt(101), true)

	// Create token (collection and address)
	err := NFTKeeper.MintNFT(ctx, denom, basenft)
	require.Nil(t, err)

	exists := NFTKeeper.IsNFT(ctx, denom, id)
	require.True(t, exists)

	// burning a non-existent NFT should fail
	failBurnNFT := types.NewMsgBurnNFT(address, id2, denom, []sdk.Int{})
	res, err := h(ctx, failBurnNFT)
	require.Error(t, err)

	// NFT should still exist
	exists = NFTKeeper.IsNFT(ctx, denom, id)
	require.True(t, exists)

	// burning the NFt should succeed
	burnNFT := types.NewMsgBurnNFT(address, id, denom, []sdk.Int{})

	res, err = h(ctx, burnNFT)
	require.NoError(t, err)

	// event events should be emitted correctly
	for _, event := range res.Events {
		for _, attribute := range event.Attributes {
			value := string(attribute.Value)
			switch key := string(attribute.Key); key {
			case module:
				require.Equal(t, value, types.ModuleName)
			case denom:
				require.Equal(t, value, denom)
			case nftID:
				require.Equal(t, value, id)
			case sender:
				require.Equal(t, value, address.String())
			default:
				require.Fail(t, fmt.Sprintf("unrecognized event %s", key))
			}
		}
	}

	// the NFT should not exist after burn
	exists = NFTKeeper.IsNFT(ctx, denom, id)
	require.False(t, exists)

	ownerReturned := NFTKeeper.GetOwner(ctx, address)
	require.Equal(t, 0, ownerReturned.Supply())

	require.True(t, CheckInvariants(NFTKeeper, ctx))
}*/
