package keeper

import (
	appTypes "bitbucket.org/decimalteam/go-node/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/supply/exported"
)

// GetBondedPool returns the bonded tokens pool's module account
func (k Keeper) GetBondedPool(ctx sdk.Context) (bondedPool exported.ModuleAccountI) {
	return k.supplyKeeper.GetModuleAccount(ctx, appTypes.BondedPoolName)
}

// GetNotBondedPool returns the not bonded tokens pool's module account
func (k Keeper) GetNotBondedPool(ctx sdk.Context) (notBondedPool exported.ModuleAccountI) {
	return k.supplyKeeper.GetModuleAccount(ctx, appTypes.NotBondedPoolName)
}

// bondedTokensToNotBonded transfers coins from the bonded to the not bonded pool within staking
func (k Keeper) bondedTokensToNotBonded(ctx sdk.Context, coins sdk.Coins) {
	err := k.supplyKeeper.SendCoinsFromModuleToModule(ctx, appTypes.BondedPoolName, appTypes.NotBondedPoolName, coins)
	if err != nil {
		panic(err)
	}
}

// notBondedTokensToBonded transfers coins from the not bonded to the bonded pool within staking
func (k Keeper) notBondedTokensToBonded(ctx sdk.Context, coins sdk.Coins) {
	err := k.supplyKeeper.SendCoinsFromModuleToModule(ctx, appTypes.NotBondedPoolName, appTypes.BondedPoolName, coins)
	if err != nil {
		panic(err)
	}
}

// burnBondedTokens removes coins from the bonded pool module account
func (k Keeper) burnBondedTokens(ctx sdk.Context, coins sdk.Coins) error {
	coinsBurn := sdk.NewCoins()
	for _, coin := range coins {
		if !coin.Amount.IsPositive() {
			continue
		}
		coinsBurn = coinsBurn.Add(sdk.NewCoins(coin)...)
	}
	err := k.supplyKeeper.BurnCoins(ctx, appTypes.BondedPoolName, coinsBurn)
	if err != nil {
		return err
	}
	return nil
}

// burnNotBondedTokens removes coins from the not bonded pool module account
func (k Keeper) burnNotBondedTokens(ctx sdk.Context, coins sdk.Coins) error {
	coinsBurn := sdk.NewCoins()
	for _, coin := range coins {
		if !coin.Amount.IsPositive() {
			continue
		}
		coinsBurn = coinsBurn.Add(sdk.NewCoins(coin)...)
	}
	err := k.supplyKeeper.BurnCoins(ctx, appTypes.NotBondedPoolName, coinsBurn)
	if err != nil {
		return err
	}
	return nil
}

// TotalBondedTokens total staking tokens supply which is bonded
func (k Keeper) TotalBondedTokens(ctx sdk.Context) sdk.Int {
	bondedPool := k.GetBondedPool(ctx)
	return bondedPool.GetCoins().AmountOf(k.BondDenom(ctx))
}

// StakingTokenSupply staking tokens from the total supply
func (k Keeper) StakingTokenSupply(ctx sdk.Context) sdk.Int {
	return k.supplyKeeper.GetSupply(ctx).GetTotal().AmountOf(k.BondDenom(ctx))
}

// BondedRatio the fraction of the staking tokens which are currently bonded
func (k Keeper) BondedRatio(ctx sdk.Context) sdk.Dec {
	bondedPool := k.GetBondedPool(ctx)

	stakeSupply := k.StakingTokenSupply(ctx)
	if stakeSupply.IsPositive() {
		return bondedPool.GetCoins().AmountOf(k.BondDenom(ctx)).ToDec().QuoInt(stakeSupply)
	}
	return sdk.ZeroDec()
}
