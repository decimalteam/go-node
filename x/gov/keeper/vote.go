package keeper

import (
	types2 "bitbucket.org/decimalteam/go-node/x/gov/types"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// AddVote adds a vote on a specific proposal
func (keeper Keeper) AddVote(ctx sdk.Context, proposalID uint64, voterAddr sdk.ValAddress, option types2.VoteOption) error {
	proposal, ok := keeper.GetProposal(ctx, proposalID)
	if !ok {
		return sdkerrors.Wrapf(types2.ErrUnknownProposal, "%d", proposalID)
	}
	if proposal.Status != types2.StatusVotingPeriod {
		return sdkerrors.Wrapf(types2.ErrInactiveProposal, "%d", proposalID)
	}

	if !types2.ValidVoteOption(option) {
		return sdkerrors.Wrap(types2.ErrInvalidVote, option.String())
	}

	vote := types2.NewVote(proposalID, voterAddr, option)
	keeper.SetVote(ctx, vote)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types2.EventTypeProposalVote,
			sdk.NewAttribute(types2.AttributeKeyOption, option.String()),
			sdk.NewAttribute(types2.AttributeKeyProposalID, fmt.Sprintf("%d", proposalID)),
		),
	)

	return nil
}

// GetVote gets the vote from an address on a specific proposal
func (keeper Keeper) GetVote(ctx sdk.Context, proposalID uint64, voterAddr sdk.ValAddress) (vote types2.Vote, found bool) {
	store := ctx.KVStore(keeper.storeKey)
	bz := store.Get(types2.VoteKey(proposalID, voterAddr))
	if bz == nil {
		return vote, false
	}

	keeper.cdc.MustUnmarshalLengthPrefixed(bz, &vote)
	return vote, true
}

// SetVote sets a Vote to the gov store
func (keeper Keeper) SetVote(ctx sdk.Context, vote types2.Vote) {
	store := ctx.KVStore(keeper.storeKey)
	bz := keeper.cdc.MustMarshalLengthPrefixed(vote)
	store.Set(types2.VoteKey(vote.ProposalID, vote.Voter), bz)
}

// GetAllVotes returns all the votes from the store
func (keeper Keeper) GetAllVotes(ctx sdk.Context) (votes types2.Votes) {
	keeper.IterateAllVotes(ctx, func(vote types2.Vote) bool {
		votes = append(votes, vote)
		return false
	})
	return
}

// GetVotes returns all the votes from a proposal
func (keeper Keeper) GetVotes(ctx sdk.Context, proposalID uint64) (votes types2.Votes) {
	keeper.IterateVotes(ctx, proposalID, func(vote types2.Vote) bool {
		votes = append(votes, vote)
		return false
	})
	return
}

// IterateAllVotes iterates over the all the stored votes and performs a callback function
func (keeper Keeper) IterateAllVotes(ctx sdk.Context, cb func(vote types2.Vote) (stop bool)) {
	store := ctx.KVStore(keeper.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types2.VotesKeyPrefix)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var vote types2.Vote
		keeper.cdc.MustUnmarshalLengthPrefixed(iterator.Value(), &vote)

		if cb(vote) {
			break
		}
	}
}

// IterateVotes iterates over the all the proposals votes and performs a callback function
func (keeper Keeper) IterateVotes(ctx sdk.Context, proposalID uint64, cb func(vote types2.Vote) (stop bool)) {
	store := ctx.KVStore(keeper.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types2.VotesKey(proposalID))

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var vote types2.Vote
		keeper.cdc.MustUnmarshalLengthPrefixed(iterator.Value(), &vote)

		if cb(vote) {
			break
		}
	}
}

// deleteVote deletes a vote from a given proposalID and voter from the store
func (keeper Keeper) deleteVote(ctx sdk.Context, proposalID uint64, voterAddr sdk.ValAddress) {
	store := ctx.KVStore(keeper.storeKey)
	store.Delete(types2.VoteKey(proposalID, voterAddr))
}
