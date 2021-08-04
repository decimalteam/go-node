package types

import (
	"bytes"
	"encoding/json"
	"fmt"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	abci "github.com/tendermint/tendermint/abci/types"
	tmprotocrypto "github.com/tendermint/tendermint/proto/tendermint/crypto"
	"sort"
	"strings"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"gopkg.in/yaml.v2"
)

// nolint
const (
	MaxMonikerLength         = 70
	MaxIdentityLength        = 3000
	MaxWebsiteLength         = 140
	MaxSecurityContactLength = 140
	MaxDetailsLength         = 280
)

//type Validator struct {
//	ValAddress              sdk.ValAddress `json:"val_address" yaml:"val_address"`
//	PubKey                  types.PubKey  `json:"pub_key" yaml:"pub_key"`
//	Tokens                  sdk.Int        `json:"stake_coins" yaml:"stake_coins"`
//	Status                  BondStatus     `json:"status" yaml:"status"`
//	Commission              sdk.Dec        `json:"commission" yaml:"commission"`
//	Jailed                  bool           `json:"jailed" yaml:"jailed"`
//	UnbondingCompletionTime time.Time      `json:"unbonding_completion_time" yaml:"unbonding_completion_time"`
//	UnbondingHeight         int64          `json:"unbonding_height" yaml:"unbonding_height"`
//	Description             Description    `json:"description" yaml:"description"`
//	AccumRewards            sdk.Int        `json:"accum_rewards" yaml:"accum_rewards"`
//	RewardAddress           sdk.AccAddress `json:"reward_address" yaml:"reward_address"`
//	Online                  bool           `json:"online" yaml:"online"`
//}

type Stake struct {
	Delegator sdk.AccAddress `json:"delegator"`
	Coin      sdk.Coin       `json:"coin"`
}

// String returns a human readable string representation of a validator.
func (v *Validator) String() string {
	pk, err := v.GetConsPubKey()
	if err != nil {
		panic(err)
	}

	return fmt.Sprintf(`Validator
  Operator Address:           %s
  Validator Consensus Pubkey: %s
  Jailed:                     %v
  Online:                     %v
  Status:                     %s
  Tokens:                     %s
  Description:                %s
  Unbonding Height:           %d
  Unbonding Completion Time:  %v
  Commission:                 %s
  Accum Rewards:              %s`, v.ValAddress, pk.String(),
		v.Jailed, v.Online, v.Status, v.Tokens,
		v.Description,
		v.UnbondingHeight, v.UnbondingCompletionTime, v.Commission, v.AccumRewards)
}

//this is a helper struct used for JSON de- and encoding only
//type bechValidator struct {
//	ValAddress              sdk.ValAddress `json:"val_address" yaml:"val_address"`
//	PubKey                  string         `json:"pub_key" yaml:"pub_key"`
//	Tokens                  sdk.Int        `json:"stake_coins" yaml:"stake_coins"`
//	Status                  BondStatus     `json:"status" yaml:"status"`
//	Commission              sdk.Dec        `json:"commission" yaml:"commission"`
//	Jailed                  bool           `json:"jailed" yaml:"jailed"`
//	UnbondingCompletionTime time.Time      `json:"unbonding_completion_time" yaml:"unbonding_completion_time"`
//	UnbondingHeight         int64          `json:"unbonding_height" yaml:"unbonding_height"`
//	Description             Description    `json:"description" yaml:"description"`
//	AccumRewards            sdk.Int        `json:"accum_rewards" yaml:"accum_rewards"`
//	RewardAddress           sdk.AccAddress `json:"reward_address" yaml:"reward_address"`
//	Online                  bool           `json:"online" yaml:"online"`
//}

// MarshalJSON marshals the validator to JSON using Bech32
func (v Validator) MarshalJSON() ([]byte, error) {
	return ModuleCdc.MarshalJSON(&BechValidator{
		ValAddress:              v.ValAddress,
		PubKey:                  v.PubKey,
		Jailed:                  v.Jailed,
		Status:                  v.Status,
		Tokens:                  v.Tokens,
		Description:             v.Description,
		UnbondingHeight:         v.UnbondingHeight,
		UnbondingCompletionTime: v.UnbondingCompletionTime,
		Commission:              v.Commission,
		RewardAddress:           v.RewardAddress,
		Online:                  v.Online,
		AccumRewards:            v.AccumRewards,
	})
}

