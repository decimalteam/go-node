package keeper

import (
	types2 "bitbucket.org/decimalteam/go-node/x/gov/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GetTallyParams returns the current TallyParam from the global param store
func (keeper Keeper) GetTallyParams(ctx sdk.Context) types2.TallyParams {
	var tallyParams types2.TallyParams
	keeper.paramSpace.Get(ctx, types2.ParamStoreKeyTallyParams, &tallyParams)
	return tallyParams
}

// SetTallyParams sets TallyParams to the global param store
func (keeper Keeper) SetTallyParams(ctx sdk.Context, tallyParams types2.TallyParams) {
	keeper.paramSpace.Set(ctx, types2.ParamStoreKeyTallyParams, &tallyParams)
}
