package keeper

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"

	decsdk "bitbucket.org/decimalteam/go-node/utils/types"
	"bitbucket.org/decimalteam/go-node/x/validator/internal/types"
)

//_______________________________________________________

func TestSetValidator(t *testing.T) {
	ctx, _, keeper, _, _ := CreateTestInput(t, false, 10)

	valPubKey := PKs[0]
	valAddr := decsdk.ValAddress(valPubKey.Address().Bytes())
	valTokens := sdk.TokensFromConsensusPower(10)

	// test how the validator is set from a purely unbonbed pool
	validator := types.NewValidator(valAddr, valPubKey, sdk.ZeroDec(), decsdk.AccAddress(valAddr), types.Description{})
	delegator := types.NewDelegation(addrDels[0], valAddr, sdk.NewCoin(keeper.BondDenom(ctx), valTokens))
	keeper.SetDelegation(ctx, delegator)
	require.Equal(t, types.Unbonded, validator.Status)
	require.Equal(t, valTokens, keeper.TotalStake(ctx, validator))
	err := keeper.SetValidator(ctx, validator)
	require.Nil(t, err)
	keeper.SetValidatorByPowerIndex(ctx, validator)

	// ensure update
	updates, err := keeper.ApplyAndReturnValidatorSetUpdates(ctx)
	require.Nil(t, err)
	validator, err = keeper.GetValidator(ctx, valAddr)
	require.Nil(t, err)
	require.Equal(t, 1, len(updates))
	require.Equal(t, validator.ABCIValidatorUpdate(keeper.TotalStake(ctx, validator)), updates[0])

	// after the save the validator should be bonded
	require.Equal(t, types.Bonded, validator.Status)
	assert.Equal(t, valTokens, keeper.TotalStake(ctx, validator))

	// Check each store for being saved
	resVal, err := keeper.GetValidator(ctx, valAddr)
	require.Nil(t, err)
	assert.True(ValEq(t, validator, resVal))

	resVals := keeper.GetLastValidators(ctx)
	require.Equal(t, 1, len(resVals))
	assert.True(ValEq(t, validator, resVals[0]))

	resVals = keeper.GetBondedValidatorsByPower(ctx)
	require.Equal(t, 1, len(resVals))
	require.True(ValEq(t, validator, resVals[0]))

	resVals = keeper.GetValidators(ctx, 1)
	require.Equal(t, 1, len(resVals))
	require.True(ValEq(t, validator, resVals[0]))

	resVals = keeper.GetValidators(ctx, 10)
	require.Equal(t, 1, len(resVals))
	require.True(ValEq(t, validator, resVals[0]))

	allVals := keeper.GetAllValidators(ctx)
	require.Equal(t, 1, len(allVals))
}

func TestUpdateValidatorByPowerIndex(t *testing.T) {
	ctx, _, keeper, _, _ := CreateTestInput(t, false, 0)

	bondedPool := keeper.GetBondedPool(ctx)
	notBondedPool := keeper.GetNotBondedPool(ctx)
	bondedPool.SetCoins(sdk.NewCoins(sdk.NewCoin(keeper.BondDenom(ctx), sdk.TokensFromConsensusPower(1234))))
	notBondedPool.SetCoins(sdk.NewCoins(sdk.NewCoin(keeper.BondDenom(ctx), sdk.TokensFromConsensusPower(10000))))
	keeper.supplyKeeper.SetModuleAccount(ctx, bondedPool)
	keeper.supplyKeeper.SetModuleAccount(ctx, notBondedPool)

	// add a validator
	validator := types.NewValidator(addrVals[0], PKs[0], sdk.ZeroDec(), decsdk.AccAddress(addrVals[0]), types.Description{})
	delegator := types.NewDelegation(addrDels[0], addrVals[0], sdk.NewCoin(keeper.BondDenom(ctx), sdk.TokensFromConsensusPower(100)))
	keeper.SetDelegation(ctx, delegator)
	require.Equal(t, types.Unbonded, validator.Status)
	require.Equal(t, sdk.TokensFromConsensusPower(100), keeper.TotalStake(ctx, validator))
	TestingUpdateValidator(keeper, ctx, validator, true)
	validator, err := keeper.GetValidator(ctx, addrVals[0])
	require.Nil(t, err)
	require.Equal(t, sdk.TokensFromConsensusPower(100), keeper.TotalStake(ctx, validator))

	power := types.GetValidatorsByPowerIndexKey(validator, keeper.TotalStake(ctx, validator))
	require.True(t, validatorByPowerIndexExists(keeper, ctx, power))

	// burn half the delegator shares
	keeper.DeleteValidatorByPowerIndex(ctx, validator)
	sdkErr := keeper.unbond(ctx, addrDels[0], addrVals[0], sdk.NewCoin(keeper.BondDenom(ctx), sdk.TokensFromConsensusPower(50)))
	require.Nil(t, sdkErr)
	TestingUpdateValidator(keeper, ctx, validator, true) // update the validator, possibly kicking it out
	require.False(t, validatorByPowerIndexExists(keeper, ctx, power))

	validator, err = keeper.GetValidator(ctx, addrVals[0])
	require.Nil(t, err)

	power = types.GetValidatorsByPowerIndexKey(validator, keeper.TotalStake(ctx, validator))
	require.True(t, validatorByPowerIndexExists(keeper, ctx, power))
}

func TestUpdateBondedValidatorsDecreaseCliff(t *testing.T) {
	numVals := 10
	maxVals := 5

	// create context, keeper, and pool for tests
	ctx, _, keeper, _, _ := CreateTestInput(t, false, 0)
	bondedPool := keeper.GetBondedPool(ctx)
	notBondedPool := keeper.GetNotBondedPool(ctx)

	// create keeper parameters
	params := keeper.GetParams(ctx)
	params.MaxValidators = uint16(maxVals)
	keeper.SetParams(ctx, params)

	// create a random pool
	bondedPool.SetCoins(sdk.NewCoins(sdk.NewCoin(keeper.BondDenom(ctx), sdk.TokensFromConsensusPower(1234))))
	notBondedPool.SetCoins(sdk.NewCoins(sdk.NewCoin(keeper.BondDenom(ctx), sdk.TokensFromConsensusPower(10000))))
	keeper.supplyKeeper.SetModuleAccount(ctx, bondedPool)
	keeper.supplyKeeper.SetModuleAccount(ctx, notBondedPool)

	validators := make([]types.Validator, numVals)
	for i := 0; i < len(validators); i++ {
		val := types.NewValidator(decsdk.ValAddress(Addrs[i]), PKs[i], sdk.ZeroDec(), Addrs[i], types.Description{})
		delTokens := sdk.TokensFromConsensusPower(int64((i + 1) * 10))
		delegator := types.NewDelegation(Addrs[i], decsdk.ValAddress(Addrs[i]), sdk.NewCoin(keeper.BondDenom(ctx), delTokens))
		keeper.SetDelegation(ctx, delegator)

		val = TestingUpdateValidator(keeper, ctx, val, true)
		validators[i] = val
	}

	nextCliffVal := validators[numVals-maxVals+1]

	// remove enough tokens to kick out the validator below the current cliff
	// validator and next in line cliff validator
	keeper.DeleteValidatorByPowerIndex(ctx, nextCliffVal)
	shares := sdk.TokensFromConsensusPower(21)
	del, found := keeper.GetDelegation(ctx, decsdk.AccAddress(nextCliffVal.ValAddress), nextCliffVal.ValAddress)
	require.True(t, found)
	del.Coin = del.Coin.Sub(sdk.NewCoin(keeper.BondDenom(ctx), shares))
	keeper.SetDelegation(ctx, del)
	nextCliffVal = TestingUpdateValidator(keeper, ctx, nextCliffVal, true)

	expectedValStatus := map[int]types.BondStatus{
		9: types.Bonded, 8: types.Bonded, 7: types.Bonded, 5: types.Bonded, 4: types.Bonded,
		0: types.Unbonded, 1: types.Unbonded, 2: types.Unbonded, 3: types.Unbonded, 6: types.Unbonded,
	}

	// require all the validators have their respective statuses
	for valIdx, status := range expectedValStatus {
		valAddr := validators[valIdx].ValAddress
		val, _ := keeper.GetValidator(ctx, valAddr)

		assert.Equal(
			t, status, val.GetStatus(),
			fmt.Sprintf("expected validator at index %v to have status: %s", valIdx, status),
		)
	}
}

