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

	CodeEmptyValidatorAddr CodeType = 101
	CodeInvalidValidator   CodeType = 102

	CodeInvalidStruct CodeType = 201
	CodeInvalidInput  CodeType = 202

	CodeInvalidDelegation CodeType = 301

	CodeEmptyPubKey CodeType = 401

	CodeCoinReserveIsNotSufficient CodeType = 501
)

func ErrEmptyPubKey(codespace string) *sdkerrors.Error {
	return sdkerrors.New(codespace, CodeEmptyPubKey, `empty PubKey`)
}

func ErrEmptyValidatorAddr(codespace string) *sdkerrors.Error {
	return sdkerrors.New(codespace, CodeEmptyValidatorAddr, "empty validator address")
}

func ErrInvalidStruct(codespace string) *sdkerrors.Error {
	return sdkerrors.New(codespace, CodeInvalidStruct, "invalid struct")
}

func ErrValidatorOwnerExists(codespace string) *sdkerrors.Error {
	return sdkerrors.New(codespace, CodeInvalidValidator, "validator already exist for this operator address, must use new validator operator address")
}

func ErrNoValidatorFound(codespace string) *sdkerrors.Error {
	return sdkerrors.New(codespace, CodeInvalidValidator, "validator does not exist for that address")
}

func ErrValidatorPubKeyExists(codespace string) *sdkerrors.Error {
	return sdkerrors.New(codespace, CodeInvalidValidator, "validator already exist for this pubkey, must use new validator pubkey")
}

func ErrInsufficientShares(codespace string) *sdkerrors.Error {
	return sdkerrors.New(codespace, CodeInvalidDelegation, "insufficient delegation shares")
}

func ErrDelegatorShareExRateInvalid(codespace string) *sdkerrors.Error {
	return sdkerrors.New(codespace, CodeInvalidDelegation,
		"cannot delegate to validators with invalid (zero) ex-rate")
}

func ErrNoUnbondingDelegation(codespace string) *sdkerrors.Error {
	return sdkerrors.New(codespace, CodeInvalidDelegation, "no unbonding delegation found")
}

func ErrBadDelegationAddr(codespace string) *sdkerrors.Error {
	return sdkerrors.New(codespace, CodeInvalidInput, "unexpected address length for this (address, validator) pair")
}

func ErrNoDelegatorForAddress(codespace string) *sdkerrors.Error {
	return sdkerrors.New(codespace, CodeInvalidDelegation, "delegator does not contain this delegation")
}

func ErrNotEnoughDelegationShares(codespace string, shares string) *sdkerrors.Error {
	return sdkerrors.New(codespace, CodeInvalidDelegation, fmt.Sprintf("not enough shares only have %v", shares))
}

func ErrNilDelegatorAddr(codespace string) *sdkerrors.Error {
	return sdkerrors.New(codespace, CodeInvalidInput, "delegator address is nil")
}

func ErrCoinReserveIsNotSufficient(codespace string, reserve string, amount string) *sdkerrors.Error {
	return sdkerrors.New(codespace, CodeCoinReserveIsNotSufficient, fmt.Sprintf("Coin reserve balance is not sufficient for transaction. Has: %s, required %s", reserve, amount))
}
