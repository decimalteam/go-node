package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// Swap message types and routes
const (
	TypeMsgBurn  = "burn"
	TypeMsgClaim = "claim"
)

type MsgBurn struct {
	From              sdk.AccAddress `json:"from"`
	Recipient         string         `json:"recipient"`
	Amount            sdk.Int        `json:"amount"`
	TokenName         string         `json:"token_name"`
	TokenSymbol       string         `json:"token_symbol"`
	TransactionNumber string         `json:"transaction_number"`
	DestChain         int            `json:"dest_chain"`
}

func NewMsgBurn(from sdk.AccAddress, recipient string, amount sdk.Int, tokenName, tokenSymbol,
	transactionNumber string, destChain int) MsgBurn {
	return MsgBurn{
		From:              from,
		Recipient:         recipient,
		Amount:            amount,
		TokenName:         tokenName,
		TokenSymbol:       tokenSymbol,
		TransactionNumber: transactionNumber,
		DestChain:         destChain,
	}
}

func (msg MsgBurn) Route() string { return RouterKey }

func (msg MsgBurn) Type() string { return TypeMsgBurn }

func (msg MsgBurn) ValidateBasic() error {
	if msg.From.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.From.String())
	}

	if !msg.Amount.IsPositive() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, "amount must be positive")
	}

	return nil
}

func (msg MsgBurn) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg MsgBurn) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.From}
}

type MsgClaim struct {
	Sender            sdk.AccAddress `json:"sender"`
	From              string         `json:"from"`
	Recipient         sdk.AccAddress `json:"recipient"`
	Amount            sdk.Int        `json:"amount"`
	TokenName         string         `json:"token_name"`
	TokenSymbol       string         `json:"token_symbol"`
	TransactionNumber string         `json:"transaction_number"`
	DestChain         int            `json:"dest_chain"`
	V                 uint8          `json:"v"`
	R                 [32]byte       `json:"r"`
	S                 [32]byte       `json:"s"`
}

func NewMsgClaim(sender, recipient sdk.AccAddress, from string, amount sdk.Int, tokenName, tokenSymbol,
	transactionNumber string, destChain int) MsgClaim {
	return MsgClaim{
		Sender:            sender,
		From:              from,
		Recipient:         recipient,
		Amount:            amount,
		TokenName:         tokenName,
		TokenSymbol:       tokenSymbol,
		TransactionNumber: transactionNumber,
		DestChain:         destChain,
	}
}

func (msg MsgClaim) Route() string { return RouterKey }

func (msg MsgClaim) Type() string { return TypeMsgClaim }

func (msg MsgClaim) ValidateBasic() error {
	if msg.Sender.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Sender.String())
	}

	if !msg.Amount.IsPositive() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, "amount must be positive")
	}

	return nil
}

func (msg MsgClaim) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg MsgClaim) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}
