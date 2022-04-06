package validator

import (
	"bytes"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"bitbucket.org/decimalteam/go-node/config"
	"bitbucket.org/decimalteam/go-node/x/coin"
	"bitbucket.org/decimalteam/go-node/x/multisig"
	"bitbucket.org/decimalteam/go-node/x/nft"
	val "bitbucket.org/decimalteam/go-node/x/validator/internal/keeper"
	"bitbucket.org/decimalteam/go-node/x/validator/internal/types"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/supply"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"github.com/tendermint/tendermint/libs/log"
	tmtypes "github.com/tendermint/tendermint/types"
	dbm "github.com/tendermint/tm-db"
)

//______________________________________________________________________

func TestValidatorByPowerIndex(t *testing.T) {
	validatorAddr, validatorAddr3 := sdk.ValAddress(val.Addrs[0]), sdk.ValAddress(val.Addrs[1])

	initPower := int64(1000000)
	initBond := TokensFromConsensusPower(initPower)
	ctx, _, keeper, supplyKeeper, coinKeeper, _ := val.CreateTestInput(t, false, initPower)

	// create validator
	msgCreateValidator := NewTestMsgDeclareCandidate(validatorAddr, val.PKs[0], initBond)
	res, err := handleMsgDeclareCandidate(ctx, keeper, msgCreateValidator)
	require.NoError(t, err)
	require.NotNil(t, res)

	// must end-block
	updates, err := keeper.ApplyAndReturnValidatorSetUpdates(ctx)
	require.NoError(t, err)
	require.Equal(t, 1, len(updates))

	// verify the self-delegation exists
	bond, found := keeper.GetDelegation(ctx, sdk.AccAddress(validatorAddr), validatorAddr, keeper.BondDenom(ctx))
	require.True(t, found)
	gotBond := bond.Coin.Amount
	require.Equal(t, initBond, gotBond)

	// verify that the by power index exists
	validator, err := keeper.GetValidator(ctx, validatorAddr)
	require.NoError(t, err)
	power := types.GetValidatorsByPowerIndexKey(validator, keeper.TotalStake(ctx, validator))
	require.True(t, val.ValidatorByPowerIndexExists(ctx, keeper, power))

	// create a second validator keep it bonded
	msgCreateValidator = NewTestMsgDeclareCandidate(validatorAddr3, val.PKs[2], initBond)
	res, err = handleMsgDeclareCandidate(ctx, keeper, msgCreateValidator)
	require.NoError(t, err)
	require.NotNil(t, res)

	// must end-block
	updates, err = keeper.ApplyAndReturnValidatorSetUpdates(ctx)
	require.NoError(t, err)
	require.Equal(t, 1, len(updates))

	// slash and jail the first validator
	consAddr0 := sdk.ConsAddress(val.PKs[0].Address())
	keeper.Slash(ctx, consAddr0, 0, sdk.NewDecWithPrec(5, 1))
	keeper.Jail(ctx, consAddr0)
	_, err = keeper.ApplyAndReturnValidatorSetUpdates(ctx)
	require.NoError(t, err)

	validator, err = keeper.GetValidator(ctx, validatorAddr)
	if err != nil {
		panic(err)
	}
	require.Equal(t, types.Unbonded, validator.Status)                      // ensure is unbonding
	require.Equal(t, initBond.QuoRaw(2), keeper.TotalStake(ctx, validator)) // ensure tokens slashed
	keeper.Unjail(ctx, consAddr0)

	// the old power record should have been deleted as the power changed
	require.False(t, val.ValidatorByPowerIndexExists(ctx, keeper, power))

	// but the new power record should have been created
	validator, err = keeper.GetValidator(ctx, validatorAddr)
	if err != nil {
		panic(err)
	}
	power2 := types.GetValidatorsByPowerIndexKey(validator, validator.Tokens)
	require.True(t, val.ValidatorByPowerIndexExists(ctx, keeper, power2))

	// now the new record power index should be the same as the original record
	power3 := types.GetValidatorsByPowerIndexKey(validator, validator.Tokens)
	require.Equal(t, power2, power3)

	// unbond self-delegation
	unbondAmt := sdk.NewCoin(keeper.BondDenom(ctx), keeper.TotalStake(ctx, validator))
	msgUndelegate := types.NewMsgUnbond(validatorAddr, sdk.AccAddress(validatorAddr), unbondAmt)

	res, err = handleMsgUnbond(ctx, keeper, msgUndelegate)
	require.NoError(t, err)
	require.NotNil(t, res)

	var finishTime time.Time
	types.ModuleCdc.MustUnmarshalBinaryLengthPrefixed(res.Data, &finishTime)

	ctx = ctx.WithBlockTime(finishTime)
	EndBlocker(ctx, keeper, coinKeeper, supplyKeeper, false)

	// verify that by power key nolonger exists
	_, err = keeper.GetValidator(ctx, validatorAddr)
	require.Error(t, err)
	require.False(t, val.ValidatorByPowerIndexExists(ctx, keeper, power3))
}

func TestDuplicatesMsgCreateValidator(t *testing.T) {
	ctx, _, keeper, _, _, _ := val.CreateTestInput(t, false, 1000)

	addr1, addr2 := sdk.ValAddress(val.Addrs[0]), sdk.ValAddress(val.Addrs[1])
	pk1, pk2 := val.PKs[0], val.PKs[1]

	valTokens := TokensFromConsensusPower(10)
	msgCreateValidator1 := NewTestMsgDeclareCandidate(addr1, pk1, valTokens)
	res, err := handleMsgDeclareCandidate(ctx, keeper, msgCreateValidator1)
	require.NoError(t, err)
	require.NotNil(t, res)

	_, err = keeper.ApplyAndReturnValidatorSetUpdates(ctx)
	require.NoError(t, err)

	validator, err := keeper.GetValidator(ctx, addr1)
	require.NoError(t, err)
	assert.Equal(t, types.Bonded, validator.Status)
	assert.Equal(t, addr1, validator.ValAddress)
	assert.Equal(t, pk1, validator.PubKey)
	assert.Equal(t, valTokens, keeper.TotalStake(ctx, validator))

	// two validators can't have the same operator address
	msgCreateValidator2 := NewTestMsgDeclareCandidate(addr1, pk2, valTokens)
	res, err = handleMsgDeclareCandidate(ctx, keeper, msgCreateValidator2)
	require.Error(t, err)
	require.Nil(t, res)

	// two validators can't have the same pubkey
	msgCreateValidator3 := NewTestMsgDeclareCandidate(addr2, pk1, valTokens)
	res, err = handleMsgDeclareCandidate(ctx, keeper, msgCreateValidator3)
	require.Error(t, err)
	require.Nil(t, res)

	// must have different pubkey and operator
	msgCreateValidator4 := NewTestMsgDeclareCandidate(addr2, pk2, valTokens)
	res, err = handleMsgDeclareCandidate(ctx, keeper, msgCreateValidator4)
	require.NoError(t, err)
	require.NotNil(t, res)

	// must end-block
	updates, err := keeper.ApplyAndReturnValidatorSetUpdates(ctx)
	require.NoError(t, err)
	require.Equal(t, 1, len(updates))

	validator, err = keeper.GetValidator(ctx, addr2)
	require.NoError(t, err)

	assert.Equal(t, types.Bonded, validator.Status)
	assert.Equal(t, addr2, validator.ValAddress)
	assert.Equal(t, pk2, validator.PubKey)
	assert.True(sdk.IntEq(t, valTokens, keeper.TotalStake(ctx, validator)))
}

