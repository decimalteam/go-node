package keeper

import (
	"bitbucket.org/decimalteam/go-node/utils/formulas"
	"bitbucket.org/decimalteam/go-node/x/multisig"
	"bitbucket.org/decimalteam/go-node/x/validator/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
)

// DAO     {"address":"dx1pk2rurh73er88p032qrd6kq5xmu53thjylflsr","owners":["dx18tay9ayumxjun9sexlq4t3nvt7zts5typnyjdr","dx1w54s4wq8atjmmu4snv0tt72qpvtg38megw5ngn","dx19ws36j00axpk0ytumc20l9wyv0ae26zygk2z0f"],"weights":["1","1","1"],"threshold":"3"}
// Develop {"address":"dx1gsa4w0cuyjqwt9j7qtc32m6n0lkyxfanphfaug","owners":["dx1fpjhs2wlaz6dd95d0lmxj5tfrmncwg437jh0y3","dx1lfleqkc39pt2jkyhr7m845x207kh5d9av3423z","dx1f46tyn4wmnvuxfj9cu5yn6vn939spfzt3yhxey"],"weights":["1","1","1"],"threshold":"3"}

var DAOCommission = sdk.NewDec(5).QuoInt64(100)
var DevelopCommission = sdk.NewDec(5).QuoInt64(100)

