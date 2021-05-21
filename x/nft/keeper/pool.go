package keeper

import (
	types2 "bitbucket.org/decimalteam/go-node/x/nft/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/supply/exported"
)

func (k Keeper) GetReservedPool(ctx sdk.Context) exported.ModuleAccountI {
	return k.supplyKeeper.GetModuleAccount(ctx, types2.ReservedPool)
}

func (k Keeper) ReserveTokens(ctx sdk.Context, amount sdk.Coins, address sdk.AccAddress) error {
	return k.supplyKeeper.SendCoinsFromAccountToModule(ctx, address, types2.ReservedPool, amount)
}

func (k Keeper) BurnTokens(ctx sdk.Context, amount sdk.Coins) error {
	return k.supplyKeeper.BurnCoins(ctx, types2.ReservedPool, amount)
}
