package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"regexp"
	"strings"
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

const regName = "^[a-zA-Z0-9_]{1,255}$"

var minReserve = sdk.NewInt(100)

// Route Implements Msg
func (msg MsgMintNFT) Route() string { return RouterKey }

// Type Implements Msg
func (msg MsgMintNFT) Type() string { return "mint_nft" }

// ValidateBasic Implements Msg.
func (msg MsgMintNFT) ValidateBasic() error {
	if strings.TrimSpace(msg.Denom) == "" || strings.TrimSpace(msg.ID) == "" {
		return ErrInvalidNFT()
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
	if !msg.Reserve.IsPositive() || msg.Reserve.LT(minReserve) {
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
	Sender   sdk.AccAddress `json:"sender"`
	ID       string         `json:"id"`
	Denom    string         `json:"denom"`
	Quantity sdk.Int        `json:"quantity"`
}

// NewMsgBurnNFT is a constructor function for MsgBurnNFT
func NewMsgBurnNFT(sender sdk.AccAddress, id string, denom string, quantity sdk.Int) MsgBurnNFT {
	return MsgBurnNFT{
		Sender:   sender,
		ID:       strings.TrimSpace(id),
		Denom:    strings.TrimSpace(denom),
		Quantity: quantity,
	}
}

// Route Implements Msg
func (msg MsgBurnNFT) Route() string { return RouterKey }

// Type Implements Msg
func (msg MsgBurnNFT) Type() string { return "burn_nft" }

// ValidateBasic Implements Msg.
func (msg MsgBurnNFT) ValidateBasic() error {
	if strings.TrimSpace(msg.ID) == "" || strings.TrimSpace(msg.Denom) == "" {
		return ErrInvalidNFT()
	}
	if msg.Sender.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "invalid sender address")
	}
	if !msg.Quantity.IsPositive() {
		return ErrInvalidQuantity(msg.Quantity.String())
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
// MsgTransferNFT
/* --------------------------------------------------------------------------- */

type MsgTransferNFT struct {
	Sender    sdk.AccAddress `json:"sender"`
	Recipient sdk.AccAddress `json:"recipient"`
	ID        string         `json:"id"`
	Denom     string         `json:"denom"`
	Quantity  sdk.Int        `json:"quantity"`
}

// NewMsgTransferNFT is a constructor function for MsgSetName
func NewMsgTransferNFT(sender, recipient sdk.AccAddress, denom, id string, quantity sdk.Int) MsgTransferNFT {
	return MsgTransferNFT{
		Sender:    sender,
		Recipient: recipient,
		Denom:     strings.TrimSpace(denom),
		ID:        strings.TrimSpace(id),
		Quantity:  quantity,
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
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "invalid sender address")
	}
	if msg.Recipient.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "invalid recipient address")
	}
	if strings.TrimSpace(msg.ID) == "" {
		return ErrInvalidCollection(msg.ID)
	}
	if !msg.Quantity.IsPositive() {
		return ErrInvalidQuantity(msg.Quantity.String())
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
	if strings.TrimSpace(msg.ID) == "" || strings.TrimSpace(msg.Denom) == "" {
		return ErrInvalidNFT()
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
// MsgDelegateNFT
/* --------------------------------------------------------------------------- */

type MsgDelegateNFT struct {
	DelegatorAddress sdk.AccAddress `json:"delegator_address"`
	ValidatorAddress sdk.ValAddress `json:"validator_address"`
	ID               string         `json:"id"`
	Denom            string         `json:"denom"`
	Quantity         sdk.Int        `json:"quantity"`
}
