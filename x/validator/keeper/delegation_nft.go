package keeper

import (
	"bitbucket.org/decimalteam/go-node/x/validator/types"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"time"
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
	creationHeight int64, minTime time.Time, tokenID, denom string, quantity sdk.Int) types.UnbondingDelegation {

	token, err := k.nftKeeper.GetNFT(ctx, denom, tokenID)
	if err != nil {
		panic(err)
	}
	balance := sdk.NewCoin(k.BondDenom(ctx), quantity.Mul(token.GetReserve()))

	ubd, found := k.GetUnbondingDelegation(ctx, delegatorAddr, validatorAddr)
	if found {
		ubd.AddNFTEntry(creationHeight, minTime, tokenID, denom, quantity, balance)
	} else {
		ubd = types.NewUnbondingDelegation(delegatorAddr, validatorAddr,
			types.NewUnbondingDelegationNFTEntry(creationHeight, minTime, denom, tokenID, quantity, balance))
	}
	k.SetUnbondingDelegation(ctx, ubd)
	return ubd
}

func (k Keeper) DelegateNFT(ctx sdk.Context, delAddr sdk.AccAddress, tokenID, denom string, quantity sdk.Int, validator types.Validator) error {
	nft, err := k.nftKeeper.GetNFT(ctx, denom, tokenID)
	if err != nil {
		return err
	}

	owner := nft.GetOwners().GetOwner(delAddr)
	if owner == nil {
		return fmt.Errorf("not found owner %s", delAddr.String())
	}

	if quantity.GT(owner.GetQuantity()) {
		return fmt.Errorf("not enough quantity: %s < %s", owner.GetQuantity().String(), quantity.String())
	}

	owner = owner.SetQuantity(owner.GetQuantity().Sub(quantity))
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
	if found {
		delegation.Quantity = delegation.Quantity.Add(quantity)
		delegation.Coin.Amount = delegation.Quantity.Mul(nft.GetReserve())
	} else {
		delegation = types.NewDelegationNFT(delAddr, validator.ValAddress, tokenID, denom, quantity,
			sdk.NewCoin(k.BondDenom(ctx), quantity.Mul(nft.GetReserve())))
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
	ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress, tokenID, denom string, quantity sdk.Int,
) (time.Time, error) {

	_, foundErr := k.GetValidator(ctx, valAddr)
	if foundErr != nil {
		return time.Time{}, types.ErrNoDelegatorForAddress()
	}

	err := k.unbondNFT(ctx, delAddr, valAddr, tokenID, denom, quantity)
	if err != nil {
		return time.Time{}, err
	}

	completionTime := ctx.BlockHeader().Time.Add(k.UnbondingTime(ctx))
	ubd := k.SetUnbondingDelegationNFTEntry(ctx, delAddr, valAddr, ctx.BlockHeight(), completionTime, tokenID, denom, quantity)
	k.InsertUBDQueue(ctx, ubd, completionTime)

	return completionTime, nil
}

// unbond a particular delegation and perform associated store operations
func (k Keeper) unbondNFT(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress, tokenID, denom string, quantity sdk.Int) error {
	// check if a delegation object exists in the store
	delegation, found := k.GetDelegationNFT(ctx, valAddr, delAddr, tokenID, denom)
	if !found {
		return types.ErrNoDelegatorForAddress()
	}

	// call the before-delegation-modified hook
	k.BeforeDelegationSharesModified(ctx, delAddr, valAddr)

	// ensure that we have enough shares to remove
	if delegation.Quantity.LT(quantity) {
		return types.ErrNotEnoughDelegationShares(delegation.Coin.Amount.String())
	}

	// get validator
	validator, err := k.GetValidator(ctx, valAddr)
	if err != nil {
		return types.ErrNoValidatorFound()
	}

	token, err := k.nftKeeper.GetNFT(ctx, denom, tokenID)
	if err != nil {
		return types.ErrInternal(err.Error())
	}

	// subtract shares from delegation
	delegation.Quantity = delegation.Quantity.Sub(quantity)
	delegation.Coin.Amount = delegation.Quantity.Mul(token.GetReserve())

	// remove the delegation
	if delegation.Quantity.IsZero() {
		k.RemoveDelegationNFT(ctx, delegation)
	} else {
		k.SetDelegationNFT(ctx, delegation)
		// call the after delegation modification hook
		k.AfterDelegationModified(ctx, delegation.DelegatorAddress, delegation.ValidatorAddress)
	}

	k.DeleteValidatorByPowerIndex(ctx, validator)

	amountBase := quantity.Mul(token.GetReserve())
	decreasedTokens := k.DecreaseValidatorTokens(ctx, validator, amountBase)

	if decreasedTokens.IsZero() && validator.IsUnbonded() {
		// if not unbonded, we must instead remove validator in EndBlocker once it finishes its unbonding period
		err = k.RemoveValidator(ctx, validator.ValAddress)
		if err != nil {
			return types.ErrInternal(err.Error())
		}
	}

	return nil
}
