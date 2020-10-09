package gov

import (
	"bitbucket.org/decimalteam/go-node/x/gov/internal/types"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// EndBlocker called every block, process inflation, update validator set.
func EndBlocker(ctx sdk.Context, keeper Keeper) {
	logger := keeper.Logger(ctx)

	// delete inactive proposal from store
	keeper.IterateInactiveProposalsQueue(ctx, uint64(ctx.BlockHeight()), func(proposal Proposal) bool {
		keeper.RemoveFromInactiveProposalQueue(ctx, proposal.ProposalID, proposal.VotingStartBlock)
		keeper.InsertActiveProposalQueue(ctx, proposal.ProposalID, proposal.VotingEndBlock)
		proposal.Status = StatusVotingPeriod
		keeper.SetProposal(ctx, proposal)

		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeInactiveProposal,
				sdk.NewAttribute(types.AttributeKeyProposalID, fmt.Sprintf("%d", proposal.ProposalID)),
			),
		)

		return false
	})

	// fetch active proposals whose voting periods have ended (are passed the block time)
	keeper.IterateActiveProposalsQueue(ctx, uint64(ctx.BlockHeight()), func(proposal Proposal) bool {
		var tagValue, logMsg string

		passes, tallyResults := keeper.Tally(ctx, proposal)

		if passes {
			proposal.Status = StatusPassed
			tagValue = types.AttributeValueProposalPassed
			logMsg = "passed"
		} else {
			proposal.Status = StatusRejected
			tagValue = types.AttributeValueProposalRejected
			logMsg = "rejected"
		}

		proposal.FinalTallyResult = tallyResults

		keeper.SetProposal(ctx, proposal)
		keeper.RemoveFromActiveProposalQueue(ctx, proposal.ProposalID, proposal.VotingEndBlock)

		logger.Info(
			fmt.Sprintf(
				"proposal %d (%s) tallied; result: %s",
				proposal.ProposalID, proposal.GetTitle(), logMsg,
			),
		)

		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeActiveProposal,
				sdk.NewAttribute(types.AttributeKeyProposalID, fmt.Sprintf("%d", proposal.ProposalID)),
				sdk.NewAttribute(types.AttributeKeyProposalResult, tagValue),
				sdk.NewAttribute(types.AttributeKeyResultVoteYes, tallyResults.Yes.String()),
				sdk.NewAttribute(types.AttributeKeyResultVoteAbstain, tallyResults.Abstain.String()),
				sdk.NewAttribute(types.AttributeKeyResultVoteNo, tallyResults.No.String()),
			),
		)
		return false
	})
}