// UnmarshalJSON unmarshals the validator from JSON using Bech32
func (v *Validator) UnmarshalJSON(data []byte) error {
	bv := &BechValidator{}
	if err := json.Unmarshal(data, bv); err != nil {
		return err
	}

	*v = Validator{
		ValAddress:              bv.ValAddress,
		PubKey:                  bv.PubKey,
		Jailed:                  bv.Jailed,
		Tokens:                  bv.Tokens,
		Status:                  bv.Status,
		Description:             bv.Description,
		UnbondingHeight:         bv.UnbondingHeight,
		UnbondingCompletionTime: bv.UnbondingCompletionTime,
		Commission:              bv.Commission,
		RewardAddress:           bv.RewardAddress,
		Online:                  bv.Online,
		AccumRewards:            bv.AccumRewards,
	}
	return nil
}

// custom marshal yaml function due to consensus pubkey
func (v Validator) MarshalYAML() (interface{}, error) {
	pk, _ := v.PubKey.GetCachedValue().(cryptotypes.PubKey)

	bs, err := yaml.Marshal(struct {
		ValAddress/*sdk.ValAddress*/ string
		RewardAddress/*sdk.AccAddress*/ string
		PubKey                  string
		Jailed                  bool
		Status                  BondStatus
		Tokens                  sdk.Int
		Description             Description
		UnbondingHeight         int64
		UnbondingCompletionTime time.Time
		Commission              sdk.Dec
		AccumRewards            sdk.Int
		Online                  bool
	}{
		ValAddress:              v.ValAddress,
		RewardAddress:           v.RewardAddress,
		PubKey:                  sdk.MustBech32ifyAddressBytes(sdk.Bech32PrefixConsPub, pk.Bytes()),
		Jailed:                  v.Jailed,
		Status:                  v.Status,
		Tokens:                  v.Tokens,
		Description:             v.Description,
		UnbondingHeight:         v.UnbondingHeight,
		UnbondingCompletionTime: v.UnbondingCompletionTime,
		Commission:              v.Commission,
		AccumRewards:            v.AccumRewards,
		Online:                  v.Online,
	})
	if err != nil {
		return nil, err
	}

	return string(bs), nil
}

func (v Validator) SharesFromTokens(tokens sdk.Int, valTokens sdk.Int, delTokens sdk.Dec) (sdk.Dec, error) {
	if v.Tokens.IsZero() {
		return sdk.Dec{}, ErrInsufficientShares()
	}

	return delTokens.MulInt(tokens).QuoInt(valTokens), nil
}

func (v Validator) IsJailed() bool        { return v.Jailed }
func (v Validator) GetMoniker() string    { return v.Description.Moniker }
func (v Validator) GetStatus() BondStatus { return v.Status }
func (v Validator) GetOperator() sdk.ValAddress {
	valAddr, err := sdk.ValAddressFromBech32(v.ValAddress)
	if err != nil {
		panic(err)
	}

	return valAddr
}
func (v Validator) GetConsPubKey() (cryptotypes.PubKey, error) {
	pk, ok := v.PubKey.GetCachedValue().(cryptotypes.PubKey)
	if !ok {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidType, "expecting cryptotypes.PubKey, got %T", pk)
	}

	return pk, nil
}

func (v Validator) TmConsPublicKey() (tmprotocrypto.PublicKey, error) {
	pk, err := v.GetConsPubKey()
	if err != nil {
		return tmprotocrypto.PublicKey{}, err
	}

	tmPk, err := cryptocodec.ToTmProtoPublicKey(pk)
	if err != nil {
		return tmprotocrypto.PublicKey{}, err
	}

	return tmPk, nil
}

func (v Validator) GetConsAddr() (sdk.ConsAddress, error) {
	pk, ok := v.PubKey.GetCachedValue().(cryptotypes.PubKey)
	if !ok {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidType, "wrong type %T", pk)
	}

	return sdk.ConsAddress(pk.Address()), nil
}
func (v Validator) GetTokens() sdk.Int       { return v.Tokens }
func (v Validator) GetBondedTokens() sdk.Int { return v.BondedTokens() }
func (v Validator) GetCommission() sdk.Dec   { return v.Commission }

func (v Validator) UnpackInterfaces(unpacker codectypes.AnyUnpacker) error {
	var pk cryptotypes.PubKey
	return unpacker.UnpackAny(v.PubKey, &pk)
}

