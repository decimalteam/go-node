package keeper

import (
	"encoding/hex"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/supply"

	"bitbucket.org/decimalteam/go-node/config"
	"bitbucket.org/decimalteam/go-node/utils/formulas"
	"bitbucket.org/decimalteam/go-node/x/coin/internal/types"
)

// Keeper of the coin store
type Keeper struct {
	storeKey      sdk.StoreKey
	cdc           *codec.Codec
	paramspace    types.ParamSubspace
	AccountKeeper auth.AccountKeeper
	BankKeeper    bank.Keeper
	SupplyKeeper  supply.Keeper
	Config        *config.Config

	coinCache      map[string]bool
	coinCacheMutex *sync.Mutex
}

// NewKeeper creates a coin keeper
func NewKeeper(cdc *codec.Codec, key sdk.StoreKey, paramspace types.ParamSubspace, accountKeeper auth.AccountKeeper, coinKeeper bank.Keeper, supplyKeeper supply.Keeper, config *config.Config) Keeper {
	keeper := Keeper{
		storeKey:       key,
		cdc:            cdc,
		paramspace:     paramspace.WithKeyTable(types.ParamKeyTable()),
		AccountKeeper:  accountKeeper,
		BankKeeper:     coinKeeper,
		SupplyKeeper:   supplyKeeper,
		Config:         config,
		coinCache:      make(map[string]bool),
		coinCacheMutex: &sync.Mutex{},
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
	key := []byte(types.CoinPrefix + strings.ToLower(symbol))
	value := store.Get(key)
	if value == nil {
		return coin, fmt.Errorf("coin %s is not found in the key-value store", strings.ToLower(symbol))
	}
	err := k.cdc.UnmarshalBinaryLengthPrefixed(value, &coin)
	if err != nil {
		return coin, err
	}
	return coin, nil
}

func (k Keeper) SetCoin(ctx sdk.Context, coin types.Coin) {
	store := ctx.KVStore(k.storeKey)
	value := k.cdc.MustMarshalBinaryLengthPrefixed(coin)
	key := []byte(types.CoinPrefix + strings.ToLower(coin.Symbol))
	store.Set(key, value)
}

// GetCoinsIterator gets an iterator over all Coins in which the keys are the symbols and the values are the coins
func (k Keeper) GetCoinsIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, []byte(types.CoinPrefix))
}

func (k Keeper) GetAllCoins(ctx sdk.Context) []types.Coin {
	var coins []types.Coin
	iterator := k.GetCoinsIterator(ctx)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var coin types.Coin
		err := k.cdc.UnmarshalBinaryLengthPrefixed(iterator.Value(), &coin)
		if err != nil {
			panic(err)
		}
		coins = append(coins, coin)
	}
	return coins
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
func (k Keeper) UpdateBalance(ctx sdk.Context, coinSymbol string, amount sdk.Int, address sdk.AccAddress) error {
	// Get account instance
	acc := k.AccountKeeper.GetAccount(ctx, address)
	if acc == nil {
		acc = k.AccountKeeper.NewAccountWithAddress(ctx, address)
	}
	// Get account coins information
	coins := acc.GetCoins()
	updAmount := Abs(amount)
	updCoin := sdk.Coins{sdk.NewCoin(strings.ToLower(coinSymbol), updAmount)}
	// Updating coin information
	if amount.IsNegative() {
		coins = coins.Sub(updCoin)
	} else {
		coins = coins.Add(updCoin...)
	}
	// Update coin information
	err := acc.SetCoins(coins)
	if err != nil {
		return err
	}
	// Update account information
	fmt.Println("abc 333")
	k.AccountKeeper.SetAccount(ctx, acc)
	return nil
}

func (k Keeper) IsCoinBase(symbol string) bool {
	return k.Config.SymbolBaseCoin == symbol
}

