package types

import (
	"testing"

	"bitbucket.org/decimalteam/go-node/x/swap/internal/keeper"
)

func TestMsgBurn(t *testing.T) {
	ctx, swapKeeper := keeper.CreateTestInput(t, false)

	amount := sdk.NewInt(5)

	burn := types.NewMsgBurn(
		keeper.Addrs[0],
		"",
		amount,
		"",
		"",
		"",
		2,
	)

	swapKeeper.UnlockFunds()
}

func TestEcrecover(t *testing.T) {
	ctx, swapKeeper := keeper.CreateTestInput(t, false)

	h := NewHandler

}