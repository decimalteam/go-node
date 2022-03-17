package keeper

import (
	"fmt"

	"bitbucket.org/decimalteam/go-node/utils/updates"
	"bitbucket.org/decimalteam/go-node/x/gov/internal/types"
	"bitbucket.org/decimalteam/go-node/x/validator/exported"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/libs/log"
)

// Keeper defines the governance module Keeper
type Keeper struct {
	// The reference to the Paramstore to get and set gov specific params
	paramSpace types.ParamSubspace

	// The SupplyKeeper to reduce the supply of the network
	supplyKeeper types.SupplyKeeper

	// The reference to the DelegationSet and ValidatorSet to get information about validators and delegators
	vk types.ValidatorKeeper

	// The (unexposed) keys used to access the stores from the Context.
	storeKey sdk.StoreKey

	// The codec codec for binary encoding/decoding.
	cdc *codec.Codec

	// Proposal router
	router types.Router

	skipUpgradeHeights map[int64]bool
}

// NewKeeper returns a governance keeper. It handles:
// - submitting governance proposals
// - depositing funds into proposals, and activating upon sufficient funds being deposited
// - users voting on proposals, with weight proportional to stake in the system
// - and tallying the result of the vote.
//
// CONTRACT: the parameter Subspace must have the param key table already initialized
func NewKeeper(
	cdc *codec.Codec, key sdk.StoreKey, paramSpace types.ParamSubspace,
	supplyKeeper types.SupplyKeeper, vk types.ValidatorKeeper, rtr types.Router,
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
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// ProposalQueues

// InsertActiveProposalQueue inserts a ProposalID into the active proposal queue at endBlock
func (keeper Keeper) InsertActiveProposalQueue(ctx sdk.Context, proposalID uint64, endBlock uint64) {
	store := ctx.KVStore(keeper.storeKey)
	bz := types.GetProposalIDBytes(proposalID)
	store.Set(types.ActiveProposalQueueKey(ctx, proposalID, endBlock), bz)
}

// RemoveFromActiveProposalQueue removes a proposalID from the Active Proposal Queue
func (keeper Keeper) RemoveFromActiveProposalQueue(ctx sdk.Context, proposalID uint64, endBlock uint64) {
	store := ctx.KVStore(keeper.storeKey)
	store.Delete(types.ActiveProposalQueueKey(ctx, proposalID, endBlock))
}

// InsertInactiveProposalQueue Inserts a ProposalID into the inactive proposal queue at endBlock
func (keeper Keeper) InsertInactiveProposalQueue(ctx sdk.Context, proposalID uint64, endBlock uint64) {
	store := ctx.KVStore(keeper.storeKey)
	bz := types.GetProposalIDBytes(proposalID)
	store.Set(types.InactiveProposalQueueKey(ctx, proposalID, endBlock), bz)
}

// RemoveFromInactiveProposalQueue removes a proposalID from the Inactive Proposal Queue
func (keeper Keeper) RemoveFromInactiveProposalQueue(ctx sdk.Context, proposalID uint64, endBlock uint64) {
	store := ctx.KVStore(keeper.storeKey)
	store.Delete(types.InactiveProposalQueueKey(ctx, proposalID, endBlock))
}

// Iterators

// IterateActiveProposalsQueue iterates over the proposals in the active proposal queue
// and performs a callback function
func (keeper Keeper) IterateActiveProposalsQueue(ctx sdk.Context, endBlock uint64, cb func(proposal types.Proposal) (stop bool)) {
	iterator := keeper.ActiveProposalQueueIterator(ctx, endBlock)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		proposalID, _ := types.SplitActiveProposalQueueKey(ctx, iterator.Key())
		proposal, found := keeper.GetProposal(ctx, proposalID)
		if !found {
			panic(fmt.Sprintf("proposal %d does not exist", proposalID))
		}

		if cb(proposal) {
			break
		}
	}
}