func TestSlashToZeroPowerRemoved(t *testing.T) {
	// initialize setup
	ctx, _, keeper, _, _ := CreateTestInput(t, false, 100)

	// add a validator
	validator := types.NewValidator(addrVals[0], PKs[0], sdk.ZeroDec(), decsdk.AccAddress(addrVals[0]), types.Description{})
	valTokens := sdk.TokensFromConsensusPower(100)

	bondedPool := keeper.GetBondedPool(ctx)
	err := bondedPool.SetCoins(sdk.NewCoins(sdk.NewCoin(keeper.BondDenom(ctx), valTokens)))
	require.NoError(t, err)
	keeper.supplyKeeper.SetModuleAccount(ctx, bondedPool)

	delegator := types.NewDelegation(addrDels[0], addrVals[0], sdk.NewCoin(keeper.BondDenom(ctx), valTokens))
	keeper.SetDelegation(ctx, delegator)
	require.Equal(t, types.Unbonded, validator.Status)
	require.Equal(t, valTokens, keeper.TotalStake(ctx, validator))
	keeper.SetValidatorByConsAddr(ctx, validator)
	validator = TestingUpdateValidator(keeper, ctx, validator, true)
	require.Equal(t, valTokens, keeper.TotalStake(ctx, validator), "\nvalidator %v\npool %v", validator, valTokens)

	// slash the validator by 100%
	consAddr0 := decsdk.ConsAddress(PKs[0].Address())
	keeper.Slash(ctx, consAddr0, 0, 100, sdk.OneDec())
	// apply TM updates
	keeper.ApplyAndReturnValidatorSetUpdates(ctx)
	// validator should be unbonding
	validator, _ = keeper.GetValidator(ctx, addrVals[0])
	require.Equal(t, validator.GetStatus(), types.Unbonded)
}

// This function tests UpdateValidator, GetValidator, GetLastValidators, RemoveValidator
func TestValidatorBasics(t *testing.T) {
	ctx, _, keeper, _, _ := CreateTestInput(t, false, 1000)

	//construct the validators
	var validators [3]types.Validator
	powers := []int64{9, 8, 7}
	for i, power := range powers {
		validators[i] = types.NewValidator(addrVals[i], PKs[i], sdk.ZeroDec(), decsdk.AccAddress(addrVals[i]), types.Description{})
		validators[i].Status = types.Unbonded
		tokens := sdk.TokensFromConsensusPower(power)

		delegator := types.NewDelegation(decsdk.AccAddress(addrVals[i]), addrVals[i], sdk.NewCoin(keeper.BondDenom(ctx), tokens))
		keeper.SetDelegation(ctx, delegator)
	}
	assert.Equal(t, sdk.TokensFromConsensusPower(9), keeper.TotalStake(ctx, validators[0]))
	assert.Equal(t, sdk.TokensFromConsensusPower(8), keeper.TotalStake(ctx, validators[1]))
	assert.Equal(t, sdk.TokensFromConsensusPower(7), keeper.TotalStake(ctx, validators[2]))

	// check the empty keeper first
	_, err := keeper.GetValidator(ctx, addrVals[0])
	require.Error(t, err)
	resVals := keeper.GetLastValidators(ctx)
	require.Zero(t, len(resVals))

	resVals = keeper.GetValidators(ctx, 2)
	require.Zero(t, len(resVals))

	// set and retrieve a record
	validators[0] = TestingUpdateValidator(keeper, ctx, validators[0], true)
	keeper.SetValidatorByConsAddr(ctx, validators[0])
	resVal, err := keeper.GetValidator(ctx, addrVals[0])
	require.Nil(t, err)
	assert.True(ValEq(t, validators[0], resVal))

	// retrieve from consensus
	resVal, err = keeper.GetValidatorByConsAddr(ctx, decsdk.ConsAddress(PKs[0].Address()))
	require.Nil(t, err)
	assert.True(ValEq(t, validators[0], resVal))
	resVal, err = keeper.GetValidatorByConsAddr(ctx, decsdk.GetConsAddress(PKs[0]))
	require.Nil(t, err)
	assert.True(ValEq(t, validators[0], resVal))

	resVals = keeper.GetLastValidators(ctx)
	require.Equal(t, 1, len(resVals))
	assert.True(ValEq(t, validators[0], resVals[0]))
	assert.Equal(t, types.Bonded, validators[0].Status)
	assert.True(sdk.IntEq(t, sdk.TokensFromConsensusPower(9), keeper.TotalStake(ctx, validators[0])))

	// modify a records, save, and retrieve
	validators[0].Status = types.Bonded
	validators[0].Tokens = sdk.TokensFromConsensusPower(10)
	validators[0].DelegatorShares = validators[0].Tokens.ToDec()
	validators[0] = TestingUpdateValidator(keeper, ctx, validators[0], true)
	resVal, err = keeper.GetValidator(ctx, addrVals[0])
	require.Nil(t, err)
	assert.True(ValEq(t, validators[0], resVal))

	resVals = keeper.GetLastValidators(ctx)
	require.Equal(t, 1, len(resVals))
	assert.True(ValEq(t, validators[0], resVals[0]))

	// add other validators
	validators[1] = TestingUpdateValidator(keeper, ctx, validators[1], true)
	validators[2] = TestingUpdateValidator(keeper, ctx, validators[2], true)
	resVal, err = keeper.GetValidator(ctx, addrVals[1])
	require.Nil(t, err)
	assert.True(ValEq(t, validators[1], resVal))
	resVal, err = keeper.GetValidator(ctx, addrVals[2])
	require.Nil(t, err)
	assert.True(ValEq(t, validators[2], resVal))

	resVals = keeper.GetLastValidators(ctx)
	require.Equal(t, 3, len(resVals))
	assert.True(ValEq(t, validators[0], resVals[0])) // order doesn't matter here
	assert.True(ValEq(t, validators[1], resVals[1]))
	assert.True(ValEq(t, validators[2], resVals[2]))

	// remove a record

	// shouldn't be able to remove if status is not unbonded
	err = keeper.RemoveValidator(ctx, validators[1].ValAddress)
	assert.EqualError(t, err,
		"cannot call RemoveValidator on bonded or unbonding validators")

	// shouldn't be able to remove if there are still tokens left
	validators[1].Status = types.Unbonded
	keeper.SetValidator(ctx, validators[1])
	err = keeper.RemoveValidator(ctx, validators[1].ValAddress)
	assert.EqualError(t, err,
		"attempting to remove a validator which still contains tokens")

	del, found := keeper.GetDelegation(ctx, decsdk.AccAddress(validators[1].ValAddress), validators[1].ValAddress)
	require.True(t, found)
	keeper.RemoveDelegation(ctx, del)                     // ...remove all tokens
	keeper.SetValidator(ctx, validators[1])               // ...set the validator
	keeper.RemoveValidator(ctx, validators[1].ValAddress) // Now it can be removed.
	_, err = keeper.GetValidator(ctx, addrVals[1])
	require.Error(t, err)
}

