package keeper

import (
	"bitbucket.org/decimalteam/go-node/x/swap/internal/types"
	"reflect"
	"testing"
)

func TestHasDestChain(t *testing.T) {
	ctx, swapKeeper := CreateTestInput(t, false)

	chain := types.NewChain("chain1", true)
	chainNumber := 1
	has := swapKeeper.HasChain(ctx, chainNumber)
	require.False(t, has)

	swapKeeper.SetChain(ctx, chainNumber, chain)

	has = swapKeeper.HasChain(ctx, chainNumber)
	require.True(t, has)
}

func TestSetDestChain(t *testing.T) {
	ctx, swapKeeper := CreateTestInput(t, false)

	chain := types.NewChain("chain1", true)
	chainNumber := 1

	swapKeeper.SetChain(ctx, chainNumber, chain)

	require.True(t, swapKeeper.HasChain(ctx, chainNumber))

	ch, found := swapKeeper.GetChain(ctx, chainNumber)
	require.True(t, found)
	require.True(t, reflect.DeepEqual(chain, ch))
}

func TestGetDestChainName(t *testing.T) {

}

func TestDeleteDestChain(t *testing.T) {

}
