package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	decsdk "bitbucket.org/decimalteam/go-node/utils/types"
)

var _ sdk.Msg = &MsgSellAllCoin{}

type MsgSellAllCoin struct {
	Seller      decsdk.AccAddress `json:"seller" yaml:"seller"`
	CoinToBuy   string            `json:"coin_to_buy" yaml:"coin_to_buy"`
	CoinToSell  string            `json:"coin_to_sell" yaml:"coin_to_sell"`
	AmountToBuy sdk.Int           `json:"amount_to_buy" yaml:"amount_to_buy"`
}

func NewMsgSellAllCoin(seller decsdk.AccAddress, coinToBuy string, coinToSell string, amountToBuy sdk.Int) MsgSellAllCoin {
	return MsgSellAllCoin{
		Seller:      seller,
		CoinToBuy:   coinToBuy,
		CoinToSell:  coinToSell,
		AmountToBuy: amountToBuy,
	}
}

const SellAllCoinConst = "SellAllCoin"

func (msg MsgSellAllCoin) Route() string { return RouterKey }
func (msg MsgSellAllCoin) Type() string  { return SellAllCoinConst }
func (msg MsgSellAllCoin) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.Seller)}
}

func (msg MsgSellAllCoin) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg MsgSellAllCoin) ValidateBasic() error {
	if msg.CoinToSell == msg.CoinToBuy {
		return sdkerrors.New(DefaultCodespace, SameCoins, "Cannot sell same coins")
	}
	return nil
}
