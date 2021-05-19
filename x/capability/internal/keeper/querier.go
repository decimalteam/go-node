package keeper

import (
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"bitbucket.org/decimalteam/go-node/x/capability/internal/types"
)

// NewQuerier creates a new querier for coin clients.
func NewQuerier(k Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		switch path[0] {
		default:
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "unknown coin query endpoint")
		}
	}
}

func listCoins(ctx sdk.Context, k Keeper) ([]byte, error) {

	var coinList types.QueryResCoins
	for it := k.GetCoinsIterator(ctx); it.Valid(); it.Next() {
		coinHash := it.Key()[len(types.CoinPrefix):]
		coinList = append(coinList, string(coinHash))
	}

	res, err := codec.MarshalJSONIndent(k.cdc, coinList)
	if err != nil {
		return res, types.ErrInternal(err.Error())
	}

	return res, nil
}

func getCoin(ctx sdk.Context, path []string, k Keeper) (res []byte, sdkError error) {
	coinHash := path[0]

	coin, err := k.GetCoin(ctx, coinHash)
	if err != nil {
		return nil, types.ErrCoinDoesNotExist(coinHash)
	}

	res, err = codec.MarshalJSONIndent(k.cdc, coin)
	if err != nil {
		return nil, types.ErrInternal(err.Error())
	}

	return res, nil
}
