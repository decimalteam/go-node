package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/supply"

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

func (k Keeper) IsMigratedToUpdatedPrefixes(ctx sdk.Context) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.LegacyMigrationKey)
}

func (k Keeper) MigrateToUpdatedPrefixes(ctx sdk.Context) error {
	k.migrateSwaps(ctx)
	k.migrateSwapsV2(ctx)
	k.migrateChains(ctx)
	k.finishMigration(ctx)
	return nil
}

func (k Keeper) migrateSwaps(ctx sdk.Context) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.LegacySwapKey)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		keyFrom, value := iterator.Key(), iterator.Value()
		if len(keyFrom) != 34 { // previous key format: 0x5001<hash_Bytes> (2+32)
			continue
		}
		var swap types.Swap
		k.cdc.MustUnmarshalBinaryLengthPrefixed(value, &swap)
		keyTo := types.GetSwapKey(swap.HashedSecret)
		store.Set(keyTo, value)
		store.Delete(keyFrom)
	}
}

func (k Keeper) migrateSwapsV2(ctx sdk.Context) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.LegacySwapV2Key)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		keyFrom, value := iterator.Key(), iterator.Value()
		if value == nil { // swap v2 just stores the key with empty value
			value = []byte{}
		}
		if len(keyFrom) != 34 { // previous key format: 0x5002<hash_Bytes> (2+32)
			continue
		}
		keyTo := append(types.SwapV2Key, keyFrom[2:]...)
		store.Set(keyTo, value)
		store.Delete(keyFrom)
	}
}

func (k Keeper) migrateChains(ctx sdk.Context) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.LegacyChainKey)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		keyFrom, value := iterator.Key(), iterator.Value()
		if len(keyFrom) != 10 { // previous key format: 0x5003<chainId_Bytes> (2+8)
			continue
		}
		var chain types.Chain
		k.cdc.MustUnmarshalBinaryLengthPrefixed(value, &chain)
		keyTo := append(types.ChainKey, keyFrom[2:]...)
		store.Set(keyTo, value)
		store.Delete(keyFrom)
	}
}

func (k Keeper) finishMigration(ctx sdk.Context) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.LegacyMigrationKey, []byte{})
}
