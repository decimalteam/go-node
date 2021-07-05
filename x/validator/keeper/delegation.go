package keeper

import (
	"bitbucket.org/decimalteam/go-node/utils/formulas"
	"bitbucket.org/decimalteam/go-node/x/nft"
	"bitbucket.org/decimalteam/go-node/x/validator/exported"
	"fmt"
	"log"
	"time"

	"bitbucket.org/decimalteam/go-node/x/validator/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// return a specific delegation
func (k Keeper) GetDelegation(ctx sdk.Context,
	delAddr sdk.AccAddress, valAddr sdk.ValAddress, coin string) (
	delegation types.Delegation, found bool) {

	store := ctx.KVStore(k.storeKey)
	key := types.GetDelegationKey(delAddr, valAddr, coin)
	value := store.Get(key)
	if value == nil {
		return delegation, false
	}

	delegation = types.MustUnmarshalDelegation(k.cdc, value)
	return delegation, true
}

// IterateAllDelegations iterate through all of the delegations
func (k Keeper) IterateAllDelegations(ctx sdk.Context, cb func(delegation exported.DelegationI) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, []byte{types.DelegationKey})
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		delegation := types.MustUnmarshalDelegation(k.cdc, iterator.Value())
		if cb(delegation) {
			break
		}
	}

	iteratorNFT := sdk.KVStorePrefixIterator(store, []byte{types.DelegationNFTKey})
	defer iteratorNFT.Close()

	for ; iteratorNFT.Valid(); iteratorNFT.Next() {
		delegation := types.MustUnmarshalDelegationNFT(k.cdc, iteratorNFT.Value())
		if cb(delegation) {
			break
		}
	}
}

// GetAllDelegations returns all delegations used during genesis dump
func (k Keeper) GetAllDelegations(ctx sdk.Context) (delegations []exported.DelegationI) {
	k.IterateAllDelegations(ctx, func(delegation exported.DelegationI) bool {
		delegations = append(delegations, delegation)
		return false
	})
	return delegations
}

// return all delegations to a specific validator. Useful for querier.
func (k Keeper) GetValidatorDelegations(ctx sdk.Context, valAddr sdk.ValAddress) (delegations []exported.DelegationI) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, []byte{types.DelegationKey})
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		delegation := types.MustUnmarshalDelegation(k.cdc, iterator.Value())
		if delegation.GetValidatorAddr().Equals(valAddr) {
			delegations = append(delegations, delegation)
		}
	}

	iteratorNFT := sdk.KVStorePrefixIterator(store, []byte{types.DelegationNFTKey})
	defer iteratorNFT.Close()

	for ; iteratorNFT.Valid(); iteratorNFT.Next() {
		delegation := types.MustUnmarshalDelegationNFT(k.cdc, iteratorNFT.Value())
		if delegation.GetValidatorAddr().Equals(valAddr) {
			delegations = append(delegations, delegation)
		}
	}
	return delegations
}

// return a given amount of all the delegations from a delegator
func (k Keeper) GetDelegatorDelegations(ctx sdk.Context, delegator sdk.AccAddress,
	maxRetrieve uint16) (delegations []exported.DelegationI) {

	delegations = make([]exported.DelegationI, maxRetrieve)

	store := ctx.KVStore(k.storeKey)
	delegatorPrefixKey := types.GetDelegationsKey(delegator)
	iterator := sdk.KVStorePrefixIterator(store, delegatorPrefixKey)
	defer iterator.Close()

	i := 0
	for ; iterator.Valid() && i < int(maxRetrieve); iterator.Next() {
		delegation := types.MustUnmarshalDelegation(k.cdc, iterator.Value())
		delegations[i] = delegation
		i++
	}
	return delegations[:i] // trim if the array length < maxRetrieve
}