func TestInvalidPubKeyTypeMsgCreateValidator(t *testing.T) {
	ctx, _, keeper, _, _, _ := val.CreateTestInput(t, false, 1000)

	addr := sdk.ValAddress(val.Addrs[0])
	invalidPk := secp256k1.GenPrivKey().PubKey()

	// invalid pukKey type should not be allowed
	msgCreateValidator := NewTestMsgDeclareCandidate(addr, invalidPk, sdk.NewInt(10))
	res, err := handleMsgDeclareCandidate(ctx, keeper, msgCreateValidator)
	require.Error(t, err)
	require.Nil(t, res)

	ctx = ctx.WithConsensusParams(&abci.ConsensusParams{
		Validator: &abci.ValidatorParams{PubKeyTypes: []string{tmtypes.ABCIPubKeyTypeSecp256k1}},
	})

	res, err = handleMsgDeclareCandidate(ctx, keeper, msgCreateValidator)
	require.NoError(t, err)
	require.NotNil(t, res)
}

func TestIncrementsMsgDelegate(t *testing.T) {
	initPower := int64(1000)
	initBond := TokensFromConsensusPower(initPower)
	ctx, accMapper, keeper, _, _, _ := val.CreateTestInput(t, false, initPower)
	params := keeper.GetParams(ctx)

	bondAmount := TokensFromConsensusPower(10)
	validatorAddr, delegatorAddr := sdk.ValAddress(val.Addrs[0]), val.Addrs[1]

	// first create validator
	msgCreateValidator := NewTestMsgDeclareCandidate(validatorAddr, val.PKs[0], bondAmount)
	res, err := handleMsgDeclareCandidate(ctx, keeper, msgCreateValidator)
	require.NoError(t, err)
	require.NotNil(t, res)

	// apply TM updates
	_, err = keeper.ApplyAndReturnValidatorSetUpdates(ctx)
	require.NoError(t, err)

	validator, err := keeper.GetValidator(ctx, validatorAddr)
	require.NoError(t, err)
	require.Equal(t, types.Bonded, validator.Status)
	require.Equal(t, bondAmount, keeper.TotalStake(ctx, validator), "validator: %v", validator)

	_, found := keeper.GetDelegation(ctx, delegatorAddr, validatorAddr, keeper.BondDenom(ctx))
	require.False(t, found)

	bond, found := keeper.GetDelegation(ctx, sdk.AccAddress(validatorAddr), validatorAddr, keeper.BondDenom(ctx))
	require.True(t, found)
	require.Equal(t, bondAmount, bond.Coin.Amount)

	bondedTokens := keeper.TotalBondedTokens(ctx)
	require.Equal(t, bondAmount, bondedTokens)

	// just send the same msgbond multiple times
	msgDelegate := NewTestMsgDelegate(delegatorAddr, validatorAddr, bondAmount)

	for i := int64(0); i < 5; i++ {
		ctx = ctx.WithBlockHeight(i)

		res, err = handleMsgDelegate(ctx, keeper, msgDelegate)
		require.NoError(t, err)
		require.NotNil(t, res)

		//Check that the accounts and the bond account have the appropriate values
		validator, err := keeper.GetValidator(ctx, validatorAddr)
		require.NoError(t, err)
		bond, found := keeper.GetDelegation(ctx, delegatorAddr, validatorAddr, keeper.BondDenom(ctx))
		require.True(t, found)

		expBond := bondAmount.MulRaw(i + 1)
		expDelegatorShares := bondAmount.MulRaw(i + 2) // (1 self delegation)
		expDelegatorAcc := initBond.Sub(expBond)

		gotBond := bond.Coin.Amount
		gotDelegatorShares := keeper.TotalStake(ctx, validator)
		gotDelegatorAcc := accMapper.GetAccount(ctx, delegatorAddr).GetCoins().AmountOf(params.BondDenom)

		require.Equal(t, expBond, gotBond,
			"i: %v\nexpBond: %v\ngotBond: %v\nvalidator: %v\nbond: %v\n",
			i, expBond, gotBond, validator, bond)
		require.Equal(t, expDelegatorShares, gotDelegatorShares,
			"i: %v\nexpDelegatorShares: %v\ngotDelegatorShares: %v\nvalidator: %v\nbond: %v\n",
			i, expDelegatorShares, gotDelegatorShares, validator, bond)
		require.Equal(t, expDelegatorAcc, gotDelegatorAcc,
			"i: %v\nexpDelegatorAcc: %v\ngotDelegatorAcc: %v\nvalidator: %v\nbond: %v\n",
			i, expDelegatorAcc, gotDelegatorAcc, validator, bond)
	}
}

