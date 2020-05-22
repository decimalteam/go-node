package types

import (
	"fmt"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// Local code type
type CodeType = uint32

const (
	// Default coin codespace
	DefaultCodespace string = ModuleName
	// Create coin
	DecodeError                     CodeType = 101
	InvalidCRR                      CodeType = 102
	InvalidCoinSymbol               CodeType = 103
	CoinAlreadyExists               CodeType = 104
	InvalidCoinTitle                CodeType = 105
	InvalidCoinInitVolume           CodeType = 106
	InvalidCoinInitReserve          CodeType = 107
	CodeInvalid                     CodeType = 108
	InsufficientCoinReserve         CodeType = 118
	InsufficientCoinToPayCommission CodeType = 120
	InsufficientCoinToCreateCoin    CodeType = 121
	// Buy/Sell Coin
	SameCoins                 CodeType = 109
	CoinToBuyNotExists        CodeType = 110
	CoinToSellNotExists       CodeType = 111
	InsufficientCoinToSell    CodeType = 112
	TxBreaksVolumeLimit       CodeType = 113
	TxBreaksMinReserveLimit   CodeType = 114
	MaximumValueToSellReached CodeType = 115
	MinimumValueToBuyReached  CodeType = 116
	UpdateBalanceError        CodeType = 117
	// Send Coin
	InvalidAmount CodeType = 119
)

func ErrorInsufficientCoinToPayCommission(commission string) *sdkerrors.Error {
	return sdkerrors.New(DefaultCodespace, InsufficientCoinToPayCommission, fmt.Sprintf("Insufficient coin to pay commission: wanted = %s", commission))
}

func ErrorInsufficientFundsToCreateCoin(funds string) *sdkerrors.Error {
	return sdkerrors.New(DefaultCodespace, InsufficientCoinToPayCommission, fmt.Sprintf("Insufficient funds to create coin: wanted = %s", funds))
}

func ErrorUpdateBalance(err error) *sdkerrors.Error {
	return sdkerrors.New(DefaultCodespace, UpdateBalanceError, err.Error())
}
