package types

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto"
	tmtypes "github.com/tendermint/tendermint/types"
	"strings"
	"time"
)

type Validator struct {
	ValAddress              sdk.ValAddress `json:"val_address"`
	PubKey                  crypto.PubKey  `json:"pub_key"`
	Tokens                  sdk.Int        `json:"stake_coins"`
	DelegatorShares         sdk.Dec        `json:"delegator_shares"`
	Status                  BondStatus     `json:"status"`
	Commission              Commission     `json:"commission"`
	Jailed                  bool           `json:"jailed"`
	UnbondingCompletionTime time.Time      `json:"unbonding_completion_time"`
	UnbondingHeight         int64          `json:"unbonding_height"`
	Description             Description    `json:"description"`
	AccumRewards            sdk.Int        `json:"accum_rewards"`
	DelegatorStakes         []Stake        `json:"delegator_stakes"`
}

type Stake struct {
	Delegator sdk.AccAddress `json:"delegator"`
	Coin      sdk.Coin       `json:"coin"`
}

// String returns a human readable string representation of a validator.
func (v Validator) String() string {
	bechConsPubKey, err := sdk.Bech32ifyConsPub(v.PubKey)
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf(`Validator
  Operator Address:           %s
  Validator Consensus Pubkey: %s
  Jailed:                     %v
  Status:                     %s
  Tokens:                     %s
  Delegator Shares:           %s
  Description:                %s
  Unbonding Height:           %d
  Unbonding Completion Time:  %v
  Commission:                 %s`, v.ValAddress, bechConsPubKey,
		v.Jailed, v.Status, v.Tokens,
		v.DelegatorShares, v.Description,
		v.UnbondingHeight, v.UnbondingCompletionTime, v.Commission)
}

func (v Validator) SharesFromTokens(tokens sdk.Int, valTokens sdk.Int, delTokens sdk.Dec) (sdk.Dec, sdk.Error) {
	if v.Tokens.IsZero() {
		return sdk.Dec{}, ErrInsufficientShares(DefaultCodespace)
	}

	return delTokens.MulInt(tokens).QuoInt(valTokens), nil
}

func (v Validator) IsJailed() bool               { return v.Jailed }
func (v Validator) GetMoniker() string           { return v.Description.Moniker }
func (v Validator) GetStatus() sdk.BondStatus    { return sdk.BondStatus(v.Status) }
func (v Validator) GetOperator() sdk.ValAddress  { return v.ValAddress }
func (v Validator) GetConsPubKey() crypto.PubKey { return v.PubKey }
func (v Validator) GetConsAddr() sdk.ConsAddress { return sdk.ConsAddress(v.PubKey.Address()) }
func (v Validator) GetTokens() sdk.Int           { return v.Tokens }
func (v Validator) GetBondedTokens() sdk.Int     { return v.BondedTokens() }
func (v Validator) GetCommission() sdk.Dec       { return v.Commission.Rate }
func (v Validator) GetDelegatorShares() sdk.Dec  { return v.DelegatorShares }

type Validators []Validator

func (v Validators) String() string {
	var out string
	for _, val := range v {
		out += val.String() + "\n"
	}
	return strings.TrimSpace(out)
}

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

func NewValidator(valAddress sdk.ValAddress, pubKey crypto.PubKey, commission Commission) Validator {
	return Validator{
		ValAddress:      valAddress,
		PubKey:          pubKey,
		Tokens:          sdk.ZeroInt(),
		Status:          Unbonded,
		Commission:      commission,
		DelegatorShares: sdk.ZeroDec(),
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

// get the bonded tokens which the validator holds
func (v Validator) BondedTokens() sdk.Int {
	if v.IsBonded() {
		return v.Tokens
	}
	return sdk.ZeroInt()
}

// AddTokensFromDel adds tokens to a validator
func (v Validator) AddTokensFromDel(token sdk.Coin, totalVal sdk.Int) (Validator, sdk.Dec) {

	// calculate the shares to issue
	issuedShares := sdk.ZeroDec()
	if v.DelegatorShares.IsZero() {
		// the first delegation to a validator sets the exchange rate to one
		issuedShares = token.Amount.ToDec()
	} else {
		shares, err := v.SharesFromTokens(token.Amount, totalVal, v.DelegatorShares)
		if err != nil {
			panic(err)
		}

		issuedShares = shares
	}

	v.Tokens = v.Tokens.Add(token.Amount)
	v.DelegatorShares = v.DelegatorShares.Add(issuedShares)

	return v, issuedShares
}

// In some situations, the exchange rate becomes invalid, e.g. if
// Validator loses all tokens due to slashing. In this case,
// make all future delegations invalid.
func (v Validator) InvalidExRate() bool {
	return v.Tokens.IsZero() && v.DelegatorShares.IsPositive()
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
