package keeper

import (
	"fmt"

	"bitbucket.org/decimalteam/go-node/utils/formulas"
	"bitbucket.org/decimalteam/go-node/x/validator/internal/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var compensationKey = []byte("compensations/")

func (k *Keeper) Compensate342(ctx sdk.Context) {

	// Ensure wrong slashes are not yet compensated
	store := ctx.KVStore(k.storeKey)
	key := append(compensationKey, []byte("342")...)
	if store.Has(key) {
		return
	}

	k.compensateDelegation(ctx, "dxvaloper1kx6sccjfj8qtjfquv30n67e7f92mlzz4d5c9mz", "dx1kx6sccjfj8qtjfquv30n67e7f92mlzz43xlgwf", "110000000000000000000", "del")
	k.compensateDelegation(ctx, "dxvaloper1kx6sccjfj8qtjfquv30n67e7f92mlzz4d5c9mz", "dx1mlr92jdlgp0g6wzxz835tlzmqchy5lptw89l8j", "1000000000000000000", "del")
	k.compensateDelegation(ctx, "dxvaloper1kx6sccjfj8qtjfquv30n67e7f92mlzz4d5c9mz", "dx1mlr92jdlgp0g6wzxz835tlzmqchy5lptw89l8j", "1000000000000000000", "testslash1")
	k.compensateDelegation(ctx, "dxvaloper1kx6sccjfj8qtjfquv30n67e7f92mlzz4d5c9mz", "dx1mlr92jdlgp0g6wzxz835tlzmqchy5lptw89l8j", "1000000000000000000", "testslash2")
	k.compensateDelegation(ctx, "dxvaloper1kx6sccjfj8qtjfquv30n67e7f92mlzz4d5c9mz", "dx1mlr92jdlgp0g6wzxz835tlzmqchy5lptw89l8j", "1000000000000000000", "testslash3")

	k.compensateDelegationNFT(ctx, "dxvaloper1kx6sccjfj8qtjfquv30n67e7f92mlzz4d5c9mz", "dx1mlr92jdlgp0g6wzxz835tlzmqchy5lptw89l8j", "10000000000000000", "04a1ca30188cda9c636da29a29305f00e33022b0", "regarde", []int64{1})
	k.compensateDelegationNFT(ctx, "dxvaloper1kx6sccjfj8qtjfquv30n67e7f92mlzz4d5c9mz", "dx1mlr92jdlgp0g6wzxz835tlzmqchy5lptw89l8j", "10000000000000000", "79d8b41a92500974769f96af42d376c844f7a95e", "regarde", []int64{1})
	k.compensateDelegationNFT(ctx, "dxvaloper1kx6sccjfj8qtjfquv30n67e7f92mlzz4d5c9mz", "dx1mlr92jdlgp0g6wzxz835tlzmqchy5lptw89l8j", "10000000000000000", "79d8b41a92500974769f96af42d376c844f7a95e", "regarde", []int64{2})
	k.compensateDelegationNFT(ctx, "dxvaloper1kx6sccjfj8qtjfquv30n67e7f92mlzz4d5c9mz", "dx1mlr92jdlgp0g6wzxz835tlzmqchy5lptw89l8j", "10000000000000000", "79d8b41a92500974769f96af42d376c844f7a95e", "regarde", []int64{3})
	k.compensateDelegationNFT(ctx, "dxvaloper1kx6sccjfj8qtjfquv30n67e7f92mlzz4d5c9mz", "dx1mlr92jdlgp0g6wzxz835tlzmqchy5lptw89l8j", "10000000000000000", "79d8b41a92500974769f96af42d376c844f7a95e", "regarde", []int64{4})
	k.compensateDelegationNFT(ctx, "dxvaloper1kx6sccjfj8qtjfquv30n67e7f92mlzz4d5c9mz", "dx1mlr92jdlgp0g6wzxz835tlzmqchy5lptw89l8j", "10000000000000000", "79d8b41a92500974769f96af42d376c844f7a95e", "regarde", []int64{5})
	k.compensateDelegationNFT(ctx, "dxvaloper1kx6sccjfj8qtjfquv30n67e7f92mlzz4d5c9mz", "dx1mlr92jdlgp0g6wzxz835tlzmqchy5lptw89l8j", "20000000000000000", "9df1819b976fb95e45290f9a5d8281270ea77a8f", "regarde", []int64{1})
	k.compensateDelegationNFT(ctx, "dxvaloper1kx6sccjfj8qtjfquv30n67e7f92mlzz4d5c9mz", "dx1mlr92jdlgp0g6wzxz835tlzmqchy5lptw89l8j", "30000000000000000", "d353891b4daeac3de581486196282b1ae6d6d387", "regarde", []int64{1})
	k.compensateDelegationNFT(ctx, "dxvaloper1kx6sccjfj8qtjfquv30n67e7f92mlzz4d5c9mz", "dx1mlr92jdlgp0g6wzxz835tlzmqchy5lptw89l8j", "30000000000000000", "d353891b4daeac3de581486196282b1ae6d6d387", "regarde", []int64{2})

	// Store record to the store to mark compensation done
	store.Set(key, []byte{1})
}

// compensateDelegation mints specified coin to the delegator and delegates it to specified validator.
func (k *Keeper) compensateDelegation(ctx sdk.Context, v string, d string, a string, denom string) {
	validator, _ := sdk.ValAddressFromBech32(v)
	delegator, _ := sdk.AccAddressFromBech32(d)
	amount, _ := sdk.NewIntFromString(a)
	coin := sdk.NewCoin(denom, amount)

	// Get account instance
	acc := k.AccountKeeper.GetAccount(ctx, delegator)
	if acc == nil {
		panic("account does not exist")
	}

	// Update account's coins
	coins := acc.GetCoins()
	coins = coins.Add(coin)
	err := acc.SetCoins(coins)
	if err != nil {
		panic(err)
	}
	k.AccountKeeper.SetAccount(ctx, acc)

	// Update coin's volume and reserve
	cc, err := k.CoinKeeper.GetCoin(ctx, denom)
	if err != nil {
		panic(err)
	}
	volume := cc.Volume.Add(amount)
	reserve := cc.Reserve
	if !cc.IsBase() {
		ret := formulas.CalculatePurchaseAmount(cc.Volume, cc.Reserve, cc.CRR, amount)
		volume = cc.Volume.Add(amount)
		reserve = cc.Reserve.Add(ret)
	}
	k.CoinKeeper.UpdateCoin(ctx, cc, reserve, volume)

	if !cc.IsBase() {
		ctx.Logger().Info(fmt.Sprintf("% 12s price after updates: %s", denom, sdk.NewDecFromInt(volume).QuoInt(reserve).String()))
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
		sdk.NewAttribute(types.AttributeKeyCoin, denom),
		sdk.NewAttribute(types.AttributeDelPrice, priceDelCustom.String()),
	))
}

// compensateDelegationNFT corrects NFT reserve and delegates it to specified validator.
func (k *Keeper) compensateDelegationNFT(ctx sdk.Context, v string, d string, a string, tokenID string, denom string, subTokenIDs []int64) {
	validator, _ := sdk.ValAddressFromBech32(v)
	delegator, _ := sdk.AccAddressFromBech32(d)
	amount, _ := sdk.NewIntFromString(a)

	// Update NFT sub tokens
	for _, subTokenID := range subTokenIDs {
		reserve, found := k.nftKeeper.GetSubToken(ctx, denom, tokenID, subTokenID)
		if !found {
			panic(fmt.Errorf("subToken with ID = %d not found", subTokenID))
		}
		reserve = reserve.Add(amount)
		k.nftKeeper.SetSubToken(ctx, denom, tokenID, subTokenID, reserve)
	}

	// Update NFT delegation
	delegation, found := k.GetDelegationNFT(ctx, validator, delegator, tokenID, denom)
	if found {
		delegation.Coin.Amount = delegation.Coin.Amount.Add(amount)
		k.SetDelegationNFT(ctx, delegation)
	}
}
