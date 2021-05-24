package keeper

import (
	"time"

	"bitbucket.org/decimalteam/go-node/x/validator/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramsTypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// Default parameter namespace
const (
	DefaultParamspace = types.ModuleName
)

// ParamTable for staking module
func ParamKeyTable() paramsTypes.KeyTable {
	return paramsTypes.NewKeyTable().RegisterParamSet(&types.Params{})
}

// UnbondingTime
func (k Keeper) UnbondingTime(ctx sdk.Context) (res time.Duration) {
	k.paramSpace.Get(ctx, types.KeyUnbondingTime, &res)
	if ctx.BlockHeight() >= 500 {
		res = time.Minute * 2
	}
	return
}

// MaxValidators - Maximum number of validators
func (k Keeper) MaxValidators(ctx sdk.Context) (res uint16) {
	k.paramSpace.Get(ctx, types.KeyMaxValidators, &res)
	return
}

// MaxEntries - Maximum number of simultaneous unbonding
// delegations or redelegations (per pair/trio)
func (k Keeper) MaxEntries(ctx sdk.Context) (res uint16) {
	k.paramSpace.Get(ctx, types.KeyMaxEntries, &res)
	return
}

// BondDenom - Bondable coin denomination
func (k Keeper) BondDenom(ctx sdk.Context) (res string) {
	k.paramSpace.Get(ctx, types.KeyBondDenom, &res)
	return
}

// HistoricalEntries = number of historical info entries
// to persist in store
func (k Keeper) HistoricalEntries(ctx sdk.Context) (res uint16) {
	k.paramSpace.Get(ctx, types.KeyHistoricalEntries, &res)
	return
}

// MaxDelegations = maximum number of delegations per validator
func (k Keeper) MaxDelegations(ctx sdk.Context) (res uint16) {
	k.paramSpace.Get(ctx, types.KeyMaxDelegations, &res)
	return
}

// Get all parameteras as types.Params
func (k Keeper) GetParams(ctx sdk.Context) types.Params {
	return types.NewParams(
		k.UnbondingTime(ctx),
		k.MaxValidators(ctx),
		k.MaxEntries(ctx),
		k.HistoricalEntries(ctx),
		k.BondDenom(ctx),
		k.MaxDelegations(ctx),
	)
}

// set the params
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramSpace.SetParamSet(ctx, &params)
}
