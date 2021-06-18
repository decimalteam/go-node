package keeper

import (
	"bitbucket.org/decimalteam/go-node/x/swap/internal/types"
	"testing"
)

func TestSetSwap(t *testing.T) {
	ctx, swapKeeper := CreateTestInput(t, false)

	msgburn := types.NewMsgBurn(Addr[0], "recipient", sdk.NewInt(2), "name", "symbol", "12", 4)

	hash, err := types.GetHash(msgburn.TransactionNumber, msgburn.TokenName, msgburn.TokenSymbol, msgburn.Amount, msgburn.Recipient, msgburn.DestChain)
	require.NoError(t, err)

	swapKeeper.SetSwapV2(ctx, hash)
	require.True(t, true, swapKeeper.HasSwapV2(ctx, hash))
}
