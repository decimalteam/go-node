package keeper

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"

	decsdk "bitbucket.org/decimalteam/go-node/utils/types"
	"bitbucket.org/decimalteam/go-node/x/validator/internal/types"
)

// tests GetDelegation, GetDelegatorDelegations, SetDelegation, RemoveDelegation, GetDelegatorDelegations
func TestDelegation(t *testing.T) {
	ctx, _, keeper, _, _ := CreateTestInput(t, false, 10)

	//construct the validators
	amts := []sdk.Int{sdk.NewInt(9), sdk.NewInt(8), sdk.NewInt(7)}
	var validators [3]types.Validator
	for i, amt := range amts {
		validators[i] = types.NewValidator(addrVals[i], PKs[i], sdk.ZeroDec(), decsdk.AccAddress(addrVals[i]), types.Description{})
		delegation := types.NewDelegation(decsdk.AccAddress(addrVals[i]), addrVals[i], sdk.NewCoin(keeper.BondDenom(ctx), amt))
		keeper.SetDelegation(ctx, delegation)
	}

	validators[0] = TestingUpdateValidator(keeper, ctx, validators[0], true)
	validators[1] = TestingUpdateValidator(keeper, ctx, validators[1], true)
	validators[2] = TestingUpdateValidator(keeper, ctx, validators[2], true)

	// first add a validators[0] to delegate too

	bond1to1 := types.NewDelegation(addrDels[0], addrVals[0], sdk.NewCoin(keeper.BondDenom(ctx), sdk.NewInt(9)))

	// check the empty keeper first
	_, found := keeper.GetDelegation(ctx, addrDels[0], addrVals[0])
	require.False(t, found)

	// set and retrieve a record
	keeper.SetDelegation(ctx, bond1to1)
	resBond, found := keeper.GetDelegation(ctx, addrDels[0], addrVals[0])
	require.True(t, found)
	require.True(t, bond1to1.Equal(resBond))

	// modify a records, save, and retrieve
	bond1to1.Coin = sdk.NewCoin(keeper.BondDenom(ctx), sdk.NewInt(99))
	keeper.SetDelegation(ctx, bond1to1)
	resBond, found = keeper.GetDelegation(ctx, addrDels[0], addrVals[0])
	require.True(t, found)
	require.True(t, bond1to1.Equal(resBond))

	// add some more records
	bond1to2 := types.NewDelegation(addrDels[0], addrVals[1], sdk.NewCoin(keeper.BondDenom(ctx), sdk.NewInt(9)))
	bond1to3 := types.NewDelegation(addrDels[0], addrVals[2], sdk.NewCoin(keeper.BondDenom(ctx), sdk.NewInt(9)))
	bond2to1 := types.NewDelegation(addrDels[1], addrVals[0], sdk.NewCoin(keeper.BondDenom(ctx), sdk.NewInt(9)))
	bond2to2 := types.NewDelegation(addrDels[1], addrVals[1], sdk.NewCoin(keeper.BondDenom(ctx), sdk.NewInt(9)))
	bond2to3 := types.NewDelegation(addrDels[1], addrVals[2], sdk.NewCoin(keeper.BondDenom(ctx), sdk.NewInt(9)))
	keeper.SetDelegation(ctx, bond1to2)
	keeper.SetDelegation(ctx, bond1to3)
	keeper.SetDelegation(ctx, bond2to1)
	keeper.SetDelegation(ctx, bond2to2)
	keeper.SetDelegation(ctx, bond2to3)

	// test all bond retrieve capabilities
	resBonds := keeper.GetDelegatorDelegations(ctx, addrDels[0], 5)
	require.Equal(t, 3, len(resBonds))
	require.True(t, bond1to1.Equal(resBonds[0]))
	require.True(t, bond1to2.Equal(resBonds[1]))
	require.True(t, bond1to3.Equal(resBonds[2]))
	resBonds = keeper.GetAllDelegatorDelegations(ctx, addrDels[0])
	require.Equal(t, 3, len(resBonds))
	resBonds = keeper.GetDelegatorDelegations(ctx, addrDels[0], 2)
	require.Equal(t, 2, len(resBonds))
	resBonds = keeper.GetDelegatorDelegations(ctx, addrDels[1], 5)
	require.Equal(t, 3, len(resBonds))
	require.True(t, bond2to1.Equal(resBonds[0]))
	require.True(t, bond2to2.Equal(resBonds[1]))
	require.True(t, bond2to3.Equal(resBonds[2]))
	allBonds := keeper.GetAllDelegations(ctx)
	require.Equal(t, 9, len(allBonds))
	require.True(t, bond1to1.Equal(allBonds[0]))
	require.True(t, bond1to2.Equal(allBonds[1]))
	require.True(t, bond1to3.Equal(allBonds[2]))
	require.True(t, bond2to1.Equal(allBonds[3]))
	require.True(t, bond2to2.Equal(allBonds[4]))
	require.True(t, bond2to3.Equal(allBonds[5]))

	resVals := keeper.GetDelegatorValidators(ctx, addrDels[0], 3)
	require.Equal(t, 3, len(resVals))
	resVals = keeper.GetDelegatorValidators(ctx, addrDels[1], 4)
	require.Equal(t, 3, len(resVals))

	for i := 0; i < 3; i++ {

		resVal, err := keeper.GetDelegatorValidator(ctx, addrDels[0], addrVals[i])
		require.Nil(t, err)
		require.Equal(t, addrVals[i], resVal.ValAddress)

		resVal, err = keeper.GetDelegatorValidator(ctx, addrDels[1], addrVals[i])
		require.Nil(t, err)
		require.Equal(t, addrVals[i], resVal.ValAddress)

		resDels := keeper.GetValidatorDelegations(ctx, addrVals[i])
		require.Len(t, resDels, 3)
	}

	// delete a record
	keeper.RemoveDelegation(ctx, bond2to3)
	_, found = keeper.GetDelegation(ctx, addrDels[1], addrVals[2])
	require.False(t, found)
	resBonds = keeper.GetDelegatorDelegations(ctx, addrDels[1], 5)
	require.Equal(t, 2, len(resBonds))
	require.True(t, bond2to1.Equal(resBonds[0]))
	require.True(t, bond2to2.Equal(resBonds[1]))

	resBonds = keeper.GetAllDelegatorDelegations(ctx, addrDels[1])
	require.Equal(t, 2, len(resBonds))

	// delete all the records from delegator 2
	keeper.RemoveDelegation(ctx, bond2to1)
	keeper.RemoveDelegation(ctx, bond2to2)
	_, found = keeper.GetDelegation(ctx, addrDels[1], addrVals[0])
	require.False(t, found)
	_, found = keeper.GetDelegation(ctx, addrDels[1], addrVals[1])
	require.False(t, found)
	resBonds = keeper.GetDelegatorDelegations(ctx, addrDels[1], 5)
	require.Equal(t, 0, len(resBonds))
}

