package keeper

import (
	"bitbucket.org/decimalteam/go-node/x/coin"
	"bitbucket.org/decimalteam/go-node/x/swap/internal/types"
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/supply"
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
		return false, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "account not found")
	}

	account := k.accountKeeper.GetAccount(ctx, accountAddr)
	if account == nil {
		return false, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "account not found")
	}

	if !account.GetCoins().IsAllGTE(coins) {
		return false, nil
	}

	return true, nil
}
