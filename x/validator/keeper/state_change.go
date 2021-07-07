package keeper

import (
	"bitbucket.org/decimalteam/go-node/x/validator/exported"
	"bytes"
	"errors"
	"fmt"
	"runtime/debug"
	"sort"

	abci "github.com/tendermint/tendermint/abci/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"bitbucket.org/decimalteam/go-node/x/validator/types"
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
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("stacktrace from panic: %s \n%s\n", r, string(debug.Stack()))
		}
	}()

	maxValidators := k.getValidatorsCountForBlock(ctx, ctx.BlockHeight())
	totalPower := sdk.ZeroInt()
	var amtFromBondedToNotBonded, amtFromNotBondedToBonded sdk.Coins

	// Retrieve the last validator set.
	// The persistent set is updated later in this function.
	// (see LastValidatorPowerKey).
	last := k.getLastValidatorsByAddr(ctx)

	validators := k.GetAllValidatorsByPowerIndexReversed(ctx)
	for _, validator := range validators {
		if validator.Jailed {
			continue
		}
		k.DeleteValidatorByPowerIndex(ctx, validator)
		k.checkDelegations(ctx, validator)
		k.SetValidatorByPowerIndex(ctx, validator)
	}

	validators = k.GetAllValidatorsByPowerIndexReversed(ctx)

	for i := 0; i < len(validators) && i < maxValidators; i++ {
		// everything that is iterated in this loop is becoming or already a
		// part of the bonded validator set

		validator := validators[i]

		valAddr, err := sdk.ValAddressFromBech32(validator.ValAddress)
		if err != nil {
			return nil, err
		}

		if validator.Jailed {
			return nil, errors.New("ApplyAndReturnValidatorSetUpdates: should never retrieve a jailed validator from the power store")
		}

		// if we get to a zero-power validator (which we don't bond),
		// there are no more possible bonded validators
		if validator.PotentialConsensusPower() == 0 {
			break
		}

		// apply the appropriate state change if necessary
		switch {
		case validator.IsUnbonded():
			if validator.Online {
				validator, err = k.unbondedToBonded(ctx, valAddr, validator)
				if err != nil {
					return nil, fmt.Errorf("ApplyAndReturnValidatorSetUpdates: %w", err)
				}
				delegations := k.GetValidatorDelegations(ctx, valAddr)
				for _, delegation := range delegations {
					if _, ok := delegation.(types.Delegation); ok {
						amtFromNotBondedToBonded = amtFromNotBondedToBonded.Add(delegation.GetCoin())
					}
				}
			}
		case validator.IsBonded():
			// no state change
		default:
			panic("unexpected validator status")
		}

		// fetch the old power bytes
		var valAddrBytes [types.AddrLen]byte
		copy(valAddrBytes[:], validator.ValAddress[:])
		oldPowerBytes, found := last[valAddrBytes]

		// calculate the new power bytes
		newPower := validator.ConsensusPower()
		newPowerBytes := k.cdc.MustMarshalLengthPrefixed(newPower)

		// update the validator set if power has changed
		if !found || !bytes.Equal(oldPowerBytes, newPowerBytes) {
			if validator.Online {
				updates = append(updates, validator.ABCIValidatorUpdate())
			}
			addr, err := validator.GetConsPubKey()
			ctx.EventManager().EmitEvent(
				sdk.NewEvent(
					types.EventTypeUpdatesValidators,
					sdk.NewAttribute(types.AttributeKeyPubKey, addr.String()),
					sdk.NewAttribute(types.AttributeKeyPower, fmt.Sprintf("%d", validator.ConsensusPower())),
					sdk.NewAttribute(types.AttributeKeyStake, validator.Tokens.String()),
					sdk.NewAttribute(types.AttributeKeyValidatorOdCandidate, "validator"),
				),
			)

			// set validator power on lookup index
			err = k.SetLastValidatorPower(ctx, valAddr, newPower)
			if err != nil {
				return nil, fmt.Errorf("ApplyAndReturnValidatorSetUpdates: %w", err)
			}
		}

		// validator still in the validator set, so delete from the copy
		delete(last, valAddrBytes)

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

		valAddr, err := sdk.ValAddressFromBech32(validator.ValAddress)
		if err != nil {
			return nil, err
		}

		pk, err := validator.GetConsPubKey()
		if validator.Jailed {
			ctx.EventManager().EmitEvent(
				sdk.NewEvent(
					types.EventTypeUpdatesValidators,
					sdk.NewAttribute(types.AttributeKeyPubKey, pk.String()),
					sdk.NewAttribute(types.AttributeKeyPower, "0"),
					sdk.NewAttribute(types.AttributeKeyStake, validator.Tokens.String()),
					sdk.NewAttribute(types.AttributeKeyValidatorOdCandidate, "candidate"),
				),
			)
		} else {
			ctx.EventManager().EmitEvent(
				sdk.NewEvent(
					types.EventTypeUpdatesValidators,
					sdk.NewAttribute(types.AttributeKeyPubKey, pk.String()),
					sdk.NewAttribute(types.AttributeKeyPower, fmt.Sprintf("%d", validator.ConsensusPower())),
					sdk.NewAttribute(types.AttributeKeyStake, validator.Tokens.String()),
					sdk.NewAttribute(types.AttributeKeyValidatorOdCandidate, "candidate"),
				),
			)
		}

		delegations := k.GetValidatorDelegations(ctx, valAddr)
		for _, delegation := range delegations {
			if _, ok := delegation.(types.Delegation); ok {
				amtFromBondedToNotBonded = amtFromBondedToNotBonded.Add(delegation.GetCoin())
			}
		}

		if validator.Tokens.IsZero() {
			validator, err = k.bondedToUnbonding(ctx, validator)
			if err != nil {
				panic(err)
			}
		} else {
			validator = validator.UpdateStatus(types.Unbonded)
		}

		err = k.SetValidator(ctx, validator)
		if err != nil {
			return nil, fmt.Errorf("ApplyAndReturnValidatorSetUpdates: %w", err)
		}
		// delete from the bonded validator index
		k.DeleteLastValidatorPower(ctx, valAddr)

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
		//log.Println(updates)
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
type validatorsByAddr map[[types.AddrLen]byte][]byte

// get the last validator set
func (k Keeper) getLastValidatorsByAddr(ctx sdk.Context) validatorsByAddr {
	last := make(validatorsByAddr)
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, []byte{types.LastValidatorPowerKey})
	defer iterator.Close()
	// iterate over the last validator set index
	for ; iterator.Valid(); iterator.Next() {
		var valAddr [types.AddrLen]byte
		// extract the validator address from the key (prefix is 1-byte)
		copy(valAddr[:], iterator.Key()[1:])
		// power bytes is just the value
		powerBytes := iterator.Value()
		last[valAddr] = make([]byte, len(powerBytes))
		copy(last[valAddr][:], powerBytes[:])
	}
	return last
}