// tests Get/Set/Remove UnbondingDelegation
func TestUnbondingDelegation(t *testing.T) {
	ctx, _, keeper, _, _ := CreateTestInput(t, false, 0)

	ubd := types.NewUnbondingDelegation(addrDels[0], addrVals[0], 0,
		time.Unix(0, 0), sdk.NewCoin(keeper.BondDenom(ctx), sdk.NewInt(5)))

	// set and retrieve a record
	keeper.SetUnbondingDelegation(ctx, ubd)
	resUnbond, found := keeper.GetUnbondingDelegation(ctx, addrDels[0], addrVals[0])
	require.True(t, found)
	require.True(t, ubd.Equal(resUnbond))

	// modify a records, save, and retrieve
	ubd.Entries[0].Balance = sdk.NewCoin(keeper.BondDenom(ctx), sdk.NewInt(21))
	keeper.SetUnbondingDelegation(ctx, ubd)

	resUnbonds := keeper.GetUnbondingDelegations(ctx, addrDels[0], 5)
	require.Equal(t, 1, len(resUnbonds))

	resUnbonds = keeper.GetAllUnbondingDelegations(ctx, addrDels[0])
	require.Equal(t, 1, len(resUnbonds))

	resUnbond, found = keeper.GetUnbondingDelegation(ctx, addrDels[0], addrVals[0])
	require.True(t, found)
	require.True(t, ubd.Equal(resUnbond))

	// delete a record
	keeper.RemoveUnbondingDelegation(ctx, ubd)
	_, found = keeper.GetUnbondingDelegation(ctx, addrDels[0], addrVals[0])
	require.False(t, found)

	resUnbonds = keeper.GetUnbondingDelegations(ctx, addrDels[0], 5)
	require.Equal(t, 0, len(resUnbonds))

	resUnbonds = keeper.GetAllUnbondingDelegations(ctx, addrDels[0])
	require.Equal(t, 0, len(resUnbonds))

}

