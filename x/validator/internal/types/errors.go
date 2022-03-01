package types

import (
	"fmt"

	"bitbucket.org/decimalteam/go-node/utils/errors"
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
	return errors.Encode(
		DefaultCodespace,
		CodeEmptyValidatorAddr,
		"empty validator address",
	)
}

func ErrValidatorOwnerExists() *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeValidatorAddressAlreadyExist,
		"validator already exist for this operator address, must use new validator operator address",
	)
}

func ErrValidatorPubKeyExists() *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeValidatorPubKeyAlreadyExist,
		"validator already exist for this pubkey, must use new validator pubkey",
	)
}

func ErrNoValidatorFound() *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeValidatorDoesNotExist,
		"validator does not exist for that address",
	)
}

func ErrCommissionNegative() *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeCommissionNegative,
		"commission must be positive",
	)
}

func ErrCommissionHuge() *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeCommissionHuge,
		"commission cannot be more than 100%",
	)
}

func ErrValidatorAlreadyOnline() *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeValidatorAlreadyOnline,
		"validator already online",
	)
}

func ErrValidatorAlreadyOffline() *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeValidatorAlreadyOffline,
		"validator already offline",
	)
}

func ErrInvalidStruct() *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeInvalidStruct,
		"invalid struct",
	)
}

func ErrAccountNotSet() *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeAccountNotSet,
		"pool accounts haven't been set",
	)
}

func ErrNoDelegation() *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeInvalidDelegation,
		"no delegation for this (address, validator) pair",
	)
}

func ErrInsufficientShares() *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeInsufficientShares,
		"insufficient delegation shares",
	)
}

func ErrDelegatorShareExRateInvalid() *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeDelegatorShareExRateInvalid,
		"cannot delegate to validators with invalid (zero) ex-rate",
	)
}

func ErrUnbondingDelegationNotFound() *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeUnbondingDelegationNotFound,
		"unbonding delegation not found",
	)
}

func ErrBadDelegationAmount() *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeBadDelegationAmount,
		"amount must be > 0",
	)
}

func ErrNoDelegatorForAddress() *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeNoDelegatorForAddress,
		"delegator does not contain this address",
	)
}

func ErrNotEnoughDelegationShares(shares string) *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeNotEnoughDelegationShares,
		fmt.Sprintf("not enough shares only have %v", shares),
	)
}

func ErrEmptyPubKey() *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeEmptyPubKey,
		"empty PubKey",
	)
}

func ErrValidatorPubKeyTypeNotSupported(PKeyType string, AllowedTypes string) *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeValidatorPubKeyTypeNotSupported,
		fmt.Sprintf("validator pubkey type is not supported. got: %s, valid: %s", PKeyType, AllowedTypes),
	)
}

func ErrCoinReserveIsNotSufficient(reserve string, amount string) *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeCoinReserveIsNotSufficient,
		fmt.Sprintf("Coin reserve balance is not sufficient for transaction. Has: %s, required %s", reserve, amount),
	)
}

func ErrInvalidHistoricalInfo() *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeErrInvalidHistoricalInfo,
		"invalid historical info",
	)
}

func ErrValidatorSetEmpty() *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeValidatorSetEmpty,
		"validator set is empty",
	)
}

func ErrValidatorSetNotSorted() *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeValidatorSetNotSorted,
		"validator set is not sorted by address",
	)
}

func ErrNoHistoricalInfo(height string) *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeErrNoHistoricalInfo,
		fmt.Sprintf("no historical info found: %s", height),
	)
}

func ErrDelegatorStakeIsTooLow() *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeDelegatorStakeIsTooLow,
		"stake is too low",
	)
}

func ErrEmptyDelegatorAddr() *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeEmptyDelegatorAddress,
		"empty delegator address",
	)
}

func ErrInsufficientCoinToPayCommission(commission string) *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeInsufficientCoinToPayCommission,
		fmt.Sprintf("Insufficient coin to pay commission: wanted = %s", commission),
	)
}

func ErrInsufficientFunds(funds string) *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeInsufficientFunds,
		fmt.Sprintf("Insufficient funds: wanted = %s", funds),
	)
}

func ErrUpdateBalance(error string) *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeUpdateBalanceError,
		fmt.Sprintf("update balance error: %s", error),
	)
}

func ErrCalculateCommission(error string) *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeErrCalculateCommission,
		fmt.Sprintf("calculate commision error: %s", error),
	)
}

func ErrCoinDoesNotExist(symbol string) *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeCoinDoesNotExist,
		fmt.Sprintf("coin %s does not exist", symbol),
	)
}

func ErrInternal(error string) *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeInternalError,
		fmt.Sprintf("internal error: %s", error),
	)
}

func ErrOwnerDoesNotOwnSubTokenID(owner string, subTokenID string) *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeOwnerDoesNotOwnSubTokenID,
		fmt.Sprintf("the owner %s does not own the token with ID = %s", owner, subTokenID),
	)
}
