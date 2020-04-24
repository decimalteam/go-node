package keeper

import (
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"bitbucket.org/decimalteam/go-node/x/coin/internal/types"
)

// NewQuerier creates a new querier for coin clients.
func NewQuerier(k Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		switch path[0] {
		case types.QueryGetCoin:
			return getCoin(ctx, path[1:], k)
		case types.QueryListCoins:
			return listCoins(ctx, k)
		default:
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "unknown coin query endpoint")
		}
	}
}

func RemovePrefixFromHash(key []byte, prefix []byte) (hash []byte) {
	hash = key[len(prefix):]
	return hash
}

func listCoins(ctx sdk.Context, k Keeper) ([]byte, error) {
	var coinList types.QueryResCoins

	iterator := k.GetCoinsIterator(ctx)

	for ; iterator.Valid(); iterator.Next() {
		coinHash := RemovePrefixFromHash(iterator.Key(), []byte(types.CoinPrefix))
		coinList = append(coinList, string(coinHash))
	}

	res, err := codec.MarshalJSONIndent(k.cdc, coinList)
	if err != nil {
		return res, sdkerrors.New(types.DefaultCodespace, types.CodeInvalid, "Could not marshal result to JSON")
	}

	return res, nil
}

func getCoin(ctx sdk.Context, path []string, k Keeper) (res []byte, sdkError error) {
	coinHash := path[0]
	coin, err := k.GetCoin(ctx, coinHash)
	if err != nil {
		return nil, sdkerrors.New(types.DefaultCodespace, types.CodeInvalid, err.Error())
	}

	res, err = codec.MarshalJSONIndent(k.cdc, coin)
	if err != nil {
		return nil, sdkerrors.New(types.DefaultCodespace, types.CodeInvalid, "Could not marshal result to JSON")
	}

	return res, nil
}