// set a delegation
func (k Keeper) SetDelegation(ctx sdk.Context, delegation types.Delegation) {
	delegatorAddr, err := sdk.AccAddressFromBech32(delegation.DelegatorAddress)
	if err != nil {
		panic(err)
	}
	validatorAddr, err := sdk.ValAddressFromBech32(delegation.ValidatorAddress)
	if err != nil {
		panic(err)
	}

	delegation.TokensBase = k.CalcTokensBase(ctx, delegation)
	err = k.set(ctx, types.GetDelegationKey(delegatorAddr, validatorAddr, delegation.Coin.Denom), delegation)
	if err != nil {
		panic(err)
	}
}

// remove a delegation
func (k Keeper) RemoveDelegation(ctx sdk.Context, delegatorAddr sdk.AccAddress, validatorAddr sdk.ValAddress, delegation types.Delegation) {
	// TODO: Consider calling hooks outside of the store wrapper functions, it's unobvious.
	k.BeforeDelegationRemoved(ctx, delegatorAddr, validatorAddr)
	k.delete(ctx, types.GetDelegationKey(delegatorAddr, validatorAddr, delegation.Coin.Denom))
}

func (k Keeper) CalcTokensBase(ctx sdk.Context, delegation exported.DelegationI) sdk.Int {
	var tokensBase sdk.Int
	if delegation.GetCoin().Denom != k.BondDenom(ctx) {
		coin, err := k.GetCoin(ctx, delegation.GetCoin().Denom)
		if err != nil {
			panic(err)
		}
		tokensBase = formulas.CalculateSaleReturn(coin.Volume, coin.Reserve, coin.CRR, delegation.GetCoin().Amount)
	} else {
		tokensBase = delegation.GetCoin().Amount
	}
	return tokensBase
}

// return a given amount of all the delegator unbonding-delegations
func (k Keeper) GetUnbondingDelegations(ctx sdk.Context, delegator sdk.AccAddress,
	maxRetrieve uint16) (unbondingDelegations []types.UnbondingDelegation) {

	unbondingDelegations = make([]types.UnbondingDelegation, maxRetrieve)

	store := ctx.KVStore(k.storeKey)
	delegatorPrefixKey := types.GetUBDsKey(delegator)
	iterator := sdk.KVStorePrefixIterator(store, delegatorPrefixKey)
	defer iterator.Close()

	i := 0
	for ; iterator.Valid() && i < int(maxRetrieve); iterator.Next() {
		unbondingDelegation := types.MustUnmarshalUBD(k.cdc, iterator.Value())
		unbondingDelegations[i] = unbondingDelegation
		i++
	}
	return unbondingDelegations[:i] // trim if the array length < maxRetrieve
}

// return a unbonding delegation
func (k Keeper) GetUnbondingDelegation(ctx sdk.Context,
	delAddr sdk.AccAddress, valAddr sdk.ValAddress) (ubd types.UnbondingDelegation, found bool) {

	store := ctx.KVStore(k.storeKey)
	key := types.GetUBDKey(delAddr, valAddr)
	value := store.Get(key)
	if value == nil {
		return ubd, false
	}

	ubd = types.MustUnmarshalUBD(k.cdc, value)
	return ubd, true
}

// return all unbonding delegations from a particular validator
func (k Keeper) GetUnbondingDelegationsFromValidator(ctx sdk.Context, valAddr sdk.ValAddress) (ubds []types.UnbondingDelegation) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.GetUBDsByValIndexKey(valAddr))
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		key := types.GetUBDKeyFromValIndexKey(iterator.Key())
		value := store.Get(key)
		ubd := types.MustUnmarshalUBD(k.cdc, value)
		ubds = append(ubds, ubd)
	}
	return ubds
}

// iterate through all of the unbonding delegations
func (k Keeper) IterateUnbondingDelegations(ctx sdk.Context, fn func(index int64, ubd types.UnbondingDelegation) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, []byte{types.UnbondingDelegationKey})
	defer iterator.Close()

	for i := int64(0); iterator.Valid(); iterator.Next() {
		ubd := types.MustUnmarshalUBD(k.cdc, iterator.Value())
		if stop := fn(i, ubd); stop {
			break
		}
		i++
	}
}

