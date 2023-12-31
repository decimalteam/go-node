package types

import (
	"regexp"
	"strings"

	"bitbucket.org/decimalteam/go-node/utils/helpers"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

/* --------------------------------------------------------------------------- */
// MsgMintNFT
/* --------------------------------------------------------------------------- */

type MsgMintNFT struct {
	Sender    sdk.AccAddress `json:"sender"`
	Recipient sdk.AccAddress `json:"recipient"`
	ID        string         `json:"id"`
	Denom     string         `json:"denom"`
	Quantity  sdk.Int        `json:"quantity"`
	TokenURI  string         `json:"token_uri"`
	Reserve   sdk.Int        `json:"reserve"`
	AllowMint bool           `json:"allow_mint"`
}

// NewMsgMintNFT is a constructor function for MsgMintNFT
func NewMsgMintNFT(sender, recipient sdk.AccAddress, id, denom, tokenURI string, quantity, reserve sdk.Int, allowMint bool) MsgMintNFT {
	return MsgMintNFT{
		Sender:    sender,
		Recipient: recipient,
		ID:        strings.TrimSpace(id),
		Denom:     strings.TrimSpace(denom),
		TokenURI:  strings.TrimSpace(tokenURI),
		Quantity:  quantity,
		Reserve:   reserve,
		AllowMint: allowMint,
	}
}

const regName = "^[a-zA-Z0-9_-]{1,255}$"

var MinReserve = sdk.NewInt(100)

var NewMinReserve = helpers.BipToPip(sdk.NewInt(100))
var NewMinReserve2 = helpers.BipToPip(sdk.NewInt(1))

// Route Implements Msg
func (msg MsgMintNFT) Route() string { return RouterKey }

// Type Implements Msg
func (msg MsgMintNFT) Type() string { return "mint_nft" }

// ValidateBasic Implements Msg.
func (msg MsgMintNFT) ValidateBasic() error {
	if strings.TrimSpace(msg.Denom) == "" {
		return ErrInvalidDenom(msg.Denom)
	}
	if strings.TrimSpace(msg.ID) == "" {
		return ErrInvalidNFT(msg.ID)
	}
	if msg.Sender.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "invalid sender address")
	}
	if msg.Recipient.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "invalid recipient address")
	}
	if !msg.Quantity.IsPositive() {
		return ErrInvalidQuantity(msg.Quantity.String())
	}

	if !msg.Reserve.IsPositive() || msg.Reserve.LT(MinReserve) {
		return ErrInvalidReserve(msg.Reserve.String())
	}
	if match, _ := regexp.MatchString(regName, msg.Denom); !match {
		return ErrInvalidDenom(msg.Denom)
	}
	if match, _ := regexp.MatchString(regName, msg.ID); !match {
		return ErrInvalidTokenID(msg.ID)
	}

	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgMintNFT) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// GetSigners Implements Msg.
func (msg MsgMintNFT) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}

/* --------------------------------------------------------------------------- */
// MsgBurnNFT
/* --------------------------------------------------------------------------- */

type MsgBurnNFT struct {
	Sender      sdk.AccAddress `json:"sender"`
	ID          string         `json:"id"`
	Denom       string         `json:"denom"`
	SubTokenIDs []int64        `json:"sub_token_ids"`
}

// NewMsgBurnNFT is a constructor function for MsgBurnNFT
func NewMsgBurnNFT(sender sdk.AccAddress, id string, denom string, subTokenIDs []int64) MsgBurnNFT {
	return MsgBurnNFT{
		Sender:      sender,
		ID:          strings.TrimSpace(id),
		Denom:       strings.TrimSpace(denom),
		SubTokenIDs: subTokenIDs,
	}
}

// Route Implements Msg
func (msg MsgBurnNFT) Route() string { return RouterKey }

// Type Implements Msg
func (msg MsgBurnNFT) Type() string { return "burn_nft" }

// ValidateBasic Implements Msg.
func (msg MsgBurnNFT) ValidateBasic() error {
	if strings.TrimSpace(msg.Denom) == "" {
		return ErrInvalidDenom(msg.Denom)
	}
	if strings.TrimSpace(msg.ID) == "" {
		return ErrInvalidNFT(msg.ID)
	}
	if msg.Sender.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "invalid sender address")
	}
	if !CheckUnique(msg.SubTokenIDs) {
		return ErrNotUniqueSubTokenIDs()
	}

	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgBurnNFT) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// GetSigners Implements Msg.
func (msg MsgBurnNFT) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}

/* --------------------------------------------------------------------------- */
// MsgUpdateReservNFT
/* --------------------------------------------------------------------------- */
type MsgUpdateReserveNFT struct {
	Sender        sdk.AccAddress `json:"sender"`
	ID            string         `json:"id"`
	Denom         string         `json:"denom"`
	SubTokenIDs   []int64        `json:"sub_token_ids"`
	NewReserveNFT sdk.Int        `json:"reserve"`
}

