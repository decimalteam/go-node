package keeper

import (
	"bitbucket.org/decimalteam/go-node/utils/updates"
	"bitbucket.org/decimalteam/go-node/x/validator/internal/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/supply"
	"github.com/cosmos/cosmos-sdk/x/supply/exported"
)

// GetBondedPool returns the bonded tokens pool's module account
func (k Keeper) GetBondedPool(ctx sdk.Context) (bondedPool exported.ModuleAccountI) {
	return k.supplyKeeper.GetModuleAccount(ctx, types.BondedPoolName)
}

// GetNotBondedPool returns the not bonded tokens pool's module account
func (k Keeper) GetNotBondedPool(ctx sdk.Context) (notBondedPool exported.ModuleAccountI) {
	return k.supplyKeeper.GetModuleAccount(ctx, types.NotBondedPoolName)
}

// bondedTokensToNotBonded transfers coins from the bonded to the not bonded pool within staking
func (k Keeper) bondedTokensToNotBonded(ctx sdk.Context, coins sdk.Coins) {
	err := k.supplyKeeper.SendCoinsFromModuleToModule(ctx, types.BondedPoolName, types.NotBondedPoolName, coins)
	if err != nil {
		panic(err)
	}
}

// notBondedTokensToBonded transfers coins from the not bonded to the bonded pool within staking
func (k Keeper) notBondedTokensToBonded(ctx sdk.Context, coins sdk.Coins) {
	err := k.supplyKeeper.SendCoinsFromModuleToModule(ctx, types.NotBondedPoolName, types.BondedPoolName, coins)
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
	if ctx.BlockHeight() >= updates.Update1Block {
		err := k.burnCoins(ctx, types.BondedPoolName, coinsBurn)
		if err != nil {
			return err
		}
	} else {
		err := k.supplyKeeper.BurnCoins(ctx, types.BondedPoolName, coinsBurn)
		if err != nil {
			return err
		}
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
	if ctx.BlockHeight() >= updates.Update1Block {
		err := k.burnCoins(ctx, types.NotBondedPoolName, coinsBurn)
		if err != nil {
			return err
		}
	} else {
		err := k.supplyKeeper.BurnCoins(ctx, types.NotBondedPoolName, coinsBurn)
		if err != nil {
			return err
		}
	}
	return nil
}

func (k Keeper) burnCoins(ctx sdk.Context, moduleAccount string, coins sdk.Coins) error {
	acc := k.supplyKeeper.GetModuleAccount(ctx, types.BondedPoolName)
	if acc == nil {
		panic(sdkerrors.Wrapf(sdkerrors.ErrUnknownAddress, "module account %s does not exist", moduleAccount))
	}

	if !acc.HasPermission(supply.Burner) {
		panic(sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "module account %s does not have permissions to burn tokens", moduleAccount))
	}

	_, err := k.CoinKeeper.BankKeeper.SubtractCoins(ctx, acc.GetAddress(), coins)
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
