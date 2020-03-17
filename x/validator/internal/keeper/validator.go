package keeper

import (
	"bitbucket.org/decimalteam/go-node/utils/formulas"
	"bitbucket.org/decimalteam/go-node/x/validator/internal/types"
	"bytes"
	"errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"time"
)

// Cache the amino decoding of validators, as it can be the case that repeated slashing calls
// cause many calls to GetValidator, which were shown to throttle the state machine in our
// simulation. Note this is quite biased though, as the simulator does more slashes than a
// live chain should, however we require the slashing to be fast as noone pays gas for it.
type cachedValidator struct {
	val        types.Validator
	marshalled string // marshalled amino bytes for the validator object (not operator address)
}

func newCachedValidator(val types.Validator, marshalled string) cachedValidator {
	return cachedValidator{
		val:        val,
		marshalled: marshalled,
	}
}

// get a single validator
func (k Keeper) GetValidator(ctx sdk.Context, addr sdk.ValAddress) (types.Validator, error) {
	store := ctx.KVStore(k.storeKey)
	value := store.Get(types.GetValidatorKey(addr))
	if value == nil {
		return types.Validator{}, errors.New("not found validator ")
	}

	// If these amino encoded bytes are in the cache, return the cached validator
	strValue := string(value)
	if val, ok := k.validatorCache[strValue]; ok {
		valToReturn := val.val
		// Doesn't mutate the cache's value
		valToReturn.ValAddress = addr
		return valToReturn, nil
	}

	// amino bytes weren't found in cache, so amino unmarshal and add it to the cache
	validator, err := types.UnmarshalValidator(k.cdc, value)
	if err != nil {
		return types.Validator{}, errors.New("error unmarshal validator ")
	}
	cachedVal := newCachedValidator(validator, strValue)
	k.validatorCache[strValue] = newCachedValidator(validator, strValue)
	k.validatorCacheList.PushBack(cachedVal)

	// if the cache is too big, pop off the last element from it
	if k.validatorCacheList.Len() > aminoCacheSize {
		valToRemove := k.validatorCacheList.Remove(k.validatorCacheList.Front()).(cachedValidator)
		delete(k.validatorCache, valToRemove.marshalled)
	}

	return validator, nil
}

func (k Keeper) SetValidator(ctx sdk.Context, validator types.Validator) error {
	return k.set(ctx, types.GetValidatorKey(validator.ValAddress), validator)
}

// validator index
func (k Keeper) SetValidatorByConsAddr(ctx sdk.Context, validator types.Validator) error {
	consAddr := sdk.ConsAddress(validator.PubKey.Address())
	return k.set(ctx, types.GetValidatorByConsAddrKey(consAddr), validator.ValAddress)
}

// validator index
func (k Keeper) SetValidatorByPowerIndex(ctx sdk.Context, validator types.Validator) error {
	// jailed validators are not kept in the power index
	if validator.Jailed {
		return nil
	}
	power := k.TotalPowerValidator(ctx, validator)
	return k.set(ctx, types.GetValidatorsByPowerIndexKey(validator, power), validator.ValAddress)
}

func (k Keeper) TotalPowerValidator(ctx sdk.Context, validator types.Validator) sdk.Int {
	power := sdk.Int{}
	for _, token := range validator.StakeCoins {
		coin, err := k.coinKeeper.GetCoin(ctx, token.Denom)
		if err != nil {
			continue
		}
		power.Add(formulas.CalculateSaleReturn(coin.Volume, coin.Reserve, coin.ConstantReserveRatio, token.Amount))
	}
	return power
}

//_______________________________________________________________________
// Validator Queue

// gets a specific validator queue timeslice. A timeslice is a slice of ValAddresses corresponding to unbonding validators
// that expire at a certain time.
func (k Keeper) GetValidatorQueueTimeSlice(ctx sdk.Context, timestamp time.Time) (valAddrs []sdk.ValAddress) {
	bz, err := k.Get(ctx, types.GetValidatorQueueTimeKey(timestamp))
	if bz == nil {
		return []sdk.ValAddress{}
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &valAddrs)
	return valAddrs
}

