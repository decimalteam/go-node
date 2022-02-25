package types

import (
	"encoding/json"
	"fmt"
	"strings"

	//"bitbucket.org/decimalteam/go-node/x/validator"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// DefaultStartingProposalID is 1
const DefaultStartingProposalID uint64 = 1

var AllowedAddresses = []string{
	"dx1lx4lvt8sjuxj8vw5dcf6knnq0pacre4w6hdh2v",
	"dx1mvqrrrlcd0gdt256jxg7n68e4neppu5t24e8h6",
	"dx1nrr6er27mmcufmaqm4dyu6c5r6489cfm35m4ft",
	"dx16rr3cvdgj8jsywhx8lfteunn9uz0xg2c7ua9nl",
	"dx1ajytg8jg8ypx0rj9p792x32fuxyezga43jd3ry",
	"dx1azre0dtclv5y05ufynkhswzh0cwh4ktzlas3mp",

	"dxvalcons1xwxakrlp9zpxyylc6udc0tun34dt70sfgwv6gk",
	"dxvalcons1q7kxlea59qu75gdr0z6g4hhsa7p5z6ekysk76z",
	"dxvalcons198dvzuj6n9hyvfkjh0255l9hx46p62j20nue97",
	"dxvalcons1hhgda24n56gwme73hxkhl6yah4mpe30k52r3t8",
	"dxvalcons17ntss5hyuutk5w4a4upptz3xc9f3f0tgwwmfw8",
	"dxvalcons15aczs4nlk06faxzxvuuevl0g9r5zve2pnualvp",
}

/*
var AllowedAddresses = []string{
	validator.DAOAddress1,
	validator.DAOAddress2,
	validator.DAOAddress3,
	validator.DevelopAddress1,
	validator.DevelopAddress2,
	validator.DevelopAddress3,
}
*/

func CheckProposalAddress(address sdk.AccAddress) bool {
	for _, allowedAddress := range AllowedAddresses {
		if allowedAddress == address.String() {
			return true
		}
	}
	return false
}

// Proposal defines a struct used by the governance module to allow for voting
// on network changes.
type Proposal struct {
	Content

	ProposalID       uint64         `json:"id" yaml:"id"`                                 //  ID of the proposal
	Status           ProposalStatus `json:"proposal_status" yaml:"proposal_status"`       // Status of the Proposal {Pending, Active, Passed, Rejected}
	FinalTallyResult TallyResult    `json:"final_tally_result" yaml:"final_tally_result"` // Result of Tallys

	VotingStartBlock uint64 `json:"voting_start_time" yaml:"voting_start_time"` // Time of the block where MinDeposit was reached. -1 if MinDeposit is not reached
	VotingEndBlock   uint64 `json:"voting_end_time" yaml:"voting_end_time"`     // Time that the VotingPeriod for this proposal will end and votes will be tallied
}

// NewProposal creates a new Proposal instance
func NewProposal(content Content, id, votingStartBlock, VotingEndBlock uint64) Proposal {
	return Proposal{
		Content:          content,
		ProposalID:       id,
		Status:           StatusWaiting,
		FinalTallyResult: EmptyTallyResult(),
		VotingStartBlock: votingStartBlock,
		VotingEndBlock:   VotingEndBlock,
	}
}

// String implements stringer interface
func (p Proposal) String() string {
	return fmt.Sprintf(`Proposal %d:
  Title:              %s
  Status:             %s
  Voting Start Time:  %d
  Voting End Time:    %d
  Description:        %s`,
		p.ProposalID, p.Title,
		p.Status, p.VotingStartBlock, p.VotingEndBlock, p.Description,
	)
}

// Proposals is an array of proposal
type Proposals []Proposal

// String implements stringer interface
func (p Proposals) String() string {
	out := "ID - (Status) Title\n"
	for _, prop := range p {
		out += fmt.Sprintf("%d - (%s) %s\n",
			prop.ProposalID, prop.Status, prop.Title)
	}
	return strings.TrimSpace(out)
}

type (
	// ProposalQueue defines a queue for proposal ids
	ProposalQueue []uint64

	// ProposalStatus is a type alias that represents a proposal status as a byte
	ProposalStatus byte
)

// Valid Proposal statuses
const (
	StatusNil          ProposalStatus = 0x00
	StatusWaiting      ProposalStatus = 0x01
	StatusVotingPeriod ProposalStatus = 0x02
	StatusPassed       ProposalStatus = 0x03
	StatusRejected     ProposalStatus = 0x04
	StatusFailed       ProposalStatus = 0x05
)

// ProposalStatusFromString turns a string into a ProposalStatus
func ProposalStatusFromString(str string) (ProposalStatus, error) {
	switch str {
	case "Waiting":
		return StatusWaiting, nil

	case "VotingPeriod":
		return StatusVotingPeriod, nil

	case "Passed":
		return StatusPassed, nil

	case "Rejected":
		return StatusRejected, nil

	case "Failed":
		return StatusFailed, nil

	case "":
		return StatusNil, nil

	default:
		return ProposalStatus(0xff), fmt.Errorf("'%s' is not a valid proposal status", str)
	}
}

// ValidProposalStatus returns true if the proposal status is valid and false
// otherwise.
func ValidProposalStatus(status ProposalStatus) bool {
	if status == StatusWaiting ||
		status == StatusVotingPeriod ||
		status == StatusPassed ||
		status == StatusRejected ||
		status == StatusFailed {
		return true
	}
	return false
}

// Marshal needed for protobuf compatibility
func (status ProposalStatus) Marshal() ([]byte, error) {
	return []byte{byte(status)}, nil
}

// Unmarshal needed for protobuf compatibility
func (status *ProposalStatus) Unmarshal(data []byte) error {
	*status = ProposalStatus(data[0])
	return nil
}

// MarshalJSON Marshals to JSON using string representation of the status
func (status ProposalStatus) MarshalJSON() ([]byte, error) {
	return json.Marshal(status.String())
}

// UnmarshalJSON Unmarshals from JSON assuming Bech32 encoding
func (status *ProposalStatus) UnmarshalJSON(data []byte) error {
	var s string
	err := json.Unmarshal(data, &s)
	if err != nil {
		return err
	}

	bz2, err := ProposalStatusFromString(s)
	if err != nil {
		return err
	}

	*status = bz2
	return nil
}

// String implements the Stringer interface.
func (status ProposalStatus) String() string {
	switch status {
	case StatusWaiting:
		return "Waiting"

	case StatusVotingPeriod:
		return "VotingPeriod"

	case StatusPassed:
		return "Passed"

	case StatusRejected:
		return "Rejected"

	case StatusFailed:
		return "Failed"

	default:
		return ""
	}
}

// Format implements the fmt.Formatter interface.
// nolint: errcheck
func (status ProposalStatus) Format(s fmt.State, verb rune) {
	switch verb {
	case 's':
		s.Write([]byte(status.String()))
	default:
		// TODO: Do this conversion more directly
		s.Write([]byte(fmt.Sprintf("%v", byte(status))))
	}
}
