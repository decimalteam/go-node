package gov

import (
	types2 "bitbucket.org/decimalteam/go-node/x/gov/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// InitGenesis - store genesis parameters
func InitGenesis(ctx sdk.Context, k Keeper, supplyKeeper types2.SupplyKeeper, data GenesisState) {

	k.SetProposalID(ctx, data.StartingProposalID)
	k.SetTallyParams(ctx, data.TallyParams)

	for _, vote := range data.Votes {
		k.SetVote(ctx, vote)
	}

	for _, proposal := range data.Proposals {
		switch proposal.Status {
		case StatusWaiting:
			k.InsertInactiveProposalQueue(ctx, proposal.ProposalID, proposal.VotingStartBlock)
		case StatusVotingPeriod:
			k.InsertActiveProposalQueue(ctx, proposal.ProposalID, proposal.VotingEndBlock)
		}
		k.SetProposal(ctx, proposal)
	}
}

// ExportGenesis - output genesis parameters
func ExportGenesis(ctx sdk.Context, k Keeper) GenesisState {
	startingProposalID, _ := k.GetProposalID(ctx)
	tallyParams := k.GetTallyParams(ctx)
	proposals := k.GetProposals(ctx)

	var proposalsVotes Votes
	for _, proposal := range proposals {

		votes := k.GetVotes(ctx, proposal.ProposalID)
		proposalsVotes = append(proposalsVotes, votes...)
	}

	return GenesisState{
		StartingProposalID: startingProposalID,
		Votes:              proposalsVotes,
		Proposals:          proposals,
		TallyParams:        tallyParams,
	}
}
