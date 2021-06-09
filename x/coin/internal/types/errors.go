package types

import (
	"fmt"
	"strconv"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"bitbucket.org/decimalteam/go-node/utils/errors"
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

func ErrInvalidCRR(crr string) *sdkerrors.Error {
	return sdkerrors.New(
		DefaultCodespace,
		CodeInvalidCRR,
		errors.Encode(
			fmt.Sprintf("coin CRR must be between 10 and 100, crr is: %s", crr),
			errors.NewParam("crr", crr)),
	)
}

func ErrCoinDoesNotExist(symbol string) *sdkerrors.Error {
	return sdkerrors.New(
		DefaultCodespace,
		CodeCoinDoesNotExist,
		errors.Encode(
			fmt.Sprintf("coin %s does not exist", symbol),
			errors.NewParam("symbol", symbol)),
	)
}

func ErrInvalidCoinSymbol(symbol string) *sdkerrors.Error {
	return sdkerrors.New(
		DefaultCodespace,
		CodeInvalidCoinSymbol,
		errors.Encode(
			fmt.Sprintf("invalid coin symbol %s. Symbol must match this regular expression: %s", symbol, allowedCoinSymbols),
			errors.NewParam("symbol", symbol)),
	)
}

func ErrForbiddenCoinSymbol(symbol string) *sdkerrors.Error {
	return sdkerrors.New(
		DefaultCodespace,
		CodeForbiddenCoinSymbol,
		errors.Encode(
			fmt.Sprintf("forbidden coin symbol %s", symbol),
			errors.NewParam("symbol", symbol)),
	)
}

func ErrRetrievedAnotherCoin(symbolWant, symbolRetrieved string) *sdkerrors.Error {
	return sdkerrors.New(
		DefaultCodespace,
		CodeRetrievedAnotherCoin,
		errors.Encode(
			fmt.Sprintf("retrieved coin %s instead %s", symbolRetrieved, symbolWant),
			errors.NewParam("symbol_want", symbolWant),
			errors.NewParam("symbol_retrieved", symbolRetrieved)),
	)
}

func ErrCoinAlreadyExist(coin string) *sdkerrors.Error {
	return sdkerrors.New(
		DefaultCodespace,
		CodeCoinAlreadyExists,
		errors.Encode(
			fmt.Sprintf("coin %s already exist", coin),
			errors.NewParam("coin", coin)),
	)
}

func ErrInvalidCoinTitle(title string) *sdkerrors.Error {
	return sdkerrors.New(
		DefaultCodespace,
		CodeInvalidCoinTitle,
		errors.Encode(
			fmt.Sprintf("invalid coin title: %s. Allowed up to %d bytes", title, maxCoinNameBytes),
			errors.NewParam("title", title)),
	)
}

func ErrInvalidCoinInitialVolume(initialVolume string) *sdkerrors.Error {
	return sdkerrors.New(
		DefaultCodespace,
		CodeInvalidCoinInitialVolume,
		errors.Encode(
			fmt.Sprintf("coin initial volume should be between %s and %s. Given %s", minCoinSupply.String(), maxCoinSupply.String(), initialVolume),
			errors.NewParam("min_coin_supply", minCoinSupply.String()),
			errors.NewParam("max_coin_supply", maxCoinSupply.String()),
			errors.NewParam("initial_volume", initialVolume)),
	)
}

func ErrInvalidCoinInitialReserve(reserve string) *sdkerrors.Error {
	return sdkerrors.New(
		DefaultCodespace,
		CodeInvalidCoinInitialReserve,
		errors.Encode(
			fmt.Sprintf("coin initial reserve should be greater than or equal to %s", reserve),
			errors.NewParam("reserve", reserve)),
	)
}

func ErrInternal(err string) *sdkerrors.Error {
	return sdkerrors.New(
		DefaultCodespace,
		CodeInternal,
		errors.Encode(
			fmt.Sprintf("Internal error: %s", err),
			errors.NewParam("err", err)),
	)
}

func ErrInsufficientCoinReserve() *sdkerrors.Error {
	return sdkerrors.New(
		DefaultCodespace,
		CodeInsufficientCoinReserve,
		errors.Encode(
			fmt.Sprintf("not enough coin to reserve")),
	)
}

func ErrInsufficientFundsToPayCommission(commission string) *sdkerrors.Error {
	return sdkerrors.New(
		DefaultCodespace,
		CodeInsufficientCoinToPayCommission,
		errors.Encode(
			fmt.Sprintf("insufficient funds to pay commission: wanted = %s", commission),
			errors.NewParam("commission", commission)),
	)
}

func ErrInsufficientFunds(fundsWant, fundsExist string) *sdkerrors.Error {
	return sdkerrors.New(
		DefaultCodespace,
		CodeInsufficientFunds,
		errors.Encode(
			fmt.Sprintf("insufficient account funds; %s < %s", fundsExist, fundsWant),
			errors.NewParam("funds_want", fundsWant),
			errors.NewParam("funds_exist", fundsExist)),
	)
}

func ErrCalculateCommission(err error) *sdkerrors.Error {
	return sdkerrors.New(
		DefaultCodespace,
		CodeCalculateCommission,
		errors.Encode(
			err.Error(),
		),
	)
}

func ErrUpdateOnlyForCreator() *sdkerrors.Error {
	return sdkerrors.New(
		DefaultCodespace,
		CodeForbiddenUpdate,
		errors.Encode(
			"updating allowed only for creator of coin",
		),
	)
}

func ErrSameCoin() *sdkerrors.Error {
	return sdkerrors.New(
		DefaultCodespace,
		CodeSameCoins,
		errors.Encode(
			"can't buy same coins",
		),
	)
}

func ErrInsufficientFundsToSellAll() *sdkerrors.Error {
	return sdkerrors.New(
		DefaultCodespace,
		CodeInsufficientFundsToSellAll,
		errors.Encode(
			"not enough coin to sell",
		),
	)
}

func ErrTxBreaksVolumeLimit(volume string, limitVolume string) *sdkerrors.Error {
	return sdkerrors.New(
		DefaultCodespace,
		CodeTxBreaksVolumeLimit,
		errors.Encode(
			fmt.Sprintf("tx breaks LimitVolume rule: %s > %s", volume, limitVolume),
			errors.NewParam("volume", volume),
			errors.NewParam("limit_volume", limitVolume)),
	)
}

func ErrTxBreaksMinReserveRule(minCoinReserve string, reserve string) *sdkerrors.Error {
	return sdkerrors.New(
		DefaultCodespace,
		CodeTxBreaksMinReserveLimit,
		errors.Encode(
			fmt.Sprintf("tx breaks MinReserveLimit rule: %s < %s", reserve, minCoinReserve),
			errors.NewParam("reserve", reserve),
			errors.NewParam("min_coin_reserve", minCoinReserve)),
	)
}

func ErrMaximumValueToSellReached(amount, max string) *sdkerrors.Error {
	return sdkerrors.New(
		DefaultCodespace,
		CodeMaximumValueToSellReached,
		errors.Encode(
			fmt.Sprintf("wanted to sell maximum %s, but required to spend %s at the moment", max, amount),
			errors.NewParam("max", max),
			errors.NewParam("amount", amount)),
	)
}

func ErrMinimumValueToBuyReached(amount, min string) *sdkerrors.Error {
	return sdkerrors.New(
		DefaultCodespace,
		CodeMinimumValueToBuyReached,
		errors.Encode(
			fmt.Sprintf("wanted to buy minimum %s, but expected to receive %s at the moment", min, amount),
			errors.NewParam("min", min),
			errors.NewParam("amount", amount)),
	)
}

func ErrUpdateBalance(account string, err string) *sdkerrors.Error {
	return sdkerrors.New(
		DefaultCodespace,
		CodeUpdateBalance,
		errors.Encode(
			fmt.Sprintf("unable to update balance of account %s: %s", account, err),
			errors.NewParam("account", account),
			errors.NewParam("err", err)),
	)
}

func ErrLimitVolumeBroken(volume string, limit string) *sdkerrors.Error {
	return sdkerrors.New(
		DefaultCodespace,
		CodeLimitVolumeBroken,
		errors.Encode(
			fmt.Sprintf("volume should be less than or equal the volume limit: %s > %s", volume, limit),
			errors.NewParam("volume", volume),
			errors.NewParam("limit", limit)),
	)
}

func ErrInvalidAmount() *sdkerrors.Error {
	return sdkerrors.New(
		DefaultCodespace,
		CodeInvalidAmount,
		errors.Encode(
			"amount should be greater than 0"),
	)
}

// Redeem check

func ErrInvalidCheck(data string) *sdkerrors.Error {
	return sdkerrors.New(
		DefaultCodespace,
		CodeInvalidCheck,
		errors.Encode(
			fmt.Sprintf("unable to parse check: %s", data),
			errors.NewParam("data", data)),
	)
}

func ErrInvalidProof(error string) *sdkerrors.Error {
	return sdkerrors.New(
		DefaultCodespace,
		CodeInvalidProof,
		errors.Encode(
			fmt.Sprintf("provided proof is invalid %s", error),
			errors.NewParam("error", error)),
	)
}

func ErrInvalidPassphrase(error string) *sdkerrors.Error {
	return sdkerrors.New(
		DefaultCodespace,
		CodeInvalidPassphrase,
		errors.Encode(
			fmt.Sprintf("unable to create private key from passphrase: %s", error),
			errors.NewParam("error", error)),
	)
}

func ErrInvalidChainID(wanted string, issued string) *sdkerrors.Error {
	return sdkerrors.New(
		DefaultCodespace,
		CodeInvalidChainID,
		errors.Encode(
			fmt.Sprintf("wanted chain ID %s, but check is issued for chain with ID %s", wanted, issued),
			errors.NewParam("wanted", wanted),
			errors.NewParam("issued", issued)),
	)
}

func ErrInvalidNonce() *sdkerrors.Error {
	return sdkerrors.New(
		DefaultCodespace,
		CodeInvalidNonce,
		errors.Encode(
			"nonce is too big (should be up to 16 bytes)"),
	)
}

func ErrCheckExpired(block uint64) *sdkerrors.Error {
	return sdkerrors.New(
		DefaultCodespace,
		CodeCheckExpired,
		errors.Encode(
			fmt.Sprintf("check was expired at block %d", block),
			errors.NewParam("block", strconv.FormatInt(int64(block), 10))),
	)
}

func ErrCheckRedeemed() *sdkerrors.Error {
	return sdkerrors.New(
		DefaultCodespace,
		CodeCheckRedeemed,
		errors.Encode(
			"check was redeemed already"),
	)
}

func ErrUnableDecodeCheck(check string) *sdkerrors.Error {
	return sdkerrors.New(
		DefaultCodespace,
		CodeUnableDecodeCheck,
		errors.Encode(
			fmt.Sprintf("unable to decode check from base58 %s", check),
			errors.NewParam("check", check)),
	)
}

func ErrUnableRPLEncodeCheck(error string) *sdkerrors.Error {
	return sdkerrors.New(
		DefaultCodespace,
		CodeUnableRPLEncodeCheck,
		errors.Encode(
			fmt.Sprintf("unable to RLP encode check receiver address: %s", error),
			errors.NewParam("error", error)),
	)
}

func ErrUnableSignCheck(error string) *sdkerrors.Error {
	return sdkerrors.New(
		DefaultCodespace,
		CodeUnableSignCheck,
		errors.Encode(
			fmt.Sprintf("unable to sign check receiver address by private key generated from passphrase: %s", error),
			errors.NewParam("error", error)),
	)
}

func ErrUnableDecodeProof() *sdkerrors.Error {
	return sdkerrors.New(
		DefaultCodespace,
		CodeUnableDecodeProof,
		errors.Encode(
			"unable to decode proof from base64"),
	)
}

func ErrUnableRecoverAddress(error string) *sdkerrors.Error {
	return sdkerrors.New(
		DefaultCodespace,
		CodeUnableRecoverAddress,
		errors.Encode(
			fmt.Sprintf("unable to recover check issuer address: %s", error),
			errors.NewParam("error", error)),
	)
}

func ErrUnableRecoverLockPkey(error string) *sdkerrors.Error {
	return sdkerrors.New(
		DefaultCodespace,
		CodeUnableRecoverLockPkey,
		errors.Encode(
			fmt.Sprintf("unable to recover lock public key from check: %s", error),
			errors.NewParam("error", error)),
	)
}

// AccountKeys Errors

func ErrInvalidPkey() *sdkerrors.Error {
	return sdkerrors.New(
		DefaultCodespace,
		CodeInvalidPkey,
		errors.Encode(
			"invalid private key"),
	)
}

func ErrUnableRetrieveArmoredPkey(name string, error string) *sdkerrors.Error {
	return sdkerrors.New(
		DefaultCodespace,
		CodeUnableRetriveArmoredPkey,
		errors.Encode(
			fmt.Sprintf("unable to retrieve armored private key for account %s: %s", name, error),
			errors.NewParam("name", name),
			errors.NewParam("error", error)),
	)
}

func ErrUnableRetrievePkey(name string, error string) *sdkerrors.Error {
	return sdkerrors.New(
		DefaultCodespace,
		CodeUnableRetrivePkey,
		errors.Encode(
			fmt.Sprintf("unable to retrieve private key for account %s: %s", name, error),
			errors.NewParam("name", name),
			errors.NewParam("error", error)),
	)
}

func ErrUnableRetrieveSECPPkey(name string, algo string) *sdkerrors.Error {
	return sdkerrors.New(
		DefaultCodespace,
		CodeUnableRetriveSECPPkey,
		errors.Encode(
			fmt.Sprintf("unable to retrieve secp256k1 private key for account %s: %s private key retrieved instead", name, algo),
			errors.NewParam("name", name),
			errors.NewParam("algo", algo)),
	)
}
