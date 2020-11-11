package gov

import (
	"bitbucket.org/decimalteam/go-node/x/gov/internal/types"
	"github.com/stretchr/testify/require"
	"testing"

	abci "github.com/tendermint/tendermint/abci/types"
)

func TestTickPassedVotingPeriod(t *testing.T) {
	input := getMockApp(t, 10, GenesisState{}, nil)
	SortAddresses(input.addrs)

	header := abci.Header{Height: input.mApp.LastBlockHeight() + 1}
	input.mApp.BeginBlock(abci.RequestBeginBlock{Header: header})

	ctx := input.mApp.BaseApp.NewContext(false, abci.Header{})
	govHandler := NewHandler(input.keeper)

	inactiveQueue := input.keeper.InactiveProposalQueueIterator(ctx, uint64(ctx.BlockHeight()))
	require.False(t, inactiveQueue.Valid())
	inactiveQueue.Close()
	activeQueue := input.keeper.ActiveProposalQueueIterator(ctx, uint64(ctx.BlockHeight()))
	require.False(t, activeQueue.Valid())
	activeQueue.Close()

	newProposalMsg := NewMsgSubmitProposal(
		types.Content{
			Title:       "title",
			Description: "desc",
		},
		input.addrs[0],
		5,
		10,
	)

	res, err := govHandler(ctx, newProposalMsg)
	require.NoError(t, err)
	require.NotNil(t, res)

	ctx = ctx.WithBlockHeight(5)

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

/*
func TestProposalPassedEndblocker(t *testing.T) {
	input := getMockApp(t, 1, GenesisState{}, nil)
	SortAddresses(input.addrs)

	handler := NewHandler(input.keeper)
	stakingHandler := validator.NewHandler(input.vk)

	header := abci.Header{Height: input.mApp.LastBlockHeight() + 1}
	input.mApp.BeginBlock(abci.RequestBeginBlock{Header: header})
	ctx := input.mApp.BaseApp.NewContext(false, abci.Header{})

	valAddr := sdk.ValAddress(input.addrs[0])

	createValidators(t, stakingHandler, ctx, []sdk.ValAddress{valAddr}, []int64{10})
	validator.EndBlocker(ctx, input.vk, input.ck, input.sk, false)

	proposal, err := input.keeper.SubmitProposal(ctx, keep.TestProposal.Content, keep.TestProposal.VotingStartBlock, keep.TestProposal.VotingEndBlock)
	require.NoError(t, err)

	deposits := initialModuleAccCoins.Add(proposal.TotalDeposit...).Add(proposalCoins...)
	require.True(t, moduleAccCoins.IsEqual(deposits))

	err = input.keeper.AddVote(ctx, proposal.ProposalID, input.addrs[0], OptionYes)
	require.NoError(t, err)

	newHeader := ctx.BlockHeader()
	newHeader.Time = uint64(ctx.BlockHeight()).Add(input.keeper.GetDepositParams(ctx).MaxDepositPeriod).Add(input.keeper.GetVotingParams(ctx).VotingPeriod)
	ctx = ctx.WithBlockHeader(newHeader)

	EndBlocker(ctx, input.keeper)

	macc = input.keeper.GetGovernanceAccount(ctx)
	require.NotNil(t, macc)
	require.True(t, macc.GetCoins().IsEqual(initialModuleAccCoins))
}

func TestEndBlockerProposalHandlerFailed(t *testing.T) {
	// hijack the router to one that will fail in a proposal's handler
	input := getMockApp(t, 1, GenesisState{}, nil)
	SortAddresses(input.addrs)

	handler := NewHandler(input.keeper)
	stakingHandler := staking.NewHandler(input.vk)

	header := abci.Header{Height: input.mApp.LastBlockHeight() + 1}
	input.mApp.BeginBlock(abci.RequestBeginBlock{Header: header})
	ctx := input.mApp.BaseApp.NewContext(false, abci.Header{})

	valAddr := sdk.ValAddress(input.addrs[0])

	createValidators(t, stakingHandler, ctx, []sdk.ValAddress{valAddr}, []int64{10})
	staking.EndBlocker(ctx, input.vk)

	// Create a proposal where the handler will pass for the test proposal
	// because the value of contextKeyBadProposal is true.
	ctx = ctx.WithValue(contextKeyBadProposal, true)
	proposal, err := input.keeper.SubmitProposal(ctx, keep.TestProposal)
	require.NoError(t, err)

	proposalCoins := sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.TokensFromConsensusPower(10)))
	newDepositMsg := NewMsgDeposit(input.addrs[0], proposal.ProposalID, proposalCoins)

	res, err := handler(ctx, newDepositMsg)
	require.NoError(t, err)
	require.NotNil(t, res)

	err = input.keeper.AddVote(ctx, proposal.ProposalID, input.addrs[0], OptionYes)
	require.NoError(t, err)

	newHeader := ctx.BlockHeader()
	newHeader.Time = uint64(ctx.BlockHeight()).Add(input.keeper.GetDepositParams(ctx).MaxDepositPeriod).Add(input.keeper.GetVotingParams(ctx).VotingPeriod)
	ctx = ctx.WithBlockHeader(newHeader)

	// Set the contextKeyBadProposal value to false so that the handler will fail
	// during the processing of the proposal in the EndBlocker.
	ctx = ctx.WithValue(contextKeyBadProposal, false)

	// validate that the proposal fails/has been rejected
	EndBlocker(ctx, input.keeper)
}
*/