// test how the validators are sorted, tests GetBondedValidatorsByPower
func TestGetValidatorSortingUnmixed(t *testing.T) {
	ctx, _, keeper, _, _ := CreateTestInput(t, false, 1000)

	// initialize some validators into the state
	amts := []sdk.Int{
		sdk.TokensFromConsensusPower(0),
		sdk.TokensFromConsensusPower(100),
		sdk.TokensFromConsensusPower(1),
		sdk.TokensFromConsensusPower(400),
		sdk.TokensFromConsensusPower(200),
	}
	n := len(amts)
	var validators [5]types.Validator
	for i, amt := range amts {
		validators[i] = types.NewValidator(decsdk.ValAddress(Addrs[i]), PKs[i], sdk.ZeroDec(), Addrs[i], types.Description{})
		validators[i].Status = types.Bonded
		delegator := types.NewDelegation(Addrs[i], decsdk.ValAddress(Addrs[i]), sdk.NewCoin(keeper.BondDenom(ctx), amt))
		keeper.SetDelegation(ctx, delegator)
		TestingUpdateValidator(keeper, ctx, validators[i], true)
	}

	// first make sure everything made it in to the gotValidator group
	resValidators := keeper.GetBondedValidatorsByPower(ctx)
	assert.Equal(t, n, len(resValidators))
	assert.Equal(t, sdk.TokensFromConsensusPower(400), keeper.TotalStake(ctx, resValidators[0]), "%v", resValidators)
	assert.Equal(t, sdk.TokensFromConsensusPower(200), keeper.TotalStake(ctx, resValidators[1]), "%v", resValidators)
	assert.Equal(t, sdk.TokensFromConsensusPower(100), keeper.TotalStake(ctx, resValidators[2]), "%v", resValidators)
	assert.Equal(t, sdk.TokensFromConsensusPower(1), keeper.TotalStake(ctx, resValidators[3]), "%v", resValidators)
	assert.Equal(t, sdk.TokensFromConsensusPower(0), keeper.TotalStake(ctx, resValidators[4]), "%v", resValidators)
	assert.Equal(t, validators[3].ValAddress, resValidators[0].ValAddress, "%v", resValidators)
	assert.Equal(t, validators[4].ValAddress, resValidators[1].ValAddress, "%v", resValidators)
	assert.Equal(t, validators[1].ValAddress, resValidators[2].ValAddress, "%v", resValidators)
	assert.Equal(t, validators[2].ValAddress, resValidators[3].ValAddress, "%v", resValidators)
	assert.Equal(t, validators[0].ValAddress, resValidators[4].ValAddress, "%v", resValidators)

	// test a basic increase in voting power
	del, found := keeper.GetDelegation(ctx, decsdk.AccAddress(validators[3].ValAddress), validators[3].ValAddress)
	require.True(t, found)
	del.Coin.Amount = sdk.TokensFromConsensusPower(500)
	keeper.SetDelegation(ctx, del)
	TestingUpdateValidator(keeper, ctx, validators[3], true)
	resValidators = keeper.GetBondedValidatorsByPower(ctx)
	require.Equal(t, len(resValidators), n)
	assert.True(ValEq(t, validators[3], resValidators[0]))

	// test a decrease in voting power
	del, found = keeper.GetDelegation(ctx, decsdk.AccAddress(validators[3].ValAddress), validators[3].ValAddress)
	require.True(t, found)
	del.Coin.Amount = sdk.TokensFromConsensusPower(300)
	keeper.SetDelegation(ctx, del)
	TestingUpdateValidator(keeper, ctx, validators[3], true)
	resValidators = keeper.GetBondedValidatorsByPower(ctx)
	require.Equal(t, len(resValidators), n)
	assert.True(ValEq(t, validators[3], resValidators[0]))
	assert.True(ValEq(t, validators[4], resValidators[1]))

	// test equal voting power, different age
	del, found = keeper.GetDelegation(ctx, decsdk.AccAddress(validators[3].ValAddress), validators[3].ValAddress)
	require.True(t, found)
	del.Coin.Amount = sdk.TokensFromConsensusPower(200)
	keeper.SetDelegation(ctx, del)
	ctx = ctx.WithBlockHeight(10)
	TestingUpdateValidator(keeper, ctx, validators[3], true)
	resValidators = keeper.GetBondedValidatorsByPower(ctx)
	require.Equal(t, len(resValidators), n)
	assert.True(ValEq(t, validators[3], resValidators[0]))
	assert.True(ValEq(t, validators[4], resValidators[1]))

	// no change in voting power - no change in sort
	ctx = ctx.WithBlockHeight(20)
	TestingUpdateValidator(keeper, ctx, validators[4], true)
	resValidators = keeper.GetBondedValidatorsByPower(ctx)
	require.Equal(t, len(resValidators), n)
	assert.True(ValEq(t, validators[3], resValidators[0]))
	assert.True(ValEq(t, validators[4], resValidators[1]))

	// change in voting power of both validators, both still in v-set, no age change
	del, found = keeper.GetDelegation(ctx, decsdk.AccAddress(validators[3].ValAddress), validators[3].ValAddress)
	require.True(t, found)
	del.Coin.Amount = sdk.TokensFromConsensusPower(300)
	keeper.SetDelegation(ctx, del)
	del, found = keeper.GetDelegation(ctx, decsdk.AccAddress(validators[4].ValAddress), validators[4].ValAddress)
	require.True(t, found)
	del.Coin.Amount = sdk.TokensFromConsensusPower(300)
	keeper.SetDelegation(ctx, del)
	TestingUpdateValidator(keeper, ctx, validators[3], true)
	resValidators = keeper.GetBondedValidatorsByPower(ctx)
	require.Equal(t, len(resValidators), n)
	ctx = ctx.WithBlockHeight(30)
	TestingUpdateValidator(keeper, ctx, validators[4], true)
	resValidators = keeper.GetBondedValidatorsByPower(ctx)
	require.Equal(t, len(resValidators), n, "%v", resValidators)
	assert.True(ValEq(t, validators[3], resValidators[0]))
	assert.True(ValEq(t, validators[4], resValidators[1]))
}

