package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"bitbucket.org/decimalteam/go-node/utils/formulas"
	"bitbucket.org/decimalteam/go-node/x/validator/internal/types"
)

var compensationKey = []byte("compensations/")

func (k *Keeper) Compensate2211(ctx sdk.Context) {

	// Ensure wrong slashes are not yet compensated
	store := ctx.KVStore(k.storeKey)
	key := append(compensationKey, []byte("2211")...)
	if store.Has(key) {
		return
	}

	k.compensateDelegation(ctx, "dxvaloper1m23hp05spzs0kzwlgfyqk3gfpdt295sxzx292n", "dx1m23hp05spzs0kzwlgfyqk3gfpdt295sx75dglc", "10000000000000000000", "del")
	k.compensateDelegation(ctx, "dxvaloper1m23hp05spzs0kzwlgfyqk3gfpdt295sxzx292n", "dx1mlr92jdlgp0g6wzxz835tlzmqchy5lptw89l8j", "1000000000000000000", "del")
	k.compensateDelegation(ctx, "dxvaloper1m23hp05spzs0kzwlgfyqk3gfpdt295sxzx292n", "dx1mlr92jdlgp0g6wzxz835tlzmqchy5lptw89l8j", "1000000000000000000", "testslash1")
	k.compensateDelegation(ctx, "dxvaloper1m23hp05spzs0kzwlgfyqk3gfpdt295sxzx292n", "dx1mlr92jdlgp0g6wzxz835tlzmqchy5lptw89l8j", "1000000000000000000", "testslash2")
	k.compensateDelegation(ctx, "dxvaloper1m23hp05spzs0kzwlgfyqk3gfpdt295sxzx292n", "dx1mlr92jdlgp0g6wzxz835tlzmqchy5lptw89l8j", "1000000000000000000", "testslash3")

	k.compensateDelegationNFT(ctx, "dxvaloper1m23hp05spzs0kzwlgfyqk3gfpdt295sxzx292n", "dx1mlr92jdlgp0g6wzxz835tlzmqchy5lptw89l8j", "10000000000000000", "03bc199ef176dfb28fe4c5cbe0acd6b9f46e7554", "slashnftsnd", 1)
	k.compensateDelegationNFT(ctx, "dxvaloper1m23hp05spzs0kzwlgfyqk3gfpdt295sxzx292n", "dx1mlr92jdlgp0g6wzxz835tlzmqchy5lptw89l8j", "10000000000000000", "03bc199ef176dfb28fe4c5cbe0acd6b9f46e7554", "slashnftsnd", 2)
	k.compensateDelegationNFT(ctx, "dxvaloper1m23hp05spzs0kzwlgfyqk3gfpdt295sxzx292n", "dx1mlr92jdlgp0g6wzxz835tlzmqchy5lptw89l8j", "10000000000000000", "03bc199ef176dfb28fe4c5cbe0acd6b9f46e7554", "slashnftsnd", 3)
	k.compensateDelegationNFT(ctx, "dxvaloper1m23hp05spzs0kzwlgfyqk3gfpdt295sxzx292n", "dx1mlr92jdlgp0g6wzxz835tlzmqchy5lptw89l8j", "10000000000000000", "03bc199ef176dfb28fe4c5cbe0acd6b9f46e7554", "slashnftsnd", 4)
	k.compensateDelegationNFT(ctx, "dxvaloper1m23hp05spzs0kzwlgfyqk3gfpdt295sxzx292n", "dx1mlr92jdlgp0g6wzxz835tlzmqchy5lptw89l8j", "10000000000000000", "03bc199ef176dfb28fe4c5cbe0acd6b9f46e7554", "slashnftsnd", 5)
	k.compensateDelegationNFT(ctx, "dxvaloper1m23hp05spzs0kzwlgfyqk3gfpdt295sxzx292n", "dx1mlr92jdlgp0g6wzxz835tlzmqchy5lptw89l8j", "10000000000000000", "0a969d518cec5ee78bf3818eeb03817202481c63", "slashnftfst", 1)
	k.compensateDelegationNFT(ctx, "dxvaloper1m23hp05spzs0kzwlgfyqk3gfpdt295sxzx292n", "dx1mlr92jdlgp0g6wzxz835tlzmqchy5lptw89l8j", "20000000000000000", "958f3ece3d7df71bd1a447e8d794fc69bb44c06d", "slashnfttrd", 3)
	k.compensateDelegationNFT(ctx, "dxvaloper1m23hp05spzs0kzwlgfyqk3gfpdt295sxzx292n", "dx1mlr92jdlgp0g6wzxz835tlzmqchy5lptw89l8j", "30000000000000000", "fdf080ab9aa18fb007af07c34fc1fb729429a011", "slashnftfth", 2)
	k.compensateDelegationNFT(ctx, "dxvaloper1m23hp05spzs0kzwlgfyqk3gfpdt295sxzx292n", "dx1mlr92jdlgp0g6wzxz835tlzmqchy5lptw89l8j", "30000000000000000", "fdf080ab9aa18fb007af07c34fc1fb729429a011", "slashnftfth", 3)

	// Store record to the store to mark compensation done
	store.Set(key, []byte{1})
}

// addCoinsToAccount adds specified amount of the coin to the account.
func (k *Keeper) addCoinsToAccount(ctx sdk.Context, address sdk.AccAddress, coin sdk.Coin) {

	// Get account instance
	acc := k.AccountKeeper.GetAccount(ctx, address)
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
	cc, err := k.CoinKeeper.GetCoin(ctx, coin.Denom)
	if err != nil {
		panic(err)
	}
	volume := cc.Volume.Add(coin.Amount)
	reserve := cc.Reserve
	if !cc.IsBase() {
		ret := formulas.CalculatePurchaseAmount(cc.Volume, cc.Reserve, cc.CRR, coin.Amount)
		volume = cc.Volume.Add(coin.Amount)
		reserve = cc.Reserve.Add(ret)
	}
	k.CoinKeeper.UpdateCoin(ctx, cc, reserve, volume)
}

// compensateDelegation mints specified coin to the delegator and delegates it to specified validator.
func (k *Keeper) compensateDelegation(ctx sdk.Context, v string, d string, a string, denom string) {
	validator, _ := sdk.ValAddressFromBech32(v)
	delegator, _ := sdk.AccAddressFromBech32(d)
	amount, _ := sdk.NewIntFromString(a)
	coin := sdk.NewCoin(denom, amount)

	// Compensate slash to the account firstly
	k.addCoinsToAccount(ctx, delegator, coin)

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
	k.addCoinsToAccount(ctx, delegator, coin)

	// Update NFT sub token
	reserve, found := k.nftKeeper.GetSubToken(ctx, denom, tokenID, subTokenID)
	if !found {
		panic(fmt.Errorf("subToken with ID = %d not found", subTokenID))
	}

	reserve = reserve.Add(amount)
	k.nftKeeper.SetSubToken(ctx, denom, tokenID, subTokenID, reserve)

	// Update NFT delegation
	delegation, found := k.GetDelegationNFT(ctx, validator, delegator, tokenID, denom)
	if found {
		delegation.Coin.Amount = delegation.Coin.Amount.Add(amount)
		k.SetDelegationNFT(ctx, delegation)

		// Send compensated coins to the reserve pool
		err := k.supplyKeeper.SendCoinsFromAccountToModule(ctx, delegator, "reserved_pool", sdk.NewCoins(coin))
		if err != nil {
			panic(fmt.Errorf("insufficient funds. required: %s", amount))
		}
	}

}
