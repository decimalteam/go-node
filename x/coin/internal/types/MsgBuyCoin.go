package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	decsdk "bitbucket.org/decimalteam/go-node/utils/types"
)

var _ sdk.Msg = &MsgBuyCoin{}

type MsgBuyCoin struct {
	Buyer        decsdk.AccAddress `json:"buyer" yaml:"buyer"`
	CoinToBuy    string            `json:"coin_to_buy" yaml:"coin_to_buy"`
	CoinToSell   string            `json:"coin_to_sell" yaml:"coin_to_sell"`
	AmountToBuy  sdk.Int           `json:"amount_to_buy" yaml:"amount_to_buy"`
	AmountToSell sdk.Int           `json:"amount_to_sell" yaml:"amount_to_sell"`
}

func NewMsgBuyCoin(buyer decsdk.AccAddress, coinToBuy string, coinToSell string, amountToBuy sdk.Int, amountToSell sdk.Int) MsgBuyCoin {
	return MsgBuyCoin{
		Buyer:        buyer,
		CoinToBuy:    coinToBuy,
		AmountToBuy:  amountToBuy,
		CoinToSell:   coinToSell,
		AmountToSell: amountToSell,
	}
}

const BuyCoinConst = "BuyCoin"

func (msg MsgBuyCoin) Route() string { return RouterKey }
func (msg MsgBuyCoin) Type() string  { return BuyCoinConst }
func (msg MsgBuyCoin) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.Buyer)}
}

func (msg MsgBuyCoin) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg MsgBuyCoin) ValidateBasic() error {
	if msg.CoinToSell == msg.CoinToBuy {
		return sdkerrors.New(DefaultCodespace, SameCoins, "Cannot buy same coins")
	}
	return nil
}