//func TestGetValidatorSortingMixed(t *testing.T) {
//	ctx, _, keeper, _, _ := CreateTestInput(t, false, 1000)
//
//	// now 2 max resValidators
//	params := keeper.GetParams(ctx)
//	params.MaxValidators = 2
//	keeper.SetParams(ctx, params)
//
//	// initialize some validators into the state
//	amts := []int64{0, 100, 1, 400, 200}
//
//	n := len(amts)
//	var validators [5]types.Validator
//	for i, amt := range amts {
//		validators[i] = types.NewValidator(decsdk.ValAddress(Addrs[i]), PKs[i], sdk.ZeroDec(), Addrs[i])
//		delegator := types.NewDelegation(Addrs[i], decsdk.ValAddress(Addrs[i]), sdk.ZeroDec(), sdk.NewCoin(keeper.BondDenom(ctx), sdk.NewInt(amt)))
//		keeper.SetDelegation(ctx, delegator)
//	}
//
//	validators[0].Status = types.Bonded
//	validators[1].Status = types.Bonded
//	validators[2].Status = types.Bonded
//	validators[3].Status = types.Bonded
//	validators[4].Status = types.Bonded
//
//	for i := range amts {
//		TestingUpdateValidator(keeper, ctx, validators[i], true)
//	}
//	val0, err := keeper.GetValidator(ctx, decsdk.ValAddress(Addrs[0]))
//	require.Nil(t, err)
//	val1, err := keeper.GetValidator(ctx, decsdk.ValAddress(Addrs[1]))
//	require.Nil(t, err)
//	val2, err := keeper.GetValidator(ctx, decsdk.ValAddress(Addrs[2]))
//	require.Nil(t, err)
//	val3, err := keeper.GetValidator(ctx, decsdk.ValAddress(Addrs[3]))
//	require.Nil(t, err)
//	val4, err := keeper.GetValidator(ctx, decsdk.ValAddress(Addrs[4]))
//	require.Nil(t, err)
//	require.Equal(t, types.Unbonded, val0.Status)
//	require.Equal(t, types.Unbonded, val1.Status)
//	require.Equal(t, types.Unbonded, val2.Status)
//	require.Equal(t, types.Bonded, val3.Status)
//	require.Equal(t, types.Bonded, val4.Status)
//
//	// first make sure everything made it in to the gotValidator group
//	resValidators := keeper.GetBondedValidatorsByPower(ctx)
//	assert.Equal(t, n, len(resValidators))
//	assert.Equal(t, sdk.NewInt(400), resValidators[0].BondedTokens(), "%v", resValidators)
//	assert.Equal(t, sdk.NewInt(200), resValidators[1].BondedTokens(), "%v", resValidators)
//	assert.Equal(t, sdk.NewInt(100), resValidators[2].BondedTokens(), "%v", resValidators)
//	assert.Equal(t, sdk.NewInt(1), resValidators[3].BondedTokens(), "%v", resValidators)
//	assert.Equal(t, sdk.NewInt(0), resValidators[4].BondedTokens(), "%v", resValidators)
//	assert.Equal(t, validators[3].ValAddress, resValidators[0].ValAddress, "%v", resValidators)
//	assert.Equal(t, validators[4].ValAddress, resValidators[1].ValAddress, "%v", resValidators)
//	assert.Equal(t, validators[1].ValAddress, resValidators[2].ValAddress, "%v", resValidators)
//	assert.Equal(t, validators[2].ValAddress, resValidators[3].ValAddress, "%v", resValidators)
//	assert.Equal(t, validators[0].ValAddress, resValidators[4].ValAddress, "%v", resValidators)
//}
//
//// TODO separate out into multiple tests
func TestGetValidatorsEdgeCases(t *testing.T) {
	ctx, _, keeper, _, _ := CreateTestInput(t, false, 1000)

	// set max validators to 2
	params := keeper.GetParams(ctx)
	nMax := uint16(2)
	params.MaxValidators = nMax
	keeper.SetParams(ctx, params)

	// initialize some validators into the state
	powers := []int64{0, 100, 400, 400}
	var validators [4]types.Validator
	for i, power := range powers {
		validators[i] = types.NewValidator(decsdk.ValAddress(Addrs[i]), PKs[i], sdk.ZeroDec(), Addrs[i], types.Description{})
		tokens := sdk.TokensFromConsensusPower(power)
		delegator := types.NewDelegation(Addrs[i], decsdk.ValAddress(Addrs[i]), sdk.NewCoin(keeper.BondDenom(ctx), tokens))
		keeper.SetDelegation(ctx, delegator)
		notBondedPool := keeper.GetNotBondedPool(ctx)
		require.NoError(t, notBondedPool.SetCoins(notBondedPool.GetCoins().Add(sdk.NewCoin(params.BondDenom, tokens))))
		keeper.supplyKeeper.SetModuleAccount(ctx, notBondedPool)
		validators[i] = TestingUpdateValidator(keeper, ctx, validators[i], true)
	}

	// ensure that the first two bonded validators are the largest validators
	resValidators := keeper.GetBondedValidatorsByPower(ctx)
	require.Equal(t, nMax, uint16(len(resValidators)))
	assert.True(ValEq(t, validators[2], resValidators[0]))
	assert.True(ValEq(t, validators[3], resValidators[1]))

	// delegate 500 tokens to validator 0
	keeper.DeleteValidatorByPowerIndex(ctx, validators[0])
	delTokens := sdk.TokensFromConsensusPower(500)
	delegation := types.NewDelegation(Addrs[0], decsdk.ValAddress(Addrs[0]), sdk.NewCoin(keeper.BondDenom(ctx), delTokens))
	keeper.SetDelegation(ctx, delegation)
	notBondedPool := keeper.GetNotBondedPool(ctx)
	newTokens := sdk.NewCoins(sdk.NewCoin(params.BondDenom, delTokens))
	require.NoError(t, notBondedPool.SetCoins(notBondedPool.GetCoins().Add(newTokens...)))
	keeper.supplyKeeper.SetModuleAccount(ctx, notBondedPool)

	// test that the two largest validators are
	//   a) validator 0 with 500 tokens
	//   b) validator 2 with 400 tokens (delegated before validator 3)
	validators[0] = TestingUpdateValidator(keeper, ctx, validators[0], true)
	resValidators = keeper.GetBondedValidatorsByPower(ctx)
	require.Equal(t, nMax, uint16(len(resValidators)))
	assert.True(ValEq(t, validators[0], resValidators[0]))
	assert.True(ValEq(t, validators[2], resValidators[1]))

	// A validator which leaves the bonded validator set due to a decrease in voting power,
	// then increases to the original voting power, does not get its spot back in the
	// case of a tie.
	//
	// Order of operations for this test:
	//  - validator 3 enter validator set with 1 new token
	//  - validator 3 removed validator set by removing 201 tokens (validator 2 enters)
	//  - validator 3 adds 200 tokens (equal to validator 2 now) and does not get its spot back

	// validator 3 enters bonded validator set
	ctx = ctx.WithBlockHeight(40)

	var err error
	validators[3], err = keeper.GetValidator(ctx, validators[3].ValAddress)
	require.Nil(t, err)
	keeper.DeleteValidatorByPowerIndex(ctx, validators[3])
	delegation = types.NewDelegation(Addrs[3], decsdk.ValAddress(Addrs[3]), sdk.NewCoin(keeper.BondDenom(ctx), sdk.TokensFromConsensusPower(1)))
	keeper.SetDelegation(ctx, delegation)

	notBondedPool = keeper.GetNotBondedPool(ctx)
	newTokens = sdk.NewCoins(sdk.NewCoin(params.BondDenom, sdk.TokensFromConsensusPower(1)))
	require.NoError(t, notBondedPool.SetCoins(notBondedPool.GetCoins().Add(newTokens...)))
	keeper.supplyKeeper.SetModuleAccount(ctx, notBondedPool)

	validators[3] = TestingUpdateValidator(keeper, ctx, validators[3], true)
	resValidators = keeper.GetBondedValidatorsByPower(ctx)
	require.Equal(t, nMax, uint16(len(resValidators)))
	assert.True(ValEq(t, validators[0], resValidators[0]))
	assert.True(ValEq(t, validators[2], resValidators[1]))

	// validator 3 kicked out temporarily
	keeper.DeleteValidatorByPowerIndex(ctx, validators[3])
	rmTokens := sdk.NewInt(201)
	del, found := keeper.GetDelegation(ctx, decsdk.AccAddress(validators[3].ValAddress), validators[3].ValAddress)
	require.True(t, found)
	del.Coin.Amount = del.Coin.Amount.Sub(sdk.NewInt(201))
	keeper.SetDelegation(ctx, del)

	bondedPool := keeper.GetBondedPool(ctx)
	require.NoError(t, bondedPool.SetCoins(bondedPool.GetCoins().Add(sdk.NewCoin(params.BondDenom, rmTokens))))
	keeper.supplyKeeper.SetModuleAccount(ctx, bondedPool)

	validators[3] = TestingUpdateValidator(keeper, ctx, validators[3], true)
	resValidators = keeper.GetBondedValidatorsByPower(ctx)
	require.Equal(t, nMax, uint16(len(resValidators)))
	assert.True(ValEq(t, validators[0], resValidators[0]))
	assert.True(ValEq(t, validators[2], resValidators[1]))

	// validator 3 does not get spot back
	keeper.DeleteValidatorByPowerIndex(ctx, validators[3])
	delegation = types.NewDelegation(Addrs[3], decsdk.ValAddress(Addrs[3]), sdk.NewCoin(keeper.BondDenom(ctx), sdk.TokensFromConsensusPower(200)))
	keeper.SetDelegation(ctx, delegation)

	notBondedPool = keeper.GetNotBondedPool(ctx)
	require.NoError(t, notBondedPool.SetCoins(notBondedPool.GetCoins().Add(sdk.NewCoin(params.BondDenom, sdk.NewInt(200)))))
	keeper.supplyKeeper.SetModuleAccount(ctx, notBondedPool)

	validators[3] = TestingUpdateValidator(keeper, ctx, validators[3], true)
	resValidators = keeper.GetBondedValidatorsByPower(ctx)
	require.Equal(t, nMax, uint16(len(resValidators)))
	assert.True(ValEq(t, validators[0], resValidators[0]))
	assert.True(ValEq(t, validators[2], resValidators[1]))
	_, err = keeper.GetValidator(ctx, validators[3].ValAddress)
	require.Nil(t, err)
}

