package types

import (
	"fmt"

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

	CodeDelegatorStakeIsTooLow CodeType = 800
	CodeEmptyDelegatorAddress  CodeType = 801

	CodeInsufficientCoinToPayCommission CodeType = 900
	CodeInsufficientFunds               CodeType = 901
	CodeUpdateBalanceError              CodeType = 902
	CodeErrCalculateCommission          CodeType = 903
	CodeCoinDoesNotExist                CodeType = 904

	CodeInternalError CodeType = 1000
)

func ErrEmptyPubKey() *sdkerrors.Error {
	return sdkerrors.New(
		DefaultCodespace,
		CodeEmptyPubKey,
		"empty PubKey",
	)
}

func ErrEmptyValidatorAddr() *sdkerrors.Error {
	return sdkerrors.New(
		DefaultCodespace,
		CodeEmptyValidatorAddr,
		"empty validator address",
	)
}

func ErrInvalidStruct() *sdkerrors.Error {
	return sdkerrors.New(
		DefaultCodespace,
		CodeInvalidStruct,
		"invalid struct",
	)
}

func ErrAccountNotSet() *sdkerrors.Error {
	return sdkerrors.New(
		DefaultCodespace,
		CodeAccountNotSet,
		"pool accounts haven't been set",
	)
}

func ErrValidatorOwnerExists() *sdkerrors.Error {
	return sdkerrors.New(
		DefaultCodespace,
		CodeValidatorAddressAlreadyExist,
		"validator already exist for this operator address, must use new validator operator address",
	)
}

func ErrNoValidatorFound() *sdkerrors.Error {
	return sdkerrors.New(
		DefaultCodespace,
		CodeValidatorDoesNotExist,
		"validator does not exist for that address",
	)
}

func ErrValidatorPubKeyExists() *sdkerrors.Error {
	return sdkerrors.New(
		DefaultCodespace,
		CodeValidatorPubKeyAlreadyExist,
		"validator already exist for this pubkey, must use new validator pubkey",
	)
}

func ErrInsufficientShares() *sdkerrors.Error {
	return sdkerrors.New(
		DefaultCodespace,
		CodeInsufficientShares,
		"insufficient delegation shares",
	)
}

func ErrDelegatorShareExRateInvalid() *sdkerrors.Error {
	return sdkerrors.New(
		DefaultCodespace,
		CodeDelegatorShareExRateInvalid,
		"cannot delegate to validators with invalid (zero) ex-rate",
	)
}

func ErrUnbondingDelegationNotFound() *sdkerrors.Error {
	return sdkerrors.New(
		DefaultCodespace,
		CodeUnbondingDelegationNotFound,
		"unbonding delegation not found",
	)
}

func ErrBadDelegationAmount() *sdkerrors.Error {
	return sdkerrors.New(
		DefaultCodespace,
		CodeBadDelegationAmount,
		"amount must be > 0",
	)
}

func ErrNoDelegatorForAddress() *sdkerrors.Error {
	return sdkerrors.New(
		DefaultCodespace,
		CodeNoDelegatorForAddress,
		"delegator does not contain this address",
	)
}

func ErrNotEnoughDelegationShares(shares string) *sdkerrors.Error {
	return sdkerrors.New(
		DefaultCodespace,
		CodeNotEnoughDelegationShares,
		fmt.Sprintf("not enough shares only have %v", shares),
	)
}

func ErrCoinReserveIsNotSufficient(reserve string, amount string) *sdkerrors.Error {
	return sdkerrors.New(
		DefaultCodespace,
		CodeCoinReserveIsNotSufficient,
		fmt.Sprintf("Coin reserve balance is not sufficient for transaction. Has: %s, required %s", reserve, amount),
	)
}

func ErrNoDelegation() *sdkerrors.Error {
	return sdkerrors.New(
		DefaultCodespace,
		CodeInvalidDelegation,
		"no delegation for this (address, validator) pair",
	)
}

