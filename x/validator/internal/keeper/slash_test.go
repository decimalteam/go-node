package keeper

import (
	"log"
	"testing"
	"time"

	ncfg "bitbucket.org/decimalteam/go-node/config"

	"bitbucket.org/decimalteam/go-node/utils/helpers"
	"bitbucket.org/decimalteam/go-node/x/validator/internal/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
)

// TODO integrate with test_common.go helper (CreateTestInput)
// setup helper function - creates two validators
func setupHelper(t *testing.T, power int64) (sdk.Context, Keeper, types.Params) {
	// setup
	ctx, _, keeper, _, _, _ := CreateTestInput(t, false, power)

	params := keeper.GetParams(ctx)
	numVals := int64(3)
	amt := types.TokensFromConsensusPower(power)

	bondedCoins := sdk.NewCoins(sdk.NewCoin(keeper.BondDenom(ctx), amt.MulRaw(numVals)))

	bondedPool := keeper.GetBondedPool(ctx)
	err := bondedPool.SetCoins(bondedCoins)
	require.NoError(t, err)
	keeper.supplyKeeper.SetModuleAccount(ctx, bondedPool)

	// add numVals validators
	for i := int64(0); i < numVals; i++ {
		validator := types.NewValidator(addrVals[i], PKs[i], sdk.ZeroDec(), sdk.AccAddress(addrVals[i]), types.Description{})
		validator.Online = true
		del := types.NewDelegation(sdk.AccAddress(validator.ValAddress), validator.ValAddress, sdk.NewCoin(keeper.BondDenom(ctx), amt))
		keeper.SetDelegation(ctx, del)
		validator = TestingUpdateValidator(keeper, ctx, validator, true)
		keeper.SetValidatorByConsAddr(ctx, validator)
	}

	return ctx, keeper, params
}

//_________________________________________________________________________________

// tests Jail, Unjail
func TestRevocation(t *testing.T) {

	// setup
	ctx, keeper, _ := setupHelper(t, 10)
	addr := addrVals[0]
	consAddr := sdk.ConsAddress(PKs[0].Address())

	// initial state
	val, err := keeper.GetValidator(ctx, addr)
	require.NoError(t, err)
	require.False(t, val.IsJailed())

	// test jail
	keeper.Jail(ctx, consAddr)
	val, err = keeper.GetValidator(ctx, addr)
	require.NoError(t, err)
	require.True(t, val.IsJailed())

	// test unjail
	keeper.Unjail(ctx, consAddr)
	val, err = keeper.GetValidator(ctx, addr)
	require.NoError(t, err)
	require.False(t, val.IsJailed())
}

func TestSlashBondedDelegationNFT(t *testing.T) {
	ctx, _, keeper, _, _, nftKeeper := CreateTestInput(t, false, 1000)

	valAddr := addrVals[0]
	delAddr := sdk.AccAddress(addrVals[0])
	amt := types.TokensFromConsensusPower(100)

	// NFT params
	const denom = "denom1"
	const tokenID = "token1"
	quantity := sdk.NewInt(5)
	reserve := helpers.BipToPip(sdk.NewInt(100))

	// create nft
	_, err := nftKeeper.MintNFT(ctx,
		denom,
		tokenID,
		reserve,
		quantity,
		delAddr,
		delAddr,
		"",
		true,
	)
	require.NoError(t, err)

	bondedCoins := sdk.NewCoins(sdk.NewCoin(keeper.BondDenom(ctx), amt))

	bondedPool := keeper.GetBondedPool(ctx)
	err = bondedPool.SetCoins(bondedCoins)
	require.NoError(t, err)
	keeper.supplyKeeper.SetModuleAccount(ctx, bondedPool)

	// set validator
	validator := types.NewValidator(addrVals[0], PKs[0], sdk.ZeroDec(), delAddr, types.Description{})
	validator.Online = true
	del := types.NewDelegation(delAddr, validator.ValAddress, sdk.NewCoin(keeper.BondDenom(ctx), amt))
	keeper.SetDelegation(ctx, del)
	validator = TestingUpdateValidator(keeper, ctx, validator, true)
	keeper.SetValidatorByConsAddr(ctx, validator)

	// set nft delegation
	delegationNFT := types.NewDelegationNFT(delAddr, valAddr, tokenID, denom, []int64{1, 2, 3},
		sdk.NewCoin(keeper.BondDenom(ctx), quantity.Mul(reserve)))
	keeper.SetDelegationNFT(ctx, delegationNFT)

	keeper.Slash(ctx, validator.GetConsAddr(), ctx.BlockHeight(), types.SlashFractionDowntime)

	delegationNFT, ok := keeper.GetDelegationNFT(ctx, valAddr, delAddr, tokenID, denom)
	require.True(t, ok)
	require.Equal(t, []int64{1, 2, 3}, delegationNFT.SubTokenIDs)
}

