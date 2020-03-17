package types

import (
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

func ErrValidatorPubKeyExists(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidValidator, "validator already exist for this pubkey, must use new validator pubkey")
}
