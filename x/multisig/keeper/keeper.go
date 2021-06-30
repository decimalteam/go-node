package keeper

import (
	"bitbucket.org/decimalteam/go-node/x/coin"
	types2 "bitbucket.org/decimalteam/go-node/x/multisig/types"
	"fmt"
	"github.com/cosmos/cosmos-sdk/x/auth/keeper"
	bankKeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Keeper of the multisig store
type Keeper struct {
	storeKey      sdk.StoreKey
	cdc           *codec.LegacyAmino
	paramspace    types2.ParamSubspace
	AccountKeeper keeper.AccountKeeper
	CoinKeeper    coin.Keeper
	BankKeeper    bankKeeper.Keeper
}

// NewKeeper creates a multisig keeper
func NewKeeper(cdc *codec.LegacyAmino, key sdk.StoreKey, paramspace types2.ParamSubspace, accountKeeper keeper.AccountKeeper, coinKeeper coin.Keeper, bankKeeper bankKeeper.Keeper) Keeper {
	keeper := Keeper{
		storeKey:      key,
		cdc:           cdc,
		paramspace:    paramspace.WithKeyTable(types2.ParamKeyTable()),
		AccountKeeper: accountKeeper,
		CoinKeeper:    coinKeeper,
		BankKeeper:    bankKeeper,
	}
	return keeper
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types2.ModuleName))
}

// GetIterator returns iterator over KVStore with specified prefix.
func (k Keeper) GetIterator(ctx sdk.Context, prefix string) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, []byte(prefix))
}

// GetWallet returns multisig wallet metadata struct with specified address.
func (k Keeper) GetWallet(ctx sdk.Context, address string) types2.Wallet {
	key := fmt.Sprintf("wallet/%s", address)
	store := ctx.KVStore(k.storeKey)
	if !store.Has([]byte(key)) {
		return types2.Wallet{}
	}
	bz := store.Get([]byte(key))
	var wallet types2.Wallet
	k.cdc.MustUnmarshal(bz, &wallet)
	return wallet
}

// SetWallet sets the entire wallet metadata struct for a multisig wallet.
func (k Keeper) SetWallet(ctx sdk.Context, wallet types2.Wallet) {
	key := fmt.Sprintf("wallet/%s", wallet.Address.String())
	store := ctx.KVStore(k.storeKey)
	store.Set([]byte(key), k.cdc.MustMarshal(wallet))
}

// GetTransaction returns multisig wallet transaction metadata with specified address transaction ID.
func (k Keeper) GetTransaction(ctx sdk.Context, txID string) types2.Transaction {
	key := fmt.Sprintf("tx/%s", txID)
	store := ctx.KVStore(k.storeKey)
	if !store.Has([]byte(key)) {
		return types2.Transaction{}
	}
	bz := store.Get([]byte(key))
	var transaction types2.Transaction
	k.cdc.MustUnmarshal(bz, &transaction)
	return transaction
}

// SetTransaction sets the entire multisig wallet transaction metadata struct for a multisig wallet.
func (k Keeper) SetTransaction(ctx sdk.Context, transaction types2.Transaction) {
	key := fmt.Sprintf("tx/%s", transaction.ID)
	store := ctx.KVStore(k.storeKey)
	store.Set([]byte(key), k.cdc.MustMarshal(transaction))
}