// HasMaxUnbondingDelegationEntries - check if unbonding delegation has maximum number of entries
func (k Keeper) HasMaxUnbondingDelegationEntries(ctx sdk.Context,
	delegatorAddr sdk.AccAddress, validatorAddr sdk.ValAddress) bool {

	ubd, found := k.GetUnbondingDelegation(ctx, delegatorAddr, validatorAddr)
	if !found {
		return false
	}
	return len(ubd.Entries) >= int(k.MaxEntries(ctx))
}

// set the unbonding delegation and associated index
func (k Keeper) SetUnbondingDelegation(ctx sdk.Context, ubd types.UnbondingDelegation) {
	delegatorAddr, err := sdk.AccAddressFromBech32(ubd.DelegatorAddress)
	if err != nil {
		panic(err)
	}
	validatorAddr, err := sdk.ValAddressFromBech32(ubd.ValidatorAddress)
	if err != nil {
		panic(err)
	}

	store := ctx.KVStore(k.storeKey)
	bz := types.MustMarshalUBD(k.cdc, ubd)
	key := types.GetUBDKey(delegatorAddr, validatorAddr)
	store.Set(key, bz)
	store.Set(types.GetUBDByValIndexKey(delegatorAddr, validatorAddr), []byte{}) // index, store empty bytes
}

// remove the unbonding delegation object and associated index
func (k Keeper) RemoveUnbondingDelegation(ctx sdk.Context, delegatorAddr sdk.AccAddress, validatorAddr sdk.ValAddress, ubd types.UnbondingDelegation) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetUBDKey(delegatorAddr, validatorAddr)
	store.Delete(key)
	store.Delete(types.GetUBDByValIndexKey(delegatorAddr, validatorAddr))
}

// SetUnbondingDelegationEntry adds an entry to the unbonding delegation at
// the given addresses. It creates the unbonding delegation if it does not exist
func (k Keeper) SetUnbondingDelegationEntry(ctx sdk.Context,
	delegatorAddr sdk.AccAddress, validatorAddr sdk.ValAddress,
	creationHeight int64, minTime time.Time, balance sdk.Coin) types.UnbondingDelegation {

	ubd, found := k.GetUnbondingDelegation(ctx, delegatorAddr, validatorAddr)
	if found {
		ubd.AddEntry(creationHeight, minTime, balance)
	} else {
		ubd = types.NewUnbondingDelegation(delegatorAddr, validatorAddr, types.NewUnbondingDelegationEntry(creationHeight, minTime, balance))
	}
	k.SetUnbondingDelegation(ctx, ubd)
	return ubd
}

// unbonding delegation queue timeslice operations

// gets a specific unbonding queue timeslice. A timeslice is a slice of DVPairs
// corresponding to unbonding delegations that expire at a certain time.
func (k Keeper) GetUBDQueueTimeSlice(ctx sdk.Context, timestamp time.Time) (dvPairs []types.DVPair) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetUnbondingDelegationTimeKey(timestamp))
	if bz == nil {
		return []types.DVPair{}
	}
	k.cdc.MustUnmarshalLengthPrefixed(bz, &dvPairs)
	return dvPairs
}

// Sets a specific unbonding queue timeslice.
func (k Keeper) SetUBDQueueTimeSlice(ctx sdk.Context, timestamp time.Time, keys []types.DVPair) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalLengthPrefixed(keys)
	store.Set(types.GetUnbondingDelegationTimeKey(timestamp), bz)
}

