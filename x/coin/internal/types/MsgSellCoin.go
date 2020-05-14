package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	decsdk "bitbucket.org/decimalteam/go-node/utils/types"
)

var _ sdk.Msg = &MsgSellCoin{}

type MsgSellCoin struct {
	Seller       decsdk.AccAddress `json:"seller" yaml:"seller"`
	CoinToBuy    string            `json:"coin_to_buy" yaml:"coin_to_buy"`
	CoinToSell   string            `json:"coin_to_sell" yaml:"coin_to_sell"`
	AmountToSell sdk.Int           `json:"amount_to_sell" yaml:"amount_to_sell"`
	AmountToBuy  sdk.Int           `json:"amount_to_buy" yaml:"amount_to_buy"`
}

func NewMsgSellCoin(seller decsdk.AccAddress, coinToBuy string, coinToSell string, amountToSell sdk.Int, amountToBuy sdk.Int) MsgSellCoin {
	return MsgSellCoin{
		Seller:       seller,
		CoinToBuy:    coinToBuy,
		AmountToSell: amountToSell,
		CoinToSell:   coinToSell,
		AmountToBuy:  amountToBuy,
	}
}

const SellCoinConst = "SellCoin"

func (msg MsgSellCoin) Route() string { return RouterKey }
func (msg MsgSellCoin) Type() string  { return SellCoinConst }
func (msg MsgSellCoin) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.Seller)}
}

func (msg MsgSellCoin) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg MsgSellCoin) ValidateBasic() error {
	if msg.CoinToSell == msg.CoinToBuy {
		return sdkerrors.New(DefaultCodespace, SameCoins, "Cannot sell same coins")
	}
	return nil
}
