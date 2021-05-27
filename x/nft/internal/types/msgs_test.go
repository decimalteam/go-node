package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"testing"

	"github.com/stretchr/testify/require"
)

// ---------------------------------------- Msgs ---------------------------------------------------

func TestNewMsgTransferNFT(t *testing.T) {
	newMsgTransferNFT := NewMsgTransferNFT(Addrs[0], Addrs[1],
		fmt.Sprintf("     %s     ", Denom1),
		fmt.Sprintf("     %s     ", ID1),
		[]sdk.Int{},
	)
	require.Equal(t, newMsgTransferNFT.Sender, Addrs[0])
	require.Equal(t, newMsgTransferNFT.Recipient, Addrs[1])
	require.Equal(t, newMsgTransferNFT.Denom, Denom1)
	require.Equal(t, newMsgTransferNFT.ID, ID1)
}

func TestMsgTransferNFTValidateBasicMethod(t *testing.T) {
	newMsgTransferNFT := NewMsgTransferNFT(Addrs[0], Addrs[1], "", ID1, []sdk.Int{})
	err := newMsgTransferNFT.ValidateBasic()
	require.Error(t, err)

	newMsgTransferNFT = NewMsgTransferNFT(Addrs[0], Addrs[1], Denom1, "", []sdk.Int{})
	err = newMsgTransferNFT.ValidateBasic()
	require.Error(t, err)

	newMsgTransferNFT = NewMsgTransferNFT(nil, Addrs[1], Denom1, "", []sdk.Int{})
	err = newMsgTransferNFT.ValidateBasic()
	require.Error(t, err)

	newMsgTransferNFT = NewMsgTransferNFT(Addrs[0], nil, Denom1, "", []sdk.Int{})
	err = newMsgTransferNFT.ValidateBasic()
	require.Error(t, err)

	newMsgTransferNFT = NewMsgTransferNFT(Addrs[0], Addrs[1], Denom1, ID1, []sdk.Int{})
	err = newMsgTransferNFT.ValidateBasic()
	require.NoError(t, err)
}

func TestMsgTransferNFTGetSignBytesMethod(t *testing.T) {
	newMsgTransferNFT := NewMsgTransferNFT(Addrs[0], Addrs[1], Denom1, ID1, []sdk.Int{})
	sortedBytes := newMsgTransferNFT.GetSignBytes()

	require.Equal(t, string(sortedBytes), fmt.Sprintf(`{"type":"nft/msg_transfer","value":{"denom":"%s","id":"%s","recipient":"%s","sender":"%s","sub_token_ids":%v}}`,
		Denom1, ID1, Addrs[1], Addrs[0], newMsgTransferNFT.SubTokenIDs,
	))
}

func TestMsgTransferNFTGetSignersMethod(t *testing.T) {
	newMsgTransferNFT := NewMsgTransferNFT(Addrs[0], Addrs[1], Denom1, ID1, []sdk.Int{})
	signers := newMsgTransferNFT.GetSigners()
	require.Equal(t, 1, len(signers))
	require.Equal(t, Addrs[0].String(), signers[0].String())
}

func TestNewMsgEditNFTMetadata(t *testing.T) {
	newMsgEditNFTMetadata := NewMsgEditNFTMetadata(Addrs[0],
		fmt.Sprintf("     %s     ", ID1),
		fmt.Sprintf("     %s     ", Denom1),
		fmt.Sprintf("     %s     ", TokenURI1))

	require.Equal(t, newMsgEditNFTMetadata.Sender.String(), Addrs[0].String())
	require.Equal(t, newMsgEditNFTMetadata.ID, ID1)
	require.Equal(t, newMsgEditNFTMetadata.Denom, Denom1)
	require.Equal(t, newMsgEditNFTMetadata.TokenURI, TokenURI1)
}

