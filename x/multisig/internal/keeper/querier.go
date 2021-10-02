package keeper

import (
	"fmt"
	"strings"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"bitbucket.org/decimalteam/go-node/x/multisig/internal/types"
)

// Query endpoints supported by the multisig querier.
const (
	QueryListWallets      = "listWallets"
	QueryGetWallet        = "getWallet"
	QueryListTransactions = "listTransactions"
	QueryGetTransaction   = "getTransaction"
)

// NewQuerier creates a new querier for multisig clients.
func NewQuerier(k Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		switch path[0] {
		case QueryListWallets:
			return listWallets(ctx, path[1:], req, k)
		case QueryGetWallet:
			return getWallet(ctx, path[1:], req, k)
		case QueryListTransactions:
			return listTransactions(ctx, path[1:], req, k)
		case QueryGetTransaction:
			return getTransaction(ctx, path[1:], req, k)
		default:
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "unknown multisig query endpoint")
		}
	}
}

func listWallets(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, error) {
	var wallets types.QueryWallets

	owner, err := sdk.AccAddressFromBech32(path[0])
	if err != nil {
		msgError := fmt.Sprintf("unable to parse owner address: %s", err.Error())
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msgError)
	}

	for iterator := keeper.GetIterator(ctx, "wallet/"); iterator.Valid(); iterator.Next() {
		address := strings.TrimPrefix(string(iterator.Key()), "wallet/")
		wallet := keeper.GetWallet(ctx, address)
		for _, o := range wallet.Owners {
			if o.Equals(owner) {
				wallets = append(wallets, wallet)
				break
			}
		}
	}

	res, err := codec.MarshalJSONIndent(keeper.cdc, wallets)
	if err != nil {
		panic("could not marshal result to JSON")
	}

	return res, nil
}

func getWallet(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, error) {

	wallet := keeper.GetWallet(ctx, path[0])

	res, err := codec.MarshalJSONIndent(keeper.cdc, wallet)
	if err != nil {
		panic("could not marshal result to JSON")
	}

	return res, nil
}

func listTransactions(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, error) {
	var transactionList types.QueryTransactions

	wallet, err := sdk.AccAddressFromBech32(path[0])
	if err != nil {
		msgError := fmt.Sprintf("unable to parse wallet address: %s", err.Error())
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msgError)
	}

	for iterator := keeper.GetIterator(ctx, "transaction/"); iterator.Valid(); iterator.Next() {
		txID := strings.TrimPrefix(string(iterator.Key()), "transaction/")
		transaction := keeper.GetTransaction(ctx, txID)
		if transaction.Wallet.Equals(wallet) {
			transactionList = append(transactionList, transaction)
		}
	}

	res, err := codec.MarshalJSONIndent(keeper.cdc, transactionList)
	if err != nil {
		panic("could not marshal result to JSON")
	}

	return res, nil
}

func getTransaction(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, error) {

	transaction := keeper.GetTransaction(ctx, path[0])

	res, err := codec.MarshalJSONIndent(keeper.cdc, transaction)
	if err != nil {
		panic("could not marshal result to JSON")
	}

	return res, nil
}
