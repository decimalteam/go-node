package formulas

import (
	"bitbucket.org/decimalteam/go-node/utils/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"math/big"
)

const (
	precision = 100
)

func newFloat(x float64) *big.Float {
	return big.NewFloat(x).SetPrec(precision)
}

func newIntPrec(x sdk.Int) sdk.Int {
	return sdk.NewInt(x.Int64()).Mul(sdk.NewInt(precision))
}

// Return = supply * ((1 + deposit / reserve) ^ (crr / 100) - 1)
// Рассчитывает сколько монет мы получим заплатив deposit DLC (Покупка формула 2)
func CalculatePurchaseReturn(supply sdk.Int, reserve sdk.Int, crr uint, deposit sdk.Int) sdk.Int {
	if deposit.Equal(sdk.NewInt(0)) {
		return sdk.NewInt(0)
	}

	if crr == 100 {
		result := sdk.NewInt(1).Mul(supply).Mul(deposit)

		return result.Quo(reserve)
	}

	//tSupply1 := newIntPrec(supply)
	//tReserve1 := newIntPrec(reserve)
	//tDeposit1 := newIntPrec(deposit)
	//
	//res1 := tDeposit1.Quo(reserve)                           // deposit * precision / reserve
	//res1 = res1.Add(newIntPrec(sdk.NewInt(1)))            // 1 * precision  + (deposit * precision / reserve)
	//res1 = math.PowInt(res1, sdk.NewInt(int64(crr / 100)))   // (1 * precision + (deposit * precision / reserve)) ^ (crr / 100)
	//res1 = res1.Add(res1, newFloat(1))                       // ((1 + deposit / reserve) ^ (crr / 100) - 1)
	//res.Mul(res, tSupply)                           // supply * ((1 + deposit / reserve) ^ (crr / 100) - 1)
	//----------
	tSupply := newFloat(0).SetInt(supply.BigInt())
	tReserve := newFloat(0).SetInt(reserve.BigInt())
	tDeposit := newFloat(0).SetInt(deposit.BigInt())

	res := newFloat(0).Quo(tDeposit, tReserve)      // deposit / reserve
	res.Add(res, newFloat(1))                       // 1 + (deposit / reserve)
	res = math.Pow(res, newFloat(float64(crr)/100)) // (1 + deposit / reserve) ^ (crr / 100)
	res.Sub(res, newFloat(1))                       // ((1 + deposit / reserve) ^ (crr / 100) - 1)
	res.Mul(res, tSupply)                           // supply * ((1 + deposit / reserve) ^ (crr / 100) - 1)

	converted, _ := res.Int(nil)
	result := sdk.NewIntFromBigInt(converted)

	return result
}

// reversed function CalculatePurchaseReturn
// deposit = reserve * (((wantReceive + supply) / supply)^(100 / crr) - 1)
// Рассчитывает сколько DLC надо заплатить , чтобы получить wantReceive монет (Покупка)

func CalculatePurchaseAmount(supply sdk.Int, reserve sdk.Int, crr uint, wantReceive sdk.Int) sdk.Int {
	if wantReceive.Equal(sdk.NewInt(0)) {
		return sdk.NewInt(0)
	}

	if crr == 100 {
		result := sdk.NewInt(1).Mul(wantReceive).Mul(reserve)

		return result.Quo(supply)
	}

	tSupply := newFloat(0).SetInt(supply.BigInt())
	tReserve := newFloat(0).SetInt(reserve.BigInt())
	tWantReceive := newFloat(0).SetInt(wantReceive.BigInt())

	res := newFloat(0).Add(tWantReceive, tSupply)   // reserve + supply
	res.Quo(res, tSupply)                           // (reserve + supply) / supply
	res = math.Pow(res, newFloat(100/float64(crr))) // ((reserve + supply) / supply)^(100/c)
	res.Sub(res, newFloat(1))                       // (((reserve + supply) / supply)^(100/c) - 1)
	res.Mul(res, tReserve)                          // reserve * (((reserve + supply) / supply)^(100/c) - 1)

	converted, _ := res.Int(nil)
	result := sdk.NewIntFromBigInt(converted)

	return result
}

// Return = reserve * (1 - (1 - sellAmount / supply) ^ (100 / crr))
// Рассчитывает сколько DCL вы получите, если продадите sellAmount монет. (Продажа)
func CalculateSaleReturn(supply sdk.Int, reserve sdk.Int, crr uint, sellAmount sdk.Int) sdk.Int {
	// special case for 0 sell amount
	if sellAmount.Equal(sdk.NewInt(0)) {
		return sdk.NewInt(0)
	}

	// special case for selling the entire supply
	if sellAmount.Equal(supply) {
		return reserve
	}

	if crr == 100 {
		ret := sdk.NewInt(1).Mul(reserve).Mul(sellAmount)
		ret.Quo(supply)

		return ret
	}

	tSupply := newFloat(0).SetInt(supply.BigInt())
	tReserve := newFloat(0).SetInt(reserve.BigInt())
	tSellAmount := newFloat(0).SetInt(sellAmount.BigInt())

	res := newFloat(0).Quo(tSellAmount, tSupply)      // sellAmount / supply
	res.Sub(newFloat(1), res)                         // (1 - sellAmount / supply)
	res = math.Pow(res, newFloat(100/(float64(crr)))) // (1 - sellAmount / supply) ^ (100 / crr)
	res.Sub(newFloat(1), res)                         // (1 - (1 - sellAmount / supply) ^ (1 / (crr / 100)))
	res.Mul(res, tReserve)                            // reserve * (1 - (1 - sellAmount / supply) ^ (1 / (crr / 100)))

	converted, _ := res.Int(nil)
	result := sdk.NewIntFromBigInt(converted)

	return result
}

// reversed function CalculateSaleReturn
// -(-1 + (-(wantReceive - reserve)/reserve)^(crr / 100)) * supply
// Рассчитывает сколько монет надо продать, чтобы получить wantReceive DCL. (Продажа 2)

func CalculateSaleAmount(supply sdk.Int, reserve sdk.Int, crr uint, wantReceive sdk.Int) sdk.Int {
	if wantReceive.Equal(sdk.NewInt(0)) {
		return sdk.NewInt(0)
	}

	if crr == 100 {
		ret := sdk.NewInt(1).Mul(wantReceive).Mul(supply)
		ret.Quo(reserve)

		return ret
	}

	tSupply := newFloat(0).SetInt(supply.BigInt())
	tReserve := newFloat(0).SetInt(reserve.BigInt())
	tWantReceive := newFloat(0).SetInt(wantReceive.BigInt())

	res := newFloat(0).Sub(tWantReceive, tReserve)  // (wantReceive - reserve)
	res.Neg(res)                                    // -(wantReceive - reserve)
	res.Quo(res, tReserve)                          // -(wantReceive - reserve)/reserve
	res = math.Pow(res, newFloat(float64(crr)/100)) // (-(wantReceive - reserve)/reserve)^(crr/100)
	res.Add(res, newFloat(-1))                      // -1 + (-(wantReceive - reserve)/reserve)^(crr/100)
	res.Neg(res)                                    // -(-1 + (-(wantReceive - reserve)/reserve)^(crr/100))
	res.Mul(res, tSupply)                           // -(-1 + (-(wantReceive - reserve)/reserve)^(crr/100)) * supply

	converted, _ := res.Int(nil)
	result := sdk.NewIntFromBigInt(converted)

	return result
}
