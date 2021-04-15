package coin

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
	"log"
)

// BeginBlocker check for infraction evidence or downtime of validators
// on every begin block
func BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock, k Keeper) {
	k.ClearCoinCache()

	if ctx.BlockHeight() == 1_205_000 {
		coin, err := k.GetCoin(ctx, "cantsell")
		if err != nil {
			log.Println(err)
			return
		}
		coin.Creator, err = sdk.AccAddressFromBech32("dx1lx4lvt8sjuxj8vw5dcf6knnq0pacre4w6hdh2v")
		if err != nil {
			log.Println(err)
			return
		}
		k.SetCoin(ctx, coin)

		coin, err = k.GetCoin(ctx, "timocoin")
		if err != nil {
			log.Println(err)
			return
		}
		coin.Creator, err = sdk.AccAddressFromBech32("dx1mxvw3d39fmn4vzq8x6xycjm79vkdefnwulrkzz")
		if err != nil {
			log.Println(err)
			return
		}
		k.SetCoin(ctx, coin)
	}

	if ctx.BlockHeight() == 2_336_260 {
		coins := k.GetAllCoins(ctx)
		for _, coin := range coins {
			k.SetCachedCoin(coin.Symbol)
		}
	}
}

// EndBlocker called every block, process inflation, update validator set.
func EndBlocker(ctx sdk.Context, k Keeper) {
}
