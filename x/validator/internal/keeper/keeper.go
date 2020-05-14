package keeper

import (
	"container/list"
	"fmt"
	"strings"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/supply"

	"bitbucket.org/decimalteam/go-node/x/coin"
	"bitbucket.org/decimalteam/go-node/x/validator/internal/types"
)

const aminoCacheSize = 500

// Keeper of the validator store
type Keeper struct {
	storeKey         sdk.StoreKey
	cdc              *codec.Codec
	paramSpace       types.ParamSubspace
	coinKeeper       coin.Keeper
	supplyKeeper     supply.Keeper
	hooks            types.ValidatorHooks
	FeeCollectorName string

	validatorCache     map[string]cachedValidator
	validatorCacheList *list.List
}

// NewKeeper creates a validator keeper
func NewKeeper(cdc *codec.Codec, key sdk.StoreKey, paramSpace types.ParamSubspace, coinKeeper coin.Keeper, supplyKeeper supply.Keeper, feeCollectorName string) Keeper {

	// ensure bonded and not bonded module accounts are set
	if addr := supplyKeeper.GetModuleAddress(types.BondedPoolName); addr == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.BondedPoolName))
	}

	if addr := supplyKeeper.GetModuleAddress(types.NotBondedPoolName); addr == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.NotBondedPoolName))
	}

	keeper := Keeper{
		storeKey:           key,
		cdc:                cdc,
		paramSpace:         paramSpace.WithKeyTable(ParamKeyTable()),
		coinKeeper:         coinKeeper,
		supplyKeeper:       supplyKeeper,
		validatorCache:     make(map[string]cachedValidator, aminoCacheSize),
		validatorCacheList: list.New(),
		FeeCollectorName:   feeCollectorName,
	}
	return keeper
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) Codespace() string {
	return types.DefaultCodespace
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

func (k Keeper) GetCoin(ctx sdk.Context, symbol string) (coin.Coin, error) {
	if symbol == "tdel" {
		symbol = "tDEL"
	} else {
		symbol = strings.ToUpper(symbol)
	}
	return k.coinKeeper.GetCoin(ctx, symbol)
}
