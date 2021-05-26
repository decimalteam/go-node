package types

/*
import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"testing"

	"github.com/stretchr/testify/require"
)

// ---------------------------------------- Msgs ---------------------------------------------------

func TestNewMsgTransferNFT(t *testing.T) {
	newMsgTransferNFT := NewMsgTransferNFT(address, address2,
		fmt.Sprintf("     %s     ", denom),
		fmt.Sprintf("     %s     ", id),
		[]sdk.Int{},
	)
	require.Equal(t, newMsgTransferNFT.Sender, address)
	require.Equal(t, newMsgTransferNFT.Recipient, address2)
	require.Equal(t, newMsgTransferNFT.Denom, denom)
	require.Equal(t, newMsgTransferNFT.ID, id)
}

func TestMsgTransferNFTValidateBasicMethod(t *testing.T) {
	newMsgTransferNFT := NewMsgTransferNFT(address, address2, "", id, []sdk.Int{})
	err := newMsgTransferNFT.ValidateBasic()
	require.Error(t, err)

	newMsgTransferNFT = NewMsgTransferNFT(address, address2, denom, "", []sdk.Int{})
	err = newMsgTransferNFT.ValidateBasic()
	require.Error(t, err)

	newMsgTransferNFT = NewMsgTransferNFT(nil, address2, denom, "", []sdk.Int{})
	err = newMsgTransferNFT.ValidateBasic()
	require.Error(t, err)

	newMsgTransferNFT = NewMsgTransferNFT(address, nil, denom, "", []sdk.Int{})
	err = newMsgTransferNFT.ValidateBasic()
	require.Error(t, err)

	newMsgTransferNFT = NewMsgTransferNFT(address, address2, denom, id, []sdk.Int{})
	err = newMsgTransferNFT.ValidateBasic()
	require.NoError(t, err)
}

func TestMsgTransferNFTGetSignBytesMethod(t *testing.T) {
	newMsgTransferNFT := NewMsgTransferNFT(address, address2, denom, id, []sdk.Int{})
	sortedBytes := newMsgTransferNFT.GetSignBytes()
	require.Equal(t, string(sortedBytes), fmt.Sprintf(`{"type":"cosmos-sdk/MsgTransferNFT","value":{"denom":"%s","id":"%s","recipient":"%s","sender":"%s"}}`,
		denom, id, address2, address,
	))
}

func TestMsgTransferNFTGetSignersMethod(t *testing.T) {
	newMsgTransferNFT := NewMsgTransferNFT(address, address2, denom, id, []sdk.Int{})
	signers := newMsgTransferNFT.GetSigners()
	require.Equal(t, 1, len(signers))
	require.Equal(t, address.String(), signers[0].String())
}

func TestNewMsgEditNFTMetadata(t *testing.T) {
	newMsgEditNFTMetadata := NewMsgEditNFTMetadata(address,
		fmt.Sprintf("     %s     ", id),
		fmt.Sprintf("     %s     ", denom),
		fmt.Sprintf("     %s     ", tokenURI))

	require.Equal(t, newMsgEditNFTMetadata.Sender.String(), address.String())
	require.Equal(t, newMsgEditNFTMetadata.ID, id)
	require.Equal(t, newMsgEditNFTMetadata.Denom, denom)
	require.Equal(t, newMsgEditNFTMetadata.TokenURI, tokenURI)
}

func TestMsgEditNFTMetadataValidateBasicMethod(t *testing.T) {
	newMsgEditNFTMetadata := NewMsgEditNFTMetadata(nil, id, denom, tokenURI)

	err := newMsgEditNFTMetadata.ValidateBasic()
	require.Error(t, err)

	newMsgEditNFTMetadata = NewMsgEditNFTMetadata(address, "", denom, tokenURI)
	err = newMsgEditNFTMetadata.ValidateBasic()
	require.Error(t, err)

	newMsgEditNFTMetadata = NewMsgEditNFTMetadata(address, id, "", tokenURI)
	err = newMsgEditNFTMetadata.ValidateBasic()
	require.Error(t, err)

	newMsgEditNFTMetadata = NewMsgEditNFTMetadata(address, id, denom, tokenURI)
	err = newMsgEditNFTMetadata.ValidateBasic()
	require.NoError(t, err)
}

func TestMsgEditNFTMetadataGetSignBytesMethod(t *testing.T) {
	newMsgEditNFTMetadata := NewMsgEditNFTMetadata(address, id, denom, tokenURI)
	sortedBytes := newMsgEditNFTMetadata.GetSignBytes()
	require.Equal(t, string(sortedBytes), fmt.Sprintf(`{"type":"cosmos-sdk/MsgEditNFTMetadata","value":{"denom":"%s","id":"%s","sender":"%s","token_uri":"%s"}}`,
		denom, id, address.String(), tokenURI,
	))
}

func TestMsgEditNFTMetadataGetSignersMethod(t *testing.T) {
	newMsgEditNFTMetadata := NewMsgEditNFTMetadata(address, id, denom, tokenURI)
	signers := newMsgEditNFTMetadata.GetSigners()
	require.Equal(t, 1, len(signers))
	require.Equal(t, address.String(), signers[0].String())
}

func TestNewMsgMintNFT(t *testing.T) {
	newMsgMintNFT := NewMsgMintNFT(address, address2,
		fmt.Sprintf("     %s     ", id),
		fmt.Sprintf("     %s     ", denom),
		fmt.Sprintf("     %s     ", tokenURI),
		sdk.NewInt(1),
		sdk.NewInt(1),
		true,
	)

	require.Equal(t, newMsgMintNFT.Sender.String(), address.String())
	require.Equal(t, newMsgMintNFT.Recipient.String(), address2.String())
	require.Equal(t, newMsgMintNFT.ID, id)
	require.Equal(t, newMsgMintNFT.Denom, denom)
	require.Equal(t, newMsgMintNFT.TokenURI, tokenURI)
}

func TestMsgMsgMintNFTValidateBasicMethod(t *testing.T) {
	newMsgMintNFT := NewMsgMintNFT(nil, address2, id, denom, tokenURI, sdk.NewInt(1), sdk.NewInt(100), true)
	err := newMsgMintNFT.ValidateBasic()
	require.Error(t, err)

	newMsgMintNFT = NewMsgMintNFT(address, nil, id, denom, tokenURI, sdk.NewInt(1), sdk.NewInt(100), true)
	err = newMsgMintNFT.ValidateBasic()
	require.Error(t, err)

	newMsgMintNFT = NewMsgMintNFT(address, address2, "", denom, tokenURI, sdk.NewInt(1), sdk.NewInt(100), true)
	err = newMsgMintNFT.ValidateBasic()
	require.Error(t, err)

	newMsgMintNFT = NewMsgMintNFT(address, address2, id, "", tokenURI, sdk.NewInt(1), sdk.NewInt(100), true)
	err = newMsgMintNFT.ValidateBasic()
	require.Error(t, err)

	newMsgMintNFT = NewMsgMintNFT(address, address2, id, denom, tokenURI, sdk.NewInt(1), sdk.NewInt(99), true)
	err = newMsgMintNFT.ValidateBasic()
	require.Error(t, err)

	newMsgMintNFT = NewMsgMintNFT(address, address2, id, denom, tokenURI, sdk.NewInt(1), sdk.NewInt(100), true)
	err = newMsgMintNFT.ValidateBasic()
	require.NoError(t, err)
}

func TestMsgMintNFTGetSignBytesMethod(t *testing.T) {
	newMsgMintNFT := NewMsgMintNFT(address, address2, id, denom, tokenURI, sdk.NewInt(1), sdk.NewInt(100), true)
	sortedBytes := newMsgMintNFT.GetSignBytes()
	require.Equal(t, string(sortedBytes), fmt.Sprintf(`{"type":"cosmos-sdk/MsgMintNFT","value":{"denom":"%s","id":"%s","recipient":"%s","sender":"%s","token_uri":"%s"}}`,
		denom, id, address2.String(), address.String(), tokenURI,
	))
}

func TestMsgMintNFTGetSignersMethod(t *testing.T) {
	newMsgMintNFT := NewMsgMintNFT(address, address2, id, denom, tokenURI, sdk.NewInt(1), sdk.NewInt(100), true)
	signers := newMsgMintNFT.GetSigners()
	require.Equal(t, 1, len(signers))
	require.Equal(t, address.String(), signers[0].String())
}

func TestNewMsgBurnNFT(t *testing.T) {
	newMsgBurnNFT := NewMsgBurnNFT(address,
		fmt.Sprintf("     %s     ", id),
		fmt.Sprintf("     %s     ", denom),
		[]sdk.Int{},
	)

	require.Equal(t, newMsgBurnNFT.Sender.String(), address.String())
	require.Equal(t, newMsgBurnNFT.ID, id)
	require.Equal(t, newMsgBurnNFT.Denom, denom)
}

func TestMsgMsgBurnNFTValidateBasicMethod(t *testing.T) {
	newMsgBurnNFT := NewMsgBurnNFT(nil, id, denom, []sdk.Int{})
	err := newMsgBurnNFT.ValidateBasic()
	require.Error(t, err)

	newMsgBurnNFT = NewMsgBurnNFT(address, "", denom, []sdk.Int{})
	err = newMsgBurnNFT.ValidateBasic()
	require.Error(t, err)

	newMsgBurnNFT = NewMsgBurnNFT(address, id, "", []sdk.Int{})
	err = newMsgBurnNFT.ValidateBasic()
	require.Error(t, err)

	newMsgBurnNFT = NewMsgBurnNFT(address, id, denom, []sdk.Int{})
	err = newMsgBurnNFT.ValidateBasic()
	require.NoError(t, err)
}

func TestMsgBurnNFTGetSignBytesMethod(t *testing.T) {
	newMsgBurnNFT := NewMsgBurnNFT(address, id, denom, []sdk.Int{})
	sortedBytes := newMsgBurnNFT.GetSignBytes()
	require.Equal(t, string(sortedBytes), fmt.Sprintf(`{"type":"cosmos-sdk/MsgBurnNFT","value":{"denom":"%s","id":"%s","sender":"%s"}}`,
		denom, id, address.String(),
	))
}

func TestMsgBurnNFTGetSignersMethod(t *testing.T) {
	newMsgBurnNFT := NewMsgBurnNFT(address, id, denom, []sdk.Int{})
	signers := newMsgBurnNFT.GetSigners()
	require.Equal(t, 1, len(signers))
	require.Equal(t, address.String(), signers[0].String())
}*/
