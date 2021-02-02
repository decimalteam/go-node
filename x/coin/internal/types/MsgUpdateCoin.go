package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"regexp"
)

var _ sdk.Msg = &MsgUpdateCoin{}

type MsgUpdateCoin struct {
	Sender      sdk.AccAddress `json:"sender" yaml:"sender"`
	Symbol      string         `json:"symbol" yaml:"symbol"`
	LimitVolume sdk.Int        `json:"limit_volume" yaml:"limit_volume"`
	Icon        string         `json:"icon" yaml:"icon"`
}

func NewMsgUpdateCoin(sender sdk.AccAddress, symbol string, limitVolume sdk.Int, icon string) MsgUpdateCoin {
	return MsgUpdateCoin{
		Sender:      sender,
		Symbol:      symbol,
		LimitVolume: limitVolume,
		Icon:        icon,
	}
}

const UpdateCoinConst = "update_coin"

func (msg MsgUpdateCoin) Route() string { return RouterKey }
func (msg MsgUpdateCoin) Type() string  { return UpdateCoinConst }
func (msg MsgUpdateCoin) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}

func (msg MsgUpdateCoin) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg MsgUpdateCoin) ValidateBasic() error {
	// Validate coin symbol
	if match, _ := regexp.MatchString(allowedCoinSymbols, msg.Symbol); !match {
		return ErrInvalidCoinSymbol(msg.Symbol)
	}

	if msg.LimitVolume.GT(maxCoinSupply) {
		return ErrLimitVolumeBroken(msg.LimitVolume.String(), maxCoinSupply.String())
	}

	return nil
}
