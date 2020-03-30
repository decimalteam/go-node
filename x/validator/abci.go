package validator

import (
	"bitbucket.org/decimalteam/go-node/x/coin"
	"bitbucket.org/decimalteam/go-node/x/validator/internal/types"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
	tmtypes "github.com/tendermint/tendermint/types"
)

// slashing begin block functionality
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
func EndBlocker(ctx sdk.Context, k Keeper, coinKeeper coin.Keeper) []abci.ValidatorUpdate {
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

	// Remove all mature unbonding delegations from the ubd queue.
	//matureUnbonds := k.DequeueAllMatureUBDQueue(ctx, ctx.BlockHeader().Time)
	//for _, dvPair := range matureUnbonds {
	//	err := k.CompleteUnbonding(ctx, dvPair.DelegatorAddress, dvPair.ValidatorAddress)
	//	if err != nil {
	//		continue
	//	}
	//
	//	ctx.EventManager().EmitEvent(
	//		sdk.NewEvent(
	//			types.EventTypeCompleteUnbonding,
	//			sdk.NewAttribute(types.AttributeKeyValidator, dvPair.ValidatorAddress.String()),
	//			sdk.NewAttribute(types.AttributeKeyDelegator, dvPair.DelegatorAddress.String()),
	//		),
	//	)
	//}

	// Remove all mature redelegations from the red queue.
	//matureRedelegations := k.DequeueAllMatureRedelegationQueue(ctx, ctx.BlockHeader().Time)
	//for _, dvvTriplet := range matureRedelegations {
	//	err := k.CompleteRedelegation(ctx, dvvTriplet.DelegatorAddress,
	//		dvvTriplet.ValidatorSrcAddress, dvvTriplet.ValidatorDstAddress)
	//	if err != nil {
	//		continue
	//	}
	//
	//	ctx.EventManager().EmitEvent(
	//		sdk.NewEvent(
	//			types.EventTypeCompleteRedelegation,
	//			sdk.NewAttribute(types.AttributeKeyDelegator, dvvTriplet.DelegatorAddress.String()),
	//			sdk.NewAttribute(types.AttributeKeySrcValidator, dvvTriplet.ValidatorSrcAddress.String()),
	//			sdk.NewAttribute(types.AttributeKeyDstValidator, dvvTriplet.ValidatorDstAddress.String()),
	//		),
	//	)
	//}

	rewards := types.GetRewardForBlock(uint64(height))
	// TODO потому что регистрозависимое
	denomCoin, err := coinKeeper.GetCoin(ctx, "tDCL")
	if err != nil {
		panic(err)
	}
	coinKeeper.UpdateCoin(ctx, denomCoin, denomCoin.Reserve, denomCoin.Volume.Add(rewards))

	if height%120 == 0 {
		err = k.PayRewards(ctx)
		if err != nil {
			panic(err)
		}
	}

	return validatorUpdates
}
