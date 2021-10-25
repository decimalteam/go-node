package types

// Governance module event types
const (
	EventTypeSubmitProposal   = "submit_proposal"
	EventTypeProposalVote     = "proposal_vote"
	EventTypeInactiveProposal = "inactive_proposal"
	EventTypeActiveProposal   = "active_proposal"

	AttributeKeyProposalResult           = "proposal_result"
	AttributeKeyOption                   = "option"
	AttributeKeyProposalID               = "proposal_id"
	AttributeKeyVotingPeriodStart        = "voting_period_start"
	AttributeValueCategory               = "governance"
	AttributeValueProposalPassed         = "proposal_passed"   // met vote quorum
	AttributeValueProposalRejected       = "proposal_rejected" // didn't meet vote quorum
	AttributeValueProposalFailed         = "proposal_failed"   // error on proposal handler
	AttributeKeyProposalType             = "proposal_type"
	AttributeKeyProposalTitle            = "proposal_title"
	AttributeKeyProposalDescription      = "proposal_description"
	AttributeKeyProposalVotingStartBlock = "proposal_voting_start_block"
	AttributeKeyProposalVotingEndBlock   = "proposal_voting_end_block"

	AttributeKeyResultVoteYes     = "result_vote_yes"
	AttributeKeyResultVoteAbstain = "result_vote_abstain"
	AttributeKeyResultVoteNo      = "result_vote_no"
	AttributeKeyTotalVotingPower  = "total_voting_power"

	AttributeKeyUpgradeHeight = "upgrade_height"
)
