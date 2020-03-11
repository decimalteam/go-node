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
)

func ErrEmptyValidatorAddr(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeEmptyValidatorAddr, "empty validator address")
}
