package keeper

import (
	types2 "bitbucket.org/decimalteam/go-node/x/swap/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/supply/exported"
)

func (k Keeper) GetPool(ctx sdk.Context) exported.ModuleAccountI {
	return k.supplyKeeper.GetModuleAccount(ctx, types2.PoolName)
}
