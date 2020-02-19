package keeper

import (
	abci "github.com/tendermint/tendermint/abci/types"

	"bitbucket.org/decimalteam/go-node/x/coin/internal/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewQuerier creates a new querier for coin clients.
func NewQuerier(k Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, sdk.Error) {
		switch path[0] {
		case types.QueryGetCoin:
			return getCoin(ctx, path[1:], k)
		case types.QueryListCoins:
			return listCoins(ctx, k)
		default:
			return nil, sdk.ErrUnknownRequest("unknown coin query endpoint")
		}
	}
}

func RemovePrefixFromHash(key []byte, prefix []byte) (hash []byte) {
	hash = key[len(prefix):]
	return hash
}

func listCoins(ctx sdk.Context, k Keeper) ([]byte, sdk.Error) {
	var coinList types.QueryResCoins

	iterator := k.GetCoinsIterator(ctx)

	for ; iterator.Valid(); iterator.Next() {
		coinHash := RemovePrefixFromHash(iterator.Key(), []byte(types.CoinPrefix))
		coinList = append(coinList, string(coinHash))
	}

	res, err := codec.MarshalJSONIndent(k.cdc, coinList)
	if err != nil {
		return res, sdk.NewError(types.DefaultCodespace, types.CodeInvalid, "Could not marshal result to JSON")
	}

	return res, nil
}

func getCoin(ctx sdk.Context, path []string, k Keeper) (res []byte, sdkError sdk.Error) {
	coinHash := path[0]
	coin, err := k.GetCoin(ctx, coinHash)
	if err != nil {
		return nil, sdk.NewError(types.DefaultCodespace, types.CodeInvalid, err.Error())
	}

	res, err = codec.MarshalJSONIndent(k.cdc, coin)
	if err != nil {
		return nil, sdk.NewError(types.DefaultCodespace, types.CodeInvalid, "Could not marshal result to JSON")
	}

	return res, nil
}
