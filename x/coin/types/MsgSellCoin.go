package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ sdk.Msg = &MsgSellCoin{}

//type MsgSellCoin struct {
//	Sender       sdk.AccAddress `json:"sender" yaml:"sender"`
//	CoinToSell   sdk.Coin       `json:"coin_to_sell" yaml:"coin_to_sell"`
//	MinCoinToBuy sdk.Coin       `json:"min_coin_to_buy" yaml:"min_coin_to_buy"`
//}

func NewMsgSellCoin(sender sdk.AccAddress, coinToSell sdk.Coin, minCoinToBuy sdk.Coin) MsgSellCoin {
	return MsgSellCoin{
		Sender:       sender.String(),
		CoinToSell:   coinToSell,
		MinCoinToBuy: minCoinToBuy,
	}
}

const SellCoinConst = "sell_coin"

func (msg MsgSellCoin) Route() string { return RouterKey }
func (msg MsgSellCoin) Type() string  { return SellCoinConst }
func (msg MsgSellCoin) GetSigners() []sdk.AccAddress {
	accAddr, _ := sdk.AccAddressFromBech32(msg.Sender)

	return []sdk.AccAddress{accAddr}
}

func (msg MsgSellCoin) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg MsgSellCoin) ValidateBasic() error {
	if msg.CoinToSell.Denom == msg.MinCoinToBuy.Denom {
		return ErrSameCoin()
	}
	return nil
}