// Insert an unbonding delegation to the appropriate timeslice in the unbonding queue
func (k Keeper) InsertUBDQueue(ctx sdk.Context, ubd types.UnbondingDelegation,
	completionTime time.Time) {

	timeSlice := k.GetUBDQueueTimeSlice(ctx, completionTime)
	dvPair := types.DVPair{DelegatorAddress: ubd.DelegatorAddress, ValidatorAddress: ubd.ValidatorAddress}
	if len(timeSlice) == 0 {
		k.SetUBDQueueTimeSlice(ctx, completionTime, []types.DVPair{dvPair})
	} else {
		timeSlice = append(timeSlice, dvPair)
		k.SetUBDQueueTimeSlice(ctx, completionTime, timeSlice)
	}
}

// Returns all the unbonding queue timeslices from time 0 until endTime
func (k Keeper) UBDQueueIterator(ctx sdk.Context, endTime time.Time) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return store.Iterator([]byte{types.UnbondingQueueKey},
		sdk.InclusiveEndBytes(types.GetUnbondingDelegationTimeKey(endTime)))
}

// Returns a concatenated list of all the timeslices inclusively previous to
// currTime, and deletes the timeslices from the queue
func (k Keeper) DequeueAllMatureUBDQueue(ctx sdk.Context,
	currTime time.Time) (matureUnbonds []types.DVPair) {

	store := ctx.KVStore(k.storeKey)
	// gets an iterator for all timeslices from time 0 until the current Blockheader time
	unbondingTimesliceIterator := k.UBDQueueIterator(ctx, ctx.BlockHeader().Time)
	for ; unbondingTimesliceIterator.Valid(); unbondingTimesliceIterator.Next() {
		timeslice := []types.DVPair{}
		value := unbondingTimesliceIterator.Value()
		k.cdc.MustUnmarshalLengthPrefixed(value, &timeslice)
		matureUnbonds = append(matureUnbonds, timeslice...)
		store.Delete(unbondingTimesliceIterator.Key())
	}
	return matureUnbonds
}

// Perform a delegation, set/update everything necessary within the store.
// tokenSrc indicates the bond status of the incoming funds.
func (k Keeper) Delegate(ctx sdk.Context, delAddr sdk.AccAddress, bondCoin sdk.Coin, tokenSrc types.BondStatus, validator types.Validator, subtractAccount bool) error {
	valAddr, _ := sdk.ValAddressFromBech32(validator.ValAddress)

	// Get or create the delegation object
	delegation, found := k.GetDelegation(ctx, delAddr, valAddr, bondCoin.Denom)
	if !found {
		delegation = types.NewDelegation(delAddr, valAddr, bondCoin)
	}

	// call the appropriate hook if present
	if found {
		k.BeforeDelegationSharesModified(ctx, delAddr, valAddr)
	} else {
		k.BeforeDelegationCreated(ctx, delAddr, valAddr)
	}

	delegatorAddr, err := sdk.AccAddressFromBech32(delegation.DelegatorAddress)
	if err != nil {
		return err
	}
	delegatorValAddr, err := sdk.ValAddressFromBech32(delegation.ValidatorAddress)
	if err != nil {
		return err
	}

	// if subtractAccount is true then we are
	// performing a delegation and not a redelegation, thus the source tokens are
	// all non bonded
	if subtractAccount {
		if tokenSrc == types.Bonded {
			panic("delegation token source cannot be bonded")
		}

		var sendName string
		switch {
		case validator.IsBonded():
			sendName = types.BondedPoolName
		case validator.IsUnbonding(), validator.IsUnbonded():
			sendName = types.NotBondedPoolName
		default:
			panic("invalid validator status")
		}

		err = k.baseKeeper.DelegateCoinsFromAccountToModule(ctx, delegatorAddr, sendName, sdk.NewCoins(bondCoin))
		if err != nil {
			return err
		}
	} else {

		// potentially transfer tokens between pools, if
		switch {
		case tokenSrc == types.Bonded && validator.IsBonded():
			// do nothing
		case (tokenSrc == types.Unbonded || tokenSrc == types.Unbonding) && !validator.IsBonded():
			// do nothing
		case (tokenSrc == types.Unbonded || tokenSrc == types.Unbonding) && validator.IsBonded():
			// transfer pools
			k.notBondedTokensToBonded(ctx, sdk.NewCoins(bondCoin))
		case tokenSrc == types.Bonded && !validator.IsBonded():
			// transfer pools
			k.bondedTokensToNotBonded(ctx, sdk.NewCoins(bondCoin))
		default:
			panic("unknown token source bond status")
		}
	}

	// Update delegation
	if !found {
		delegation.Coin = bondCoin
	} else {
		if delegation.Coin.Denom == bondCoin.Denom {
			delegation.Coin = delegation.Coin.Add(bondCoin)
		} else {
			delegation = types.NewDelegation(delAddr, valAddr, bondCoin)
		}
	}

	k.SetDelegation(ctx, delegation)

	k.DeleteValidatorByPowerIndex(ctx, validator)
	if bondCoin.Denom == k.BondDenom(ctx) {
		validator.Tokens = validator.Tokens.Add(bondCoin.Amount)
	} else {
		coin, err := k.GetCoin(ctx, bondCoin.Denom)
		if err != nil {
			return err
		}
		validator.Tokens = validator.Tokens.Add(formulas.CalculateSaleReturn(coin.Volume, coin.Reserve, coin.CRR, bondCoin.Amount))
	}
	err = k.SetValidator(ctx, validator)
	if err != nil {
		return err
	}
	k.SetValidatorByPowerIndexWithoutCalc(ctx, valAddr, validator)

	// Call the after-modification hook
	k.AfterDelegationModified(ctx, delegatorAddr, delegatorValAddr)

	return nil
}

