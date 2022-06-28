package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ sdk.Msg = &MsgBurnCoin{}

type MsgBurnCoin struct {
	Sender sdk.AccAddress `json:"sender" yaml:"sender"`
	Coin   sdk.Coin       `json:"coin" yaml:"coin"`
}

func NewMsgBurnCoin(sender sdk.AccAddress, coin sdk.Coin) MsgBurnCoin {
	return MsgBurnCoin{
		Sender: sender,
		Coin:   coin,
	}
}

const BurnCoinConst = "burn_coin"

func (msg MsgBurnCoin) Route() string { return RouterKey }
func (msg MsgBurnCoin) Type() string  { return BurnCoinConst }
func (msg MsgBurnCoin) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}

func (msg MsgBurnCoin) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg MsgBurnCoin) ValidateBasic() error {
	return ValidateBurn(msg)
}

func ValidateBurn(msg MsgBurnCoin) error {
	if msg.Coin.Amount.LTE(sdk.NewInt(0)) {
		return ErrInvalidAmount()
	}
	return nil
}
