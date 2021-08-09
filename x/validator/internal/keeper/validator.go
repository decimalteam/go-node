package keeper

import (
	"bitbucket.org/decimalteam/go-node/utils/updates"
	"bytes"
	"errors"
	"fmt"
	"sync"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"bitbucket.org/decimalteam/go-node/utils/formulas"
	"bitbucket.org/decimalteam/go-node/x/validator/internal/types"
)

// get a single validator
func (k Keeper) GetValidator(ctx sdk.Context, addr sdk.ValAddress) (types.Validator, error) {
	store := ctx.KVStore(k.storeKey)
	value := store.Get(types.GetValidatorKey(addr))
	if value == nil {
		return types.Validator{}, errors.New("not found validator ")
	}

	// amino bytes weren't found in cache, so amino unmarshal and add it to the cache
	validator, err := types.UnmarshalValidator(k.cdc, value)
	if err != nil {
		return types.Validator{}, errors.New("error unmarshal validator ")
	}

	return validator, nil
}

// get the set of all validators with no limits, used during genesis dump
func (k Keeper) GetAllValidators(ctx sdk.Context) (validators []types.Validator) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, []byte{types.ValidatorsKey})
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		validator, err := types.UnmarshalValidator(k.cdc, iterator.Value())
		if err != nil {
			panic(err)
		}
		validators = append(validators, validator)
	}
	return validators
}

func (k Keeper) HasValidator(ctx sdk.Context, addr sdk.ValAddress) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.GetValidatorKey(addr))
}

func (k Keeper) SetValidator(ctx sdk.Context, validator types.Validator) error {
	return k.set(ctx, types.GetValidatorKey(validator.ValAddress), validator)
}

// validator index
func (k Keeper) SetValidatorByConsAddr(ctx sdk.Context, validator types.Validator) {
	store := ctx.KVStore(k.storeKey)
	consAddr := sdk.GetConsAddress(validator.PubKey)
	store.Set(types.GetValidatorByConsAddrKey(consAddr), validator.ValAddress)
}

// get a single validator by consensus address
func (k Keeper) GetValidatorByConsAddr(ctx sdk.Context, consAddr sdk.ConsAddress) (types.Validator, error) {
	store := ctx.KVStore(k.storeKey)
	valAddr := store.Get(types.GetValidatorByConsAddrKey(consAddr))
	if valAddr == nil {
		return types.Validator{}, errors.New("not found validator ")
	}
	return k.GetValidator(ctx, valAddr)
}

// validator index
func (k Keeper) SetValidatorByPowerIndex(ctx sdk.Context, validator types.Validator) {
	// jailed validators are not kept in the power index
	if validator.Jailed {
		return
	}
	power := k.TotalStake(ctx, validator)
	validator.Tokens = power
	err := k.SetValidator(ctx, validator)
	if err != nil {
		panic(err)
	}
	ctx.KVStore(k.storeKey).Set(types.GetValidatorsByPowerIndexKey(validator, power), validator.ValAddress)
}

// validator index
func (k Keeper) SetValidatorByPowerIndexWithoutCalc(ctx sdk.Context, validator types.Validator) {
	// jailed validators are not kept in the power index
	if validator.Jailed {
		return
	}
	ctx.KVStore(k.storeKey).Set(types.GetValidatorsByPowerIndexKey(validator, validator.Tokens), validator.ValAddress)
}

// validator index
func (k Keeper) SetNewValidatorByPowerIndex(ctx sdk.Context, validator types.Validator) {
	ctx.KVStore(k.storeKey).Set(types.GetValidatorsByPowerIndexKey(validator, k.TotalStake(ctx, validator)), validator.ValAddress)
}

func (k Keeper) GetAllValidatorsByPowerIndex(ctx sdk.Context) []types.Validator {
	var validators []types.Validator
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, []byte{types.ValidatorsByPowerIndexKey})
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		validator, err := k.GetValidator(ctx, iterator.Value())
		if err != nil {
			panic(err)
		}
		validators = append(validators, validator)
	}
	return validators
}

