package gov

import (
	"fmt"
	"os"
	"time"

	ncfg "bitbucket.org/decimalteam/go-node/config"
	"bitbucket.org/decimalteam/go-node/x/gov/internal/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	updateStart = make(chan bool)
	// blockerStop  = make(chan bool)
	downloadStat = make(map[string]bool)
)

func BeginBlocker(ctx sdk.Context, keeper Keeper) {
	select {
	case <-updateStart:
		// blockerStop <- true
		time.Sleep(10 * time.Second)
	default:
	}
}

// EndBlocker called every block, process inflation, update validator set.
func EndBlocker(ctx sdk.Context, keeper Keeper) {
	// if ctx.BlockHeight() < updates.Update1Block {
	// 	return
	// }

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
					types.EventTypeInactiveProposal,
					sdk.NewAttribute(types.AttributeKeyProposalID, fmt.Sprintf("%d", proposal.ProposalID)),
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
					sdk.NewAttribute(types.AttributeKeyTotalVotingPower, totalVotingPower.String()),
				),
			)
		}

		return false
	})

	plan, found := keeper.GetUpgradePlan(ctx)
	if !found {
		return
	}

	if ctx.BlockHeight() > plan.Height {
		keeper.ClearUpgradePlan(ctx)
		return
	}

	allBlocks := ncfg.UpdatesInfo.AllBlocks
	if _, ok := allBlocks[plan.Name]; ok {
		return
	}

	go checkUpdate(ctx, keeper, plan)
}

func checkUpdate(ctx sdk.Context, k Keeper, plan types.Plan) {
	_, ok := downloadStat[plan.Name]

	if ctx.BlockHeight() > (plan.Height-plan.ToDownload) && ctx.BlockHeight() < plan.Height && !ok {
		for _, name := range ncfg.NameFiles {
			// example:
			// from "http://127.0.0.1/95000/decd"
			// to "http://127.0.0.1/95000/linux/ubuntu/20.04/decd"
			newUrl := k.GenerateUrl(fmt.Sprintf("%s/%s", plan.Name, name))
			if newUrl == "" {
				return
			}

			if !k.UrlPageExist(newUrl) {
				return
			}

			downloadStat[plan.Name] = true
			go k.DownloadBinary(k.GetDownloadName(name), newUrl)
		}
	}

	// To make sure clear upgrade is executed at./de the same block
	if plan.ShouldExecute(ctx) {
		// If skip upgrade has been set for current height, we clear the upgrade plan
		if k.IsSkipHeight(ctx.BlockHeight()) {
			skipUpgradeMsg := fmt.Sprintf("UPGRADE \"%s\" SKIPPED at %d: %s", plan.Name, plan.Height, plan.Info)
			ctx.Logger().Info(skipUpgradeMsg)

			// Clear the upgrade plan at current height
			k.ClearUpgradePlan(ctx)
			return
		}

		// We have an upgrade handler for this upgrade name, so apply the upgrade
		ctx.Logger().Info(fmt.Sprintf("applying upgrade \"%s\" at %s", plan.Name, plan.DueAt()))
		ctx = ctx.WithBlockGasMeter(sdk.NewInfiniteGasMeter())

		updateStart <- true
		// <-blockerStop
		err := k.ApplyUpgrade(ctx, plan)
		if err != nil {
			ctx.Logger().Error(fmt.Sprintf("upgrade \"%s\" with %s", plan.Name, err.Error()))
			return
		}

		ncfg.UpdatesInfo.Push(plan.Name, ctx.BlockHeight())
		os.Exit(0)
		return
	}
}
