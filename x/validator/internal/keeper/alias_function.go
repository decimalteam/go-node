package keeper

import (
	"bitbucket.org/decimalteam/go-node/x/validator/exported"
	"bitbucket.org/decimalteam/go-node/x/validator/internal/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// iterate through the bonded validator set and perform the provided function
func (k Keeper) IterateBondedValidatorsByPower(ctx sdk.Context, fn func(index int64, validator exported.ValidatorI) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	maxValidators := k.MaxValidators(ctx)

	iterator := sdk.KVStoreReversePrefixIterator(store, []byte{types.ValidatorsByPowerIndexKey})
	defer iterator.Close()

	i := int64(0)
	for ; iterator.Valid() && i < int64(maxValidators); iterator.Next() {
		address := iterator.Value()
		validator, err := k.GetValidator(ctx, address)
		if err != nil {
			panic(err)
		}

		if validator.IsBonded() {
			stop := fn(i, validator) // XXX is this safe will the validator unexposed fields be able to get written to?
			if stop {
				break
			}
			i++
		}
	}
}

// iterate through all of the delegations from a delegator
func (k Keeper) IterateDelegations(ctx sdk.Context, delAddr sdk.AccAddress,
	fn func(index int64, del exported.DelegationI) (stop bool)) {

	store := ctx.KVStore(k.storeKey)
	delegatorPrefixKey := types.GetDelegationsKey(delAddr)
	iterator := sdk.KVStorePrefixIterator(store, delegatorPrefixKey) // smallest to largest
	defer iterator.Close()
	for i := int64(0); iterator.Valid(); iterator.Next() {
		del := types.MustUnmarshalDelegation(k.cdc, iterator.Value())
		stop := fn(i, del)
		if stop {
			break
		}
		i++
	}
}
