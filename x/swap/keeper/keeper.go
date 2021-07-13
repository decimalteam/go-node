package keeper

import (
	"bitbucket.org/decimalteam/go-node/x/coin"
	types2 "bitbucket.org/decimalteam/go-node/x/swap/types"
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	authKeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	bankKeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
)

// Keeper of the validator store
type Keeper struct {
	storeKey      sdk.StoreKey
	cdc           codec.BinaryCodec
	paramSpace    types2.ParamSubspace
	coinKeeper    coin.Keeper
	accountKeeper authKeeper.AccountKeeper
	baseKeeper    bankKeeper.BaseKeeper
}

func NewKeeper(cdc codec.BinaryCodec, storeKey sdk.StoreKey, paramSpace types2.ParamSubspace, coinKeeper coin.Keeper,
	accountKeeper authKeeper.AccountKeeper, baseKeeper bankKeeper.BaseKeeper) Keeper {
	if addr := accountKeeper.GetModuleAddress(types2.PoolName); addr == nil {
		panic(fmt.Sprintf("%s module account has not been set", types2.PoolName))
	}
	return Keeper{
		storeKey:      storeKey,
		cdc:           cdc,
		paramSpace:    paramSpace.WithKeyTable(ParamKeyTable()),
		coinKeeper:    coinKeeper,
		accountKeeper: accountKeeper,
		baseKeeper:    baseKeeper,
	}
}

func (k Keeper) CheckBalance(ctx sdk.Context, address sdk.AccAddress, coins sdk.Coins) (bool, error) {
	account := k.accountKeeper.GetAccount(ctx, address)
	if account == nil {
		return false, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "account not found")
	}
	if !k.baseKeeper.GetAllBalances(ctx, account.GetAddress()).IsAllGTE(coins) {
		return false, nil
	}

	return true, nil
}

func (k Keeper) UnlockFunds(ctx sdk.Context, address sdk.AccAddress, coins sdk.Coins) error {
	return k.baseKeeper.SendCoinsFromModuleToAccount(ctx, types2.PoolName, address, coins)
}

func (k Keeper) LockFunds(ctx sdk.Context, address sdk.AccAddress, coins sdk.Coins) error {
	return k.baseKeeper.SendCoinsFromAccountToModule(ctx, address, types2.PoolName, coins)
}

func (k Keeper) CheckPoolFunds(ctx sdk.Context, coins sdk.Coins) (bool, error) {
	accountAddr := k.accountKeeper.GetModuleAddress(types2.PoolName)
	if accountAddr.Empty() {
		return false, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "account not found")
	}

	account := k.accountKeeper.GetAccount(ctx, accountAddr)
	if account == nil {
		return false, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "account not found")
	}

	if !k.baseKeeper.GetAllBalances(ctx, account.GetAddress()).IsAllGTE(coins) {
		return false, nil
	}

	return true, nil
}
