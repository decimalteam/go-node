package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto"
)

var _ sdk.Msg = &MsgCreateValidator{}

type MsgCreateValidator struct {
	Commission    Commission     `json:"commission" yaml:"commission"`
	DelegatorAddr sdk.AccAddress `json:"delegator_addr" yaml:"delegator_addr"`
	ValidatorAddr sdk.ValAddress `json:"validator_addr" yaml:"validator_addr"`
	PubKey        crypto.PubKey  `json:"pub_key" yaml:"pub_key"`
	Stake         sdk.Coin       `json:"value" yaml:"value"`
}

func NewMsgCreateValidator(validatorAddr sdk.ValAddress, pubKey crypto.PubKey, commission Commission, stake sdk.Coin) MsgCreateValidator {
	return MsgCreateValidator{
		Commission:    commission,
		DelegatorAddr: sdk.AccAddress(validatorAddr),
		ValidatorAddr: validatorAddr,
		PubKey:        pubKey,
		Stake:         stake,
	}
}

const CreateValidatorConst = "CreateValidator"

// nolint
func (msg MsgCreateValidator) Route() string { return RouterKey }
func (msg MsgCreateValidator) Type() string  { return CreateValidatorConst }
func (msg MsgCreateValidator) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.ValidatorAddr)}
}

func (msg MsgCreateValidator) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg MsgCreateValidator) ValidateBasic() sdk.Error {
	if msg.ValidatorAddr.Empty() {
		return ErrEmptyValidatorAddr(DefaultCodespace)
	}
	return nil
}
