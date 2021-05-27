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
	quantity := sdk.NewInt(1)
	basenft := types.NewBaseNFT(ID1, Addrs[0], Addrs[0], TokenURI1, quantity, sdk.NewInt(101), true)

	// Define MsgTransferNft
	transferNftMsg := types.NewMsgTransferNFT(Addrs[0], Addrs[1], Denom1, ID1, []sdk.Int{})

	// handle should fail trying to transfer NFT that doesn't exist
	res, err := h(ctx, transferNftMsg)
	require.Error(t, err)

	// Create token (collection and owner)
	_, err = NFTKeeper.MintNFT(ctx, Denom1, basenft)
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

	transferNftMsg = types.NewMsgTransferNFT(Addrs[1], Addrs[2], Denom1, ID1, []sdk.Int{})

	// handle should succeed when nft exists and is transferred by owner
	res, err = h(ctx, transferNftMsg)
	require.NoError(t, err)
	require.True(t, CheckInvariants(NFTKeeper, ctx))

	// Create token (collection and owner)
	_, err = NFTKeeper.MintNFT(ctx, Denom2, basenft)
	require.Nil(t, err)
	require.True(t, CheckInvariants(NFTKeeper, ctx))

	transferNftMsg = types.NewMsgTransferNFT(Addrs[1], Addrs[2], Denom2, ID1, []sdk.Int{})

	// handle should succeed when nft exists and is transferred by owner
	res, err = h(ctx, transferNftMsg)
	require.NoError(t, err)
	require.True(t, CheckInvariants(NFTKeeper, ctx))
}

func TestEditNFTMetadataMsg(t *testing.T) {
	ctx, _, NFTKeeper := createTestApp(t, false)
	h := GenericHandler(NFTKeeper)

	// An NFT to be edited
	basenft := types.NewBaseNFT(ID1, Addrs[0], Addrs[0], TokenURI1, sdk.NewInt(1), sdk.NewInt(101), true)

	// Create token (collection and address)
	_, err := NFTKeeper.MintNFT(ctx, Denom1, basenft)
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
			case amount:
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
	mintNFT := types.NewMsgMintNFT(Addrs[0], Addrs[0], ID1, Denom1, TokenURI1, sdk.NewInt(1), sdk.NewInt(101), false)

	s1 := mintNFT.Sender.String()
	s2 := mintNFT.Recipient.String()
	fmt.Printf("%s %s", s1, s2)
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
			case recipient:
			case amount:
			case subTokenIdStartRange:
				//require.Equal(t, value, )
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

	// An NFT to be burned
	basenft := types.NewBaseNFT(ID1, Addrs[0], Addrs[0], TokenURI1, sdk.NewInt(1), sdk.NewInt(101), true)

	// Create token (collection and address)
	_, err := NFTKeeper.MintNFT(ctx, Denom1, basenft)
	require.Nil(t, err)

	exists := NFTKeeper.IsNFT(ctx, Denom1, ID1)
	require.True(t, exists)

	// burning a non-existent NFT should fail
	failBurnNFT := types.NewMsgBurnNFT(Addrs[0], ID2, Denom1, []sdk.Int{})
	res, err := h(ctx, failBurnNFT)
	require.Error(t, err)

	// NFT should still exist
	exists = NFTKeeper.IsNFT(ctx, Denom1, ID1)
	require.True(t, exists)

	// burning the NFt should succeed
	burnNFT := types.NewMsgBurnNFT(Addrs[0], ID1, Denom1, []sdk.Int{})

	res, err = h(ctx, burnNFT)
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
			case recipient:
			case amount:
			default:
				require.Fail(t, fmt.Sprintf("unrecognized event %s", key))
			}
		}
	}

	// the NFT should not exist after burn
	exists = NFTKeeper.IsNFT(ctx, Denom1, ID1)
	require.False(t, exists)

	ownerReturned := NFTKeeper.GetOwner(ctx, Addrs[0])
	require.Equal(t, 0, ownerReturned.Supply())

	require.True(t, CheckInvariants(NFTKeeper, ctx))
}