func TestValidatorBondHeight(t *testing.T) {
	ctx, _, keeper, _, _ := CreateTestInput(t, false, 1000)

	// now 2 max resValidators
	params := keeper.GetParams(ctx)
	params.MaxValidators = 2
	keeper.SetParams(ctx, params)

	// initialize some validators into the state
	var validators [3]types.Validator
	validators[0] = types.NewValidator(decsdk.ValAddress(Addrs[0]), PKs[0], sdk.ZeroDec(), Addrs[0], types.Description{})
	validators[1] = types.NewValidator(decsdk.ValAddress(Addrs[1]), PKs[1], sdk.ZeroDec(), Addrs[1], types.Description{})
	validators[2] = types.NewValidator(decsdk.ValAddress(Addrs[2]), PKs[2], sdk.ZeroDec(), Addrs[2], types.Description{})

	tokens0 := sdk.TokensFromConsensusPower(200)
	tokens1 := sdk.TokensFromConsensusPower(100)
	tokens2 := sdk.TokensFromConsensusPower(100)
	delegation := types.NewDelegation(Addrs[0], decsdk.ValAddress(Addrs[0]), sdk.NewCoin(keeper.BondDenom(ctx), tokens0))
	keeper.SetDelegation(ctx, delegation)
	delegation = types.NewDelegation(Addrs[1], decsdk.ValAddress(Addrs[1]), sdk.NewCoin(keeper.BondDenom(ctx), tokens1))
	keeper.SetDelegation(ctx, delegation)
	delegation = types.NewDelegation(Addrs[2], decsdk.ValAddress(Addrs[2]), sdk.NewCoin(keeper.BondDenom(ctx), tokens2))
	keeper.SetDelegation(ctx, delegation)

	validators[0] = TestingUpdateValidator(keeper, ctx, validators[0], true)

	////////////////////////////////////////
	// If two validators both increase to the same voting power in the same block,
	// the one with the first transaction should become bonded
	validators[1] = TestingUpdateValidator(keeper, ctx, validators[1], true)
	validators[2] = TestingUpdateValidator(keeper, ctx, validators[2], true)

	resValidators := keeper.GetBondedValidatorsByPower(ctx)
	require.Equal(t, uint16(len(resValidators)), params.MaxValidators)

	assert.True(ValEq(t, validators[0], resValidators[0]))
	assert.True(ValEq(t, validators[1], resValidators[1]))
	keeper.DeleteValidatorByPowerIndex(ctx, validators[1])
	keeper.DeleteValidatorByPowerIndex(ctx, validators[2])
	delTokens := sdk.TokensFromConsensusPower(50)
	delegation = types.NewDelegation(Addrs[1], decsdk.ValAddress(Addrs[1]), sdk.NewCoin(keeper.BondDenom(ctx), delTokens))
	keeper.SetDelegation(ctx, delegation)
	delegation = types.NewDelegation(Addrs[2], decsdk.ValAddress(Addrs[2]), sdk.NewCoin(keeper.BondDenom(ctx), delTokens))
	keeper.SetDelegation(ctx, delegation)
	validators[2] = TestingUpdateValidator(keeper, ctx, validators[2], true)
	resValidators = keeper.GetBondedValidatorsByPower(ctx)
	require.Equal(t, params.MaxValidators, uint16(len(resValidators)))
	validators[1] = TestingUpdateValidator(keeper, ctx, validators[1], true)
	assert.True(ValEq(t, validators[0], resValidators[0]))
	assert.True(ValEq(t, validators[2], resValidators[1]))
}

func TestFullValidatorSetPowerChange(t *testing.T) {
	ctx, _, keeper, _, _ := CreateTestInput(t, false, 1000)
	params := keeper.GetParams(ctx)
	max := 2
	params.MaxValidators = uint16(2)
	keeper.SetParams(ctx, params)

	// initialize some validators into the state
	powers := []int64{0, 100, 400, 400, 200}
	var validators [5]types.Validator
	for i, power := range powers {
		validators[i] = types.NewValidator(decsdk.ValAddress(Addrs[i]), PKs[i], sdk.ZeroDec(), Addrs[i], types.Description{})
		tokens := sdk.TokensFromConsensusPower(power)
		delegation := types.NewDelegation(Addrs[i], decsdk.ValAddress(Addrs[i]), sdk.NewCoin(keeper.BondDenom(ctx), tokens))
		keeper.SetDelegation(ctx, delegation)
		TestingUpdateValidator(keeper, ctx, validators[i], true)
	}
	var err error
	for i := range powers {
		validators[i], err = keeper.GetValidator(ctx, validators[i].ValAddress)
		require.NoError(t, err)
	}
	assert.Equal(t, types.Unbonded, validators[0].Status)
	assert.Equal(t, types.Unbonded, validators[1].Status)
	assert.Equal(t, types.Bonded, validators[2].Status)
	assert.Equal(t, types.Bonded, validators[3].Status)
	assert.Equal(t, types.Unbonded, validators[4].Status)
	resValidators := keeper.GetBondedValidatorsByPower(ctx)
	assert.Equal(t, max, len(resValidators))
	assert.True(ValEq(t, validators[2], resValidators[0])) // in the order of txs
	assert.True(ValEq(t, validators[3], resValidators[1]))

	// test a swap in voting power

	tokens := sdk.TokensFromConsensusPower(600)
	delegation := types.NewDelegation(Addrs[0], decsdk.ValAddress(Addrs[0]), sdk.NewCoin(keeper.BondDenom(ctx), tokens))
	keeper.SetDelegation(ctx, delegation)
	validators[0] = TestingUpdateValidator(keeper, ctx, validators[0], true)
	resValidators = keeper.GetBondedValidatorsByPower(ctx)
	assert.Equal(t, max, len(resValidators))
	assert.True(ValEq(t, validators[0], resValidators[0]))
	assert.True(ValEq(t, validators[2], resValidators[1]))
}

func TestApplyAndReturnValidatorSetUpdatesAllNone(t *testing.T) {
	ctx, _, keeper, _, _ := CreateTestInput(t, false, 1000)

	powers := []int64{10, 20}
	var validators [2]types.Validator
	for i, power := range powers {
		valPubKey := PKs[i+1]
		valAddr := decsdk.ValAddress(valPubKey.Address().Bytes())

		validators[i] = types.NewValidator(valAddr, valPubKey, sdk.ZeroDec(), decsdk.AccAddress(valAddr), types.Description{})
		tokens := sdk.TokensFromConsensusPower(power)
		delegation := types.NewDelegation(decsdk.AccAddress(valAddr), valAddr, sdk.NewCoin(keeper.BondDenom(ctx), tokens))
		keeper.SetDelegation(ctx, delegation)
	}

	// test from nothing to something
	//  tendermintUpdate set: {} -> {c1, c3}
	updates, err := keeper.ApplyAndReturnValidatorSetUpdates(ctx)
	require.NoError(t, err)
	require.Equal(t, 0, len(updates))
	keeper.SetValidator(ctx, validators[0])
	keeper.SetValidatorByPowerIndex(ctx, validators[0])
	keeper.SetValidator(ctx, validators[1])
	keeper.SetValidatorByPowerIndex(ctx, validators[1])

	updates, err = keeper.ApplyAndReturnValidatorSetUpdates(ctx)
	require.NoError(t, err)
	assert.Equal(t, 2, len(updates))
	validators[0], _ = keeper.GetValidator(ctx, validators[0].ValAddress)
	validators[1], _ = keeper.GetValidator(ctx, validators[1].ValAddress)
	assert.Equal(t, validators[0].ABCIValidatorUpdate(keeper.TotalStake(ctx, validators[0])), updates[1])
	assert.Equal(t, validators[1].ABCIValidatorUpdate(keeper.TotalStake(ctx, validators[1])), updates[0])
}

