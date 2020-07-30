package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ sdk.Msg = &MsgSellAllCoin{}

type MsgSellAllCoin struct {
	Sender       sdk.AccAddress `json:"sender" yaml:"sender"`
	CoinToSell   sdk.Coin       `json:"coin_to_sell" yaml:"coin_to_sell"`
	MinCoinToBuy sdk.Coin       `json:"min_coin_to_buy" yaml:"min_coin_to_buy"`
}

func NewMsgSellAllCoin(sender sdk.AccAddress, coinToSell sdk.Coin, minCoinToBuy sdk.Coin) MsgSellAllCoin {
	return MsgSellAllCoin{
		Sender:       sender,
		CoinToSell:   coinToSell,
		MinCoinToBuy: minCoinToBuy,
	}
}

const SellAllCoinConst = "sell_all_coin"

func (msg MsgSellAllCoin) Route() string { return RouterKey }
func (msg MsgSellAllCoin) Type() string  { return SellAllCoinConst }
func (msg MsgSellAllCoin) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}

func (msg MsgSellAllCoin) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg MsgSellAllCoin) ValidateBasic() error {
	if msg.CoinToSell.Denom == msg.MinCoinToBuy.Denom {
		return ErrSameCoin()
	}
	return nil
}