func TestIncrementsMsgUnbond(t *testing.T) {
	initPower := int64(1000)
	initBond := TokensFromConsensusPower(initPower)
	ctx, accMapper, keeper, supplyKeeper, coinKeeper, _ := val.CreateTestInput(t, false, initPower)

	params := keeper.GetParams(ctx)
	denom := params.BondDenom

	// create validator, delegate
	validatorAddr, delegatorAddr := sdk.ValAddress(val.Addrs[0]), val.Addrs[1]

	msgCreateValidator := NewTestMsgDeclareCandidate(validatorAddr, val.PKs[0], initBond)
	res, err := handleMsgDeclareCandidate(ctx, keeper, msgCreateValidator)
	require.NoError(t, err)
	require.NotNil(t, res)

	// initial balance
	amt1 := accMapper.GetAccount(ctx, delegatorAddr).GetCoins().AmountOf(denom)

	msgDelegate := NewTestMsgDelegate(delegatorAddr, validatorAddr, initBond)
	res, err = handleMsgDelegate(ctx, keeper, msgDelegate)
	require.NoError(t, err)
	require.NotNil(t, res)

	// balance should have been subtracted after delegation
	amt2 := accMapper.GetAccount(ctx, delegatorAddr).GetCoins().AmountOf(denom)
	require.True(sdk.IntEq(t, amt1.Sub(initBond), amt2))

	// apply TM updates
	_, err = keeper.ApplyAndReturnValidatorSetUpdates(ctx)
	require.NoError(t, err)

	validator, err := keeper.GetValidator(ctx, validatorAddr)
	require.NoError(t, err)
	require.Equal(t, initBond.MulRaw(2), validator.Tokens)

	// just send the same msgUnbond multiple times
	// TODO use decimals here
	unbondAmt := sdk.NewCoin(keeper.BondDenom(ctx), sdk.NewInt(10))
	msgUndelegate := types.NewMsgUnbond(validatorAddr, delegatorAddr, unbondAmt)
	numUnbonds := int64(5)

	for i := int64(0); i < numUnbonds; i++ {
		res, err := handleMsgUnbond(ctx, keeper, msgUndelegate)
		require.NoError(t, err)
		require.NotNil(t, res)

		var finishTime time.Time
		types.ModuleCdc.MustUnmarshalBinaryLengthPrefixed(res.Data, &finishTime)

		ctx = ctx.WithBlockTime(finishTime)
		EndBlocker(ctx, keeper, coinKeeper, supplyKeeper, false)

		// check that the accounts and the bond account have the appropriate values
		validator, err = keeper.GetValidator(ctx, validatorAddr)
		require.NoError(t, err)
		bond, found := keeper.GetDelegation(ctx, delegatorAddr, validatorAddr, keeper.BondDenom(ctx))
		require.True(t, found)

		expBond := initBond.Sub(unbondAmt.Amount.Mul(sdk.NewInt(i + 1)))
		expDelegatorShares := initBond.MulRaw(2).Sub(unbondAmt.Amount.Mul(sdk.NewInt(i + 1)))
		expDelegatorAcc := initBond.Sub(expBond)

		gotBond := bond.Coin.Amount
		gotDelegatorShares := keeper.TotalStake(ctx, validator)
		gotDelegatorAcc := accMapper.GetAccount(ctx, delegatorAddr).GetCoins().AmountOf(params.BondDenom)

		require.Equal(t, expBond, gotBond,
			"i: %v\nexpBond: %v\ngotBond: %v\nvalidator: %v\nbond: %v\n",
			i, expBond, gotBond, validator, bond)
		require.Equal(t, expDelegatorShares, gotDelegatorShares,
			"i: %v\nexpDelegatorShares: %v\ngotDelegatorShares: %v\nvalidator: %v\nbond: %v\n",
			i, expDelegatorShares, gotDelegatorShares, validator, bond)
		require.Equal(t, expDelegatorAcc, gotDelegatorAcc,
			"i: %v\nexpDelegatorAcc: %v\ngotDelegatorAcc: %v\nvalidator: %v\nbond: %v\n",
			i, expDelegatorAcc, gotDelegatorAcc, validator, bond)
	}

	// these are more than we have bonded now
	errorCases := []sdk.Int{
		//1<<64 - 1, // more than int64 power
		//1<<63 + 1, // more than int64 power
		TokensFromConsensusPower(1<<63 - 1),
		TokensFromConsensusPower(1 << 31),
		initBond,
	}

	for _, c := range errorCases {
		unbondAmt := sdk.NewCoin(keeper.BondDenom(ctx), c)
		msgUndelegate := types.NewMsgUnbond(validatorAddr, delegatorAddr, unbondAmt)
		res, err = handleMsgUnbond(ctx, keeper, msgUndelegate)
		require.Error(t, err)
		require.Nil(t, res)
	}

	leftBonded := initBond.Sub(unbondAmt.Amount.Mul(sdk.NewInt(numUnbonds)))

	// should be able to unbond remaining
	unbondAmt = sdk.NewCoin(keeper.BondDenom(ctx), leftBonded)
	msgUndelegate = types.NewMsgUnbond(validatorAddr, delegatorAddr, unbondAmt)
	res, err = handleMsgUnbond(ctx, keeper, msgUndelegate)
	require.NoError(t, err, "msgUnbond: %v\nshares: %s\nleftBonded: %s\n", msgUndelegate, unbondAmt, leftBonded)
	require.NotNil(t, res, "msgUnbond: %v\nshares: %s\nleftBonded: %s\n", msgUndelegate, unbondAmt, leftBonded)
}

func TestMultipleMsgCreateValidator(t *testing.T) {
	initPower := int64(1000)
	initTokens := TokensFromConsensusPower(initPower)
	ctx, accMapper, keeper, supplyKeeper, coinKeeper, _ := val.CreateTestInput(t, false, initPower)

	params := keeper.GetParams(ctx)
	blockTime := time.Now().UTC()
	ctx = ctx.WithBlockTime(blockTime)

	validatorAddrs := []sdk.ValAddress{
		sdk.ValAddress(val.Addrs[0]),
		sdk.ValAddress(val.Addrs[1]),
		sdk.ValAddress(val.Addrs[2]),
	}
	delegatorAddrs := []sdk.AccAddress{
		val.Addrs[0],
		val.Addrs[1],
		val.Addrs[2],
	}

	// bond them all
	for i, validatorAddr := range validatorAddrs {
		valTokens := TokensFromConsensusPower(10)
		msgCreateValidatorOnBehalfOf := NewTestMsgDeclareCandidate(validatorAddr, val.PKs[i], valTokens)

		res, err := handleMsgDeclareCandidate(ctx, keeper, msgCreateValidatorOnBehalfOf)
		require.NoError(t, err)
		require.NotNil(t, res)

		// verify that the account is bonded
		validators := keeper.GetValidators(ctx, 100)
		require.Equal(t, i+1, len(validators))

		validator := validators[i]
		balanceExpd := initTokens.Sub(valTokens)
		balanceGot := accMapper.GetAccount(ctx, delegatorAddrs[i]).GetCoins().AmountOf(params.BondDenom)

		require.Equal(t, i+1, len(validators), "expected %d validators res, err%d, validators: %v", i+1, len(validators), validators)
		require.Equal(t, valTokens, keeper.TotalStake(ctx, validator), "expected %d shares, res, err%d", 10, keeper.TotalStake(ctx, validator))
		require.Equal(t, balanceExpd, balanceGot, "expected account to have %d, res, err%d", balanceExpd, balanceGot)
	}

	EndBlocker(ctx, keeper, coinKeeper, supplyKeeper, false)

	// unbond them all by removing delegation
	for i, validatorAddr := range validatorAddrs {
		_, err := keeper.GetValidator(ctx, validatorAddr)
		require.NoError(t, err)

		unbondAmt := sdk.NewCoin(keeper.BondDenom(ctx), TokensFromConsensusPower(10))
		msgUndelegate := types.NewMsgUnbond(validatorAddr, sdk.AccAddress(validatorAddr), unbondAmt) // remove delegation
		res, err := handleMsgUnbond(ctx, keeper, msgUndelegate)
		require.NoError(t, err)
		require.NotNil(t, res)

		var finishTime time.Time
		types.ModuleCdc.MustUnmarshalBinaryLengthPrefixed(res.Data, &finishTime)

		// adds validator into unbonding queue
		EndBlocker(ctx, keeper, coinKeeper, supplyKeeper, false)

		// removes validator from queue and set
		EndBlocker(ctx.WithBlockTime(blockTime.Add(params.UnbondingTime)), keeper, coinKeeper, supplyKeeper, false)

		// Check that the validator is deleted from state
		validators := keeper.GetValidators(ctx, 100)
		require.Equal(t, len(validatorAddrs)-(i+1), len(validators),
			"expected %d validators got %d", len(validatorAddrs)-(i+1), len(validators))

		gotBalance := accMapper.GetAccount(ctx, delegatorAddrs[i]).GetCoins().AmountOf(params.BondDenom)
		require.True(t, initTokens.Equal(gotBalance))
	}
}

