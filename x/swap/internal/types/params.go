package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	"time"
)

const (
	DefaultLockedTimeOut = time.Hour * 24
	DefaultLockedTimeIn  = time.Hour * 12
)

const ServiceAddress = "dx1p844kydt9eljvuef4nk52dm6lcgj5c42q4zmvd"
const ChainActivatorAddress = "dx16aeq4ypsx5ar4076v507ch5z8ryd6usx32tnru"
const CheckingAddress = "d2d9207a88982ecffec424709ff2b02f6c95a9ba"

var ServiceAccAddress, _ = sdk.AccAddressFromBech32("dx1jqx7chw0faswfmw78cdejzzery5akzmk5zc5x5")

var (
	KeyLockedTimeOut = []byte("LockedTimeOut")
	KeyLockedTimeIn  = []byte("LockedTimeIn")
)

type Params struct {
	LockedTimeOut time.Duration `json:"locked_time_out"`
	LockedTimeIn  time.Duration `json:"locked_time_in"`
}

func NewParams(lockedTimeOut, lockedTimeIn time.Duration) Params {
	return Params{
		LockedTimeOut: lockedTimeOut,
		LockedTimeIn:  lockedTimeIn,
	}
}

// Implements params.ParamSet
func (p *Params) ParamSetPairs() params.ParamSetPairs {
	return params.ParamSetPairs{
		params.NewParamSetPair(KeyLockedTimeOut, &p.LockedTimeOut, validateLockedTime),
		params.NewParamSetPair(KeyLockedTimeIn, &p.LockedTimeIn, validateLockedTime),
	}
}

func validateLockedTime(i interface{}) error {
	v, ok := i.(time.Duration)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v <= 0 {
		return fmt.Errorf("locked time must be positive: %d", v)
	}

	return nil
}

func DefaultParams() Params {
	return NewParams(DefaultLockedTimeOut, DefaultLockedTimeIn)
}

func (p Params) Validate() error {
	return nil
}
