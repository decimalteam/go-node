package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto"
	tmtypes "github.com/tendermint/tendermint/types"
	"time"
)

type Validator struct {
	ValAddress              sdk.ValAddress `json:"val_address"`
	PubKey                  crypto.PubKey  `json:"pub_key"`
	StakeCoins              []sdk.Coin     `json:"stake_coins"`
	DelegatorShares         sdk.Coins      `json:"delegator_shares"`
	Status                  BondStatus     `json:"status"`
	Commission              Commission     `json:"commission"`
	Jailed                  bool           `json:"jailed"`
	UnbondingCompletionTime time.Time      `json:"unbonding_completion_time"`
	UnbondingHeight         int64          `json:"unbonding_height"`
	Description             Description    `json:"description"`
}

type Validators []Validator

// Description - description fields for a validator
type Description struct {
	Moniker  string `json:"moniker" yaml:"moniker"`   // name
	Identity string `json:"identity" yaml:"identity"` // optional identity signature (ex. UPort or Keybase)
	Website  string `json:"website" yaml:"website"`   // optional website link
	Details  string `json:"details" yaml:"details"`   // optional details
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

func NewValidator(valAddress sdk.ValAddress, pubKey crypto.PubKey, coin sdk.Coin, commission Commission) Validator {
	return Validator{
		ValAddress: valAddress,
		PubKey:     pubKey,
		StakeCoins: []sdk.Coin{coin},
		Status:     Unbonded,
		Commission: commission,
	}
}

// unmarshal a redelegation from a store value
func UnmarshalValidator(cdc *codec.Codec, value []byte) (validator Validator, err error) {
	err = cdc.UnmarshalBinaryLengthPrefixed(value, &validator)
	return validator, err
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
func (v Validator) ConsensusPower(power sdk.Int) int64 {
	if v.IsBonded() {
		return v.PotentialConsensusPower(power)
	}
	return 0
}

// potential consensus-engine power
func (v Validator) PotentialConsensusPower(power sdk.Int) int64 {
	return sdk.TokensToConsensusPower(power)
}

// ABCIValidatorUpdate returns an abci.ValidatorUpdate from a staking validator type
// with the full validator power
func (v Validator) ABCIValidatorUpdate(power sdk.Int) abci.ValidatorUpdate {
	return abci.ValidatorUpdate{
		PubKey: tmtypes.TM2PB.PubKey(v.PubKey),
		Power:  v.ConsensusPower(power),
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