func TestApplyAndReturnValidatorSetUpdatesIdentical(t *testing.T) {
	ctx, _, keeper, _, _ := CreateTestInput(t, false, 1000)

	powers := []int64{10, 20}
	var validators [2]types.Validator
	for i, power := range powers {
		validators[i] = types.NewValidator(decsdk.ValAddress(Addrs[i]), PKs[i], sdk.ZeroDec(), Addrs[i], types.Description{})

		tokens := sdk.TokensFromConsensusPower(power)
		delegation := types.NewDelegation(Addrs[i], decsdk.ValAddress(Addrs[i]), sdk.NewCoin(keeper.BondDenom(ctx), tokens))
		keeper.SetDelegation(ctx, delegation)
	}
	validators[0] = TestingUpdateValidator(keeper, ctx, validators[0], false)
	validators[1] = TestingUpdateValidator(keeper, ctx, validators[1], false)
	updates, err := keeper.ApplyAndReturnValidatorSetUpdates(ctx)
	require.NoError(t, err)
	require.Equal(t, 2, len(updates))

	// test identical,
	//  tendermintUpdate set: {} -> {}
	validators[0] = TestingUpdateValidator(keeper, ctx, validators[0], false)
	validators[1] = TestingUpdateValidator(keeper, ctx, validators[1], false)
	updates, err = keeper.ApplyAndReturnValidatorSetUpdates(ctx)
	require.NoError(t, err)
	require.Equal(t, 0, len(updates))
}

func TestApplyAndReturnValidatorSetUpdatesSingleValueChange(t *testing.T) {
	ctx, _, keeper, _, _ := CreateTestInput(t, false, 1000)

	powers := []int64{10, 20}
	var validators [2]types.Validator
	for i, power := range powers {

		validators[i] = types.NewValidator(decsdk.ValAddress(Addrs[i]), PKs[i], sdk.ZeroDec(), Addrs[i], types.Description{})

		tokens := sdk.TokensFromConsensusPower(power)
		delegation := types.NewDelegation(Addrs[i], decsdk.ValAddress(Addrs[i]), sdk.NewCoin(keeper.BondDenom(ctx), tokens))
		keeper.SetDelegation(ctx, delegation)
	}
	validators[0] = TestingUpdateValidator(keeper, ctx, validators[0], false)
	validators[1] = TestingUpdateValidator(keeper, ctx, validators[1], false)
	updates, err := keeper.ApplyAndReturnValidatorSetUpdates(ctx)
	require.NoError(t, err)
	require.Equal(t, 2, len(updates))

	// test single value change
	//  tendermintUpdate set: {} -> {c1'}
	validators[0].Status = types.Bonded
	del, found := keeper.GetDelegation(ctx, decsdk.AccAddress(validators[0].ValAddress), validators[0].ValAddress)
	require.True(t, found)
	del.Coin.Amount = sdk.TokensFromConsensusPower(600)
	keeper.SetDelegation(ctx, del)
	validators[0] = TestingUpdateValidator(keeper, ctx, validators[0], false)

	updates, err = keeper.ApplyAndReturnValidatorSetUpdates(ctx)
	require.NoError(t, err)

	require.Equal(t, 1, len(updates))
	require.Equal(t, validators[0].ABCIValidatorUpdate(keeper.TotalStake(ctx, validators[0])), updates[0])
}

func TestApplyAndReturnValidatorSetUpdatesMultipleValueChange(t *testing.T) {
	ctx, _, keeper, _, _ := CreateTestInput(t, false, 1000)

	powers := []int64{10, 20}
	var validators [2]types.Validator
	for i, power := range powers {
		validators[i] = types.NewValidator(decsdk.ValAddress(Addrs[i]), PKs[i], sdk.ZeroDec(), Addrs[i], types.Description{})

		tokens := sdk.TokensFromConsensusPower(power)
		delegation := types.NewDelegation(Addrs[i], decsdk.ValAddress(Addrs[i]), sdk.NewCoin(keeper.BondDenom(ctx), tokens))
		keeper.SetDelegation(ctx, delegation)
	}
	validators[0] = TestingUpdateValidator(keeper, ctx, validators[0], false)
	validators[1] = TestingUpdateValidator(keeper, ctx, validators[1], false)
	updates, err := keeper.ApplyAndReturnValidatorSetUpdates(ctx)
	require.NoError(t, err)
	require.Equal(t, 2, len(updates))

	// test multiple value change
	//  tendermintUpdate set: {c1, c3} -> {c1', c3'}
	delTokens1 := sdk.TokensFromConsensusPower(190)
	delTokens2 := sdk.TokensFromConsensusPower(80)
	delegation := types.NewDelegation(Addrs[0], decsdk.ValAddress(Addrs[0]), sdk.NewCoin(keeper.BondDenom(ctx), delTokens1))
	keeper.SetDelegation(ctx, delegation)
	delegation = types.NewDelegation(Addrs[1], decsdk.ValAddress(Addrs[1]), sdk.NewCoin(keeper.BondDenom(ctx), delTokens2))
	keeper.SetDelegation(ctx, delegation)
	validators[0] = TestingUpdateValidator(keeper, ctx, validators[0], false)
	validators[1] = TestingUpdateValidator(keeper, ctx, validators[1], false)

	updates, err = keeper.ApplyAndReturnValidatorSetUpdates(ctx)
	require.NoError(t, err)
	require.Equal(t, 2, len(updates))
	require.Equal(t, validators[0].ABCIValidatorUpdate(keeper.TotalStake(ctx, validators[0])), updates[0])
	require.Equal(t, validators[1].ABCIValidatorUpdate(keeper.TotalStake(ctx, validators[1])), updates[1])
}

func TestApplyAndReturnValidatorSetUpdatesInserted(t *testing.T) {
	ctx, _, keeper, _, _ := CreateTestInput(t, false, 1000)

	powers := []int64{10, 20, 5, 15, 25}
	var validators [5]types.Validator
	for i, power := range powers {
		validators[i] = types.NewValidator(decsdk.ValAddress(Addrs[i]), PKs[i], sdk.ZeroDec(), Addrs[i], types.Description{})

		tokens := sdk.TokensFromConsensusPower(power)
		delegation := types.NewDelegation(Addrs[i], decsdk.ValAddress(Addrs[i]), sdk.NewCoin(keeper.BondDenom(ctx), tokens))
		keeper.SetDelegation(ctx, delegation)
	}

	validators[0] = TestingUpdateValidator(keeper, ctx, validators[0], false)
	validators[1] = TestingUpdateValidator(keeper, ctx, validators[1], false)
	updates, err := keeper.ApplyAndReturnValidatorSetUpdates(ctx)
	require.NoError(t, err)
	require.Equal(t, 2, len(updates))

	// test validtor added at the beginning
	//  tendermintUpdate set: {} -> {c0}
	keeper.SetValidator(ctx, validators[2])
	keeper.SetValidatorByPowerIndex(ctx, validators[2])
	updates, err = keeper.ApplyAndReturnValidatorSetUpdates(ctx)
	require.NoError(t, err)
	validators[2], _ = keeper.GetValidator(ctx, validators[2].ValAddress)
	require.Equal(t, 1, len(updates))
	require.Equal(t, validators[2].ABCIValidatorUpdate(keeper.TotalStake(ctx, validators[2])), updates[0])

	// test validtor added at the beginning
	//  tendermintUpdate set: {} -> {c0}
	keeper.SetValidator(ctx, validators[3])
	keeper.SetValidatorByPowerIndex(ctx, validators[3])
	updates, err = keeper.ApplyAndReturnValidatorSetUpdates(ctx)
	require.NoError(t, err)
	validators[3], _ = keeper.GetValidator(ctx, validators[3].ValAddress)
	require.Equal(t, 1, len(updates))
	require.Equal(t, validators[3].ABCIValidatorUpdate(keeper.TotalStake(ctx, validators[3])), updates[0])

	// test validtor added at the end
	//  tendermintUpdate set: {} -> {c0}
	keeper.SetValidator(ctx, validators[4])
	keeper.SetValidatorByPowerIndex(ctx, validators[4])
	updates, err = keeper.ApplyAndReturnValidatorSetUpdates(ctx)
	require.NoError(t, err)
	validators[4], _ = keeper.GetValidator(ctx, validators[4].ValAddress)
	require.Equal(t, 1, len(updates))
	require.Equal(t, validators[4].ABCIValidatorUpdate(keeper.TotalStake(ctx, validators[4])), updates[0])
}

