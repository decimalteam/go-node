package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	decsdk "bitbucket.org/decimalteam/go-node/utils/types"
)

var _ sdk.Msg = &MsgSendCoin{}

type MsgSendCoin struct {
	Sender   decsdk.AccAddress `json:"sender" yaml:"sender"`
	Coin     string            `json:"coin" yaml:"coin"`
	Amount   sdk.Int           `json:"amount" yaml:"amount"`
	Receiver decsdk.AccAddress `json:"receiver" yaml:"receiver"`
}

func NewMsgSendCoin(sender decsdk.AccAddress, coin string, amount sdk.Int, receiver decsdk.AccAddress) MsgSendCoin {
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
	return []sdk.AccAddress{sdk.AccAddress(msg.Sender)}
}

func (msg MsgSendCoin) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg MsgSendCoin) ValidateBasic() error {
	return ValidateSendCoin(msg)
}

func ValidateSendCoin(msg MsgSendCoin) error {
	if msg.Amount.LTE(sdk.NewInt(0)) {
		return sdkerrors.New(DefaultCodespace, InvalidAmount, "Amount should be greater than 0")
	}
	return nil
}
