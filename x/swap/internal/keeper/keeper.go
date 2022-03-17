package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/supply"

	"bitbucket.org/decimalteam/go-node/utils/updates"
	"bitbucket.org/decimalteam/go-node/x/coin"
	"bitbucket.org/decimalteam/go-node/x/swap/internal/types"
)

// Keeper of the validator store
type Keeper struct {
	storeKey      sdk.StoreKey
	cdc           *codec.Codec
	paramSpace    types.ParamSubspace
	coinKeeper    coin.Keeper
	accountKeeper auth.AccountKeeper
	supplyKeeper  supply.Keeper
}

func NewKeeper(cdc *codec.Codec, storeKey sdk.StoreKey, paramSpace types.ParamSubspace, coinKeeper coin.Keeper,
	accountKeeper auth.AccountKeeper, supplyKeeper supply.Keeper) Keeper {
	if addr := supplyKeeper.GetModuleAddress(types.PoolName); addr == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.PoolName))
	}
	return Keeper{
		storeKey:      storeKey,
		cdc:           cdc,
		paramSpace:    paramSpace.WithKeyTable(ParamKeyTable()),
		coinKeeper:    coinKeeper,
		accountKeeper: accountKeeper,
		supplyKeeper:  supplyKeeper,
	}
}

func (k Keeper) CheckBalance(ctx sdk.Context, address sdk.AccAddress, coins sdk.Coins) (bool, error) {
	account := k.accountKeeper.GetAccount(ctx, address)
	if account == nil {
		return false, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "account not found")
	}
	if !account.GetCoins().IsAllGTE(coins) {
		return false, nil
	}

	return true, nil
}

func (k Keeper) UnlockFunds(ctx sdk.Context, address sdk.AccAddress, coins sdk.Coins) error {
	return k.supplyKeeper.SendCoinsFromModuleToAccount(ctx, types.PoolName, address, coins)
}

func (k Keeper) LockFunds(ctx sdk.Context, address sdk.AccAddress, coins sdk.Coins) error {
	return k.supplyKeeper.SendCoinsFromAccountToModule(ctx, address, types.PoolName, coins)
}

func (k Keeper) CheckPoolFunds(ctx sdk.Context, coins sdk.Coins) (bool, error) {
	accountAddr := k.supplyKeeper.GetModuleAddress(types.PoolName)
	if accountAddr.Empty() {
		return false, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "module account not found")
	}

	account := k.accountKeeper.GetAccount(ctx, accountAddr)
	if account == nil {
		return false, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "module account not found")
	}

	if !account.GetCoins().IsAllGTE(coins) {
		return false, nil
	}

	return true, nil
}

func (k Keeper) GetLockedFunds(ctx sdk.Context) sdk.Coins {
	account := k.supplyKeeper.GetModuleAccount(ctx, types.PoolName)
	if account == nil {
		panic("account not found")
	}

	return account.GetCoins()
}

func (k Keeper) MigrateToUpdatedPrefixes(ctx sdk.Context) error {
	if ctx.BlockHeight() != updates.Update14Block {
		panic(fmt.Sprintf("wrong time for data migration (called at block %d instead of %d)", ctx.BlockHeight(), updates.Update14Block))
	}
	k.migrateSwaps(ctx)
	k.migrateSwapsV2(ctx)
	k.migrateChains(ctx)
	return nil
}

func (k Keeper) migrateSwaps(ctx sdk.Context) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, []byte{0x50, 0x01})
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		keyFrom, value := iterator.Key(), iterator.Value()
		if len(keyFrom) != 34 {
			continue
		}
		var swap types.Swap
		k.cdc.MustUnmarshalBinaryLengthPrefixed(value, &swap)
		keyTo := types.GetSwapKey(ctx, swap.HashedSecret)
		store.Set(keyTo, value)
		store.Delete(keyFrom)
	}
}

func (k Keeper) migrateSwapsV2(ctx sdk.Context) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, []byte{0x50, 0x02})
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		keyFrom, value := iterator.Key(), iterator.Value()
		if value == nil {
			value = []byte{}
		}
		if len(keyFrom) != 34 {
			continue
		}
		keyTo := append(types.SwapV2Key, keyFrom[2:]...)
		store.Set(keyTo, value)
		store.Delete(keyFrom)
	}
}

func (k Keeper) migrateChains(ctx sdk.Context) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, []byte{0x50, 0x03})
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		keyFrom, value := iterator.Key(), iterator.Value()
		if len(keyFrom) != 10 {
			continue
		}
		var chain types.Chain
		k.cdc.MustUnmarshalBinaryLengthPrefixed(value, &chain)
		keyTo := append(types.ChainKey, keyFrom[2:]...)
		store.Set(keyTo, value)
		store.Delete(keyFrom)
	}
}