func TestMultipleMsgDelegate(t *testing.T) {
	ctx, _, keeper, supplyKeeper, coinKeeper, _ := val.CreateTestInput(t, false, 1000)
	validatorAddr, delegatorAddrs := sdk.ValAddress(val.Addrs[0]), val.Addrs[1:]

	// first make a validator
	msgCreateValidator := NewTestMsgDeclareCandidate(validatorAddr, val.PKs[0], sdk.NewInt(10))
	res, err := handleMsgDeclareCandidate(ctx, keeper, msgCreateValidator)
	require.NoError(t, err)
	require.NotNil(t, res)

	// delegate multiple parties
	for _, delegatorAddr := range delegatorAddrs {
		msgDelegate := NewTestMsgDelegate(delegatorAddr, validatorAddr, sdk.NewInt(10))
		res, err := handleMsgDelegate(ctx, keeper, msgDelegate)
		require.NoError(t, err)
		require.NotNil(t, res)

		// check that the account is bonded
		bond, found := keeper.GetDelegation(ctx, delegatorAddr, validatorAddr, keeper.BondDenom(ctx))
		require.True(t, found)
		require.NotNil(t, bond, "expected delegatee bond %d to exist", bond)
	}

	// unbond them all
	for _, delegatorAddr := range delegatorAddrs {
		unbondAmt := sdk.NewCoin(keeper.BondDenom(ctx), sdk.NewInt(10))
		msgUndelegate := types.NewMsgUnbond(validatorAddr, delegatorAddr, unbondAmt)

		res, err := handleMsgUnbond(ctx, keeper, msgUndelegate)
		require.NoError(t, err)
		require.NotNil(t, res)

		var finishTime time.Time
		types.ModuleCdc.MustUnmarshalBinaryLengthPrefixed(res.Data, &finishTime)

		ctx = ctx.WithBlockTime(finishTime)
		EndBlocker(ctx, keeper, coinKeeper, supplyKeeper, false)

		// check that the account is unbonded
		_, found := keeper.GetDelegation(ctx, delegatorAddr, validatorAddr, keeper.BondDenom(ctx))
		require.False(t, found)
	}
}

func TestValidatorQueue(t *testing.T) {
	ctx, _, keeper, supplyKeeper, coinKeeper, _ := val.CreateTestInput(t, false, 1000)
	validatorAddr := sdk.ValAddress(val.Addrs[0])

	// set the unbonding time
	params := keeper.GetParams(ctx)
	params.UnbondingTime = 7 * time.Second
	keeper.SetParams(ctx, params)

	// create the validator
	valTokens := TokensFromConsensusPower(10)
	msgCreateValidator := NewTestMsgDeclareCandidate(validatorAddr, val.PKs[0], valTokens)
	res, err := handleMsgDeclareCandidate(ctx, keeper, msgCreateValidator)
	require.NoError(t, err)
	require.NotNil(t, res)

	EndBlocker(ctx, keeper, coinKeeper, supplyKeeper, false)

	// unbond the all self-delegation to put validator in unbonding state
	unbondAmt := sdk.NewCoin(keeper.BondDenom(ctx), valTokens)
	msgUndelegateValidator := types.NewMsgUnbond(validatorAddr, sdk.AccAddress(validatorAddr), unbondAmt)
	res, err = handleMsgUnbond(ctx, keeper, msgUndelegateValidator)
	require.NoError(t, err)
	require.NotNil(t, res)

	var finishTime time.Time
	types.ModuleCdc.MustUnmarshalBinaryLengthPrefixed(res.Data, &finishTime)

	ctx = ctx.WithBlockTime(finishTime)
	EndBlocker(ctx, keeper, coinKeeper, supplyKeeper, false)

	var reqEvent sdk.Event
	for _, event := range ctx.EventManager().Events() {
		if event.Type == types.EventTypeCompleteUnbonding {
			reqEvent = event
			break
		}
	}
	require.Equal(t, reqEvent, sdk.NewEvent(
		types.EventTypeCompleteUnbonding,
		sdk.NewAttribute(types.AttributeKeyValidator, validatorAddr.String()),
		sdk.NewAttribute(types.AttributeKeyDelegator, sdk.AccAddress(validatorAddr).String()),
		sdk.NewAttribute(types.AttributeKeyCoin, unbondAmt.String()),
	))

	origHeader := ctx.BlockHeader()

	validator, err := keeper.GetValidator(ctx, validatorAddr)
	require.NoError(t, err)
	require.True(t, validator.IsUnbonding(), "%v", validator)

	// should still be unbonding at time 6 seconds later
	ctx = ctx.WithBlockTime(origHeader.Time.Add(time.Second * 6))
	EndBlocker(ctx, keeper, coinKeeper, supplyKeeper, false)

	validator, err = keeper.GetValidator(ctx, validatorAddr)
	require.NoError(t, err)
	require.True(t, validator.IsUnbonding(), "%v", validator)

	// should be in unbonded state at time 7 seconds later
	ctx = ctx.WithBlockTime(origHeader.Time.Add(time.Second * 7))
	EndBlocker(ctx, keeper, coinKeeper, supplyKeeper, false)

	validator, err = keeper.GetValidator(ctx, validatorAddr)
	require.Error(t, err)
}

func TestUnbondingPeriod(t *testing.T) {
	ctx, _, keeper, supplyKeeper, coinKeeper, _ := val.CreateTestInput(t, false, 1000)
	validatorAddr := sdk.ValAddress(val.Addrs[0])

	// set the unbonding time
	params := keeper.GetParams(ctx)
	params.UnbondingTime = 7 * time.Second
	keeper.SetParams(ctx, params)

	// create the validator
	valTokens := TokensFromConsensusPower(10)
	msgCreateValidator := NewTestMsgDeclareCandidate(validatorAddr, val.PKs[0], valTokens)
	res, err := handleMsgDeclareCandidate(ctx, keeper, msgCreateValidator)
	require.NoError(t, err)
	require.NotNil(t, res)

	EndBlocker(ctx, keeper, coinKeeper, supplyKeeper, false)

	// begin unbonding
	unbondAmt := sdk.NewCoin(keeper.BondDenom(ctx), TokensFromConsensusPower(10))
	msgUndelegate := types.NewMsgUnbond(validatorAddr, sdk.AccAddress(validatorAddr), unbondAmt)
	res, err = handleMsgUnbond(ctx, keeper, msgUndelegate)
	require.NoError(t, err)
	require.NotNil(t, res)

	origHeader := ctx.BlockHeader()

	_, found := keeper.GetUnbondingDelegation(ctx, sdk.AccAddress(validatorAddr), validatorAddr)
	require.True(t, found, "should not have unbonded")

	// cannot complete unbonding at same time
	EndBlocker(ctx, keeper, coinKeeper, supplyKeeper, false)
	_, found = keeper.GetUnbondingDelegation(ctx, sdk.AccAddress(validatorAddr), validatorAddr)
	require.True(t, found, "should not have unbonded")

	// cannot complete unbonding at time 6 seconds later
	ctx = ctx.WithBlockTime(origHeader.Time.Add(time.Second * 6))
	EndBlocker(ctx, keeper, coinKeeper, supplyKeeper, false)
	_, found = keeper.GetUnbondingDelegation(ctx, sdk.AccAddress(validatorAddr), validatorAddr)
	require.True(t, found, "should not have unbonded")

	// can complete unbonding at time 7 seconds later
	ctx = ctx.WithBlockTime(origHeader.Time.Add(time.Second * 7))
	EndBlocker(ctx, keeper, coinKeeper, supplyKeeper, false)
	_, found = keeper.GetUnbondingDelegation(ctx, sdk.AccAddress(validatorAddr), validatorAddr)
	require.False(t, found, "should have unbonded")
}