func TestMsgEditNFTMetadataValidateBasicMethod(t *testing.T) {
	newMsgEditNFTMetadata := NewMsgEditNFTMetadata(nil, ID1, Denom1, TokenURI1)

	err := newMsgEditNFTMetadata.ValidateBasic()
	require.Error(t, err)

	newMsgEditNFTMetadata = NewMsgEditNFTMetadata(Addrs[0], "", Denom1, TokenURI1)
	err = newMsgEditNFTMetadata.ValidateBasic()
	require.Error(t, err)

	newMsgEditNFTMetadata = NewMsgEditNFTMetadata(Addrs[0], ID1, "", TokenURI1)
	err = newMsgEditNFTMetadata.ValidateBasic()
	require.Error(t, err)

	newMsgEditNFTMetadata = NewMsgEditNFTMetadata(Addrs[0], ID1, Denom1, TokenURI1)
	err = newMsgEditNFTMetadata.ValidateBasic()
	require.NoError(t, err)
}

func TestMsgEditNFTMetadataGetSignBytesMethod(t *testing.T) {
	newMsgEditNFTMetadata := NewMsgEditNFTMetadata(Addrs[0], ID1, Denom1, TokenURI1)
	sortedBytes := newMsgEditNFTMetadata.GetSignBytes()
	require.Equal(t, string(sortedBytes), fmt.Sprintf(`{"type":"nft/msg_edit_metadata","value":{"denom":"%s","id":"%s","sender":"%s","token_uri":"%s"}}`,
		Denom1, ID1, Addrs[0].String(), TokenURI1,
	))
}

func TestMsgEditNFTMetadataGetSignersMethod(t *testing.T) {
	newMsgEditNFTMetadata := NewMsgEditNFTMetadata(Addrs[0], ID1, Denom1, TokenURI1)
	signers := newMsgEditNFTMetadata.GetSigners()
	require.Equal(t, 1, len(signers))
	require.Equal(t, Addrs[0].String(), signers[0].String())
}

func TestNewMsgMintNFT(t *testing.T) {
	newMsgMintNFT := NewMsgMintNFT(Addrs[0], Addrs[1],
		fmt.Sprintf("     %s     ", ID1),
		fmt.Sprintf("     %s     ", Denom1),
		fmt.Sprintf("     %s     ", TokenURI1),
		sdk.NewInt(1),
		sdk.NewInt(1),
		true,
	)

	require.Equal(t, newMsgMintNFT.Sender.String(), Addrs[0].String())
	require.Equal(t, newMsgMintNFT.Recipient.String(), Addrs[1].String())
	require.Equal(t, newMsgMintNFT.ID, ID1)
	require.Equal(t, newMsgMintNFT.Denom, Denom1)
	require.Equal(t, newMsgMintNFT.TokenURI, TokenURI1)
}

func TestMsgMsgMintNFTValidateBasicMethod(t *testing.T) {
	newMsgMintNFT := NewMsgMintNFT(nil, Addrs[1], ID1, Denom1, TokenURI1, sdk.NewInt(1), sdk.NewInt(100), true)
	err := newMsgMintNFT.ValidateBasic()
	require.Error(t, err)

	newMsgMintNFT = NewMsgMintNFT(Addrs[0], nil, ID1, Denom1, TokenURI1, sdk.NewInt(1), sdk.NewInt(100), true)
	err = newMsgMintNFT.ValidateBasic()
	require.Error(t, err)

	newMsgMintNFT = NewMsgMintNFT(Addrs[0], Addrs[1], "", Denom1, TokenURI1, sdk.NewInt(1), sdk.NewInt(100), true)
	err = newMsgMintNFT.ValidateBasic()
	require.Error(t, err)

	newMsgMintNFT = NewMsgMintNFT(Addrs[0], Addrs[1], ID1, "", TokenURI1, sdk.NewInt(1), sdk.NewInt(100), true)
	err = newMsgMintNFT.ValidateBasic()
	require.Error(t, err)

	newMsgMintNFT = NewMsgMintNFT(Addrs[0], Addrs[1], ID1, Denom1, TokenURI1, sdk.NewInt(1), sdk.NewInt(99), true)
	err = newMsgMintNFT.ValidateBasic()
	require.Error(t, err)

	newMsgMintNFT = NewMsgMintNFT(Addrs[0], Addrs[1], ID1, Denom1, TokenURI1, sdk.NewInt(1), sdk.NewInt(100), true)
	err = newMsgMintNFT.ValidateBasic()
	require.NoError(t, err)
}

