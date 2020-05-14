package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"bitbucket.org/decimalteam/go-node/utils/formulas"
	decsdk "bitbucket.org/decimalteam/go-node/utils/types"
	"bitbucket.org/decimalteam/go-node/x/validator/internal/types"
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
				sdk.NewAttribute("accum_rewards_validator", val.GetOperator().String()),
			),
		)

		rewardsVal := rewards.ToDec().Mul(val.Commission).TruncateInt()
		err := k.coinKeeper.UpdateBalance(ctx, k.BondDenom(ctx), rewardsVal, decsdk.AccAddress(val.ValAddress))
		if err != nil {
			return err
		}

		rewards = rewards.Sub(rewardsVal)
		remainder := rewards
		totalStake := k.TotalStake(ctx, val)
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

			err := k.coinKeeper.UpdateBalance(ctx, k.BondDenom(ctx), reward, del.DelegatorAddress)
			if err != nil {
				continue
			}
			remainder.Sub(reward)

			ctx.EventManager().EmitEvent(
				sdk.NewEvent(
					types.EventTypeProposerReward,
					sdk.NewAttribute(sdk.AttributeKeyAmount, reward.String()),
					sdk.NewAttribute(types.AttributeKeyValidator, val.GetOperator().String()),
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
