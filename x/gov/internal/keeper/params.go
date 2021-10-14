package keeper

import (
	"bitbucket.org/decimalteam/go-node/utils/updates"
	"bitbucket.org/decimalteam/go-node/x/gov/internal/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GetTallyParams returns the current TallyParam from the global param store
func (keeper Keeper) GetTallyParams(ctx sdk.Context) types.TallyParams {
	var tallyParams types.TallyParams
	keeper.paramSpace.Get(ctx, types.ParamStoreKeyTallyParams, &tallyParams)
	if ctx.BlockHeight() >= updates.Update13Block {
		tallyParams.Quorum = sdk.NewDec(2).QuoInt64(3)
	}
	if ctx.BlockHeight() >= updates.Update13Block {
		tallyParams.Quorum = sdk.NewDecWithPrec(667, 3)
	}
	return tallyParams
}

// SetTallyParams sets TallyParams to the global param store
func (keeper Keeper) SetTallyParams(ctx sdk.Context, tallyParams types.TallyParams) {
	keeper.paramSpace.Set(ctx, types.ParamStoreKeyTallyParams, &tallyParams)
}
