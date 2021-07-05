package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ sdk.Msg = &MsgRedeemCheck{}

//type MsgRedeemCheck struct {
//	Sender sdk.AccAddress `json:"sender" yaml:"sender"`
//	Check  string         `json:"check" yaml:"check"`
//	Proof  string         `json:"proof" yaml:"proof"`
//}

func NewMsgRedeemCheck(sender sdk.AccAddress, check string, proof string) MsgRedeemCheck {
	return MsgRedeemCheck{
		Sender: sender.String(),
		Check:  check,
		Proof:  proof,
	}
}

const RedeemCheckConst = "redeem_check"

func (msg *MsgRedeemCheck) Route() string { return RouterKey }
func (msg *MsgRedeemCheck) Type() string  { return RedeemCheckConst }
func (msg *MsgRedeemCheck) GetSigners() []sdk.AccAddress {
	accAddr, _ := sdk.AccAddressFromBech32(msg.Sender)

	return []sdk.AccAddress{accAddr}
}

func (msg *MsgRedeemCheck) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgRedeemCheck) ValidateBasic() error {
	return ValidateRedeemCheck(*msg)
}

func ValidateRedeemCheck(msg MsgRedeemCheck) error {
	// TODO
	return nil
}
