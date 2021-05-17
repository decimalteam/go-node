package types

import (
	"bitbucket.org/decimalteam/go-node/x/validator/exported"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"strings"
	"time"
)

// DVPair is struct that just has a delegator-validator pair with no other data.
// It is intended to be used as a marshalable pointer. For example, a DVPair can be used to construct the
// key to getting an UnbondingDelegation from state.
type DVPair struct {
	DelegatorAddress sdk.AccAddress
	ValidatorAddress sdk.ValAddress
}

// DVVTriplet is struct that just has a delegator-validator-validator triplet with no other data.
// It is intended to be used as a marshalable pointer. For example, a DVVTriplet can be used to construct the
// key to getting a Redelegation from state.
type DVVTriplet struct {
	DelegatorAddress    sdk.AccAddress
	ValidatorSrcAddress sdk.ValAddress
	ValidatorDstAddress sdk.ValAddress
}

// Delegation represents the bond with tokens held by an account. It is
// owned by one delegator, and is associated with the voting power of one
// validator.
type Delegation struct {
	DelegatorAddress sdk.AccAddress `json:"delegator_address" yaml:"delegator_address"`
	ValidatorAddress sdk.ValAddress `json:"validator_address" yaml:"validator_address"`
	Coin             sdk.Coin       `json:"coin" yaml:"coin"`
	TokensBase       sdk.Int        `json:"tokens_base" yaml:"tokens_base"`
}

// NewDelegation creates a new delegation object
func NewDelegation(delegatorAddr sdk.AccAddress, validatorAddr sdk.ValAddress, coin sdk.Coin) Delegation {
	return Delegation{
		DelegatorAddress: delegatorAddr,
		ValidatorAddress: validatorAddr,
		Coin:             coin,
	}
}

// return the delegation
func MustMarshalDelegation(cdc *codec.Codec, delegation Delegation) []byte {
	return cdc.MustMarshalBinaryLengthPrefixed(delegation)
}

// return the delegation
func MustUnmarshalDelegation(cdc *codec.Codec, value []byte) Delegation {
	delegation, err := UnmarshalDelegation(cdc, value)
	if err != nil {
		panic(err)
	}
	return delegation
}

// return the delegation
func UnmarshalDelegation(cdc *codec.Codec, value []byte) (delegation Delegation, err error) {
	err = cdc.UnmarshalBinaryLengthPrefixed(value, &delegation)
	return delegation, err
}

// nolint
func (d Delegation) Equal(d2 Delegation) bool {
	return bytes.Equal(d.DelegatorAddress, d2.DelegatorAddress) &&
		bytes.Equal(d.ValidatorAddress, d2.ValidatorAddress)
}

// nolint - for Delegation
func (d Delegation) GetDelegatorAddr() sdk.AccAddress { return d.DelegatorAddress }
func (d Delegation) GetValidatorAddr() sdk.ValAddress { return d.ValidatorAddress }
func (d Delegation) GetCoin() sdk.Coin                { return d.Coin }
func (d Delegation) GetTokensBase() sdk.Int           { return d.TokensBase }
func (d Delegation) SetTokensBase(tokensBase sdk.Int) exported.DelegationI {
	d.TokensBase = tokensBase
	return d
}

// String returns a human readable string representation of a Delegation.
func (d Delegation) String() string {
	return fmt.Sprintf(`Delegation:
  Delegator:  %s
  Validator:  %s
  Coin:       %s%s
  TokensBase: %s`, d.DelegatorAddress,
		d.ValidatorAddress, d.Coin.Amount, d.Coin.Denom, d.TokensBase)
}

// Delegations is a collection of delegations
type Delegations []Delegation

func (d Delegations) String() (out string) {
	for _, del := range d {
		out += del.String() + "\n"
	}
	return strings.TrimSpace(out)
}

