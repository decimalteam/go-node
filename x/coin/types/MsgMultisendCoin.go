package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ sdk.Msg = &MsgMultiSendCoin{}

type Send struct {
	Coin     sdk.Coin       `json:"coin" yaml:"coin"`
	Receiver sdk.AccAddress `json:"receiver" yaml:"receiver"`
}

type MsgMultiSendCoin struct {
	Sender sdk.AccAddress `json:"sender" yaml:"sender"`
	Sends  []Send         `json:"sends"`
}

func (msg MsgMultiSendCoin) Reset() {
	panic("implement me")
}

func (msg MsgMultiSendCoin) String() string {
	panic("implement me")
}

func (msg MsgMultiSendCoin) ProtoMessage() {
	panic("implement me")
}

func NewMsgMultiSendCoin(sender sdk.AccAddress, sends []Send) MsgMultiSendCoin {
	return MsgMultiSendCoin{
		Sender: sender,
		Sends:  sends,
	}
}

const MultiSendCoinConst = "multi_send_coin"

func (msg MsgMultiSendCoin) Route() string { return RouterKey }
func (msg MsgMultiSendCoin) Type() string  { return MultiSendCoinConst }
func (msg MsgMultiSendCoin) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}

func (msg MsgMultiSendCoin) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg MsgMultiSendCoin) ValidateBasic() error {
	for i := range msg.Sends {
		err := ValidateSend(MsgSendCoin{
			Sender:   msg.Sender,
			Coin:     msg.Sends[i].Coin,
			Receiver: msg.Sends[i].Receiver,
		})

		if err != nil {
			return err
		}
	}
	return nil
}