func (k Keeper) unbondedToBonded(ctx sdk.Context, valAddr sdk.ValAddress, validator types.Validator) (types.Validator, error) {
	if !validator.IsUnbonded() {
		return types.Validator{}, fmt.Errorf("bad state transition unbondedToBonded, validator: %v\n", validator)
	}
	return k.bondValidator(ctx, valAddr, validator)
}

func (k Keeper) unbondingToBonded(ctx sdk.Context, valAddr sdk.ValAddress, validator types.Validator) (types.Validator, error) {
	if !validator.IsUnbonding() {
		return types.Validator{}, fmt.Errorf("bad state transition unbondingToBonded, validator: %v\n", validator)
	}
	return k.bondValidator(ctx, valAddr, validator)
}

// perform all the store operations for when a validator status becomes bonded
func (k Keeper) bondValidator(ctx sdk.Context, valAddr sdk.ValAddress, validator types.Validator) (types.Validator, error) {

	// delete the validator by power index, as the key will change
	k.DeleteValidatorByPowerIndex(ctx, validator)

	// set the status
	validator = validator.UpdateStatus(types.Bonded)

	// save the now bonded validator record to the two referenced stores
	err := k.SetValidator(ctx, validator)
	if err != nil {
		return types.Validator{}, err
	}
	k.SetValidatorByPowerIndex(ctx, validator)

	// delete from queue if present
	err = k.DeleteValidatorQueue(ctx, validator)
	if err != nil {
		return types.Validator{}, err
	}

	addr, err := validator.GetConsAddr()
	k.AfterValidatorBonded(ctx, addr, valAddr)

	return validator, nil
}

