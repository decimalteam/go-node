package types

import (
	"encoding/hex"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// Local code type
type CodeType = uint32

const (
	// Default validator codespace
	DefaultCodespace string = ModuleName

	CodeSwapNotFound      = 100
	CodeSwapAlreadyExist  = 101
	CodeFromFieldNotEqual = 102
	CodeAlreadyRedeem     = 103
)

func ErrSwapNotFound() *sdkerrors.Error {
	return sdkerrors.New(DefaultCodespace, CodeSwapNotFound, `swap not found`)
}

func ErrSwapAlreadyExist(hash Hash) *sdkerrors.Error {
	return sdkerrors.New(DefaultCodespace, CodeSwapAlreadyExist, fmt.Sprintf(`swap with hash %s already exist`, hex.EncodeToString(hash[:])))
}

func ErrFromFieldNotEqual(fromMsg, fromSwap sdk.AccAddress) *sdkerrors.Error {
	return sdkerrors.New(DefaultCodespace, CodeFromFieldNotEqual, fmt.Sprintf(`'from' field not equal: %s != %s`, fromMsg.String(), fromSwap.String()))
}

func ErrAlreadyRedeem() *sdkerrors.Error {
	return sdkerrors.New(DefaultCodespace, CodeAlreadyRedeem, "already redeem")
}
