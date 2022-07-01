package coin

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// InitGenesis initialize default parameters
// and the keeper's address to pubkey map
func InitGenesis(ctx sdk.Context, k Keeper, data GenesisState) []abci.ValidatorUpdate {
	coin := Coin{
		Title:  data.Title,
		Symbol: data.Symbol,
		Volume: data.InitialVolume,
	}

	k.SetCoin(ctx, coin)
	return []abci.ValidatorUpdate{}
}

// ExportGenesis writes the current store values to a genesis file, which can be imported again with InitGenesis
func ExportGenesis(ctx sdk.Context, k Keeper) (data GenesisState) {
	coins := k.GetAllCoins(ctx)

	return NewGenesisState(k.Config.TitleBaseCoin, k.Config.SymbolBaseCoin, k.Config.InitialVolumeBaseCoin, coins)
}
