package coin

import (
	"bitbucket.org/decimalteam/go-node/x/coin/internal/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// InitGenesis initialize default parameters
// and the keeper's address to pubkey map
func InitGenesis(ctx sdk.Context, k Keeper, data types.GenesisState) []abci.ValidatorUpdate {
	var coin = types.Coin{
		Title:  k.Config.TitleBaseCoin,
		Symbol: k.Config.SymbolBaseCoin,
		Volume: k.Config.InitialVolumeBaseCoin,
	}
	k.SetCoin(ctx, coin)

	//test
	var testCoin = types.Coin{
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
func ExportGenesis(ctx sdk.Context, k Keeper) (data types.GenesisState) {
	// TODO: Define logic for exporting state
	return types.NewGenesisState(k.Config.TitleBaseCoin, k.Config.SymbolBaseCoin, k.Config.InitialVolumeBaseCoin)
}