func (k Keeper) UpdateCoin(ctx sdk.Context, coin types.Coin, reserve sdk.Int, volume sdk.Int) {
	if !coin.IsBase() {
		k.SetCachedCoin(coin.Symbol)
	}
	coin.Reserve = reserve
	coin.Volume = volume
	k.SetCoin(ctx, coin)
	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeUpdateCoin,
		sdk.NewAttribute(types.AttributeSymbol, coin.Symbol),
		sdk.NewAttribute(types.AttributeVolume, coin.Volume.String()),
		sdk.NewAttribute(types.AttributeReserve, coin.Reserve.String()),
	))
}

func (k Keeper) IsCheckRedeemed(ctx sdk.Context, check *types.Check) bool {
	checkHash := check.HashFull()
	store := ctx.KVStore(k.storeKey)
	key := []byte(types.CheckPrefix + hex.EncodeToString(checkHash[:]))
	return len(store.Get(key)) > 0
}

func (k Keeper) SetCheckRedeemed(ctx sdk.Context, check *types.Check) {
	checkHash := check.HashFull()
	store := ctx.KVStore(k.storeKey)
	key := []byte(types.CheckPrefix + hex.EncodeToString(checkHash[:]))
	store.Set(key, []byte{1})
	return
}

func (k Keeper) GetCommission(ctx sdk.Context, commissionInBaseCoin sdk.Int) (sdk.Int, string, error) {
	var feeCoin string
	fee, ok := ctx.Value("fee").(sdk.Coins)
	if !ok || fee == nil {
		feeCoin = k.GetBaseCoin(ctx)
		return commissionInBaseCoin, feeCoin, nil
	}

	commission := sdk.ZeroInt()

	coin := fee[0]

	feeCoin = coin.Denom
	if feeCoin != k.GetBaseCoin(ctx) {
		coinInfo, err := k.GetCoin(ctx, feeCoin)
		if err != nil {
			return sdk.Int{}, "", err
		}

		if coinInfo.Reserve.LT(commissionInBaseCoin) {
			return sdk.Int{}, "", fmt.Errorf("coin reserve balance is not sufficient for transaction. Has: %s, required %s",
				coinInfo.Reserve.String(),
				commissionInBaseCoin.String())
		}

		commission = formulas.CalculateSaleAmount(coinInfo.Volume, coinInfo.Reserve, coinInfo.CRR, commissionInBaseCoin)
	}

	return commission, feeCoin, nil
}

func (k *Keeper) SetCachedCoin(coin string) {
	defer k.coinCacheMutex.Unlock()
	k.coinCacheMutex.Lock()

	_, _ = os.Stdout.WriteString(fmt.Sprintf("COINCACHE SetCachedCoin %s", coin))
	k.coinCache[coin] = true
}

func (k *Keeper) ClearCoinCache() {
	defer k.coinCacheMutex.Unlock()
	k.coinCacheMutex.Lock()
	_, _ = os.Stdout.WriteString(fmt.Sprintf("COINCACHE ClearCoinCache"))
	for key := range k.coinCache {
		delete(k.coinCache, key)
	}
}

func (k Keeper) GetCoinsCache() map[string]bool {
	defer k.coinCacheMutex.Unlock()
	k.coinCacheMutex.Lock()
	_, _ = os.Stdout.WriteString(fmt.Sprintf("COINCACHE GetCoinsCacheAAA"))
	return k.coinCache
}

func (k Keeper) GetCoinCache(symbol string) bool {
	defer k.coinCacheMutex.Unlock()
	k.coinCacheMutex.Lock()
	_, ok := k.coinCache[symbol]
	_, _ = os.Stdout.WriteString(fmt.Sprintf("COINCACHE GetCoinCache %s %v", symbol, ok))
	return ok
}

func (k Keeper) GetBaseCoin(ctx sdk.Context) string {
	if ctx.BlockHeight() < 33000 {
		return config.SymbolTestBaseCoin
	}

	return k.Config.SymbolBaseCoin
}
