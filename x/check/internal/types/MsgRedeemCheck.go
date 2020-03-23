package types

import sdk "github.com/cosmos/cosmos-sdk/types"

var _ sdk.Msg = &MsgRedeemCheck{}

type MsgRedeemCheck struct {
	Check    string         `json:"check" yaml:"check"`
	Receiver sdk.AccAddress `json:"receiver" yaml:"receiver"`
	Proof    string         `json:"proof" yaml:"proof"`
}

func NewMsgRedeemCheck(check string) MsgRedeemCheck {
	return MsgRedeemCheck{
		Check: check,
	}
}

const RedeemCheckConst = "RedeemCheck"

// nolint
func (msg MsgRedeemCheck) Route() string { return RouterKey }
func (msg MsgRedeemCheck) Type() string  { return RedeemCheckConst }
func (msg MsgRedeemCheck) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Receiver}
}

func (msg MsgRedeemCheck) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg MsgRedeemCheck) ValidateBasic() sdk.Error {

	return nil
}
