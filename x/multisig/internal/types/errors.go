package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// CodeType defines the local code type.
type CodeType = uint32

// DefaultCodespace defines default multisig codespace.
const DefaultCodespace string = ModuleName

// Custom errors codes.
const (
	CodeInvalidSender         CodeType = 101
	CodeInvalidOwnerCount     CodeType = 102
	CodeInvalidOwner          CodeType = 103
	CodeInvalidWeightCount    CodeType = 104
	CodeInvalidWeight         CodeType = 105
	CodeInvalidCoinToSend     CodeType = 106
	CodeInvalidAmountToSend   CodeType = 107
	CodeWalletAccountNotFound CodeType = 108
	CodeInsufficientFunds     CodeType = 109
	CodeDuplicateOwner        CodeType = 110
)

func ErrInvalidSender() *sdkerrors.Error {
	return sdkerrors.New(
		DefaultCodespace,
		CodeInvalidSender,
		"Invalid sender address: sender address cannot be empty",
	)
}

func ErrInvalidOwnerCount(count int, more bool) *sdkerrors.Error {
	var AppendWord string
	if more {
		AppendWord = "allowed no more"
	} else {
		AppendWord = "need at least"
	}
	return sdkerrors.New(
		DefaultCodespace,
		CodeInvalidOwnerCount,
		fmt.Sprintf("Invalid owner count: %s %d owners", AppendWord, MinOwnerCount),
	)
}

func ErrInvalidOwner() *sdkerrors.Error {
	return sdkerrors.New(
		DefaultCodespace,
		CodeInvalidOwner,
		"Invalid owner address: owner address cannot be empty",
	)
}

func ErrInvalidWeightCount(LenMsgWeights int, LenMsgOwners int) *sdkerrors.Error {
	return sdkerrors.New(
		DefaultCodespace,
		CodeInvalidWeightCount,
		fmt.Sprintf("Invalid weight count: weight count (%d) is not equal to owner count (%d)", LenMsgWeights, LenMsgOwners),
	)
}

func ErrInvalidWeight(weight int, data string) *sdkerrors.Error {
	return sdkerrors.New(
		DefaultCodespace,
		CodeInvalidWeight,
		fmt.Sprintf("Invalid weight: weight cannot be %s than %d", data, weight),
	)
}
func ErrInvalidCoinToSend(denom string) *sdkerrors.Error {
	return sdkerrors.New(
		DefaultCodespace,
		CodeInvalidCoinToSend,
		fmt.Sprintf("Coin to send with symbol %s does not exist", denom),
	)
}

func ErrWalletAccountNotFound() *sdkerrors.Error {
	return sdkerrors.New(
		DefaultCodespace,
		CodeWalletAccountNotFound,
		"wallet account not found",
	)
}

func ErrInsufficientFunds(funds string) *sdkerrors.Error {
	return sdkerrors.New(
		DefaultCodespace,
		CodeInsufficientFunds,
		fmt.Sprintf("Insufficient funds: wanted = %s", funds),
	)
}

func ErrDuplicateOwner(address sdk.AccAddress) *sdkerrors.Error {
	return sdkerrors.New(
		DefaultCodespace,
		CodeDuplicateOwner,
		fmt.Sprintf("Invalid owners: owner with address %s is duplicated", address),
	)
}
