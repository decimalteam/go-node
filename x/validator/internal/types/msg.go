package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto"
)

var _ sdk.Msg = &MsgDeclareCandidate{}

type MsgDeclareCandidate struct {
	Commission    Commission     `json:"commission" yaml:"commission"`
	ValidatorAddr sdk.ValAddress `json:"validator_addr" yaml:"validator_addr"`
	PubKey        crypto.PubKey  `json:"pub_key" yaml:"pub_key"`
	Stake         sdk.Coin       `json:"value" yaml:"value"`
	Description   Description    `json:"description"`
}

func NewMsgDeclareCandidate(validatorAddr sdk.ValAddress, pubKey crypto.PubKey, commission Commission, stake sdk.Coin, description Description) MsgDeclareCandidate {
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
	return []sdk.AccAddress{sdk.AccAddress(msg.ValidatorAddr)}
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

type MsgDelegate struct {
	DelegatorAddress sdk.AccAddress `json:"delegator_address"`
	ValidatorAddress sdk.ValAddress `json:"validator_address"`
	Amount           sdk.Coin       `json:"amount"`
}

func NewMsgDelegate(validatorAddr sdk.ValAddress, delegatorAddr sdk.AccAddress, amount sdk.Coin) MsgDelegate {
	return MsgDelegate{
		DelegatorAddress: delegatorAddr,
		ValidatorAddress: validatorAddr,
		Amount:           amount,
	}
}

const DelegateConst = "delegate"

func (msg MsgDelegate) Route() string { return RouterKey }
func (msg MsgDelegate) Type() string  { return DeclareCandidateConst }
func (msg MsgDelegate) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.DelegatorAddress}
}

func (msg MsgDelegate) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg MsgDelegate) ValidateBasic() sdk.Error {
	if msg.ValidatorAddress.Empty() {
		return ErrEmptyValidatorAddr(DefaultCodespace)
	}
	if msg.DelegatorAddress.Empty() {
		return ErrNilDelegatorAddr(DefaultCodespace)
	}
	return nil
}

type MsgSetOnline struct {
	ValidatorAddress sdk.ValAddress `json:"validator_address"`
}

func NewMsgSetOnline(validatorAddr sdk.ValAddress) MsgSetOnline {
	return MsgSetOnline{ValidatorAddress: validatorAddr}
}

const SetOnlineConst = "set_online"

func (msg MsgSetOnline) Route() string { return RouterKey }
func (msg MsgSetOnline) Type() string  { return SetOnlineConst }
func (msg MsgSetOnline) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.ValidatorAddress)}
}

func (msg MsgSetOnline) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg MsgSetOnline) ValidateBasic() sdk.Error {
	if msg.ValidatorAddress.Empty() {
		return ErrEmptyValidatorAddr(DefaultCodespace)
	}
	return nil
}

type MsgSetOffline struct {
	ValidatorAddress sdk.ValAddress `json:"validator_address"`
}

func NewMsgSetOffline(validatorAddr sdk.ValAddress) MsgSetOffline {
	return MsgSetOffline{ValidatorAddress: validatorAddr}
}

const SetOfflineConst = "set_offline"

func (msg MsgSetOffline) Route() string { return RouterKey }
func (msg MsgSetOffline) Type() string  { return SetOfflineConst }
func (msg MsgSetOffline) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.ValidatorAddress)}
}

func (msg MsgSetOffline) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg MsgSetOffline) ValidateBasic() sdk.Error {
	if msg.ValidatorAddress.Empty() {
		return ErrEmptyValidatorAddr(DefaultCodespace)
	}
	return nil
}

type MsgUnbond struct {
	ValidatorAddress sdk.ValAddress `json:"validator_address"`
	DelegatorAddress sdk.AccAddress `json:"delegator_address"`
	Amount           sdk.Coin       `json:"amount"`
}

func NewMsgUnbond(validatorAddr sdk.ValAddress, delegatorAddr sdk.AccAddress, amount sdk.Coin) MsgUnbond {
	return MsgUnbond{
		ValidatorAddress: validatorAddr,
		DelegatorAddress: delegatorAddr,
		Amount:           amount,
	}
}

const UnbondConst = "unbond"

func (msg MsgUnbond) Route() string { return RouterKey }
func (msg MsgUnbond) Type() string  { return UnbondConst }
func (msg MsgUnbond) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.DelegatorAddress}
}

func (msg MsgUnbond) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg MsgUnbond) ValidateBasic() sdk.Error {
	if msg.ValidatorAddress.Empty() {
		return ErrEmptyValidatorAddr(DefaultCodespace)
	}
	if msg.DelegatorAddress.Empty() {
		return ErrNilDelegatorAddr(DefaultCodespace)
	}
	return nil
}
