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

	CodeUnknownProposal         CodeType = 100
	CodeInactiveProposal        CodeType = 200
	CodeAlreadyActiveProposal   CodeType = 300
	CodeInvalidProposalContent  CodeType = 400
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

var (
//ErrUnknownProposal         = sdkerrors.Register(ModuleName, 100, "unknown proposal")
//ErrInactiveProposal        = sdkerrors.Register(ModuleName, 200, "inactive proposal")
//ErrAlreadyActiveProposal   = sdkerrors.Register(ModuleName, 300, "proposal already active")
//ErrInvalidProposalContent  = sdkerrors.Register(ModuleName, 400, "invalid proposal content")
//ErrInvalidProposalType     = sdkerrors.Register(ModuleName, 500, "invalid proposal type")
//ErrInvalidVote             = sdkerrors.Register(ModuleName, 600, "invalid vote option")
//ErrInvalidGenesis          = sdkerrors.Register(ModuleName, 700, "invalid genesis state")
//ErrNoProposalHandlerExists = sdkerrors.Register(ModuleName, 800, "no handler exists for proposal type")
//ErrInvalidStartEndBlocks   = sdkerrors.Register(ModuleName, 900, "invalid start or end blocks")
//ErrSubmitProposal          = sdkerrors.Register(ModuleName, 1000, "error submit proposal")
//ErrStartBlock              = sdkerrors.Register(ModuleName, 1100, "start block must greater then current block height")
//ErrDurationTooLong         = sdkerrors.Register(ModuleName, 1200, "duration too long. Max duration = 3 month (1296000 blocks)")
//ErrNotAllowed              = sdkerrors.Register(ModuleName, 1300, "not allowed to create the proposal from this address")
)

func ErrUnknownProposal(proposalID uint64) *sdkerrors.Error {
	return sdkerrors.New(
		DefaultCodespace,
		CodeUnknownProposal,
		//todo proposalID
		"unknown proposal",
	)
}

func ErrInactiveProposal() *sdkerrors.Error {
	return sdkerrors.New(
		DefaultCodespace,
		CodeInactiveProposal,
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
		"invalid proposal content",
	)
}

func ErrInvalidProposalType() *sdkerrors.Error {
	return sdkerrors.New(
		DefaultCodespace,
		CodeInvalidProposalType,
		"invalid proposal type",
	)
}

func ErrInvalidVote() *sdkerrors.Error {
	return sdkerrors.New(
		DefaultCodespace,
		CodeInvalidVote,
		"invalid vote option",
	)
}

func ErrInvalidGenesis() *sdkerrors.Error {
	return sdkerrors.New(
		DefaultCodespace,
		CodeInvalidGenesis,
		"invalid genesis state",
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

func ErrSubmitProposal() *sdkerrors.Error {
	return sdkerrors.New(
		DefaultCodespace,
		CodeSubmitProposal,
		"error submit proposal",
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
