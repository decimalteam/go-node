package keeper

import (
	"bitbucket.org/decimalteam/go-node/utils/formulas"
	cliUtils "bitbucket.org/decimalteam/go-node/x/coin/client/utils"
	"encoding/hex"
	"fmt"
	"strings"

	"bitbucket.org/decimalteam/go-node/config"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/tendermint/tendermint/libs/log"

	"bitbucket.org/decimalteam/go-node/x/coin/internal/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Keeper of the coin store
type Keeper struct {
	storeKey      sdk.StoreKey
	cdc           *codec.Codec
	paramspace    types.ParamSubspace
	AccountKeeper auth.AccountKeeper
	BankKeeper    bank.Keeper
	Config        *config.Config
}

// NewKeeper creates a coin keeper
func NewKeeper(cdc *codec.Codec, key sdk.StoreKey, paramspace types.ParamSubspace, accountKeeper auth.AccountKeeper, coinKeeper bank.Keeper, config *config.Config) Keeper {
	keeper := Keeper{
		storeKey:      key,
		cdc:           cdc,
		paramspace:    paramspace.WithKeyTable(types.ParamKeyTable()),
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
func (k Keeper) UpdateBalance(ctx sdk.Context, coinSymbol string, amount sdk.Int, address sdk.AccAddress) error {
	// Get account instance
	acc := k.AccountKeeper.GetAccount(ctx, address)
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
	k.AccountKeeper.SetAccount(ctx, acc)
	return nil
}

func (k Keeper) IsCoinBase(symbol string) bool {
	return k.Config.SymbolBaseCoin == symbol
}

func (k Keeper) UpdateCoin(ctx sdk.Context, coin types.Coin, reserve sdk.Int, volume sdk.Int) {
	coin.Reserve = reserve
	coin.Volume = volume
	k.SetCoin(ctx, coin)
}

func (k Keeper) GetCommission(ctx sdk.Context, commissionInBaseCoin sdk.Int) (sdk.Int, string, error) {
	feeCoin, ok := ctx.Value(types.FeeCoin).(string)
	if !ok {
		feeCoin = cliUtils.GetBaseCoin()
	}

	commission := commissionInBaseCoin

	if feeCoin != cliUtils.GetBaseCoin() {
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
