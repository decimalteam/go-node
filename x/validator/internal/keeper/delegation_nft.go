package keeper

import (
	"bitbucket.org/decimalteam/go-node/x/validator/internal/types"
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
	//ubd := k.SetUnbondingDelegationEntry(ctx, delAddr, valAddr, ctx.BlockHeight(), completionTime, amount)
	//k.InsertUBDQueue(ctx, ubd, completionTime)

	return completionTime, nil
}

// unbond a particular delegation and perform associated store operations
func (k Keeper) unbondNFT(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress, tokenID, denom string, quantity sdk.Int) error {
	//	// check if a delegation object exists in the store
	//	delegation, found := k.GetDelegation(ctx, delAddr, valAddr, coin.Denom)
	//	if !found {
	//		return types.ErrNoDelegatorForAddress()
	//	}
	//
	//	// call the before-delegation-modified hook
	//	k.BeforeDelegationSharesModified(ctx, delAddr, valAddr)
	//
	//	// ensure that we have enough shares to remove
	//	if delegation.Coin.Amount.LT(coin.Amount) {
	//		return types.ErrNotEnoughDelegationShares(delegation.Coin.Amount.String())
	//	}
	//
	//	// get validator
	//	validator, err := k.GetValidator(ctx, valAddr)
	//	if err != nil {
	//		return types.ErrNoValidatorFound()
	//	}
	//
	//	// subtract shares from delegation
	//	delegation.Coin = delegation.Coin.Sub(coin)
	//
	//	// remove the delegation
	//	if delegation.Coin.IsZero() {
	//		k.RemoveDelegation(ctx, delegation)
	//	} else {
	//		k.SetDelegation(ctx, delegation)
	//		// call the after delegation modification hook
	//		k.AfterDelegationModified(ctx, delegation.DelegatorAddress, delegation.ValidatorAddress)
	//	}
	//
	//	k.DeleteValidatorByPowerIndex(ctx, validator)
	//
	//	amountBase := coin.Amount
	//	if coin.Denom != k.BondDenom(ctx) {
	//		c, err := k.GetCoin(ctx, coin.Denom)
	//		if err != nil {
	//			return types.ErrInternal(err.Error())
	//		}
	//		amountBase = formulas.CalculateSaleReturn(c.Volume, c.Reserve, c.CRR, coin.Amount)
	//	}
	//	decreasedTokens := k.DecreaseValidatorTokens(ctx, validator, amountBase)
	//
	//	if decreasedTokens.IsZero() && validator.IsUnbonded() {
	//		// if not unbonded, we must instead remove validator in EndBlocker once it finishes its unbonding period
	//		err = k.RemoveValidator(ctx, validator.ValAddress)
	//		if err != nil {
	//			return types.ErrInternal(err.Error())
	//		}
	//	}
	//
	return nil
}