func TestUnbondDelegation(t *testing.T) {
	ctx, _, keeper, _, _ := CreateTestInput(t, false, 0)

	startTokens := sdk.TokensFromConsensusPower(10)

	notBondedPool := keeper.GetNotBondedPool(ctx)
	err := notBondedPool.SetCoins(sdk.NewCoins(sdk.NewCoin(keeper.BondDenom(ctx), startTokens)))
	require.NoError(t, err)
	keeper.supplyKeeper.SetModuleAccount(ctx, notBondedPool)

	// create a validator and a delegator to that validator
	// note this validator starts not-bonded
	validator := types.NewValidator(addrVals[0], PKs[0], sdk.ZeroDec(), decsdk.AccAddress(addrVals[0]), types.Description{})

	validator = TestingUpdateValidator(keeper, ctx, validator, true)

	delegation := types.NewDelegation(addrDels[0], addrVals[0], sdk.NewCoin(keeper.BondDenom(ctx), startTokens))
	keeper.SetDelegation(ctx, delegation)

	bondTokens := sdk.TokensFromConsensusPower(6)
	err = keeper.unbond(ctx, addrDels[0], addrVals[0], sdk.NewCoin(keeper.BondDenom(ctx), bondTokens))
	require.NoError(t, err)

	delegation, found := keeper.GetDelegation(ctx, addrDels[0], addrVals[0])
	require.True(t, found)
	validator, err = keeper.GetValidator(ctx, addrVals[0])
	require.NoError(t, err)

	remainingTokens := startTokens.Sub(bondTokens)
	require.Equal(t, remainingTokens, delegation.Coin.Amount)
	require.Equal(t, remainingTokens, keeper.TotalStake(ctx, validator))
}

func TestUndelegateFromUnbondedValidator(t *testing.T) {
	ctx, _, keeper, _, _ := CreateTestInput(t, false, 1)
	delTokens := sdk.TokensFromConsensusPower(10)
	delCoins := sdk.NewCoins(sdk.NewCoin(keeper.BondDenom(ctx), delTokens))

	// add bonded tokens to pool for delegations
	notBondedPool := keeper.GetNotBondedPool(ctx)
	err := notBondedPool.SetCoins(notBondedPool.GetCoins().Add(delCoins...))
	require.NoError(t, err)
	keeper.supplyKeeper.SetModuleAccount(ctx, notBondedPool)

	// create a validator with a self-delegation
	validator := types.NewValidator(addrVals[0], PKs[0], sdk.ZeroDec(), decsdk.AccAddress(addrVals[0]), types.Description{})

	valTokens := sdk.TokensFromConsensusPower(10)
	delegation := types.NewDelegation(decsdk.AccAddress(addrVals[0]), addrVals[0], sdk.NewCoin(keeper.BondDenom(ctx), valTokens))
	keeper.SetDelegation(ctx, delegation)
	validator = TestingUpdateValidator(keeper, ctx, validator, true)
	require.True(t, validator.IsBonded())

	bondedPool := keeper.GetBondedPool(ctx)
	err = bondedPool.SetCoins(bondedPool.GetCoins().Add(delCoins...))
	require.NoError(t, err)
	keeper.supplyKeeper.SetModuleAccount(ctx, bondedPool)

	// create a second delegation to this validator
	keeper.DeleteValidatorByPowerIndex(ctx, validator)
	delegation = types.NewDelegation(addrDels[0], addrVals[0], sdk.NewCoin(keeper.BondDenom(ctx), delTokens))
	keeper.SetDelegation(ctx, delegation)
	validator = TestingUpdateValidator(keeper, ctx, validator, true)
	require.True(t, validator.IsBonded())
	delegation = types.NewDelegation(addrDels[0], addrVals[0], sdk.NewCoin(keeper.BondDenom(ctx), delTokens))
	keeper.SetDelegation(ctx, delegation)

	ctx = ctx.WithBlockHeight(10)
	ctx = ctx.WithBlockTime(time.Unix(333, 0))

	validator, err = keeper.GetValidator(ctx, addrVals[0])
	require.NoError(t, err)
	validator = validator.UpdateStatus(types.Unbonded)
	err = keeper.SetValidator(ctx, validator)
	require.NoError(t, err)

	// unbond some of the other delegation's shares
	unbondTokens := sdk.TokensFromConsensusPower(6)
	_, err = keeper.Undelegate(ctx, addrDels[0], addrVals[0], sdk.NewCoin(keeper.BondDenom(ctx), unbondTokens))
	require.NoError(t, err)

	// unbond rest of the other delegation's shares
	remainingTokens := delTokens.Sub(unbondTokens)
	_, err = keeper.Undelegate(ctx, addrDels[0], addrVals[0], sdk.NewCoin(keeper.BondDenom(ctx), remainingTokens))
	require.NoError(t, err)
}

