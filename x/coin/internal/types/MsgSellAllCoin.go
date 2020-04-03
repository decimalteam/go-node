package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ sdk.Msg = &MsgSellAllCoin{}

type MsgSellAllCoin struct {
	Seller      sdk.AccAddress `json:"seller" yaml:"seller"`
	CoinToBuy   string         `json:"coin_to_buy" yaml:"coin_to_buy"`
	CoinToSell  string         `json:"coin_to_sell" yaml:"coin_to_sell"`
	AmountToBuy sdk.Int        `json:"amount_to_buy" yaml:"amount_to_buy"`
}

func NewMsgSellAllCoin(seller sdk.AccAddress, coinToBuy string, coinToSell string) MsgSellAllCoin {
	return MsgSellAllCoin{
		Seller:     seller,
		CoinToBuy:  coinToBuy,
		CoinToSell: coinToSell,
	}
}

const SellAllCoinConst = "SellAllCoin"

func (msg MsgSellAllCoin) Route() string { return RouterKey }
func (msg MsgSellAllCoin) Type() string  { return SellAllCoinConst }
func (msg MsgSellAllCoin) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Seller}
}

func (msg MsgSellAllCoin) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg MsgSellAllCoin) ValidateBasic() sdk.Error {
	if msg.CoinToSell == msg.CoinToBuy {
		return sdk.NewError(DefaultCodespace, SameCoins, "Cannot sell same coins")
	}
	return nil
}