func (k Keeper) GetAllValidatorsByPowerIndexReversed(ctx sdk.Context) []types.Validator {
	var validators []types.Validator
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStoreReversePrefixIterator(store, []byte{types.ValidatorsByPowerIndexKey})
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		validator, err := k.GetValidator(ctx, iterator.Value())
		if err != nil {
			panic(err)
		}
		validators = append(validators, validator)
	}
	return validators
}

func (k Keeper) TotalStake(ctx sdk.Context, validator types.Validator) sdk.Int {
	total := sdk.ZeroInt()
	delegations := k.GetValidatorDelegations(ctx, validator.ValAddress)
	mutex := sync.Mutex{}
	eventMutex := sync.Mutex{}
	wg := sync.WaitGroup{}
	wg.Add(len(delegations))
	for _, del := range delegations {
		if del.GetCoin().Denom != k.BondDenom(ctx) {
			coin, err := k.GetCoin(ctx, del.GetCoin().Denom)
			if err != nil {
				panic(err)
			}
			if ctx.BlockHeight() >= updates.Update2Block {
				fmt.Println(coin.String())
				delegatedCoin := k.GetDelegatedCoin(ctx, del.GetCoin().Denom)
				totalAmountCoin := formulas.CalculateSaleReturn(coin.Volume, coin.Reserve, coin.CRR, delegatedCoin)
				del = del.SetTokensBase(totalAmountCoin.Mul(del.GetCoin().Amount.ToDec().Quo(delegatedCoin.ToDec()).TruncateInt()))
			} else {
				del = del.SetTokensBase(formulas.CalculateSaleReturn(coin.Volume, coin.Reserve, coin.CRR, del.GetCoin().Amount))
			}
			tokenBase := formulas.CalculateSaleReturn(coin.Volume, coin.Reserve, coin.CRR, del.GetCoin().Amount)
			eventMutex.Lock()
			ctx.EventManager().EmitEvent(sdk.NewEvent(
				types.EventTypeCalcStake,
				sdk.NewAttribute(types.AttributeKeyValidator, validator.ValAddress.String()),
				sdk.NewAttribute(types.AttributeKeyDelegator, del.GetDelegatorAddr().String()),
				sdk.NewAttribute(types.AttributeKeyCoin, del.GetCoin().String()),
				sdk.NewAttribute(types.AttributeKeyStake, tokenBase.String()),
			))
			eventMutex.Unlock()
			switch del := del.(type) {
			case types.Delegation:
				k.SetDelegation(ctx, del)
			case types.DelegationNFT:
				k.SetDelegationNFT(ctx, del)
			}
		}
		mutex.Lock()
		total = total.Add(del.GetTokensBase())
		mutex.Unlock()
	}
	return total
}

//_______________________________________________________________________
// Validator Queue

// gets a specific validator queue timeSlice. A timeSlice is a slice of ValAddresses corresponding to unbonding validators
// that expire at a certain time.
func (k Keeper) GetValidatorQueueTimeSlice(ctx sdk.Context, timestamp time.Time) []sdk.ValAddress {
	var valAddr []sdk.ValAddress
	err := k.Get(ctx, types.GetValidatorQueueTimeKey(timestamp), &valAddr)
	if valAddr == nil || err != nil {
		return []sdk.ValAddress{}
	}
	return valAddr
}

// Sets a specific validator queue timeSlice.
func (k Keeper) SetValidatorQueueTimeSlice(ctx sdk.Context, timestamp time.Time, keys []sdk.ValAddress) error {
	return k.set(ctx, types.GetValidatorQueueTimeKey(timestamp), keys)
}

// Deletes a specific validator queue timeSlice.
func (k Keeper) DeleteValidatorQueueTimeSlice(ctx sdk.Context, timestamp time.Time) {
	k.delete(ctx, types.GetValidatorQueueTimeKey(timestamp))
}

