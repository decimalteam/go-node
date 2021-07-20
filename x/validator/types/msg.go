package types

import (
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	_ sdk.Msg                            = &MsgDeclareCandidate{}
	_ codectypes.UnpackInterfacesMessage = (*MsgDeclareCandidate)(nil)
	_ sdk.Msg                            = &MsgDelegate{}
	_ sdk.Msg                            = &MsgDelegateNFT{}
	_ sdk.Msg                            = &MsgUnbondNFT{}
	_ sdk.Msg                            = &MsgUnbond{}
	_ sdk.Msg                            = &MsgEditCandidate{}
)

//type MsgDeclareCandidate struct {
//	Commission    sdk.Dec        `json:"commission" yaml:"commission"`
//	ValidatorAddr sdk.ValAddress `json:"validator_addr" yaml:"validator_addr"`
//	RewardAddr    sdk.AccAddress `json:"reward_addr" yaml:"reward_addr"`
//	PubKey        crypto.PubKey  `json:"pub_key" yaml:"pub_key"`
//	Stake         sdk.Coin       `json:"stake" yaml:"stake"`
//	Description   Description    `json:"description"`
//}

func NewMsgDeclareCandidate(validatorAddr sdk.ValAddress, pubKey types.PubKey, commission sdk.Dec, stake sdk.Coin, description Description, rewardAddress sdk.AccAddress) (*MsgDeclareCandidate, error) {
	var pkAny *codectypes.Any
	if pubKey != nil {
		var err error
		if pkAny, err = codectypes.NewAnyWithValue(pubKey); err != nil {
			return nil, err
		}
	}

	return &MsgDeclareCandidate{
		Commission:    commission,
		ValidatorAddr: validatorAddr.String(),
		RewardAddr:    rewardAddress.String(),
		PubKey:        pkAny,
		Stake:         stake,
		Description:   description,
	}, nil
}

const DeclareCandidateConst = "declare_candidate"

func (msg MsgDeclareCandidate) Route() string { return RouterKey }
func (msg MsgDeclareCandidate) Type() string  { return DeclareCandidateConst }
func (msg MsgDeclareCandidate) GetSigners() []sdk.AccAddress {
	valaddr, err := sdk.ValAddressFromBech32(msg.ValidatorAddr)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{sdk.AccAddress(valaddr)}
}

func (msg MsgDeclareCandidate) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// quick validity check
func (msg MsgDeclareCandidate) ValidateBasic() error {
	// note that unmarshaling from bech32 ensures either empty or valid
	if len(msg.ValidatorAddr) == 0 {
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
	if msg.PubKey == nil {
		return ErrEmptyPubKey()
	}

	return nil
}

func (msg MsgDeclareCandidate) UnpackInterfaces(c codectypes.AnyUnpacker) error {
	var pk types.PubKey
	return c.UnpackAny(msg.PubKey, &pk)
}

// -----------------------------------------------------------------------------------------

//type MsgDelegate struct {
//	DelegatorAddress sdk.AccAddress `json:"delegator_address"`
//	ValidatorAddress sdk.ValAddress `json:"validator_address"`
//	Coin             sdk.Coin       `json:"coin"`
//}

func NewMsgDelegate(validatorAddr sdk.ValAddress, delegatorAddr sdk.AccAddress, coin sdk.Coin) MsgDelegate {
	return MsgDelegate{
		DelegatorAddress: delegatorAddr.String(),
		ValidatorAddress: validatorAddr.String(),
		Coin:             coin,
	}
}

const DelegateConst = "delegate"

func (msg MsgDelegate) Route() string { return RouterKey }
func (msg MsgDelegate) Type() string  { return DelegateConst }
func (msg MsgDelegate) GetSigners() []sdk.AccAddress {
	delAddr, err := sdk.AccAddressFromBech32(msg.DelegatorAddress)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{delAddr}
}

func (msg MsgDelegate) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

func (msg MsgDelegate) ValidateBasic() error {
	if len(msg.ValidatorAddress) == 0 {
		return ErrEmptyValidatorAddr()
	}
	if len(msg.DelegatorAddress) == 0 {
		return ErrEmptyDelegatorAddr()
	}
	return nil
}

// -----------------------------------------------------------------------------------------

//type MsgDelegateNFT struct {
//	DelegatorAddress sdk.AccAddress `json:"delegator_address"`
//	ValidatorAddress sdk.ValAddress `json:"validator_address"`
//	TokenID          string         `json:"id"`
//	Denom            string         `json:"denom"`
//	Quantity         sdk.Int        `json:"quantity"`
//}

func NewMsgDelegateNFT(validatorAddr sdk.ValAddress, delegatorAddr sdk.AccAddress, tokenID, denom string, quantity sdk.Int) MsgDelegateNFT {
	return MsgDelegateNFT{
		DelegatorAddress: delegatorAddr.String(),
		ValidatorAddress: validatorAddr.String(),
		TokenID:          tokenID,
		Denom:            denom,
		Quantity:         quantity,
	}
}

const DelegateNFTConst = "delegate_nft"

func (msg MsgDelegateNFT) Route() string { return RouterKey }
func (msg MsgDelegateNFT) Type() string  { return DelegateNFTConst }
func (msg MsgDelegateNFT) GetSigners() []sdk.AccAddress {
	delAddr, err := sdk.AccAddressFromBech32(msg.DelegatorAddress)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{delAddr}
}

func (msg MsgDelegateNFT) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

func (msg MsgDelegateNFT) ValidateBasic() error {
	if len(msg.ValidatorAddress) == 0 {
		return ErrEmptyValidatorAddr()
	}
	if len(msg.DelegatorAddress) == 0 {
		return ErrEmptyDelegatorAddr()
	}
	return nil
}

// -----------------------------------------------------------------------------------------

//type MsgSetOnline struct {
//	ValidatorAddress sdk.ValAddress `json:"validator_address"`
//}

func NewMsgSetOnline(validatorAddr sdk.ValAddress) MsgSetOnline {
	return MsgSetOnline{ValidatorAddress: validatorAddr.String()}
}

const SetOnlineConst = "set_online"

func (msg MsgSetOnline) Route() string { return RouterKey }
func (msg MsgSetOnline) Type() string  { return SetOnlineConst }
func (msg MsgSetOnline) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.ValidatorAddress)}
}

