package ante

import (
	"errors"

	decerrors "bitbucket.org/decimalteam/go-node/utils/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// List of errors from Cosmos SDK cosmos-sdk@v0.39.3/types/errors/errors.go
var convertList = []struct {
	sdkError *sdkerrors.Error
	code     uint32
}{
	// ErrTxDecode is returned if we cannot parse a transaction
	{sdkerrors.ErrTxDecode, 102},

	// ErrInvalidSequence is used the sequence number (nonce) is incorrect
	// for the signature
	{sdkerrors.ErrInvalidSequence, 103},

	// ErrUnauthorized is used whenever a request without sufficient
	// authorization is handled.
	{sdkerrors.ErrUnauthorized, 104},

	// ErrInsufficientFunds is used when the account cannot pay requested amount.
	{sdkerrors.ErrInsufficientFunds, 105},

	// ErrUnknownRequest to doc
	{sdkerrors.ErrUnknownRequest, 106},

	// ErrInvalidAddress to doc
	{sdkerrors.ErrInvalidAddress, 107},

	// ErrInvalidPubKey to doc
	{sdkerrors.ErrInvalidPubKey, 108},

	// ErrUnknownAddress to doc
	{sdkerrors.ErrUnknownAddress, 109},

	// ErrInvalidCoins to doc
	{sdkerrors.ErrInvalidCoins, 110},

	// ErrOutOfGas to doc
	{sdkerrors.ErrOutOfGas, 111},

	// ErrMemoTooLarge to doc
	{sdkerrors.ErrMemoTooLarge, 112},

	// ErrInsufficientFee to doc
	{sdkerrors.ErrInsufficientFee, 113},

	// ErrTooManySignatures to doc
	{sdkerrors.ErrTooManySignatures, 114},

	// ErrNoSignatures to doc
	{sdkerrors.ErrNoSignatures, 115},

	// ErrJSONMarshal defines an ABCI typed JSON marshalling error
	{sdkerrors.ErrJSONMarshal, 116},

	// ErrJSONUnmarshal defines an ABCI typed JSON unmarshalling error
	{sdkerrors.ErrJSONUnmarshal, 117},

	// ErrInvalidRequest defines an ABCI typed error where the request contains
	// invalid data.
	{sdkerrors.ErrInvalidRequest, 118},

	// ErrTxInMempoolCache defines an ABCI typed error where a tx already exists
	// in the mempool.
	{sdkerrors.ErrTxInMempoolCache, 119},

	// ErrMempoolIsFull defines an ABCI typed error where the mempool is full.
	{sdkerrors.ErrMempoolIsFull, 120},

	// ErrTxTooLarge defines an ABCI typed error where tx is too large.
	{sdkerrors.ErrTxTooLarge, 121},
}

// Function sdkConverter decorate AnteHandler and convert errors from
// Cosmos SDK to appropriate error with codespace, code and message
func sdkErrorConverter(handler sdk.AnteHandler) sdk.AnteHandler {
	return func(ctx sdk.Context, tx sdk.Tx, simulate bool) (newCtx sdk.Context, err error) {
		newCtx, err = handler(ctx, tx, simulate)
		if err == nil {
			return newCtx, err
		}
		for _, convErr := range convertList {
			if errors.Is(err, convErr.sdkError) {
				return newCtx, decerrors.Encode("sdk", convErr.code, err.Error())
			}
		}
		return newCtx, err
	}
}
