package types

import (
	"errors"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// Local code type
type CodeType = uint32

const (
	// Default coin codespace
	DefaultCodespace string = ModuleName
	// Create coin
	CodeInvalidCRR                      CodeType = 100
	CodeCoinDoesNotExist                CodeType = 101
	CodeInvalidCoinSymbol               CodeType = 102
	CodeForbiddenCoinSymbol             CodeType = 103
	CodeRetrievedAnotherCoin            CodeType = 104
	CodeCoinAlreadyExists               CodeType = 105
	CodeInvalidCoinTitle                CodeType = 106
	CodeInvalidCoinInitialVolume        CodeType = 107
	CodeInvalidCoinInitialReserve       CodeType = 108
	CodeInternal                        CodeType = 109
	CodeInsufficientCoinReserve         CodeType = 110
	CodeInsufficientCoinToPayCommission CodeType = 111
	CodeInsufficientFunds               CodeType = 112
	CodeCalculateCommission             CodeType = 113
	CodeForbiddenUpdate                 CodeType = 114

	// Buy/Sell coin
	CodeSameCoins                  CodeType = 200
	CodeInsufficientFundsToSellAll CodeType = 201
	CodeTxBreaksVolumeLimit        CodeType = 202
	CodeTxBreaksMinReserveLimit    CodeType = 203
	CodeMaximumValueToSellReached  CodeType = 204
	CodeMinimumValueToBuyReached   CodeType = 205
	CodeUpdateBalance              CodeType = 206
	CodeLimitVolumeBroken          CodeType = 207
	// Send coin
	CodeInvalidAmount CodeType = 300
	CodeInvalidReceiverAddress CodeType = 301
	// Redeem check
	InvalidCheck      CodeType = 400
	InvalidProof      CodeType = 401
	InvalidPassphrase CodeType = 402
	InvalidChainID    CodeType = 403
	InvalidNonce      CodeType = 404
	CheckExpired      CodeType = 405
	CheckRedeemed     CodeType = 406
)

var ErrInvalidLimitVolume = errors.New("invalid limitVolume")

func ErrInvalidCRR() *sdkerrors.Error {
	return sdkerrors.New(DefaultCodespace, CodeInvalidCRR, "coin CRR must be between 10 and 100")
}

func ErrInvalidCoinSymbol(symbol string) *sdkerrors.Error {
	return sdkerrors.New(DefaultCodespace, CodeInvalidCoinSymbol, fmt.Sprintf("invalid coin symbol %s. Symbol must match this regular expression: %s", symbol, allowedCoinSymbols))
}

func ErrForbiddenCoinSymbol(symbol string) *sdkerrors.Error {
	return sdkerrors.New(DefaultCodespace, CodeInvalidCoinSymbol, fmt.Sprintf("forbidden coin symbol %s", symbol))
}

func ErrUpdateOnlyForCreator() *sdkerrors.Error {
	return sdkerrors.New(DefaultCodespace, CodeForbiddenUpdate, fmt.Sprintf("updating allowed only for creator of coin"))
}

func ErrCoinDoesNotExist(symbol string) *sdkerrors.Error {
	return sdkerrors.New(DefaultCodespace, CodeCoinDoesNotExist, fmt.Sprintf("coin %s does not exist", symbol))
}

func ErrRetrievedAnotherCoin(symbolWant, symbolRetrieved string) *sdkerrors.Error {
	return sdkerrors.New(DefaultCodespace, CodeRetrievedAnotherCoin, fmt.Sprintf("retrieved coin %s instead %s", symbolRetrieved, symbolWant))
}

func ErrInvalidCoinTitle() *sdkerrors.Error {
	return sdkerrors.New(DefaultCodespace, CodeInvalidCoinTitle, fmt.Sprintf("invalid coin title. Allowed up to %d bytes", maxCoinNameBytes))
}

func ErrInvalidCoinInitialVolume(initialVolume string) *sdkerrors.Error {
	return sdkerrors.New(DefaultCodespace, CodeInvalidCoinInitialVolume, fmt.Sprintf("coin initial volume should be between %s and %s. Given %s", minCoinSupply.String(), maxCoinSupply.String(), initialVolume))
}

func ErrInvalidCoinInitialReserve(ctx sdk.Context) *sdkerrors.Error {
	return sdkerrors.New(DefaultCodespace, CodeInvalidCoinInitialReserve, fmt.Sprintf("coin initial reserve should be greater than or equal to %s", MinCoinReserve(ctx).String()))
}

func ErrInternal(err string) *sdkerrors.Error {
	return sdkerrors.New(DefaultCodespace, CodeInternal, err)
}

func ErrInsufficientCoinReserve() *sdkerrors.Error {
	return sdkerrors.New(DefaultCodespace, CodeInsufficientCoinReserve, "not enough coin to reserve")
}

func ErrInsufficientFundsToPayCommission(commission string) *sdkerrors.Error {
	return sdkerrors.New(DefaultCodespace, CodeInsufficientCoinToPayCommission, fmt.Sprintf("insufficient funds to pay commission: wanted = %s", commission))
}

func ErrInsufficientFunds(fundsWant, fundsExist string) *sdkerrors.Error {
	return sdkerrors.New(DefaultCodespace, CodeInsufficientFunds, fmt.Sprintf("insufficient account funds; %s < %s", fundsExist, fundsWant))
}

func ErrInsufficientFundsToSellAll() *sdkerrors.Error {
	return sdkerrors.New(DefaultCodespace, CodeInsufficientFundsToSellAll, fmt.Sprintf("not enough coin to sell"))
}

func ErrUpdateBalance(account string, err string) *sdkerrors.Error {
	return sdkerrors.New(DefaultCodespace, CodeUpdateBalance, fmt.Sprintf("unable to update balance of account %s: %s", account, err))
}

func ErrCalculateCommission(err error) *sdkerrors.Error {
	return sdkerrors.New(DefaultCodespace, CodeCalculateCommission, err.Error())
}

func ErrCoinAlreadyExist(coin string) *sdkerrors.Error {
	return sdkerrors.New(DefaultCodespace, CodeCoinAlreadyExists, fmt.Sprintf("coin %s already exist", coin))
}

func ErrLimitVolumeBroken(volume string, limit string) *sdkerrors.Error {
	return sdkerrors.New(DefaultCodespace, CodeLimitVolumeBroken, fmt.Sprintf("volume should be less than or equal the volume limit: %s > %s", volume, limit))
}

func ErrSameCoin() *sdkerrors.Error {
	return sdkerrors.New(DefaultCodespace, CodeSameCoins, "can't buy same coins")
}

func ErrTxBreaksVolumeLimit(volume, limitVolume string) *sdkerrors.Error {
	return sdkerrors.New(DefaultCodespace, CodeTxBreaksVolumeLimit, fmt.Sprintf("tx breaks LimitVolume rule: %s > %s", volume, limitVolume))
}

func ErrTxBreaksMinReserveRule(ctx sdk.Context, reserve string) *sdkerrors.Error {
	return sdkerrors.New(DefaultCodespace, CodeTxBreaksMinReserveLimit, fmt.Sprintf("tx breaks MinReserveLimit rule: %s < %s", reserve, MinCoinReserve(ctx).String()))
}

func ErrMaximumValueToSellReached(amount, max string) *sdkerrors.Error {
	return sdkerrors.New(DefaultCodespace, CodeMaximumValueToSellReached, fmt.Sprintf("wanted to sell maximum %s, but required to spend %s at the moment", max, amount))
}

func ErrMinimumValueToBuyReached(amount, min string) *sdkerrors.Error {
	return sdkerrors.New(DefaultCodespace, CodeMinimumValueToBuyReached, fmt.Sprintf("wanted to buy minimum %s, but expected to receive %s at the moment", min, amount))
}

func ErrInvalidAmount() *sdkerrors.Error {
	return sdkerrors.New(DefaultCodespace, CodeInvalidAmount, "amount should be greater than 0")
}

func ErrReceiverEmpty() *sdkerrors.Error {
	return sdkerrors.New(DefaultCodespace, CodeInvalidReceiverAddress, "Receiver cannot be empty")

}