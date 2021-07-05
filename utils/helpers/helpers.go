package helpers

import (
	"math/big"
	"strconv"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func BipToPip(bip sdk.Int) sdk.Int {
	return bip.Mul(sdk.NewIntFromBigInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)))
}

func UnitToPip(unit sdk.Int) sdk.Int {
	return unit.Mul(sdk.NewIntFromBigInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(15), nil)))
}

func PipToUnit(pip sdk.Int) sdk.Int {
	return pip.Quo(sdk.NewIntFromBigInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(15), nil)))
}

// JoinAccAddresses returns string containing all provided address joined with ",".
func JoinAccAddresses(values []string) string {
	var sb strings.Builder
	for i, v := range values {
		if i > 0 {
			sb.WriteString(",")
		}
		sb.WriteString(v)
	}
	return sb.String()
}

// JoinUints returns string containing all provided uint values joined with ",".
func JoinUints(values []uint64) string {
	var sb strings.Builder
	for i, v := range values {
		if i > 0 {
			sb.WriteString(",")
		}
		sb.WriteString(strconv.FormatUint(v, 10))
	}
	return sb.String()
}
