package gov

import (
	appTypes "bitbucket.org/decimalteam/go-node/types"
	keep "bitbucket.org/decimalteam/go-node/x/gov/internal/keeper"
	"bitbucket.org/decimalteam/go-node/x/gov/internal/types"
	"bitbucket.org/decimalteam/go-node/x/validator"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestTickPassedVotingPeriod(t *testing.T) {
	input := getTestInput(t, 10, GenesisState{}, nil)
	SortAddresses(input.addrs)

	ctx := input.ctx

	govHandler := NewHandler(input.keeper)

	ctx = ctx.WithBlockHeight(1000000000000)

	inactiveQueue := input.keeper.InactiveProposalQueueIterator(ctx, uint64(ctx.BlockHeight()))
	require.False(t, inactiveQueue.Valid())
	inactiveQueue.Close()
	activeQueue := input.keeper.ActiveProposalQueueIterator(ctx, uint64(ctx.BlockHeight()))
	require.False(t, activeQueue.Valid())
	activeQueue.Close()

	proposer, err := sdk.AccAddressFromBech32(validator.DAOAddress1)
	require.NoError(t, err)
	newProposalMsg := NewMsgSubmitProposal(
		types.Content{
			Title:       "title",
			Description: "desc",
		},
		proposer,
		1000000000000+5,
		1000000000000+10,
	)

	res, err := govHandler(ctx, newProposalMsg)
	require.NoError(t, err)
	require.NotNil(t, res)

	ctx = ctx.WithBlockHeight(1000000000000 + 5)

	EndBlocker(ctx, input.keeper)

	inactiveQueue = input.keeper.InactiveProposalQueueIterator(ctx, uint64(ctx.BlockHeight()))
	require.False(t, inactiveQueue.Valid())

	activeQueue = input.keeper.ActiveAllProposalQueueIterator(ctx)
	require.True(t, activeQueue.Valid())

	activeProposalID := types.GetProposalIDFromBytes(activeQueue.Value())
	proposal, ok := input.keeper.GetProposal(ctx, activeProposalID)
	require.True(t, ok)
	require.Equal(t, StatusVotingPeriod, proposal.Status)

	activeQueue.Close()

	EndBlocker(ctx, input.keeper)

	activeQueue = input.keeper.ActiveProposalQueueIterator(ctx, uint64(ctx.BlockHeight()))
	require.False(t, activeQueue.Valid())
	activeQueue.Close()
}

func TestProposalPassedEndBlocker(t *testing.T) {
	input := getTestInput(t, 1, GenesisState{}, nil)
	SortAddresses(input.addrs)

	validatorHandler := validator.NewHandler(input.vk)

	ctx := input.ctx

	valAddr := sdk.ValAddress(input.addrs[0])

	createValidators(t, validatorHandler, ctx, []sdk.ValAddress{valAddr}, []int64{10})
	valUpdates := validator.EndBlocker(ctx, input.vk, input.ck, input.sk, false)
	require.Equal(t, 1, len(valUpdates))

	proposal, err := input.keeper.SubmitProposal(ctx, keep.TestProposal.Content, keep.TestProposal.VotingStartBlock, keep.TestProposal.VotingEndBlock)
	require.NoError(t, err)

	ctx = ctx.WithBlockHeight(int64(proposal.VotingStartBlock))

	EndBlocker(ctx, input.keeper)

	err = input.keeper.AddVote(ctx, proposal.ProposalID, sdk.ValAddress(input.addrs[0]), types.OptionYes)
	require.NoError(t, err)

	ctx = ctx.WithBlockHeight(int64(proposal.VotingEndBlock))

	EndBlocker(ctx, input.keeper)

	proposal, ok := input.keeper.GetProposal(ctx, proposal.ProposalID)
	require.True(t, ok)

	require.Equal(t, types.StatusPassed, proposal.Status)

	t.Log(proposal.FinalTallyResult)
	require.Equal(t, appTypes.TokensFromConsensusPower(10), proposal.FinalTallyResult.Yes)
}

func TestProposalPassedEndBlocker2(t *testing.T) {
	input := getTestInput(t, 2, GenesisState{}, nil)
	SortAddresses(input.addrs)

	validatorHandler := validator.NewHandler(input.vk)

	ctx := input.ctx

	valAddr1 := sdk.ValAddress(input.addrs[0])
	valAddr2 := sdk.ValAddress(input.addrs[1])

	createValidators(t, validatorHandler, ctx, []sdk.ValAddress{valAddr1, valAddr2}, []int64{21, 10})
	valUpdates := validator.EndBlocker(ctx, input.vk, input.ck, input.sk, false)
	require.Equal(t, 2, len(valUpdates))

	proposal, err := input.keeper.SubmitProposal(ctx, keep.TestProposal.Content, keep.TestProposal.VotingStartBlock, keep.TestProposal.VotingEndBlock)
	require.NoError(t, err)

	ctx = ctx.WithBlockHeight(int64(proposal.VotingStartBlock))

	EndBlocker(ctx, input.keeper)

	err = input.keeper.AddVote(ctx, proposal.ProposalID, sdk.ValAddress(input.addrs[0]), types.OptionYes)
	require.NoError(t, err)

	err = input.keeper.AddVote(ctx, proposal.ProposalID, sdk.ValAddress(input.addrs[1]), types.OptionNo)
	require.NoError(t, err)

	ctx = ctx.WithBlockHeight(int64(proposal.VotingEndBlock))

	EndBlocker(ctx, input.keeper)

	proposal, ok := input.keeper.GetProposal(ctx, proposal.ProposalID)
	require.True(t, ok)

	require.Equal(t, types.StatusPassed, proposal.Status)

	t.Log(proposal.FinalTallyResult)
	require.Equal(t, appTypes.TokensFromConsensusPower(21), proposal.FinalTallyResult.Yes)
	require.Equal(t, appTypes.TokensFromConsensusPower(10), proposal.FinalTallyResult.No)
	require.Equal(t, sdk.ZeroInt(), proposal.FinalTallyResult.Abstain)
}

func TestEndBlockerProposalRejected(t *testing.T) {
	input := getTestInput(t, 1, GenesisState{}, nil)
	SortAddresses(input.addrs)

	validatorHandler := validator.NewHandler(input.vk)

	ctx := input.ctx

	valAddr := sdk.ValAddress(input.addrs[0])

	createValidators(t, validatorHandler, ctx, []sdk.ValAddress{valAddr}, []int64{10})
	valUpdates := validator.EndBlocker(ctx, input.vk, input.ck, input.sk, false)
	require.Equal(t, 1, len(valUpdates))

	proposal, err := input.keeper.SubmitProposal(ctx, keep.TestProposal.Content, keep.TestProposal.VotingStartBlock, keep.TestProposal.VotingEndBlock)
	require.NoError(t, err)

	ctx = ctx.WithBlockHeight(int64(proposal.VotingStartBlock))

	EndBlocker(ctx, input.keeper)

	err = input.keeper.AddVote(ctx, proposal.ProposalID, sdk.ValAddress(input.addrs[0]), types.OptionNo)
	require.NoError(t, err)

	ctx = ctx.WithBlockHeight(int64(proposal.VotingEndBlock))

	EndBlocker(ctx, input.keeper)

	proposal, ok := input.keeper.GetProposal(ctx, proposal.ProposalID)
	require.True(t, ok)

	require.Equal(t, types.StatusRejected, proposal.Status)

	t.Log(proposal.FinalTallyResult)
	require.Equal(t, appTypes.TokensFromConsensusPower(10), proposal.FinalTallyResult.No)
}

func TestEndBlockerProposalRejected2(t *testing.T) {
	input := getTestInput(t, 2, GenesisState{}, nil)
	SortAddresses(input.addrs)

	validatorHandler := validator.NewHandler(input.vk)

	ctx := input.ctx

	valAddr1 := sdk.ValAddress(input.addrs[0])
	valAddr2 := sdk.ValAddress(input.addrs[1])

	createValidators(t, validatorHandler, ctx, []sdk.ValAddress{valAddr1, valAddr2}, []int64{10, 10})
	valUpdates := validator.EndBlocker(ctx, input.vk, input.ck, input.sk, false)
	require.Equal(t, 2, len(valUpdates))

	proposal, err := input.keeper.SubmitProposal(ctx, keep.TestProposal.Content, keep.TestProposal.VotingStartBlock, keep.TestProposal.VotingEndBlock)
	require.NoError(t, err)

	ctx = ctx.WithBlockHeight(int64(proposal.VotingStartBlock))

	EndBlocker(ctx, input.keeper)

	err = input.keeper.AddVote(ctx, proposal.ProposalID, sdk.ValAddress(input.addrs[0]), types.OptionYes)
	require.NoError(t, err)

	err = input.keeper.AddVote(ctx, proposal.ProposalID, sdk.ValAddress(input.addrs[1]), types.OptionNo)
	require.NoError(t, err)

	ctx = ctx.WithBlockHeight(int64(proposal.VotingEndBlock))

	EndBlocker(ctx, input.keeper)

	proposal, ok := input.keeper.GetProposal(ctx, proposal.ProposalID)
	require.True(t, ok)

	require.Equal(t, types.StatusRejected, proposal.Status)

	t.Log(proposal.FinalTallyResult)
	require.Equal(t, appTypes.TokensFromConsensusPower(10), proposal.FinalTallyResult.Yes)
	require.Equal(t, appTypes.TokensFromConsensusPower(10), proposal.FinalTallyResult.No)
	require.Equal(t, sdk.ZeroInt(), proposal.FinalTallyResult.Abstain)
}

func TestEndBlockerProposalRejected3(t *testing.T) {
	input := getTestInput(t, 2, GenesisState{}, nil)
	SortAddresses(input.addrs)

	validatorHandler := validator.NewHandler(input.vk)

	ctx := input.ctx

	valAddr1 := sdk.ValAddress(input.addrs[0])
	valAddr2 := sdk.ValAddress(input.addrs[1])

	createValidators(t, validatorHandler, ctx, []sdk.ValAddress{valAddr1, valAddr2}, []int64{20, 10})
	valUpdates := validator.EndBlocker(ctx, input.vk, input.ck, input.sk, false)
	require.Equal(t, 2, len(valUpdates))

	proposal, err := input.keeper.SubmitProposal(ctx, keep.TestProposal.Content, keep.TestProposal.VotingStartBlock, keep.TestProposal.VotingEndBlock)
	require.NoError(t, err)

	ctx = ctx.WithBlockHeight(int64(proposal.VotingStartBlock))

	EndBlocker(ctx, input.keeper)

	err = input.keeper.AddVote(ctx, proposal.ProposalID, sdk.ValAddress(input.addrs[0]), types.OptionYes)
	require.NoError(t, err)

	ctx = ctx.WithBlockHeight(int64(proposal.VotingEndBlock))

	EndBlocker(ctx, input.keeper)

	proposal, ok := input.keeper.GetProposal(ctx, proposal.ProposalID)
	require.True(t, ok)

	require.Equal(t, types.StatusRejected, proposal.Status)

	t.Log(proposal.FinalTallyResult)
	require.Equal(t, appTypes.TokensFromConsensusPower(20), proposal.FinalTallyResult.Yes)
	require.Equal(t, appTypes.TokensFromConsensusPower(10), proposal.FinalTallyResult.Abstain)
	require.Equal(t, sdk.ZeroInt(), proposal.FinalTallyResult.No)
}
