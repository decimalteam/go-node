package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ sdk.Msg = &MsgSendCoin{}

type MsgSendCoin struct {
	Sender   sdk.AccAddress `json:"sender" yaml:"sender"`
	Coin     string         `json:"coin" yaml:"coin"`
	Amount   sdk.Int        `json:"amount" yaml:"amount"`
	Receiver sdk.AccAddress `json:"receiver" yaml:"receiver"`
}

func NewMsgSendCoin(sender sdk.AccAddress, coin string, amount sdk.Int, receiver sdk.AccAddress) MsgSendCoin {
	return MsgSendCoin{
		Sender:   sender,
		Coin:     coin,
		Amount:   amount,
		Receiver: receiver,
	}
}

const SendCoinConst = "SendCoin"

func (msg MsgSendCoin) Route() string { return RouterKey }
func (msg MsgSendCoin) Type() string  { return SendCoinConst }
func (msg MsgSendCoin) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}

func (msg MsgSendCoin) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg MsgSendCoin) ValidateBasic() sdk.Error {
	return ValidateSendCoin(msg)
}

func ValidateSendCoin(msg MsgSendCoin) sdk.Error {
	if msg.Amount.LTE(sdk.NewInt(0)) {
		return sdk.NewError(DefaultCodespace, InvalidAmount, "Amount should be greater than 0")
	}
	return nil
}
