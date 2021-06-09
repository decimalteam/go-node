package types

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// Local code type
type CodeType = uint32

const (
	// Default validator codespace
	DefaultCodespace string = ModuleName

	CodeEmptyValidatorAddr           CodeType = 100
	CodeValidatorAddressAlreadyExist CodeType = 101
	CodeValidatorPubKeyAlreadyExist  CodeType = 102
	CodeValidatorDoesNotExist        CodeType = 103
	CodeCommissionNegative           CodeType = 104
	CodeCommissionHuge               CodeType = 105
	CodeValidatorAlreadyOnline       CodeType = 106
	CodeValidatorAlreadyOffline      CodeType = 107

	CodeInvalidStruct CodeType = 200
	CodeAccountNotSet CodeType = 201

	CodeInvalidDelegation           CodeType = 300
	CodeInsufficientShares          CodeType = 301
	CodeDelegatorShareExRateInvalid CodeType = 302
	CodeUnbondingDelegationNotFound CodeType = 303
	CodeBadDelegationAmount         CodeType = 304
	CodeNoDelegatorForAddress       CodeType = 305
	CodeNotEnoughDelegationShares   CodeType = 306

	CodeEmptyPubKey                     CodeType = 400
	CodeValidatorPubKeyTypeNotSupported CodeType = 401

	CodeCoinReserveIsNotSufficient CodeType = 500

	CodeErrInvalidHistoricalInfo CodeType = 600
	CodeValidatorSetEmpty        CodeType = 601
	CodeValidatorSetNotSorted    CodeType = 602

	CodeErrNoHistoricalInfo CodeType = 700

	CodeDelegatorStakeIsTooLow    CodeType = 800
	CodeEmptyDelegatorAddress     CodeType = 801
	CodeOwnerDoesNotOwnSubTokenID CodeType = 802

	CodeInsufficientCoinToPayCommission CodeType = 900
	CodeInsufficientFunds               CodeType = 901
	CodeUpdateBalanceError              CodeType = 902
	CodeErrCalculateCommission          CodeType = 903
	CodeCoinDoesNotExist                CodeType = 904

	CodeInternalError CodeType = 1000
)

func ErrEmptyValidatorAddr() *sdkerrors.Error {
	jsonData, _ := json.Marshal(
		map[string]string{
			"code":      getCodeString(CodeEmptyValidatorAddr),
			"codespace": DefaultCodespace,
			"desc":      "empty validator address",
		},
	)
	return sdkerrors.New(
		DefaultCodespace,
		CodeEmptyValidatorAddr,
		string(jsonData),
	)
}

func ErrValidatorOwnerExists() *sdkerrors.Error {
	jsonData, _ := json.Marshal(
		map[string]string{
			"code":      getCodeString(CodeValidatorAddressAlreadyExist),
			"codespace": DefaultCodespace,
			"desc":      "validator already exist for this operator address, must use new validator operator address",
		},
	)
	return sdkerrors.New(
		DefaultCodespace,
		CodeValidatorAddressAlreadyExist,
		string(jsonData),
	)
}

func ErrValidatorPubKeyExists() *sdkerrors.Error {
	jsonData, _ := json.Marshal(
		map[string]string{
			"code":      getCodeString(CodeValidatorPubKeyAlreadyExist),
			"codespace": DefaultCodespace,
			"desc":      "validator already exist for this pubkey, must use new validator pubkey",
		},
	)
	return sdkerrors.New(
		DefaultCodespace,
		CodeValidatorPubKeyAlreadyExist,
		string(jsonData),
	)
}

func ErrNoValidatorFound() *sdkerrors.Error {
	jsonData, _ := json.Marshal(
		map[string]string{
			"code":      getCodeString(CodeValidatorDoesNotExist),
			"codespace": DefaultCodespace,
			"desc":      "validator does not exist for that address",
		},
	)
	return sdkerrors.New(
		DefaultCodespace,
		CodeValidatorDoesNotExist,
		string(jsonData),
	)
}

func ErrCommissionNegative() *sdkerrors.Error {
	jsonData, _ := json.Marshal(
		map[string]string{
			"code":      getCodeString(CodeCommissionNegative),
			"codespace": DefaultCodespace,
			"desc":      "commission must be positive",
		},
	)
	return sdkerrors.New(
		DefaultCodespace,
		CodeCommissionNegative,
		string(jsonData),
	)
}