func TestUnbondingAllDelegationFromValidator(t *testing.T) {
	ctx, _, keeper, _, _ := CreateTestInput(t, false, 0)
	delTokens := sdk.TokensFromConsensusPower(10)
	delCoins := sdk.NewCoins(sdk.NewCoin(keeper.BondDenom(ctx), delTokens))

	// add bonded tokens to pool for delegations
	notBondedPool := keeper.GetNotBondedPool(ctx)
	err := notBondedPool.SetCoins(notBondedPool.GetCoins().Add(delCoins...))
	require.NoError(t, err)
	keeper.supplyKeeper.SetModuleAccount(ctx, notBondedPool)

	//create a validator with a self-delegation
	validator := types.NewValidator(addrVals[0], PKs[0], sdk.ZeroDec(), decsdk.AccAddress(addrVals[0]), types.Description{})

	valTokens := sdk.TokensFromConsensusPower(10)
	val0AccAddr := decsdk.AccAddress(addrVals[0].Bytes())

	selfDelegation := types.NewDelegation(val0AccAddr, addrVals[0], sdk.NewCoin(keeper.BondDenom(ctx), valTokens))
	keeper.SetDelegation(ctx, selfDelegation)

	// create a second delegation to this validator
	keeper.DeleteValidatorByPowerIndex(ctx, validator)
	delegation := types.NewDelegation(val0AccAddr, addrVals[0], sdk.NewCoin(keeper.BondDenom(ctx), delTokens))
	keeper.SetDelegation(ctx, delegation)

	bondedPool := keeper.GetBondedPool(ctx)
	err = bondedPool.SetCoins(bondedPool.GetCoins().Add(delCoins...))
	require.NoError(t, err)
	keeper.supplyKeeper.SetModuleAccount(ctx, bondedPool)

	validator = TestingUpdateValidator(keeper, ctx, validator, true)
	require.True(t, validator.IsBonded())

	delegation = types.NewDelegation(val0AccAddr, addrVals[0], sdk.NewCoin(keeper.BondDenom(ctx), delTokens))
	keeper.SetDelegation(ctx, delegation)

	ctx = ctx.WithBlockHeight(10)
	ctx = ctx.WithBlockTime(time.Unix(333, 0))

	// unbond all the remaining delegation
	_, err = keeper.Undelegate(ctx, val0AccAddr, addrVals[0], sdk.NewCoin(keeper.BondDenom(ctx), delTokens))
	require.NoError(t, err)

	validator = TestingUpdateValidator(keeper, ctx, validator, true)

	// validator should still be in state and still be in unbonding state
	validator, err = keeper.GetValidator(ctx, addrVals[0])
	require.NoError(t, err)
	require.Equal(t, validator.Status, types.Unbonded)
}
