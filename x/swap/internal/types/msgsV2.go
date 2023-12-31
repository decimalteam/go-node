package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// Swap message types and routes
const (
	TypeMsgSwapInitialize  = "swap_initialize"
	TypeMsgRedeemV2        = "redeem_v2"
	TypeMsgChainActivate   = "chain_activate"
	TypeMsgChainDeactivate = "chain_deactivate"
)

type MsgSwapInitialize struct {
	From              sdk.AccAddress `json:"from"`
	Recipient         string         `json:"recipient"`
	Amount            sdk.Int        `json:"amount"`
	TokenSymbol       string         `json:"token_symbol"`
	TransactionNumber string         `json:"transaction_number"`
	FromChain         int            `json:"from_chain"`
	DestChain         int            `json:"dest_chain"`
}

func NewMsgSwapInitialize(from sdk.AccAddress, recipient string, amount sdk.Int, tokenSymbol,
	transactionNumber string, fromChain, destChain int) MsgSwapInitialize {
	return MsgSwapInitialize{
		From:              from,
		Recipient:         recipient,
		Amount:            amount,
		TokenSymbol:       tokenSymbol,
		TransactionNumber: transactionNumber,
		FromChain:         fromChain,
		DestChain:         destChain,
	}
}

func (msg MsgSwapInitialize) Route() string { return RouterKey }

func (msg MsgSwapInitialize) Type() string { return TypeMsgSwapInitialize }

func (msg MsgSwapInitialize) ValidateBasic() error {
	if msg.From.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.From.String())
	}

	if !msg.Amount.IsPositive() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, "amount must be positive")
	}

	return nil
}

func (msg MsgSwapInitialize) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg MsgSwapInitialize) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.From}
}

type MsgRedeemV2 struct {
	Sender            sdk.AccAddress `json:"sender"`
	From              string         `json:"from"`
	Recipient         sdk.AccAddress `json:"recipient"`
	Amount            sdk.Int        `json:"amount"`
	TokenSymbol       string         `json:"token_symbol"`
	TransactionNumber string         `json:"transaction_number"`
	FromChain         int            `json:"from_chain"`
	DestChain         int            `json:"dest_chain"`
	V                 uint8          `json:"v"`
	R                 Hash           `json:"r"`
	S                 Hash           `json:"s"`
}

func NewMsgRedeemV2(sender, recipient sdk.AccAddress, from string, amount sdk.Int, tokenSymbol,
	transactionNumber string, fromChain, destChain int, v uint8, r, s [32]byte) MsgRedeemV2 {
	return MsgRedeemV2{
		Sender:            sender,
		From:              from,
		Recipient:         recipient,
		Amount:            amount,
		TokenSymbol:       tokenSymbol,
		TransactionNumber: transactionNumber,
		FromChain:         fromChain,
		DestChain:         destChain,
		V:                 v,
		R:                 r,
		S:                 s,
	}
}

func (msg MsgRedeemV2) Route() string { return RouterKey }

func (msg MsgRedeemV2) Type() string { return TypeMsgRedeemV2 }

func (msg MsgRedeemV2) ValidateBasic() error {
	if msg.Sender.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Sender.String())
	}

	if !msg.Amount.IsPositive() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, "amount must be positive")
	}

	return nil
}

func (msg MsgRedeemV2) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg MsgRedeemV2) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}

type MsgChainActivate struct {
	From        sdk.AccAddress `json:"from"`
	ChainNumber int            `json:"chain_number"`
	ChainName   string         `json:"chain_name"`
}

func NewMsgChainActivate(from sdk.AccAddress, chainNumber int, chainName string) MsgChainActivate {
	return MsgChainActivate{From: from, ChainNumber: chainNumber, ChainName: chainName}
}

func (msg MsgChainActivate) Route() string { return RouterKey }

func (msg MsgChainActivate) Type() string { return TypeMsgChainActivate }

func (msg MsgChainActivate) ValidateBasic() error {
	if msg.From.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.From.String())
	}

	swapServiceAddress, err := sdk.AccAddressFromBech32(ChainActivatorAddress)
	if err != nil {
		return err
	}
	if !msg.From.Equals(swapServiceAddress) {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.From.String())
	}

	if msg.ChainNumber <= 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, "chain number must be positive")
	}

	return nil
}

func (msg MsgChainActivate) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg MsgChainActivate) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.From}
}

type MsgChainDeactivate struct {
	From        sdk.AccAddress `json:"from"`
	ChainNumber int            `json:"chain_number"`
}

func NewMsgChainDeactivate(from sdk.AccAddress, chainNumber int) MsgChainDeactivate {
	return MsgChainDeactivate{From: from, ChainNumber: chainNumber}
}

func (msg MsgChainDeactivate) Route() string { return RouterKey }

func (msg MsgChainDeactivate) Type() string { return TypeMsgChainDeactivate }

func (msg MsgChainDeactivate) ValidateBasic() error {
	if msg.From.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.From.String())
	}

	swapServiceAddress, err := sdk.AccAddressFromBech32(ChainActivatorAddress)
	if err != nil {
		return err
	}
	if !msg.From.Equals(swapServiceAddress) {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.From.String())
	}

	if msg.ChainNumber <= 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, "chain number must be positive")
	}

	return nil
}

func (msg MsgChainDeactivate) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg MsgChainDeactivate) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.From}
}
