package types

import (
	"bitbucket.org/decimalteam/go-node/utils/errors"
	"fmt"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// Local code type
type CodeType = uint32

const (
	// Default ante codespace
	DefaultRootCodespace string = "ante"

	CodeFeePayerAddressDoesNotExist    CodeType = 101
	CodeFeeLessThanCommission          CodeType = 102
	CodeFailedToSendCoins              CodeType = 103
	CodeInsufficientFundsToPayFee      CodeType = 104
	CodeInvalidFeeAmount               CodeType = 105
	CodeCoinDoesNotExist               CodeType = 106
	CodeNotStdTxType                   CodeType = 107
	CodeNotFeeTxType                   CodeType = 108
	CodeOutOfGas                       CodeType = 109
	CodeNotGasTxType                   CodeType = 110
	CodeInvalidAddressOfCreatedAccount CodeType = 111
	CodeUnableToFindCreatedAccount     CodeType = 112
)

func ErrFeePayerAddressDoesNotExist(feePayer string) *sdkerrors.Error {
	return errors.Encode(
		DefaultRootCodespace,
		CodeFeePayerAddressDoesNotExist,
		fmt.Sprintf("fee payer address does not exist: %s", feePayer),
		errors.NewParam("feePayer", feePayer),
	)
}

func ErrFeeLessThanCommission(feeInBaseCoin, commissionInBaseCoin string) *sdkerrors.Error {
	return errors.Encode(
		DefaultRootCodespace,
		CodeFeeLessThanCommission,
		fmt.Sprintf("insufficient funds to pay for fees; %s < %s", feeInBaseCoin, commissionInBaseCoin),
		errors.NewParam("feeInBaseCoin", feeInBaseCoin),
		errors.NewParam("commissionInBaseCoin", commissionInBaseCoin),
	)
}

func ErrFailedToSendCoins(err string) *sdkerrors.Error {
	return errors.Encode(
		DefaultRootCodespace,
		CodeFailedToSendCoins,
		fmt.Sprintf("failed to send coins: %s", err),
		errors.NewParam("error", err),
	)
}

func ErrInsufficientFundsToPayFee(coins, fee string) *sdkerrors.Error {
	return errors.Encode(
		DefaultRootCodespace,
		CodeInsufficientFundsToPayFee,
		fmt.Sprintf("insufficient funds to pay for fee; %s < %s", coins, fee),
		errors.NewParam("coins", coins),
		errors.NewParam("fee", fee),
	)
}

func ErrInvalidFeeAmount(fee string) *sdkerrors.Error {
	return errors.Encode(
		DefaultRootCodespace,
		CodeInvalidFeeAmount,
		"invalid fee amount",
		errors.NewParam("fee", fee),
	)
}

func ErrCoinDoesNotExist(feeDenom string) *sdkerrors.Error {
	return errors.Encode(
		DefaultRootCodespace,
		CodeCoinDoesNotExist,
		fmt.Sprintf("coin not exist: %s", feeDenom),
		errors.NewParam("feeDenom", feeDenom),
	)
}

func ErrNotStdTxType() *sdkerrors.Error {
	return errors.Encode(
		DefaultRootCodespace,
		CodeNotStdTxType,
		"Tx must be StdTx",
	)
}

func ErrNotFeeTxType() *sdkerrors.Error {
	return errors.Encode(
		DefaultRootCodespace,
		CodeNotFeeTxType,
		"x must be a FeeTx",
	)
}

func ErrNotGasTxType() *sdkerrors.Error {
	return errors.Encode(
		DefaultRootCodespace,
		CodeNotGasTxType,
		"Tx must be GasTx",
	)
}

func ErrOutOfGas(location, gasWanted, gasUsed string) *sdkerrors.Error {
	return errors.Encode(
		DefaultRootCodespace,
		CodeOutOfGas,
		fmt.Sprintf("out of gas in location: %v; gasWanted: %s, gasUsed: %s", location, gasWanted, gasUsed),
		errors.NewParam("location", location),
		errors.NewParam("gasWanted", gasWanted),
		errors.NewParam("gasUsed", gasUsed),
	)
}

func ErrInvalidAddressOfCreatedAccount() *sdkerrors.Error {
	return errors.Encode(
		DefaultRootCodespace,
		CodeInvalidAddressOfCreatedAccount,
		"invalid address of created account",
	)
}

func ErrUnableToFindCreatedAccount() *sdkerrors.Error {
	return errors.Encode(
		DefaultRootCodespace,
		CodeUnableToFindCreatedAccount,
		"unable to find created account",
	)
}
