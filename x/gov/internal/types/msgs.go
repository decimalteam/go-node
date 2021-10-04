package types

import (
	"errors"
	"fmt"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// Governance message types and routes
const (
	TypeMsgVote            = "vote"
	TypeMsgSubmitProposal  = "submit_proposal"
	TypeMsgSoftwareUpgrade = "software_upgrade"
)

var _, _ sdk.Msg = MsgSubmitProposal{}, MsgVote{}

// MsgSubmitProposal defines a message to create a governance proposal with a
// given content and initial deposit
type MsgSubmitProposal struct {
	Content          Content        `json:"content" yaml:"content"`
	Proposer         sdk.AccAddress `json:"proposer" yaml:"proposer"` //  Address of the proposer
	VotingStartBlock uint64         `json:"voting_start_block" yaml:"voting_start_block"`
	VotingEndBlock   uint64         `json:"voting_end_block" yaml:"voting_end_block"`
}

// NewMsgSubmitProposal creates a new MsgSubmitProposal instance
func NewMsgSubmitProposal(content Content, proposer sdk.AccAddress, votingStartBlock, votingEndBlock uint64) MsgSubmitProposal {
	return MsgSubmitProposal{
		Content:          content,
		Proposer:         proposer,
		VotingStartBlock: votingStartBlock,
		VotingEndBlock:   votingEndBlock,
	}
}

// Route implements Msg
func (msg MsgSubmitProposal) Route() string { return RouterKey }

// Type implements Msg
func (msg MsgSubmitProposal) Type() string { return TypeMsgSubmitProposal }

// ValidateBasic implements Msg
func (msg MsgSubmitProposal) ValidateBasic() error {
	if msg.Content.Title == "" || msg.Content.Description == "" {
		return ErrInvalidProposalContent()
	}
	if msg.Proposer.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Proposer.String())
	}

	if msg.VotingStartBlock >= msg.VotingEndBlock {
		return ErrInvalidStartEndBlocks(
			strconv.FormatUint(msg.VotingStartBlock, 10), strconv.FormatUint(msg.VotingEndBlock, 10))
	}

	if msg.VotingEndBlock-msg.VotingStartBlock > 1296000 {
		return ErrDurationTooLong()
	}

	return nil
}

// String implements the Stringer interface
func (msg MsgSubmitProposal) String() string {
	return fmt.Sprintf(`Submit Proposal Message:
  Title:          %s
  Description:    %s
`, msg.Content.Title, msg.Content.Description)
}

// GetSignBytes implements Msg
func (msg MsgSubmitProposal) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// GetSigners implements Msg
func (msg MsgSubmitProposal) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Proposer}
}

// MsgVote defines a message to cast a vote
type MsgVote struct {
	ProposalID uint64         `json:"proposal_id" yaml:"proposal_id"` // ID of the proposal
	Voter      sdk.ValAddress `json:"voter" yaml:"voter"`             //  address of the voter
	Option     VoteOption     `json:"option" yaml:"option"`           //  option from OptionSet chosen by the voter
}

// NewMsgVote creates a message to cast a vote on an active proposal
func NewMsgVote(voter sdk.ValAddress, proposalID uint64, option VoteOption) MsgVote {
	return MsgVote{proposalID, voter, option}
}

// Route implements Msg
func (msg MsgVote) Route() string { return RouterKey }

// Type implements Msg
func (msg MsgVote) Type() string { return TypeMsgVote }

// ValidateBasic implements Msg
func (msg MsgVote) ValidateBasic() error {
	if msg.Voter.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Voter.String())
	}
	if !ValidVoteOption(msg.Option) {
		return ErrInvalidVote(msg.Option.String())
	}

	return nil
}

// String implements the Stringer interface
func (msg MsgVote) String() string {
	return fmt.Sprintf(`Vote Message:
  Proposal ID: %d
  Option:      %s
`, msg.ProposalID, msg.Option)
}

// GetSignBytes implements Msg
func (msg MsgVote) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// GetSigners implements Msg
func (msg MsgVote) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.Voter)}
}

const AddressForSoftwareUpgrade = "dx134fskg774x3ch2hz9n7phkzxyxvl645z5dsel5"

// Software Upgrade Proposals
type MsgSoftwareUpgradeProposal struct {
	Title       string         `json:"title" yaml:"title"`
	Description string         `json:"description" yaml:"description"`
	Plan        Plan           `json:"plan" yaml:"plan"`
	Proposer    sdk.AccAddress `json:"proposer" yaml:"proposer"` //  Address of the proposer
}

func NewSoftwareUpgradeProposal(title, description string, plan Plan, proposer sdk.AccAddress) MsgSoftwareUpgradeProposal {
	return MsgSoftwareUpgradeProposal{title, description, plan, proposer}
}

const ProposalTypeSoftwareUpgrade = "SoftwareUpgrade"

// Route implements Msg
func (msg MsgSoftwareUpgradeProposal) Route() string { return RouterKey }

// Type implements Msg
func (msg MsgSoftwareUpgradeProposal) Type() string { return ProposalTypeSoftwareUpgrade }

// ValidateBasic implements Msg
func (msg MsgSoftwareUpgradeProposal) ValidateBasic() error {
	address, err := sdk.AccAddressFromBech32(AddressForSoftwareUpgrade)
	if err != nil {
		return err
	}
	// ADD NOT
	if !msg.Proposer.Equals(address) {
		return errors.New("not allowed")
	}
	return nil
}

// String implements the Stringer interface
func (msg MsgSoftwareUpgradeProposal) String() string {
	return fmt.Sprintf(`Submit Proposal Message:
  Title:          %s
  Description:    %s
`, msg.Title, msg.Description)
}

// GetSignBytes implements Msg
func (msg MsgSoftwareUpgradeProposal) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// GetSigners implements Msg
func (msg MsgSoftwareUpgradeProposal) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Proposer}
}