func TestUnbondingFromUnbondingValidator(t *testing.T) {
	ctx, _, keeper, supplyKeeper, coinKeeper, _ := val.CreateTestInput(t, false, 1000)
	validatorAddr, delegatorAddr := sdk.ValAddress(val.Addrs[0]), val.Addrs[1]

	// create the validator
	msgCreateValidator := NewTestMsgDeclareCandidate(validatorAddr, val.PKs[0], sdk.NewInt(10))
	res, err := handleMsgDeclareCandidate(ctx, keeper, msgCreateValidator)
	require.NoError(t, err)
	require.NotNil(t, res)

	// bond a delegator
	msgDelegate := NewTestMsgDelegate(delegatorAddr, validatorAddr, sdk.NewInt(10))
	res, err = handleMsgDelegate(ctx, keeper, msgDelegate)
	require.NoError(t, err)
	require.NotNil(t, res)

	// unbond the validators bond portion
	unbondAmt := sdk.NewCoin(keeper.BondDenom(ctx), sdk.NewInt(10))
	msgUndelegateValidator := types.NewMsgUnbond(validatorAddr, sdk.AccAddress(validatorAddr), unbondAmt)
	res, err = handleMsgUnbond(ctx, keeper, msgUndelegateValidator)
	require.NoError(t, err)
	require.NotNil(t, res)

	// change the ctx to Block Time one second before the validator would have unbonded
	var finishTime time.Time
	types.ModuleCdc.MustUnmarshalBinaryLengthPrefixed(res.Data, &finishTime)
	ctx = ctx.WithBlockTime(finishTime.Add(time.Second * -1))

	// unbond the delegator from the validator
	msgUndelegateDelegator := types.NewMsgUnbond(validatorAddr, delegatorAddr, unbondAmt)
	res, err = handleMsgUnbond(ctx, keeper, msgUndelegateDelegator)
	require.NoError(t, err)
	require.NotNil(t, res)

	ctx = ctx.WithBlockTime(ctx.BlockHeader().Time.Add(keeper.UnbondingTime(ctx)))

	// Run the EndBlocker
	EndBlocker(ctx, keeper, coinKeeper, supplyKeeper, false)

	// Check to make sure that the unbonding delegation is no longer in state
	// (meaning it was deleted in the above EndBlocker)
	_, found := keeper.GetUnbondingDelegation(ctx, delegatorAddr, validatorAddr)
	require.False(t, found, "should be removed from state")
}

func TestMultipleUnbondingDelegationAtSameTime(t *testing.T) {
	ctx, _, keeper, supplyKeeper, coinKeeper, _ := val.CreateTestInput(t, false, 1000)
	valAddr := sdk.ValAddress(val.Addrs[0])

	// set the unbonding time
	params := keeper.GetParams(ctx)
	params.UnbondingTime = 1 * time.Second
	keeper.SetParams(ctx, params)

	// create the validator
	valTokens := TokensFromConsensusPower(10)
	msgCreateValidator := NewTestMsgDeclareCandidate(valAddr, val.PKs[0], valTokens)
	res, err := handleMsgDeclareCandidate(ctx, keeper, msgCreateValidator)
	require.NoError(t, err)
	require.NotNil(t, res)

	// end block to bond
	EndBlocker(ctx, keeper, coinKeeper, supplyKeeper, false)

	// begin an unbonding delegation
	selfDelAddr := sdk.AccAddress(valAddr) // (the validator is it's own delegator)
	unbondAmt := sdk.NewCoin(keeper.BondDenom(ctx), valTokens.QuoRaw(2))
	msgUndelegate := types.NewMsgUnbond(valAddr, selfDelAddr, unbondAmt)
	res, err = handleMsgUnbond(ctx, keeper, msgUndelegate)
	require.NoError(t, err)
	require.NotNil(t, res)

	// there should only be one entry in the ubd object
	ubd, found := keeper.GetUnbondingDelegation(ctx, selfDelAddr, valAddr)
	require.True(t, found)
	require.Len(t, ubd.Entries, 1)

	// start a second ubd at this same time as the first
	res, err = handleMsgUnbond(ctx, keeper, msgUndelegate)
	require.NoError(t, err)
	require.NotNil(t, res)

	// now there should be two entries
	ubd, found = keeper.GetUnbondingDelegation(ctx, selfDelAddr, valAddr)
	require.True(t, found)
	require.Len(t, ubd.Entries, 2)

	// move forwaubd in time, should complete both ubds
	ctx = ctx.WithBlockTime(ctx.BlockHeader().Time.Add(2 * time.Second))
	EndBlocker(ctx, keeper, coinKeeper, supplyKeeper, false)

	ubd, found = keeper.GetUnbondingDelegation(ctx, selfDelAddr, valAddr)
	require.False(t, found)
}

func TestMultipleUnbondingDelegationAtUniqueTimes(t *testing.T) {
	ctx, _, keeper, supplyKeeper, coinKeeper, _ := val.CreateTestInput(t, false, 1000)
	valAddr := sdk.ValAddress(val.Addrs[0])

	// set the unbonding time
	params := keeper.GetParams(ctx)
	params.UnbondingTime = 10 * time.Second
	keeper.SetParams(ctx, params)

	// create the validator
	valTokens := TokensFromConsensusPower(10)
	msgCreateValidator := NewTestMsgDeclareCandidate(valAddr, val.PKs[0], valTokens)
	res, err := handleMsgDeclareCandidate(ctx, keeper, msgCreateValidator)
	require.NoError(t, err)
	require.NotNil(t, res)

	// end block to bond
	EndBlocker(ctx, keeper, coinKeeper, supplyKeeper, false)

	// begin an unbonding delegation
	selfDelAddr := sdk.AccAddress(valAddr) // (the validator is it's own delegator)
	unbondAmt := sdk.NewCoin(keeper.BondDenom(ctx), valTokens.QuoRaw(2))
	msgUndelegate := types.NewMsgUnbond(valAddr, selfDelAddr, unbondAmt)
	res, err = handleMsgUnbond(ctx, keeper, msgUndelegate)
	require.NoError(t, err)
	require.NotNil(t, res)

	// there should only be one entry in the ubd object
	ubd, found := keeper.GetUnbondingDelegation(ctx, selfDelAddr, valAddr)
	require.True(t, found)
	require.Len(t, ubd.Entries, 1)

	// move forwaubd in time and start a second redelegation
	ctx = ctx.WithBlockTime(ctx.BlockHeader().Time.Add(5 * time.Second))
	res, err = handleMsgUnbond(ctx, keeper, msgUndelegate)
	require.NoError(t, err)
	require.NotNil(t, res)

	// now there should be two entries
	ubd, found = keeper.GetUnbondingDelegation(ctx, selfDelAddr, valAddr)
	require.True(t, found)
	require.Len(t, ubd.Entries, 2)

	// move forwaubd in time, should complete the first redelegation, but not the second
	ctx = ctx.WithBlockTime(ctx.BlockHeader().Time.Add(5 * time.Second))
	EndBlocker(ctx, keeper, coinKeeper, supplyKeeper, false)
	ubd, found = keeper.GetUnbondingDelegation(ctx, selfDelAddr, valAddr)
	require.True(t, found)
	require.Len(t, ubd.Entries, 1)

	// move forwaubd in time, should complete the second redelegation
	ctx = ctx.WithBlockTime(ctx.BlockHeader().Time.Add(5 * time.Second))
	EndBlocker(ctx, keeper, coinKeeper, supplyKeeper, false)
	ubd, found = keeper.GetUnbondingDelegation(ctx, selfDelAddr, valAddr)
	require.False(t, found)
}

