package validator

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// EndBlocker called every block, process inflation, update validator set.
func EndBlocker(ctx sdk.Context, k Keeper) []abci.ValidatorUpdate {
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

	return validatorUpdates
}
