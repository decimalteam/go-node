package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto"
)

type Validator struct {
	ValAddress sdk.ValAddress `json:"val_address"`
	PubKey     crypto.PubKey  `json:"pub_key"`
	StakeCoins []sdk.Coin     `json:"stake_coins"`
	Status     BondStatus     `json:"status"`
	Commission Commission     `json:"commission"`
	Jailed     bool           `json:"jailed"`
}

type Validators []Validator

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
	return v.GetStatus().Equal(Bonded)
}

// IsUnbonded checks if the validator status equals Unbonded
func (v Validator) IsUnbonded() bool {
	return v.GetStatus().Equal(Unbonded)
}

// IsUnbonding checks if the validator status equals Unbonding
func (v Validator) IsUnbonding() bool {
	return v.GetStatus().Equal(Unbonding)
}

func (v Validator) GetStatus() BondStatus { return v.Status }
