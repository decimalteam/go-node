package keeper

import (
	"bitbucket.org/decimalteam/go-node/x/validator/internal/types"
	"bytes"
	"errors"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
	"sort"
)

// Apply and return accumulated updates to the bonded validator set. Also,
// * Updates the active valset as keyed by LastValidatorPowerKey.
// * Updates the total power as keyed by LastTotalPowerKey.
// * Updates validator status' according to updated powers.
// * Updates the fee pool bonded vs not-bonded tokens.
// * Updates relevant indices.
// It gets called once after genesis, another time maybe after genesis transactions,
// then once at every EndBlock.
//
// CONTRACT: Only validators with non-zero power or zero-power that were bonded
// at the previous block height or were removed from the validator set entirely
// are returned to Tendermint.
func (k Keeper) ApplyAndReturnValidatorSetUpdates(ctx sdk.Context) ([]abci.ValidatorUpdate, error) {
	var updates []abci.ValidatorUpdate

	store := ctx.KVStore(k.storeKey)
	maxValidators := k.GetParams(ctx).MaxValidators
	totalPower := sdk.ZeroInt()
	var amtFromBondedToNotBonded, amtFromNotBondedToBonded sdk.Coins

	// Retrieve the last validator set.
	// The persistent set is updated later in this function.
	// (see LastValidatorPowerKey).
	last := k.getLastValidatorsByAddr(ctx)

	// Iterate over validators, highest power to lowest.
	iterator := sdk.KVStoreReversePrefixIterator(store, []byte{types.ValidatorsByPowerIndexKey})
	defer iterator.Close()
	for count := 0; iterator.Valid() && count < int(maxValidators); iterator.Next() {

		// everything that is iterated in this loop is becoming or already a
		// part of the bonded validator set

		// fetch the validator
		valAddr := sdk.ValAddress(iterator.Value())
		validator, err := k.GetValidator(ctx, valAddr[2:])
		if err != nil {
			return nil, fmt.Errorf("ApplyAndReturnValidatorSetUpdates: %w", err)
		}

		if validator.Jailed {
			return nil, errors.New("ApplyAndReturnValidatorSetUpdates: should never retrieve a jailed validator from the power store")
		}

		// if we get to a zero-power validator (which we don't bond),
		// there are no more possible bonded validators
		if validator.PotentialConsensusPower(k.TotalCoinsValidator(ctx, validator)) == 0 {
			break
		}

		// apply the appropriate state change if necessary
		switch {
		case validator.IsUnbonded():
			validator, err = k.unbondedToBonded(ctx, validator)
			if err != nil {
				return nil, fmt.Errorf("ApplyAndReturnValidatorSetUpdates: %w", err)
			}
			amtFromNotBondedToBonded = amtFromNotBondedToBonded.Add(validator.StakeCoins)
		case validator.IsUnbonding():
			validator, err = k.unbondingToBonded(ctx, validator)
			if err != nil {
				return nil, fmt.Errorf("ApplyAndReturnValidatorSetUpdates: %w", err)
			}
			amtFromNotBondedToBonded = amtFromNotBondedToBonded.Add(validator.StakeCoins)
		case validator.IsBonded():
			// no state change
		default:
			panic("unexpected validator status")
		}

		// fetch the old power bytes
		var valAddrBytes [sdk.AddrLen]byte
		copy(valAddrBytes[:], valAddr[:])
		oldPowerBytes, found := last[valAddrBytes]

		// calculate the new power bytes
		newPower := validator.ConsensusPower(k.TotalCoinsValidator(ctx, validator))
		newPowerBytes := k.cdc.MustMarshalBinaryLengthPrefixed(newPower)

		// update the validator set if power has changed
		if !found || !bytes.Equal(oldPowerBytes, newPowerBytes) {
			updates = append(updates, validator.ABCIValidatorUpdate(k.TotalCoinsValidator(ctx, validator)))

			// set validator power on lookup index
			err = k.SetLastValidatorPower(ctx, valAddr, newPower)
			if err != nil {
				return nil, fmt.Errorf("ApplyAndReturnValidatorSetUpdates: %w", err)
			}
		}

		// validator still in the validator set, so delete from the copy
		delete(last, valAddrBytes)

		// keep count
		count++
		totalPower = totalPower.Add(sdk.NewInt(newPower))
	}

	// sort the no-longer-bonded validators
	noLongerBonded := sortNoLongerBonded(last)

	// iterate through the sorted no-longer-bonded validators
	for _, valAddrBytes := range noLongerBonded {
		// fetch the validator
		validator, err := k.GetValidator(ctx, valAddrBytes)
		if err != nil {
			return nil, fmt.Errorf("ApplyAndReturnValidatorSetUpdates: %w", err)
		}

		// bonded to unbonding
		validator, err = k.bondedToUnbonding(ctx, validator)
		if err != nil {
			return nil, fmt.Errorf("ApplyAndReturnValidatorSetUpdates: %w", err)
		}
		amtFromBondedToNotBonded = amtFromBondedToNotBonded.Add(validator.StakeCoins)

		// delete from the bonded validator index
		k.DeleteLastValidatorPower(ctx, validator.ValAddress)

		// update the validator set
		updates = append(updates, validator.ABCIValidatorUpdateZero())
	}

	// Update the pools based on the recent updates in the validator set:
	// - The tokens from the non-bonded candidates that enter the new validator set need to be transferred
	// to the Bonded pool.
	// - The tokens from the bonded validators that are being kicked out from the validator set
	// need to be transferred to the NotBonded pool.
	k.notBondedTokensToBonded(ctx, amtFromNotBondedToBonded)
	k.bondedTokensToNotBonded(ctx, amtFromBondedToNotBonded)

	// set total power on lookup index if there are any updates
	if len(updates) > 0 {
		err := k.SetLastTotalPower(ctx, totalPower)
		if err != nil {
			return nil, fmt.Errorf("ApplyAndReturnValidatorSetUpdates: %w", err)
		}
	}

	return updates, nil
}