func TestMsgMintNFTGetSignBytesMethod(t *testing.T) {
	newMsgMintNFT := NewMsgMintNFT(Addrs[0], Addrs[1], ID1, Denom1, TokenURI1, sdk.NewInt(1), sdk.NewInt(100), true)
	sortedBytes := newMsgMintNFT.GetSignBytes()
	require.Equal(t, string(sortedBytes), fmt.Sprintf(`{"type":"nft/msg_mint","value":{"allow_mint":%t,"denom":"%s","id":"%s","quantity":"%v","recipient":"%s","reserve":"%v","sender":"%s","token_uri":"%s"}}`,
		newMsgMintNFT.AllowMint, Denom1, ID1, newMsgMintNFT.Quantity, Addrs[1].String(), newMsgMintNFT.Reserve, Addrs[0].String(), TokenURI1,
	))
}

func TestMsgMintNFTGetSignersMethod(t *testing.T) {
	newMsgMintNFT := NewMsgMintNFT(Addrs[0], Addrs[1], ID1, Denom1, TokenURI1, sdk.NewInt(1), sdk.NewInt(100), true)
	signers := newMsgMintNFT.GetSigners()
	require.Equal(t, 1, len(signers))
	require.Equal(t, Addrs[0].String(), signers[0].String())
}

func TestNewMsgBurnNFT(t *testing.T) {
	newMsgBurnNFT := NewMsgBurnNFT(Addrs[0],
		fmt.Sprintf("     %s     ", ID1),
		fmt.Sprintf("     %s     ", Denom1),
		[]sdk.Int{},
	)

	require.Equal(t, newMsgBurnNFT.Sender.String(), Addrs[0].String())
	require.Equal(t, newMsgBurnNFT.ID, ID1)
	require.Equal(t, newMsgBurnNFT.Denom, Denom1)
}

func TestMsgMsgBurnNFTValidateBasicMethod(t *testing.T) {
	newMsgBurnNFT := NewMsgBurnNFT(nil, ID1, Denom1, []sdk.Int{})
	err := newMsgBurnNFT.ValidateBasic()
	require.Error(t, err)

	newMsgBurnNFT = NewMsgBurnNFT(Addrs[0], "", Denom1, []sdk.Int{})
	err = newMsgBurnNFT.ValidateBasic()
	require.Error(t, err)

	newMsgBurnNFT = NewMsgBurnNFT(Addrs[0], ID1, "", []sdk.Int{})
	err = newMsgBurnNFT.ValidateBasic()
	require.Error(t, err)

	newMsgBurnNFT = NewMsgBurnNFT(Addrs[0], ID1, Denom1, []sdk.Int{})
	err = newMsgBurnNFT.ValidateBasic()
	require.NoError(t, err)
}

func TestMsgBurnNFTGetSignBytesMethod(t *testing.T) {
	newMsgBurnNFT := NewMsgBurnNFT(Addrs[0], ID1, Denom1, []sdk.Int{})
	sortedBytes := newMsgBurnNFT.GetSignBytes()
	require.Equal(t, string(sortedBytes), fmt.Sprintf(`{"type":"nft/msg_burn","value":{"denom":"%s","id":"%s","sender":"%s","sub_token_ids":%v}}`,
		Denom1, ID1, Addrs[0].String(), newMsgBurnNFT.SubTokenIDs,
	))
}

func TestMsgBurnNFTGetSignersMethod(t *testing.T) {
	newMsgBurnNFT := NewMsgBurnNFT(Addrs[0], ID1, Denom1, []sdk.Int{})
	signers := newMsgBurnNFT.GetSigners()
	require.Equal(t, 1, len(signers))
	require.Equal(t, Addrs[0].String(), signers[0].String())
}