func GetBaseDelegations(delegations []exported.DelegationI) Delegations {
	var dels Delegations
	for _, delegation := range delegations {
		switch delegation := delegation.(type) {
		case Delegation:
			dels = append(dels, delegation)
		}
	}
	return dels
}

// UnbondingDelegation stores all of a single delegator's unbonding bonds
// for a single validator in an time-ordered list
type UnbondingDelegation struct {
	DelegatorAddress sdk.AccAddress                       `json:"delegator_address" yaml:"delegator_address"` // delegator
	ValidatorAddress sdk.ValAddress                       `json:"validator_address" yaml:"validator_address"` // validator unbonding from operator addr
	Entries          []exported.UnbondingDelegationEntryI `json:"entries" yaml:"entries"`                     // unbonding delegation entries
}

// UnbondingDelegationEntry - entry to an UnbondingDelegation
type UnbondingDelegationEntry struct {
	CreationHeight int64     `json:"creation_height" yaml:"creation_height"` // height which the unbonding took place
	CompletionTime time.Time `json:"completion_time" yaml:"completion_time"` // time at which the unbonding delegation will complete
	InitialBalance sdk.Coin  `json:"initial_balance" yaml:"initial_balance"` // atoms initially scheduled to receive at completion
	Balance        sdk.Coin  `json:"balance" yaml:"balance"`                 // atoms to receive at completion
}

func (e UnbondingDelegationEntry) GetCreationHeight() int64     { return e.CreationHeight }
func (e UnbondingDelegationEntry) GetCompletionTime() time.Time { return e.CompletionTime }
func (e UnbondingDelegationEntry) GetBalance() sdk.Coin         { return e.Balance }
func (e UnbondingDelegationEntry) GetInitialBalance() sdk.Coin  { return e.InitialBalance }
func (e UnbondingDelegationEntry) GetEvent() sdk.Event {
	return sdk.NewEvent(
		EventTypeCompleteUnbonding,
	)
}

// IsMature - is the current entry mature
func (e UnbondingDelegationEntry) IsMature(currentTime time.Time) bool {
	return !e.CompletionTime.After(currentTime)
}

func (e UnbondingDelegationEntry) String() string {
	return fmt.Sprintf(`      Creation Height:           %v
      Min time to unbond (unix): %v
      Expected balance:          %s`,
		e.CreationHeight,
		e.CompletionTime,
		e.Balance.String())
}

// NewUnbondingDelegation - create a new unbonding delegation object
func NewUnbondingDelegation(delegatorAddr sdk.AccAddress,
	validatorAddr sdk.ValAddress,
	entry exported.UnbondingDelegationEntryI) UnbondingDelegation {

	return UnbondingDelegation{
		DelegatorAddress: delegatorAddr,
		ValidatorAddress: validatorAddr,
		Entries:          []exported.UnbondingDelegationEntryI{entry},
	}
}

// NewUnbondingDelegation - create a new unbonding delegation object
func NewUnbondingDelegationEntry(creationHeight int64, completionTime time.Time,
	balance sdk.Coin) UnbondingDelegationEntry {

	return UnbondingDelegationEntry{
		CreationHeight: creationHeight,
		CompletionTime: completionTime,
		InitialBalance: balance,
		Balance:        balance,
	}
}

// AddEntry - append entry to the unbonding delegation
func (d *UnbondingDelegation) AddEntry(creationHeight int64,
	minTime time.Time, balance sdk.Coin) {

	d.Entries = append(d.Entries, NewUnbondingDelegationEntry(creationHeight, minTime, balance))
}

func (d *UnbondingDelegation) AddNFTEntry(creationHeight int64, minTime time.Time, tokenID, denom string, quantity sdk.Int, balance sdk.Coin) {
	d.Entries = append(d.Entries, NewUnbondingDelegationNFTEntry(creationHeight, minTime, denom, tokenID, quantity, balance))
}

// RemoveEntry - remove entry at index i to the unbonding delegation
func (d *UnbondingDelegation) RemoveEntry(i int64) {
	d.Entries = append(d.Entries[:i], d.Entries[i+1:]...)
}

