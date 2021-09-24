package gov

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"bitbucket.org/decimalteam/go-node/x/gov/internal/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func printJson(data interface{}) {
	res, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(res))
}

func fileExist(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

func PushNewSkipPlan(planfile, name string) {
	skipPlans := LoadSkipPlans(planfile)
	skipPlans[name] = true

	bytes, err := json.Marshal(skipPlans)
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(planfile, bytes, 0644)
	if err != nil {
		panic(err)
	}
}

func LoadSkipPlans(planfile string) map[string]bool {
	plans := make(map[string]bool)

	if !fileExist(planfile) {
		err := ioutil.WriteFile(planfile, []byte("{}"), 0600)
		if err != nil {
			panic(err)
		}
	}

	bytes, err := ioutil.ReadFile(planfile)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(bytes, &plans)
	if err != nil {
		panic(err)
	}

	return plans
}

func BeginBlocker(ctx sdk.Context, k Keeper) {
	planfile := os.Getenv("HOME") + "/.decimal/daemon/config/skip_plans.json"
	fmt.Println("VERSION 1!")

	plan, found := k.GetUpgradePlan(ctx)
	if !found {
		return
	}

	skipPlans := LoadSkipPlans(planfile)
	_, ok := skipPlans[plan.Name]
	if ok {
		return
	}

	// printJson(plan)

	// To make sure clear upgrade is executed at the same block
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
		k.ApplyUpgrade(ctx, plan)

		PushNewSkipPlan(planfile, plan.Name)
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
}