// Insert an validator address to the appropriate timeslice in the validator queue
func (k Keeper) InsertValidatorQueue(ctx sdk.Context, val types.Validator) error {
	timeSlice := k.GetValidatorQueueTimeSlice(ctx, val.UnbondingCompletionTime)
	var keys []sdk.ValAddress
	if len(timeSlice) == 0 {
		keys = []sdk.ValAddress{val.ValAddress}
	} else {
		keys = append(timeSlice, val.ValAddress)
	}
	return k.SetValidatorQueueTimeSlice(ctx, val.UnbondingCompletionTime, keys)
}

// Delete a validator address from the validator queue
func (k Keeper) DeleteValidatorQueue(ctx sdk.Context, val types.Validator) error {
	timeSlice := k.GetValidatorQueueTimeSlice(ctx, val.UnbondingCompletionTime)
	var newTimeSlice []sdk.ValAddress
	for _, addr := range timeSlice {
		if !bytes.Equal(addr, val.ValAddress) {
			newTimeSlice = append(newTimeSlice, addr)
		}
	}
	if len(newTimeSlice) == 0 {
		k.DeleteValidatorQueueTimeSlice(ctx, val.UnbondingCompletionTime)
	} else {
		return k.SetValidatorQueueTimeSlice(ctx, val.UnbondingCompletionTime, newTimeSlice)
	}
	return nil
}

// Returns all the validator queue timeslices from time 0 until endTime
func (k Keeper) ValidatorQueueIterator(ctx sdk.Context, endTime time.Time) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return store.Iterator([]byte{types.ValidatorQueueKey}, sdk.InclusiveEndBytes(types.GetValidatorQueueTimeKey(endTime)))
}

// Returns a concatenated list of all the timeslices before currTime, and deletes the timeslices from the queue
func (k Keeper) GetAllMatureValidatorQueue(ctx sdk.Context) (matureValsAddrs []sdk.ValAddress) {
	// gets an iterator for all timeslices from time 0 until the current Blockheader time
	validatorTimesliceIterator := k.ValidatorQueueIterator(ctx, ctx.BlockHeader().Time)
	defer validatorTimesliceIterator.Close()

	for ; validatorTimesliceIterator.Valid(); validatorTimesliceIterator.Next() {
		var timeslice []sdk.ValAddress
		k.cdc.MustUnmarshalBinaryLengthPrefixed(validatorTimesliceIterator.Value(), &timeslice)
		matureValsAddrs = append(matureValsAddrs, timeslice...)
	}

	return matureValsAddrs
}

// Unbonds all the unbonding validators that have finished their unbonding period
func (k Keeper) UnbondAllMatureValidatorQueue(ctx sdk.Context) {
	validatorTimesliceIterator := k.ValidatorQueueIterator(ctx, ctx.BlockHeader().Time)
	defer validatorTimesliceIterator.Close()

	for ; validatorTimesliceIterator.Valid(); validatorTimesliceIterator.Next() {
		var timeslice []sdk.ValAddress
		k.cdc.MustUnmarshalBinaryLengthPrefixed(validatorTimesliceIterator.Value(), &timeslice)

		for _, valAddr := range timeslice {
			val, err := k.GetValidator(ctx, valAddr)
			if err != nil {
				continue
			}

			val, err = k.unbondingToUnbonded(ctx, val)
			if err != nil {
				continue
			}
			if val.Tokens.IsZero() {
				err = k.RemoveValidator(ctx, val.ValAddress)
			}
		}

		k.delete(ctx, validatorTimesliceIterator.Key())
	}
}

func (k Keeper) RemoveValidator(ctx sdk.Context, address sdk.ValAddress) error {
	// first retrieve the old validator record
	validator, err := k.GetValidator(ctx, address)
	if err != nil {
		return err
	}

	if !validator.IsUnbonded() {
		return errors.New("cannot call RemoveValidator on bonded or unbonding validators")
	}
	k.DeleteValidatorByPowerIndex(ctx, validator)
	if !k.TotalStake(ctx, validator).IsZero() {
		return errors.New("attempting to remove a validator which still contains tokens")
	}

	// delete the old validator record
	k.delete(ctx, types.GetValidatorKey(address))
	k.delete(ctx, types.GetValidatorByConsAddrKey(sdk.ConsAddress(validator.PubKey.Address())))
	k.delete(ctx, types.GetValidatorsByPowerIndexKey(validator, validator.Tokens))
	return nil
}

