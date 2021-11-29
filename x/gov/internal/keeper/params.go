package keeper

import (
	"bitbucket.org/decimalteam/go-node/x/gov/internal/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GetTallyParams returns the current TallyParam from the global param store
func (keeper Keeper) GetTallyParams(ctx sdk.Context) types.TallyParams {
	var tallyParams types.TallyParams
	keeper.paramSpace.Get(ctx, types.ParamStoreKeyTallyParams, &tallyParams)

	tallyParams.Quorum = sdk.NewDecWithPrec(667, 3)

	return tallyParams
}

// SetTallyParams sets TallyParams to the global param store
func (keeper Keeper) SetTallyParams(ctx sdk.Context, tallyParams types.TallyParams) {
	keeper.paramSpace.Set(ctx, types.ParamStoreKeyTallyParams, &tallyParams)
}
