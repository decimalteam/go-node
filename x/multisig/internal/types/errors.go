package types

import (
	"bitbucket.org/decimalteam/go-node/utils/errors"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"strconv"
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
	return errors.Encode(
		DefaultCodespace,
		CodeInvalidSender,
		fmt.Sprintf("Invalid sender address: sender address cannot be empty"),
	)
}

func ErrInvalidOwnerCount(more bool) *sdkerrors.Error {
	var AppendWord string = "need at least"
	ownerCount := fmt.Sprintf("%d", MinOwnerCount)
	if more {
		AppendWord = "allowed no more"
		ownerCount = fmt.Sprintf("%d", MaxOwnerCount)
	}

	return errors.Encode(
		DefaultCodespace,
		CodeInvalidOwnerCount,
		fmt.Sprintf("Invalid owner count: %s %s owners", AppendWord, ownerCount),
		errors.NewParam("AppendWord", AppendWord),
		errors.NewParam("ownerCount", ownerCount),
	)
}

func ErrInvalidOwner() *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeInvalidOwner,
		"Invalid owner address: owner address cannot be empty",
	)
}

func ErrInvalidWeightCount(LenMsgWeights string, LenMsgOwners string) *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeInvalidWeightCount,
		fmt.Sprintf("Invalid weight count: weight count (%s) is not equal to owner count (%s)", LenMsgWeights, LenMsgOwners),
		errors.NewParam("LenMsgWeights", LenMsgWeights),
		errors.NewParam("LenMsgOwners", LenMsgOwners),
	)
}

func ErrInvalidWeight(weight string, data string) *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeInvalidWeight,
		fmt.Sprintf("Invalid weight: weight cannot be %s than %s", data, weight),
		errors.NewParam("data", data),
		errors.NewParam("weight", weight),
	)
}
func ErrInvalidCoinToSend(denom string) *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeInvalidCoinToSend,
		fmt.Sprintf("Coin to send with symbol %s does not exist", denom),
		errors.NewParam("denom", denom),
	)
}

func ErrWalletAccountNotFound() *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeWalletAccountNotFound,
		fmt.Sprintf("wallet account not found"),
	)
}

func ErrInsufficientFunds(funds string) *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeInsufficientFunds,
		fmt.Sprintf("Insufficient funds: wanted = %s", funds),
		errors.NewParam("funds", funds),
	)
}

func ErrDuplicateOwner(address sdk.AccAddress) *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeDuplicateOwner,
		fmt.Sprintf("Invalid owners: owner with address %s is duplicated", address),
		errors.NewParam("address", address.String()),
	)
}

func getCodeString(code CodeType) string {
	return strconv.FormatInt(int64(code), 10)
}