type Validators []Validator

func (v Validators) String() string {
	var out string
	for _, val := range v {
		out += val.String() + "\n"
	}
	return strings.TrimSpace(out)
}

// Sort Validators sorts validator array in ascending operator address order
func (v Validators) Sort() {
	sort.Sort(v)
}

// Implements sort interface
func (v Validators) Len() int {
	return len(v)
}

// Implements sort interface
func (v Validators) Less(i, j int) bool {
	return bytes.Compare([]byte(v[i].ValAddress), []byte(v[j].ValAddress)) == -1
}

// Implements sort interface
func (v Validators) Swap(i, j int) {
	it := v[i]
	v[i] = v[j]
	v[j] = it
}

func (v Validators) UnpackInterfaces(u codectypes.AnyUnpacker) error {
	for i := range v {
		if err := v[i].UnpackInterfaces(u); err != nil {
			return err
		}
	}
	return nil
}

// constant used in flags to indicate that description field should not be updated
const DoNotModifyDesc = "[do-not-modify]"

// Description - description fields for a validator
//type Description struct {
//	Moniker         string `json:"moniker" yaml:"moniker"`                   // name
//	Identity        string `json:"identity" yaml:"identity"`                 // optional identity signature (ex. UPort or Keybase)
//	Website         string `json:"website" yaml:"website"`                   // optional website link
//	SecurityContact string `json:"security_contact" yaml:"security_contact"` // optional security contact info
//	Details         string `json:"details" yaml:"details"`                   // optional details
//}

// NewDescription returns a new Description with the provided values.
func NewDescription(moniker, identity, website, securityContact, details string) Description {
	return Description{
		Moniker:         moniker,
		Identity:        identity,
		Website:         website,
		SecurityContact: securityContact,
		Details:         details,
	}
}

// UpdateDescription updates the fields of a given description. An error is
// returned if the resulting description contains an invalid length.
func (d Description) UpdateDescription(d2 Description) (Description, error) {
	if d2.Moniker == DoNotModifyDesc {
		d2.Moniker = d.Moniker
	}
	if d2.Identity == DoNotModifyDesc {
		d2.Identity = d.Identity
	}
	if d2.Website == DoNotModifyDesc {
		d2.Website = d.Website
	}
	if d2.SecurityContact == DoNotModifyDesc {
		d2.SecurityContact = d.SecurityContact
	}
	if d2.Details == DoNotModifyDesc {
		d2.Details = d.Details
	}

	return NewDescription(
		d2.Moniker,
		d2.Identity,
		d2.Website,
		d2.SecurityContact,
		d2.Details,
	).EnsureLength()
}

// EnsureLength ensures the length of a validator's description.
func (d Description) EnsureLength() (Description, error) {
	if len(d.Moniker) > MaxMonikerLength {
		return d, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid moniker length; got: %d, max: %d", len(d.Moniker), MaxMonikerLength)
	}
	if len(d.Identity) > MaxIdentityLength {
		return d, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid identity length; got: %d, max: %d", len(d.Identity), MaxIdentityLength)
	}
	if len(d.Website) > MaxWebsiteLength {
		return d, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid website length; got: %d, max: %d", len(d.Website), MaxWebsiteLength)
	}
	if len(d.SecurityContact) > MaxSecurityContactLength {
		return d, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid security contact length; got: %d, max: %d", len(d.SecurityContact), MaxSecurityContactLength)
	}
	if len(d.Details) > MaxDetailsLength {
		return d, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid details length; got: %d, max: %d", len(d.Details), MaxDetailsLength)
	}

	return d, nil
}

// Equal compares two BondStatus instances
func (b BondStatus) Equal(b2 BondStatus) bool {
	return byte(b) == byte(b2)
}

func NewValidator(valAddress string, pubKey cryptotypes.PubKey, commission sdk.Dec, rewardAddress string, description Description) (Validator, error) {
	pkAny, err := codectypes.NewAnyWithValue(pubKey)
	if err != nil {
		return Validator{}, err
	}

	return Validator{
		ValAddress:              valAddress,
		PubKey:                  pkAny,
		Jailed:                  false,
		Tokens:                  sdk.ZeroInt(),
		Description:             description,
		Status:                  Unbonded,
		Commission:              commission,
		RewardAddress:           rewardAddress,
		UnbondingCompletionTime: time.Unix(0, 0).UTC(),
		AccumRewards:            sdk.ZeroInt(),
		Online:                  true,
	}, nil
}

