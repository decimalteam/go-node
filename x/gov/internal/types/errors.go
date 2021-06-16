package types

import (
	"bitbucket.org/decimalteam/go-node/utils/errors"
	"fmt"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"strconv"
)

// x/gov module sentinel errors

type CodeType = uint32

const (
	// Default coin codespace
	DefaultCodespace string = ModuleName

	DurationInMonth  int = 3
	DurationInBlocks int = 1296000

	CodeUnknownProposal       CodeType = 100
	CodeInactiveProposal      CodeType = 200
	CodeAlreadyActiveProposal CodeType = 300
	// proposal content
	CodeInvalidProposalContent           CodeType = 400
	CodeInvalidProposalContentTitleBlank CodeType = 401
	CodeInvalidProposalContentTitleLong  CodeType = 402
	CodeInvalidProposalContentDescrBlank CodeType = 403
	CodeInvalidProposalContentDescrLong  CodeType = 404

	CodeInvalidProposalType     CodeType = 500
	CodeInvalidVote             CodeType = 600
	CodeInvalidGenesis          CodeType = 700
	CodeNoProposalHandlerExists CodeType = 800
	CodeInvalidStartEndBlocks   CodeType = 900
	CodeSubmitProposal          CodeType = 1000
	CodeStartBlock              CodeType = 1100
	CodeDurationTooLong         CodeType = 1200
	CodeNotAllowed              CodeType = 1300
)

func ErrUnknownProposal(proposalID string) *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeUnknownProposal,
		fmt.Sprintf("unknown proposal: %s", proposalID),
		errors.NewParam("proposalID", proposalID),
	)
}

func ErrInactiveProposal(proposalID string) *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeInactiveProposal,
		fmt.Sprintf("inactive proposal: %s", proposalID),
		errors.NewParam("proposalID", proposalID),
	)
}

func ErrAlreadyActiveProposal() *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeAlreadyActiveProposal,
		fmt.Sprintf("proposal already active"),
	)
}

func ErrInvalidProposalContent() *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeInvalidProposalContent,
		fmt.Sprintf("missing content"),
	)
}

func ErrInvalidProposalContentTitleBlank() *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeInvalidProposalContentTitleBlank,
		fmt.Sprintf("proposal title cannot be blank"),
	)
}

func ErrInvalidProposalContentTitleLong(MaxTitleLength string) *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeInvalidProposalContentTitleLong,
		fmt.Sprintf("proposal title is longer than max length of %s", MaxTitleLength),
		errors.NewParam("MaxTitleLength", MaxTitleLength),
	)
}

func ErrInvalidProposalContentDescrBlank() *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeInvalidProposalContentDescrBlank,
		fmt.Sprintf("proposal description cannot be blank"),
	)
}

func ErrInvalidProposalContentDescrLong(MaxDescriptionLength string) *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeInvalidProposalContentDescrLong,
		fmt.Sprintf("proposal description is longer than max length of %d", MaxDescriptionLength),
		errors.NewParam("MaxDescriptionLength", MaxDescriptionLength),
	)
}

func ErrInvalidProposalType(ProposalType string) *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeInvalidProposalType,
		fmt.Sprintf("invalid proposal type: %s", ProposalType),
		errors.NewParam("ProposalType", ProposalType),
	)
}

func ErrInvalidVote(option string) *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeInvalidVote,
		fmt.Sprintf("invalid vote option: %s", option),
		errors.NewParam("option", option),
	)
}

func ErrInvalidGenesis() *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeInvalidGenesis,
		fmt.Sprintf("initial proposal ID hasn't been set"),
	)
}

func ErrNoProposalHandlerExists() *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeNoProposalHandlerExists,
		fmt.Sprintf("no handler exists for proposal type"),
	)
}

func ErrInvalidStartEndBlocks(StartBlock string, EndBlock string) *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeInvalidStartEndBlocks,
		fmt.Sprintf("invalid start or end blocks: start %s,  end %s ", StartBlock, EndBlock),
		errors.NewParam("StartBlock", StartBlock),
		errors.NewParam("EndBlock", EndBlock),
	)
}

func ErrSubmitProposal(error string) *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeSubmitProposal,
		fmt.Sprintf("error submit proposal: %s", error),
		errors.NewParam("error", error),
	)
}

func ErrStartBlock() *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeStartBlock,
		fmt.Sprintf("start block must greater then current block height"),
	)
}

func ErrDurationTooLong() *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeDurationTooLong,
		fmt.Sprintf("start block must greater then current block height"),
		errors.NewParam("maxDurationInMonth", fmt.Sprintf("%d", DurationInMonth)),
		errors.NewParam("maxDurationInBlocks", fmt.Sprintf("%d", DurationInBlocks)),
	)
}

func ErrNotAllowed() *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeNotAllowed,
		fmt.Sprintf("not allowed to create the proposal from this address"),
	)
}

func getCodeString(code CodeType) string {
	return strconv.FormatInt(int64(code), 10)
}