func (keeper Keeper) IterateAllActiveProposalsQueue(ctx sdk.Context, cb func(proposal types.Proposal) (stop bool)) {
	iterator := keeper.ActiveAllProposalQueueIterator(ctx)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		if len(iterator.Key()) != 17 {
			continue
		}
		proposalID, _ := types.SplitActiveProposalQueueKey(ctx, iterator.Key())
		proposal, found := keeper.GetProposal(ctx, proposalID)
		if !found {
			continue
			// panic(fmt.Sprintf("proposal %d does not exist", proposalID))
		}

		if cb(proposal) {
			break
		}
	}
}

// IterateInactiveProposalsQueue iterates over the proposals in the inactive proposal queue
// and performs a callback function
func (keeper Keeper) IterateInactiveProposalsQueue(ctx sdk.Context, endBlock uint64, cb func(proposal types.Proposal) (stop bool)) {
	iterator := keeper.InactiveProposalQueueIterator(ctx, endBlock)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		proposalID, _ := types.SplitInactiveProposalQueueKey(ctx, iterator.Key())
		proposal, found := keeper.GetProposal(ctx, proposalID)
		if !found {
			panic(fmt.Sprintf("proposal %d does not exist", proposalID))
		}

		if cb(proposal) {
			break
		}
	}
}

func (keeper Keeper) IterateAllInactiveProposalsQueue(ctx sdk.Context, cb func(proposal types.Proposal) (stop bool)) {
	iterator := keeper.InactiveAllProposalQueueIterator(ctx)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		if len(iterator.Key()) != 17 {
			continue
		}
		proposalID, _ := types.SplitInactiveProposalQueueKey(ctx, iterator.Key())
		proposal, found := keeper.GetProposal(ctx, proposalID)
		if !found {
			continue
			// panic(fmt.Sprintf("proposal %d does not exist", proposalID))
		}

		if cb(proposal) {
			break
		}
	}
}

// ActiveProposalQueueIterator returns an sdk.Iterator for all the proposals in the Active Queue that expire by endTime
func (keeper Keeper) ActiveProposalQueueIterator(ctx sdk.Context, endBlock uint64) sdk.Iterator {
	store := ctx.KVStore(keeper.storeKey)
	keyPrefix := types.ActiveProposalQueuePrefix
	if ctx.BlockHeight() < updates.Update14Block {
		keyPrefix = types.LegacyActiveProposalQueuePrefix
	}
	return store.Iterator(keyPrefix, sdk.PrefixEndBytes(types.ActiveProposalByTimeKey(ctx, endBlock)))
}

// ActiveAllProposalQueueIterator returns an sdk.Iterator for all the proposals in the Active Queue
func (keeper Keeper) ActiveAllProposalQueueIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(keeper.storeKey)
	keyPrefix := types.ActiveProposalQueuePrefix
	if ctx.BlockHeight() < updates.Update14Block {
		keyPrefix = types.LegacyActiveProposalQueuePrefix
	}
	return store.Iterator(keyPrefix, nil)
}

// InactiveProposalQueueIterator returns an sdk.Iterator for all the proposals in the Inactive Queue that expire by endTime
func (keeper Keeper) InactiveProposalQueueIterator(ctx sdk.Context, endBlock uint64) sdk.Iterator {
	store := ctx.KVStore(keeper.storeKey)
	keyPrefix := types.InactiveProposalQueuePrefix
	if ctx.BlockHeight() < updates.Update14Block {
		keyPrefix = types.LegacyInactiveProposalQueuePrefix
	}
	return store.Iterator(keyPrefix, sdk.PrefixEndBytes(types.InactiveProposalByTimeKey(ctx, endBlock)))
}

// InactiveAllProposalQueueIterator returns an sdk.Iterator for all the proposals in the Inactive Queue
func (keeper Keeper) InactiveAllProposalQueueIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(keeper.storeKey)
	keyPrefix := types.InactiveProposalQueuePrefix
	if ctx.BlockHeight() < updates.Update14Block {
		keyPrefix = types.LegacyInactiveProposalQueuePrefix
	}
	return store.Iterator(keyPrefix, nil)
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

func (k Keeper) Get(ctx sdk.Context, key []byte, value *int64) error {
	store := ctx.KVStore(k.storeKey)
	err := k.cdc.UnmarshalBinaryLengthPrefixed(store.Get(key), value)
	if err != nil {
		return err
	}
	return nil
}