// tests slashUnbondingDelegation
func TestSlashUnbondingDelegation(t *testing.T) {
	ctx, keeper, _ := setupHelper(t, 10)
	fraction := sdk.NewDecWithPrec(5, 1)

	// set an unbonding delegation with expiration timestamp (beyond which the
	// unbonding delegation shouldn't be slashed)
	ubd := types.NewUnbondingDelegation(addrDels[0], addrVals[0], types.NewUnbondingDelegationEntry(0,
		time.Unix(5, 0), sdk.NewCoin(keeper.BondDenom(ctx), sdk.NewInt(10))))

	keeper.SetUnbondingDelegation(ctx, ubd)

	// unbonding started prior to the infraction height, stake didn't contribute
	slashAmount := keeper.slashUnbondingDelegation(ctx, ubd, 1, fraction)
	require.Equal(t, int64(0), slashAmount.AmountOf(keeper.BondDenom(ctx)).Int64())

	// after the expiration time, no longer eligible for slashing
	ctx = ctx.WithBlockHeader(abci.Header{Time: time.Unix(10, 0)})
	keeper.SetUnbondingDelegation(ctx, ubd)
	slashAmount = keeper.slashUnbondingDelegation(ctx, ubd, 0, fraction)
	require.Equal(t, int64(0), slashAmount.AmountOf(keeper.BondDenom(ctx)).Int64())

	// test valid slash, before expiration timestamp and to which stake contributed
	oldUnbondedPool := keeper.GetNotBondedPool(ctx)
	ctx = ctx.WithBlockHeader(abci.Header{Time: time.Unix(0, 0)})
	keeper.SetUnbondingDelegation(ctx, ubd)
	slashAmount = keeper.slashUnbondingDelegation(ctx, ubd, 0, fraction)
	require.Equal(t, int64(5), slashAmount.AmountOf(keeper.BondDenom(ctx)).Int64())
	ubd, found := keeper.GetUnbondingDelegation(ctx, addrDels[0], addrVals[0])
	require.True(t, found)
	require.Len(t, ubd.Entries, 1)

	// initial balance unchanged
	require.Equal(t, sdk.NewCoin(keeper.BondDenom(ctx), sdk.NewInt(10)), ubd.Entries[0].GetInitialBalance())

	// balance decreased
	require.Equal(t, sdk.NewCoin(keeper.BondDenom(ctx), sdk.NewInt(5)), ubd.Entries[0].GetBalance())
	newUnbondedPool := keeper.GetNotBondedPool(ctx)
	diffTokens := oldUnbondedPool.GetCoins().Sub(newUnbondedPool.GetCoins()).AmountOf(keeper.BondDenom(ctx))
	require.Equal(t, int64(5), diffTokens.Int64())
}

// tests Slash at a future height (must panic)
func TestSlashAtFutureHeight(t *testing.T) {
	ctx, keeper, _ := setupHelper(t, 10)
	consAddr := sdk.ConsAddress(PKs[0].Address())
	fraction := sdk.NewDecWithPrec(5, 1)
	require.Panics(t, func() { keeper.Slash(ctx, consAddr, 1, fraction) })
}

