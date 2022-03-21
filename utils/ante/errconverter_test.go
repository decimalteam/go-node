package ante

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/stretchr/testify/require"
)

type checkFunctionType func(oldCtx sdk.Context, newCtx sdk.Context, err error)

func TestErrorConverting(t *testing.T) {
	testSuite := []struct {
		functionToWrap sdk.AnteHandler
		resultChecker  checkFunctionType
	}{
		{ // 1 no error
			func(ctx sdk.Context, tx sdk.Tx, simulate bool) (newCtx sdk.Context, err error) {
				return ctx.WithBlockHeight(ctx.BlockHeight() + 1), nil
			},
			func(oldCtx sdk.Context, newCtx sdk.Context, err error) {
				require.True(t, (oldCtx.BlockHeight()+1 == newCtx.BlockHeight()) && err == nil, "expect no changes")
			},
		},
		{ // 2 simple error
			func(ctx sdk.Context, tx sdk.Tx, simulate bool) (newCtx sdk.Context, err error) {
				return ctx, sdkerrors.ErrInvalidPubKey
			},
			func(oldCtx sdk.Context, newCtx sdk.Context, err error) {
				cnvErr, ok := err.(*sdkerrors.Error)
				require.True(t, ok, "expect error is sdk error")
				require.Equal(t, "sdk", cnvErr.Codespace(), "wrong codespace")
				require.Equal(t, sdkerrors.ErrInvalidPubKey.ABCICode()+100, cnvErr.ABCICode(), "wrong error code")
				require.Equal(t, "{\"description\":\"invalid pubkey\",\"codespace\":\"sdk\"}", cnvErr.Error(), "wrong error message")
			},
		},
		{ // 3 wrapped error
			func(ctx sdk.Context, tx sdk.Tx, simulate bool) (newCtx sdk.Context, err error) {
				return ctx, sdkerrors.Wrap(sdkerrors.ErrInvalidPubKey, "this is message")
			},
			func(oldCtx sdk.Context, newCtx sdk.Context, err error) {
				cnvErr, ok := err.(*sdkerrors.Error)
				require.True(t, ok, "expect error is sdk error")
				require.Equal(t, "sdk", cnvErr.Codespace(), "wrong codespace")
				require.Equal(t, sdkerrors.ErrInvalidPubKey.ABCICode()+100, cnvErr.ABCICode(), "wrong error code")
				require.Equal(t, "{\"description\":\"invalid pubkey: this is message\",\"codespace\":\"sdk\"}", cnvErr.Error(), "wrong error message")
			},
		},
	}
	for _, suite := range testSuite {
		newFunc := sdkErrorConverter(suite.functionToWrap)
		oldCtx := sdk.Context{}.WithBlockHeight(10)
		newCtx, err := newFunc(oldCtx, nil, false)
		suite.resultChecker(oldCtx, newCtx, err)
	}

}
