package utils

import (
	"bitbucket.org/decimalteam/go-node/x/gov/internal/types"
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
)

const (
	defaultPage  = 1
	defaultLimit = 30 // should be consistent with tendermint/tendermint/rpc/core/pipe.go:19
)

// Proposer contains metadata of a governance proposal used for querying a
// proposer.
type Proposer struct {
	ProposalID uint64 `json:"proposal_id" yaml:"proposal_id"`
	Proposer   string `json:"proposer" yaml:"proposer"`
}

// NewProposer returns a new Proposer given id and proposer
func NewProposer(proposalID uint64, proposer string) Proposer {
	return Proposer{proposalID, proposer}
}

func (p Proposer) String() string {
	return fmt.Sprintf("Proposal with ID %d was proposed by %s", p.ProposalID, p.Proposer)
}

// QueryVotesByTxQuery will query for votes via a direct txs tags query. It
// will fetch and build votes directly from the returned txs and return a JSON
// marshalled result or any error that occurred.
func QueryVotesByTxQuery(cliCtx context.CLIContext, params types.QueryProposalVotesParams) ([]byte, error) {
	var (
		events = []string{
			fmt.Sprintf("%s.%s='%s'", sdk.EventTypeMessage, sdk.AttributeKeyAction, types.TypeMsgVote),
			fmt.Sprintf("%s.%s='%s'", types.EventTypeProposalVote, types.AttributeKeyProposalID, []byte(fmt.Sprintf("%d", params.ProposalID))),
		}
		votes      []types.Vote
		nextTxPage = defaultPage
		totalLimit = params.Limit * params.Page
	)
	// query interrupted either if we collected enough votes or tx indexer run out of relevant txs
	for len(votes) < totalLimit {
		searchResult, err := utils.QueryTxsByEvents(cliCtx, events, nextTxPage, defaultLimit)
		if err != nil {
			return nil, err
		}
		nextTxPage++
		for _, info := range searchResult.Txs {
			for _, msg := range info.Tx.GetMsgs() {
				if msg.Type() == types.TypeMsgVote {
					voteMsg := msg.(types.MsgVote)

					votes = append(votes, types.Vote{
						Voter:      voteMsg.Voter,
						ProposalID: params.ProposalID,
						Option:     voteMsg.Option,
					})
				}
			}
		}
		if len(searchResult.Txs) != defaultLimit {
			break
		}
	}
	start, end := client.Paginate(len(votes), params.Page, params.Limit, 100)
	if start < 0 || end < 0 {
		votes = []types.Vote{}
	} else {
		votes = votes[start:end]
	}
	if cliCtx.Indent {
		return cliCtx.Codec.MarshalJSONIndent(votes, "", "  ")
	}
	return cliCtx.Codec.MarshalJSON(votes)
}

// QueryVoteByTxQuery will query for a single vote via a direct txs tags query.
func QueryVoteByTxQuery(cliCtx context.CLIContext, params types.QueryVoteParams) ([]byte, error) {
	events := []string{
		fmt.Sprintf("%s.%s='%s'", sdk.EventTypeMessage, sdk.AttributeKeyAction, types.TypeMsgVote),
		fmt.Sprintf("%s.%s='%s'", types.EventTypeProposalVote, types.AttributeKeyProposalID, []byte(fmt.Sprintf("%d", params.ProposalID))),
		fmt.Sprintf("%s.%s='%s'", sdk.EventTypeMessage, sdk.AttributeKeySender, []byte(params.Voter.String())),
	}

	// NOTE: SearchTxs is used to facilitate the txs query which does not currently
	// support configurable pagination.
	searchResult, err := utils.QueryTxsByEvents(cliCtx, events, defaultPage, defaultLimit)
	if err != nil {
		return nil, err
	}
	for _, info := range searchResult.Txs {
		for _, msg := range info.Tx.GetMsgs() {
			// there should only be a single vote under the given conditions
			if msg.Type() == types.TypeMsgVote {
				voteMsg := msg.(types.MsgVote)

				vote := types.Vote{
					Voter:      voteMsg.Voter,
					ProposalID: params.ProposalID,
					Option:     voteMsg.Option,
				}

				if cliCtx.Indent {
					return cliCtx.Codec.MarshalJSONIndent(vote, "", "  ")
				}

				return cliCtx.Codec.MarshalJSON(vote)
			}
		}
	}

	return nil, fmt.Errorf("address '%s' did not vote on proposalID %d", params.Voter, params.ProposalID)
}

// QueryProposerByTxQuery will query for a proposer of a governance proposal by
// ID.
func QueryProposerByTxQuery(cliCtx context.CLIContext, proposalID uint64) (Proposer, error) {
	events := []string{
		fmt.Sprintf("%s.%s='%s'", sdk.EventTypeMessage, sdk.AttributeKeyAction, types.TypeMsgSubmitProposal),
		fmt.Sprintf("%s.%s='%s'", types.EventTypeSubmitProposal, types.AttributeKeyProposalID, []byte(fmt.Sprintf("%d", proposalID))),
	}

	// NOTE: SearchTxs is used to facilitate the txs query which does not currently
	// support configurable pagination.
	searchResult, err := utils.QueryTxsByEvents(cliCtx, events, defaultPage, defaultLimit)
	if err != nil {
		return Proposer{}, err
	}

	for _, info := range searchResult.Txs {
		for _, msg := range info.Tx.GetMsgs() {
			// there should only be a single proposal under the given conditions
			if msg.Type() == types.TypeMsgSubmitProposal {
				subMsg := msg.(types.MsgSubmitProposal)
				return NewProposer(proposalID, subMsg.Proposer.String()), nil
			}
		}
	}

	return Proposer{}, fmt.Errorf("failed to find the proposer for proposalID %d", proposalID)
}

// QueryProposalByID takes a proposalID and returns a proposal
func QueryProposalByID(proposalID uint64, cliCtx context.CLIContext, queryRoute string) ([]byte, error) {
	params := types.NewQueryProposalParams(proposalID)
	bz, err := cliCtx.Codec.MarshalJSON(params)
	if err != nil {
		return nil, err
	}

	res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/proposal", queryRoute), bz)
	if err != nil {
		return nil, err
	}

	return res, err
}

// NormalizeVoteOption - normalize user specified vote option
func NormalizeVoteOption(option string) string {
	switch option {
	case "Yes", "yes":
		return types.OptionYes.String()

	case "Abstain", "abstain":
		return types.OptionAbstain.String()

	case "No", "no":
		return types.OptionNo.String()

	default:
		return ""
	}
}

//NormalizeProposalStatus - normalize user specified proposal status
func NormalizeProposalStatus(status string) string {
	switch status {
	case "DepositPeriod", "deposit_period":
		return "DepositPeriod"
	case "VotingPeriod", "voting_period":
		return "VotingPeriod"
	case "Passed", "passed":
		return "Passed"
	case "Rejected", "rejected":
		return "Rejected"
	}
	return ""
}
