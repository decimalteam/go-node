package keeper

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"strconv"
	"time"

	nftTypes "bitbucket.org/decimalteam/go-node/x/nft"
	"bitbucket.org/decimalteam/go-node/x/validator/internal/types"
)

func (k Keeper) GetDelegationNFT(ctx sdk.Context, valAddr sdk.ValAddress, delAddr sdk.AccAddress, tokenID, denom string) (types.DelegationNFT, bool) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetDelegationNFTKey(delAddr, valAddr, tokenID, denom)
	value := store.Get(key)
	if value == nil {
		return types.DelegationNFT{}, false
	}

	delegation := types.MustUnmarshalDelegationNFT(k.cdc, value)
	return delegation, true
}

// set a delegation
func (k Keeper) SetDelegationNFT(ctx sdk.Context, delegation types.DelegationNFT) {
	err := k.set(ctx, types.GetDelegationNFTKey(delegation.DelegatorAddress, delegation.ValidatorAddress, delegation.TokenID, delegation.Denom), delegation)
	if err != nil {
		panic(err)
	}
}

func (k Keeper) RemoveDelegationNFT(ctx sdk.Context, delegation types.DelegationNFT) {
	k.delete(ctx, types.GetDelegationNFTKey(delegation.DelegatorAddress, delegation.ValidatorAddress, delegation.TokenID, delegation.Denom))
}

func (k Keeper) SetUnbondingDelegationNFTEntry(ctx sdk.Context,
	delegatorAddr sdk.AccAddress, validatorAddr sdk.ValAddress,
	creationHeight int64, minTime time.Time, tokenID, denom string, subTokenIDs []int64) types.UnbondingDelegation {

	balance := sdk.NewCoin(k.BondDenom(ctx), sdk.ZeroInt())

	for _, id := range subTokenIDs {
		subToken, found := k.nftKeeper.GetSubToken(ctx, denom, tokenID, id)
		if !found {
			panic(fmt.Sprintf("subToken with ID = %d not found", id))
		}
		balance.Amount = balance.Amount.Add(subToken)
	}

	ubd, found := k.GetUnbondingDelegation(ctx, delegatorAddr, validatorAddr)
	if found {
		ubd.AddNFTEntry(creationHeight, minTime, tokenID, denom, subTokenIDs, balance)
	} else {
		ubd = types.NewUnbondingDelegation(delegatorAddr, validatorAddr,
			types.NewUnbondingDelegationNFTEntry(creationHeight, minTime, denom, tokenID, subTokenIDs, balance))
	}
	k.SetUnbondingDelegation(ctx, ubd)
	return ubd
}

func (k Keeper) DelegateNFT(ctx sdk.Context, delAddr sdk.AccAddress, tokenID, denom string, subTokenIDs []int64, validator types.Validator) error {
	nft, err := k.nftKeeper.GetNFT(ctx, denom, tokenID)
	if err != nil {
		return err
	}

	owner := nft.GetOwners().GetOwner(delAddr)
	if owner == nil {
		return fmt.Errorf("not found owner %s", delAddr.String())
	}

	subTokenIDs = nftTypes.SortedIntArray(subTokenIDs).Sort()

	for _, id := range subTokenIDs {
		if nftTypes.SortedIntArray(owner.GetSubTokenIDs()).Find(id) == -1 {
			return fmt.Errorf("the owner %s does not own the token with ID = %d", owner.GetAddress().String(), id)
		}
		owner = owner.RemoveSubTokenID(id)
	}

	nft = nft.SetOwners(nft.GetOwners().SetOwner(owner))

	collection, ok := k.nftKeeper.GetCollection(ctx, denom)
	if !ok {
		return fmt.Errorf("collection %s not found", collection)
	}

	collection, err = collection.UpdateNFT(nft)
	if err != nil {
		return err
	}

	k.nftKeeper.SetCollection(ctx, denom, collection)

	delegation, found := k.GetDelegationNFT(ctx, validator.ValAddress, delAddr, tokenID, denom)
	if !found {
		delegation = types.NewDelegationNFT(delAddr, validator.ValAddress, tokenID, denom, []int64{},
			sdk.NewCoin(k.BondDenom(ctx), sdk.ZeroInt()))
	}
	for _, id := range subTokenIDs {
		subToken, found := k.nftKeeper.GetSubToken(ctx, denom, tokenID, id)
		if !found {
			return fmt.Errorf("subToken with ID = %d not found", id)
		}
		delegation.SubTokenIDs = append(delegation.SubTokenIDs, id)
		delegation.Coin.Amount = delegation.Coin.Amount.Add(subToken)
	}

	k.SetDelegationNFT(ctx, delegation)

	k.DeleteValidatorByPowerIndex(ctx, validator)
	validator.Tokens = validator.Tokens.Add(delegation.Coin.Amount)
	err = k.SetValidator(ctx, validator)
	if err != nil {
		return err
	}
	k.SetValidatorByPowerIndexWithoutCalc(ctx, validator)

	return nil
}