// NewUpdateReservNFT is a constructor function for MsgUpdateReservNFT
func NewMsgUpdateReserveNFT(sender sdk.AccAddress, id string, denom string, subTokenIDs []int64, newReserveNFT sdk.Int) MsgUpdateReserveNFT {
	return MsgUpdateReserveNFT{
		Sender:        sender,
		ID:            strings.TrimSpace(id),
		Denom:         strings.TrimSpace(denom),
		SubTokenIDs:   subTokenIDs,
		NewReserveNFT: newReserveNFT,
	}
}

// Route Implements Msg
func (msg MsgUpdateReserveNFT) Route() string { return RouterKey }

// Type Implements Msg
func (msg MsgUpdateReserveNFT) Type() string { return "update_nft_reserve" }

// ValidateBasic Implements Msg.
func (msg MsgUpdateReserveNFT) ValidateBasic() error {
	if strings.TrimSpace(msg.Denom) == "" {

		return ErrInvalidDenom(msg.Denom)
	}
	if strings.TrimSpace(msg.ID) == "" {
		return ErrInvalidNFT(msg.ID)
	}
	if msg.Sender.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "invalid sender address")
	}
	if !CheckUnique(msg.SubTokenIDs) {
		return ErrNotUniqueSubTokenIDs()
	}

	if msg.NewReserveNFT.IsZero() {
		return ErrInvalidReserve("Reserv can not be equal to zero")
	}

	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgUpdateReserveNFT) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// GetSigners Implements Msg.
func (msg MsgUpdateReserveNFT) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}

/* --------------------------------------------------------------------------- */
// MsgTransferNFT
/* --------------------------------------------------------------------------- */

type MsgTransferNFT struct {
	Sender      sdk.AccAddress `json:"sender"`
	Recipient   sdk.AccAddress `json:"recipient"`
	ID          string         `json:"id"`
	Denom       string         `json:"denom"`
	SubTokenIDs []int64        `json:"sub_token_ids"`
}

// NewMsgTransferNFT is a constructor function for MsgSetName
func NewMsgTransferNFT(sender, recipient sdk.AccAddress, denom, id string, subTokenIDs []int64) MsgTransferNFT {
	return MsgTransferNFT{
		Sender:      sender,
		Recipient:   recipient,
		Denom:       strings.TrimSpace(denom),
		ID:          strings.TrimSpace(id),
		SubTokenIDs: subTokenIDs,
	}
}

// Route Implements Msg
func (msg MsgTransferNFT) Route() string { return RouterKey }

// Type Implements Msg
func (msg MsgTransferNFT) Type() string { return "transfer_nft" }

// ValidateBasic Implements Msg.
func (msg MsgTransferNFT) ValidateBasic() error {
	if strings.TrimSpace(msg.Denom) == "" {
		return ErrInvalidCollection(msg.Denom)
	}
	if msg.Sender.Empty() {
		return ErrInvalidRecipientAddress(msg.Sender.String())
	}
	if msg.Recipient.Empty() {
		return ErrInvalidRecipientAddress(msg.Recipient.String())
	}
	if msg.Sender.Equals(msg.Recipient) {
		return ErrForbiddenToTransferToYourself()
	}
	if strings.TrimSpace(msg.ID) == "" {
		return ErrInvalidCollection(msg.ID)
	}
	if !CheckUnique(msg.SubTokenIDs) {
		return ErrNotUniqueSubTokenIDs()
	}

	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgTransferNFT) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// GetSigners Implements Msg.
func (msg MsgTransferNFT) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}

/* --------------------------------------------------------------------------- */
// MsgEditNFTMetadata
/* --------------------------------------------------------------------------- */

type MsgEditNFTMetadata struct {
	Sender   sdk.AccAddress `json:"sender"`
	ID       string         `json:"id"`
	Denom    string         `json:"denom"`
	TokenURI string         `json:"token_uri"`
}

// NewMsgEditNFTMetadata is a constructor function for MsgSetName
func NewMsgEditNFTMetadata(sender sdk.AccAddress, id,
	denom, tokenURI string,
) MsgEditNFTMetadata {
	return MsgEditNFTMetadata{
		Sender:   sender,
		ID:       strings.TrimSpace(id),
		Denom:    strings.TrimSpace(denom),
		TokenURI: strings.TrimSpace(tokenURI),
	}
}

// Route Implements Msg
func (msg MsgEditNFTMetadata) Route() string { return RouterKey }

// Type Implements Msg
func (msg MsgEditNFTMetadata) Type() string { return "edit_nft_metadata" }

// ValidateBasic Implements Msg.
func (msg MsgEditNFTMetadata) ValidateBasic() error {
	if msg.Sender.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "invalid sender address")
	}
	if strings.TrimSpace(msg.Denom) == "" {
		return ErrInvalidDenom(msg.Denom)
	}
	if strings.TrimSpace(msg.ID) == "" {
		return ErrInvalidNFT(msg.ID)
	}
	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgEditNFTMetadata) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// GetSigners Implements Msg.
func (msg MsgEditNFTMetadata) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}

/* --------------------------------------------------------------------------- */
func CheckUnique(arr []int64) bool {
	for i, el := range arr {
		for j, el2 := range arr {
			if i != j && el == el2 {
				return false
			}
		}
	}
	return true
}
