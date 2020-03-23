package keeper

import (
	"bitbucket.org/decimalteam/go-node/x/coin"
	"fmt"
	"github.com/cosmos/cosmos-sdk/x/auth"

	"github.com/tendermint/tendermint/libs/log"

	"bitbucket.org/decimalteam/go-node/x/check/internal/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Keeper of the check store
type Keeper struct {
	storeKey      sdk.StoreKey
	cdc           *codec.Codec
	paramspace    types.ParamSubspace
	codespace     sdk.CodespaceType
	coinKeeper    coin.Keeper
	accountKeeper auth.AccountKeeper
}

// NewKeeper creates a check keeper
func NewKeeper(cdc *codec.Codec, key sdk.StoreKey, paramspace types.ParamSubspace, codespace sdk.CodespaceType, coinKeeper coin.Keeper, accKeeper auth.AccountKeeper) Keeper {
	keeper := Keeper{
		storeKey:      key,
		cdc:           cdc,
		paramspace:    paramspace.WithKeyTable(types.ParamKeyTable()),
		codespace:     codespace,
		coinKeeper:    coinKeeper,
		accountKeeper: accKeeper,
	}
	return keeper
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// Get returns the pubkey from the address-pubkey relation
func (k Keeper) Get(ctx sdk.Context, key []byte) (types.Check, error) {
	store := ctx.KVStore(k.storeKey)
	var item types.Check
	err := k.cdc.UnmarshalBinaryLengthPrefixed(store.Get(key), &item)
	if err != nil {
		return types.Check{}, err
	}
	return item, nil
}

func (k Keeper) set(ctx sdk.Context, key []byte, value types.Check) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(value)
	store.Set(key, bz)
}