// test slash at a negative height
// this just represents pre-genesis and should have the same effect as slashing at height 0
func TestSlashAtNegativeHeight(t *testing.T) {
	ctx, keeper, _ := setupHelper(t, 10)
	consAddr := sdk.ConsAddress(PKs[0].Address())
	fraction := sdk.NewDecWithPrec(5, 1)

	oldBondedPool := keeper.GetBondedPool(ctx)
	log.Println(oldBondedPool)
	validator, err := keeper.GetValidatorByConsAddr(ctx, consAddr)
	require.NoError(t, err)
	keeper.Slash(ctx, consAddr, -2, fraction)

	// read updated state
	validator, err = keeper.GetValidatorByConsAddr(ctx, consAddr)
	require.NoError(t, err)
	newBondedPool := keeper.GetBondedPool(ctx)
	log.Println(newBondedPool)

	// end block
	updates, err := keeper.ApplyAndReturnValidatorSetUpdates(ctx)
	require.NoError(t, err)
	require.Equal(t, 1, len(updates), "cons addr: %v, updates: %v", []byte(consAddr), updates)

	validator, err = keeper.GetValidator(ctx, validator.ValAddress)
	require.NoError(t, err)
	// power decreased
	require.Equal(t, int64(5), validator.ConsensusPower())
	// pool bonded shares decreased
	diffTokens := oldBondedPool.GetCoins().Sub(newBondedPool.GetCoins()).AmountOf(keeper.BondDenom(ctx))
	log.Println(diffTokens)
	require.Equal(t, types.TokensFromConsensusPower(5), diffTokens)
}

//tests Slash at the current height
func TestSlashValidatorAtCurrentHeight(t *testing.T) {
	ctx, keeper, _ := setupHelper(t, 10)
	consAddr := sdk.ConsAddress(PKs[0].Address())
	fraction := sdk.NewDecWithPrec(5, 1)

	oldBondedPool := keeper.GetBondedPool(ctx)
	validator, err := keeper.GetValidatorByConsAddr(ctx, consAddr)
	require.NoError(t, err)
	keeper.Slash(ctx, consAddr, ctx.BlockHeight(), fraction)

	// read updated state
	validator, err = keeper.GetValidatorByConsAddr(ctx, consAddr)
	require.NoError(t, err)
	newBondedPool := keeper.GetBondedPool(ctx)
	log.Println(newBondedPool.GetCoins())

	// end block
	updates, err := keeper.ApplyAndReturnValidatorSetUpdates(ctx)
	require.Equal(t, 1, len(updates), "cons addr: %v, updates: %v", []byte(consAddr), updates)

	validator, err = keeper.GetValidator(ctx, validator.ValAddress)
	require.NoError(t, err)
	// power decreased
	require.Equal(t, int64(5), validator.ConsensusPower())
	// pool bonded shares decreased
	diffTokens := oldBondedPool.GetCoins().Sub(newBondedPool.GetCoins()).AmountOf(keeper.BondDenom(ctx))
	require.Equal(t, types.TokensFromConsensusPower(5), diffTokens)
}