func (k Keeper) bondedToUnbonding(ctx sdk.Context, validator types.Validator) (types.Validator, error) {
	if !validator.IsBonded() {
		return types.Validator{}, fmt.Errorf("bad state transition bondedToUnbonding, validator: %v\n", validator)
	}
	return k.beginUnbondingValidator(ctx, validator)
}

func (k Keeper) checkDelegations(ctx sdk.Context, validator types.Validator) {
	valAddr, err := sdk.ValAddressFromBech32(validator.ValAddress)
	if err != nil {
		panic(err)
	}

	delegations := k.GetValidatorDelegations(ctx, valAddr)
	if len(delegations) <= int(k.MaxDelegations(ctx)) {
		return

	}

	var delegationsInBase []exported.DelegationI

	for _, delegation := range delegations {
		if _, ok := k.CoinKeeper.GetCoinsCache()[delegation.GetCoin().Denom]; ok {
			delegation = delegation.SetTokensBase(k.CalcTokensBase(ctx, delegation))
		}
		delegationsInBase = append(delegationsInBase, delegation)
	}

	sort.SliceStable(delegations, func(i, j int) bool {
		amountI := delegations[i].GetTokensBase()
		amountJ := delegations[j].GetTokensBase()
		return amountI.GT(amountJ)
	})

	for i := int(k.MaxDelegations(ctx)); i < len(delegations); i++ {
		switch validator.Status {
		case types.Bonded:
			err := k.baseKeeper.UndelegateCoinsFromModuleToAccount(ctx, types.BondedPoolName, delegations[i].GetDelegatorAddr(), sdk.NewCoins(delegations[i].GetCoin()))
			if err != nil {
				panic(err)
			}
		case types.Unbonded:
			err := k.baseKeeper.UndelegateCoinsFromModuleToAccount(ctx, types.NotBondedPoolName, delegations[i].GetDelegatorAddr(), sdk.NewCoins(delegations[i].GetCoin()))
			if err != nil {
				panic(err)
			}
		}
		err := k.unbond(ctx, delegations[i].GetDelegatorAddr(), delegations[i].GetValidatorAddr(), delegations[i].GetCoin())
		if err != nil {
			panic(err)
		}
	}
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
	k.SetValidatorByPowerIndex(ctx, validator)

	valAddr, err := sdk.ValAddressFromBech32(validator.ValAddress)
	if err != nil {
		return types.Validator{}, err
	}

	// Adds to unbonding validator queue
	err = k.InsertValidatorQueue(ctx, valAddr, validator)
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
		valAddr := make([]byte, types.AddrLen)
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

	k.DeleteValidatorByPowerIndex(ctx, validator)
	validator.Jailed = true
	validator.Online = false
	err := k.SetValidator(ctx, validator)
	if err != nil {
		return err
	}
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
	k.DeleteValidatorByPowerIndex(ctx, validator)
	k.SetValidatorByPowerIndex(ctx, validator)
	return nil
}

func (k Keeper) getValidatorsCountForBlock(ctx sdk.Context, block int64) int {
	// count := 16 + (block/518400)*4
	count := 5 + (block/7200)*1
	if uint16(count) > uint16(k.GetParams(ctx).MaxValidators) {
		return int(k.GetParams(ctx).MaxValidators)
	}

	return int(count)
}
