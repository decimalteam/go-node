package keeper

import (
	types2 "bitbucket.org/decimalteam/go-node/x/gov/types"
	"bitbucket.org/decimalteam/go-node/x/validator/exported"
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/libs/log"
)

// Keeper defines the governance module Keeper
type Keeper struct {
	// The reference to the Paramstore to get and set gov specific params
	paramSpace types2.ParamSubspace

	// The SupplyKeeper to reduce the supply of the network
	supplyKeeper types2.SupplyKeeper

	// The reference to the DelegationSet and ValidatorSet to get information about validators and delegators
	vk types2.ValidatorKeeper

	// The (unexposed) keys used to access the stores from the Context.
	storeKey sdk.StoreKey

	// The codec codec for binary encoding/decoding.
	cdc *codec.LegacyAmino

	// Proposal router
	router types2.Router
}

// NewKeeper returns a governance keeper. It handles:
// - submitting governance proposals
// - depositing funds into proposals, and activating upon sufficient funds being deposited
// - users voting on proposals, with weight proportional to stake in the system
// - and tallying the result of the vote.
//
// CONTRACT: the parameter Subspace must have the param key table already initialized
func NewKeeper(
	cdc *codec.LegacyAmino, key sdk.StoreKey, paramSpace types2.ParamSubspace,
	supplyKeeper types2.SupplyKeeper, vk types2.ValidatorKeeper, rtr types2.Router,
) Keeper {

	// It is vital to seal the governance proposal router here as to not allow
	// further handlers to be registered after the keeper is created since this
	// could create invalid or non-deterministic behavior.
	rtr.Seal()

	return Keeper{
		storeKey:     key,
		paramSpace:   paramSpace,
		supplyKeeper: supplyKeeper,
		vk:           vk,
		cdc:          cdc,
		router:       rtr,
	}
}

// Logger returns a module-specific logger.
func (keeper Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types2.ModuleName))
}

// ProposalQueues

// InsertActiveProposalQueue inserts a ProposalID into the active proposal queue at endBlock
func (keeper Keeper) InsertActiveProposalQueue(ctx sdk.Context, proposalID uint64, endBlock uint64) {
	store := ctx.KVStore(keeper.storeKey)
	bz := types2.GetProposalIDBytes(proposalID)
	store.Set(types2.ActiveProposalQueueKey(proposalID, endBlock), bz)
}

// RemoveFromActiveProposalQueue removes a proposalID from the Active Proposal Queue
func (keeper Keeper) RemoveFromActiveProposalQueue(ctx sdk.Context, proposalID uint64, endBlock uint64) {
	store := ctx.KVStore(keeper.storeKey)
	store.Delete(types2.ActiveProposalQueueKey(proposalID, endBlock))
}

// InsertInactiveProposalQueue Inserts a ProposalID into the inactive proposal queue at endBlock
func (keeper Keeper) InsertInactiveProposalQueue(ctx sdk.Context, proposalID uint64, endBlock uint64) {
	store := ctx.KVStore(keeper.storeKey)
	bz := types2.GetProposalIDBytes(proposalID)
	store.Set(types2.InactiveProposalQueueKey(proposalID, endBlock), bz)
}

// RemoveFromInactiveProposalQueue removes a proposalID from the Inactive Proposal Queue
func (keeper Keeper) RemoveFromInactiveProposalQueue(ctx sdk.Context, proposalID uint64, endBlock uint64) {
	store := ctx.KVStore(keeper.storeKey)
	store.Delete(types2.InactiveProposalQueueKey(proposalID, endBlock))
}

// Iterators

// IterateActiveProposalsQueue iterates over the proposals in the active proposal queue
// and performs a callback function
func (keeper Keeper) IterateActiveProposalsQueue(ctx sdk.Context, endBlock uint64, cb func(proposal types2.Proposal) (stop bool)) {
	iterator := keeper.ActiveProposalQueueIterator(ctx, endBlock)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		proposalID, _ := types2.SplitActiveProposalQueueKey(iterator.Key())
		proposal, found := keeper.GetProposal(ctx, proposalID)
		if !found {
			panic(fmt.Sprintf("proposal %d does not exist", proposalID))
		}

		if cb(proposal) {
			break
		}
	}
}

