package keeper

import (
	"bitbucket.org/decimalteam/go-node/utils/updates"
	"bitbucket.org/decimalteam/go-node/x/gov/internal/types"
	"bitbucket.org/decimalteam/go-node/x/validator/exported"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Tally iterates over the votes and updates the tally of a proposal based on the voting power of the
// voters
func (keeper Keeper) Tally(ctx sdk.Context, proposal types.Proposal) (passes bool, tallyResults types.TallyResult, totalVotingPower sdk.Dec) {
	results := make(map[types.VoteOption]sdk.Dec)
	results[types.OptionYes] = sdk.ZeroDec()
	results[types.OptionAbstain] = sdk.ZeroDec()
	results[types.OptionNo] = sdk.ZeroDec()

	totalVotingPower = sdk.ZeroDec()
	currValidators := make(map[string]types.ValidatorGovInfo)

	// fetch all the bonded validators, insert them into currValidators
	keeper.vk.IterateBondedValidatorsByPower(ctx, func(index int64, validator exported.ValidatorI) (stop bool) {
		if ctx.BlockHeight() >= updates.Update1Block {
			if index == 9 {
				return true
			}
		}

		currValidators[validator.GetOperator().String()] = types.NewValidatorGovInfo(
			validator.GetOperator(),
			validator.GetBondedTokens(),
			types.OptionEmpty,
		)

		return false
	})

	keeper.IterateVotes(ctx, proposal.ProposalID, func(vote types.Vote) bool {
		// if validator, just record it in the map
		valAddrStr := vote.Voter.String()
		if val, ok := currValidators[valAddrStr]; ok {
			val.Vote = vote.Option
			currValidators[valAddrStr] = val
		}

		keeper.deleteVote(ctx, vote.ProposalID, vote.Voter)
		return false
	})

	// iterate over the validators again to tally their voting power
	for _, val := range currValidators {
		if val.Vote == types.OptionEmpty {
			if ctx.BlockHeight() >= updates.Update1Block {
				val.Vote = types.OptionAbstain
			} else {
				continue
			}
		}

		votingPower := val.BondedTokens

		results[val.Vote] = results[val.Vote].Add(sdk.NewDecFromInt(votingPower))
		totalVotingPower = totalVotingPower.Add(sdk.NewDecFromInt(votingPower))
	}

	tallyParams := keeper.GetTallyParams(ctx)
	tallyResults = types.NewTallyResultFromMap(results)

	// TODO: Upgrade the spec to cover all of these cases & remove pseudocode.
	// If there is no staked coins, the proposal fails
	if keeper.vk.TotalBondedTokens(ctx).IsZero() {
		return false, tallyResults, totalVotingPower
	}

	if ctx.BlockHeight() >= updates.Update1Block {
		// If no one votes (everyone abstains), proposal fails
		if totalVotingPower.Sub(results[types.OptionAbstain]).Equal(sdk.ZeroDec()) {
			return false, tallyResults, totalVotingPower
		}

		if results[types.OptionYes].Quo(totalVotingPower).GT(tallyParams.Quorum) {
			return true, tallyResults, totalVotingPower
		}

		return false, tallyResults, totalVotingPower
	} else {
		// If there is not enough quorum of votes, the proposal fails
		percentVoting := totalVotingPower.Quo(keeper.vk.TotalBondedTokens(ctx).ToDec())
		if percentVoting.LT(tallyParams.Quorum) {
			return false, tallyResults, totalVotingPower
		}

		// If no one votes (everyone abstains), proposal fails
		if totalVotingPower.Sub(results[types.OptionAbstain]).Equal(sdk.ZeroDec()) {
			return false, tallyResults, totalVotingPower
		}

		return true, tallyResults, totalVotingPower
	}
}
