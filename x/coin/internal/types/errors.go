package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Local code type
type CodeType = sdk.CodeType

const (
	// Default coin codespace
	DefaultCodespace sdk.CodespaceType = ModuleName

	DecodeError            CodeType = 101
	InvalidCRR             CodeType = 102
	InvalidCoinSymbol      CodeType = 103
	CoinAlreadyExists      CodeType = 104
	InvalidCoinTitle       CodeType = 105
	InvalidCoinInitAmount  CodeType = 106
	InvalidCoinInitReserve CodeType = 107
	CodeInvalid            CodeType = 108
)
