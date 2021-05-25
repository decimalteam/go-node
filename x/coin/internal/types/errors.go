package types

import (
	"encoding/json"
	"errors"
	"fmt"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"strconv"
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

func ErrInvalidCRR(crr string) *sdkerrors.Error {

	jsonData, _ := json.Marshal(
		map[string]string{
			"code":      getCodeString(CodeInvalidCRR),
			"codespace": DefaultCodespace,
			"desc":      fmt.Sprintf("coin CRR must be between 10 and 100, crr is: %s", crr),
			"crr":       crr,
		},
	)

	return sdkerrors.New(
		DefaultCodespace,
		CodeInvalidCRR,
		string(jsonData),
	)
}

func ErrCoinDoesNotExist(symbol string) *sdkerrors.Error {
	jsonData, _ := json.Marshal(
		map[string]string{
			"code":      getCodeString(CodeCoinDoesNotExist),
			"codespace": DefaultCodespace,
			"desc":      fmt.Sprintf("coin %s does not exist", symbol),
			"symbol":    symbol,
		},
	)
	return sdkerrors.New(
		DefaultCodespace,
		CodeCoinDoesNotExist,
		string(jsonData),
	)
}

func ErrInvalidCoinSymbol(symbol string) *sdkerrors.Error {
	jsonData, _ := json.Marshal(
		map[string]string{
			"code":      getCodeString(CodeInvalidCoinSymbol),
			"codespace": DefaultCodespace,
			"desc":      fmt.Sprintf("invalid coin symbol %s. Symbol must match this regular expression: %s", symbol, allowedCoinSymbols),
			"symbol":    symbol,
		},
	)

	return sdkerrors.New(
		DefaultCodespace,
		CodeInvalidCoinSymbol,
		string(jsonData),
	)
}

func ErrForbiddenCoinSymbol(symbol string) *sdkerrors.Error {
	jsonData, _ := json.Marshal(
		map[string]string{
			"code":      getCodeString(CodeForbiddenCoinSymbol),
			"codespace": DefaultCodespace,
			"desc":      fmt.Sprintf("forbidden coin symbol %s", symbol),
			"symbol":    symbol,
		},
	)
	return sdkerrors.New(
		DefaultCodespace,
		CodeForbiddenCoinSymbol,
		string(jsonData),
	)
}

func ErrRetrievedAnotherCoin(symbolWant, symbolRetrieved string) *sdkerrors.Error {
	jsonData, _ := json.Marshal(
		map[string]string{
			"code":            getCodeString(CodeRetrievedAnotherCoin),
			"codespace":       DefaultCodespace,
			"desc":            fmt.Sprintf("retrieved coin %s instead %s", symbolRetrieved, symbolWant),
			"symbolWant":      symbolWant,
			"symbolRetrieved": symbolRetrieved,
		},
	)
	return sdkerrors.New(
		DefaultCodespace,
		CodeRetrievedAnotherCoin,
		string(jsonData),
	)
}

func ErrCoinAlreadyExist(coin string) *sdkerrors.Error {
	jsonData, _ := json.Marshal(
		map[string]string{
			"code":      getCodeString(CodeCoinAlreadyExists),
			"codespace": DefaultCodespace,
			"desc":      fmt.Sprintf("coin %s already exist", coin),
			"coin":      coin,
		},
	)
	return sdkerrors.New(
		DefaultCodespace,
		CodeCoinAlreadyExists,
		string(jsonData),
	)
}

func ErrInvalidCoinTitle(title string) *sdkerrors.Error {
	jsonData, _ := json.Marshal(
		map[string]string{
			"code":      getCodeString(CodeInvalidCoinTitle),
			"codespace": DefaultCodespace,
			"desc":      fmt.Sprintf("invalid coin title: %s. Allowed up to %d bytes", title, maxCoinNameBytes),
			"title":     title,
		},
	)
	return sdkerrors.New(
		DefaultCodespace,
		CodeInvalidCoinTitle,
		string(jsonData),
	)
}

func ErrInvalidCoinInitialVolume(initialVolume string) *sdkerrors.Error {
	jsonData, _ := json.Marshal(
		map[string]string{
			"code":          getCodeString(CodeInvalidCoinInitialVolume),
			"codespace":     DefaultCodespace,
			"desc":          fmt.Sprintf("coin initial volume should be between %s and %s. Given %s", minCoinSupply.String(), maxCoinSupply.String(), initialVolume),
			"minCoinSupply": minCoinSupply.String(),
			"maxCoinSupply": maxCoinSupply.String(),
			"initialVolume": initialVolume,
		},
	)
	return sdkerrors.New(
		DefaultCodespace,
		CodeInvalidCoinInitialVolume,
		string(jsonData),
	)
}

func ErrInvalidCoinInitialReserve(reserve string) *sdkerrors.Error {
	jsonData, _ := json.Marshal(
		map[string]string{
			"code":      getCodeString(CodeInvalidCoinInitialVolume),
			"codespace": DefaultCodespace,
			"desc":      fmt.Sprintf("coin initial reserve should be greater than or equal to %s", reserve),
			"reserve":   reserve,
		},
	)
	return sdkerrors.New(
		DefaultCodespace,
		CodeInvalidCoinInitialReserve,
		string(jsonData),
	)
}

func ErrInternal(err string) *sdkerrors.Error {
	jsonData, _ := json.Marshal(
		map[string]string{
			"code":      getCodeString(CodeInternal),
			"codespace": DefaultCodespace,
			"desc":      fmt.Sprintf("Internal error: %s", err),
			"err":       err,
		},
	)
	return sdkerrors.New(
		DefaultCodespace,
		CodeInternal,
		string(jsonData),
	)
}

func ErrInsufficientCoinReserve() *sdkerrors.Error {
	jsonData, _ := json.Marshal(
		map[string]string{
			"code":      getCodeString(CodeInsufficientCoinReserve),
			"codespace": DefaultCodespace,
			"desc":      fmt.Sprintf("not enough coin to reserve"),
		},
	)
	return sdkerrors.New(
		DefaultCodespace,
		CodeInsufficientCoinReserve,
		string(jsonData),
	)
}

func ErrInsufficientFundsToPayCommission(commission string) *sdkerrors.Error {
	jsonData, _ := json.Marshal(
		map[string]string{
			"code":       getCodeString(CodeInsufficientCoinToPayCommission),
			"codespace":  DefaultCodespace,
			"desc":       fmt.Sprintf("insufficient funds to pay commission: wanted = %s", commission),
			"commission": commission,
		},
	)
	return sdkerrors.New(
		DefaultCodespace,
		CodeInsufficientCoinToPayCommission,
		string(jsonData),
	)
}

func ErrInsufficientFunds(fundsWant, fundsExist string) *sdkerrors.Error {
	jsonData, _ := json.Marshal(
		map[string]string{
			"code":       getCodeString(CodeInsufficientFunds),
			"codespace":  DefaultCodespace,
			"desc":       fmt.Sprintf("insufficient account funds; %s < %s", fundsExist, fundsWant),
			"fundsWant":  fundsWant,
			"fundsExist": fundsExist,
		},
	)
	return sdkerrors.New(
		DefaultCodespace,
		CodeInsufficientFunds,
		string(jsonData),
	)
}

func ErrCalculateCommission(err error) *sdkerrors.Error {
	jsonData, _ := json.Marshal(
		map[string]string{
			"code":      getCodeString(CodeCalculateCommission),
			"codespace": DefaultCodespace,
			"desc":      fmt.Sprintf(err.Error()),
			"err":       err.Error(),
		},
	)
	return sdkerrors.New(
		DefaultCodespace,
		CodeCalculateCommission,
		string(jsonData),
	)
}

func ErrUpdateOnlyForCreator() *sdkerrors.Error {
	jsonData, _ := json.Marshal(
		map[string]string{
			"code":      getCodeString(CodeForbiddenUpdate),
			"codespace": DefaultCodespace,
			"desc":      fmt.Sprintf("updating allowed only for creator of coin"),
		},
	)
	return sdkerrors.New(
		DefaultCodespace,
		CodeForbiddenUpdate,
		string(jsonData),
	)
}

func ErrSameCoin() *sdkerrors.Error {
	jsonData, _ := json.Marshal(
		map[string]string{
			"code":      getCodeString(CodeSameCoins),
			"codespace": DefaultCodespace,
			"desc":      fmt.Sprintf("can't buy same coins"),
		},
	)
	return sdkerrors.New(
		DefaultCodespace,
		CodeSameCoins,
		string(jsonData),
	)
}

func ErrInsufficientFundsToSellAll() *sdkerrors.Error {
	jsonData, _ := json.Marshal(
		map[string]string{
			"code":      getCodeString(CodeInsufficientFundsToSellAll),
			"codespace": DefaultCodespace,
			"desc":      fmt.Sprintf("not enough coin to sell"),
		},
	)
	return sdkerrors.New(
		DefaultCodespace,
		CodeInsufficientFundsToSellAll,
		string(jsonData),
	)
}

func ErrTxBreaksVolumeLimit(volume string, limitVolume string) *sdkerrors.Error {
	jsonData, _ := json.Marshal(
		map[string]string{
			"code":        getCodeString(CodeTxBreaksVolumeLimit),
			"codespace":   DefaultCodespace,
			"desc":        fmt.Sprintf("tx breaks LimitVolume rule: %s > %s", volume, limitVolume),
			"volume":      volume,
			"limitVolume": limitVolume,
		},
	)
	return sdkerrors.New(
		DefaultCodespace,
		CodeTxBreaksVolumeLimit,
		string(jsonData),
	)
}

func ErrTxBreaksMinReserveRule(minCoinReserve string, reserve string) *sdkerrors.Error {
	jsonData, _ := json.Marshal(
		map[string]string{
			"code":           getCodeString(CodeTxBreaksMinReserveLimit),
			"codespace":      DefaultCodespace,
			"desc":           fmt.Sprintf("tx breaks MinReserveLimit rule: %s < %s", reserve, minCoinReserve),
			"reserve":        reserve,
			"minCoinReserve": minCoinReserve,
		},
	)
	return sdkerrors.New(
		DefaultCodespace,
		CodeTxBreaksMinReserveLimit,
		string(jsonData),
	)
}

func ErrMaximumValueToSellReached(amount, max string) *sdkerrors.Error {
	jsonData, _ := json.Marshal(
		map[string]string{
			"code":      getCodeString(CodeMaximumValueToSellReached),
			"codespace": DefaultCodespace,
			"desc":      fmt.Sprintf("wanted to sell maximum %s, but required to spend %s at the moment", max, amount),
			"max":       max,
			"amount":    amount,
		},
	)
	return sdkerrors.New(
		DefaultCodespace,
		CodeMaximumValueToSellReached,
		string(jsonData),
	)
}

func ErrMinimumValueToBuyReached(amount, min string) *sdkerrors.Error {
	jsonData, _ := json.Marshal(
		map[string]string{
			"code":      getCodeString(CodeMinimumValueToBuyReached),
			"codespace": DefaultCodespace,
			"desc":      fmt.Sprintf("wanted to buy minimum %s, but expected to receive %s at the moment", min, amount),
			"min":       min,
			"amount":    amount,
		},
	)
	return sdkerrors.New(
		DefaultCodespace,
		CodeMinimumValueToBuyReached,
		string(jsonData),
	)
}

func ErrUpdateBalance(account string, err string) *sdkerrors.Error {
	jsonData, _ := json.Marshal(
		map[string]string{
			"code":      getCodeString(CodeUpdateBalance),
			"codespace": DefaultCodespace,
			"desc":      fmt.Sprintf("unable to update balance of account %s: %s", account, err),
			"account":   account,
			"err":       err,
		},
	)
	return sdkerrors.New(
		DefaultCodespace,
		CodeUpdateBalance,
		string(jsonData),
	)
}

func ErrLimitVolumeBroken(volume string, limit string) *sdkerrors.Error {
	jsonData, _ := json.Marshal(
		map[string]string{
			"code":      getCodeString(CodeLimitVolumeBroken),
			"codespace": DefaultCodespace,
			"desc":      fmt.Sprintf("volume should be less than or equal the volume limit: %s > %s", volume, limit),
			"volume":    volume,
			"limit":     limit,
		},
	)
	return sdkerrors.New(
		DefaultCodespace,
		CodeLimitVolumeBroken,
		string(jsonData),
	)
}

func ErrInvalidAmount() *sdkerrors.Error {
	jsonData, _ := json.Marshal(
		map[string]string{
			"code":      getCodeString(CodeInvalidAmount),
			"codespace": DefaultCodespace,
			"desc":      fmt.Sprintf("amount should be greater than 0"),
		},
	)
	return sdkerrors.New(
		DefaultCodespace,
		CodeInvalidAmount,
		string(jsonData),
	)
}

// Redeem check

func ErrInvalidCheck(data string) *sdkerrors.Error {
	jsonData, _ := json.Marshal(
		map[string]string{
			"code":      getCodeString(CodeInvalidCheck),
			"codespace": DefaultCodespace,
			"desc":      fmt.Sprintf("amount should be greater than 0"),
		},
	)
	return sdkerrors.New(
		DefaultCodespace,
		CodeInvalidCheck,
		string(jsonData),
	)
}

func ErrInvalidProof(error string) *sdkerrors.Error {
	jsonData, _ := json.Marshal(
		map[string]string{
			"code":      getCodeString(CodeInvalidProof),
			"codespace": DefaultCodespace,
			"desc":      fmt.Sprintf("provided proof is invalid %s", error),
			"error":     error,
		},
	)
	return sdkerrors.New(
		DefaultCodespace,
		CodeInvalidProof,
		string(jsonData),
	)
}

func ErrInvalidPassphrase(error string) *sdkerrors.Error {
	jsonData, _ := json.Marshal(
		map[string]string{
			"code":      getCodeString(CodeInvalidAmount),
			"codespace": DefaultCodespace,
			"desc":      fmt.Sprintf("unable to create private key from passphrase: %s", error),
			"error":     error,
		},
	)
	return sdkerrors.New(
		DefaultCodespace,
		CodeInvalidPassphrase,
		string(jsonData),
	)
}

func ErrInvalidChainID(wanted string, issued string) *sdkerrors.Error {
	jsonData, _ := json.Marshal(
		map[string]string{
			"code":      getCodeString(CodeInvalidChainID),
			"codespace": DefaultCodespace,
			"desc":      fmt.Sprintf("wanted chain ID %s, but check is issued for chain with ID %s", wanted, issued),
			"wanted":    wanted,
			"issued":    issued,
		},
	)
	return sdkerrors.New(
		DefaultCodespace,
		CodeInvalidChainID,
		string(jsonData),
	)
}

func ErrInvalidNonce() *sdkerrors.Error {
	jsonData, _ := json.Marshal(
		map[string]string{
			"code":      getCodeString(CodeInvalidNonce),
			"codespace": DefaultCodespace,
			"desc":      fmt.Sprintf("nonce is too big (should be up to 16 bytes)"),
		},
	)
	return sdkerrors.New(
		DefaultCodespace,
		CodeInvalidNonce,
		string(jsonData),
	)
}

func ErrCheckExpired(block uint64) *sdkerrors.Error {
	jsonData, _ := json.Marshal(
		map[string]string{
			"code":      getCodeString(CodeCheckExpired),
			"codespace": DefaultCodespace,
			"desc":      fmt.Sprintf("check was expired at block %d", block),
			"block":     strconv.FormatInt(int64(block), 10),
		},
	)
	return sdkerrors.New(
		DefaultCodespace,
		CodeCheckExpired,
		string(jsonData),
	)
}

func ErrCheckRedeemed() *sdkerrors.Error {
	jsonData, _ := json.Marshal(
		map[string]string{
			"code":      getCodeString(CodeCheckRedeemed),
			"codespace": DefaultCodespace,
			"desc":      fmt.Sprintf("check was redeemed already"),
		},
	)
	return sdkerrors.New(
		DefaultCodespace,
		CodeCheckRedeemed,
		string(jsonData),
	)
}

func ErrUnableDecodeCheck(check string) *sdkerrors.Error {
	jsonData, _ := json.Marshal(
		map[string]string{
			"code":      getCodeString(CodeUnableDecodeCheck),
			"codespace": DefaultCodespace,
			"desc":      fmt.Sprintf("unable to decode check from base58 %s", check),
			"check":     check,
		},
	)
	return sdkerrors.New(
		DefaultCodespace,
		CodeUnableDecodeCheck,
		string(jsonData),
	)
}

func ErrUnableRPLEncodeCheck(error string) *sdkerrors.Error {
	jsonData, _ := json.Marshal(
		map[string]string{
			"code":      getCodeString(CodeUnableRPLEncodeCheck),
			"codespace": DefaultCodespace,
			"desc":      fmt.Sprintf("unable to RLP encode check receiver address: %s", error),
			"error":     error,
		},
	)
	return sdkerrors.New(
		DefaultCodespace,
		CodeUnableRPLEncodeCheck,
		string(jsonData),
	)
}

func ErrUnableSignCheck(error string) *sdkerrors.Error {
	jsonData, _ := json.Marshal(
		map[string]string{
			"code":      getCodeString(CodeUnableSignCheck),
			"codespace": DefaultCodespace,
			"desc":      fmt.Sprintf("unable to sign check receiver address by private key generated from passphrase: %s", error),
			"error":     error,
		},
	)
	return sdkerrors.New(
		DefaultCodespace,
		CodeUnableSignCheck,
		string(jsonData),
	)
}

func ErrUnableDecodeProof() *sdkerrors.Error {
	jsonData, _ := json.Marshal(
		map[string]string{
			"code":      getCodeString(CodeUnableDecodeProof),
			"codespace": DefaultCodespace,
			"desc":      fmt.Sprintf("unable to decode proof from base64"),
		},
	)
	return sdkerrors.New(
		DefaultCodespace,
		CodeUnableDecodeProof,
		string(jsonData),
	)
}

func ErrUnableRecoverAddress(error string) *sdkerrors.Error {
	jsonData, _ := json.Marshal(
		map[string]string{
			"code":      getCodeString(CodeUnableRecoverAddress),
			"codespace": DefaultCodespace,
			"desc":      fmt.Sprintf("unable to recover check issuer address: %s", error),
			"error":     error,
		},
	)
	return sdkerrors.New(
		DefaultCodespace,
		CodeUnableRecoverAddress,
		string(jsonData),
	)
}

func ErrUnableRecoverLockPkey(error string) *sdkerrors.Error {
	jsonData, _ := json.Marshal(
		map[string]string{
			"code":      getCodeString(CodeUnableRecoverLockPkey),
			"codespace": DefaultCodespace,
			"desc":      fmt.Sprintf("unable to recover lock public key from check: %s", error),
			"error":     error,
		},
	)
	return sdkerrors.New(
		DefaultCodespace,
		CodeUnableRecoverLockPkey,
		string(jsonData),
	)
}

// AccountKeys Errors

func ErrInvalidPkey() *sdkerrors.Error {
	jsonData, _ := json.Marshal(
		map[string]string{
			"code":      getCodeString(CodeInvalidPkey),
			"codespace": DefaultCodespace,
			"desc":      fmt.Sprintf("invalid private key"),
		},
	)
	return sdkerrors.New(
		DefaultCodespace,
		CodeInvalidPkey,
		string(jsonData),
	)
}

func ErrUnableRetriveArmoredPkey(name string, error string) *sdkerrors.Error {
	jsonData, _ := json.Marshal(
		map[string]string{
			"code":      getCodeString(CodeUnableRetriveArmoredPkey),
			"codespace": DefaultCodespace,
			"desc":      fmt.Sprintf("unable to retrieve armored private key for account %s: %s", name, error),
			"name":      name,
			"error":     error,
		},
	)
	return sdkerrors.New(
		DefaultCodespace,
		CodeUnableRetriveArmoredPkey,
		string(jsonData),
	)
}

func ErrUnableRetrivePkey(name string, error string) *sdkerrors.Error {
	jsonData, _ := json.Marshal(
		map[string]string{
			"code":      getCodeString(CodeUnableRetrivePkey),
			"codespace": DefaultCodespace,
			"desc":      fmt.Sprintf("unable to retrieve private key for account %s: %s", name, error),
			"name":      name,
			"error":     error,
		},
	)
	return sdkerrors.New(
		DefaultCodespace,
		CodeUnableRetrivePkey,
		string(jsonData),
	)
}

func ErrUnableRetriveSECPPkey(name string, algo string) *sdkerrors.Error {
	jsonData, _ := json.Marshal(
		map[string]string{
			"code":      getCodeString(CodeUnableRetriveSECPPkey),
			"codespace": DefaultCodespace,
			"desc":      fmt.Sprintf("unable to retrieve secp256k1 private key for account %s: %s private key retrieved instead", name, algo),
			"name":      name,
			"algo":      algo,
		},
	)
	return sdkerrors.New(
		DefaultCodespace,
		CodeUnableRetriveSECPPkey,
		string(jsonData),
	)
}

func getCodeString(code CodeType) string {
	return strconv.FormatInt(int64(code), 10)
}
