package keeper

import (
	"bitbucket.org/decimalteam/go-node/x/swap/internal/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	"time"
)

const (
	DefaultParamspace = types.ModuleName
)

func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable().RegisterParamSet(&types.Params{})
}

func (k Keeper) LockedTimeOut(ctx sdk.Context) (res time.Duration) {
	k.paramSpace.Get(ctx, types.KeyLockedTimeOut, &res)
	return
}

func (k Keeper) LockedTimeIn(ctx sdk.Context) (res time.Duration) {
	k.paramSpace.Get(ctx, types.KeyLockedTimeIn, &res)
	return
}

func (k Keeper) GetParams(ctx sdk.Context) types.Params {
	return types.NewParams(k.LockedTimeOut(ctx), k.LockedTimeIn(ctx))
}

// set the params
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramSpace.SetParamSet(ctx, &params)
}