func TestUnbondingWhenExcessValidators(t *testing.T) {
	ctx, _, keeper, _, _, _ := val.CreateTestInput(t, false, 1000)
	validatorAddr1 := sdk.ValAddress(val.Addrs[0])
	validatorAddr2 := sdk.ValAddress(val.Addrs[1])
	validatorAddr3 := sdk.ValAddress(val.Addrs[2])

	// set the unbonding time
	params := keeper.GetParams(ctx)
	params.MaxValidators = 2
	keeper.SetParams(ctx, params)

	// add three validators
	valTokens1 := TokensFromConsensusPower(50)
	msgCreateValidator := NewTestMsgDeclareCandidate(validatorAddr1, val.PKs[0], valTokens1)
	res, err := handleMsgDeclareCandidate(ctx, keeper, msgCreateValidator)
	require.NoError(t, err)
	require.NotNil(t, res)
	// apply TM updates
	_, err = keeper.ApplyAndReturnValidatorSetUpdates(ctx)
	require.NoError(t, err)
	require.Equal(t, 1, len(keeper.GetLastValidators(ctx)))

	valTokens2 := TokensFromConsensusPower(30)
	msgCreateValidator = NewTestMsgDeclareCandidate(validatorAddr2, val.PKs[1], valTokens2)
	res, err = handleMsgDeclareCandidate(ctx, keeper, msgCreateValidator)
	require.NoError(t, err)
	require.NotNil(t, res)
	// apply TM updates
	_, err = keeper.ApplyAndReturnValidatorSetUpdates(ctx)
	require.NoError(t, err)
	require.Equal(t, 2, len(keeper.GetLastValidators(ctx)))

	valTokens3 := TokensFromConsensusPower(10)
	msgCreateValidator = NewTestMsgDeclareCandidate(validatorAddr3, val.PKs[2], valTokens3)
	res, err = handleMsgDeclareCandidate(ctx, keeper, msgCreateValidator)
	require.NoError(t, err)
	require.NotNil(t, res)
	// apply TM updates
	_, err = keeper.ApplyAndReturnValidatorSetUpdates(ctx)
	require.NoError(t, err)
	require.Equal(t, 2, len(keeper.GetLastValidators(ctx)))

	// unbond the validator-2
	unbondAmt := sdk.NewCoin(keeper.BondDenom(ctx), valTokens2)
	msgUndelegate := types.NewMsgUnbond(validatorAddr2, sdk.AccAddress(validatorAddr2), unbondAmt)
	res, err = handleMsgUnbond(ctx, keeper, msgUndelegate)
	require.NoError(t, err)
	require.NotNil(t, res)

	// apply TM updates
	_, err = keeper.ApplyAndReturnValidatorSetUpdates(ctx)
	require.NoError(t, err)

	// because there are extra validators waiting to get in, the queued
	// validator (aka. validator-1) should make it into the bonded group, thus
	// the total number of validators should stay the same
	vals := keeper.GetLastValidators(ctx)
	require.Equal(t, 2, len(vals), "vals %v", vals)
	val1, err := keeper.GetValidator(ctx, validatorAddr1)
	require.NoError(t, err)
	require.Equal(t, types.Bonded, val1.Status, "%v", val1)
}

func TestInvalidMsg(t *testing.T) {
	k := val.Keeper{}
	h := NewHandler(k)

	_, err := h(sdk.NewContext(nil, abci.Header{}, false, nil), sdk.NewTestMsg())
	require.Errorf(t, err, "unrecognized staking message type")
}

func TestEditCandidate(t *testing.T) {
	ctx, _, keeper, _, _, _ := val.CreateTestInput(t, false, 1000)
	validatorAddr1 := sdk.ValAddress(val.Addrs[0])
	rewardAddr1 := val.Addrs[1]

	valTokens1 := TokensFromConsensusPower(50)
	msgCreateValidator := NewTestMsgDeclareCandidate(validatorAddr1, val.PKs[0], valTokens1)
	res, err := handleMsgDeclareCandidate(ctx, keeper, msgCreateValidator)
	require.NoError(t, err)
	require.NotNil(t, res)

	_, err = keeper.ApplyAndReturnValidatorSetUpdates(ctx)
	require.NoError(t, err)
	require.Equal(t, 1, len(keeper.GetLastValidators(ctx)))

	msgEditValidator := NewMsgEditCandidate(validatorAddr1, rewardAddr1, types.Description{})
	res, err = handleMsgEditCandidate(ctx, keeper, msgEditValidator)
	require.NoError(t, err)
	require.NotNil(t, res)

	_, err = keeper.ApplyAndReturnValidatorSetUpdates(ctx)
	require.NoError(t, err)

	validator, err := keeper.GetValidator(ctx, validatorAddr1)
	require.NoError(t, err)

	validator.AccumRewards = valTokens1

	err = keeper.SetValidator(ctx, validator)
	require.NoError(t, err)

	err = keeper.PayRewards(ctx)
	require.NoError(t, err)
}

func TestSetOnline(t *testing.T) {
	ctx, _, keeper, _, _, _ := val.CreateTestInput(t, false, 1000)
	validatorAddr1 := sdk.ValAddress(val.Addrs[0])
	validatorAddr2 := sdk.ValAddress(val.Addrs[1])
	validatorAddr3 := sdk.ValAddress(val.Addrs[2])

	// add three validators
	valTokens1 := TokensFromConsensusPower(50)
	msgCreateValidator := NewTestMsgDeclareCandidate(validatorAddr1, val.PKs[0], valTokens1)
	res, err := handleMsgDeclareCandidate(ctx, keeper, msgCreateValidator)
	require.NoError(t, err)
	require.NotNil(t, res)
	// apply TM updates
	_, err = keeper.ApplyAndReturnValidatorSetUpdates(ctx)
	require.NoError(t, err)
	require.Equal(t, 1, len(keeper.GetLastValidators(ctx)))

	valTokens2 := TokensFromConsensusPower(30)
	msgCreateValidator = NewTestMsgDeclareCandidate(validatorAddr2, val.PKs[1], valTokens2)
	res, err = handleMsgDeclareCandidate(ctx, keeper, msgCreateValidator)
	require.NoError(t, err)
	require.NotNil(t, res)
	// apply TM updates
	_, err = keeper.ApplyAndReturnValidatorSetUpdates(ctx)
	require.NoError(t, err)
	require.Equal(t, 2, len(keeper.GetLastValidators(ctx)))

	valTokens3 := TokensFromConsensusPower(10)
	msgCreateValidator = NewTestMsgDeclareCandidate(validatorAddr3, val.PKs[2], valTokens3)
	res, err = handleMsgDeclareCandidate(ctx, keeper, msgCreateValidator)
	require.NoError(t, err)
	require.NotNil(t, res)
	// apply TM updates
	_, err = keeper.ApplyAndReturnValidatorSetUpdates(ctx)
	require.NoError(t, err)
	require.Equal(t, 3, len(keeper.GetLastValidators(ctx)))

	val3, err := keeper.GetValidator(ctx, validatorAddr3)
	require.NoError(t, err)

	val3.Online = false
	val3.Status = types.Unbonded

	err = keeper.SetValidator(ctx, val3)
	require.NoError(t, err)

	_, err = keeper.ApplyAndReturnValidatorSetUpdates(ctx)
	require.NoError(t, err)
	require.Equal(t, 3, len(keeper.GetLastValidators(ctx)))

	res, err = handleMsgSetOnline(ctx, keeper, NewMsgSetOnline(validatorAddr3))
	require.NoError(t, err)
	require.NotNil(t, res)

	updates, err := keeper.ApplyAndReturnValidatorSetUpdates(ctx)
	require.NoError(t, err)
	require.Equal(t, 1, len(updates))
	require.Equal(t, 3, len(keeper.GetLastValidators(ctx)))
}

