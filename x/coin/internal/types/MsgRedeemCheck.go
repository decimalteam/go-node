package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ sdk.Msg = &MsgRedeemCheck{}

type MsgRedeemCheck struct {
	Receiver sdk.AccAddress `json:"receiver" yaml:"receiver"`
	Check    string         `json:"check" yaml:"check"`
	Proof    string         `json:"proof" yaml:"proof"`
}

func NewMsgRedeemCheck(receiver sdk.AccAddress, check string, proof string) MsgRedeemCheck {
	return MsgRedeemCheck{
		Receiver: receiver,
		Check:    check,
		Proof:    proof,
	}
}

const RedeemCheckConst = "RedeemCheck"

func (msg MsgRedeemCheck) Route() string { return RouterKey }
func (msg MsgRedeemCheck) Type() string  { return RedeemCheckConst }
func (msg MsgRedeemCheck) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Receiver}
}

func (msg MsgRedeemCheck) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg MsgRedeemCheck) ValidateBasic() error {
	return ValidateRedeemCheck(msg)
}

func ValidateRedeemCheck(msg MsgRedeemCheck) error {
	// TODO
	return nil
}
