package types

import (
	"github.com/tendermint/tendermint/crypto"

	sdk "github.com/cosmos/cosmos-sdk/types"

	decsdk "bitbucket.org/decimalteam/go-node/utils/types"
)

var _ sdk.Msg = &MsgDeclareCandidate{}

type MsgDeclareCandidate struct {
	Commission    sdk.Dec           `json:"commission" yaml:"commission"`
	ValidatorAddr decsdk.ValAddress `json:"validator_addr" yaml:"validator_addr"`
	RewardAddr    decsdk.AccAddress `json:"reward_addr" yaml:"reward_addr"`
	PubKey        crypto.PubKey     `json:"pub_key" yaml:"pub_key"`
	Stake         sdk.Coin          `json:"value" yaml:"value"`
	Description   Description       `json:"description"`
}

func NewMsgDeclareCandidate(validatorAddr decsdk.ValAddress, pubKey crypto.PubKey, commission sdk.Dec, stake sdk.Coin, description Description, rewardAddress decsdk.AccAddress) MsgDeclareCandidate {
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
		return ErrNilValidatorAddr(DefaultCodespace)
	}
	if msg.Stake.Amount.LTE(sdk.ZeroInt()) {
		return ErrBadDelegationAmount(DefaultCodespace)
	}
	if msg.Commission.LT(sdk.ZeroDec()) {
		return ErrCommissionNegative(DefaultCodespace)
	}
	if msg.Commission.GT(sdk.OneDec()) {
		return ErrCommissionHuge(DefaultCodespace)
	}

	return nil
}

// -----------------------------------------------------------------------------------------

type MsgDelegate struct {
	DelegatorAddress decsdk.AccAddress `json:"delegator_address"`
	ValidatorAddress decsdk.ValAddress `json:"validator_address"`
	Amount           sdk.Coin          `json:"amount"`
}

func NewMsgDelegate(validatorAddr decsdk.ValAddress, delegatorAddr decsdk.AccAddress, amount sdk.Coin) MsgDelegate {
	return MsgDelegate{
		DelegatorAddress: delegatorAddr,
		ValidatorAddress: validatorAddr,
		Amount:           amount,
	}
}

const DelegateConst = "delegate"

func (msg MsgDelegate) Route() string { return RouterKey }
func (msg MsgDelegate) Type() string  { return DelegateConst }
func (msg MsgDelegate) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.DelegatorAddress)}
}

func (msg MsgDelegate) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg MsgDelegate) ValidateBasic() error {
	if msg.ValidatorAddress.Empty() {
		return ErrEmptyValidatorAddr(DefaultCodespace)
	}
	if msg.DelegatorAddress.Empty() {
		return ErrNilDelegatorAddr(DefaultCodespace)
	}
	return nil
}

// -----------------------------------------------------------------------------------------

type MsgSetOnline struct {
	ValidatorAddress decsdk.ValAddress `json:"validator_address"`
}

func NewMsgSetOnline(validatorAddr decsdk.ValAddress) MsgSetOnline {
	return MsgSetOnline{ValidatorAddress: validatorAddr}
}

const SetOnlineConst = "set-online"

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
		return ErrEmptyValidatorAddr(DefaultCodespace)
	}
	return nil
}

// -----------------------------------------------------------------------------------------

type MsgSetOffline struct {
	ValidatorAddress decsdk.ValAddress `json:"validator_address"`
}

func NewMsgSetOffline(validatorAddr decsdk.ValAddress) MsgSetOffline {
	return MsgSetOffline{ValidatorAddress: validatorAddr}
}

const SetOfflineConst = "set-offline"

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
		return ErrEmptyValidatorAddr(DefaultCodespace)
	}
	return nil
}

// -----------------------------------------------------------------------------------------

type MsgUnbond struct {
	ValidatorAddress decsdk.ValAddress `json:"validator_address"`
	DelegatorAddress decsdk.AccAddress `json:"delegator_address"`
	Amount           sdk.Coin          `json:"amount"`
}

func NewMsgUnbond(validatorAddr decsdk.ValAddress, delegatorAddr decsdk.AccAddress, amount sdk.Coin) MsgUnbond {
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
	return []sdk.AccAddress{sdk.AccAddress(msg.DelegatorAddress)}
}

func (msg MsgUnbond) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg MsgUnbond) ValidateBasic() error {
	if msg.ValidatorAddress.Empty() {
		return ErrEmptyValidatorAddr(DefaultCodespace)
	}
	if msg.DelegatorAddress.Empty() {
		return ErrNilDelegatorAddr(DefaultCodespace)
	}
	return nil
}

// -----------------------------------------------------------------------------------------

type MsgEditCandidate struct {
	PubKey           crypto.PubKey     `json:"pub_key"`
	ValidatorAddress decsdk.ValAddress `json:"validator_address"`
	RewardAddress    decsdk.AccAddress `json:"reward_address"`
}

func NewMsgEditCandidate(pubKey crypto.PubKey, validatorAddress decsdk.ValAddress, rewardAddress decsdk.AccAddress) MsgEditCandidate {
	return MsgEditCandidate{
		PubKey:           pubKey,
		ValidatorAddress: validatorAddress,
		RewardAddress:    rewardAddress,
	}
}

const EditCandidateConst = "edit-candidate"

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
	if msg.PubKey == nil {
		return ErrEmptyPubKey(DefaultCodespace)
	}
	if msg.ValidatorAddress.Empty() {
		return ErrEmptyValidatorAddr(DefaultCodespace)
	}
	return nil
}
