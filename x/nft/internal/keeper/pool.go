package keeper

import (
	"bitbucket.org/decimalteam/go-node/x/nft/internal/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/supply/exported"
)

func (k Keeper) GetReservedPool(ctx sdk.Context) exported.ModuleAccountI {
	return k.supplyKeeper.GetModuleAccount(ctx, types.ReservedPool)
}

func (k Keeper) ReserveTokens(ctx sdk.Context, amount sdk.Coins, address sdk.AccAddress) error {
	return k.supplyKeeper.SendCoinsFromAccountToModule(ctx, address, types.ReservedPool, amount)
}

func (k Keeper) BurnTokens(ctx sdk.Context, amount sdk.Coins) error {
	return k.supplyKeeper.BurnCoins(ctx, types.ReservedPool, amount)
}
