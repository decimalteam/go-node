package keeper

import (
	"bitbucket.org/decimalteam/go-node/x/validator/internal/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	abci "github.com/tendermint/tendermint/abci/types"
	"strings"
)

/// creates a querier for staking REST endpoints
func NewQuerier(k Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err error) {
		switch path[0] {
		case types.QueryValidators:
			return queryValidators(ctx, req, k)
		case types.QueryValidator:
			return queryValidator(ctx, req, k)
		case types.QueryValidatorDelegations:
			return queryValidatorDelegations(ctx, req, k)
		case types.QueryValidatorUnbondingDelegations:
			return queryValidatorUnbondingDelegations(ctx, req, k)
		case types.QueryDelegation:
			return queryDelegation(ctx, req, k)
		case types.QueryUnbondingDelegation:
			return queryUnbondingDelegation(ctx, req, k)
		case types.QueryDelegatorDelegations:
			return queryDelegatorDelegations(ctx, req, k)
		case types.QueryDelegatorUnbondingDelegations:
			return queryDelegatorUnbondingDelegations(ctx, req, k)
		case types.QueryDelegatorValidators:
			return queryDelegatorValidators(ctx, req, k)
		case types.QueryDelegatorValidator:
			return queryDelegatorValidator(ctx, req, k)
		case types.QueryHistoricalInfo:
			return queryHistoricalInfo(ctx, req, k)
		case types.QueryPool:
			return queryPool(ctx, k)
		case types.QueryParameters:
			return queryParameters(ctx, k)
		default:
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "unknown validator query endpoint")
		}
	}
}

