package keeper

import (
	types2 "bitbucket.org/decimalteam/go-node/x/swap/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
)

func (k Keeper) GetPool(ctx sdk.Context) types.ModuleAccountI {
	return k.accountKeeper.GetModuleAccount(ctx, types2.PoolName)
}
