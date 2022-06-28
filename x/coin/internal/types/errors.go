package types

import (
	"fmt"

	"bitbucket.org/decimalteam/go-node/utils/errors"
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
	CodeMinVolumeBroken            CodeType = 208

	// Send coin
	CodeInvalidAmount          CodeType = 300
	CodeInvalidReceiverAddress CodeType = 301
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
	return errors.Encode(
		DefaultCodespace,
		CodeInvalidCRR,
		fmt.Sprintf("coin CRR must be between 10 and 100, crr is: %s", crr),
	)
}

func ErrCoinDoesNotExist(symbol string) *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeCoinDoesNotExist,
		fmt.Sprintf("coin %s does not exist", symbol),
	)
}

func ErrInvalidCoinSymbol(symbol string) *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeInvalidCoinSymbol,
		fmt.Sprintf("invalid coin symbol %s. Symbol must match this regular expression: %s", symbol, allowedCoinSymbols),
	)
}

func ErrForbiddenCoinSymbol(symbol string) *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeForbiddenCoinSymbol,
		fmt.Sprintf("forbidden coin symbol %s", symbol),
	)
}

func ErrRetrievedAnotherCoin(symbolWant string, symbolRetrieved string) *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeRetrievedAnotherCoin,
		fmt.Sprintf("retrieved coin %s instead %s", symbolRetrieved, symbolWant),
	)
}

func ErrCoinAlreadyExist(coin string) *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeCoinAlreadyExists,
		fmt.Sprintf("coin %s already exist", coin),
	)
}

func ErrInvalidCoinTitle(title string) *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeInvalidCoinTitle,
		fmt.Sprintf("invalid coin title: %s. Allowed up to %d bytes", title, maxCoinNameBytes),
	)
}

func ErrInvalidCoinInitialVolume(initialVolume string) *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeInvalidCoinInitialVolume,
		fmt.Sprintf("coin initial volume should be between %s and %s. Given %s", MinCoinSupply.String(), maxCoinSupply.String(), initialVolume),
	)
}

func ErrInvalidCoinInitialReserve(reserve string) *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeInvalidCoinInitialReserve,
		fmt.Sprintf("coin initial reserve should be greater than or equal to %s", reserve),
	)
}

func ErrInternal(err string) *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeInternal,
		fmt.Sprintf("Internal error: %s", err),
	)
}

func ErrInsufficientCoinReserve() *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeInsufficientCoinReserve,
		"not enough coin to reserve",
	)
}

func ErrInsufficientFundsToPayCommission(commission string) *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeInsufficientCoinToPayCommission,
		fmt.Sprintf("insufficient funds to pay commission: wanted = %s", commission),
	)
}

func ErrInsufficientFunds(fundsWant string, fundsExist string) *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeInsufficientFunds,
		fmt.Sprintf("insufficient account funds; %s < %s", fundsExist, fundsWant),
	)
}

func ErrCalculateCommission(err string) *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeCalculateCommission,
		err,
	)
}

func ErrUpdateOnlyForCreator() *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeForbiddenUpdate,
		"updating allowed only for creator of coin",
	)
}

func ErrSameCoin() *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeSameCoins,
		"can't buy same coins",
	)
}

func ErrInsufficientFundsToSellAll() *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeInsufficientFundsToSellAll,
		"not enough coin to sell",
	)
}

func ErrTxBreaksVolumeLimit(volume string, limitVolume string) *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeTxBreaksVolumeLimit,
		fmt.Sprintf("tx breaks LimitVolume rule: %s > %s", volume, limitVolume),
	)
}

func ErrTxBreaksMinReserveRule(minCoinReserve string, reserve string) *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeTxBreaksMinReserveLimit,
		fmt.Sprintf("tx breaks MinReserveLimit rule: %s < %s", reserve, minCoinReserve),
	)
}

func ErrMaximumValueToSellReached(amount string, max string) *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeMaximumValueToSellReached,
		fmt.Sprintf("wanted to sell maximum %s, but required to spend %s at the moment", max, amount),
	)
}

