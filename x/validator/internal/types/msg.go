package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto"
)

var _ sdk.Msg = &MsgDeclareCandidate{}

type MsgDeclareCandidate struct {
	Commission    sdk.Dec        `json:"commission" yaml:"commission"`
	ValidatorAddr sdk.ValAddress `json:"validator_addr" yaml:"validator_addr"`
	RewardAddr    sdk.AccAddress `json:"reward_addr" yaml:"reward_addr"`
	PubKey        crypto.PubKey  `json:"pub_key" yaml:"pub_key"`
	Stake         sdk.Coin       `json:"stake" yaml:"stake"`
	Description   Description    `json:"description"`
}

func NewMsgDeclareCandidate(validatorAddr sdk.ValAddress, pubKey crypto.PubKey, commission sdk.Dec, stake sdk.Coin, description Description, rewardAddress sdk.AccAddress) MsgDeclareCandidate {
	return MsgDeclareCandidate{
		Commission:    commission,
		ValidatorAddr: validatorAddr,
		RewardAddr:    rewardAddress,
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

// quick validity check
func (msg MsgDeclareCandidate) ValidateBasic() error {
	// note that unmarshaling from bech32 ensures either empty or valid
	if msg.ValidatorAddr.Empty() {
		return ErrEmptyValidatorAddr()
	}
	if msg.Stake.Amount.LTE(sdk.ZeroInt()) {
		return ErrBadDelegationAmount()
	}
	if msg.Commission.LT(sdk.ZeroDec()) {
		return ErrCommissionNegative()
	}
	if msg.Commission.GT(sdk.OneDec()) {
		return ErrCommissionHuge()
	}

	return nil
}

// -----------------------------------------------------------------------------------------

type MsgDelegate struct {
	DelegatorAddress sdk.AccAddress `json:"delegator_address"`
	ValidatorAddress sdk.ValAddress `json:"validator_address"`
	Coin             sdk.Coin       `json:"coin"`
}

func NewMsgDelegate(validatorAddr sdk.ValAddress, delegatorAddr sdk.AccAddress, coin sdk.Coin) MsgDelegate {
	return MsgDelegate{
		DelegatorAddress: delegatorAddr,
		ValidatorAddress: validatorAddr,
		Coin:             coin,
	}
}

const DelegateConst = "delegate"

func (msg MsgDelegate) Route() string { return RouterKey }
func (msg MsgDelegate) Type() string  { return DelegateConst }
func (msg MsgDelegate) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.DelegatorAddress}
}

func (msg MsgDelegate) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg MsgDelegate) ValidateBasic() error {
	if msg.ValidatorAddress.Empty() {
		return ErrEmptyValidatorAddr()
	}
	if msg.DelegatorAddress.Empty() {
		return ErrEmptyDelegatorAddr()
	}
	return nil
}

// -----------------------------------------------------------------------------------------

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

func (msg MsgSetOnline) ValidateBasic() error {
	if msg.ValidatorAddress.Empty() {
		return ErrEmptyValidatorAddr()
	}
	return nil
}

// -----------------------------------------------------------------------------------------

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

func (msg MsgSetOffline) ValidateBasic() error {
	if msg.ValidatorAddress.Empty() {
		return ErrEmptyValidatorAddr()
	}
	return nil
}

// -----------------------------------------------------------------------------------------

type MsgUnbond struct {
	ValidatorAddress sdk.ValAddress `json:"validator_address"`
	DelegatorAddress sdk.AccAddress `json:"delegator_address"`
	Coin             sdk.Coin       `json:"coin"`
}

func NewMsgUnbond(validatorAddr sdk.ValAddress, delegatorAddr sdk.AccAddress, coin sdk.Coin) MsgUnbond {
	return MsgUnbond{
		ValidatorAddress: validatorAddr,
		DelegatorAddress: delegatorAddr,
		Coin:             coin,
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

func (msg MsgUnbond) ValidateBasic() error {
	if msg.ValidatorAddress.Empty() {
		return ErrEmptyValidatorAddr()
	}
	if msg.DelegatorAddress.Empty() {
		return ErrEmptyDelegatorAddr()
	}
	return nil
}

// -----------------------------------------------------------------------------------------

type MsgEditCandidate struct {
	ValidatorAddress sdk.ValAddress `json:"validator_address"`
	RewardAddress    sdk.AccAddress `json:"reward_address"`
	Description      Description    `json:"description"`
}

func NewMsgEditCandidate(validatorAddress sdk.ValAddress, rewardAddress sdk.AccAddress, description Description) MsgEditCandidate {
	return MsgEditCandidate{
		ValidatorAddress: validatorAddress,
		RewardAddress:    rewardAddress,
		Description:      description,
	}
}

const EditCandidateConst = "edit_candidate"

func (msg MsgEditCandidate) Route() string { return RouterKey }
func (msg MsgEditCandidate) Type() string  { return EditCandidateConst }
func (msg MsgEditCandidate) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.ValidatorAddress)}
}

func (msg MsgEditCandidate) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg MsgEditCandidate) ValidateBasic() error {
	if msg.ValidatorAddress.Empty() {
		return ErrEmptyValidatorAddr()
	}
	return nil
}
