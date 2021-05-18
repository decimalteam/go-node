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
	// Redeem check
	CodeInvalidCheck          CodeType = 400
	CodeInvalidProof          CodeType = 401
	CodeInvalidPassphrase     CodeType = 402
	CodeInvalidChainID        CodeType = 403
	CodeInvalidNonce          CodeType = 404
	CodeCheckExpired          CodeType = 405
	CodeCheckRedeemed         CodeType = 406
	CodeUnableDecodeCheck     CodeType = 407
	CodeUnableRPLEncodeCheck  CodeType = 408
	CodeUnableSignCheck       CodeType = 409
	CodeUnableDecodeProof     CodeType = 410
	CodeUnableRecoverAddress  CodeType = 411
	CodeUnableRecoverLockPkey CodeType = 412
	// AccountKeys
	CodeInvalidPkey              CodeType = 500
	CodeUnableRetriveArmoredPkey CodeType = 501
	CodeUnableRetrivePkey        CodeType = 502
	CodeUnableRetriveSECPPkey    CodeType = 503
)

var ErrInvalidLimitVolume = errors.New("invalid limitVolume")

func ErrInvalidCRR() *sdkerrors.Error {
	return sdkerrors.New(
		DefaultCodespace,
		CodeInvalidCRR,
		"coin CRR must be between 10 and 100",
	)
}

func ErrInvalidCoinSymbol(symbol string) *sdkerrors.Error {
	return sdkerrors.New(
		DefaultCodespace,
		CodeInvalidCoinSymbol,
		fmt.Sprintf("invalid coin symbol %s. Symbol must match this regular expression: %s", symbol, allowedCoinSymbols),
	)
}

func ErrForbiddenCoinSymbol(symbol string) *sdkerrors.Error {
	return sdkerrors.New(
		DefaultCodespace,
		CodeForbiddenCoinSymbol,
		fmt.Sprintf("forbidden coin symbol %s", symbol),
	)
}

func ErrUpdateOnlyForCreator() *sdkerrors.Error {
	return sdkerrors.New(
		DefaultCodespace,
		CodeForbiddenUpdate,
		fmt.Sprintf("updating allowed only for creator of coin"),
	)
}

func ErrCoinDoesNotExist(symbol string) *sdkerrors.Error {
	return sdkerrors.New(
		DefaultCodespace,
		CodeCoinDoesNotExist,
		fmt.Sprintf("coin %s does not exist", symbol),
	)
}

func ErrRetrievedAnotherCoin(symbolWant, symbolRetrieved string) *sdkerrors.Error {
	return sdkerrors.New(
		DefaultCodespace,
		CodeRetrievedAnotherCoin,
		fmt.Sprintf("retrieved coin %s instead %s", symbolRetrieved, symbolWant),
	)
}

func ErrInvalidCoinTitle() *sdkerrors.Error {
	return sdkerrors.New(
		DefaultCodespace,
		CodeInvalidCoinTitle,
		fmt.Sprintf("invalid coin title. Allowed up to %d bytes", maxCoinNameBytes),
	)
}

func ErrInvalidCoinInitialVolume(initialVolume string) *sdkerrors.Error {
	return sdkerrors.New(
		DefaultCodespace,
		CodeInvalidCoinInitialVolume,
		fmt.Sprintf("coin initial volume should be between %s and %s. Given %s", minCoinSupply.String(), maxCoinSupply.String(), initialVolume),
	)
}

func ErrInvalidCoinInitialReserve(ctx sdk.Context) *sdkerrors.Error {
	return sdkerrors.New(
		DefaultCodespace,
		CodeInvalidCoinInitialReserve,
		fmt.Sprintf("coin initial reserve should be greater than or equal to %s", MinCoinReserve(ctx).String()),
	)
}

func ErrInternal(err string) *sdkerrors.Error {
	return sdkerrors.New(
		DefaultCodespace,
		CodeInternal,
		err,
	)
}

