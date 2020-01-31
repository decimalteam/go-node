package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Local code type
type CodeType = sdk.CodeType

const (
	// Default coin codespace
	DefaultCodespace sdk.CodespaceType = ModuleName
	// Create coin
	DecodeError            CodeType = 101
	InvalidCRR             CodeType = 102
	InvalidCoinSymbol      CodeType = 103
	CoinAlreadyExists      CodeType = 104
	InvalidCoinTitle       CodeType = 105
	InvalidCoinInitVolume  CodeType = 106
	InvalidCoinInitReserve CodeType = 107
	CodeInvalid            CodeType = 108

	// Buy Coin
	SameCoins              CodeType = 109
	CoinToBuyNotExists     CodeType = 110
	CoinToSellNotExists    CodeType = 111
	InsufficientCoinToSell CodeType = 112
)
