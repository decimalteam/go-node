package keeper

import (
	types2 "bitbucket.org/decimalteam/go-node/x/swap/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	abci "github.com/tendermint/tendermint/abci/types"
	"time"
)

/// creates a querier for swap REST endpoints
func NewQuerier(k Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err error) {
		switch path[0] {
		case types2.QuerySwap:
			return querySwap(ctx, req, k)
		case types2.QueryActiveSwaps:
			return queryActiveSwaps(ctx, k)
		case types2.QueryPool:
			return queryPool(ctx, k)
		default:
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "unknown swap query endpoint")
		}
	}
}

func querySwap(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, error) {
	var params types2.QuerySwapParams

	err := types2.ModuleCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	swap, ok := k.GetSwap(ctx, params.HashedSecret)
	if !ok {
		return nil, types2.ErrSwapNotFound()
	}

	res, err := codec.MarshalJSONIndent(types2.ModuleCdc, swap)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

func queryActiveSwaps(ctx sdk.Context, k Keeper) ([]byte, error) {
	var activeSwaps types2.Swaps
	swaps := k.GetAllSwaps(ctx)
	for _, swap := range swaps {
		if ctx.BlockTime().Sub(time.Unix(0, int64(swap.Timestamp))) <= k.LockedTimeOut(ctx) {
			activeSwaps = append(activeSwaps, swap)
		}
	}

	res, err := codec.MarshalJSONIndent(types2.ModuleCdc, activeSwaps)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

func queryPool(ctx sdk.Context, k Keeper) ([]byte, error) {
	pool := k.GetPool(ctx)

	res, err := codec.MarshalJSONIndent(types2.ModuleCdc, pool)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}
