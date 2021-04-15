package coin

import (
	"bitbucket.org/decimalteam/go-node/utils/updates"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// BeginBlocker check for infraction evidence or downtime of validators
// on every begin block
func BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock, k Keeper) {
	k.ClearCoinCache()

	if ctx.BlockHeight() == updates.Update6Block {
		coins := map[string]string{
			"btt":        "dx1gxxcvyr27xa9g03f0pp3k9cmyrwdxq3u9sd0mg",
			"chugacoin":  "dx1watwunjhgd8jzzw8q77fgtu9654p2daer3gk34",
			"crypton":    "dx1naeup2d7gc30tw0wqxgt4c5yau0atmxvtl4vrg",
			"dar":        "dx18kuajadwqwgp7xyltmkhez7n3t54xsmx4jz5nu",
			"diamond":    "dx1wcmq8ply2la3duzegav2hw4r363hlj48we7ysk",
			"legiondocs": "dx1m8xj0f8snwg45xqejere52l406yx85myr52htj",
			"legiongame": "dx1nafxm7gn4kmyjtctya7cshj4nj956k5tq5p9wu",
			"purplish":   "dx159ztx4yg3ex3yghzd7q8pudzsa726e0hu5c09g",
			"rrunion":    "dx18tuyg2dqq8uvljj3hytfph02ahdkgt4yju6dne",
			"sertifikat": "dx1nafxm7gn4kmyjtctya7cshj4nj956k5tq5p9wu",
		}
		for symbol, address := range coins {
			coin, err := k.GetCoin(ctx, symbol)
			if err != nil {
				panic(err)
			}
			coin.Creator, err = sdk.AccAddressFromBech32(address)
			if err != nil {
				panic(err)
			}
			k.SetCoin(ctx, coin)
		}
	}

	if ctx.BlockHeight() == updates.Update7Block {
		coins := k.GetAllCoins(ctx)
		for _, coin := range coins {
			k.SetCachedCoin(coin.Symbol)
		}
	}
}

// EndBlocker called every block, process inflation, update validator set.
func EndBlocker(ctx sdk.Context, k Keeper) {
}
