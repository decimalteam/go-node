package keeper

import (
	"bitbucket.org/decimalteam/go-node/x/validator/internal/types"
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// NewQuerier creates a new querier for validator clients.
func NewQuerier(k Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, sdk.Error) {
		switch path[0] {
		case types.QueryParams:
			return queryParams(ctx, k)
		case types.QueryValidatorDelegations:
			return queryValidatorDelegations(ctx, req, k)
		case types.QueryPool:
			return queryPool(ctx, k)
		default:
			return nil, sdk.ErrUnknownRequest("unknown validator query endpoint")
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

func queryValidatorDelegations(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, sdk.Error) {
	var params types.QueryValidatorParams

	err := types.ModuleCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", err))
	}

	delegations := k.GetValidatorDelegations(ctx, params.ValidatorAddr)
	delegationResps, err := delegationsToDelegationResponses(ctx, k, delegations)
	if err != nil {
		return nil, sdk.ErrInternal(err.Error())
	}

	res, err := codec.MarshalJSONIndent(types.ModuleCdc, delegationResps)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("failed to marshal result to JSON", err.Error()))
	}

	return res, nil
}

// util

func delegationsToDelegationResponses(
	ctx sdk.Context, k Keeper, delegations types.Delegations,
) (types.DelegationResponses, sdk.Error) {

	resp := make(types.DelegationResponses, len(delegations))
	for i, del := range delegations {
		delResp, err := delegationToDelegationResponse(del)
		if err != nil {
			return nil, err
		}

		resp[i] = delResp
	}

	return resp, nil
}

func delegationToDelegationResponse(del types.Delegation) (types.DelegationResponse, sdk.Error) {
	return types.NewDelegationResp(
		del.DelegatorAddress,
		del.ValidatorAddress,
		del.Shares,
		del.Coin,
	), nil
}

func queryPool(ctx sdk.Context, k Keeper) ([]byte, sdk.Error) {
	bondedPool := k.GetBondedPool(ctx)
	notBondedPool := k.GetNotBondedPool(ctx)
	if bondedPool == nil || notBondedPool == nil {
		return nil, sdk.ErrInternal("pool accounts haven't been set")
	}

	pool := types.NewPool(
		notBondedPool.GetCoins(),
		bondedPool.GetCoins(),
	)

	res, err := codec.MarshalJSONIndent(types.ModuleCdc, pool)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}

	return res, nil
}
