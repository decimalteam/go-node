package keeper

import (
	"bitbucket.org/decimalteam/go-node/x/validator/internal/types"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParams(t *testing.T) {
	ctx, _, keeper, _, _, _ := CreateTestInput(t, false, 0)
	expParams := types.DefaultParams()

	//check that the empty keeper loads the default
	resParams := keeper.GetParams(ctx)
	require.True(t, expParams.Equal(resParams))

	//modify a params, save, and retrieve
	expParams.MaxValidators = 777
	keeper.SetParams(ctx, expParams)
	resParams = keeper.GetParams(ctx)
	require.True(t, expParams.Equal(resParams))
}