// switches a validator from unbonding state to unbonded state
func (k Keeper) unbondingToUnbonded(ctx sdk.Context, validator types.Validator) (types.Validator, error) {
	if !validator.IsUnbonding() {
		return types.Validator{}, errors.New(fmt.Sprintf("bad state transition unbondingToBonded, validator: %v\n", validator))
	}
	return k.completeUnbondingValidator(ctx, validator)
}

// perform all the store operations for when a validator status becomes unbonded
func (k Keeper) completeUnbondingValidator(ctx sdk.Context, validator types.Validator) (types.Validator, error) {
	validator = validator.UpdateStatus(types.Unbonded)
	err := k.SetValidator(ctx, validator)
	if err != nil {
		return types.Validator{}, err
	}
	return validator, nil
}

// map of operator addresses to serialized power
type validatorsByAddr map[[sdk.AddrLen]byte][]byte

// get the last validator set
func (k Keeper) getLastValidatorsByAddr(ctx sdk.Context) validatorsByAddr {
	last := make(validatorsByAddr)
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, []byte{types.LastValidatorPowerKey})
	defer iterator.Close()
	// iterate over the last validator set index
	for ; iterator.Valid(); iterator.Next() {
		var valAddr [sdk.AddrLen]byte
		// extract the validator address from the key (prefix is 1-byte)
		copy(valAddr[:], iterator.Key()[1:])
		// power bytes is just the value
		powerBytes := iterator.Value()
		last[valAddr] = make([]byte, len(powerBytes))
		copy(last[valAddr][:], powerBytes[:])
	}
	return last
}

func (k Keeper) unbondedToBonded(ctx sdk.Context, validator types.Validator) (types.Validator, error) {
	if !validator.IsUnbonded() {
		return types.Validator{}, fmt.Errorf("bad state transition unbondedToBonded, validator: %v\n", validator)
	}
	return k.bondValidator(ctx, validator)
}

func (k Keeper) unbondingToBonded(ctx sdk.Context, validator types.Validator) (types.Validator, error) {
	if !validator.IsUnbonding() {
		return types.Validator{}, fmt.Errorf("bad state transition unbondingToBonded, validator: %v\n", validator)
	}
	return k.bondValidator(ctx, validator)
}

