package keeper

import (
	"bitbucket.org/decimalteam/go-node/config"
	"fmt"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/tendermint/tendermint/libs/log"
	"strings"

	"bitbucket.org/decimalteam/go-node/x/coin/internal/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Keeper of the coin store
type Keeper struct {
	storeKey      sdk.StoreKey
	cdc           *codec.Codec
	paramspace    types.ParamSubspace
	codespace     sdk.CodespaceType
	AccountKeeper auth.AccountKeeper
	BankKeeper    bank.Keeper
	Config        *config.Config
}

// NewKeeper creates a coin keeper
func NewKeeper(cdc *codec.Codec, key sdk.StoreKey, paramspace types.ParamSubspace, codespace sdk.CodespaceType, accountKeeper auth.AccountKeeper, coinKeeper bank.Keeper, config *config.Config) Keeper {
	keeper := Keeper{
		storeKey:      key,
		cdc:           cdc,
		paramspace:    paramspace.WithKeyTable(types.ParamKeyTable()),
		codespace:     codespace,
		AccountKeeper: accountKeeper,
		BankKeeper:    coinKeeper,
		Config:        config,
	}
	return keeper
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// GetCoin returns types.Coin instance if exists in KVStore
func (k Keeper) GetCoin(ctx sdk.Context, symbol string) (types.Coin, error) {
	store := ctx.KVStore(k.storeKey)
	var coin types.Coin
	byteKey := []byte(types.CoinPrefix + symbol)
	err := k.cdc.UnmarshalBinaryLengthPrefixed(store.Get(byteKey), &coin)
	if err != nil {
		return coin, err
	}
	return coin, nil
}

func (k Keeper) SetCoin(ctx sdk.Context, value types.Coin) { // CreateCoin
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(value)
	key := value.Symbol
	store.Set([]byte(types.CoinPrefix+key), bz)
}

// GetCoinsIterator gets an iterator over all Coins in which the keys are the symbols and the values are the coins
func (k Keeper) GetCoinsIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, []byte(types.CoinPrefix))
}

// Returns integer abs
func Abs(x sdk.Int) sdk.Int {
	if x.IsNegative() {
		return x.Neg()
	} else {
		return x
	}
}

// Updating balances
func (k Keeper) UpdateBalance(ctx sdk.Context, coinSymbol string, amount sdk.Int, address sdk.AccAddress) {
	// Get account instance
	acc := k.AccountKeeper.GetAccount(ctx, address)
	// Get account coins information
	coins := acc.GetCoins()
	isNeg := amount.IsNegative()
	updAmount := Abs(amount)
	updCoin := sdk.Coins{sdk.NewCoin(strings.ToLower(coinSymbol), updAmount)}
	// Updating coin information
	if isNeg {
		coins = coins.Sub(updCoin)
	} else {
		coins = coins.Add(updCoin)
	}
	// Update coin information
	_ = acc.SetCoins(coins)
	// Update account information
	k.AccountKeeper.SetAccount(ctx, acc)
}

func (k Keeper) IsCoinBase(symbol string) bool {
	return k.Config.SymbolBaseCoin == symbol
}