func TestUnbondNFT(t *testing.T) {
	ctx, _, keeper, supplyKeeper, coinKeeper, nftKeeper := val.CreateTestInput(t, false, 1000)
	validatorAddr := sdk.ValAddress(val.Addrs[0])
	delegatorAddr := sdk.AccAddress(validatorAddr)

	// set the unbonding time
	params := keeper.GetParams(ctx)
	params.UnbondingTime = 7 * time.Second
	keeper.SetParams(ctx, params)

	// create the validator
	valTokens := types.TokensFromConsensusPower(10)
	msgCreateValidator := NewTestMsgDeclareCandidate(validatorAddr, val.PKs[0], valTokens)
	res, err := handleMsgDeclareCandidate(ctx, keeper, msgCreateValidator)
	require.NoError(t, err)
	require.NotNil(t, res)

	EndBlocker(ctx, keeper, coinKeeper, supplyKeeper, false)

	// create nft
	const denom = "denom1"
	const tokenID = "token1"
	quantity := sdk.NewInt(100)
	reserve := sdk.NewInt(100)

	nftHandler := nft.GenericHandler(nftKeeper)
	_, err = nftHandler(ctx, nft.NewMsgMintNFT(delegatorAddr, delegatorAddr, tokenID, denom, "", quantity, reserve, true))
	require.NoError(t, err)

	// delegate nft
	msgDelegateNft := types.NewMsgDelegateNFT(validatorAddr, delegatorAddr, tokenID, denom, []int64{1, 2, 3})
	res, err = handleMsgDelegateNFT(ctx, keeper, msgDelegateNft)
	require.NoError(t, err)
	require.NotNil(t, res)

	// unbond the half of delegations nft
	unbondQuantity := sdk.NewInt(2)
	msgUnbondNFT := types.NewMsgUnbondNFT(validatorAddr, delegatorAddr, tokenID, denom, []int64{1, 2})
	res, err = handleMsgUnbondNFT(ctx, keeper, msgUnbondNFT)
	require.NoError(t, err)
	require.NotNil(t, res)

	unbondingDelegation, ok := keeper.GetUnbondingDelegation(ctx, delegatorAddr, validatorAddr)
	require.True(t, ok)
	require.Equal(t, unbondQuantity.Mul(reserve), unbondingDelegation.Entries[0].GetBalance().Amount)

	//validator, err := keeper.GetValidator(ctx, validatorAddr)
	//require.NoError(t, err)
	//require.Equal(t, validator.Tokens, valTokens.Add(quantity.Sub(unbondQuantity).Mul(reserve)))

	var finishTime time.Time
	types.ModuleCdc.MustUnmarshalBinaryLengthPrefixed(res.Data, &finishTime)

	ctx = ctx.WithBlockTime(finishTime)
	EndBlocker(ctx, keeper, coinKeeper, supplyKeeper, false)

	var reqEvent sdk.Event
	for _, event := range ctx.EventManager().Events() {
		if event.Type == types.EventTypeCompleteUnbondingNFT {
			reqEvent = event
			break
		}
	}
	require.Equal(t, sdk.NewEvent(
		types.EventTypeCompleteUnbondingNFT,
		sdk.NewAttribute(types.AttributeKeyDenom, denom),
		sdk.NewAttribute(types.AttributeKeyID, tokenID),
		sdk.NewAttribute(types.AttributeKeySubTokenIDs, "1,2"),
		sdk.NewAttribute(types.AttributeKeyValidator, validatorAddr.String()),
		sdk.NewAttribute(types.AttributeKeyDelegator, delegatorAddr.String()),
		sdk.NewAttribute(types.AttributeKeyCoin, sdk.NewCoin(keeper.BondDenom(ctx), unbondQuantity.Mul(reserve)).String()),
	), reqEvent)
}

func TestConvertAddr(t *testing.T) {
	_config := sdk.GetConfig()
	_config.SetBech32PrefixForConsensusNode(config.DecimalPrefixConsAddr, config.DecimalPrefixConsPub)

	pubkey, err := sdk.GetPubKeyFromBech32(sdk.Bech32PubKeyTypeConsPub, "dxvalconspub1zcjduepqvgm3w28q3nx8vle4p578lau9zx749glqtyrmnqfeaufzqlhy3x5sgg8523")
	require.NoError(t, err)

	consAddr := sdk.GetConsAddress(pubkey)

	t.Log(consAddr)
	t.Log(consAddr.String())

	addr, err := sdk.ConsAddressFromHex("30B3070F0F4203A88A75EE36FABB2DE09F50625A")
	require.NoError(t, err)

	t.Log(addr.String())
}

