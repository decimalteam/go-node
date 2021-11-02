package gov

import (
	types2 "bitbucket.org/decimalteam/go-node/x/gov/types"
	"fmt"
	"os"
	"strings"

	ncfg "bitbucket.org/decimalteam/go-node/config"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	downloadStat = make(map[string]bool)
)

func BeginBlocker(ctx sdk.Context, k Keeper) {
	plan, found := k.GetUpgradePlan(ctx)
	if !found {
		return
	}

	// "http://127.0.0.1/95000@v1.2.1"
	splited := strings.Split(plan.Name, "@")
	if len(splited) != 2 {
		k.ClearUpgradePlan(ctx)
		return
	}

	planURL := splited[0]
	version := splited[1]

	if ctx.BlockHeight() > plan.Height {
		if ncfg.DecimalVersion != version {
			ctx.Logger().Error(fmt.Sprintf("failed upgrade \"%s\" at height %d", plan.Name, plan.Height))
			os.Exit(2)
		}
		k.ClearUpgradePlan(ctx)
		return
	}

	allBlocks := ncfg.UpdatesInfo.AllBlocks
	if _, ok := allBlocks[planURL]; ok {
		return
	}

	_, ok := downloadStat[planURL]

	if ctx.BlockHeight() > (plan.Height-plan.ToDownload) && ctx.BlockHeight() < plan.Height && !ok {
		for _, name := range ncfg.NameFiles {
			// example:
			// from "http://127.0.0.1/95000/decd"
			// to "http://127.0.0.1/95000/linux/ubuntu/20.04/decd"
			newUrl := k.GenerateUrl(fmt.Sprintf("%s/%s", planURL, name))
			if newUrl == "" {
				return
			}

			if !k.UrlPageExist(newUrl) {
				return
			}

			downloadStat[planURL] = true
			go k.DownloadBinary(k.GetDownloadName(name), newUrl)
		}
	}

	// To make sure clear upgrade is executed at./de the same block
	if plan.ShouldExecute(ctx) {
		// If skip upgrade has been set for current height, we clear the upgrade plan
		if k.IsSkipHeight(ctx.BlockHeight()) {
			skipUpgradeMsg := fmt.Sprintf("UPGRADE \"%s\" SKIPPED at %d: %s", planURL, plan.Height, plan.Info)
			ctx.Logger().Info(skipUpgradeMsg)

			// Clear the upgrade plan at current height
			k.ClearUpgradePlan(ctx)
			return
		}

		// We have an upgrade handler for this upgrade name, so apply the upgrade
		ctx.Logger().Info(fmt.Sprintf("applying upgrade \"%s\" at %s", planURL, plan.DueAt()))
		ctx = ctx.WithBlockGasMeter(sdk.NewInfiniteGasMeter())

		err := k.ApplyUpgrade(ctx, plan)
		if err != nil {
			ctx.Logger().Error(fmt.Sprintf("upgrade \"%s\" with %s", planURL, err.Error()))
			return
		}

		ncfg.UpdatesInfo.Push(planURL, ctx.BlockHeight())
		os.Exit(0)
		return
	}
}
// EndBlocker called every block, process inflation, update validator set.
func EndBlocker(ctx sdk.Context, keeper Keeper) {
	logger := keeper.Logger(ctx)

	// delete inactive proposal from store
	keeper.IterateAllInactiveProposalsQueue(ctx, func(proposal Proposal) bool {
		if ctx.BlockHeight() == int64(proposal.VotingStartBlock) {
			keeper.RemoveFromInactiveProposalQueue(ctx, proposal.ProposalID, proposal.VotingStartBlock)
			keeper.InsertActiveProposalQueue(ctx, proposal.ProposalID, proposal.VotingEndBlock)
			proposal.Status = StatusVotingPeriod
			keeper.SetProposal(ctx, proposal)

			ctx.EventManager().EmitEvent(
				sdk.NewEvent(
					types2.EventTypeInactiveProposal,
					sdk.NewAttribute(types2.AttributeKeyProposalID, fmt.Sprintf("%d", proposal.ProposalID)),
				),
			)
		}

		return false
	})

	// fetch active proposals whose voting periods have ended (are passed the block time)
	keeper.IterateAllActiveProposalsQueue(ctx, func(proposal Proposal) bool {
		if int64(proposal.VotingEndBlock) == ctx.BlockHeight() {
			var tagValue, logMsg string

			passes, tallyResults, totalVotingPower := keeper.Tally(ctx, proposal)

			if passes {
				proposal.Status = StatusPassed
				tagValue = types2.AttributeValueProposalPassed
				logMsg = "passed"
			} else {
				proposal.Status = StatusRejected
				tagValue = types2.AttributeValueProposalRejected
				logMsg = "rejected"
			}

			proposal.FinalTallyResult = tallyResults

			keeper.SetProposal(ctx, proposal)
			keeper.RemoveFromActiveProposalQueue(ctx, proposal.ProposalID, proposal.VotingEndBlock)

			logger.Info(
				fmt.Sprintf(
					"proposal %d (%s) tallied; result: %s",
					proposal.ProposalID, /*proposal.GetTitle(),*/ logMsg,
				),
			)

			ctx.EventManager().EmitEvent(
				sdk.NewEvent(
					types2.EventTypeActiveProposal,
					sdk.NewAttribute(types2.AttributeKeyProposalID, fmt.Sprintf("%d", proposal.ProposalID)),
					sdk.NewAttribute(types2.AttributeKeyProposalResult, tagValue),
					sdk.NewAttribute(types2.AttributeKeyResultVoteYes, tallyResults.Yes.String()),
					sdk.NewAttribute(types2.AttributeKeyResultVoteAbstain, tallyResults.Abstain.String()),
					sdk.NewAttribute(types2.AttributeKeyResultVoteNo, tallyResults.No.String()),
					sdk.NewAttribute(types2.AttributeKeyTotalVotingPower, totalVotingPower.String()),
				),
			)
		}

		return false
	})
}