func ErrCommissionHuge() *sdkerrors.Error {
	jsonData, _ := json.Marshal(
		map[string]string{
			"code":      getCodeString(CodeCommissionHuge),
			"codespace": DefaultCodespace,
			"desc":      "commission cannot be more than 100%",
		},
	)
	return sdkerrors.New(
		DefaultCodespace,
		CodeCommissionHuge,
		string(jsonData),
	)
}

func ErrValidatorAlreadyOnline() *sdkerrors.Error {
	jsonData, _ := json.Marshal(
		map[string]string{
			"code":      getCodeString(CodeValidatorAlreadyOnline),
			"codespace": DefaultCodespace,
			"desc":      "validator already online",
		},
	)
	return sdkerrors.New(
		DefaultCodespace,
		CodeValidatorAlreadyOnline,
		string(jsonData),
	)
}

func ErrValidatorAlreadyOffline() *sdkerrors.Error {
	jsonData, _ := json.Marshal(
		map[string]string{
			"code":      getCodeString(CodeValidatorAlreadyOffline),
			"codespace": DefaultCodespace,
			"desc":      "validator already offline",
		},
	)
	return sdkerrors.New(
		DefaultCodespace,
		CodeValidatorAlreadyOffline,
		string(jsonData),
	)
}

func ErrInvalidStruct() *sdkerrors.Error {
	jsonData, _ := json.Marshal(
		map[string]string{
			"code":      getCodeString(CodeInvalidStruct),
			"codespace": DefaultCodespace,
			"desc":      "invalid struct",
		},
	)
	return sdkerrors.New(
		DefaultCodespace,
		CodeInvalidStruct,
		string(jsonData),
	)
}

func ErrAccountNotSet() *sdkerrors.Error {
	jsonData, _ := json.Marshal(
		map[string]string{
			"code":      getCodeString(CodeAccountNotSet),
			"codespace": DefaultCodespace,
			"desc":      "pool accounts haven't been set",
		},
	)
	return sdkerrors.New(
		DefaultCodespace,
		CodeAccountNotSet,
		string(jsonData),
	)
}

func ErrNoDelegation() *sdkerrors.Error {
	jsonData, _ := json.Marshal(
		map[string]string{
			"code":      getCodeString(CodeInvalidDelegation),
			"codespace": DefaultCodespace,
			"desc":      "no delegation for this (address, validator) pair",
		},
	)
	return sdkerrors.New(
		DefaultCodespace,
		CodeInvalidDelegation,
		string(jsonData),
	)
}

func ErrInsufficientShares() *sdkerrors.Error {
	jsonData, _ := json.Marshal(
		map[string]string{
			"code":      getCodeString(CodeInsufficientShares),
			"codespace": DefaultCodespace,
			"desc":      "insufficient delegation shares",
		},
	)
	return sdkerrors.New(
		DefaultCodespace,
		CodeInsufficientShares,
		string(jsonData),
	)
}

func ErrDelegatorShareExRateInvalid() *sdkerrors.Error {
	jsonData, _ := json.Marshal(
		map[string]string{
			"code":      getCodeString(CodeDelegatorShareExRateInvalid),
			"codespace": DefaultCodespace,
			"desc":      "cannot delegate to validators with invalid (zero) ex-rate",
		},
	)
	return sdkerrors.New(
		DefaultCodespace,
		CodeDelegatorShareExRateInvalid,
		string(jsonData),
	)
}

func ErrUnbondingDelegationNotFound() *sdkerrors.Error {
	jsonData, _ := json.Marshal(
		map[string]string{
			"code":      getCodeString(CodeUnbondingDelegationNotFound),
			"codespace": DefaultCodespace,
			"desc":      "unbonding delegation not found",
		},
	)
	return sdkerrors.New(
		DefaultCodespace,
		CodeUnbondingDelegationNotFound,
		string(jsonData),
	)
}

func ErrBadDelegationAmount() *sdkerrors.Error {
	jsonData, _ := json.Marshal(
		map[string]string{
			"code":      getCodeString(CodeBadDelegationAmount),
			"codespace": DefaultCodespace,
			"desc":      "amount must be > 0",
		},
	)
	return sdkerrors.New(
		DefaultCodespace,
		CodeBadDelegationAmount,
		string(jsonData),
	)
}