//_______________________________________________________________________
// Last Validator Index

// Set the last validator power.
func (k Keeper) SetLastValidatorPower(ctx sdk.Context, operator sdk.ValAddress, power int64) error {
	return k.set(ctx, types.GetLastValidatorPowerKey(operator), power)
}

// Delete the last validator power.
func (k Keeper) DeleteLastValidatorPower(ctx sdk.Context, operator sdk.ValAddress) {
	k.delete(ctx, types.GetLastValidatorPowerKey(operator))
}

// validator index
func (k Keeper) DeleteValidatorByPowerIndex(ctx sdk.Context, validator types.Validator) {
	k.delete(ctx, types.GetValidatorsByPowerIndexKey(validator, validator.Tokens))
}

// Iterate over last validator powers.
func (k Keeper) IterateLastValidatorPowers(ctx sdk.Context, handler func(operator sdk.ValAddress, power int64) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, []byte{types.LastValidatorPowerKey})
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		addr := sdk.ValAddress(iter.Key()[1:])
		var power int64
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iter.Value(), &power)
		if handler(addr, power) {
			break
		}
	}
}

// iterate through the active validator set and perform the provided function
func (k Keeper) IterateLastValidators(ctx sdk.Context, fn func(index int64, validator types.Validator) (stop bool)) {
	iterator := k.LastValidatorsIterator(ctx)
	defer iterator.Close()
	i := int64(0)
	for ; iterator.Valid(); iterator.Next() {
		address := types.AddressFromLastValidatorPowerKey(iterator.Key())
		validator, err := k.GetValidator(ctx, address)
		if err != nil {
			panic(fmt.Sprintf("validator record not found for address: %v\n", address))
		}

		stop := fn(i, validator) // XXX is this safe will the validator unexposed fields be able to get written to?
		if stop {
			break
		}
		i++
	}
}

// returns an iterator for the consensus validators in the last block
func (k Keeper) LastValidatorsIterator(ctx sdk.Context) (iterator sdk.Iterator) {
	store := ctx.KVStore(k.storeKey)
	iterator = sdk.KVStorePrefixIterator(store, []byte{types.LastValidatorPowerKey})
	return iterator
}

// Delegation get the delegation interface for a particular set of delegator and validator addresses
func (k Keeper) Delegation(ctx sdk.Context, addrDel sdk.AccAddress, addrVal sdk.ValAddress, coin string) types.Delegation {
	bond, ok := k.GetDelegation(ctx, addrDel, addrVal, coin)
	if !ok {
		return types.Delegation{}
	}

	return bond
}

// Update the tokens of an existing validator, update the validators power index key
func (k Keeper) RemoveValidatorTokens(ctx sdk.Context,
	validator types.Validator, tokensToRemove sdk.Int) types.Validator {

	k.DeleteValidatorByPowerIndex(ctx, validator)
	validator = validator.RemoveTokens(tokensToRemove)
	err := k.SetValidator(ctx, validator)
	if err != nil {
		panic(err)
	}
	k.SetValidatorByPowerIndex(ctx, validator)
	return validator
}

// get the current group of bonded validators sorted by power-rank
func (k Keeper) GetBondedValidatorsByPower(ctx sdk.Context) []types.Validator {
	store := ctx.KVStore(k.storeKey)
	maxValidators := k.MaxValidators(ctx)
	validators := make([]types.Validator, maxValidators)

	iterator := sdk.KVStoreReversePrefixIterator(store, []byte{types.ValidatorsByPowerIndexKey})
	defer iterator.Close()

	i := 0
	for ; iterator.Valid() && i < int(maxValidators); iterator.Next() {
		address := iterator.Value()
		validator, err := k.GetValidator(ctx, address)
		if err != nil {
			panic(err)
		}

		if validator.IsBonded() {
			validators[i] = validator
			i++
		}
	}
	return validators[:i] // trim
}