func (k Keeper) CheckTotalStake(ctx sdk.Context, validator types.Validator) bool {
	defer func() {
		if r := recover(); r != nil {
			log.Println(r)
			return
		}
	}()
	k.DeleteValidatorByPowerIndex(ctx, validator)
	k.SetValidatorByPowerIndex(ctx, validator)
	return true
}

// Undelegate unbonds an amount of delegator shares from a given validator. It
// will verify that the unbonding entries between the delegator and validator
// are not exceeded and unbond the staked tokens (based on shares) by creating
// an unbonding object and inserting it into the unbonding queue which will be
// processed during the staking EndBlocker.
func (k Keeper) Undelegate(
	ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress, amount sdk.Coin,
) (time.Time, error) {

	validator, foundErr := k.GetValidator(ctx, valAddr)
	if foundErr != nil {
		return time.Time{}, types.ErrNoDelegatorForAddress()
	}

	err := k.unbond(ctx, delAddr, valAddr, amount)
	if err != nil {
		return time.Time{}, err
	}

	// transfer the validator tokens to the not bonded pool
	if validator.IsBonded() {
		k.bondedTokensToNotBonded(ctx, sdk.NewCoins(amount))
	}

	completionTime := ctx.BlockHeader().Time.Add(k.UnbondingTime(ctx))
	ubd := k.SetUnbondingDelegationEntry(ctx, delAddr, valAddr, ctx.BlockHeight(), completionTime, amount)
	k.InsertUBDQueue(ctx, ubd, completionTime)

	return completionTime, nil
}