func ErrInsufficientCoinReserve() *sdkerrors.Error {
	return sdkerrors.New(
		DefaultCodespace,
		CodeInsufficientCoinReserve,
		"not enough coin to reserve",
	)
}

func ErrInsufficientFundsToPayCommission(commission string) *sdkerrors.Error {
	return sdkerrors.New(
		DefaultCodespace,
		CodeInsufficientCoinToPayCommission,
		fmt.Sprintf("insufficient funds to pay commission: wanted = %s", commission),
	)
}

func ErrInsufficientFunds(fundsWant, fundsExist string) *sdkerrors.Error {
	return sdkerrors.New(
		DefaultCodespace,
		CodeInsufficientFunds,
		fmt.Sprintf("insufficient account funds; %s < %s", fundsExist, fundsWant),
	)
}

func ErrInsufficientFundsToSellAll() *sdkerrors.Error {
	return sdkerrors.New(
		DefaultCodespace,
		CodeInsufficientFundsToSellAll,
		fmt.Sprintf("not enough coin to sell"),
	)
}

func ErrUpdateBalance(account string, err string) *sdkerrors.Error {
	return sdkerrors.New(
		DefaultCodespace,
		CodeUpdateBalance,
		fmt.Sprintf("unable to update balance of account %s: %s", account, err),
	)
}

func ErrCalculateCommission(err error) *sdkerrors.Error {
	return sdkerrors.New(
		DefaultCodespace,
		CodeCalculateCommission,
		err.Error(),
	)
}

func ErrCoinAlreadyExist(coin string) *sdkerrors.Error {
	return sdkerrors.New(
		DefaultCodespace,
		CodeCoinAlreadyExists,
		fmt.Sprintf("coin %s already exist", coin),
	)
}

func ErrLimitVolumeBroken(volume string, limit string) *sdkerrors.Error {
	return sdkerrors.New(
		DefaultCodespace,
		CodeLimitVolumeBroken,
		fmt.Sprintf("volume should be less than or equal the volume limit: %s > %s", volume, limit),
	)
}

func ErrSameCoin() *sdkerrors.Error {
	return sdkerrors.New(
		DefaultCodespace,
		CodeSameCoins,
		"can't buy same coins",
	)
}

func ErrTxBreaksVolumeLimit(volume, limitVolume string) *sdkerrors.Error {
	return sdkerrors.New(
		DefaultCodespace,
		CodeTxBreaksVolumeLimit,
		fmt.Sprintf("tx breaks LimitVolume rule: %s > %s", volume, limitVolume),
	)
}

func ErrTxBreaksMinReserveRule(ctx sdk.Context, reserve string) *sdkerrors.Error {
	return sdkerrors.New(
		DefaultCodespace,
		CodeTxBreaksMinReserveLimit,
		fmt.Sprintf("tx breaks MinReserveLimit rule: %s < %s", reserve, MinCoinReserve(ctx).String()),
	)
}

func ErrMaximumValueToSellReached(amount, max string) *sdkerrors.Error {
	return sdkerrors.New(
		DefaultCodespace,
		CodeMaximumValueToSellReached,
		fmt.Sprintf("wanted to sell maximum %s, but required to spend %s at the moment", max, amount),
	)
}

func ErrMinimumValueToBuyReached(amount, min string) *sdkerrors.Error {
	return sdkerrors.New(
		DefaultCodespace,
		CodeMinimumValueToBuyReached,
		fmt.Sprintf("wanted to buy minimum %s, but expected to receive %s at the moment", min, amount),
	)
}

func ErrInvalidAmount() *sdkerrors.Error {
	return sdkerrors.New(
		DefaultCodespace,
		CodeInvalidAmount,
		"amount should be greater than 0",
	)
}

// Redeem check

func ErrInvalidCheck(data string) *sdkerrors.Error {
	return sdkerrors.New(
		DefaultCodespace,
		CodeInvalidCheck,
		fmt.Sprintf("unable to parse check: %s", data),
	)
}

func ErrInvalidProof(data string) *sdkerrors.Error {
	return sdkerrors.New(
		DefaultCodespace,
		CodeInvalidProof,
		fmt.Sprintf("provided proof is invalid %s", data),
	)
}