func ErrCommissionNegative() *sdkerrors.Error {
	return sdkerrors.New(
		DefaultCodespace,
		CodeCommissionNegative,
		"commission must be positive",
	)
}

func ErrCommissionHuge() *sdkerrors.Error {
	return sdkerrors.New(DefaultCodespace,
		CodeCommissionHuge,
		"commission cannot be more than 100%",
	)
}

func ErrValidatorAlreadyOnline() *sdkerrors.Error {
	return sdkerrors.New(
		DefaultCodespace,
		CodeValidatorAlreadyOnline,
		"validator already online",
	)
}

func ErrValidatorAlreadyOffline() *sdkerrors.Error {
	return sdkerrors.New(
		DefaultCodespace,
		CodeValidatorAlreadyOffline,
		"validator already offline",
	)
}

func ErrValidatorPubKeyTypeNotSupported(PKeyType string, AllowedTypes []string) *sdkerrors.Error {
	return sdkerrors.New(
		DefaultCodespace,
		CodeValidatorPubKeyTypeNotSupported,
		fmt.Sprintf("validator pubkey type is not supported. got: %s, valid: %s", PKeyType, AllowedTypes),
	)
}

func ErrInvalidHistoricalInfo() *sdkerrors.Error {
	return sdkerrors.New(
		DefaultCodespace,
		CodeErrInvalidHistoricalInfo,
		"invalid historical info",
	)
}

func ErrValidatorSetEmpty() *sdkerrors.Error {
	return sdkerrors.New(
		DefaultCodespace,
		CodeValidatorSetEmpty,
		"validator set is empty",
	)
}

func ErrValidatorSetNotSorted() *sdkerrors.Error {
	return sdkerrors.New(
		DefaultCodespace,
		CodeValidatorSetNotSorted,
		"validator set is not sorted by address",
	)
}

func ErrNoHistoricalInfo() *sdkerrors.Error {
	return sdkerrors.New(
		DefaultCodespace,
		CodeErrNoHistoricalInfo,
		"no historical info found",
	)
}

func ErrDelegatorStakeIsTooLow() *sdkerrors.Error {
	return sdkerrors.New(
		DefaultCodespace,
		CodeDelegatorStakeIsTooLow,
		"stake is too low",
	)
}

func ErrEmptyDelegatorAddr() *sdkerrors.Error {
	return sdkerrors.New(
		DefaultCodespace,
		CodeEmptyDelegatorAddress,
		"empty delegator address",
	)
}

func ErrInsufficientCoinToPayCommission(commission string) *sdkerrors.Error {
	return sdkerrors.New(
		DefaultCodespace,
		CodeInsufficientCoinToPayCommission,
		fmt.Sprintf("Insufficient coin to pay commission: wanted = %s", commission),
	)
}

func ErrInsufficientFunds(funds string) *sdkerrors.Error {
	return sdkerrors.New(
		DefaultCodespace,
		CodeInsufficientFunds,
		fmt.Sprintf("Insufficient funds: wanted = %s", funds),
	)
}

func ErrUpdateBalance(err error) *sdkerrors.Error {
	return sdkerrors.New(
		DefaultCodespace,
		CodeUpdateBalanceError,
		err.Error(),
	)
}

func ErrCalculateCommission(err error) *sdkerrors.Error {
	return sdkerrors.New(
		DefaultCodespace,
		CodeErrCalculateCommission,
		err.Error(),
	)
}

func ErrCoinDoesNotExist(symbol string) *sdkerrors.Error {
	return sdkerrors.New(
		DefaultCodespace,
		CodeCoinDoesNotExist,
		fmt.Sprintf("coin %s does not exist", symbol),
	)
}

func ErrInternal(err string) *sdkerrors.Error {
	return sdkerrors.New(
		DefaultCodespace,
		CodeInternalError,
		fmt.Sprintf("internal error: %s", err),
	)
}