// Sets a specific validator queue timeslice.
func (k Keeper) SetValidatorQueueTimeSlice(ctx sdk.Context, timestamp time.Time, keys []sdk.ValAddress) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(keys)
	store.Set(types.GetValidatorQueueTimeKey(timestamp), bz)
}

// Deletes a specific validator queue timeslice.
func (k Keeper) DeleteValidatorQueueTimeSlice(ctx sdk.Context, timestamp time.Time) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetValidatorQueueTimeKey(timestamp))
}

// Insert an validator address to the appropriate timeslice in the validator queue
func (k Keeper) InsertValidatorQueue(ctx sdk.Context, val types.Validator) {
	timeSlice := k.GetValidatorQueueTimeSlice(ctx, val.UnbondingCompletionTime)
	var keys []sdk.ValAddress
	if len(timeSlice) == 0 {
		keys = []sdk.ValAddress{val.OperatorAddress}
	} else {
		keys = append(timeSlice, val.OperatorAddress)
	}
	k.SetValidatorQueueTimeSlice(ctx, val.UnbondingCompletionTime, keys)
}

// Delete a validator address from the validator queue
func (k Keeper) DeleteValidatorQueue(ctx sdk.Context, val types.Validator) {
	timeSlice := k.GetValidatorQueueTimeSlice(ctx, val.UnbondingCompletionTime)
	newTimeSlice := []sdk.ValAddress{}
	for _, addr := range timeSlice {
		if !bytes.Equal(addr, val.OperatorAddress) {
			newTimeSlice = append(newTimeSlice, addr)
		}
	}
	if len(newTimeSlice) == 0 {
		k.DeleteValidatorQueueTimeSlice(ctx, val.UnbondingCompletionTime)
	} else {
		k.SetValidatorQueueTimeSlice(ctx, val.UnbondingCompletionTime, newTimeSlice)
	}
}

// Returns all the validator queue timeslices from time 0 until endTime
func (k Keeper) ValidatorQueueIterator(ctx sdk.Context, endTime time.Time) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return store.Iterator(types.ValidatorQueueKey, sdk.InclusiveEndBytes(types.GetValidatorQueueTimeKey(endTime)))
}

// Returns a concatenated list of all the timeslices before currTime, and deletes the timeslices from the queue
func (k Keeper) GetAllMatureValidatorQueue(ctx sdk.Context, currTime time.Time) (matureValsAddrs []sdk.ValAddress) {
	// gets an iterator for all timeslices from time 0 until the current Blockheader time
	validatorTimesliceIterator := k.ValidatorQueueIterator(ctx, ctx.BlockHeader().Time)
	defer validatorTimesliceIterator.Close()

	for ; validatorTimesliceIterator.Valid(); validatorTimesliceIterator.Next() {
		timeslice := []sdk.ValAddress{}
		k.cdc.MustUnmarshalBinaryLengthPrefixed(validatorTimesliceIterator.Value(), &timeslice)
		matureValsAddrs = append(matureValsAddrs, timeslice...)
	}

	return matureValsAddrs
}

// Unbonds all the unbonding validators that have finished their unbonding period
func (k Keeper) UnbondAllMatureValidatorQueue(ctx sdk.Context) {
	store := ctx.KVStore(k.storeKey)
	validatorTimesliceIterator := k.ValidatorQueueIterator(ctx, ctx.BlockHeader().Time)
	defer validatorTimesliceIterator.Close()

	for ; validatorTimesliceIterator.Valid(); validatorTimesliceIterator.Next() {
		timeslice := []sdk.ValAddress{}
		k.cdc.MustUnmarshalBinaryLengthPrefixed(validatorTimesliceIterator.Value(), &timeslice)

		for _, valAddr := range timeslice {
			val, found := k.GetValidator(ctx, valAddr)
			if !found {
				panic("validator in the unbonding queue was not found")
			}

			if !val.IsUnbonding() {
				panic("unexpected validator in unbonding queue; status was not unbonding")
			}
			val = k.unbondingToUnbonded(ctx, val)
			if val.GetDelegatorShares().IsZero() {
				k.RemoveValidator(ctx, val.OperatorAddress)
			}
		}

		store.Delete(validatorTimesliceIterator.Key())
	}
}
