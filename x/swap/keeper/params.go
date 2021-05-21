package keeper

import (
	types2 "bitbucket.org/decimalteam/go-node/x/swap/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	"time"
)

const (
	DefaultParamspace = types2.ModuleName
)

func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable().RegisterParamSet(&types2.Params{})
}

func (k Keeper) LockedTimeOut(ctx sdk.Context) (res time.Duration) {
	k.paramSpace.Get(ctx, types2.KeyLockedTimeOut, &res)
	return time.Minute * 4
}

func (k Keeper) LockedTimeIn(ctx sdk.Context) (res time.Duration) {
	k.paramSpace.Get(ctx, types2.KeyLockedTimeIn, &res)
	return time.Minute * 2
}

func (k Keeper) GetParams(ctx sdk.Context) types2.Params {
	return types2.NewParams(k.LockedTimeOut(ctx), k.LockedTimeIn(ctx))
}

// set the params
func (k Keeper) SetParams(ctx sdk.Context, params types2.Params) {
	k.paramSpace.SetParamSet(ctx, &params)
}
