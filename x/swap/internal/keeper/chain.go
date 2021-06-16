package keeper

import (
	"bitbucket.org/decimalteam/go-node/x/swap/internal/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) HasChain(ctx sdk.Context, chainNumber int) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.GetChainKey(chainNumber))
}

func (k Keeper) SetChain(ctx sdk.Context, chainNumber int, chain types.Chain) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(types.NewChain(chain.Name, chain.Active))
	store.Set(types.GetChainKey(chainNumber), bz)
}

func (k Keeper) GetChain(ctx sdk.Context, chainNumber int) (types.Chain, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetChainKey(chainNumber))
	if bz == nil {
		return types.Chain{}, false
	}

	var chain types.Chain
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &chain)
	return chain, true
}
