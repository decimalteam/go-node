package keeper

//import (
//	"bitbucket.org/decimalteam/go-node/x/validator/types"
//	"fmt"
//	"github.com/cosmos/cosmos-sdk/codec"
//	sdk "github.com/cosmos/cosmos-sdk/types"
//	"github.com/stretchr/testify/require"
//	abci "github.com/tendermint/tendermint/abci/types"
//	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
//	"testing"
//)
//
//var (
//	addrAcc1, addrAcc2 = Addrs[0], Addrs[1]
//	addrVal1, addrVal2 = sdk.ValAddress(Addrs[0]), sdk.ValAddress(Addrs[1])
//	pk1, pk2           = PKs[0], PKs[1]
//)
//
//func TestNewQuerier(t *testing.T) {
//	cdc := codec.New()
//	ctx, _, keeper, _, _, _ := CreateTestInput(t, false, 1000)
//	// Create Validators
//	amts := []sdk.Int{sdk.NewInt(9), sdk.NewInt(8)}
//	var validators [2]types.Validator
//	for i, amt := range amts {
//		validators[i] = types.NewValidator(sdk.ValAddress(Addrs[i]), PKs[i], sdk.ZeroDec(), Addrs[i], types.Description{})
//		del := types.NewDelegation(sdk.AccAddress(validators[i].ValAddress), validators[i].ValAddress, sdk.NewCoin(keeper.BondDenom(ctx), amt))
//		keeper.SetDelegation(ctx, del)
//		err := keeper.SetValidator(ctx, validators[i])
//		require.NoError(t, err)
//		keeper.SetValidatorByPowerIndex(ctx, validators[i])
//	}
//
//	header := tmproto.Header{
//		ChainID: "HelloChain",
//		Height:  5,
//	}
//	hi := types.NewHistoricalInfo(header, validators[:])
//	keeper.SetHistoricalInfo(ctx, 5, hi)
//
//	query := abci.RequestQuery{
//		Path: "",
//		Data: []byte{},
//	}
//
//	querier := NewQuerier(keeper)
//
//	bz, err := querier(ctx, []string{"other"}, query)
//	require.Error(t, err)
//	require.Nil(t, bz)
//
//	_, err = querier(ctx, []string{"pool"}, query)
//	require.NoError(t, err)
//
//	_, err = querier(ctx, []string{"parameters"}, query)
//	require.NoError(t, err)
//
//	queryValParams := types.NewQueryValidatorParams(addrVal1)
//	bz, errRes := cdc.MarshalJSON(queryValParams)
//	require.NoError(t, errRes)
//
//	query.Path = "/custom/validator/validator"
//	query.Data = bz
//
//	_, err = querier(ctx, []string{"validator"}, query)
//	require.NoError(t, err)
//
//	_, err = querier(ctx, []string{"validatorDelegations"}, query)
//	require.NoError(t, err)
//
//	_, err = querier(ctx, []string{"validatorUnbondingDelegations"}, query)
//	require.NoError(t, err)
//
//	queryDelParams := types.NewQueryDelegatorParams(addrAcc2)
//	bz, errRes = cdc.MarshalJSON(queryDelParams)
//	require.NoError(t, errRes)
//
//	query.Path = "/custom/validator/validator"
//	query.Data = bz
//
//	_, err = querier(ctx, []string{"delegatorDelegations"}, query)
//	require.NoError(t, err)
//
//	_, err = querier(ctx, []string{"delegatorUnbondingDelegations"}, query)
//	require.NoError(t, err)
//
//	_, err = querier(ctx, []string{"delegatorValidators"}, query)
//	require.NoError(t, err)
//
//	queryHisParams := types.NewQueryHistoricalInfoParams(5)
//	bz, errRes = cdc.MarshalJSON(queryHisParams)
//	require.NoError(t, errRes)
//
//	query.Path = "/custom/validator/historicalInfo"
//	query.Data = bz
//
//	_, err = querier(ctx, []string{"historicalInfo"}, query)
//	require.NoError(t, err)
//}
//
//func TestQueryParametersPool(t *testing.T) {
//	cdc := codec.New()
//	ctx, _, keeper, _, _, _ := CreateTestInput(t, false, 1000)
//
//	res, err := queryParameters(ctx, keeper)
//	require.NoError(t, err)
//
//	var params types.Params
//	errRes := cdc.UnmarshalJSON(res, &params)
//	require.NoError(t, errRes)
//	require.Equal(t, keeper.GetParams(ctx), params)
//
//	res, err = queryPool(ctx, keeper)
//	require.NoError(t, err)
//
//	var pool types.Pool
//	bondedPool := keeper.GetBondedPool(ctx)
//	notBondedPool := keeper.GetNotBondedPool(ctx)
//	errRes = cdc.UnmarshalJSON(res, &pool)
//	require.NoError(t, errRes)
//	require.Equal(t, bondedPool.GetCoins(), pool.BondedTokens)
//	require.Equal(t, notBondedPool.GetCoins(), pool.NotBondedTokens)
//}
//
//func TestQueryValidators(t *testing.T) {
//	cdc := codec.New()
//	ctx, _, keeper, _, _, _ := CreateTestInput(t, false, 10000)
//	params := keeper.GetParams(ctx)
//
//	// Create Validators
//	amts := []sdk.Int{sdk.NewInt(9), sdk.NewInt(8), sdk.NewInt(7)}
//	status := []types.BondStatus{types.Bonded, types.Unbonded, types.Unbonding}
//	var validators [3]types.Validator
//	for i, amt := range amts {
//		validators[i] = types.NewValidator(sdk.ValAddress(Addrs[i]), PKs[i], sdk.ZeroDec(), Addrs[i], types.Description{})
//		del := types.NewDelegation(sdk.AccAddress(validators[i].ValAddress), validators[i].ValAddress, sdk.NewCoin(keeper.BondDenom(ctx), amt))
//		keeper.SetDelegation(ctx, del)
//		validators[i] = validators[i].UpdateStatus(status[i])
//	}
//
//	err := keeper.SetValidator(ctx, validators[0])
//	require.NoError(t, err)
//	err = keeper.SetValidator(ctx, validators[1])
//	require.NoError(t, err)
//	err = keeper.SetValidator(ctx, validators[2])
//	require.NoError(t, err)
//
//	// Query Validators
//	queriedValidators := keeper.GetValidators(ctx, params.MaxValidators)
//
//	for i, s := range status {
//		queryValsParams := types.NewQueryValidatorsParams(1, int(params.MaxValidators), s.String())
//		bz, err := cdc.MarshalJSON(queryValsParams)
//		require.NoError(t, err)
//
//		req := tmproto.RequestQuery{
//			Path: fmt.Sprintf("/custom/%s/%s", types.QuerierRoute, types.QueryValidators),
//			Data: bz,
//		}
//
//		res, err := queryValidators(ctx, req, keeper)
//		require.NoError(t, err)
//
//		var validatorsResp []types.Validator
//		err = cdc.UnmarshalJSON(res, &validatorsResp)
//		require.NoError(t, err)
//
//		require.Equal(t, 1, len(validatorsResp))
//		require.ElementsMatch(t, validators[i].ValAddress, validatorsResp[0].ValAddress)
//
//	}
//
//	// Query each validator
//	queryParams := types.NewQueryValidatorParams(addrVal1)
//	bz, err := cdc.MarshalJSON(queryParams)
//	require.NoError(t, err)
//
//	query := tmproto.RequestQuery{
//		Path: "/custom/validator/validator",
//		Data: bz,
//	}
//	res, err := queryValidator(ctx, query, keeper)
//	require.NoError(t, err)
//
//	var validator types.Validator
//	err = cdc.UnmarshalJSON(res, &validator)
//	require.NoError(t, err)
//
//	require.Equal(t, queriedValidators[0], validator)
//}
//
///*
//func TestQueryDelegation(t *testing.T) {
//	cdc := MakeTestCodec()
//	ctx, _, keeper, _, _ := CreateTestInput(t, false, 10000)
//	params := keeper.GetParams(ctx)
//
//	// Create Validators and Delegation
//	val1 := types.NewValidator(addrVal1, pk1, sdk.ZeroDec(), sdk.AccAddress(addrVal1), types.Description{})
//	err := keeper.SetValidator(ctx, val1)
//	require.NoError(t, err)
//	keeper.SetValidatorByPowerIndex(ctx, val1)
//
//	val2 := types.NewValidator(addrVal2, pk2, sdk.ZeroDec(), sdk.AccAddress(addrVal2), types.Description{})
//	err = keeper.SetValidator(ctx, val2)
//	require.NoError(t, err)
//	keeper.SetValidatorByPowerIndex(ctx, val2)
//
//	delTokens := types.TokensFromConsensusPower(20)
//	err = keeper.Delegate(ctx, addrAcc2, sdk.NewCoin(keeper.BondDenom(ctx), delTokens), types.Unbonded, val1, true)
//	require.NoError(t, err)
//
//	// apply TM updates
//	_, err = keeper.ApplyAndReturnValidatorSetUpdates(ctx)
//	require.NoError(t, err)
//
//	// Query Delegator bonded validators
//	queryParams := types.NewQueryDelegatorParams(addrAcc2)
//	bz, errRes := cdc.MarshalJSON(queryParams)
//	require.NoError(t, errRes)
//
//	query := tmproto.RequestQuery{
//		Path: "/custom/validator/delegatorValidators",
//		Data: bz,
//	}
//
//	delValidators := keeper.GetDelegatorValidators(ctx, addrAcc2, params.MaxValidators)
//
//	res, err := queryDelegatorValidators(ctx, query, keeper)
//	require.NoError(t, err)
//
//	var validatorsResp []types.Validator
//	errRes = cdc.UnmarshalJSON(res, &validatorsResp)
//	require.NoError(t, errRes)
//
//	require.Equal(t, len(delValidators), len(validatorsResp))
//	require.ElementsMatch(t, delValidators, validatorsResp)
//
//	// error unknown request
//	query.Data = bz[:len(bz)-1]
//
//	_, err = queryDelegatorValidators(ctx, query, keeper)
//	require.Error(t, err)
//
//	// Query bonded validator
//	queryBondParams := types.NewQueryBondsParams(addrAcc2, addrVal1, keeper.BondDenom(ctx))
//	bz, errRes = cdc.MarshalJSON(queryBondParams)
//	require.NoError(t, errRes)
//
//	query = tmproto.RequestQuery{
//		Path: "/custom/validator/delegatorValidator",
//		Data: bz,
//	}
//
//	res, err = queryDelegatorValidator(ctx, query, keeper)
//	require.NoError(t, err)
//
//	var validator types.Validator
//	errRes = cdc.UnmarshalJSON(res, &validator)
//	require.NoError(t, errRes)
//
//	require.Equal(t, delValidators[0], validator)
//
//	// error unknown request
//	query.Data = bz[:len(bz)-1]
//
//	_, err = queryDelegatorValidator(ctx, query, keeper)
//	require.Error(t, err)
//
//	// Query delegation
//
//	query = tmproto.RequestQuery{
//		Path: "/custom/validator/delegation",
//		Data: bz,
//	}
//
//	delegation, found := keeper.GetDelegation(ctx, addrAcc2, addrVal1, keeper.BondDenom(ctx))
//	require.True(t, found)
//
//	res, err = queryDelegation(ctx, query, keeper)
//	require.NoError(t, err)
//
//	var delegationRes types.DelegationResponse
//	errRes = cdc.UnmarshalJSON(res, &delegationRes)
//	require.NoError(t, errRes)
//
//	require.Equal(t, delegation.ValidatorAddress, delegationRes.GetValidatorAddr())
//	require.Equal(t, delegation.DelegatorAddress, delegationRes.GetDelegatorAddr())
//	require.Equal(t, delegation.Coin, delegationRes.GetCoin())
//
//	// Query Delegator Delegations
//	query = tmproto.RequestQuery{
//		Path: "/custom/validator/delegatorDelegations",
//		Data: bz,
//	}
//
//	res, err = queryDelegatorDelegations(ctx, query, keeper)
//	require.NoError(t, err)
//
//	var delegatorDelegations types.DelegationResponses
//	errRes = cdc.UnmarshalJSON(res, &delegatorDelegations)
//	require.NoError(t, errRes)
//	require.Len(t, delegatorDelegations, 1)
//	require.Equal(t, delegation.ValidatorAddress, delegatorDelegations[0].GetValidatorAddr())
//	require.Equal(t, delegation.DelegatorAddress, delegatorDelegations[0].GetDelegatorAddr())
//	require.Equal(t, delegation.Coin, delegatorDelegations[0].GetCoin())
//
//	// error unknown request
//	query.Data = bz[:len(bz)-1]
//
//	_, err = queryDelegation(ctx, query, keeper)
//	require.Error(t, err)
//
//	// Query validator delegations
//
//	bz, errRes = cdc.MarshalJSON(types.NewQueryValidatorParams(addrVal1))
//	require.NoError(t, errRes)
//
//	query = tmproto.RequestQuery{
//		Path: "custom/validator/validatorDelegations",
//		Data: bz,
//	}
//
//	res, err = queryValidatorDelegations(ctx, query, keeper)
//	require.NoError(t, err)
//
//	var delegationsRes types.DelegationResponses
//	errRes = cdc.UnmarshalJSON(res, &delegationsRes)
//	require.NoError(t, errRes)
//	require.Len(t, delegatorDelegations, 1)
//	require.Equal(t, delegation.ValidatorAddress, delegationsRes[0].ValidatorAddress)
//	require.Equal(t, delegation.DelegatorAddress, delegationsRes[0].DelegatorAddress)
//	require.Equal(t, delegation.Coin, delegationsRes[0].Coin)
//
//	// Query unbonging delegation
//	unbondingTokens := types.TokensFromConsensusPower(10)
//	_, err = keeper.Undelegate(ctx, addrAcc2, val1.ValAddress, sdk.NewCoin(keeper.BondDenom(ctx), unbondingTokens))
//	require.NoError(t, err)
//
//	queryBondParams = types.NewQueryBondsParams(addrAcc2, addrVal1, keeper.BondDenom(ctx))
//	bz, errRes = cdc.MarshalJSON(queryBondParams)
//	require.NoError(t, errRes)
//
//	query = tmproto.RequestQuery{
//		Path: "/custom/validator/unbondingDelegation",
//		Data: bz,
//	}
//
//	unbond, found := keeper.GetUnbondingDelegation(ctx, addrAcc2, addrVal1)
//	require.True(t, found)
//
//	res, err = queryUnbondingDelegation(ctx, query, keeper)
//	require.NoError(t, err)
//
//	var unbondRes types.UnbondingDelegation
//	errRes = cdc.UnmarshalJSON(res, &unbondRes)
//	require.NoError(t, errRes)
//
//	require.Equal(t, unbond, unbondRes)
//
//	// error unknown request
//	query.Data = bz[:len(bz)-1]
//
//	_, err = queryUnbondingDelegation(ctx, query, keeper)
//	require.Error(t, err)
//
//	// Query Delegator Delegations
//
//	query = tmproto.RequestQuery{
//		Path: "/custom/validator/delegatorUnbondingDelegations",
//		Data: bz,
//	}
//
//	res, err = queryDelegatorUnbondingDelegations(ctx, query, keeper)
//	require.NoError(t, err)
//
//	var delegatorUbds []types.UnbondingDelegation
//	errRes = cdc.UnmarshalJSON(res, &delegatorUbds)
//	require.NoError(t, errRes)
//	require.Equal(t, unbond, delegatorUbds[0])
//
//	// error unknown request
//	query.Data = bz[:len(bz)-1]
//
//	_, err = queryDelegatorUnbondingDelegations(ctx, query, keeper)
//	require.Error(t, err)
//}
//*/
//func TestQueryUnbondingDelegation(t *testing.T) {
//	cdc := MakeTestCodec()
//	ctx, _, keeper, _, _, _ := CreateTestInput(t, false, 10000)
//
//	// Create Validators and Delegation
//	val1 := types.NewValidator(addrVal1, pk1, sdk.ZeroDec(), sdk.AccAddress(addrVal1), types.Description{})
//	err := keeper.SetValidator(ctx, val1)
//	require.NoError(t, err)
//
//	// delegate
//	delAmount := types.TokensFromConsensusPower(100)
//	err = keeper.Delegate(ctx, addrAcc1, sdk.NewCoin(keeper.BondDenom(ctx), delAmount), types.Unbonded, val1, true)
//	require.NoError(t, err)
//	_, err = keeper.ApplyAndReturnValidatorSetUpdates(ctx)
//	require.NoError(t, err)
//
//	// undelegate
//	undelAmount := types.TokensFromConsensusPower(20)
//	_, err = keeper.Undelegate(ctx, addrAcc1, val1.GetOperator(), sdk.NewCoin(keeper.BondDenom(ctx), undelAmount))
//	require.NoError(t, err)
//	_, err = keeper.ApplyAndReturnValidatorSetUpdates(ctx)
//	require.NoError(t, err)
//
//	_, found := keeper.GetUnbondingDelegation(ctx, addrAcc1, val1.ValAddress)
//	require.True(t, found)
//
//	//
//	// found: query unbonding delegation by delegator and validator
//	//
//	queryValidatorParams := types.NewQueryBondsParams(addrAcc1, val1.GetOperator(), keeper.BondDenom(ctx))
//	bz, errRes := cdc.MarshalJSON(queryValidatorParams)
//	require.NoError(t, errRes)
//	query := tmproto.RequestQuery{
//		Path: "/custom/validator/unbondingDelegation",
//		Data: bz,
//	}
//	res, err := queryUnbondingDelegation(ctx, query, keeper)
//	require.NoError(t, err)
//	require.NotNil(t, res)
//	var ubDel types.UnbondingDelegation
//	require.NoError(t, cdc.UnmarshalJSON(res, &ubDel))
//	require.Equal(t, addrAcc1, ubDel.DelegatorAddress)
//	require.Equal(t, val1.ValAddress, ubDel.ValidatorAddress)
//	require.Equal(t, 1, len(ubDel.Entries))
//
//	//
//	// not found: query unbonding delegation by delegator and validator
//	//
//	queryValidatorParams = types.NewQueryBondsParams(addrAcc2, val1.GetOperator(), keeper.BondDenom(ctx))
//	bz, errRes = cdc.MarshalJSON(queryValidatorParams)
//	require.NoError(t, errRes)
//	query = tmproto.RequestQuery{
//		Path: "/custom/validator/unbondingDelegation",
//		Data: bz,
//	}
//	_, err = queryUnbondingDelegation(ctx, query, keeper)
//	require.Error(t, err)
//
//	//
//	// found: query unbonding delegation by delegator and validator
//	//
//	queryDelegatorParams := types.NewQueryDelegatorParams(addrAcc1)
//	bz, errRes = cdc.MarshalJSON(queryDelegatorParams)
//	require.NoError(t, errRes)
//	query = tmproto.RequestQuery{
//		Path: "/custom/validator/delegatorUnbondingDelegations",
//		Data: bz,
//	}
//	res, err = queryDelegatorUnbondingDelegations(ctx, query, keeper)
//	require.NoError(t, err)
//	require.NotNil(t, res)
//	var ubDels []types.UnbondingDelegation
//	require.NoError(t, cdc.UnmarshalJSON(res, &ubDels))
//	require.Equal(t, 1, len(ubDels))
//	require.Equal(t, addrAcc1, ubDels[0].DelegatorAddress)
//	require.Equal(t, val1.ValAddress, ubDels[0].ValidatorAddress)
//
//	//
//	// not found: query unbonding delegation by delegator and validator
//	//
//	queryDelegatorParams = types.NewQueryDelegatorParams(addrAcc2)
//	bz, errRes = cdc.MarshalJSON(queryDelegatorParams)
//	require.NoError(t, errRes)
//	query = tmproto.RequestQuery{
//		Path: "/custom/validator/delegatorUnbondingDelegations",
//		Data: bz,
//	}
//	res, err = queryDelegatorUnbondingDelegations(ctx, query, keeper)
//	require.NoError(t, err)
//	require.NotNil(t, res)
//	require.NoError(t, cdc.UnmarshalJSON(res, &ubDels))
//	require.Equal(t, 0, len(ubDels))
//}
//
//func TestQueryHistoricalInfo(t *testing.T) {
//	cdc := codec.New()
//	ctx, _, keeper, _, _, _ := CreateTestInput(t, false, 10000)
//
//	// Create Validators and Delegation
//	val1 := types.NewValidator(addrVal1, pk1, sdk.ZeroDec(), sdk.AccAddress(addrVal1), types.Description{})
//	val2 := types.NewValidator(addrVal2, pk2, sdk.ZeroDec(), sdk.AccAddress(addrVal2), types.Description{})
//	vals := []types.Validator{val1, val2}
//	err := keeper.SetValidator(ctx, val1)
//	require.NoError(t, err)
//	err = keeper.SetValidator(ctx, val2)
//	require.NoError(t, err)
//
//	header := tmproto.Header{
//		ChainID: "HelloChain",
//		Height:  5,
//	}
//	hi := types.NewHistoricalInfo(header, vals)
//	keeper.SetHistoricalInfo(ctx, 5, hi)
//
//	queryHistoricalParams := types.NewQueryHistoricalInfoParams(4)
//	bz, errRes := cdc.MarshalJSON(queryHistoricalParams)
//	require.NoError(t, errRes)
//	query := tmproto.RequestQuery{
//		Path: "/custom/validator/historicalInfo",
//		Data: bz,
//	}
//	res, err := queryHistoricalInfo(ctx, query, keeper)
//	require.Error(t, err, "Invalid query passed")
//	require.Nil(t, res, "Invalid query returned non-nil result")
//
//	queryHistoricalParams = types.NewQueryHistoricalInfoParams(5)
//	bz, errRes = cdc.MarshalJSON(queryHistoricalParams)
//	require.NoError(t, errRes)
//	query.Data = bz
//	res, err = queryHistoricalInfo(ctx, query, keeper)
//	require.NoError(t, err, "Valid query passed")
//	require.NotNil(t, res, "Valid query returned nil result")
//
//	var recv types.HistoricalInfo
//	require.NoError(t, cdc.UnmarshalJSON(res, &recv))
//	require.Equal(t, hi, recv, "HistoricalInfo query returned wrong result")
//}