func ErrMinimumValueToBuyReached(amount string, min string) *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeMinimumValueToBuyReached,
		fmt.Sprintf("wanted to buy minimum %s, but expected to receive %s at the moment", min, amount),
	)
}

func ErrUpdateBalance(account string, err string) *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeUpdateBalance,
		fmt.Sprintf("unable to update balance of account %s: %s", account, err),
	)
}

func ErrLimitVolumeBroken(volume string, limit string) *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeLimitVolumeBroken,
		fmt.Sprintf("volume should be less than or equal the volume limit: %s > %s", volume, limit),
	)
}

func ErrMinVolumeBroken(volume string) *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeMinVolumeBroken,
		fmt.Sprintf("volume should be greater than or equal to the min supply: %s < %s", volume, MinCoinSupply.String()),
	)
}

func ErrInvalidAmount() *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeInvalidAmount,
		"amount should be greater than 0",
	)
}

// Redeem check

func ErrInvalidCheck(data string) *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeInvalidCheck,
		fmt.Sprintf("unable to parse check: %s", data),
	)
}

func ErrInvalidProof(error string) *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeInvalidProof,
		fmt.Sprintf("provided proof is invalid %s", error),
	)
}

func ErrInvalidPassphrase(error string) *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeInvalidPassphrase,
		fmt.Sprintf("unable to create private key from passphrase: %s", error),
	)
}

func ErrInvalidChainID(wanted string, issued string) *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeInvalidChainID,
		fmt.Sprintf("wanted chain ID %s, but check is issued for chain with ID %s", wanted, issued),
	)
}

func ErrInvalidNonce() *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeInvalidNonce,
		"nonce is too big (should be up to 16 bytes)",
	)
}

func ErrCheckExpired(block string) *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeCheckExpired,
		fmt.Sprintf("check was expired at block %s", block),
	)
}

func ErrCheckRedeemed() *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeCheckRedeemed,
		"check was redeemed already",
	)
}

func ErrUnableDecodeCheck(check string) *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeUnableDecodeCheck,
		fmt.Sprintf("unable to decode check from base58 %s", check),
	)
}

func ErrUnableRPLEncodeCheck(error string) *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeUnableRPLEncodeCheck,
		fmt.Sprintf("unable to RLP encode check receiver address: %s", error),
	)
}

func ErrUnableSignCheck(error string) *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeUnableSignCheck,
		fmt.Sprintf("unable to sign check receiver address by private key generated from passphrase: %s", error),
	)
}

func ErrUnableDecodeProof() *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeUnableDecodeProof,
		"unable to decode proof from base64",
	)
}

func ErrUnableRecoverAddress(error string) *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeUnableRecoverAddress,
		fmt.Sprintf("unable to recover check issuer address: %s", error),
	)
}

func ErrUnableRecoverLockPkey(error string) *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeUnableRecoverLockPkey,
		fmt.Sprintf("unable to recover lock public key from check: %s", error),
	)
}

// AccountKeys Errors

func ErrInvalidPkey() *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeInvalidPkey,
		"invalid private key",
	)
}

func ErrUnableRetrieveArmoredPkey(name string, error string) *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeUnableRetriveArmoredPkey,
		fmt.Sprintf("unable to retrieve armored private key for account %s: %s", name, error),
	)
}

func ErrUnableRetrievePkey(name string, error string) *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeUnableRetrivePkey,
		fmt.Sprintf("unable to retrieve private key for account %s: %s", name, error),
	)
}

func ErrUnableRetrieveSECPPkey(name string, algo string) *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeUnableRetriveSECPPkey,
		fmt.Sprintf("unable to retrieve secp256k1 private key for account %s: %s private key retrieved instead", name, algo),
	)
}

func ErrReceiverEmpty() *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeInvalidReceiverAddress,
		"Receiver cannot be empty ",
	)
}