/////////////////////////////////////
// Hogpodge of all sorts of input required for testing.
// `initPower` is converted to an amount of tokens.
// If `initPower` is 0, no addrs get created.
func CreateTestInput(t *testing.T, isCheckTx bool, initPower int64, addrList []sdk.AccAddress) (sdk.Context, auth.AccountKeeper, Keeper, supply.Keeper, coin.Keeper, nft.Keeper) {
	keyStaking := sdk.NewKVStoreKey(types.StoreKey)
	tkeyStaking := sdk.NewTransientStoreKey(types.TStoreKey)
	keyAcc := sdk.NewKVStoreKey(auth.StoreKey)
	keyParams := sdk.NewKVStoreKey(params.StoreKey)
	tkeyParams := sdk.NewTransientStoreKey(params.TStoreKey)
	keySupply := sdk.NewKVStoreKey(supply.StoreKey)
	keyCoin := sdk.NewKVStoreKey(coin.StoreKey)
	keyMultisig := sdk.NewKVStoreKey(multisig.StoreKey)

	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(tkeyStaking, sdk.StoreTypeTransient, nil)
	ms.MountStoreWithDB(keyStaking, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyAcc, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyParams, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(tkeyParams, sdk.StoreTypeTransient, db)
	ms.MountStoreWithDB(keySupply, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyCoin, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyMultisig, sdk.StoreTypeIAVL, db)
	err := ms.LoadLatestVersion()
	require.Nil(t, err)

	_config := sdk.GetConfig()
	_config.SetCoinType(60)
	_config.SetFullFundraiserPath("44'/60'/0'/0/0")
	_config.SetBech32PrefixForAccount(config.DecimalPrefixAccAddr, config.DecimalPrefixAccPub)
	_config.SetBech32PrefixForValidator(config.DecimalPrefixValAddr, config.DecimalPrefixValPub)
	_config.SetBech32PrefixForConsensusNode(config.DecimalPrefixConsAddr, config.DecimalPrefixConsPub)

	ctx := sdk.NewContext(ms, abci.Header{ChainID: "foochainid"}, isCheckTx, log.NewNopLogger())
	ctx = ctx.WithConsensusParams(
		&abci.ConsensusParams{
			Validator: &abci.ValidatorParams{
				PubKeyTypes: []string{tmtypes.ABCIPubKeyTypeEd25519},
			},
		},
	)
	cdc := val.MakeTestCodec()

	feeCollectorAcc := supply.NewEmptyModuleAccount(auth.FeeCollectorName, supply.Burner, supply.Minter)
	notBondedPool := supply.NewEmptyModuleAccount(types.NotBondedPoolName, supply.Burner, supply.Staking)
	bondPool := supply.NewEmptyModuleAccount(types.BondedPoolName, supply.Burner, supply.Staking)

	blacklistedAddrs := make(map[string]bool)
	blacklistedAddrs[feeCollectorAcc.String()] = true
	blacklistedAddrs[notBondedPool.String()] = true
	blacklistedAddrs[bondPool.String()] = true

	pk := params.NewKeeper(cdc, keyParams, tkeyParams)

	accountKeeper := auth.NewAccountKeeper(
		cdc,    // amino codec
		keyAcc, // target store
		pk.Subspace(auth.DefaultParamspace),
		auth.ProtoBaseAccount, // prototype
	)

	bk := bank.NewBaseKeeper(
		accountKeeper,
		pk.Subspace(bank.DefaultParamspace),
		blacklistedAddrs,
	)

	maccPerms := map[string][]string{
		auth.FeeCollectorName:   nil,
		types.NotBondedPoolName: {supply.Burner, supply.Staking},
		types.BondedPoolName:    {supply.Burner, supply.Staking},
		nft.ReservedPool:        {supply.Burner},
	}
	supplyKeeper := supply.NewKeeper(cdc, keySupply, accountKeeper, bk, maccPerms)

	initTokens := types.TokensFromConsensusPower(initPower)
	initCoins := sdk.NewCoins(sdk.NewCoin(types.DefaultBondDenom, initTokens))
	totalSupply := sdk.NewCoins(sdk.NewCoin(types.DefaultBondDenom, initTokens.MulRaw(int64(len(addrList)))))

	supplyKeeper.SetSupply(ctx, supply.NewSupply(totalSupply))

	coinKeeper := coin.NewKeeper(cdc, keyCoin, pk.Subspace(coin.DefaultParamspace), accountKeeper, bk, config.GetDefaultConfig(config.ChainID))

	coinConfig := config.GetDefaultConfig(config.ChainID)
	coinKeeper.SetCoin(ctx, coin.Coin{
		Title:  coinConfig.TitleBaseCoin,
		Symbol: coinConfig.SymbolBaseCoin,
		Volume: coinConfig.InitialVolumeBaseCoin,
	})

	multisigKeeper := multisig.NewKeeper(cdc, keyMultisig, pk.Subspace(multisig.DefaultParamspace), accountKeeper, coinKeeper, bk)

	nftKeeper := nft.NewKeeper(cdc, keyCoin, supplyKeeper, types.DefaultBondDenom)

	keeper := NewKeeper(cdc, keyStaking, pk.Subspace(val.DefaultParamspace), coinKeeper, accountKeeper, supplyKeeper, multisigKeeper, nftKeeper, auth.FeeCollectorName)
	keeper.SetParams(ctx, types.DefaultParams())

	// set module accounts
	err = notBondedPool.SetCoins(totalSupply)
	require.NoError(t, err)

	supplyKeeper.SetModuleAccount(ctx, feeCollectorAcc)
	supplyKeeper.SetModuleAccount(ctx, bondPool)
	supplyKeeper.SetModuleAccount(ctx, notBondedPool)

	// fill all the addresses with some coins, set the loose pool tokens simultaneously
	for _, addr := range addrList {
		_, err := bk.AddCoins(ctx, addr, initCoins)
		if err != nil {
			panic(err)
		}
	}

	return ctx, accountKeeper, keeper, supplyKeeper, coinKeeper, nftKeeper
}

// default test create no more 999 addresses
func createTestAddr(numAddrs int) []sdk.AccAddress {
	var addresses []sdk.AccAddress
	var buffer bytes.Buffer
	var source = "0123456789ABCDEF"
	for i := 0; i < numAddrs; i++ {
		for j := 0; j < 40; j++ {
			buffer.WriteByte(source[rand.Intn(len(source))])
		}
		res, _ := sdk.AccAddressFromHex(buffer.String())
		bech := res.String()
		addresses = append(addresses, val.TestAddr(buffer.String(), bech))
		buffer.Reset()
	}
	return addresses
}

func TestMaximumSlots(t *testing.T) {
	const N = 4
	const AddrCount = 2000
	// init
	delegators := createTestAddr(AddrCount)
	ctx, _, keeper, _, _, _ := CreateTestInput(t, false, 20000, delegators)

	// params
	// set the unbonding time
	params := keeper.GetParams(ctx)
	params.UnbondingTime = 1 * time.Second
	keeper.SetParams(ctx, params)

	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)
	// 1. declare validators
	// first N address is validators
	for i := 0; i < N; i++ {
		validatorAddr := sdk.ValAddress(delegators[i])
		//delegatorAddr := sdk.AccAddress(validatorAddr)
		// create the validator
		valTokens := types.TokensFromConsensusPower(10000)
		msgCreateValidator := NewTestMsgDeclareCandidate(validatorAddr, val.PKs[i], valTokens)
		res, err := handleMsgDeclareCandidate(ctx, keeper, msgCreateValidator)
		require.NoError(t, err)
		require.NotNil(t, res)
	}
	_, err := keeper.ApplyAndReturnValidatorSetUpdates(ctx)
	require.NoError(t, err)
	require.Equal(t, N, len(keeper.GetLastValidators(ctx)))

	//delegate
	validatorAddr := sdk.ValAddress(delegators[0])
	validator, _ := keeper.GetValidator(ctx, validatorAddr)
	for i := 0; i < 1000; i++ {
		delegatorAddr := delegators[N+i]
		_, err := keeper.Delegate(ctx, delegatorAddr, sdk.NewCoin(DefaultBondDenom, TokensFromConsensusPower(10)), types.Bonded, validator, false)
		if err != nil {
			fmt.Printf("delegate err=%s\n", err.Error())
		}
	}
	//bond
	delegatorAddr := delegators[N+1002]
	msgDelegate := NewTestMsgDelegate(delegatorAddr, validatorAddr, TokensFromConsensusPower(20))
	handleMsgDelegate(ctx, keeper, msgDelegate)
	//unbond
	unbondAmt := sdk.NewCoin(keeper.BondDenom(ctx), TokensFromConsensusPower(5))
	msgUndelegate := types.NewMsgUnbond(validatorAddr, delegators[N+1], unbondAmt)
	handleMsgUnbond(ctx, keeper, msgUndelegate)
	// check validators count
	keeper.ApplyAndReturnValidatorSetUpdates(ctx)
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)
	updates, err := keeper.ApplyAndReturnValidatorSetUpdates(ctx)
	require.NoError(t, err)
	//check updates for duplicates
	updateDups := make(map[string]int)
	for _, u := range updates {
		updateDups[u.PubKey.String()] = updateDups[u.PubKey.String()] + 1
	}
	for _, v := range updateDups {
		require.Equal(t, 1, v, "duplicate in updates")
	}
	//
	require.Equal(t, N, len(keeper.GetLastValidators(ctx)))
}