func TestApplyAndReturnValidatorSetUpdatesWithCliffValidator(t *testing.T) {
	ctx, _, keeper, _, _ := CreateTestInput(t, false, 1000)
	params := types.DefaultParams()
	params.MaxValidators = 2
	keeper.SetParams(ctx, params)

	powers := []int64{10, 20, 5}
	var validators [5]types.Validator
	for i, power := range powers {
		validators[i] = types.NewValidator(decsdk.ValAddress(Addrs[i]), PKs[i], sdk.ZeroDec(), Addrs[i], types.Description{})

		tokens := sdk.TokensFromConsensusPower(power)
		delegation := types.NewDelegation(Addrs[i], decsdk.ValAddress(Addrs[i]), sdk.NewCoin(keeper.BondDenom(ctx), tokens))
		keeper.SetDelegation(ctx, delegation)
	}
	validators[0] = TestingUpdateValidator(keeper, ctx, validators[0], false)
	validators[1] = TestingUpdateValidator(keeper, ctx, validators[1], false)
	updates, err := keeper.ApplyAndReturnValidatorSetUpdates(ctx)
	require.NoError(t, err)
	require.Equal(t, 2, len(updates))

	// test validator added at the end but not inserted in the valset
	//  tendermintUpdate set: {} -> {}
	TestingUpdateValidator(keeper, ctx, validators[2], false)
	updates, err = keeper.ApplyAndReturnValidatorSetUpdates(ctx)
	require.NoError(t, err)
	require.Equal(t, 0, len(updates))

	// test validator change its power and become a gotValidator (pushing out an existing)
	//  tendermintUpdate set: {}     -> {c0, c4}
	updates, err = keeper.ApplyAndReturnValidatorSetUpdates(ctx)
	require.NoError(t, err)
	require.Equal(t, 0, len(updates))

	tokens := sdk.TokensFromConsensusPower(10)
	del, found := keeper.GetDelegation(ctx, decsdk.AccAddress(validators[2].ValAddress), validators[2].ValAddress)
	require.True(t, found)
	del.Coin.Amount = del.Coin.Amount.Add(tokens)
	keeper.SetDelegation(ctx, del)
	keeper.SetValidator(ctx, validators[2])
	keeper.SetValidatorByPowerIndex(ctx, validators[2])
	updates, err = keeper.ApplyAndReturnValidatorSetUpdates(ctx)
	require.NoError(t, err)
	validators[2], _ = keeper.GetValidator(ctx, validators[2].ValAddress)
	require.Equal(t, 2, len(updates), "%v", updates)
	require.Equal(t, validators[0].ABCIValidatorUpdateZero(), updates[1])
	require.Equal(t, validators[2].ABCIValidatorUpdate(keeper.TotalStake(ctx, validators[2])), updates[0])
}

func TestApplyAndReturnValidatorSetUpdatesPowerDecrease(t *testing.T) {
	ctx, _, keeper, _, _ := CreateTestInput(t, false, 1000)

	powers := []int64{100, 100}
	var validators [2]types.Validator
	for i, power := range powers {
		validators[i] = types.NewValidator(decsdk.ValAddress(Addrs[i]), PKs[i], sdk.ZeroDec(), Addrs[i], types.Description{})

		tokens := sdk.TokensFromConsensusPower(power)
		delegation := types.NewDelegation(Addrs[i], decsdk.ValAddress(Addrs[i]), sdk.NewCoin(keeper.BondDenom(ctx), tokens))
		keeper.SetDelegation(ctx, delegation)
	}
	validators[0] = TestingUpdateValidator(keeper, ctx, validators[0], false)
	validators[1] = TestingUpdateValidator(keeper, ctx, validators[1], false)
	updates, err := keeper.ApplyAndReturnValidatorSetUpdates(ctx)
	require.NoError(t, err)
	require.Equal(t, 2, len(updates))

	// check initial power
	require.Equal(t, int64(100), validators[0].ConsensusPower(keeper.TotalStake(ctx, validators[0])))
	require.Equal(t, int64(100), validators[1].ConsensusPower(keeper.TotalStake(ctx, validators[1])))

	// test multiple value change
	//  tendermintUpdate set: {c1, c3} -> {c1', c3'}
	delTokens1 := sdk.TokensFromConsensusPower(20)
	delTokens2 := sdk.TokensFromConsensusPower(30)
	del, found := keeper.GetDelegation(ctx, decsdk.AccAddress(validators[0].ValAddress), validators[0].ValAddress)
	require.True(t, found)
	del.Coin.Amount = del.Coin.Amount.Sub(delTokens1)
	keeper.SetDelegation(ctx, del)
	del, found = keeper.GetDelegation(ctx, decsdk.AccAddress(validators[1].ValAddress), validators[1].ValAddress)
	require.True(t, found)
	del.Coin.Amount = del.Coin.Amount.Sub(delTokens2)
	keeper.SetDelegation(ctx, del)
	validators[0] = TestingUpdateValidator(keeper, ctx, validators[0], false)
	validators[1] = TestingUpdateValidator(keeper, ctx, validators[1], false)

	// power has changed
	require.Equal(t, int64(80), validators[0].ConsensusPower(keeper.TotalStake(ctx, validators[0])))
	require.Equal(t, int64(70), validators[1].ConsensusPower(keeper.TotalStake(ctx, validators[1])))

	// Tendermint updates should reflect power change
	updates, err = keeper.ApplyAndReturnValidatorSetUpdates(ctx)
	require.NoError(t, err)
	require.Equal(t, 2, len(updates))
	require.Equal(t, validators[0].ABCIValidatorUpdate(keeper.TotalStake(ctx, validators[0])), updates[0])
	require.Equal(t, validators[1].ABCIValidatorUpdate(keeper.TotalStake(ctx, validators[1])), updates[1])
}

