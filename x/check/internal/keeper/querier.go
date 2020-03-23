package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// NewQuerier creates a new querier for check clients.
func NewQuerier(k Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, sdk.Error) {
		return nil, sdk.ErrUnknownRequest("unknown check query endpoint")
	}
}

//
//func queryParams(ctx sdk.Context, k Keeper) ([]byte, sdk.Error) {
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
//func issueCheck(ctx sdk.Context, keeper Keeper ){
//
//}