func (msg MsgSetOnline) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

func (msg MsgSetOnline) ValidateBasic() error {
	if len(msg.ValidatorAddress) == 0 {
		return ErrEmptyValidatorAddr()
	}
	return nil
}

// -----------------------------------------------------------------------------------------

//type MsgSetOffline struct {
//	ValidatorAddress sdk.ValAddress `json:"validator_address"`
//}

func NewMsgSetOffline(validatorAddr sdk.ValAddress) MsgSetOffline {
	return MsgSetOffline{ValidatorAddress: validatorAddr.String()}
}

const SetOfflineConst = "set_offline"

func (msg MsgSetOffline) Route() string { return RouterKey }
func (msg MsgSetOffline) Type() string  { return SetOfflineConst }
func (msg MsgSetOffline) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.ValidatorAddress)}
}

func (msg MsgSetOffline) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

func (msg MsgSetOffline) ValidateBasic() error {
	if len(msg.ValidatorAddress) == 0 {
		return ErrEmptyValidatorAddr()
	}
	return nil
}

// -----------------------------------------------------------------------------------------

//type MsgUnbond struct {
//	ValidatorAddress sdk.ValAddress `json:"validator_address"`
//	DelegatorAddress sdk.AccAddress `json:"delegator_address"`
//	Coin             sdk.Coin       `json:"coin"`
//}

func NewMsgUnbond(validatorAddr sdk.ValAddress, delegatorAddr sdk.AccAddress, coin sdk.Coin) MsgUnbond {
	return MsgUnbond{
		ValidatorAddress: validatorAddr.String(),
		DelegatorAddress: delegatorAddr.String(),
		Coin:             coin,
	}
}

const UnbondConst = "unbond"

func (msg MsgUnbond) Route() string { return RouterKey }
func (msg MsgUnbond) Type() string  { return UnbondConst }
func (msg MsgUnbond) GetSigners() []sdk.AccAddress {
	delAddr, err := sdk.AccAddressFromBech32(msg.DelegatorAddress)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{delAddr}
}

func (msg MsgUnbond) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

func (msg MsgUnbond) ValidateBasic() error {
	if len(msg.ValidatorAddress) == 0 {
		return ErrEmptyValidatorAddr()
	}
	if len(msg.DelegatorAddress) == 0 {
		return ErrEmptyDelegatorAddr()
	}
	return nil
}

// -----------------------------------------------------------------------------------------

//type MsgUnbondNFT struct {
//	ValidatorAddress sdk.ValAddress `json:"validator_address"`
//	DelegatorAddress sdk.AccAddress `json:"delegator_address"`
//	TokenID          string         `json:"id"`
//	Denom            string         `json:"denom"`
//	Quantity         sdk.Int        `json:"quantity"`
//}

func NewMsgUnbondNFT(validatorAddr sdk.ValAddress, delegatorAddr sdk.AccAddress, tokenID, denom string, quantity sdk.Int) MsgUnbondNFT {
	return MsgUnbondNFT{
		ValidatorAddress: validatorAddr.String(),
		DelegatorAddress: delegatorAddr.String(),
		TokenID:          tokenID,
		Denom:            denom,
		Quantity:         quantity,
	}
}

const UnbondNFTConst = "unbond_nft"

func (msg MsgUnbondNFT) Route() string { return RouterKey }
func (msg MsgUnbondNFT) Type() string  { return UnbondNFTConst }
func (msg MsgUnbondNFT) GetSigners() []sdk.AccAddress {
	delAddr, err := sdk.AccAddressFromBech32(msg.DelegatorAddress)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{delAddr}
}

func (msg MsgUnbondNFT) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

func (msg MsgUnbondNFT) ValidateBasic() error {
	if len(msg.ValidatorAddress) == 0 {
		return ErrEmptyValidatorAddr()
	}
	if len(msg.DelegatorAddress) == 0 {
		return ErrEmptyDelegatorAddr()
	}
	return nil
}

// -----------------------------------------------------------------------------------------

//type MsgEditCandidate struct {
//	ValidatorAddress sdk.ValAddress `json:"validator_address"`
//	RewardAddress    sdk.AccAddress `json:"reward_address"`
//	Description      Description    `json:"description"`
//}

func NewMsgEditCandidate(validatorAddress sdk.ValAddress, rewardAddress sdk.AccAddress, description Description) MsgEditCandidate {
	return MsgEditCandidate{
		ValidatorAddress: validatorAddress.String(),
		RewardAddress:    rewardAddress.String(),
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
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

func (msg MsgEditCandidate) ValidateBasic() error {
	if len(msg.ValidatorAddress) == 0 {
		return ErrEmptyValidatorAddr()
	}
	return nil
}
