package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ sdk.Msg = &MsgMultiSendCoin{}

type SendCoin struct {
	Coin     string         `json:"coin" yaml:"coin"`
	Amount   sdk.Int        `json:"amount" yaml:"amount"`
	Receiver sdk.AccAddress `json:"receiver" yaml:"receiver"`
}

type MsgMultiSendCoin struct {
	Sender sdk.AccAddress `json:"sender" yaml:"sender"`
	Coins  []SendCoin     `json:"send_coin"`
}

func NewMsgMultiSendCoin(sender sdk.AccAddress, coins []SendCoin) MsgMultiSendCoin {
	return MsgMultiSendCoin{
		Sender: sender,
		Coins:  coins,
	}
}

const MultiSendCoinConst = "MultiSendCoin"

func (msg MsgMultiSendCoin) Route() string { return RouterKey }
func (msg MsgMultiSendCoin) Type() string  { return MultiSendCoinConst }
func (msg MsgMultiSendCoin) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}

func (msg MsgMultiSendCoin) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg MsgMultiSendCoin) ValidateBasic() sdk.Error {
	for i := range msg.Coins {
		err := ValidateSendCoin(MsgSendCoin{
			Sender:   msg.Sender,
			Coin:     msg.Coins[i].Coin,
			Amount:   msg.Coins[i].Amount,
			Receiver: msg.Coins[i].Receiver,
		})

		if err != nil {
			return err
		}
	}
	return nil
}
