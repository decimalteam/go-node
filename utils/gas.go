package utils

import (
	"math"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type gasMeter struct {
	limit    sdk.Gas
	consumed sdk.Gas
}

// NewGasMeter returns a reference to a new gasMeter.
func NewGasMeter(limit sdk.Gas) sdk.GasMeter {
	return &gasMeter{
		limit:    limit,
		consumed: 0,
	}
}

func (g *gasMeter) GasConsumed() sdk.Gas {
	return g.consumed
}

func (g *gasMeter) Limit() sdk.Gas {
	return g.limit
}

func (g *gasMeter) GasConsumedToLimit() sdk.Gas {
	if g.IsPastLimit() {
		return g.limit
	}
	return g.consumed
}

// addUint64Overflow performs the addition operation on two uint64 integers and
// returns a boolean on whether or not the result overflows.
func addUint64Overflow(a, b uint64) (uint64, bool) {
	if math.MaxUint64-a < b {
		return 0, true
	}

	return a + b, false
}

func (g *gasMeter) ConsumeGas(amount sdk.Gas, descriptor string) {
	// fmt.Printf("####### Consume %d gas: %s\n", amount, descriptor)
	if descriptor != "commission" {
		return
	}
	var overflow bool
	// TODO: Should we set the consumed field after overflow checking?
	g.consumed, overflow = addUint64Overflow(g.consumed, amount)
	if overflow {
		panic(sdk.ErrorGasOverflow{Descriptor: descriptor})
	}

	if g.consumed > g.limit {
		panic(sdk.ErrorOutOfGas{Descriptor: descriptor})
	}

}

func (g *gasMeter) IsPastLimit() bool {
	return g.consumed > g.limit
}

func (g *gasMeter) IsOutOfGas() bool {
	return g.consumed >= g.limit
}