func ErrNoDelegatorForAddress() *sdkerrors.Error {
	jsonData, _ := json.Marshal(
		map[string]string{
			"code":      getCodeString(CodeNoDelegatorForAddress),
			"codespace": DefaultCodespace,
			"desc":      "delegator does not contain this address",
		},
	)
	return sdkerrors.New(
		DefaultCodespace,
		CodeNoDelegatorForAddress,
		string(jsonData),
	)
}

func ErrNotEnoughDelegationShares(shares string) *sdkerrors.Error {
	jsonData, _ := json.Marshal(
		map[string]string{
			"code":      getCodeString(CodeNotEnoughDelegationShares),
			"codespace": DefaultCodespace,
			"desc":      fmt.Sprintf("not enough shares only have %v", shares),
			"shares":    shares,
		},
	)
	return sdkerrors.New(
		DefaultCodespace,
		CodeNotEnoughDelegationShares,
		string(jsonData),
	)
}

func ErrEmptyPubKey() *sdkerrors.Error {
	jsonData, _ := json.Marshal(
		map[string]string{
			"code":      getCodeString(CodeEmptyPubKey),
			"codespace": DefaultCodespace,
			"desc":      "empty PubKey",
		},
	)
	return sdkerrors.New(
		DefaultCodespace,
		CodeEmptyPubKey,
		string(jsonData),
	)
}

func ErrValidatorPubKeyTypeNotSupported(PKeyType string, AllowedTypes []string) *sdkerrors.Error {
	jsonData, _ := json.Marshal(
		map[string]string{
			"code":         getCodeString(CodeValidatorPubKeyTypeNotSupported),
			"codespace":    DefaultCodespace,
			"desc":         fmt.Sprintf("validator pubkey type is not supported. got: %s, valid: %s", PKeyType, strings.Join(AllowedTypes, ",")),
			"PKeyType":     PKeyType,
			"AllowedTypes": strings.Join(AllowedTypes, ","),
		},
	)
	return sdkerrors.New(
		DefaultCodespace,
		CodeValidatorPubKeyTypeNotSupported,
		string(jsonData),
	)
}

func ErrCoinReserveIsNotSufficient(reserve string, amount string) *sdkerrors.Error {
	jsonData, _ := json.Marshal(
		map[string]string{
			"code":      getCodeString(CodeCoinReserveIsNotSufficient),
			"codespace": DefaultCodespace,
			"desc":      fmt.Sprintf("Coin reserve balance is not sufficient for transaction. Has: %s, required %s", reserve, amount),
			"reserve":   reserve,
			"amount":    amount,
		},
	)
	return sdkerrors.New(
		DefaultCodespace,
		CodeCoinReserveIsNotSufficient,
		string(jsonData),
	)
}

func ErrInvalidHistoricalInfo() *sdkerrors.Error {
	jsonData, _ := json.Marshal(
		map[string]string{
			"code":      getCodeString(CodeErrInvalidHistoricalInfo),
			"codespace": DefaultCodespace,
			"desc":      "invalid historical info",
		},
	)
	return sdkerrors.New(
		DefaultCodespace,
		CodeErrInvalidHistoricalInfo,
		string(jsonData),
	)
}

func ErrValidatorSetEmpty() *sdkerrors.Error {
	jsonData, _ := json.Marshal(
		map[string]string{
			"code":      getCodeString(CodeValidatorSetEmpty),
			"codespace": DefaultCodespace,
			"desc":      "validator set is empty",
		},
	)
	return sdkerrors.New(
		DefaultCodespace,
		CodeValidatorSetEmpty,
		string(jsonData),
	)
}

func ErrValidatorSetNotSorted() *sdkerrors.Error {
	jsonData, _ := json.Marshal(
		map[string]string{
			"code":      getCodeString(CodeValidatorSetEmpty),
			"codespace": DefaultCodespace,
			"desc":      "validator set is not sorted by address",
		},
	)
	return sdkerrors.New(
		DefaultCodespace,
		CodeValidatorSetNotSorted,
		string(jsonData),
	)
}

