package keeper

import (
	"bitbucket.org/decimalteam/go-node/x/coin"
	"container/list"
	"fmt"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/cosmos/cosmos-sdk/x/supply"

	"github.com/tendermint/tendermint/libs/log"

	"bitbucket.org/decimalteam/go-node/x/validator/internal/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const aminoCacheSize = 500

// Keeper of the validator store
type Keeper struct {
	storeKey      sdk.StoreKey
	cdc           *codec.Codec
	paramSpace    types.ParamSubspace
	codespace     sdk.CodespaceType
	stakingKeeper staking.Keeper
	coinKeeper    coin.Keeper
	supplyKeeper  supply.Keeper

	validatorCache     map[string]cachedValidator
	validatorCacheList *list.List
}

// NewKeeper creates a validator keeper
func NewKeeper(cdc *codec.Codec, key sdk.StoreKey, paramSpace types.ParamSubspace, codespace sdk.CodespaceType, coinKeeper coin.Keeper) Keeper {
	keeper := Keeper{
		storeKey:   key,
		cdc:        cdc,
		paramSpace: paramSpace.WithKeyTable(ParamKeyTable()),
		codespace:  codespace,
		coinKeeper: coinKeeper,
	}
	return keeper
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) Codespace() sdk.CodespaceType {
	return k.codespace
}

// Get returns the pubkey from the adddress-pubkey relation
func (k Keeper) Get(ctx sdk.Context, key []byte, value interface{}) error {
	store := ctx.KVStore(k.storeKey)
	err := k.cdc.UnmarshalBinaryLengthPrefixed(store.Get(key), &value)
	if err != nil {
		return err
	}
	return nil
}

func (k Keeper) set(ctx sdk.Context, key []byte, value interface{}) error {
	store := ctx.KVStore(k.storeKey)
	bz, err := k.cdc.MarshalBinaryLengthPrefixed(value)
	if err != nil {
		return err
	}
	store.Set(key, bz)
	return nil
}

func (k Keeper) delete(ctx sdk.Context, key []byte) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(key)
}

// Load the last total validator power.
func (k Keeper) GetLastTotalPower(ctx sdk.Context) sdk.Int {
	power := sdk.Int{}
	err := k.Get(ctx, []byte{types.LastTotalPowerKey}, &power)
	if err != nil {
		return sdk.ZeroInt()
	}
	return power
}

// Set the last total validator power.
func (k Keeper) SetLastTotalPower(ctx sdk.Context, power sdk.Int) error {
	return k.set(ctx, []byte{types.LastTotalPowerKey}, power)
}
