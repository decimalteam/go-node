package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgBuyCoin{}

type MsgBuyCoin struct {
	Sender        sdk.AccAddress `json:"sender" yaml:"sender"`
	CoinToBuy     sdk.Coin       `json:"coin_to_buy" yaml:"coin_to_buy"`
	MaxCoinToSell sdk.Coin       `json:"max_coin_to_sell" yaml:"max_coin_to_sell"`
}

func NewMsgBuyCoin(sender sdk.AccAddress, coinToBuy sdk.Coin, maxCoinToSell sdk.Coin) MsgBuyCoin {
	return MsgBuyCoin{
		Sender:        sender,
		CoinToBuy:     coinToBuy,
		MaxCoinToSell: maxCoinToSell,
	}
}

const BuyCoinConst = "buy_coin"

func (msg MsgBuyCoin) Route() string { return RouterKey }
func (msg MsgBuyCoin) Type() string  { return BuyCoinConst }
func (msg MsgBuyCoin) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}

func (msg MsgBuyCoin) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg MsgBuyCoin) ValidateBasic() error {
	if msg.CoinToBuy.Denom == msg.MaxCoinToSell.Denom {
		return sdkerrors.New(DefaultCodespace, SameCoins, "Cannot buy same coins")
	}
	return nil
}
