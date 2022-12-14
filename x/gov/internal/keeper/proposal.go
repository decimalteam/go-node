package keeper

import (
	"fmt"
	"strings"

	"bitbucket.org/decimalteam/go-node/x/gov/internal/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SubmitProposal create new proposal given a content
func (keeper Keeper) SubmitProposal(ctx sdk.Context, content types.Content, votingStartBlock, votingEndBlock uint64) (types.Proposal, error) {
	proposalID, err := keeper.GetProposalID(ctx)
	if err != nil {
		return types.Proposal{}, err
	}

	proposal := types.NewProposal(content, proposalID, votingStartBlock, votingEndBlock)

	keeper.SetProposal(ctx, proposal)
	keeper.InsertInactiveProposalQueue(ctx, proposalID, proposal.VotingStartBlock)
	keeper.SetProposalID(ctx, proposalID+1)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeSubmitProposal,
			sdk.NewAttribute(types.AttributeKeyProposalID, fmt.Sprintf("%d", proposalID)),
		),
	)

	return proposal, nil
}

// GetProposalID gets the highest proposal ID
func (keeper Keeper) GetProposalID(ctx sdk.Context) (proposalID uint64, err error) {
	store := ctx.KVStore(keeper.storeKey)
	bz := store.Get(types.ProposalIDKey)
	if bz == nil {
		return 0, types.ErrInvalidGenesis()
	}

	proposalID = types.GetProposalIDFromBytes(bz)
	return proposalID, nil
}

// SetProposalID sets the new proposal ID to the store
func (keeper Keeper) SetProposalID(ctx sdk.Context, proposalID uint64) {
	store := ctx.KVStore(keeper.storeKey)
	store.Set(types.ProposalIDKey, types.GetProposalIDBytes(proposalID))
}

// GetProposal get proposal from store by ProposalID
func (keeper Keeper) GetProposal(ctx sdk.Context, proposalID uint64) (proposal types.Proposal, ok bool) {
	store := ctx.KVStore(keeper.storeKey)
	bz := store.Get(types.ProposalKey(proposalID))
	if bz == nil {
		return
	}
	keeper.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &proposal)
	return proposal, true
}

// SetProposal set a proposal to store
func (keeper Keeper) SetProposal(ctx sdk.Context, proposal types.Proposal) {
	store := ctx.KVStore(keeper.storeKey)
	bz := keeper.cdc.MustMarshalBinaryLengthPrefixed(proposal)
	store.Set(types.ProposalKey(proposal.ProposalID), bz)
}

// DeleteProposal deletes a proposal from store
func (keeper Keeper) DeleteProposal(ctx sdk.Context, proposalID uint64) {
	store := ctx.KVStore(keeper.storeKey)
	proposal, ok := keeper.GetProposal(ctx, proposalID)
	if !ok {
		panic(fmt.Sprintf("couldn't find proposal with id#%d", proposalID))
	}
	keeper.RemoveFromInactiveProposalQueue(ctx, proposalID, proposal.VotingStartBlock)
	keeper.RemoveFromActiveProposalQueue(ctx, proposalID, proposal.VotingEndBlock)
	store.Delete(types.ProposalKey(proposalID))
}

// IterateProposals iterates over the all the proposals and performs a callback function
func (keeper Keeper) IterateProposals(ctx sdk.Context, cb func(proposal types.Proposal) (stop bool)) {
	store := ctx.KVStore(keeper.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.ProposalsKeyPrefix)
	logger := keeper.Logger(ctx)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		proposal, err := UnmarshalProposal(keeper.cdc, iterator.Value())
		switch {
		case err != nil && strings.HasPrefix(err.Error(), "Bytes left over in UnmarshalBinaryLengthPrefixed,"):
			logger.Error(fmt.Sprintf("cannot parse proposal with key %s, and value %s", string(iterator.Key()), string(iterator.Value())))
		case err != nil:
			panic(err)
		}

		if cb(proposal) {
			break
		}
	}
}

// GetProposals returns all the proposals from store
func (keeper Keeper) GetProposals(ctx sdk.Context) (proposals types.Proposals) {
	keeper.IterateProposals(ctx, func(proposal types.Proposal) bool {
		proposals = append(proposals, proposal)
		return false
	})
	return
}

// amino methods

func MustMarshaProposal(cdc *codec.Codec, ubd types.Proposal) []byte {
	return cdc.MustMarshalBinaryLengthPrefixed(ubd)
}

func UnmarshalProposal(cdc *codec.Codec, value []byte) (prop types.Proposal, err error) {
	err = cdc.UnmarshalBinaryLengthPrefixed(value, &prop)
	return prop, err
}

func MustUnmarshalProposal(cdc *codec.Codec, value []byte) types.Proposal {
	prop, err := UnmarshalProposal(cdc, value)
	if err != nil {
		panic(err)
	}
	return prop
}