// tests Slash at a previous height with an unbonding delegation
func TestSlashWithUnbondingDelegation(t *testing.T) {
	ctx, keeper, _ := setupHelper(t, 10)
	consAddr := sdk.ConsAddress(PKs[0].Address())
	fraction := sdk.NewDecWithPrec(5, 1)

	// set an unbonding delegation with expiration timestamp beyond which the
	// unbonding delegation shouldn't be slashed
	ubdTokens := sdk.NewCoin(keeper.BondDenom(ctx), types.TokensFromConsensusPower(4))
	ubd := types.NewUnbondingDelegation(addrDels[0], addrVals[0], types.NewUnbondingDelegationEntry(11,
		time.Unix(0, 0), ubdTokens))
	keeper.SetUnbondingDelegation(ctx, ubd)

	// slash validator for the first time
	ctx = ctx.WithBlockHeight(4)
	oldBondedPool := keeper.GetBondedPool(ctx)
	validator, err := keeper.GetValidatorByConsAddr(ctx, consAddr)
	require.NoError(t, err)
	keeper.Slash(ctx, consAddr, 3, fraction)

	// end block
	updates, err := keeper.ApplyAndReturnValidatorSetUpdates(ctx)
	require.NoError(t, err)
	require.Equal(t, 1, len(updates))

	// read updating unbonding delegation
	ubd, found := keeper.GetUnbondingDelegation(ctx, addrDels[0], addrVals[0])
	require.True(t, found)
	require.Len(t, ubd.Entries, 1)
	// balance decreased
	require.Equal(t, sdk.NewCoin(keeper.BondDenom(ctx), types.TokensFromConsensusPower(2)), ubd.Entries[0].GetBalance())
	// read updated pool
	newBondedPool := keeper.GetBondedPool(ctx)
	// bonded tokens burned
	diffTokens := oldBondedPool.GetCoins().Sub(newBondedPool.GetCoins()).AmountOf(keeper.BondDenom(ctx))
	require.Equal(t, types.TokensFromConsensusPower(5), diffTokens)
	// read updated validator
	validator, err = keeper.GetValidatorByConsAddr(ctx, consAddr)
	require.NoError(t, err)
	// power decreased by 3 - 6 stake originally bonded at the time of infraction
	// was still bonded at the time of discovery and was slashed by half, 4 stake
	// bonded at the time of discovery hadn't been bonded at the time of infraction
	// and wasn't slashed
	require.Equal(t, int64(5), validator.ConsensusPower())

	// slash validator again
	ctx = ctx.WithBlockHeight(6)
	keeper.Slash(ctx, consAddr, 5, fraction)
	ubd, found = keeper.GetUnbondingDelegation(ctx, addrDels[0], addrVals[0])
	require.True(t, found)
	require.Len(t, ubd.Entries, 1)
	// balance decreased again
	require.Equal(t, sdk.NewCoin(keeper.BondDenom(ctx), sdk.ZeroInt()), ubd.Entries[0].GetBalance())
	// read updated pool
	newBondedPool = keeper.GetBondedPool(ctx)
	// bonded tokens burned again
	diffTokens = oldBondedPool.GetCoins().Sub(newBondedPool.GetCoins()).AmountOf(keeper.BondDenom(ctx))
	require.Equal(t, types.TokensFromConsensusPower(5).QuoRaw(2).Add(types.TokensFromConsensusPower(5)), diffTokens)
	// read updated validator
	validator, err = keeper.GetValidatorByConsAddr(ctx, consAddr)
	require.NoError(t, err)
	// power decreased by 3 again
	require.Equal(t, int64(2), validator.ConsensusPower())

	// slash validator again
	// all originally bonded stake has been slashed, so this will have no effect
	// on the unbonding delegation, but it will slash stake bonded since the infraction
	// this may not be the desirable behaviour, ref https://github.com/cosmos/cosmos-sdk/issues/1440
	ctx = ctx.WithBlockHeight(6)
	keeper.Slash(ctx, consAddr, 5, fraction)
	ubd, found = keeper.GetUnbondingDelegation(ctx, addrDels[0], addrVals[0])
	require.True(t, found)
	require.Len(t, ubd.Entries, 1)
	// balance unchanged
	require.Equal(t, sdk.NewCoin(keeper.BondDenom(ctx), sdk.NewInt(0)), ubd.Entries[0].GetBalance())
	// read updated pool
	newBondedPool = keeper.GetBondedPool(ctx)
	// bonded tokens burned again
	diffTokens = oldBondedPool.GetCoins().Sub(newBondedPool.GetCoins()).AmountOf(keeper.BondDenom(ctx))
	require.Equal(t, types.TokensFromConsensusPower(5).QuoRaw(4).Add(types.TokensFromConsensusPower(5).QuoRaw(2)).Add(types.TokensFromConsensusPower(5)), diffTokens)
	// read updated validator
	validator, err = keeper.GetValidatorByConsAddr(ctx, consAddr)
	require.NoError(t, err)
	// power decreased by 3 again
	require.Equal(t, int64(1), validator.ConsensusPower())

	// slash validator again
	// all originally bonded stake has been slashed, so this will have no effect
	// on the unbonding delegation, but it will slash stake bonded since the infraction
	// this may not be the desirable behaviour, ref https://github.com/cosmos/cosmos-sdk/issues/1440
	ctx = ctx.WithBlockHeight(6)
	keeper.Slash(ctx, consAddr, 5, fraction)
	ubd, found = keeper.GetUnbondingDelegation(ctx, addrDels[0], addrVals[0])
	require.True(t, found)
	require.Len(t, ubd.Entries, 1)
	// balance unchanged
	require.Equal(t, sdk.NewCoin(keeper.BondDenom(ctx), sdk.NewInt(0)), ubd.Entries[0].GetBalance())
	// read updated pool
	newBondedPool = keeper.GetBondedPool(ctx)
	// just 1 bonded token burned again since that's all the validator now has
	diffTokens = oldBondedPool.GetCoins().Sub(newBondedPool.GetCoins()).AmountOf(keeper.BondDenom(ctx))
	require.Equal(t, types.TokensFromConsensusPower(5).QuoRaw(8).Add(types.TokensFromConsensusPower(5).QuoRaw(4)).Add(types.TokensFromConsensusPower(5).QuoRaw(2)).Add(types.TokensFromConsensusPower(5)), diffTokens)
	// apply TM updates
	_, err = keeper.ApplyAndReturnValidatorSetUpdates(ctx)
	require.NoError(t, err)
	// read updated validator
	// power decreased by 1 again, validator is out of stake
	// validator should be in unbonding period
	validator, _ = keeper.GetValidatorByConsAddr(ctx, consAddr)
	require.Equal(t, validator.GetStatus(), types.Unbonded)
}

