package swap

import (
	"bitbucket.org/decimalteam/go-node/x/swap/internal/keeper"
	"bitbucket.org/decimalteam/go-node/x/swap/internal/types"
	"testing"
)

func TestInitGenesis(t *testing.T) {
	ctx, swapKeeper := keeper.CreateTestInput(t, false)

	genesisState := DefaultGenesisState()
	require.Nil(t, genesisState.Swaps)
	require.Equal(t, genesisState.Params.LockedTimeOut, types.DefaultLockedTimeOut)
	require.Equal(t, genesisState.Params.LockedTimeOut, types.DefaultLockedTimeIn)

	//genesisState.Params = types.NewParams(time.Duration(30), time.Duration(15))
	//genesisState.Swaps = types.Swaps{types.NewSwapV2()}

	InitGenesis(ctx, swapKeeper, genesisState)

	//swaps :=  swapKeeper.GetAllSwaps(ctx)
	//params := swapKeeper.GetParams(ctx)
	//
	//require.Equal(t, 1, len(swaps))
}