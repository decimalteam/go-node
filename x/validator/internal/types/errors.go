package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Local code type
type CodeType = sdk.CodeType

const (
	// Default validator codespace
	DefaultCodespace sdk.CodespaceType = ModuleName

	CodeEmptyValidatorAddr CodeType = 101
	CodeInvalidValidator   CodeType = 102

	CodeInvalidStruct CodeType = 201
	CodeInvalidInput  CodeType = 202

	CodeInvalidDelegation = 301
)

func ErrEmptyValidatorAddr(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeEmptyValidatorAddr, "empty validator address")
}

func ErrInvalidStruct(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidStruct, "invalid struct")
}

func ErrValidatorOwnerExists(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidValidator, "validator already exist for this operator address, must use new validator operator address")
}

func ErrNoValidatorFound(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidValidator, "validator does not exist for that address")
}

func ErrValidatorPubKeyExists(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidValidator, "validator already exist for this pubkey, must use new validator pubkey")
}

func ErrInsufficientShares(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidDelegation, "insufficient delegation shares")
}

func ErrDelegatorShareExRateInvalid(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidDelegation,
		"cannot delegate to validators with invalid (zero) ex-rate")
}

func ErrBadDelegationAddr(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidInput, "unexpected address length for this (address, validator) pair")
}

func ErrNoDelegatorForAddress(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidDelegation, "delegator does not contain this delegation")
}

func ErrNotEnoughDelegationShares(codespace sdk.CodespaceType, shares string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidDelegation, fmt.Sprintf("not enough shares only have %v", shares))
}

func ErrNilDelegatorAddr(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidInput, "delegator address is nil")
}
