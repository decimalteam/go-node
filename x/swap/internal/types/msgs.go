package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// Swap message types and routes
const (
	TypeMsgHTLT   = "htlt"
	TypeMsgRedeem = "redeem"
	TypeMsgRefund = "refund"
)

var _ sdk.Msg = MsgHTLT{}
var _ sdk.Msg = MsgRedeem{}
var _ sdk.Msg = MsgRefund{}

type TransferType int

const (
	TransferTypeOut = 1
	TransferTypeIn  = 2
)

func TransferTypeFromString(transferType string) (TransferType, error) {
	switch transferType {
	case "out":
		return TransferTypeOut, nil
	case "in":
		return TransferTypeIn, nil
	default:
		return TransferType(0xff), fmt.Errorf("'%s' is not a valid transfer type", transferType)
	}
}

func (t TransferType) String() string {
	switch t {
	case TransferTypeOut:
		return "out"
	case TransferTypeIn:
		return "in"
	default:
		return fmt.Sprintf("'%d' is not a valid transfer type", t)
	}
}

type MsgHTLT struct {
	TransferType TransferType   `json:"transfer_type"`
	From         sdk.AccAddress `json:"from"`
	Recipient    string         `json:"recipient"`
	HashedSecret Hash           `json:"hashed_secret"`
	Amount       sdk.Coins      `json:"amount"`
}

func NewMsgHTLT(transferType TransferType, from sdk.AccAddress, recipient string, hashedSecret Hash, amount sdk.Coins) MsgHTLT {
	return MsgHTLT{TransferType: transferType, From: from, Recipient: recipient, HashedSecret: hashedSecret, Amount: amount}
}

func (msg MsgHTLT) Route() string { return RouterKey }

func (msg MsgHTLT) Type() string { return TypeMsgHTLT }

func (msg MsgHTLT) ValidateBasic() error {
	if msg.From.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.From.String())
	}

	if msg.Amount.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, "out amount is empty")
	}

	return nil
}

func (msg MsgHTLT) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg MsgHTLT) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.From}
}

type MsgRedeem struct {
	From   sdk.AccAddress `json:"from"`
	Secret Hash           `json:"secret"`
}

func NewMsgRedeem(from sdk.AccAddress, secret [32]byte) MsgRedeem {
	return MsgRedeem{From: from, Secret: secret}
}

func (msg MsgRedeem) Route() string { return RouterKey }

func (msg MsgRedeem) Type() string { return TypeMsgRedeem }

func (msg MsgRedeem) ValidateBasic() error {
	if msg.From.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.From.String())
	}
	return nil
}

func (msg MsgRedeem) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg MsgRedeem) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.From}
}

type MsgRefund struct {
	From         sdk.AccAddress `json:"from"`
	HashedSecret Hash           `json:"hashed_secret"`
}

func NewMsgRefund(from sdk.AccAddress, hash [32]byte) MsgRefund {
	return MsgRefund{From: from, HashedSecret: hash}
}

func (msg MsgRefund) Route() string { return RouterKey }

func (msg MsgRefund) Type() string { return TypeMsgRefund }

func (msg MsgRefund) ValidateBasic() error {
	if msg.From.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.From.String())
	}
	return nil
}

func (msg MsgRefund) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg MsgRefund) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.From}
}
