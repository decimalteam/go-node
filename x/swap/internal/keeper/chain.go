package keeper

import (
	"bitbucket.org/decimalteam/go-node/x/swap/internal/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) HasDestChain(ctx sdk.Context, destChain int) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.GetDestChainKey(destChain))
}

func (k Keeper) SetDestChain(ctx sdk.Context, destChain int, name string) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetDestChainKey(destChain), []byte(name))
}

func (k Keeper) GetDestChainName(ctx sdk.Context, destChain int) string {
	store := ctx.KVStore(k.storeKey)
	return string(store.Get(types.GetDestChainKey(destChain)))
}

func (k Keeper) DeleteDestChain(ctx sdk.Context, destChain int) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetDestChainKey(destChain))
}
