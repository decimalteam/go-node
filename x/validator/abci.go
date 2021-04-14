package validator

import (
	"bitbucket.org/decimalteam/go-node/utils/formulas"
	"bitbucket.org/decimalteam/go-node/x/coin"
	"bitbucket.org/decimalteam/go-node/x/validator/internal/types"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/supply"
	abci "github.com/tendermint/tendermint/abci/types"
	tmtypes "github.com/tendermint/tendermint/types"
)

// BeginBlocker check for infraction evidence or downtime of validators
// on every begin block
func BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock, k Keeper) {
	// Iterate over all the validators which *should* have signed this block
	// store whether or not they have actually signed it and slash/unbond any
	// which have missed too many blocks in a row (downtime slashing)
	for _, voteInfo := range req.LastCommitInfo.GetVotes() {
		k.HandleValidatorSignature(ctx, voteInfo.Validator.Address, voteInfo.Validator.Power, voteInfo.SignedLastBlock)
	}

	// Iterate through any newly discovered evidence of infraction
	// Slash any validators (and since-unbonded stake within the unbonding period)
	// who contributed to valid infractions
	for _, evidence := range req.ByzantineValidators {
		switch evidence.Type {
		case tmtypes.ABCIEvidenceTypeDuplicateVote:
			k.HandleDoubleSign(ctx, evidence.Validator.Address, evidence.Height, evidence.Time, evidence.Validator.Power)
		default:
			k.Logger(ctx).Error(fmt.Sprintf("ignored unknown evidence type: %s", evidence.Type))
		}
	}
}

// EndBlocker called every block, process inflation, update validator set.
func EndBlocker(ctx sdk.Context, k Keeper, coinKeeper coin.Keeper, supplyKeeper supply.Keeper, withRewards bool) []abci.ValidatorUpdate {
	if ctx.BlockHeight() == 2_335_675 {
		SyncPools(ctx, k, supplyKeeper)
	}

	// Calculate validator set changes.
	//
	// NOTE: ApplyAndReturnValidatorSetUpdates has to come before
	// UnbondAllMatureValidatorQueue.
	// This fixes a bug when the unbonding period is instant (is the case in
	// some of the tests). The test expected the validator to be completely
	// unbonded after the Endblocker (go from Bonded -> Unbonding during
	// ApplyAndReturnValidatorSetUpdates and then Unbonding -> Unbonded during
	// UnbondAllMatureValidatorQueue).
	validatorUpdates, err := k.ApplyAndReturnValidatorSetUpdates(ctx)
	if err != nil {
		panic(err)
	}

	height := ctx.BlockHeight()

	// Unbond all mature validators from the unbonding queue.
	k.UnbondAllMatureValidatorQueue(ctx)

	//Remove all mature unbonding delegations from the ubd queue.
	matureUnbonds := k.DequeueAllMatureUBDQueue(ctx, ctx.BlockHeader().Time)
	for _, dvPair := range matureUnbonds {
		delegation, found := k.GetUnbondingDelegation(ctx, dvPair.DelegatorAddress, dvPair.ValidatorAddress)
		if !found {
			continue
		}
		err := k.CompleteUnbonding(ctx, dvPair.DelegatorAddress, dvPair.ValidatorAddress)
		if err != nil {
			continue
		}

		ctxTime := ctx.BlockHeader().Time

		for _, entry := range delegation.Entries {
			if entry.IsMature(ctxTime) {
				ctx.EventManager().EmitEvent(
					sdk.NewEvent(
						types.EventTypeCompleteUnbonding,
						sdk.NewAttribute(types.AttributeKeyValidator, dvPair.ValidatorAddress.String()),
						sdk.NewAttribute(types.AttributeKeyDelegator, dvPair.DelegatorAddress.String()),
						sdk.NewAttribute(types.AttributeKeyCoin, entry.Balance.String()),
					),
				)
			}
		}
	}

	rewards := types.GetRewardForBlock(uint64(height))
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeEmission,
			sdk.NewAttribute(sdk.AttributeKeyAmount, rewards.String()),
		),
	)
	denomCoin, err := k.GetCoin(ctx, k.BondDenom(ctx))
	if err != nil {
		panic(err)
	}
	if ctx.BlockHeight() <= 79350 {
		coinKeeper.UpdateCoin(ctx, denomCoin, denomCoin.Reserve, denomCoin.Volume.Add(rewards))
	}

	feeCollector := supplyKeeper.GetModuleAccount(ctx, k.FeeCollectorName)
	feesCollectedInt := feeCollector.GetCoins()
	for _, fee := range feesCollectedInt {
		if fee.Denom == k.BondDenom(ctx) {
			rewards = rewards.Add(fee.Amount)
		} else {
			feeCoin, err := k.GetCoin(ctx, fee.Denom)
			if err != nil {
				panic(err)
			}
			feeInBaseCoin := formulas.CalculateSaleReturn(feeCoin.Volume, feeCoin.Reserve, feeCoin.CRR, fee.Amount)
			rewards = rewards.Add(feeInBaseCoin)
		}
	}
	err = supplyKeeper.BurnCoins(ctx, k.FeeCollectorName, feesCollectedInt)
	if err != nil {
		panic(err)
	}
	if ctx.BlockHeight() > 79350 {
		coinKeeper.UpdateCoin(ctx, denomCoin, denomCoin.Reserve, denomCoin.Volume.Add(rewards))
	}

	remainder := sdk.NewIntFromBigInt(rewards.BigInt())

	vals := k.GetAllValidatorsByPowerIndex(ctx)

	totalPower := sdk.ZeroInt()

	for _, val := range vals {
		totalPower = totalPower.Add(val.Tokens)
	}

	for _, val := range vals {
		if val.Tokens.IsZero() || !val.Online {
			continue
		}
		r := sdk.ZeroInt()
		r = rewards.Mul(val.Tokens).Quo(totalPower)
		remainder = remainder.Sub(r)
		val = val.AddAccumReward(r)
		err = k.SetValidator(ctx, val)
		if err != nil {
			panic(err)
		}
	}

	err = supplyKeeper.MintCoins(ctx, k.FeeCollectorName, sdk.NewCoins(sdk.NewCoin(k.BondDenom(ctx), remainder)))
	if err != nil {
		panic(err)
	}

	if height%120 == 0 && withRewards {
		err = k.PayRewards(ctx)
		if err != nil {
			panic(err)
		}
	}

	return validatorUpdates
}

func SyncPools(ctx sdk.Context, k Keeper, supplyKeeper supply.Keeper) {
	bondedTokens, notBondedTokens := sdk.NewCoins(), sdk.NewCoins()

	validators := k.GetAllValidators(ctx)
	for _, val := range validators {
		delegations := k.GetValidatorDelegations(ctx, val.ValAddress)
		for _, delegation := range delegations {
			if val.Status == Bonded {
				bondedTokens = bondedTokens.Add(delegation.Coin)
			} else {
				notBondedTokens = notBondedTokens.Add(delegation.Coin)
			}
		}
	}

	bondedPool := supplyKeeper.GetModuleAccount(ctx, BondedPoolName)

	err := bondedPool.SetCoins(bondedTokens)
	if err != nil {
		panic(err)
	}

	supplyKeeper.SetModuleAccount(ctx, bondedPool)

	notBondedPool := supplyKeeper.GetModuleAccount(ctx, NotBondedPoolName)

	err = bondedPool.SetCoins(bondedTokens)
	if err != nil {
		panic(err)
	}

	supplyKeeper.SetModuleAccount(ctx, notBondedPool)
}