func TestInGracePeriod(t *testing.T) {
	ctxWithHeight := func(height int64) sdk.Context {
		ctx := sdk.Context{}
		return ctx.WithBlockHeight(height)
	}

	//test overlapping grace periods
	{
		updatesInfo := ncfg.NewUpdatesInfo("")
		p0start := int64(10000)
		p0end := p0start + ncfg.GracePeriod
		p1start := p0end - ncfg.GracePeriod/2
		p1end := p1start + ncfg.GracePeriod
		updatesInfo.AddExecutedPlan("0", p0start)
		updatesInfo.PushNewPlanHeight(p1start)
		//
		require.False(t, inGracePeriod(ctxWithHeight(p0start-1), updatesInfo))
		require.False(t, inGracePeriod(ctxWithHeight(p1end+1), updatesInfo))
		require.True(t, inGracePeriod(ctxWithHeight(p0start+1), updatesInfo))
		require.True(t, inGracePeriod(ctxWithHeight(p0end+1), updatesInfo))
		require.True(t, inGracePeriod(ctxWithHeight(p1start-1), updatesInfo))
		require.True(t, inGracePeriod(ctxWithHeight(p1end-1), updatesInfo))
	}

	//test non-overlapping grace periods
	{
		updatesInfo := ncfg.NewUpdatesInfo("")
		p0start := int64(10000)
		p0end := p0start + ncfg.GracePeriod
		p1start := p0end + ncfg.GracePeriod/2
		p1end := p1start + ncfg.GracePeriod
		updatesInfo.AddExecutedPlan("0", p0start)
		updatesInfo.AddExecutedPlan("1", p1start)
		updatesInfo.PushNewPlanHeight(p1start)
		//
		require.False(t, inGracePeriod(ctxWithHeight(p0start-1), updatesInfo))
		require.False(t, inGracePeriod(ctxWithHeight(p0end+1), updatesInfo))
		require.True(t, inGracePeriod(ctxWithHeight(p0start+1), updatesInfo))
		require.True(t, inGracePeriod(ctxWithHeight(p0end-1), updatesInfo))
		require.False(t, inGracePeriod(ctxWithHeight(p1start-1), updatesInfo))
		require.False(t, inGracePeriod(ctxWithHeight(p1end+1), updatesInfo))
		require.True(t, inGracePeriod(ctxWithHeight(p1start+1), updatesInfo))
		require.True(t, inGracePeriod(ctxWithHeight(p1end-1), updatesInfo))
	}
}

func sumArray(ctx sdk.Context, k Keeper, consAddr sdk.ConsAddress) int64 {
	var res int64
	for i := int64(0); i < types.SignedBlocksWindow; i++ {
		v := k.getValidatorMissedBlockBitArray(ctx, consAddr, i)
		if v {
			res++
		}
	}
	return res
}

