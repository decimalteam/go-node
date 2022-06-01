package keeper

import (
	"errors"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"bitbucket.org/decimalteam/go-node/utils/formulas"
	"bitbucket.org/decimalteam/go-node/x/validator/internal/types"
)

var compensationKey = []byte("compensations/")

func (k *Keeper) Compensate32076(ctx sdk.Context) {

	// Ensure wrong slashes are not yet compensated
	store := ctx.KVStore(k.storeKey)
	key := append(compensationKey, []byte("32076")...)
	if store.Has(key) {
		return
	}

	k.compensateDelegation(ctx, "dxvaloper15jcmwaagmtuytxlzhukdvzhylx4slcwctk96qa", "dx15jcmwaagmtuytxlzhukdvzhylx4slcwchyzh4k", "100000000000000000000", "del")
	k.compensateDelegation(ctx, "dxvaloper15jcmwaagmtuytxlzhukdvzhylx4slcwctk96qa", "dx1dyj4g65fhecmvg9esyyvvnxrdjxa2s2lg9tvtg", "1000000000000000000", "del")
	k.compensateDelegation(ctx, "dxvaloper15jcmwaagmtuytxlzhukdvzhylx4slcwctk96qa", "dx1dyj4g65fhecmvg9esyyvvnxrdjxa2s2lg9tvtg", "1000000000000000000", "testslash1")
	k.compensateDelegation(ctx, "dxvaloper15jcmwaagmtuytxlzhukdvzhylx4slcwctk96qa", "dx1dyj4g65fhecmvg9esyyvvnxrdjxa2s2lg9tvtg", "1000000000000000000", "testslash2")
	k.compensateDelegation(ctx, "dxvaloper15jcmwaagmtuytxlzhukdvzhylx4slcwctk96qa", "dx1dyj4g65fhecmvg9esyyvvnxrdjxa2s2lg9tvtg", "1000000000000000000", "testslash3")

	k.compensateDelegationNFT(ctx, "dxvaloper15jcmwaagmtuytxlzhukdvzhylx4slcwctk96qa", "dx1dyj4g65fhecmvg9esyyvvnxrdjxa2s2lg9tvtg", "20000000000000000", "2dbde2e22c8a30fe91aa93e12509f34942a068f3", "newtestthree", 1)
	k.compensateDelegationNFT(ctx, "dxvaloper15jcmwaagmtuytxlzhukdvzhylx4slcwctk96qa", "dx1dyj4g65fhecmvg9esyyvvnxrdjxa2s2lg9tvtg", "10000000000000000", "b937b0f77654d0c24bb29c49afc563c4e8009751", "newtesttwo", 1)
	k.compensateDelegationNFT(ctx, "dxvaloper15jcmwaagmtuytxlzhukdvzhylx4slcwctk96qa", "dx1dyj4g65fhecmvg9esyyvvnxrdjxa2s2lg9tvtg", "10000000000000000", "b937b0f77654d0c24bb29c49afc563c4e8009751", "newtesttwo", 2)
	k.compensateDelegationNFT(ctx, "dxvaloper15jcmwaagmtuytxlzhukdvzhylx4slcwctk96qa", "dx1dyj4g65fhecmvg9esyyvvnxrdjxa2s2lg9tvtg", "10000000000000000", "b937b0f77654d0c24bb29c49afc563c4e8009751", "newtesttwo", 3)
	k.compensateDelegationNFT(ctx, "dxvaloper15jcmwaagmtuytxlzhukdvzhylx4slcwctk96qa", "dx1dyj4g65fhecmvg9esyyvvnxrdjxa2s2lg9tvtg", "10000000000000000", "b937b0f77654d0c24bb29c49afc563c4e8009751", "newtesttwo", 4)
	k.compensateDelegationNFT(ctx, "dxvaloper15jcmwaagmtuytxlzhukdvzhylx4slcwctk96qa", "dx1dyj4g65fhecmvg9esyyvvnxrdjxa2s2lg9tvtg", "10000000000000000", "b937b0f77654d0c24bb29c49afc563c4e8009751", "newtesttwo", 5)
	k.compensateDelegationNFT(ctx, "dxvaloper15jcmwaagmtuytxlzhukdvzhylx4slcwctk96qa", "dx1dyj4g65fhecmvg9esyyvvnxrdjxa2s2lg9tvtg", "10000000000000000", "c5ba43fa15c8c1779e1878a96db4e99e34242512", "newtestone", 1)
	k.compensateDelegationNFT(ctx, "dxvaloper15jcmwaagmtuytxlzhukdvzhylx4slcwctk96qa", "dx1dyj4g65fhecmvg9esyyvvnxrdjxa2s2lg9tvtg", "30000000000000000", "c8fef4cdd21f5bca2c285d56f17a6894d271adb4", "newtestfour", 3)
	k.compensateDelegationNFT(ctx, "dxvaloper15jcmwaagmtuytxlzhukdvzhylx4slcwctk96qa", "dx1dyj4g65fhecmvg9esyyvvnxrdjxa2s2lg9tvtg", "30000000000000000", "c8fef4cdd21f5bca2c285d56f17a6894d271adb4", "newtestfour", 4)

	// Store record to the store to mark compensation done
	store.Set(key, []byte{1})
}

