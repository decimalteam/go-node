package types

import (
	"fmt"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// CodeType defines the local code type.
type CodeType = uint32

// DefaultCodespace defines default multisig codespace.
const DefaultCodespace string = ModuleName

// Custom errors codes.
const (
	InvalidSender             CodeType = 101
	InvalidOwnerCount         CodeType = 102
	InvalidOwner              CodeType = 103
	InvalidWeightCount        CodeType = 104
	InvalidWeight             CodeType = 105
	InvalidCoinToSend         CodeType = 106
	InvalidAmountToSend       CodeType = 107
	CodeWalletAccountNotFound CodeType = 108
	CodeInsufficientFunds     CodeType = 109
)

func ErrWalletAccountNotFound() *sdkerrors.Error {
	return sdkerrors.New(DefaultCodespace, CodeWalletAccountNotFound, "wallet account not found")
}

func ErrInsufficientFunds(funds string) *sdkerrors.Error {
	return sdkerrors.New(DefaultCodespace, CodeInsufficientFunds, fmt.Sprintf("Insufficient funds: wanted = %s", funds))
}
