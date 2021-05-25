package types

import (
	"encoding/hex"
	"encoding/json"
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
	jsonData, _ := json.Marshal(
		map[string]string{
			"code":      getCodeString(CodeSwapNotFound),
			"codespace": DefaultCodespace,
			"desc":      `swap not found`,
		},
	)
	return sdkerrors.New(
		DefaultCodespace,
		CodeSwapNotFound,
		string(jsonData),
	)
}

func ErrSwapAlreadyExist(hash Hash) *sdkerrors.Error {
	jsonData, _ := json.Marshal(
		map[string]string{
			"code":      getCodeString(CodeSwapAlreadyExist),
			"codespace": DefaultCodespace,
			"desc":      fmt.Sprintf(`swap with hash %s already exist`, hex.EncodeToString(hash[:])),
			"hash":      hex.EncodeToString(hash[:]),
		},
	)
	return sdkerrors.New(
		DefaultCodespace,
		CodeSwapAlreadyExist,
		string(jsonData),
	)
}

func ErrFromFieldNotEqual(fromMsg, fromSwap sdk.AccAddress) *sdkerrors.Error {
	jsonData, _ := json.Marshal(
		map[string]string{
			"code":      getCodeString(CodeFromFieldNotEqual),
			"codespace": DefaultCodespace,
			"desc":      fmt.Sprintf(`'from' field not equal: %s != %s`, fromMsg.String(), fromSwap.String()),
			"fromMsg":   fromMsg.String(),
			"fromSwap":  fromSwap.String(),
		},
	)
	return sdkerrors.New(
		DefaultCodespace,
		CodeFromFieldNotEqual,
		string(jsonData),
	)
}

func ErrAlreadyRefunded() *sdkerrors.Error {
	jsonData, _ := json.Marshal(
		map[string]string{
			"code":      getCodeString(CodeAlreadyRefunded),
			"codespace": DefaultCodespace,
			"desc":      "start block must greater then current block height",
		},
	)
	return sdkerrors.New(
		DefaultCodespace,
		CodeAlreadyRefunded,
		string(jsonData),
	)
}

func ErrAlreadyRedeemed() *sdkerrors.Error {
	jsonData, _ := json.Marshal(
		map[string]string{
			"code":      getCodeString(CodeAlreadyRedeemed),
			"codespace": DefaultCodespace,
			"desc":      "already redeemed",
		},
	)
	return sdkerrors.New(
		DefaultCodespace,
		CodeAlreadyRedeemed,
		string(jsonData),
	)
}

func ErrNotExpired() *sdkerrors.Error {
	jsonData, _ := json.Marshal(
		map[string]string{
			"code":      getCodeString(CodeNotExpired),
			"codespace": DefaultCodespace,
			"desc":      "swap not expired",
		},
	)
	return sdkerrors.New(
		DefaultCodespace,
		CodeNotExpired,
		string(jsonData),
	)
}

func ErrExpired() *sdkerrors.Error {
	jsonData, _ := json.Marshal(
		map[string]string{
			"code":      getCodeString(CodeExpired),
			"codespace": DefaultCodespace,
			"desc":      "swap expired",
		},
	)
	return sdkerrors.New(
		DefaultCodespace,
		CodeExpired,
		string(jsonData),
	)
}

func ErrWrongSecret() *sdkerrors.Error {
	jsonData, _ := json.Marshal(
		map[string]string{
			"code":      getCodeString(CodeWrongSecret),
			"codespace": DefaultCodespace,
			"desc":      "wrong secret",
		},
	)
	return sdkerrors.New(
		DefaultCodespace,
		CodeWrongSecret,
		string(jsonData),
	)
}

func getCodeString(code CodeType) string {
	return strconv.FormatInt(int64(code), 10)
}