func queryValidators(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, error) {
	var params types.QueryValidatorsParams

	err := types.ModuleCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	validators := k.GetAllValidators(ctx)
	filteredVals := make([]types.Validator, 0, len(validators))

	for _, val := range validators {
		if strings.EqualFold(val.GetStatus().String(), params.Status) {
			filteredVals = append(filteredVals, val)
		}
	}

	start, end := client.Paginate(len(filteredVals), params.Page, params.Limit, int(k.GetParams(ctx).MaxValidators))
	if start < 0 || end < 0 {
		filteredVals = []types.Validator{}
	} else {
		filteredVals = filteredVals[start:end]
	}

	res, err := codec.MarshalJSONIndent(types.ModuleCdc, filteredVals)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

func queryValidator(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, error) {
	var params types.QueryValidatorParams

	err := types.ModuleCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	validator, err := k.GetValidator(ctx, params.ValidatorAddr)
	if err != nil {
		return nil, types.ErrNoValidatorFound()
	}

	res, err := codec.MarshalJSONIndent(types.ModuleCdc, validator)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

func queryValidatorDelegations(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, error) {
	var params types.QueryValidatorParams

	err := types.ModuleCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	delegations := k.GetValidatorDelegations(ctx, params.ValidatorAddr)

	resDelegations := types.DelegationResponse{}

	for _, delegation := range delegations {
		switch delegation := delegation.(type) {
		case types.Delegation:
			resDelegations.Delegations = append(resDelegations.Delegations, delegation)
		case types.DelegationNFT:
			resDelegations.DelegationsNFT = append(resDelegations.DelegationsNFT, delegation)
		}
	}

	res, err := codec.MarshalJSONIndent(types.ModuleCdc, resDelegations)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

func queryValidatorUnbondingDelegations(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, error) {
	var params types.QueryValidatorParams

	err := types.ModuleCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	unbondingDelegations := k.GetUnbondingDelegationsFromValidator(ctx, params.ValidatorAddr)
	if unbondingDelegations == nil {
		unbondingDelegations = types.UnbondingDelegations{}
	}

	baseUBDs := types.BaseUnbondingDelegations{}
	nftUBDs := types.NFTUnbondingDelegations{}

	for _, unbondingDelegation := range unbondingDelegations {
		baseUBD := types.BaseUnbondingDelegation{
			ValidatorAddress: unbondingDelegation.ValidatorAddress,
			DelegatorAddress: unbondingDelegation.DelegatorAddress,
			Entries:          []types.UnbondingDelegationEntry{},
		}

		nftUBD := types.NFTUnbondingDelegation{
			ValidatorAddress: unbondingDelegation.ValidatorAddress,
			DelegatorAddress: unbondingDelegation.DelegatorAddress,
			Entries:          []types.UnbondingDelegationNFTEntry{},
		}

		for _, entry := range unbondingDelegation.Entries {
			switch entry := entry.(type) {
			case types.UnbondingDelegationEntry:
				baseUBD.Entries = append(baseUBD.Entries, entry)
			case types.UnbondingDelegationNFTEntry:
				nftUBD.Entries = append(nftUBD.Entries, entry)
			}
		}

		baseUBDs = append(baseUBDs, baseUBD)
		nftUBDs = append(nftUBDs, nftUBD)
	}

	ubdResp := types.NewUnbondingDelegationResp(baseUBDs, nftUBDs)

	res, err := codec.MarshalJSONIndent(types.ModuleCdc, ubdResp)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

func queryDelegatorDelegations(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, error) {
	var params types.QueryDelegatorParams

	err := types.ModuleCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	delegations := k.GetAllDelegatorDelegations(ctx, params.DelegatorAddr)

	resDelegations := types.DelegationResponse{}

	for _, delegation := range delegations {
		switch delegation := delegation.(type) {
		case types.Delegation:
			resDelegations.Delegations = append(resDelegations.Delegations, delegation)
		case types.DelegationNFT:
			resDelegations.DelegationsNFT = append(resDelegations.DelegationsNFT, delegation)
		}
	}

	res, err := codec.MarshalJSONIndent(types.ModuleCdc, resDelegations)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

func queryDelegatorUnbondingDelegations(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, error) {
	var params types.QueryDelegatorParams

	err := types.ModuleCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	unbondingDelegations := k.GetUnbondingDelegationsByDelegator(ctx, params.DelegatorAddr)
	if unbondingDelegations == nil {
		unbondingDelegations = types.UnbondingDelegations{}
	}

	baseUBDs := types.BaseUnbondingDelegations{}
	nftUBDs := types.NFTUnbondingDelegations{}

	for _, unbondingDelegation := range unbondingDelegations {
		baseUBD := types.BaseUnbondingDelegation{
			ValidatorAddress: unbondingDelegation.ValidatorAddress,
			DelegatorAddress: unbondingDelegation.DelegatorAddress,
			Entries:          []types.UnbondingDelegationEntry{},
		}

		nftUBD := types.NFTUnbondingDelegation{
			ValidatorAddress: unbondingDelegation.ValidatorAddress,
			DelegatorAddress: unbondingDelegation.DelegatorAddress,
			Entries:          []types.UnbondingDelegationNFTEntry{},
		}

		for _, entry := range unbondingDelegation.Entries {
			switch entry := entry.(type) {
			case types.UnbondingDelegationEntry:
				baseUBD.Entries = append(baseUBD.Entries, entry)
			case types.UnbondingDelegationNFTEntry:
				nftUBD.Entries = append(nftUBD.Entries, entry)
			}
		}

		baseUBDs = append(baseUBDs, baseUBD)
		nftUBDs = append(nftUBDs, nftUBD)
	}

	ubdResp := types.NewUnbondingDelegationResp(baseUBDs, nftUBDs)

	res, err := codec.MarshalJSONIndent(types.ModuleCdc, ubdResp)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

func queryDelegatorValidators(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, error) {
	var params types.QueryDelegatorParams

	stakingParams := k.GetParams(ctx)

	err := types.ModuleCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	validators := k.GetDelegatorValidators(ctx, params.DelegatorAddr, stakingParams.MaxValidators)
	if validators == nil {
		validators = types.Validators{}
	}

	res, err := codec.MarshalJSONIndent(types.ModuleCdc, validators)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

func queryDelegatorValidator(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, error) {
	var params types.QueryBondsParams

	err := types.ModuleCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	validator, err := k.GetDelegatorValidator(ctx, params.DelegatorAddr, params.ValidatorAddr, params.Coin)
	if err != nil {
		return nil, err
	}

	res, err := codec.MarshalJSONIndent(types.ModuleCdc, validator)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

func queryDelegation(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, error) {
	var params types.QueryBondsParams

	err := types.ModuleCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	delegation, found := k.GetDelegation(ctx, params.DelegatorAddr, params.ValidatorAddr, params.Coin)
	if !found {
		return nil, types.ErrNoDelegation()
	}

	delegationResp := types.NewDelegationResp(types.Delegations{delegation}, types.DelegationsNFT{})

	res, err := codec.MarshalJSONIndent(types.ModuleCdc, delegationResp)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

func queryUnbondingDelegation(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, error) {
	var params types.QueryBondsParams

	err := types.ModuleCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	unbond, found := k.GetUnbondingDelegation(ctx, params.DelegatorAddr, params.ValidatorAddr)
	if !found {
		return nil, types.ErrUnbondingDelegationNotFound()
	}

	baseUBD := types.BaseUnbondingDelegation{
		ValidatorAddress: unbond.ValidatorAddress,
		DelegatorAddress: unbond.DelegatorAddress,
		Entries:          []types.UnbondingDelegationEntry{},
	}

	for _, entry := range unbond.Entries {
		switch entry := entry.(type) {
		case types.UnbondingDelegationEntry:
			baseUBD.Entries = append(baseUBD.Entries, entry)
		}
	}

	ubdResp := types.NewUnbondingDelegationResp(types.BaseUnbondingDelegations{baseUBD}, types.NFTUnbondingDelegations{})

	res, err := codec.MarshalJSONIndent(types.ModuleCdc, ubdResp)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

func queryHistoricalInfo(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, error) {
	var params types.QueryHistoricalInfoParams

	err := types.ModuleCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	hi, found := k.GetHistoricalInfo(ctx, params.Height)
	if !found {
		return nil, types.ErrNoHistoricalInfo(params.Height)
	}

	res, err := codec.MarshalJSONIndent(types.ModuleCdc, hi)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

func queryPool(ctx sdk.Context, k Keeper) ([]byte, error) {
	bondedPool := k.GetBondedPool(ctx)
	notBondedPool := k.GetNotBondedPool(ctx)
	if bondedPool == nil || notBondedPool == nil {
		return nil, types.ErrAccountNotSet()
	}

	pool := types.NewPool(
		notBondedPool.GetCoins(),
		bondedPool.GetCoins(),
	)

	res, err := codec.MarshalJSONIndent(types.ModuleCdc, pool)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

func queryParameters(ctx sdk.Context, k Keeper) ([]byte, error) {
	params := k.GetParams(ctx)

	res, err := codec.MarshalJSONIndent(types.ModuleCdc, params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}
