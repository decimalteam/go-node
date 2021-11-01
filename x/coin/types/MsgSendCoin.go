package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ sdk.Msg = &MsgSendCoin{}

//type MsgSendCoin struct {
//	Sender   sdk.AccAddress `json:"sender" yaml:"sender"`
//	Coin     sdk.Coin       `json:"coin" yaml:"coin"`
//	Receiver sdk.AccAddress `json:"receiver" yaml:"receiver"`
//}

func NewMsgSendCoin(sender sdk.AccAddress, coin sdk.Coin, receiver sdk.AccAddress) MsgSendCoin {
	return MsgSendCoin{
		Sender:   sender.String(),
		Coin:     coin,
		Receiver: receiver.String(),
	}
}

const SendCoinConst = "send_coin"

func (msg *MsgSendCoin) Route() string { return RouterKey }
func (msg *MsgSendCoin) Type() string  { return SendCoinConst }
func (msg *MsgSendCoin) GetSigners() []sdk.AccAddress {
	accAddr, _ := sdk.AccAddressFromBech32(msg.Sender)

	return []sdk.AccAddress{accAddr}
}

func (msg *MsgSendCoin) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgSendCoin) ValidateBasic() error {
	return ValidateSend(*msg)
}

func ValidateSend(msg MsgSendCoin) error {
	if msg.Coin.Amount.LTE(sdk.NewInt(0)) {
		return ErrInvalidAmount()
	}
	receiver , err := sdk.AccAddressFromBech32(msg.Receiver)
	if err!= nil {
		return err
	}

	if receiver.Empty() {
		return ErrReceiverEmpty()
	}
	return nil
}