func ErrNoHistoricalInfo(height int64) *sdkerrors.Error {
	jsonData, _ := json.Marshal(
		map[string]string{
			"code":      getCodeString(CodeErrNoHistoricalInfo),
			"codespace": DefaultCodespace,
			"desc":      fmt.Sprintf("no historical info found: %d", height),
			"height":    strconv.FormatInt(int64(height), 10),
		},
	)
	return sdkerrors.New(
		DefaultCodespace,
		CodeErrNoHistoricalInfo,
		string(jsonData),
	)
}

func ErrDelegatorStakeIsTooLow() *sdkerrors.Error {
	jsonData, _ := json.Marshal(
		map[string]string{
			"code":      getCodeString(CodeDelegatorStakeIsTooLow),
			"codespace": DefaultCodespace,
			"desc":      "stake is too low",
		},
	)
	return sdkerrors.New(
		DefaultCodespace,
		CodeDelegatorStakeIsTooLow,
		string(jsonData),
	)
}

func ErrEmptyDelegatorAddr() *sdkerrors.Error {
	jsonData, _ := json.Marshal(
		map[string]string{
			"code":      getCodeString(CodeEmptyDelegatorAddress),
			"codespace": DefaultCodespace,
			"desc":      "empty delegator address",
		},
	)
	return sdkerrors.New(
		DefaultCodespace,
		CodeEmptyDelegatorAddress,
		string(jsonData),
	)
}

func ErrInsufficientCoinToPayCommission(commission string) *sdkerrors.Error {
	jsonData, _ := json.Marshal(
		map[string]string{
			"code":       getCodeString(CodeInsufficientCoinToPayCommission),
			"codespace":  DefaultCodespace,
			"desc":       fmt.Sprintf("Insufficient coin to pay commission: wanted = %s", commission),
			"commission": commission,
		},
	)
	return sdkerrors.New(
		DefaultCodespace,
		CodeInsufficientCoinToPayCommission,
		string(jsonData),
	)
}

func ErrInsufficientFunds(funds string) *sdkerrors.Error {
	jsonData, _ := json.Marshal(
		map[string]string{
			"code":      getCodeString(CodeInsufficientFunds),
			"codespace": DefaultCodespace,
			"desc":      fmt.Sprintf("Insufficient funds: wanted = %s", funds),
			"funds":     funds,
		},
	)
	return sdkerrors.New(
		DefaultCodespace,
		CodeInsufficientFunds,
		string(jsonData),
	)
}

func ErrUpdateBalance(error error) *sdkerrors.Error {
	jsonData, _ := json.Marshal(
		map[string]string{
			"code":      getCodeString(CodeUpdateBalanceError),
			"codespace": DefaultCodespace,
			"desc":      fmt.Sprintf("update balance error: %s", error.Error()),
			"error":     error.Error(),
		},
	)
	return sdkerrors.New(
		DefaultCodespace,
		CodeUpdateBalanceError,
		string(jsonData),
	)
}

func ErrCalculateCommission(error error) *sdkerrors.Error {
	jsonData, _ := json.Marshal(
		map[string]string{
			"code":      getCodeString(CodeErrCalculateCommission),
			"codespace": DefaultCodespace,
			"desc":      fmt.Sprintf("calculate commision error: %s", error.Error()),
			"error":     error.Error(),
		},
	)
	return sdkerrors.New(
		DefaultCodespace,
		CodeErrCalculateCommission,
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

func ErrInternal(error string) *sdkerrors.Error {
	jsonData, _ := json.Marshal(
		map[string]string{
			"code":      getCodeString(CodeInternalError),
			"codespace": DefaultCodespace,
			"desc":      fmt.Sprintf("internal error: %s", error),
			"error":     error,
		},
	)
	return sdkerrors.New(
		DefaultCodespace,
		CodeInternalError,
		string(jsonData),
	)
}

func ErrOwnerDoesNotOwnSubTokenID(owner sdk.AccAddress, subTokenID int64) *sdkerrors.Error {
	return sdkerrors.New(DefaultCodespace, CodeOwnerDoesNotOwnSubTokenID, fmt.Sprintf("the owner %s does not own the token with ID = %d", owner.String(), subTokenID))
}

func getCodeString(code CodeType) string {
	return strconv.FormatInt(int64(code), 10)
}


