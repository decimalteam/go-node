package types

import (
	"fmt"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
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

func ErrUnknownProposal(proposalID uint64) *sdkerrors.Error {
	return sdkerrors.New(
		DefaultCodespace,
		CodeUnknownProposal,
		//todo proposalID
		"unknown proposal",
	)
}

func ErrInactiveProposal(proposalID uint64) *sdkerrors.Error {
	return sdkerrors.New(
		DefaultCodespace,
		CodeInactiveProposal,
		//todo proposalID
		"inactive proposal",
	)
}

func ErrAlreadyActiveProposal() *sdkerrors.Error {
	return sdkerrors.New(
		DefaultCodespace,
		CodeAlreadyActiveProposal,
		"proposal already active",
	)
}

func ErrInvalidProposalContent() *sdkerrors.Error {
	return sdkerrors.New(
		DefaultCodespace,
		CodeInvalidProposalContent,
		"missing content",
	)
}

func ErrInvalidProposalContentTitleBlank() *sdkerrors.Error {
	return sdkerrors.New(
		DefaultCodespace,
		CodeInvalidProposalContentTitleBlank,
		"proposal title cannot be blank",
	)
}

func ErrInvalidProposalContentTitleLong(MaxTitleLength int) *sdkerrors.Error {
	return sdkerrors.New(
		DefaultCodespace,
		CodeInvalidProposalContentTitleLong,
		//todo MaxTitleLength
		"proposal title is longer than max length of %d",
	)
}

func ErrInvalidProposalContentDescrBlank() *sdkerrors.Error {
	return sdkerrors.New(
		DefaultCodespace,
		CodeInvalidProposalContentDescrBlank,
		"proposal description cannot be blank",
	)
}

func ErrInvalidProposalContentDescrLong(MaxDescriptionLength int) *sdkerrors.Error {
	return sdkerrors.New(
		DefaultCodespace,
		CodeInvalidProposalContentDescrLong,
		//todo MaxDescriptionLength
		"proposal description is longer than max length of %d",
	)
}

func ErrInvalidProposalType(ProposalType string) *sdkerrors.Error {
	return sdkerrors.New(
		DefaultCodespace,
		CodeInvalidProposalType,
		//todo ProposalType
		"invalid proposal type",
	)
}

func ErrInvalidVote(option string) *sdkerrors.Error {
	return sdkerrors.New(
		DefaultCodespace,
		CodeInvalidVote,
		//todo option
		"invalid vote option",
	)
}

func ErrInvalidGenesis() *sdkerrors.Error {
	return sdkerrors.New(
		DefaultCodespace,
		CodeInvalidGenesis,
		"initial proposal ID hasn't been set",
	)
}

func ErrNoProposalHandlerExists() *sdkerrors.Error {
	return sdkerrors.New(
		DefaultCodespace,
		CodeNoProposalHandlerExists,
		"no handler exists for proposal type",
	)
}

func ErrInvalidStartEndBlocks() *sdkerrors.Error {
	return sdkerrors.New(
		DefaultCodespace,
		CodeInvalidStartEndBlocks,
		"invalid start or end blocks",
	)
}

func ErrSubmitProposal(err string) *sdkerrors.Error {
	return sdkerrors.New(
		DefaultCodespace,
		CodeSubmitProposal,
		//todo error
		"error submit proposal: %s",
	)
}

func ErrStartBlock() *sdkerrors.Error {
	return sdkerrors.New(
		DefaultCodespace,
		CodeStartBlock,
		"start block must greater then current block height",
	)
}

func ErrDurationTooLong() *sdkerrors.Error {
	return sdkerrors.New(
		DefaultCodespace,
		CodeDurationTooLong,
		fmt.Sprintf("duration too long. Max duration = %d month (%d blocks)", DurationInMonth, DurationInBlocks),
	)
}

func ErrNotAllowed() *sdkerrors.Error {
	return sdkerrors.New(
		DefaultCodespace,
		CodeNotAllowed,
		//todo address view
		fmt.Sprintf("not allowed to create the proposal from this address"),
	)
}
