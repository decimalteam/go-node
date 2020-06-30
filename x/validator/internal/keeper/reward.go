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
		accumRewards := rewards

		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeProposerReward,
				sdk.NewAttribute("accum_rewards", accumRewards.String()),
				sdk.NewAttribute("accum_rewards_validator", val.ValAddress.String()),
			),
		)

		rewardsVal := rewards.ToDec().Mul(val.Commission).TruncateInt()
		err := k.CoinKeeper.UpdateBalance(ctx, k.BondDenom(ctx), rewardsVal, val.RewardAddress)
		if err != nil {
			return err
		}
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeCommissionReward,
				sdk.NewAttribute(sdk.AttributeKeyAmount, rewardsVal.String()),
				sdk.NewAttribute(types.AttributeKeyValidator, val.ValAddress.String()),
				sdk.NewAttribute(types.AttributeKeyRewardAddress, val.RewardAddress.String()),
			),
		)

		rewards = rewards.Sub(rewardsVal)
		remainder := rewards
		totalStake := val.Tokens
		delegations := k.GetValidatorDelegations(ctx, val.ValAddress)
		for _, del := range delegations {
			reward := sdk.NewIntFromBigInt(rewards.BigInt())
			if del.Coin.Denom != k.BondDenom(ctx) {
				coinDel, err := k.GetCoin(ctx, del.Coin.Denom)
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

			err := k.CoinKeeper.UpdateBalance(ctx, k.BondDenom(ctx), reward, del.DelegatorAddress)
			if err != nil {
				continue
			}
			remainder.Sub(reward)

			ctx.EventManager().EmitEvent(
				sdk.NewEvent(
					types.EventTypeProposerReward,
					sdk.NewAttribute(sdk.AttributeKeyAmount, reward.String()),
					sdk.NewAttribute(types.AttributeKeyValidator, val.ValAddress.String()),
					sdk.NewAttribute(types.AttributeKeyDelegator, del.DelegatorAddress.String()),
				),
			)
		}
		val.AccumRewards = sdk.ZeroInt()
		err = k.SetValidator(ctx, val)
		if err != nil {
			panic(err)
		}
	}
	return nil
}