func (k Keeper) Set(ctx sdk.Context, key []byte, value *int64) error {
	store := ctx.KVStore(k.storeKey)
	bz, err := k.cdc.MarshalBinaryLengthPrefixed(value)
	if err != nil {
		return err
	}
	store.Set(key, bz)
	return nil
}

func (k Keeper) MigrateToUpdatedPrefixes(ctx sdk.Context) error {
	if ctx.BlockHeight() != updates.Update14Block {
		panic(fmt.Sprintf("wrong time for data migration (called at block %d instead of %d)", ctx.BlockHeight(), updates.Update14Block))
	}
	k.migrateProposals(ctx)
	k.migrateActiveProposals(ctx)
	k.migrateInactiveProposals(ctx)
	k.migrateVotes(ctx)
	k.migrateSingleRecord(ctx, types.LegacyProposalIDKey, types.ProposalIDKey)
	k.migrateSingleRecord(ctx, types.LegacyPlanPrefix, types.PlanKey(ctx))
	k.migrateSingleRecord(ctx, types.LegacyDonePrefix, types.DoneKey(ctx))
	return nil
}

func (k Keeper) migrateProposals(ctx sdk.Context) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.LegacyProposalsKeyPrefix)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		keyFrom, value := iterator.Key(), iterator.Value()
		if len(keyFrom) != 9 { // previous key format: 0x00<proposalID_Bytes> (1+8)
			continue
		}
		var proposal types.Proposal
		k.cdc.MustUnmarshalBinaryLengthPrefixed(value, &proposal)
		keyTo := types.ProposalKey(ctx, proposal.ProposalID)
		store.Set(keyTo, value)
		store.Delete(keyFrom)
	}
}

func (k Keeper) migrateActiveProposals(ctx sdk.Context) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.LegacyActiveProposalQueuePrefix)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		keyFrom, value := iterator.Key(), iterator.Value()
		if len(keyFrom) != 17 { // previous key format: 0x01<endTime_Bytes><proposalID_Bytes> (1+8+8)
			continue
		}
		var proposal types.Proposal
		k.cdc.MustUnmarshalBinaryLengthPrefixed(value, &proposal)
		keyTo := types.ActiveProposalQueueKey(ctx, proposal.ProposalID, proposal.VotingEndBlock)
		store.Set(keyTo, value)
		store.Delete(keyFrom)
	}
}

func (k Keeper) migrateInactiveProposals(ctx sdk.Context) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.LegacyInactiveProposalQueuePrefix)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		keyFrom, value := iterator.Key(), iterator.Value()
		if len(keyFrom) != 17 { // previous key format: 0x02<endTime_Bytes><proposalID_Bytes> (1+8+8)
			continue
		}
		var proposal types.Proposal
		k.cdc.MustUnmarshalBinaryLengthPrefixed(value, &proposal)
		keyTo := types.InactiveProposalQueueKey(ctx, proposal.ProposalID, proposal.VotingEndBlock)
		store.Set(keyTo, value)
		store.Delete(keyFrom)
	}
}

func (k Keeper) migrateVotes(ctx sdk.Context) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.LegacyVotesKeyPrefix)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		keyFrom, value := iterator.Key(), iterator.Value()
		if len(keyFrom) != 29 { // previous key format: 0x10<proposalID_Bytes><voterAddr_Bytes> (1+8+20)
			continue
		}
		var vote types.Vote
		k.cdc.MustUnmarshalBinaryLengthPrefixed(value, &vote)
		keyTo := types.VoteKey(ctx, vote.ProposalID, vote.Voter)
		store.Set(keyTo, value)
		store.Delete(keyFrom)
	}
}

func (k Keeper) migrateSingleRecord(ctx sdk.Context, keyFrom []byte, keyTo []byte) {
	store := ctx.KVStore(k.storeKey)
	value := store.Get(keyFrom)
	if value != nil {
		store.Set(keyTo, value)
		store.Delete(keyFrom)
	}
}
