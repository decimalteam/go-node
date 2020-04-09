package keeper

import (
	"bitbucket.org/decimalteam/go-node/utils/formulas"
	"bitbucket.org/decimalteam/go-node/x/validator/internal/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) PayRewards(ctx sdk.Context) error {
	validators := k.GetAllValidators(ctx)
	for _, val := range validators {
		if val.AccumRewards.IsZero() {
			continue
		}
		rewards := val.AccumRewards

		rewardsVal := rewards.ToDec().Mul(val.Commission.Rate).TruncateInt()
		err := k.coinKeeper.UpdateBalance(ctx, types.DefaultBondDenom, rewardsVal, sdk.AccAddress(val.ValAddress))
		if err != nil {
			return err
		}

		rewards = rewards.Sub(rewardsVal)
		remainder := rewards
		totalStake := k.TotalStake(ctx, val)
		delegations := k.GetValidatorDelegations(ctx, val.ValAddress)
		for _, del := range delegations {
			reward := sdk.NewIntFromBigInt(rewards.BigInt())
			if del.Coin.Denom != types.DefaultBondDenom {
				coinDel, err := k.coinKeeper.GetCoin(ctx, del.Coin.Denom)
				if err != nil {
					return err
				}
				defAmount := formulas.CalculateSaleReturn(coinDel.Volume, coinDel.Reserve, coinDel.CRR, del.Coin.Amount)

				reward = reward.Mul(defAmount).Quo(totalStake)
				if reward.LT(sdk.NewInt(1)) {
					continue
				}
			} else {
				reward = reward.Mul(del.Coin.Amount).Quo(totalStake)
				if reward.LT(sdk.NewInt(1)) {
					continue
				}
			}

			err := k.coinKeeper.UpdateBalance(ctx, types.DefaultBondDenom, reward, del.DelegatorAddress)
			if err != nil {
				continue
			}
			remainder.Sub(reward)
		}
		val.AccumRewards = sdk.ZeroInt()
		err = k.SetValidator(ctx, val)
		if err != nil {
			panic(err)
		}
	}
	return nil
}
