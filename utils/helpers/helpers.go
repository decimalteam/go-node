package helpers

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"math/big"
)

func BipToPip(bip sdk.Int) sdk.Int {
	return bip.Mul(sdk.NewIntFromBigInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)))
}