func (k Keeper) PayRewards(ctx sdk.Context) error {
	validators := k.GetAllValidators(ctx)
	for _, val := range validators {
		if val.AccumRewards.IsZero() {
			continue
		}
		validatorAddr, err := sdk.ValAddressFromBech32(val.ValAddress)
		if err != nil {
			return err
		}

		rewards := val.AccumRewards
		accumRewards := rewards

		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeProposerReward,
				sdk.NewAttribute("accum_rewards", accumRewards.String()),
				sdk.NewAttribute("accum_rewards_validator", validatorAddr.String()),
			),
		)

		daoWallet, err := k.getDAO(ctx)
		if err != nil {
			return err
		}
		developWallet, err := k.getDevelop(ctx)
		if err != nil {
			return err
		}

		daoVal := rewards.ToDec().Mul(DAOCommission).TruncateInt()

		//err = k.CoinKeeper.BankKeeper.AddCoins(ctx, daoWallet, sdk.NewCoins(sdk.NewCoin(k.BondDenom(ctx), daoVal)))
		coins := sdk.NewCoins(sdk.NewCoin(k.BondDenom(ctx), daoVal))
		if err := k.CoinKeeper.BankKeeper.MintCoins(ctx, minttypes.ModuleName, coins); err != nil {
			panic(err)
		}
		if err != nil {
			return err
		}

		err = k.CoinKeeper.BankKeeper.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, daoWallet, coins)
		if err != nil {
			return err
		}

		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeDAOReward,
				sdk.NewAttribute(sdk.AttributeKeyAmount, daoVal.String()),
				sdk.NewAttribute(types.AttributeKeyDAOAddress, daoWallet.String()),
				sdk.NewAttribute(types.AttributeKeyValidator, validatorAddr.String()),
			),
		)

		developVal := rewards.ToDec().Mul(DevelopCommission).TruncateInt()
		newCoins := sdk.NewCoins(sdk.NewCoin(k.BondDenom(ctx), developVal))
		//err = k.CoinKeeper.BankKeeper.AddCoins(ctx, developWallet, sdk.NewCoins(sdk.NewCoin(k.BondDenom(ctx), developVal)))

		if err := k.CoinKeeper.BankKeeper.MintCoins(ctx, minttypes.ModuleName, newCoins); err != nil {
			panic(err)
		}
		if err != nil {
			return err
		}

		err = k.CoinKeeper.BankKeeper.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, developWallet, newCoins)
		if err != nil {
			return err
		}

		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeDevelopReward,
				sdk.NewAttribute(sdk.AttributeKeyAmount, developVal.String()),
				sdk.NewAttribute(types.AttributeKeyDevelopAddress, developWallet.String()),
				sdk.NewAttribute(types.AttributeKeyValidator, validatorAddr.String()),
			),
		)

		rewards = rewards.Sub(daoVal)
		rewards = rewards.Sub(developVal)

		rewardsVal := rewards.ToDec().Mul(val.Commission).TruncateInt()

		rewardAddr, err := sdk.AccAddressFromBech32(val.RewardAddress)
		if err != nil {
			return err
		}

		err = k.CoinKeeper.UpdateBalance(ctx, k.BondDenom(ctx), rewardsVal, rewardAddr)
		if err != nil {
			return err
		}
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeCommissionReward,
				sdk.NewAttribute(sdk.AttributeKeyAmount, rewardsVal.String()),
				sdk.NewAttribute(types.AttributeKeyValidator, validatorAddr.String()),
				sdk.NewAttribute(types.AttributeKeyRewardAddress, rewardAddr.String()),
			),
		)

		rewards = rewards.Sub(rewardsVal)
		remainder := rewards
		totalStake := val.Tokens
		delegations := k.GetValidatorDelegations(ctx, validatorAddr)
		for _, del := range delegations {
			reward := sdk.NewIntFromBigInt(rewards.BigInt())
			if del.GetCoin().Denom != k.BondDenom(ctx) {
				coinDel, err := k.GetCoin(ctx, del.GetCoin().Denom)
				if err != nil {
					return err
				}
				defAmount := formulas.CalculateSaleReturn(coinDel.Volume, coinDel.Reserve, coinDel.CRR, del.GetCoin().Amount)

				reward = reward.Mul(defAmount).Quo(totalStake)
				if reward.LT(sdk.NewInt(1)) {
					continue
				}
			} else {
				reward = reward.Mul(del.GetCoin().Amount).Quo(totalStake)
				if reward.LT(sdk.NewInt(1)) {
					continue
				}
			}

			err := k.CoinKeeper.UpdateBalance(ctx, k.BondDenom(ctx), reward, del.GetDelegatorAddr())
			if err != nil {
				continue
			}
			remainder.Sub(reward)

			ctx.EventManager().EmitEvent(
				sdk.NewEvent(
					types.EventTypeProposerReward,
					sdk.NewAttribute(sdk.AttributeKeyAmount, reward.String()),
					sdk.NewAttribute(types.AttributeKeyValidator, validatorAddr.String()),
					sdk.NewAttribute(types.AttributeKeyDelegator, del.GetDelegatorAddr().String()),
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

const DAOAddress1 = "dx18tay9ayumxjun9sexlq4t3nvt7zts5typnyjdr"
const DAOAddress2 = "dx1w54s4wq8atjmmu4snv0tt72qpvtg38megw5ngn"
const DAOAddress3 = "dx19ws36j00axpk0ytumc20l9wyv0ae26zygk2z0f"

func (k Keeper) getDAO(ctx sdk.Context) (sdk.AccAddress, error) {
	address, err := sdk.AccAddressFromBech32("dx1pk2rurh73er88p032qrd6kq5xmu53thjylflsr")
	if err != nil {
		return nil, err
	}

	wallet := k.multisigKeeper.GetWallet(ctx, address.String())
	//if wallet.Address != nil {
	//	return address, nil
	//}

	owners := []string{
		DAOAddress1, DAOAddress2, DAOAddress3,
	}

	wallet = multisig.Wallet{
		Address:   address.String(),
		Owners:    owners,
		Weights:   []uint64{1, 1, 1},
		Threshold: 3}

	k.multisigKeeper.SetWallet(ctx, wallet)
	return address, nil
}

const DevelopAddress1 = "dx1fpjhs2wlaz6dd95d0lmxj5tfrmncwg437jh0y3"
const DevelopAddress2 = "dx1lfleqkc39pt2jkyhr7m845x207kh5d9av3423z"
const DevelopAddress3 = "dx1f46tyn4wmnvuxfj9cu5yn6vn939spfzt3yhxey"

func (k Keeper) getDevelop(ctx sdk.Context) (sdk.AccAddress, error) {
	address, err := sdk.AccAddressFromBech32("dx1gsa4w0cuyjqwt9j7qtc32m6n0lkyxfanphfaug")
	if err != nil {
		return nil, err
	}

	wallet := k.multisigKeeper.GetWallet(ctx, address.String())
	//if wallet.Address != nil {
	//	return address, nil
	//}

	owners := []string{
		DevelopAddress1, DevelopAddress2, DevelopAddress3,
	}

	wallet = multisig.Wallet{
		Address:   address.String(),
		Owners:    owners,
		Weights:   []uint64{1, 1, 1},
		Threshold: 3}

	k.multisigKeeper.SetWallet(ctx, wallet)
	return address, nil
}