// addCoinsToAccount adds specified amount of the coin to the account.
func (k *Keeper) addCoinsToAccount(ctx sdk.Context, address sdk.AccAddress, coin sdk.Coin) error {

	// Get account instance
	acc := k.AccountKeeper.GetAccount(ctx, address)
	if acc == nil {
		return errors.New("account does not exist")
	}

	// Update account's coins
	coins := acc.GetCoins()
	coins = coins.Add(coin)
	err := acc.SetCoins(coins)
	if err != nil {
		return err
	}
	k.AccountKeeper.SetAccount(ctx, acc)

	// Update coin's supply
	supply := k.supplyKeeper.GetSupply(ctx)
	supply = supply.Inflate(sdk.NewCoins(coin))
	k.supplyKeeper.SetSupply(ctx, supply)

	// Update coin's volume and reserve
	cc, err := k.CoinKeeper.GetCoin(ctx, coin.Denom)
	if err != nil {
		return err
	}
	volume := cc.Volume.Add(coin.Amount)
	reserve := cc.Reserve
	if !cc.IsBase() {
		ret := formulas.CalculatePurchaseAmount(cc.Volume, cc.Reserve, cc.CRR, coin.Amount)
		volume = cc.Volume.Add(coin.Amount)
		reserve = cc.Reserve.Add(ret)
	}
	k.CoinKeeper.UpdateCoin(ctx, cc, reserve, volume)

	return nil
}

// compensateDelegation mints specified coin to the delegator and delegates it to specified validator.
func (k *Keeper) compensateDelegation(ctx sdk.Context, v string, d string, a string, denom string) {
	validator, _ := sdk.ValAddressFromBech32(v)
	delegator, _ := sdk.AccAddressFromBech32(d)
	amount, _ := sdk.NewIntFromString(a)
	coin := sdk.NewCoin(denom, amount)

	// Compensate slash to the account firstly
	err := k.addCoinsToAccount(ctx, delegator, coin)
	if err != nil {
		// TODO: Workaround somehow other way to do not break tests
		return
	}

	// Get validator
	val, err := k.GetValidator(ctx, validator)
	if err != nil {
		panic(err)
	}

	// Delegate this compensation back to the validator
	priceDelCustom, err := k.Delegate(ctx, delegator, coin, types.Unbonded, val, true)
	if err != nil {
		panic(err)
	}

	// Also it is important to emit delegation event
	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeDelegate,
		sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		sdk.NewAttribute(sdk.AttributeKeySender, delegator.String()),
		sdk.NewAttribute(types.AttributeKeyValidator, validator.String()),
		sdk.NewAttribute(types.AttributeKeyCoin, coin.String()),
		sdk.NewAttribute(types.AttributeDelPrice, priceDelCustom.String()),
	))
}

// compensateDelegationNFT corrects NFT reserve and delegates it to specified validator.
func (k *Keeper) compensateDelegationNFT(ctx sdk.Context, v string, d string, a string, tokenID string, denom string, subTokenID int64) {
	validator, _ := sdk.ValAddressFromBech32(v)
	delegator, _ := sdk.AccAddressFromBech32(d)
	amount, _ := sdk.NewIntFromString(a)
	coin := sdk.NewCoin("del", amount)

	// Compensate slash to the account firstly
	err := k.addCoinsToAccount(ctx, delegator, coin)
	if err != nil {
		// TODO: Workaround somehow other way to do not break tests
		return
	}

	// Update NFT sub token
	reserve, found := k.nftKeeper.GetSubToken(ctx, denom, tokenID, subTokenID)
	if !found {
		panic(fmt.Errorf("subToken with ID = %d not found", subTokenID))
	}

	reserve = reserve.Add(amount)
	k.nftKeeper.SetSubToken(ctx, denom, tokenID, subTokenID, reserve)

	// Send compensated coins to the reserve pool
	err = k.supplyKeeper.SendCoinsFromAccountToModule(ctx, delegator, "reserved_pool", sdk.NewCoins(coin))
	if err != nil {
		panic(fmt.Errorf("insufficient funds. required: %s", coin))
	}

	// Update NFT delegation
	delegation, found := k.GetDelegationNFT(ctx, validator, delegator, tokenID, denom)
	if found {
		delegation.Coin.Amount = delegation.Coin.Amount.Add(amount)
		k.SetDelegationNFT(ctx, delegation)
	} else {
		delegation, found := k.GetUnbondingDelegation(ctx, delegator, validator)
		if found {
			for _, entry := range delegation.Entries {
				done := false
				if e, ok := entry.(types.UnbondingDelegationNFTEntry); ok {
					if e.TokenID == tokenID && e.Denom == denom {
						for _, id := range e.SubTokenIDs {
							if id == subTokenID {
								e.Balance.Amount = e.Balance.Amount.Add(amount)
								done = true
								break
							}
						}
					}
				}
				if done {
					break
				}
			}
			k.SetUnbondingDelegation(ctx, delegation)
		}
	}
}
