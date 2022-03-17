package helpers

import (
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"time"

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
func JoinAccAddresses(values []sdk.AccAddress) string {
	var sb strings.Builder
	for i, v := range values {
		if i > 0 {
			sb.WriteString(",")
		}
		sb.WriteString(v.String())
	}
	return sb.String()
}

// JoinUints returns string containing all provided uint values joined with ",".
func JoinUints(values []uint) string {
	var sb strings.Builder
	for i, v := range values {
		if i > 0 {
			sb.WriteString(",")
		}
		sb.WriteString(strconv.FormatUint(uint64(v), 10))
	}
	return sb.String()
}

func TimeTrack(ctx sdk.Context, msg string) (sdk.Context, string, time.Time) {
	return ctx, msg, time.Now()
}

func TimeDuration(ctx sdk.Context, msg string, start time.Time) {
	ctx.Logger().Info(fmt.Sprintf("%v: %v ms", msg, time.Since(start).Milliseconds()))
}
