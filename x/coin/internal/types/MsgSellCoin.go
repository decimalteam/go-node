package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ sdk.Msg = &MsgSellCoin{}

type MsgSellCoin struct {
	Seller             sdk.AccAddress `json:"seller" yaml:"seller"`
	CoinToBuy          string         `json:"coin_to_buy" yaml:"coin_to_buy"`
	CoinToSell         string         `json:"coin_to_sell" yaml:"coin_to_sell"`
	AmountToSell       sdk.Int        `json:"amount_to_sell" yaml:"amount_to_sell"`
	MinimumAmountToBuy sdk.Int        `json:"minimum_amount_to_buy" yaml:"minimum_amount_to_buy"`
}

func NewMsgSellCoin(seller sdk.AccAddress, coinToBuy string, coinToSell string, amountToSell sdk.Int, minAmountToBuy sdk.Int) MsgSellCoin {
	return MsgSellCoin{
		Seller:             seller,
		CoinToBuy:          coinToBuy,
		AmountToSell:       amountToSell,
		CoinToSell:         coinToSell,
		MinimumAmountToBuy: minAmountToBuy,
	}
}

const SellCoinConst = "SellCoin"

func (msg MsgSellCoin) Route() string { return RouterKey }
func (msg MsgSellCoin) Type() string  { return SellCoinConst }
func (msg MsgSellCoin) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Seller}
}

func (msg MsgSellCoin) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg MsgSellCoin) ValidateBasic() sdk.Error {
	if msg.CoinToSell == msg.CoinToBuy {
		return sdk.NewError(DefaultCodespace, SameCoins, "Cannot sell same coins")
	}
	return nil
}