func (d UnbondingDelegation) GetEvents(ctxTime time.Time) sdk.Events {
	events := sdk.Events{}
	for _, entry := range d.Entries {
		if entry.IsMature(ctxTime) {
			events = events.AppendEvent(entry.GetEvent().AppendAttributes(
				sdk.NewAttribute(AttributeKeyValidator, d.ValidatorAddress.String()),
				sdk.NewAttribute(AttributeKeyDelegator, d.DelegatorAddress.String()),
				sdk.NewAttribute(AttributeKeyCoin, entry.GetBalance().String()),
			))
		}
	}
	return events
}

// return the unbonding delegation
func MustMarshalUBD(cdc *codec.Codec, ubd UnbondingDelegation) []byte {
	return cdc.MustMarshalBinaryLengthPrefixed(ubd)
}

// unmarshal a unbonding delegation from a store value
func MustUnmarshalUBD(cdc *codec.Codec, value []byte) UnbondingDelegation {
	ubd, err := UnmarshalUBD(cdc, value)
	if err != nil {
		panic(err)
	}
	return ubd
}

// unmarshal a unbonding delegation from a store value
func UnmarshalUBD(cdc *codec.Codec, value []byte) (ubd UnbondingDelegation, err error) {
	err = cdc.UnmarshalBinaryLengthPrefixed(value, &ubd)
	return ubd, err
}

// nolint
// inefficient but only used in testing
func (d UnbondingDelegation) Equal(d2 UnbondingDelegation) bool {
	bz1 := ModuleCdc.MustMarshalBinaryLengthPrefixed(&d)
	bz2 := ModuleCdc.MustMarshalBinaryLengthPrefixed(&d2)
	return bytes.Equal(bz1, bz2)
}

// String returns a human readable string representation of an UnbondingDelegation.
func (d UnbondingDelegation) String() string {
	out := fmt.Sprintf(`Unbonding Delegations between:
  Delegator:                 %s
  Validator:                 %s
	Entries:`, d.DelegatorAddress, d.ValidatorAddress)
	for i, entry := range d.Entries {
		out += fmt.Sprintf(`    Unbonding Delegation %d:
      %s`, i, entry)
	}
	return out
}

// UnbondingDelegations is a collection of UnbondingDelegation
type UnbondingDelegations []UnbondingDelegation

func (ubds UnbondingDelegations) String() (out string) {
	for _, u := range ubds {
		out += u.String() + "\n"
	}
	return strings.TrimSpace(out)
}

// ----------------------------------------------------------------------------
// Client Types

// DelegationResponse is equivalent to Delegation except that it contains a balance
// in addition to shares which is more suitable for client responses.
type DelegationResponse struct {
	Delegations    `json:"delegations"`
	DelegationsNFT `json:"delegations_nft"`
}

func NewDelegationResp(delegations Delegations, delegationsNFT DelegationsNFT) DelegationResponse {
	return DelegationResponse{
		Delegations:    delegations,
		DelegationsNFT: delegationsNFT,
	}
}

type delegationRespAlias DelegationResponse

// MarshalJSON implements the json.Marshaler interface. This is so we can
// achieve a flattened structure while embedding other types.
func (d DelegationResponse) MarshalJSON() ([]byte, error) {
	return json.Marshal((delegationRespAlias)(d))
}

// UnmarshalJSON implements the json.Unmarshaler interface. This is so we can
// achieve a flattened structure while embedding other types.
func (d *DelegationResponse) UnmarshalJSON(bz []byte) error {
	return json.Unmarshal(bz, (*delegationRespAlias)(d))
}

// DelegationResponses is a collection of DelegationResp
type DelegationResponses []DelegationResponse

// String implements the Stringer interface for DelegationResponses.
func (d DelegationResponses) String() (out string) {
	for _, del := range d {
		out += del.String() + "\n"
	}
	return strings.TrimSpace(out)
}
