package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto"
)

var _ sdk.Msg = &MsgDeclareCandidate{}

type MsgDeclareCandidate struct {
	Commission    Commission     `json:"commission" yaml:"commission"`
	ValidatorAddr sdk.AccAddress `json:"validator_addr" yaml:"validator_addr"`
	PubKey        crypto.PubKey  `json:"pub_key" yaml:"pub_key"`
	Stake         sdk.Coin       `json:"value" yaml:"value"`
	Description   Description    `json:"description"`
}

func NewMsgDeclareCandidate(validatorAddr sdk.AccAddress, pubKey crypto.PubKey, commission Commission, stake sdk.Coin, description Description) MsgDeclareCandidate {
	return MsgDeclareCandidate{
		Commission:    commission,
		ValidatorAddr: validatorAddr,
		PubKey:        pubKey,
		Stake:         stake,
		Description:   description,
	}
}

const DeclareCandidateConst = "declare_candidate"

func (msg MsgDeclareCandidate) Route() string { return RouterKey }
func (msg MsgDeclareCandidate) Type() string  { return DeclareCandidateConst }
func (msg MsgDeclareCandidate) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.ValidatorAddr}
}

func (msg MsgDeclareCandidate) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg MsgDeclareCandidate) ValidateBasic() sdk.Error {
	if msg.ValidatorAddr.Empty() {
		return ErrEmptyValidatorAddr(DefaultCodespace)
	}
	return nil
}
