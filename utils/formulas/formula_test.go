package formulas

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"log"
	"testing"
)

func TestCalculateSaleReturn(t *testing.T) {
	supply, _ := sdk.NewIntFromString("1000000000000000000000")
	reserve, _ := sdk.NewIntFromString("1000000000000000000000")
	oneCoin, _ := sdk.NewIntFromString("1000000000000000000")
	amount, _ := sdk.NewIntFromString("10000000000")

	full := CalculateSaleReturn(supply, reserve, 50, amount)

	price := CalculateSaleReturn(supply, reserve, 50, oneCoin)

	short := amount.Mul(price).Quo(oneCoin)

	log.Println(price, full, short)
}

func BenchmarkCalculateSaleReturn(b *testing.B) {
	supply, _ := sdk.NewIntFromString("1000000000000000000000")
	reserve, _ := sdk.NewIntFromString("1000000000000000000000")
	amount, _ := sdk.NewIntFromString("10000000000")
	res := sdk.ZeroInt()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		res = CalculateSaleReturn(supply, reserve, 50, amount)
	}
	fmt.Println(res)
}

func TestCalculateSaleAmount(t *testing.T) {
	supply, _ := sdk.NewIntFromString("1030421762542736433327868")
	reserve, _ := sdk.NewIntFromString("1044216001082372924902")
	amount, _ := sdk.NewIntFromString("1000000000000000000")

	result := CalculateSaleAmount(supply, reserve, 50, amount)

	t.Log(result)
}

func TestCalculateSaleReturn2(t *testing.T) {
	supply, _ := sdk.NewIntFromString("1030421762542736433327868")
	reserve, _ := sdk.NewIntFromString("1044216001082372924902")
	amount, _ := sdk.NewIntFromString("3380759108999599642")

	result := CalculateSaleReturn(supply, reserve, 50, amount)

	t.Log(result)
}
