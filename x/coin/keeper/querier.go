package keeper

import (
	types2 "bitbucket.org/decimalteam/go-node/x/coin/types"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// NewQuerier creates a new querier for coin clients.
func NewQuerier(k Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		switch path[0] {
		case types2.QueryGetCoin:
			return getCoin(ctx, path[1:], k)
		case types2.QueryListCoins:
			return listCoins(ctx, k)
		default:
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "unknown coin query endpoint")
		}
	}
}

func listCoins(ctx sdk.Context, k Keeper) ([]byte, error) {

	var coinList types2.QueryResCoins
	for it := k.GetCoinsIterator(ctx); it.Valid(); it.Next() {
		coinHash := it.Key()[len(types2.CoinPrefix):]
		coinList = append(coinList, string(coinHash))
	}

	res, err := codec.MarshalJSONIndent(k.cdc, coinList)
	if err != nil {
		return res, types2.ErrInternal(err.Error())
	}

	return res, nil
}

func getCoin(ctx sdk.Context, path []string, k Keeper) (res []byte, sdkError error) {
	coinHash := path[0]

	coin, err := k.GetCoin(ctx, coinHash)
	if err != nil {
		return nil, types2.ErrCoinDoesNotExist(coinHash)
	}

	res, err = codec.MarshalJSONIndent(k.cdc, coin)
	if err != nil {
		return nil, types2.ErrInternal(err.Error())
	}

	return res, nil
}