func (keeper Keeper) IterateAllActiveProposalsQueue(ctx sdk.Context, cb func(proposal types2.Proposal) (stop bool)) {
	iterator := keeper.ActiveAllProposalQueueIterator(ctx)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		if len(iterator.Key()) != 17 {
			continue
		}
		proposalID, _ := types2.SplitActiveProposalQueueKey(iterator.Key())
		proposal, found := keeper.GetProposal(ctx, proposalID)
		if !found {
			panic(fmt.Sprintf("proposal %d does not exist", proposalID))
		}

		if cb(proposal) {
			break
		}
	}
}

// IterateInactiveProposalsQueue iterates over the proposals in the inactive proposal queue
// and performs a callback function
func (keeper Keeper) IterateInactiveProposalsQueue(ctx sdk.Context, endBlock uint64, cb func(proposal types2.Proposal) (stop bool)) {
	iterator := keeper.InactiveProposalQueueIterator(ctx, endBlock)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		proposalID, _ := types2.SplitInactiveProposalQueueKey(iterator.Key())
		proposal, found := keeper.GetProposal(ctx, proposalID)
		if !found {
			panic(fmt.Sprintf("proposal %d does not exist", proposalID))
		}

		if cb(proposal) {
			break
		}
	}
}

func (keeper Keeper) IterateAllInactiveProposalsQueue(ctx sdk.Context, cb func(proposal types2.Proposal) (stop bool)) {
	iterator := keeper.InactiveAllProposalQueueIterator(ctx)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		if len(iterator.Key()) != 17 {
			continue
		}
		proposalID, _ := types2.SplitInactiveProposalQueueKey(iterator.Key())
		proposal, found := keeper.GetProposal(ctx, proposalID)
		if !found {
			panic(fmt.Sprintf("proposal %d does not exist", proposalID))
		}

		if cb(proposal) {
			break
		}
	}
}

// ActiveProposalQueueIterator returns an sdk.Iterator for all the proposals in the Active Queue that expire by endTime
func (keeper Keeper) ActiveProposalQueueIterator(ctx sdk.Context, endBlock uint64) sdk.Iterator {
	store := ctx.KVStore(keeper.storeKey)
	return store.Iterator(types2.ActiveProposalQueuePrefix, sdk.PrefixEndBytes(types2.ActiveProposalByTimeKey(endBlock)))
}

// ActiveAllProposalQueueIterator returns an sdk.Iterator for all the proposals in the Active Queue
func (keeper Keeper) ActiveAllProposalQueueIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(keeper.storeKey)
	return store.Iterator(types2.ActiveProposalQueuePrefix, nil)
}

// InactiveProposalQueueIterator returns an sdk.Iterator for all the proposals in the Inactive Queue that expire by endTime
func (keeper Keeper) InactiveProposalQueueIterator(ctx sdk.Context, endBlock uint64) sdk.Iterator {
	store := ctx.KVStore(keeper.storeKey)
	return store.Iterator(types2.InactiveProposalQueuePrefix, sdk.PrefixEndBytes(types2.InactiveProposalByTimeKey(endBlock)))
}

// InactiveAllProposalQueueIterator returns an sdk.Iterator for all the proposals in the Inactive Queue
func (keeper Keeper) InactiveAllProposalQueueIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(keeper.storeKey)
	return store.Iterator(types2.InactiveProposalQueuePrefix, nil)
}

func (keeper Keeper) CheckValidator(ctx sdk.Context, address sdk.ValAddress) error {
	if !keeper.vk.HasValidator(ctx, address) {
		return fmt.Errorf("voter is not a validator")
	}

	var val exported.ValidatorI

	keeper.vk.IterateBondedValidatorsByPower(ctx, func(index int64, validator exported.ValidatorI) bool {
		if index == 10 {
			return true
		}

		if validator.GetOperator().Equals(address) {
			val = validator
			return true
		}

		return false
	})

	if val == nil {
		return fmt.Errorf("voter doesn't have enough power voting")
	}

	return nil
}
