package coin

import (
	"strings"

	"bitbucket.org/decimalteam/go-node/config"
	"bitbucket.org/decimalteam/go-node/utils/updates"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// BeginBlocker check for infraction evidence or downtime of validators
// on every begin block
func BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock, k Keeper) {
	k.ClearCoinCache()

	if ctx.BlockHeight() == updates.Update3Block {
		coins := k.GetAllCoins(ctx)
		for _, coin := range coins {
			if strings.ToLower(coin.Symbol) == config.SymbolBaseCoin {
				continue
			}
			k.SetCachedCoin(coin.Symbol)
		}
	}
}

// EndBlocker called every block, process inflation, update validator set.
func EndBlocker(ctx sdk.Context, k Keeper) {
}
