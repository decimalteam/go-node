package types

import (
	"encoding/json"
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
	jsonData, _ := json.Marshal(
		map[string]string{
			"code":      getCodeString(CodeInvalidSender),
			"codespace": DefaultCodespace,
			"desc":      "Invalid sender address: sender address cannot be empty",
		},
	)
	return sdkerrors.New(
		DefaultCodespace,
		CodeInvalidSender,
		string(jsonData),
	)
}

func ErrInvalidOwnerCount(more bool) *sdkerrors.Error {
	var AppendWord string = "need at least"
	ownerCount := MinOwnerCount
	if more {
		AppendWord = "allowed no more"
		ownerCount = MaxOwnerCount
	}
	jsonData, _ := json.Marshal(
		map[string]string{
			"code":       getCodeString(CodeInvalidOwnerCount),
			"codespace":  DefaultCodespace,
			"desc":       fmt.Sprintf("Invalid owner count: %s %d owners", AppendWord, ownerCount),
			"AppendWord": AppendWord,
		},
	)
	return sdkerrors.New(
		DefaultCodespace,
		CodeInvalidOwnerCount,
		string(jsonData),
	)
}

func ErrInvalidOwner() *sdkerrors.Error {
	jsonData, _ := json.Marshal(
		map[string]string{
			"code":      getCodeString(CodeInvalidOwner),
			"codespace": DefaultCodespace,
			"desc":      "Invalid owner address: owner address cannot be empty",
		},
	)
	return sdkerrors.New(
		DefaultCodespace,
		CodeInvalidOwner,
		string(jsonData),
	)
}

func ErrInvalidWeightCount(LenMsgWeights int, LenMsgOwners int) *sdkerrors.Error {
	jsonData, _ := json.Marshal(
		map[string]string{
			"code":          getCodeString(CodeInvalidWeightCount),
			"codespace":     DefaultCodespace,
			"desc":          fmt.Sprintf("Invalid weight count: weight count (%d) is not equal to owner count (%d)", LenMsgWeights, LenMsgOwners),
			"LenMsgWeights": fmt.Sprintf("%d", LenMsgWeights),
			"LenMsgOwners":  fmt.Sprintf("%d", LenMsgOwners),
		},
	)
	return sdkerrors.New(
		DefaultCodespace,
		CodeInvalidWeightCount,
		string(jsonData),
	)
}

func ErrInvalidWeight(weight int, data string) *sdkerrors.Error {
	jsonData, _ := json.Marshal(
		map[string]string{
			"code":      getCodeString(CodeInvalidWeight),
			"codespace": DefaultCodespace,
			"desc":      fmt.Sprintf("Invalid weight: weight cannot be %s than %d", data, weight),
			"weight":    fmt.Sprintf("%d", weight),
			"data":      data,
		},
	)
	return sdkerrors.New(
		DefaultCodespace,
		CodeInvalidWeight,
		string(jsonData),
	)
}
func ErrInvalidCoinToSend(denom string) *sdkerrors.Error {
	jsonData, _ := json.Marshal(
		map[string]string{
			"code":      getCodeString(CodeInvalidCoinToSend),
			"codespace": DefaultCodespace,
			"desc":      fmt.Sprintf("Coin to send with symbol %s does not exist", denom),
			"denom":     denom,
		},
	)
	return sdkerrors.New(
		DefaultCodespace,
		CodeInvalidCoinToSend,
		string(jsonData),
	)
}

func ErrWalletAccountNotFound() *sdkerrors.Error {
	jsonData, _ := json.Marshal(
		map[string]string{
			"code":      getCodeString(CodeWalletAccountNotFound),
			"codespace": DefaultCodespace,
			"desc":      "wallet account not found",
		},
	)
	return sdkerrors.New(
		DefaultCodespace,
		CodeWalletAccountNotFound,
		string(jsonData),
	)
}

func ErrInsufficientFunds(funds string) *sdkerrors.Error {
	jsonData, _ := json.Marshal(
		map[string]string{
			"code":      getCodeString(CodeInsufficientFunds),
			"codespace": DefaultCodespace,
			"desc":      fmt.Sprintf("Insufficient funds: wanted = %s", funds),
			"funds":     funds,
		},
	)
	return sdkerrors.New(
		DefaultCodespace,
		CodeInsufficientFunds,
		string(jsonData),
	)
}

func ErrDuplicateOwner(address sdk.AccAddress) *sdkerrors.Error {
	jsonData, _ := json.Marshal(
		map[string]string{
			"code":      getCodeString(CodeDuplicateOwner),
			"codespace": DefaultCodespace,
			"desc":      fmt.Sprintf("Invalid owners: owner with address %s is duplicated", address),
			"address":   address.String(),
		},
	)
	return sdkerrors.New(
		DefaultCodespace,
		CodeDuplicateOwner,
		string(jsonData),
	)
}

func getCodeString(code CodeType) string {
	return strconv.FormatInt(int64(code), 10)
}
