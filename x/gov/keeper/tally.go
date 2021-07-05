package keeper

import (
	types2 "bitbucket.org/decimalteam/go-node/x/gov/types"
	"bitbucket.org/decimalteam/go-node/x/validator/exported"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Tally iterates over the votes and updates the tally of a proposal based on the voting power of the
// voters
func (keeper Keeper) Tally(ctx sdk.Context, proposal types2.Proposal) (passes bool, tallyResults types2.TallyResult, totalVotingPower sdk.Dec) {
	results := make(map[types2.VoteOption]sdk.Dec)
	results[types2.OptionYes] = sdk.ZeroDec()
	results[types2.OptionAbstain] = sdk.ZeroDec()
	results[types2.OptionNo] = sdk.ZeroDec()

	totalVotingPower = sdk.ZeroDec()
	currValidators := make(map[string]types2.ValidatorGovInfo)

	// fetch all the bonded validators, insert them into currValidators
	keeper.vk.IterateBondedValidatorsByPower(ctx, func(index int64, validator exported.ValidatorI) (stop bool) {
		if index == 9 {
			return true
		}
		currValidators[validator.GetOperator().String()] = types2.NewValidatorGovInfo(
			validator.GetOperator(),
			validator.GetBondedTokens(),
			types2.OptionEmpty,
		)

		return false
	})

	keeper.IterateVotes(ctx, proposal.ProposalID, func(vote types2.Vote) bool {
		// if validator, just record it in the map
		valAddrStr := vote.Voter
		if val, ok := currValidators[valAddrStr]; ok {
			val.Vote = vote.Option
			currValidators[valAddrStr] = val
		}

		voterAddr, _ := sdk.ValAddressFromBech32(vote.Voter)

		keeper.deleteVote(ctx, vote.ProposalID, voterAddr)
		return false
	})

	// iterate over the validators again to tally their voting power
	for _, val := range currValidators {
		if val.Vote == types2.OptionEmpty {
			val.Vote = types2.OptionAbstain
		}

		votingPower := val.BondedTokens

		results[val.Vote] = results[val.Vote].Add(sdk.NewDecFromInt(votingPower))
		totalVotingPower = totalVotingPower.Add(sdk.NewDecFromInt(votingPower))
	}

	tallyParams := keeper.GetTallyParams(ctx)
	tallyResults = types2.NewTallyResultFromMap(results)

	// TODO: Upgrade the spec to cover all of these cases & remove pseudocode.
	// If there is no staked coins, the proposal fails
	if keeper.vk.TotalBondedTokens(ctx).IsZero() {
		return false, tallyResults, totalVotingPower
	}

	// If no one votes (everyone abstains), proposal fails
	if totalVotingPower.Sub(results[types2.OptionAbstain]).Equal(sdk.ZeroDec()) {
		return false, tallyResults, totalVotingPower
	}

	if results[types2.OptionYes].Quo(totalVotingPower).GT(tallyParams.Quorum) {
		return true, tallyResults, totalVotingPower
	}

	return false, tallyResults, totalVotingPower
}
