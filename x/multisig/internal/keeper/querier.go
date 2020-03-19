package keeper

import (
	abci "github.com/tendermint/tendermint/abci/types"

	"bitbucket.org/decimalteam/go-node/x/multisig/internal/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewQuerier creates a new querier for multisig clients.
func NewQuerier(k Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, sdk.Error) {
		switch path[0] {
		case types.QueryParams:
			return queryParams(ctx, k)
			// TODO: Put the modules query routes
		default:
			return nil, sdk.ErrUnknownRequest("unknown multisig query endpoint")
		}
	}
}

func queryParams(ctx sdk.Context, k Keeper) ([]byte, sdk.Error) {
	params := k.GetParams(ctx)

	res, err := codec.MarshalJSONIndent(types.ModuleCdc, params)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("failed to marshal JSON", err.Error()))
	}

	return res, nil
}

// TODO: Add the modules query functions
// They will be similar to the above one: queryParams()