// perform all the store operations for when a validator status becomes bonded
func (k Keeper) bondValidator(ctx sdk.Context, validator types.Validator) (types.Validator, error) {

	// delete the validator by power index, as the key will change
	k.DeleteValidatorByPowerIndex(ctx, validator)

	// set the status
	validator = validator.UpdateStatus(types.Bonded)

	// save the now bonded validator record to the two referenced stores
	err := k.SetValidator(ctx, validator)
	if err != nil {
		return types.Validator{}, err
	}
	err = k.SetValidatorByPowerIndex(ctx, validator)
	if err != nil {
		return types.Validator{}, err
	}

	// delete from queue if present
	err = k.DeleteValidatorQueue(ctx, validator)
	if err != nil {
		return types.Validator{}, err
	}

	return validator, nil
}

func (k Keeper) bondedToUnbonding(ctx sdk.Context, validator types.Validator) (types.Validator, error) {
	if !validator.IsBonded() {
		return types.Validator{}, fmt.Errorf("bad state transition bondedToUnbonding, validator: %v\n", validator)
	}
	return k.beginUnbondingValidator(ctx, validator)
}

// perform all the store operations for when a validator begins unbonding
func (k Keeper) beginUnbondingValidator(ctx sdk.Context, validator types.Validator) (types.Validator, error) {

	params := k.GetParams(ctx)

	// delete the validator by power index, as the key will change
	k.DeleteValidatorByPowerIndex(ctx, validator)

	// sanity check
	if validator.Status != types.Bonded {
		panic(fmt.Sprintf("should not already be unbonded or unbonding, validator: %v\n", validator))
	}

	// set the status
	validator = validator.UpdateStatus(types.Unbonding)

	// set the unbonding completion time and completion height appropriately
	validator.UnbondingCompletionTime = ctx.BlockHeader().Time.Add(params.UnbondingTime)
	validator.UnbondingHeight = ctx.BlockHeader().Height

	// save the now unbonded validator record and power index
	err := k.SetValidator(ctx, validator)
	if err != nil {
		return types.Validator{}, err
	}
	err = k.SetValidatorByPowerIndex(ctx, validator)
	if err != nil {
		return types.Validator{}, err
	}

	// Adds to unbonding validator queue
	err = k.InsertValidatorQueue(ctx, validator)
	if err != nil {
		return types.Validator{}, err
	}

	return validator, nil
}

// given a map of remaining validators to previous bonded power
// returns the list of validators to be unbonded, sorted by operator address
func sortNoLongerBonded(last validatorsByAddr) [][]byte {
	// sort the map keys for determinism
	noLongerBonded := make([][]byte, len(last))
	index := 0
	for valAddrBytes := range last {
		valAddr := make([]byte, sdk.AddrLen)
		copy(valAddr[:], valAddrBytes[:])
		noLongerBonded[index] = valAddr
		index++
	}
	// sorted by address - order doesn't matter
	sort.SliceStable(noLongerBonded, func(i, j int) bool {
		// -1 means strictly less than
		return bytes.Compare(noLongerBonded[i], noLongerBonded[j]) == -1
	})
	return noLongerBonded
}

// send a validator to jail
func (k Keeper) jailValidator(ctx sdk.Context, validator types.Validator) error {
	if validator.Jailed {
		return fmt.Errorf("cannot jail already jailed validator, validator: %v\n", validator)
	}

	validator.Jailed = true
	err := k.SetValidator(ctx, validator)
	if err != nil {
		return err
	}
	k.DeleteValidatorByPowerIndex(ctx, validator)
	return nil
}

// remove a validator from jail
func (k Keeper) unjailValidator(ctx sdk.Context, validator types.Validator) error {
	if !validator.Jailed {
		return fmt.Errorf("cannot unjail already unjailed validator, validator: %v\n", validator)
	}

	validator.Jailed = false
	err := k.SetValidator(ctx, validator)
	if err != nil {
		return err
	}
	err = k.SetValidatorByPowerIndex(ctx, validator)
	if err != nil {
		return err
	}
	return nil
}
