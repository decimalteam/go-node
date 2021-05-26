package types

import (
	"github.com/cosmos/cosmos-sdk/types"
	"math/big"
	"time"
)

// Staking params default values
const (
	// DefaultUnbondingTime reflects three weeks in seconds as the default
	// unbonding time.
	DefaultUnbondingTime = time.Hour * 24 * 30

	// Default maximum number of bonded validators
	DefaultMaxValidators uint16 = 256

	// Default maximum entries in a UBD/RED pair
	DefaultMaxEntries uint16 = 7

	// DefaultHistorical entries is 0 since it must only be non-zero for
	// IBC connected chains
	DefaultHistoricalEntries uint16 = 0

	DefaultBondDenom string = "del"

	DefaultMaxDelegations uint16 = 1000
)

// PowerReduction is the amount of staking tokens required for 1 unit of consensus-engine power
var PowerReduction = types.NewIntFromBigInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil))

// TokensToConsensusPower - convert input tokens to potential consensus-engine power
func TokensToConsensusPower(tokens types.Int) int64 {
	return (tokens.Quo(PowerReduction)).Int64()
}

// TokensFromConsensusPower - convert input power to tokens
func TokensFromConsensusPower(power int64) types.Int {
	return types.NewInt(power).Mul(PowerReduction)
}