func TestApplyAndReturnValidatorSetUpdatesNewValidator(t *testing.T) {
	ctx, _, keeper, _, _ := CreateTestInput(t, false, 1000)
	params := keeper.GetParams(ctx)
	params.MaxValidators = uint16(3)

	keeper.SetParams(ctx, params)

	powers := []int64{100, 100}
	var validators [2]types.Validator

	// initialize some validators into the state
	for i, power := range powers {

		valPubKey := PKs[i+1]
		valAddr := decsdk.ValAddress(valPubKey.Address().Bytes())

		validators[i] = types.NewValidator(valAddr, valPubKey, sdk.ZeroDec(), decsdk.AccAddress(valAddr), types.Description{})

		tokens := sdk.TokensFromConsensusPower(power)
		delegation := types.NewDelegation(decsdk.AccAddress(valAddr), valAddr, sdk.NewCoin(keeper.BondDenom(ctx), tokens))
		keeper.SetDelegation(ctx, delegation)

		keeper.SetValidator(ctx, validators[i])
		keeper.SetValidatorByPowerIndex(ctx, validators[i])
	}

	// verify initial Tendermint updates are correct
	updates, err := keeper.ApplyAndReturnValidatorSetUpdates(ctx)
	require.NoError(t, err)
	require.Equal(t, len(validators), len(updates))
	validators[0], _ = keeper.GetValidator(ctx, validators[0].ValAddress)
	validators[1], _ = keeper.GetValidator(ctx, validators[1].ValAddress)
	require.Equal(t, validators[0].ABCIValidatorUpdate(keeper.TotalStake(ctx, validators[0])), updates[0])
	require.Equal(t, validators[1].ABCIValidatorUpdate(keeper.TotalStake(ctx, validators[1])), updates[1])

	updates, err = keeper.ApplyAndReturnValidatorSetUpdates(ctx)
	require.NoError(t, err)
	require.Equal(t, 0, len(updates))

	// update initial validator set
	for i, power := range powers {
		keeper.DeleteValidatorByPowerIndex(ctx, validators[i])
		tokens := sdk.TokensFromConsensusPower(power)
		del, found := keeper.GetDelegation(ctx, decsdk.AccAddress(validators[i].ValAddress), validators[i].ValAddress)
		require.True(t, found)
		del.Coin.Amount = del.Coin.Amount.Add(tokens)
		keeper.SetDelegation(ctx, del)

		keeper.SetValidator(ctx, validators[i])
		keeper.SetValidatorByPowerIndex(ctx, validators[i])
	}

	// add a new validator that goes from zero power, to non-zero power, back to
	// zero power
	valPubKey := PKs[len(validators)+1]
	valAddr := decsdk.ValAddress(valPubKey.Address().Bytes())
	amt := sdk.NewInt(100)

	validator := types.NewValidator(valAddr, valPubKey, sdk.ZeroDec(), decsdk.AccAddress(valAddr), types.Description{})
	delegation := types.NewDelegation(decsdk.AccAddress(valAddr), valAddr, sdk.NewCoin(keeper.BondDenom(ctx), amt))
	keeper.SetDelegation(ctx, delegation)

	keeper.SetValidator(ctx, validator)

	del, found := keeper.GetDelegation(ctx, decsdk.AccAddress(validator.ValAddress), validator.ValAddress)
	require.True(t, found)
	del.Coin.Amount = del.Coin.Amount.Sub(amt)
	keeper.SetDelegation(ctx, del)

	keeper.SetValidator(ctx, validator)
	keeper.SetValidatorByPowerIndex(ctx, validator)

	// add a new validator that increases in power
	valPubKey = PKs[len(validators)+2]
	valAddr = valPubKey.Address().Bytes()

	validator = types.NewValidator(valAddr, valPubKey, sdk.ZeroDec(), decsdk.AccAddress(valAddr), types.Description{})
	tokens := sdk.TokensFromConsensusPower(500)
	delegation = types.NewDelegation(decsdk.AccAddress(valAddr), valAddr, sdk.NewCoin(keeper.BondDenom(ctx), tokens))
	keeper.SetDelegation(ctx, delegation)
	keeper.SetValidator(ctx, validator)
	keeper.SetValidatorByPowerIndex(ctx, validator)

	// verify initial Tendermint updates are correct
	updates, err = keeper.ApplyAndReturnValidatorSetUpdates(ctx)
	require.NoError(t, err)
	validator, _ = keeper.GetValidator(ctx, validator.ValAddress)
	validators[0], _ = keeper.GetValidator(ctx, validators[0].ValAddress)
	validators[1], _ = keeper.GetValidator(ctx, validators[1].ValAddress)
	require.Equal(t, len(validators)+1, len(updates))
	require.Equal(t, validator.ABCIValidatorUpdate(keeper.TotalStake(ctx, validator)), updates[0])
	require.Equal(t, validators[0].ABCIValidatorUpdate(keeper.TotalStake(ctx, validators[0])), updates[1])
	require.Equal(t, validators[1].ABCIValidatorUpdate(keeper.TotalStake(ctx, validators[1])), updates[2])
}

func TestApplyAndReturnValidatorSetUpdatesBondTransition(t *testing.T) {
	ctx, _, keeper, _, _ := CreateTestInput(t, false, 1000)
	params := keeper.GetParams(ctx)
	params.MaxValidators = uint16(2)

	keeper.SetParams(ctx, params)

	powers := []int64{100, 200, 300}
	var validators [3]types.Validator

	// initialize some validators into the state
	for i, power := range powers {
		valPubKey := PKs[i+1]
		valAddr := decsdk.ValAddress(valPubKey.Address().Bytes())

		validators[i] = types.NewValidator(valAddr, valPubKey, sdk.ZeroDec(), decsdk.AccAddress(valAddr), types.Description{})
		tokens := sdk.TokensFromConsensusPower(power)
		delegation := types.NewDelegation(decsdk.AccAddress(valAddr), valAddr, sdk.NewCoin(keeper.BondDenom(ctx), tokens))
		keeper.SetDelegation(ctx, delegation)
		keeper.SetValidator(ctx, validators[i])
		keeper.SetValidatorByPowerIndex(ctx, validators[i])
	}

	// verify initial Tendermint updates are correct
	updates, err := keeper.ApplyAndReturnValidatorSetUpdates(ctx)
	require.NoError(t, err)
	require.Equal(t, 2, len(updates))
	validators[2], _ = keeper.GetValidator(ctx, validators[2].ValAddress)
	validators[1], _ = keeper.GetValidator(ctx, validators[1].ValAddress)
	require.Equal(t, validators[2].ABCIValidatorUpdate(keeper.TotalStake(ctx, validators[2])), updates[0])
	require.Equal(t, validators[1].ABCIValidatorUpdate(keeper.TotalStake(ctx, validators[1])), updates[1])

	updates, err = keeper.ApplyAndReturnValidatorSetUpdates(ctx)
	require.NoError(t, err)
	require.Equal(t, 0, len(updates))

	// delegate to validator with lowest power but not enough to bond
	ctx = ctx.WithBlockHeight(1)

	validators[0], err = keeper.GetValidator(ctx, validators[0].ValAddress)
	require.NoError(t, err)

	keeper.DeleteValidatorByPowerIndex(ctx, validators[0])
	tokens := sdk.TokensFromConsensusPower(1)
	del, found := keeper.GetDelegation(ctx, decsdk.AccAddress(validators[0].ValAddress), validators[0].ValAddress)
	require.True(t, found)
	del.Coin.Amount = del.Coin.Amount.Add(tokens)
	keeper.SetDelegation(ctx, del)
	keeper.SetValidator(ctx, validators[0])
	keeper.SetValidatorByPowerIndex(ctx, validators[0])

	// verify initial Tendermint updates are correct
	updates, err = keeper.ApplyAndReturnValidatorSetUpdates(ctx)
	require.NoError(t, err)
	require.Equal(t, 0, len(updates))

	// create a series of events that will bond and unbond the validator with
	// lowest power in a single block context (height)
	ctx = ctx.WithBlockHeight(2)

	validators[1], err = keeper.GetValidator(ctx, validators[1].ValAddress)
	require.NoError(t, err)

	keeper.DeleteValidatorByPowerIndex(ctx, validators[0])
	del, found = keeper.GetDelegation(ctx, decsdk.AccAddress(validators[0].ValAddress), validators[0].ValAddress)
	require.True(t, found)
	del.Coin.Amount = del.Coin.Amount.Sub(tokens)
	keeper.SetDelegation(ctx, del)
	keeper.SetValidator(ctx, validators[0])
	keeper.SetValidatorByPowerIndex(ctx, validators[0])
	updates, err = keeper.ApplyAndReturnValidatorSetUpdates(ctx)
	require.NoError(t, err)
	require.Equal(t, 0, len(updates))

	keeper.DeleteValidatorByPowerIndex(ctx, validators[1])
	tokens = sdk.TokensFromConsensusPower(250)
	del, found = keeper.GetDelegation(ctx, decsdk.AccAddress(validators[1].ValAddress), validators[1].ValAddress)
	require.True(t, found)
	del.Coin.Amount = del.Coin.Amount.Sub(tokens)
	keeper.SetDelegation(ctx, del)
	keeper.SetValidator(ctx, validators[1])
	keeper.SetValidatorByPowerIndex(ctx, validators[1])

	// verify initial Tendermint updates are correct
	updates, err = keeper.ApplyAndReturnValidatorSetUpdates(ctx)
	require.NoError(t, err)
	require.Equal(t, 1, len(updates))
	require.Equal(t, validators[1].ABCIValidatorUpdate(keeper.TotalStake(ctx, validators[1])), updates[0])

	updates, err = keeper.ApplyAndReturnValidatorSetUpdates(ctx)
	require.NoError(t, err)
	require.Equal(t, 0, len(updates))
}