func TestHandleValidatorSignatureMissCount(t *testing.T) {
	ctx, keeper, _ := setupHelper(t, 10)
	validator := keeper.GetValidators(ctx, 100)[0]
	adr := validator.PubKey.Address()
	keeper.addPubkey(ctx, validator.PubKey)
	consAddr := validator.GetConsAddr()

	ncfg.UpdatesInfo.AllBlocks = make(map[string]int64)
	ncfg.UpdatesInfo.LastBlock = -10000000000

	//use different patterns for signed block and grace period
	//to trying make mess for MissedBlockArray and MissedBlocksCounter
	testsuite := []struct {
		patternSign  []bool
		patternGrace []bool
	}{
		{[]bool{true, true, false}, []bool{true, false}},
		{[]bool{true, true, true, false}, []bool{true, false, false}},
		{[]bool{true, true, true, true, false}, []bool{true, true, true, false}},
	}
	for _, pattern := range testsuite {
		//prepare sequences
		sequenceSign := make([]bool, 0)
		sequenceGrace := make([]bool, 0)
		for i := int64(0); i < 100; i++ {
			sequenceSign = append(sequenceSign, pattern.patternSign...)
		}
		for i := int64(0); i < 100; i++ {
			sequenceGrace = append(sequenceSign, pattern.patternGrace...)
		}
		count := len(sequenceGrace)
		if count > len(sequenceSign) {
			count = len(sequenceSign)
		}
		//reset
		validator, _ = keeper.GetValidatorByConsAddr(ctx, consAddr)
		validator.Online = false
		keeper.SetValidator(ctx, validator)
		keeper.HandleValidatorSignature(ctx, adr, 10, false)
		validator.Online = true
		keeper.SetValidator(ctx, validator)

		for i := 0; i < count; i++ {
			ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)
			sign := sequenceSign[i]
			ncfg.UpdatesInfo.LastBlock = -10000000000 //grace period = false
			if sequenceGrace[i] {
				ncfg.UpdatesInfo.LastBlock = ctx.BlockHeight() //grace period = true
			}
			keeper.HandleValidatorSignature(ctx, adr, 10, sign)
			signInfo, _ := keeper.getValidatorSigningInfo(ctx, consAddr)
			// CHECK: right missed block counter
			// old code will not pass
			require.Equal(t, signInfo.MissedBlocksCounter, sumArray(ctx, keeper, consAddr))
			//
			validator, _ = keeper.GetValidatorByConsAddr(ctx, consAddr)
			if validator.IsJailed() {
				keeper.Unjail(ctx, consAddr)
			}
		}
	}
}

func TestHandleValidatorSignatureJail(t *testing.T) {
	for _, inGrace := range []bool{false, true} {
		var oldBond, newBond sdk.Int
		ctx, keeper, _ := setupHelper(t, 10)
		val := keeper.GetValidators(ctx, 100)[0]
		adr := val.PubKey.Address()
		keeper.addPubkey(ctx, val.PubKey)
		consAddr := val.GetConsAddr()

		// get old delegation
		for _, del := range keeper.GetAllDelegations(ctx) {
			if del.GetValidatorAddr().String() == val.ValAddress.String() {
				oldBond = del.GetCoin().Amount
			}
		}

		if inGrace {
			ncfg.UpdatesInfo.AllBlocks = make(map[string]int64)
			ncfg.UpdatesInfo.LastBlock = 0
		} else {
			ncfg.UpdatesInfo.AllBlocks = make(map[string]int64)
			ncfg.UpdatesInfo.LastBlock = -ncfg.GracePeriod - 1000
		}

		maxMissed := types.SignedBlocksWindow - types.MinSignedPerWindow
		//not jail
		for i := int64(0); i < maxMissed; i++ {
			ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)
			keeper.HandleValidatorSignature(ctx, adr, 10, false)
			keeper.ApplyAndReturnValidatorSetUpdates(ctx)
		}
		val, _ = keeper.GetValidatorByConsAddr(ctx, consAddr)
		require.False(t, val.IsJailed(), "validator must be free")

		//reset
		val, _ = keeper.GetValidatorByConsAddr(ctx, consAddr)
		val.Online = false
		keeper.SetValidator(ctx, val)
		keeper.HandleValidatorSignature(ctx, adr, 10, false)
		val.Online = true
		keeper.SetValidator(ctx, val)

		//jail
		for i := int64(0); i < (maxMissed + 1); i++ {
			ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)
			keeper.HandleValidatorSignature(ctx, adr, 10, false)
			keeper.ApplyAndReturnValidatorSetUpdates(ctx)
		}
		val, _ = keeper.GetValidatorByConsAddr(ctx, consAddr)
		require.True(t, val.IsJailed(), "validator must be jailed")

		_, err := keeper.ApplyAndReturnValidatorSetUpdates(ctx)
		require.NoError(t, err)

		// get new delegation
		for _, del := range keeper.GetAllDelegations(ctx) {
			if del.GetValidatorAddr().String() == val.ValAddress.String() {
				newBond = del.GetCoin().Amount
			}
		}

		if inGrace {
			require.True(t, oldBond.Equal(newBond), "must no slash in grace period: %s ? %s", oldBond, newBond)
		} else {
			require.True(t, oldBond.GT(newBond), "must slash without grace period: %s ? %s", oldBond, newBond)
		}
	}

}