// unbond a particular delegation and perform associated store operations
func (k Keeper) unbond(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress, coin sdk.Coin) error {
	// check if a delegation object exists in the store
	delegation, found := k.GetDelegation(ctx, delAddr, valAddr, coin.Denom)
	if !found {
		return types.ErrNoDelegatorForAddress()
	}

	delegatorAddr, err := sdk.AccAddressFromBech32(delegation.DelegatorAddress)
	if err != nil {
		return err
	}
	delegatorValAddr, err := sdk.ValAddressFromBech32(delegation.ValidatorAddress)
	if err != nil {
		return err
	}

	// call the before-delegation-modified hook
	k.BeforeDelegationSharesModified(ctx, delAddr, valAddr)

	// ensure that we have enough shares to remove
	if delegation.Coin.Amount.LT(coin.Amount) {
		return types.ErrNotEnoughDelegationShares(delegation.Coin.Amount.String())
	}

	// get validator
	validator, err := k.GetValidator(ctx, valAddr)
	if err != nil {
		return types.ErrNoValidatorFound()
	}

	// subtract shares from delegation
	delegation.Coin = delegation.Coin.Sub(coin)

	// remove the delegation
	if delegation.Coin.IsZero() {
		k.RemoveDelegation(ctx, delegatorAddr, delegatorValAddr, delegation)
	} else {
		k.SetDelegation(ctx, delegation)
		// call the after delegation modification hook
		k.AfterDelegationModified(ctx, delegatorAddr, delegatorValAddr)
	}

	k.DeleteValidatorByPowerIndex(ctx, validator)

	amountBase := coin.Amount
	if coin.Denom != k.BondDenom(ctx) {
		c, err := k.GetCoin(ctx, coin.Denom)
		if err != nil {
			return types.ErrInternal(err.Error())
		}
		amountBase = formulas.CalculateSaleReturn(c.Volume, c.Reserve, c.CRR, coin.Amount)
	}
	decreasedTokens := k.DecreaseValidatorTokens(ctx, valAddr, validator, amountBase)

	if decreasedTokens.IsZero() && validator.IsUnbonded() {
		// if not unbonded, we must instead remove validator in EndBlocker once it finishes its unbonding period
		err = k.RemoveValidator(ctx, valAddr)
		if err != nil {
			return types.ErrInternal(err.Error())
		}
	}

	return nil
}

// CompleteUnbonding completes the unbonding of all mature entries in the
// retrieved unbonding delegation object.
func (k Keeper) CompleteUnbonding(ctx sdk.Context, delAddr sdk.AccAddress,
	valAddr sdk.ValAddress) error {

	ubd, found := k.GetUnbondingDelegation(ctx, delAddr, valAddr)
	if !found {
		return types.ErrUnbondingDelegationNotFound()
	}

	ctxTime := ctx.BlockHeader().Time

	delegatorAddr, err := sdk.AccAddressFromBech32(ubd.DelegatorAddress)
	if err != nil {
		return err
	}

	// loop through all the entries and complete unbonding mature entries
	for i := 0; i < len(ubd.Entries); i++ {
		entry := ubd.Entries[i]
		if entry.IsMature(ctxTime) {
			ubd.RemoveEntry(int64(i))
			i--

			// track undelegation only when remaining or truncated shares are non-zero
			if !entry.GetBalance().IsZero() {
				switch entry := entry.(type) {
				case types.UnbondingDelegationEntry:
					amt := sdk.NewCoins(entry.Balance)
					err := k.baseKeeper.UndelegateCoinsFromModuleToAccount(ctx, types.NotBondedPoolName, delegatorAddr, amt)
					if err != nil {
						return err
					}
				case types.UnbondingDelegationNFTEntry:
					collection, ok := k.nftKeeper.GetCollection(ctx, entry.Denom)
					if !ok {
						return fmt.Errorf("collection not found")
					}

					token, err := collection.GetNFT(entry.TokenID)
					if err != nil {
						return err
					}

					owner := token.GetOwners().GetOwner(delAddr)
					if owner == nil {
						token = token.SetOwners(token.GetOwners().SetOwner(&nft.TokenOwner{
							Address:  delAddr,
							Quantity: entry.Quantity,
						}))
					} else {
						token = token.SetOwners(token.GetOwners().SetOwner(owner.SetQuantity(owner.GetQuantity().Add(entry.Quantity))))
					}

					collection, err = collection.UpdateNFT(token)
					if err != nil {
						return err
					}

					k.nftKeeper.SetCollection(ctx, entry.Denom, collection)
				default:
					panic(fmt.Sprintf("%T", entry))
				}
			}
		}
	}

	// set the unbonding delegation or remove it if there are no more entries
	if len(ubd.Entries) == 0 {
		k.RemoveUnbondingDelegation(ctx, delegatorAddr, valAddr, ubd)
	} else {
		k.SetUnbondingDelegation(ctx, ubd)
	}

	return nil
}

