package types

import (
	"fmt"
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

	CodeInvalidDelegation CodeType = 300

	CodeEmptyPubKey                     CodeType = 400
	CodeValidatorPubKeyTypeNotSupported CodeType = 401

	CodeCoinReserveIsNotSufficient CodeType = 500

	CodeErrInvalidHistoricalInfo CodeType = 600

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

func ErrEmptyPubKey() *sdkerrors.Error {
	return sdkerrors.New(DefaultCodespace, CodeEmptyPubKey, `empty PubKey`)
}

func ErrEmptyValidatorAddr() *sdkerrors.Error {
	return sdkerrors.New(DefaultCodespace, CodeEmptyValidatorAddr, "empty validator address")
}

func ErrInvalidStruct() *sdkerrors.Error {
	return sdkerrors.New(DefaultCodespace, CodeInvalidStruct, "invalid struct")
}

func ErrValidatorOwnerExists() *sdkerrors.Error {
	return sdkerrors.New(DefaultCodespace, CodeValidatorAddressAlreadyExist, "validator already exist for this operator address, must use new validator operator address")
}

func ErrNoValidatorFound() *sdkerrors.Error {
	return sdkerrors.New(DefaultCodespace, CodeValidatorDoesNotExist, "validator does not exist for that address")
}

func ErrValidatorPubKeyExists() *sdkerrors.Error {
	return sdkerrors.New(DefaultCodespace, CodeValidatorPubKeyAlreadyExist, "validator already exist for this pubkey, must use new validator pubkey")
}

func ErrInsufficientShares() *sdkerrors.Error {
	return sdkerrors.New(DefaultCodespace, CodeInvalidDelegation, "insufficient delegation shares")
}

func ErrDelegatorShareExRateInvalid() *sdkerrors.Error {
	return sdkerrors.New(DefaultCodespace, CodeInvalidDelegation,
		"cannot delegate to validators with invalid (zero) ex-rate")
}

func ErrDelegatorStakeIsTooLow() *sdkerrors.Error {
	return sdkerrors.New(DefaultCodespace, CodeDelegatorStakeIsTooLow, "stake is too low")
}

func ErrUnbondingDelegationNotFound() *sdkerrors.Error {
	return sdkerrors.New(DefaultCodespace, CodeInvalidDelegation, "unbonding delegation not found")
}

func ErrBadDelegationAmount() *sdkerrors.Error {
	return sdkerrors.New(DefaultCodespace, CodeInvalidDelegation, "amount must be > 0")
}

func ErrNoDelegatorForAddress() *sdkerrors.Error {
	return sdkerrors.New(DefaultCodespace, CodeInvalidDelegation, "delegator does not contain this delegation")
}

func ErrNotEnoughDelegationShares(shares string) *sdkerrors.Error {
	return sdkerrors.New(DefaultCodespace, CodeInvalidDelegation, fmt.Sprintf("not enough shares only have %v", shares))
}

func ErrEmptyDelegatorAddr() *sdkerrors.Error {
	return sdkerrors.New(DefaultCodespace, CodeEmptyDelegatorAddress, "empty delegator address")
}

func ErrCoinReserveIsNotSufficient(reserve string, amount string) *sdkerrors.Error {
	return sdkerrors.New(DefaultCodespace, CodeCoinReserveIsNotSufficient, fmt.Sprintf("Coin reserve balance is not sufficient for transaction. Has: %s, required %s", reserve, amount))
}

func ErrNoDelegation() *sdkerrors.Error {
	return sdkerrors.New(DefaultCodespace, CodeInvalidDelegation, "no delegation for this (address, validator) pair")
}

func ErrCommissionNegative() *sdkerrors.Error {
	return sdkerrors.New(DefaultCodespace, CodeCommissionNegative, "commission must be positive")
}

func ErrCommissionHuge() *sdkerrors.Error {
	return sdkerrors.New(DefaultCodespace, CodeCommissionHuge, "commission cannot be more than 100%")
}

func ErrValidatorPubKeyTypeNotSupported() *sdkerrors.Error {
	return sdkerrors.New(DefaultCodespace, CodeValidatorPubKeyTypeNotSupported, "validator pubkey type is not supported")
}

func ErrInvalidHistoricalInfo() *sdkerrors.Error {
	return sdkerrors.New(DefaultCodespace, CodeErrInvalidHistoricalInfo, "invalid historical info")
}

func ErrNoHistoricalInfo() *sdkerrors.Error {
	return sdkerrors.New(DefaultCodespace, CodeErrNoHistoricalInfo, "no historical info found")
}

func ErrInsufficientCoinToPayCommission(commission string) *sdkerrors.Error {
	return sdkerrors.New(DefaultCodespace, CodeInsufficientCoinToPayCommission, fmt.Sprintf("Insufficient coin to pay commission: wanted = %s", commission))
}

func ErrInsufficientFunds(funds string) *sdkerrors.Error {
	return sdkerrors.New(DefaultCodespace, CodeInsufficientFunds, fmt.Sprintf("Insufficient funds: wanted = %s", funds))
}

func ErrUpdateBalance(err error) *sdkerrors.Error {
	return sdkerrors.New(DefaultCodespace, CodeUpdateBalanceError, err.Error())
}

func ErrCalculateCommission(err error) *sdkerrors.Error {
	return sdkerrors.New(DefaultCodespace, CodeErrCalculateCommission, err.Error())
}

func ErrCoinDoesNotExist(symbol string) *sdkerrors.Error {
	return sdkerrors.New(DefaultCodespace, CodeCoinDoesNotExist, fmt.Sprintf("coin %s does not exist", symbol))
}

func ErrValidatorAlreadyOnline() *sdkerrors.Error {
	return sdkerrors.New(DefaultCodespace, CodeValidatorAlreadyOnline, "validator already online")
}

func ErrValidatorAlreadyOffline() *sdkerrors.Error {
	return sdkerrors.New(DefaultCodespace, CodeValidatorAlreadyOffline, "validator already offline")
}

func ErrOwnerDoesNotOwnSubTokenID(owner sdk.AccAddress, subTokenID int64) *sdkerrors.Error {
	return sdkerrors.New(DefaultCodespace, CodeOwnerDoesNotOwnSubTokenID, fmt.Sprintf("the owner %s does not own the token with ID = %d", owner.String(), subTokenID))
}

func ErrInternal(err string) *sdkerrors.Error {
	return sdkerrors.New(DefaultCodespace, CodeInternalError, fmt.Sprintf("internal error: %s", err))
}