func ErrInvalidPassphrase(data string) *sdkerrors.Error {
	return sdkerrors.New(
		DefaultCodespace,
		CodeInvalidPassphrase,
		fmt.Sprintf("unable to create private key from passphrase: %s", data),
	)
}

func ErrInvalidChainID(wanted string, issued string) *sdkerrors.Error {
	return sdkerrors.New(
		DefaultCodespace,
		CodeInvalidChainID,
		fmt.Sprintf("wanted chain ID %s, but check is issued for chain with ID %s", wanted, issued),
	)
}

func ErrInvalidNonce() *sdkerrors.Error {
	return sdkerrors.New(
		DefaultCodespace,
		CodeInvalidNonce,
		fmt.Sprintf("nonce is too big (should be up to 16 bytes)"),
	)
}

func ErrCheckExpired(date uint64) *sdkerrors.Error {
	return sdkerrors.New(
		DefaultCodespace,
		CodeCheckExpired,
		fmt.Sprintf("check was expired at block %s", date),
	)
}

func ErrCheckRedeemed() *sdkerrors.Error {
	return sdkerrors.New(
		DefaultCodespace,
		CodeCheckRedeemed,
		fmt.Sprintf("check was redeemed already"),
	)
}

func ErrUnableDecodeCheck() *sdkerrors.Error {
	return sdkerrors.New(
		DefaultCodespace,
		CodeUnableDecodeCheck,
		fmt.Sprintf("unable to decode check from base58"),
	)
}

func ErrUnableRPLEncodeCheck(err string) *sdkerrors.Error {
	return sdkerrors.New(
		DefaultCodespace,
		CodeUnableRPLEncodeCheck,
		fmt.Sprintf("unable to RLP encode check receiver address: %s", err),
	)
}

func ErrUnableSignCheck(err string) *sdkerrors.Error {
	return sdkerrors.New(
		DefaultCodespace,
		CodeUnableSignCheck,
		fmt.Sprintf("unable to sign check receiver address by private key generated from passphrase: %s", err),
	)
}

func ErrUnableDecodeProof() *sdkerrors.Error {
	return sdkerrors.New(
		DefaultCodespace,
		CodeUnableDecodeProof,
		fmt.Sprintf("unable to decode proof from base64"),
	)
}

func ErrUnableRecoverAddress(err string) *sdkerrors.Error {
	return sdkerrors.New(
		DefaultCodespace,
		CodeUnableRecoverAddress,
		fmt.Sprintf("unable to recover check issuer address: %s", err),
	)
}

func ErrUnableRecoverLockPkey(err string) *sdkerrors.Error {
	return sdkerrors.New(
		DefaultCodespace,
		CodeUnableRecoverLockPkey,
		fmt.Sprintf("unable to recover lock public key from check: %s", err),
	)
}

// AccountKeys Errors

func ErrInvalidPkey() *sdkerrors.Error {
	return sdkerrors.New(
		DefaultCodespace,
		CodeInvalidPkey,
		fmt.Sprintf("invalid private key"),
	)
}

func ErrUnableRetriveArmoredPkey(name string, err string) *sdkerrors.Error {
	return sdkerrors.New(
		DefaultCodespace,
		CodeUnableRetriveArmoredPkey,
		fmt.Sprintf("unable to retrieve armored private key for account %s: %s", name, err),
	)
}

func ErrUnableRetrivePkey(name string, err string) *sdkerrors.Error {
	return sdkerrors.New(
		DefaultCodespace,
		CodeUnableRetrivePkey,
		fmt.Sprintf("unable to retrieve private key for account %s: %s", name, err),
	)
}

func ErrUnableRetriveSECPPkey(name string, algo string) *sdkerrors.Error {
	return sdkerrors.New(
		DefaultCodespace,
		CodeUnableRetriveSECPPkey,
		fmt.Sprintf("unable to retrieve secp256k1 private key for account %s: %s private key retrieved instead", name, algo),
	)
}
