package nft

import (
	"bitbucket.org/decimalteam/go-node/utils/updates"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestUpdateBondDenom(t *testing.T) {
	ctx, _, NFTKeeper := createTestApp(t, false)

	*NFTKeeper.BaseDenom = "tdel"

	require.Equal(t, "tdel", *NFTKeeper.BaseDenom)

	ctx = ctx.WithBlockHeight(updates.Update10Block)

	BeginBlocker(ctx, NFTKeeper)

	require.Equal(t, "del", *NFTKeeper.BaseDenom)
}
