package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	"time"
)

const (
	DefaultLockedTime = time.Hour * 12
)

var SwapServiceAddress, _ = sdk.AccAddressFromBech32("dx1jqx7chw0faswfmw78cdejzzery5akzmk5zc5x5")

var (
	KeyLockedTime = []byte("LockedTime")
)

type Params struct {
	LockedTime time.Duration `json:"locked_time"`
}

func NewParams(lockedTime time.Duration) Params {
	return Params{LockedTime: lockedTime}
}

// Implements params.ParamSet
func (p *Params) ParamSetPairs() params.ParamSetPairs {
	return params.ParamSetPairs{
		params.NewParamSetPair(KeyLockedTime, &p.LockedTime, validateLockedTime),
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
	return NewParams(DefaultLockedTime)
}

func (p Params) Validate() error {
	return nil
}
