package types

import (
	"bytes"
	"fmt"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"sort"
	"strings"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto"
	tmtypes "github.com/tendermint/tendermint/types"
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

type Validator struct {
	ValAddress              sdk.ValAddress `json:"val_address" yaml:"val_address"`
	PubKey                  crypto.PubKey  `json:"pub_key" yaml:"pub_key"`
	Tokens                  sdk.Int        `json:"stake_coins" yaml:"stake_coins"`
	Status                  BondStatus     `json:"status" yaml:"status"`
	Commission              sdk.Dec        `json:"commission" yaml:"commission"`
	Jailed                  bool           `json:"jailed" yaml:"jailed"`
	UnbondingCompletionTime time.Time      `json:"unbonding_completion_time" yaml:"unbonding_completion_time"`
	UnbondingHeight         int64          `json:"unbonding_height" yaml:"unbonding_height"`
	Description             Description    `json:"description" yaml:"description"`
	AccumRewards            sdk.Int        `json:"accum_rewards" yaml:"accum_rewards"`
	RewardAddress           sdk.AccAddress `json:"reward_address" yaml:"reward_address"`
	Online                  bool           `json:"online" yaml:"online"`
}

type Stake struct {
	Delegator sdk.AccAddress `json:"delegator"`
	Coin      sdk.Coin       `json:"coin"`
}

// String returns a human readable string representation of a validator.
func (v Validator) String() string {
	bechConsPubKey, err := sdk.Bech32ifyPubKey(sdk.Bech32PubKeyTypeConsPub, v.PubKey)
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
  Accum Rewards:              %s`, v.ValAddress, bechConsPubKey,
		v.Jailed, v.Online, v.Status, v.Tokens,
		v.Description,
		v.UnbondingHeight, v.UnbondingCompletionTime, v.Commission, v.AccumRewards)
}

// this is a helper struct used for JSON de- and encoding only
type bechValidator struct {
	ValAddress              sdk.ValAddress `json:"val_address" yaml:"val_address"`
	PubKey                  string         `json:"pub_key" yaml:"pub_key"`
	Tokens                  sdk.Int        `json:"stake_coins" yaml:"stake_coins"`
	Status                  BondStatus     `json:"status" yaml:"status"`
	Commission              sdk.Dec        `json:"commission" yaml:"commission"`
	Jailed                  bool           `json:"jailed" yaml:"jailed"`
	UnbondingCompletionTime time.Time      `json:"unbonding_completion_time" yaml:"unbonding_completion_time"`
	UnbondingHeight         int64          `json:"unbonding_height" yaml:"unbonding_height"`
	Description             Description    `json:"description" yaml:"description"`
	AccumRewards            sdk.Int        `json:"accum_rewards" yaml:"accum_rewards"`
	RewardAddress           sdk.AccAddress `json:"reward_address" yaml:"reward_address"`
	Online                  bool           `json:"online" yaml:"online"`
}

// MarshalJSON marshals the validator to JSON using Bech32
func (v Validator) MarshalJSON() ([]byte, error) {
	bechConsPubKey, err := sdk.Bech32ifyPubKey(sdk.Bech32PubKeyTypeConsPub, v.PubKey)
	if err != nil {
		return nil, err
	}

	return codec.Cdc.MarshalJSON(bechValidator{
		ValAddress:              v.ValAddress,
		PubKey:                  bechConsPubKey,
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
	bv := &bechValidator{}
	if err := codec.Cdc.UnmarshalJSON(data, bv); err != nil {
		return err
	}
	consPubKey, err := sdk.GetPubKeyFromBech32(sdk.Bech32PubKeyTypeConsPub, bv.PubKey)
	if err != nil {
		return err
	}
	*v = Validator{
		ValAddress:              bv.ValAddress,
		PubKey:                  consPubKey,
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
	bs, err := yaml.Marshal(struct {
		ValAddress              sdk.ValAddress
		RewardAddress           sdk.AccAddress
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
		PubKey:                  sdk.MustBech32ifyPubKey(sdk.Bech32PubKeyTypeConsPub, v.PubKey),
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

func (v Validator) IsJailed() bool               { return v.Jailed }
func (v Validator) GetMoniker() string           { return v.Description.Moniker }
func (v Validator) GetStatus() BondStatus        { return v.Status }
func (v Validator) GetOperator() sdk.ValAddress  { return v.ValAddress }
func (v Validator) GetConsPubKey() crypto.PubKey { return v.PubKey }
func (v Validator) GetConsAddr() sdk.ConsAddress { return sdk.ConsAddress(v.PubKey.Address()) }
func (v Validator) GetTokens() sdk.Int           { return v.Tokens }
func (v Validator) GetBondedTokens() sdk.Int     { return v.BondedTokens() }
func (v Validator) GetCommission() sdk.Dec       { return v.Commission }

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
	return bytes.Compare(v[i].ValAddress, v[j].ValAddress) == -1
}

// Implements sort interface
func (v Validators) Swap(i, j int) {
	it := v[i]
	v[i] = v[j]
	v[j] = it
}

// constant used in flags to indicate that description field should not be updated
const DoNotModifyDesc = "[do-not-modify]"

// Description - description fields for a validator
type Description struct {
	Moniker         string `json:"moniker" yaml:"moniker"`                   // name
	Identity        string `json:"identity" yaml:"identity"`                 // optional identity signature (ex. UPort or Keybase)
	Website         string `json:"website" yaml:"website"`                   // optional website link
	SecurityContact string `json:"security_contact" yaml:"security_contact"` // optional security contact info
	Details         string `json:"details" yaml:"details"`                   // optional details
}

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

// BondStatus is the status of a validator
type BondStatus byte

// staking constants
const (
	Unbonded  BondStatus = 0x00
	Unbonding BondStatus = 0x01
	Bonded    BondStatus = 0x02

	BondStatusUnbonded  = "Unbonded"
	BondStatusUnbonding = "Unbonding"
	BondStatusBonded    = "Bonded"
)

// Equal compares two BondStatus instances
func (b BondStatus) Equal(b2 BondStatus) bool {
	return byte(b) == byte(b2)
}

// String implements the Stringer interface for BondStatus.
func (b BondStatus) String() string {
	switch b {
	case 0x00:
		return BondStatusUnbonded
	case 0x01:
		return BondStatusUnbonding
	case 0x02:
		return BondStatusBonded
	default:
		panic("invalid bond status")
	}
}

func NewValidator(valAddress sdk.ValAddress, pubKey crypto.PubKey, commission sdk.Dec, rewardAddress sdk.AccAddress, description Description) Validator {
	return Validator{
		ValAddress:              valAddress,
		PubKey:                  pubKey,
		Jailed:                  false,
		Tokens:                  sdk.ZeroInt(),
		Description:             description,
		Status:                  Unbonded,
		Commission:              commission,
		RewardAddress:           rewardAddress,
		UnbondingCompletionTime: time.Unix(0, 0).UTC(),
		AccumRewards:            sdk.ZeroInt(),
		Online:                  true,
	}
}

// unmarshal a validator from a store value
func UnmarshalValidator(cdc *codec.Codec, value []byte) (validator Validator, err error) {
	err = cdc.UnmarshalBinaryLengthPrefixed(value, &validator)
	return validator, err
}

// unmarshal a validator from a store value
func MustUnmarshalValidator(cdc *codec.Codec, value []byte) Validator {
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

// potential consensus-engine power
func (v Validator) PotentialConsensusPower() int64 {
	return TokensToConsensusPower(v.Tokens)
}

// for exported
func (v Validator) GetConsensusPower() int64 {
	return v.ConsensusPower()
}

// ABCIValidatorUpdate returns an abci.ValidatorUpdate from a staking validator type
// with the full validator power
func (v Validator) ABCIValidatorUpdate() abci.ValidatorUpdate {
	return abci.ValidatorUpdate{
		PubKey: tmtypes.TM2PB.PubKey(v.PubKey),
		Power:  v.ConsensusPower(),
	}
}

// ABCIValidatorUpdateZero returns an abci.ValidatorUpdate from a staking validator type
// with zero power used for validator updates.
func (v Validator) ABCIValidatorUpdateZero() abci.ValidatorUpdate {
	return abci.ValidatorUpdate{
		PubKey: tmtypes.TM2PB.PubKey(v.PubKey),
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
	return v.PubKey.Equals(v2.PubKey) &&
		bytes.Equal(v.ValAddress, v2.ValAddress) &&
		v.Status.Equal(v2.Status) &&
		v.Tokens.Equal(v2.Tokens) &&
		v.Description == v2.Description &&
		v.Commission.Equal(v2.Commission)
}
