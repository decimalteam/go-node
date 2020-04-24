package keeper

import (
	abci "github.com/tendermint/tendermint/abci/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// NewQuerier creates a new querier for multisig clients.
func NewQuerier(k Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "unknown multisig query endpoint")
	}
}

//func queryParams(ctx sdk.Context, k Keeper) ([]byte, error) {
//	params := k.GetParams(ctx)
//
//	res, err := codec.MarshalJSONIndent(types.ModuleCdc, params)
//	if err != nil {
//		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("failed to marshal JSON", err.Error()))
//	}
//
//	return res, nil
//}

// TODO: Add the modules query functions
// They will be similar to the above one: queryParams()