//_____________________________________________________________________________________

// return all delegations for a delegator
func (k Keeper) GetAllDelegatorDelegations(ctx sdk.Context, delegator sdk.AccAddress) []exported.DelegationI {
	delegations := make([]exported.DelegationI, 0)

	store := ctx.KVStore(k.storeKey)
	delegatorPrefixKey := types.GetDelegationsKey(delegator)
	iterator := sdk.KVStorePrefixIterator(store, delegatorPrefixKey) //smallest to largest
	defer iterator.Close()

	i := 0
	for ; iterator.Valid(); iterator.Next() {
		delegation := types.MustUnmarshalDelegation(k.cdc, iterator.Value())
		delegations = append(delegations, delegation)
		i++
	}

	return delegations
}

// Return all validators that a delegator is bonded to. If maxRetrieve is supplied, the respective amount will be returned.
func (k Keeper) GetDelegatorValidators(ctx sdk.Context, delegatorAddr sdk.AccAddress,
	maxRetrieve uint16) (validators []types.Validator) {
	validators = make([]types.Validator, maxRetrieve)

	store := ctx.KVStore(k.storeKey)
	delegatorPrefixKey := types.GetDelegationsKey(delegatorAddr)
	iterator := sdk.KVStorePrefixIterator(store, delegatorPrefixKey) // smallest to largest
	defer iterator.Close()

	i := 0
	for ; iterator.Valid() && i < int(maxRetrieve); iterator.Next() {
		delegation := types.MustUnmarshalDelegation(k.cdc, iterator.Value())

		valAddr, err := sdk.ValAddressFromBech32(delegation.ValidatorAddress)
		if err != nil {
			panic(err)
		}

		validator, err := k.GetValidator(ctx, valAddr)
		if err != nil {
			panic(types.ErrNoValidatorFound())
		}
		validators[i] = validator
		i++
	}
	return validators[:i] // trim
}

// return a validator that a delegator is bonded to
func (k Keeper) GetDelegatorValidator(ctx sdk.Context, delegatorAddr sdk.AccAddress,
	validatorAddr sdk.ValAddress, coin string) (types.Validator, error) {

	var err error
	validator := types.Validator{}

	delegation, found := k.GetDelegation(ctx, delegatorAddr, validatorAddr, coin)
	if !found {
		return validator, types.ErrNoDelegation()
	}

	valAddr, err := sdk.ValAddressFromBech32(delegation.ValidatorAddress)
	if err != nil {
		return types.Validator{}, err
	}

	validator, err = k.GetValidator(ctx, valAddr)
	if err != nil {
		panic(types.ErrNoValidatorFound())
	}
	return validator, nil
}

// return all unbonding-delegations for a delegator
func (k Keeper) GetUnbondingDelegationsByDelegator(ctx sdk.Context, delegator sdk.AccAddress) []types.UnbondingDelegation {
	var unbondingDelegations []types.UnbondingDelegation

	store := ctx.KVStore(k.storeKey)
	delegatorPrefixKey := types.GetUBDsKey(delegator)
	iterator := sdk.KVStorePrefixIterator(store, delegatorPrefixKey) // smallest to largest
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		unbondingDelegation := types.MustUnmarshalUBD(k.cdc, iterator.Value())
		unbondingDelegations = append(unbondingDelegations, unbondingDelegation)
	}

	return unbondingDelegations
}

func (k Keeper) GetAllUnbondingDelegations(ctx sdk.Context) []types.UnbondingDelegation {
	var unbondingDelegations []types.UnbondingDelegation

	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, []byte{types.UnbondingDelegationKey}) // smallest to largest
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		unbondingDelegation := types.MustUnmarshalUBD(k.cdc, iterator.Value())
		unbondingDelegations = append(unbondingDelegations, unbondingDelegation)
	}

	return unbondingDelegations
}