// unmarshal a validator from a store value
func UnmarshalValidator(cdc *codec.LegacyAmino, value []byte) (validator Validator, err error) {
	err = cdc.UnmarshalBinaryLengthPrefixed(value, &validator)
	return validator, err
}

// unmarshal a validator from a store value
func MustUnmarshalValidator(cdc *codec.LegacyAmino, value []byte) Validator {
	validator, err := UnmarshalValidator(cdc, value)
	if err != nil {
		panic(err)
	}
	return validator
}

// IsBonded checks if the validator status equals Bonded
func (v Validator) IsBonded() bool {
	return v.Status.Equal(Bonded)
}

// IsUnbonded checks if the validator status equals Unbonded
func (v Validator) IsUnbonded() bool {
	return v.Status.Equal(Unbonded)
}

// IsUnbonding checks if the validator status equals Unbonding
func (v Validator) IsUnbonding() bool {
	return v.Status.Equal(Unbonding)
}

// UpdateStatus updates the location of the shares within a validator
// to reflect the new status
func (v Validator) UpdateStatus(newStatus BondStatus) Validator {
	v.Status = newStatus
	return v
}

// get the consensus-engine power
// a reduction of 10^6 from validator tokens is applied
func (v Validator) ConsensusPower() int64 {
	if v.IsBonded() {
		return v.PotentialConsensusPower()
	}
	return 0
}

// for exported
func (v Validator) GetConsensusPower() int64 {
	return v.ConsensusPower()
}

// potential consensus-engine power
func (v Validator) PotentialConsensusPower() int64 {
	return TokensToConsensusPower(v.Tokens)
}

// ABCIValidatorUpdate returns an abci.ValidatorUpdate from a staking validator type
// with the full validator power
func (v Validator) ABCIValidatorUpdate() abci.ValidatorUpdate {
	tmPk, err := v.TmConsPublicKey()
	if err != nil {
		panic(err)
	}

	return abci.ValidatorUpdate{
		PubKey: tmPk,
		Power:  v.ConsensusPower(),
	}
}

// ABCIValidatorUpdateZero returns an abci.ValidatorUpdate from a staking validator type
// with zero power used for validator updates.
func (v Validator) ABCIValidatorUpdateZero() abci.ValidatorUpdate {
	tmPk, err := v.TmConsPublicKey()
	if err != nil {
		panic(err)
	}

	return abci.ValidatorUpdate{
		PubKey: tmPk,
		Power:  0,
	}
}

// get the bonded tokens which the validator holds
func (v Validator) BondedTokens() sdk.Int {
	if v.IsBonded() {
		return v.Tokens
	}
	return sdk.ZeroInt()
}

// In some situations, the exchange rate becomes invalid, e.g. if
// Validator loses all tokens due to slashing. In this case,
// make all future delegations invalid.
func (v Validator) InvalidExRate() bool {
	return v.Tokens.IsZero()
}

// RemoveTokens removes tokens from a validator
func (v Validator) RemoveTokens(tokens sdk.Int) Validator {
	if tokens.IsNegative() {
		panic(fmt.Sprintf("should not happen: trying to remove negative tokens %v", tokens))
	}
	if v.Tokens.LT(tokens) {
		panic(fmt.Sprintf("should not happen: only have %v tokens, trying to remove %v", v.Tokens, tokens))
	}
	v.Tokens = v.Tokens.Sub(tokens)
	return v
}

func (v Validator) AddAccumReward(reward sdk.Int) Validator {
	v.AccumRewards = v.AccumRewards.Add(reward)
	return v
}

func (v Validator) TestEquivalent(v2 Validator) bool {
	pk1, err := v.GetConsPubKey()
	if err != nil {
		panic(err)
	}

	pk2, err := v2.GetConsPubKey()
	if err != nil {
		panic(err)
	}

	return pk1.Equals(pk2) &&
		bytes.Equal([]byte(v.ValAddress), []byte(v2.ValAddress)) &&
		v.Status.Equal(v2.Status) &&
		v.Tokens.Equal(v2.Tokens) &&
		v.Description == v2.Description &&
		v.Commission.Equal(v2.Commission)
}