// get the group of the bonded validators
func (k Keeper) GetLastValidators(ctx sdk.Context) (validators []types.Validator) {
	store := ctx.KVStore(k.storeKey)

	// add the actual validator power sorted store
	maxValidators := k.MaxValidators(ctx)
	validators = make([]types.Validator, maxValidators)

	iterator := sdk.KVStorePrefixIterator(store, []byte{types.LastValidatorPowerKey})
	defer iterator.Close()

	i := 0
	for ; iterator.Valid(); iterator.Next() {

		// sanity check
		if i >= int(maxValidators) {
			panic("more validators than maxValidators found")
		}
		address := types.AddressFromLastValidatorPowerKey(iterator.Key())
		validator, err := k.GetValidator(ctx, address)
		if err != nil {
			panic(err)
		}

		validators[i] = validator
		i++
	}
	return validators[:i] // trim
}

// return a given amount of all the validators
func (k Keeper) GetValidators(ctx sdk.Context, maxRetrieve uint16) (validators []types.Validator) {
	store := ctx.KVStore(k.storeKey)
	validators = make([]types.Validator, maxRetrieve)

	iterator := sdk.KVStorePrefixIterator(store, []byte{types.ValidatorsKey})
	defer iterator.Close()

	i := 0
	for ; iterator.Valid() && i < int(maxRetrieve); iterator.Next() {
		validator, err := types.UnmarshalValidator(k.cdc, iterator.Value())
		if err != nil {
			panic(err)
		}
		validators[i] = validator
		i++
	}
	return validators[:i] // trim if the array length < maxRetrieve
}

func (k Keeper) IsDelegatorStakeSufficient(ctx sdk.Context, validator types.Validator, delAddr sdk.AccAddress, stake sdk.Coin) (bool, error) {
	delegations := k.GetValidatorDelegations(ctx, validator.ValAddress)
	if uint16(len(delegations)) < k.MaxDelegations(ctx) {
		return true, nil
	}

	stakeValue, err := k.CalculateBipValue(ctx, stake, true)
	if err != nil {
		return false, err
	}

	for _, delegation := range delegations {
		if delegation.GetCoin().Amount.LT(stakeValue) || (delAddr.Equals(delegation.GetDelegatorAddr()) && stake.Denom == delegation.GetCoin().Denom) {
			return true, nil
		}
	}

	return false, nil
}

func (k Keeper) CalculateBipValue(ctx sdk.Context, value sdk.Coin, includeSelf bool) (sdk.Int, error) {
	if value.Denom == k.BondDenom(ctx) {
		return value.Amount, nil
	}

	totalAmount := sdk.ZeroInt()

	if includeSelf {
		totalAmount = totalAmount.Add(value.Amount)
	}

	validators := k.GetAllValidators(ctx)
	for _, validator := range validators {
		stakes := k.GetValidatorDelegations(ctx, validator.ValAddress)
		for _, stake := range stakes {
			if stake.GetCoin().Denom == value.Denom {
				totalAmount = totalAmount.Add(stake.GetCoin().Amount)
			}
		}
	}

	coin, err := k.GetCoin(ctx, value.Denom)
	if err != nil {
		return sdk.Int{}, err
	}

	totalPower := formulas.CalculateSaleReturn(coin.Volume, coin.Reserve, coin.CRR, totalAmount)

	return totalPower.Mul(value.Amount).Quo(totalAmount), nil
}

func (k Keeper) DecreaseValidatorTokens(ctx sdk.Context, validator types.Validator, amount sdk.Int) sdk.Int {
	validator.Tokens = validator.Tokens.Sub(amount)
	err := k.SetValidator(ctx, validator)
	if err != nil {
		panic(err)
	}
	if !validator.Jailed {
		ctx.KVStore(k.storeKey).Set(types.GetValidatorsByPowerIndexKey(validator, validator.Tokens), validator.ValAddress)
	}
	return validator.Tokens
}
