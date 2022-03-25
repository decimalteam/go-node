package types

import (
	"fmt"

	"bitbucket.org/decimalteam/go-node/utils/errors"
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
	CodeAlreadyRefunded   = 103
	CodeAlreadyRedeemed   = 104
	CodeNotExpired        = 105
	CodeExpired           = 106
	CodeWrongSecret       = 107

	CodeChainNotExist            = 200
	CodeInvalidServiceAddress    = 201
	CodeInsufficientPoolFunds    = 202
	CodeInvalidTransactionNumber = 203

	CodeDeprecated = 300
)

func ErrSwapNotFound() *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeSwapNotFound,
		`swap not found`,
	)
}

func ErrSwapAlreadyExist(hash string) *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeSwapAlreadyExist,
		fmt.Sprintf(`swap with hash %s already exist`, hash),
	)
}

func ErrFromFieldNotEqual(fromMsg string, fromSwap string) *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeFromFieldNotEqual,
		fmt.Sprintf(`'from' field not equal: %s != %s`, fromMsg, fromSwap),
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

func ErrChainNotExist(chain string) *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeChainNotExist,
		fmt.Sprintf("chain %s does not exist", chain),
	)
}

func ErrInvalidServiceAddress(want string, receive string) *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeInvalidServiceAddress,
		fmt.Sprintf("invalid service address: want = %s, receive = %s", want, receive),
	)
}

func ErrInsufficientPoolFunds(want string, exists string) *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeInsufficientPoolFunds,
		fmt.Sprintf("insufficient pool funds: want = %s, exists = %s", want, exists),
	)
}

func ErrInvalidTransactionNumber() *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeInvalidTransactionNumber,
		"invalid transaction number",
	)
}

func ErrDeprecated() *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeDeprecated,
		"msg deprecated",
	)
}
