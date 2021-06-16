package types

import (
	"bitbucket.org/decimalteam/go-node/utils/errors"
	"encoding/hex"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"strconv"
)

// Local code type
type CodeType = uint32

const (
	// Default validator codespace
	DefaultCodespace string = ModuleName

	CodeSwapNotFound      = 100
	CodeSwapAlreadyExist  = 101
	CodeFromFieldNotEqual = 102
	CodeAlreadyRefunded   = 103
	CodeAlreadyRedeemed   = 104
	CodeNotExpired        = 105
	CodeExpired           = 106
	CodeWrongSecret       = 107
)

func ErrSwapNotFound() *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeSwapNotFound,
		`swap not found`,
	)
}

func ErrSwapAlreadyExist(hash Hash) *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeSwapAlreadyExist,
		fmt.Sprintf(`swap with hash %s already exist`, hex.EncodeToString(hash[:])),
		errors.NewParam("hash", hex.EncodeToString(hash[:])),
	)
}

func ErrFromFieldNotEqual(fromMsg, fromSwap sdk.AccAddress) *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeFromFieldNotEqual,
		fmt.Sprintf(`'from' field not equal: %s != %s`, fromMsg.String(), fromSwap.String()),
		errors.NewParam("fromMsg", fromMsg.String()),
		errors.NewParam("fromSwap", fromSwap.String()),
	)
}

func ErrAlreadyRefunded() *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeAlreadyRefunded,
		"start block must greater then current block height",
	)
}

func ErrAlreadyRedeemed() *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeAlreadyRedeemed,
		"already redeemed",
	)
}

func ErrNotExpired() *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeNotExpired,
		"swap not expired",
	)
}

func ErrExpired() *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeExpired,
		"swap expired",
	)
}

func ErrWrongSecret() *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeWrongSecret,
		"wrong secret",
	)
}

func getCodeString(code CodeType) string {
	return strconv.FormatInt(int64(code), 10)
}
