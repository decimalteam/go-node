package keeper

import (
	"bitbucket.org/decimalteam/go-node/x/coin"
	"fmt"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"

	"bitbucket.org/decimalteam/go-node/x/multisig/internal/types"
)

// Keeper of the multisig store
type Keeper struct {
	storeKey      sdk.StoreKey
	cdc           *codec.LegacyAmino
	paramspace    types.ParamSubspace
	AccountKeeper auth.AccountKeeper
	CoinKeeper    coin.Keeper
	BankKeeper    bank.Keeper
}

// NewKeeper creates a multisig keeper
func NewKeeper(cdc *codec.LegacyAmino, key sdk.StoreKey, paramspace types.ParamSubspace, accountKeeper auth.AccountKeeper, coinKeeper coin.Keeper, bankKeeper bank.Keeper) Keeper {
	keeper := Keeper{
		storeKey:      key,
		cdc:           cdc,
		paramspace:    paramspace.WithKeyTable(types.ParamKeyTable()),
		AccountKeeper: accountKeeper,
		CoinKeeper:    coinKeeper,
		BankKeeper:    bankKeeper,
	}
	return keeper
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// GetIterator returns iterator over KVStore with specified prefix.
func (k Keeper) GetIterator(ctx sdk.Context, prefix string) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, []byte(prefix))
}

// GetWallet returns multisig wallet metadata struct with specified address.
func (k Keeper) GetWallet(ctx sdk.Context, address string) types.Wallet {
	key := fmt.Sprintf("wallet/%s", address)
	store := ctx.KVStore(k.storeKey)
	if !store.Has([]byte(key)) {
		return types.Wallet{}
	}
	bz := store.Get([]byte(key))
	var wallet types.Wallet
	k.cdc.MustUnmarshalBinaryBare(bz, &wallet)
	return wallet
}

// SetWallet sets the entire wallet metadata struct for a multisig wallet.
func (k Keeper) SetWallet(ctx sdk.Context, wallet types.Wallet) {
	key := fmt.Sprintf("wallet/%s", wallet.Address.String())
	store := ctx.KVStore(k.storeKey)
	store.Set([]byte(key), k.cdc.MustMarshalBinaryBare(wallet))
}

// GetTransaction returns multisig wallet transaction metadata with specified address transaction ID.
func (k Keeper) GetTransaction(ctx sdk.Context, txID string) types.Transaction {
	key := fmt.Sprintf("tx/%s", txID)
	store := ctx.KVStore(k.storeKey)
	if !store.Has([]byte(key)) {
		return types.Transaction{}
	}
	bz := store.Get([]byte(key))
	var transaction types.Transaction
	k.cdc.MustUnmarshalBinaryBare(bz, &transaction)
	return transaction
}

// SetTransaction sets the entire multisig wallet transaction metadata struct for a multisig wallet.
func (k Keeper) SetTransaction(ctx sdk.Context, transaction types.Transaction) {
	key := fmt.Sprintf("tx/%s", transaction.ID)
	store := ctx.KVStore(k.storeKey)
	store.Set([]byte(key), k.cdc.MustMarshalBinaryBare(transaction))
}