func (k Keeper) UndelegateNFT(
	ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress, tokenID, denom string, subTokenIDs []int64,
) (time.Time, error) {

	_, foundErr := k.GetValidator(ctx, valAddr)
	if foundErr != nil {
		return time.Time{}, types.ErrNoDelegatorForAddress()
	}

	err := k.unbondNFT(ctx, delAddr, valAddr, tokenID, denom, subTokenIDs)
	if err != nil {
		return time.Time{}, err
	}

	completionTime := ctx.BlockHeader().Time.Add(k.UnbondingTime(ctx))
	ubd := k.SetUnbondingDelegationNFTEntry(ctx, delAddr, valAddr, ctx.BlockHeight(), completionTime, tokenID, denom, subTokenIDs)
	k.InsertUBDQueue(ctx, ubd, completionTime)

	return completionTime, nil
}

// unbond a particular delegation and perform associated store operations
func (k Keeper) unbondNFT(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress, tokenID, denom string, subTokenIDs []int64) error {
	// check if a delegation object exists in the store
	delegation, found := k.GetDelegationNFT(ctx, valAddr, delAddr, tokenID, denom)
	if !found {
		return types.ErrNoDelegatorForAddress()
	}

	// call the before-delegation-modified hook
	k.BeforeDelegationSharesModified(ctx, delAddr, valAddr)

	// ensure that we have enough shares to remove
	for _, id := range subTokenIDs {
		if nftTypes.SortedIntArray(delegation.SubTokenIDs).Find(id) == -1 {
			return types.ErrOwnerDoesNotOwnSubTokenID(
				delAddr.String(), strconv.FormatInt(int64(id), 10))
		}
	}

	// get validator
	validator, err := k.GetValidator(ctx, valAddr)
	if err != nil {
		return types.ErrNoValidatorFound()
	}

	decreasedAmount := sdk.ZeroInt()
	// subtract shares from delegation
	for _, id := range subTokenIDs {
		index := nftTypes.SortedIntArray(delegation.SubTokenIDs).Find(id)
		delegation.SubTokenIDs = append(delegation.SubTokenIDs[:index], delegation.SubTokenIDs[index+1:]...)

		subToken, found := k.nftKeeper.GetSubToken(ctx, denom, tokenID, id)
		if !found {
			return fmt.Errorf("subToken with ID = %d not found", id)
		}
		delegation.Coin.Amount = delegation.Coin.Amount.Sub(subToken)
		decreasedAmount = decreasedAmount.Add(subToken)
	}

	// remove the delegation
	if len(delegation.SubTokenIDs) == 0 {
		k.RemoveDelegationNFT(ctx, delegation)
	} else {
		k.SetDelegationNFT(ctx, delegation)
		// call the after delegation modification hook
		k.AfterDelegationModified(ctx, delegation.DelegatorAddress, delegation.ValidatorAddress)
	}

	k.DeleteValidatorByPowerIndex(ctx, validator)

	decreasedTokens := k.DecreaseValidatorTokens(ctx, validator, decreasedAmount)

	if decreasedTokens.IsZero() && validator.IsUnbonded() {
		// if not unbonded, we must instead remove validator in EndBlocker once it finishes its unbonding period
		err = k.RemoveValidator(ctx, validator.ValAddress)
		if err != nil {
			return types.ErrInternal(err.Error())
		}
	}

	return nil
}
