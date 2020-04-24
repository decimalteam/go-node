package coin

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// InitGenesis initialize default parameters
// and the keeper's address to pubkey map
func InitGenesis(ctx sdk.Context, k Keeper, data GenesisState) []abci.ValidatorUpdate {
	var coin = Coin{
		Title:  k.Config.TitleBaseCoin,
		Symbol: k.Config.SymbolBaseCoin,
		Volume: k.Config.InitialVolumeBaseCoin,
	}
	k.SetCoin(ctx, coin)

	//test
	var testCoin = Coin{
		Title:   k.Config.TitleTestCoin,
		Symbol:  k.Config.SymbolTestCoin,
		Volume:  k.Config.InitialVolumeTestCoin,
		Reserve: k.Config.InitialReserveTestCoin,
		CRR:     k.Config.ConstantReserveRatioTest,
	}

	k.SetCoin(ctx, testCoin)
	return []abci.ValidatorUpdate{}
}

// ExportGenesis writes the current store values to a genesis file, which can be imported again with InitGenesis
func ExportGenesis(ctx sdk.Context, k Keeper) (data GenesisState) {

	return NewGenesisState(k.Config.TitleTestCoin, k.Config.SymbolTestCoin, k.Config.InitialVolumeTestCoin)
}
