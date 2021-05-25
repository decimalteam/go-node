package types

import (
	"encoding/json"
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

func ErrUnknownProposal(proposalID uint64) *sdkerrors.Error {
	jsonData, _ := json.Marshal(
		map[string]string{
			"code":       getCodeString(CodeUnknownProposal),
			"codespace":  DefaultCodespace,
			"desc":       fmt.Sprintf("unknown proposal: %d", proposalID),
			"proposalID": fmt.Sprintf("%d", proposalID),
		},
	)
	return sdkerrors.New(
		DefaultCodespace,
		CodeUnknownProposal,
		string(jsonData),
	)
}

func ErrInactiveProposal(proposalID uint64) *sdkerrors.Error {
	jsonData, _ := json.Marshal(
		map[string]string{
			"code":       getCodeString(CodeInactiveProposal),
			"codespace":  DefaultCodespace,
			"desc":       fmt.Sprintf("inactive proposal: %d", proposalID),
			"proposalID": fmt.Sprintf("%d", proposalID),
		},
	)
	return sdkerrors.New(
		DefaultCodespace,
		CodeInactiveProposal,
		string(jsonData),
	)
}

func ErrAlreadyActiveProposal() *sdkerrors.Error {
	jsonData, _ := json.Marshal(
		map[string]string{
			"code":      getCodeString(CodeAlreadyActiveProposal),
			"codespace": DefaultCodespace,
			"desc":      "proposal already active",
		},
	)
	return sdkerrors.New(
		DefaultCodespace,
		CodeAlreadyActiveProposal,
		string(jsonData),
	)
}

func ErrInvalidProposalContent() *sdkerrors.Error {
	jsonData, _ := json.Marshal(
		map[string]string{
			"code":      getCodeString(CodeInvalidProposalContent),
			"codespace": DefaultCodespace,
			"desc":      "missing content",
		},
	)
	return sdkerrors.New(
		DefaultCodespace,
		CodeInvalidProposalContent,
		string(jsonData),
	)
}

func ErrInvalidProposalContentTitleBlank() *sdkerrors.Error {
	jsonData, _ := json.Marshal(
		map[string]string{
			"code":      getCodeString(CodeInvalidProposalContentTitleBlank),
			"codespace": DefaultCodespace,
			"desc":      "proposal title cannot be blank",
		},
	)
	return sdkerrors.New(
		DefaultCodespace,
		CodeInvalidProposalContentTitleBlank,
		string(jsonData),
	)
}

func ErrInvalidProposalContentTitleLong(MaxTitleLength int) *sdkerrors.Error {
	jsonData, _ := json.Marshal(
		map[string]string{
			"code":           getCodeString(CodeInvalidProposalContentTitleLong),
			"codespace":      DefaultCodespace,
			"desc":           fmt.Sprintf("proposal title is longer than max length of %d", MaxTitleLength),
			"MaxTitleLength": fmt.Sprintf("%d", MaxTitleLength),
		},
	)
	return sdkerrors.New(
		DefaultCodespace,
		CodeInvalidProposalContentTitleLong,
		string(jsonData),
	)
}

func ErrInvalidProposalContentDescrBlank() *sdkerrors.Error {
	jsonData, _ := json.Marshal(
		map[string]string{
			"code":      getCodeString(CodeInvalidProposalContentDescrBlank),
			"codespace": DefaultCodespace,
			"desc":      "proposal description cannot be blank",
		},
	)
	return sdkerrors.New(
		DefaultCodespace,
		CodeInvalidProposalContentDescrBlank,
		string(jsonData),
	)
}

func ErrInvalidProposalContentDescrLong(MaxDescriptionLength int) *sdkerrors.Error {
	jsonData, _ := json.Marshal(
		map[string]string{
			"code":                 getCodeString(CodeInvalidProposalContentDescrLong),
			"codespace":            DefaultCodespace,
			"desc":                 fmt.Sprintf("proposal description is longer than max length of %d", MaxDescriptionLength),
			"MaxDescriptionLength": fmt.Sprintf("%d", MaxDescriptionLength),
		},
	)
	return sdkerrors.New(
		DefaultCodespace,
		CodeInvalidProposalContentDescrLong,
		string(jsonData),
	)
}

func ErrInvalidProposalType(ProposalType string) *sdkerrors.Error {
	jsonData, _ := json.Marshal(
		map[string]string{
			"code":         getCodeString(CodeInvalidProposalType),
			"codespace":    DefaultCodespace,
			"desc":         fmt.Sprintf("invalid proposal type: %s", ProposalType),
			"ProposalType": ProposalType,
		},
	)
	return sdkerrors.New(
		DefaultCodespace,
		CodeInvalidProposalType,
		string(jsonData),
	)
}

func ErrInvalidVote(option string) *sdkerrors.Error {
	jsonData, _ := json.Marshal(
		map[string]string{
			"code":      getCodeString(CodeInvalidVote),
			"codespace": DefaultCodespace,
			"desc":      fmt.Sprintf("invalid vote option: %s", option),
			"option":    option,
		},
	)
	return sdkerrors.New(
		DefaultCodespace,
		CodeInvalidVote,
		string(jsonData),
	)
}

func ErrInvalidGenesis() *sdkerrors.Error {
	jsonData, _ := json.Marshal(
		map[string]string{
			"code":      getCodeString(CodeInvalidGenesis),
			"codespace": DefaultCodespace,
			"desc":      "initial proposal ID hasn't been set",
		},
	)
	return sdkerrors.New(
		DefaultCodespace,
		CodeInvalidGenesis,
		string(jsonData),
	)
}

func ErrNoProposalHandlerExists() *sdkerrors.Error {
	jsonData, _ := json.Marshal(
		map[string]string{
			"code":      getCodeString(CodeNoProposalHandlerExists),
			"codespace": DefaultCodespace,
			"desc":      "no handler exists for proposal type",
		},
	)
	return sdkerrors.New(
		DefaultCodespace,
		CodeNoProposalHandlerExists,
		string(jsonData),
	)
}

func ErrInvalidStartEndBlocks(StartBlock uint64, EndBlock uint64) *sdkerrors.Error {
	jsonData, _ := json.Marshal(
		map[string]string{
			"code":       getCodeString(CodeInvalidStartEndBlocks),
			"codespace":  DefaultCodespace,
			"desc":       fmt.Sprintf("invalid start or end blocks: start %d,  end %d ", StartBlock, EndBlock),
			"StartBlock": fmt.Sprintf("%d", StartBlock),
			"EndBlock":   fmt.Sprintf("%d", EndBlock),
		},
	)
	return sdkerrors.New(
		DefaultCodespace,
		CodeInvalidStartEndBlocks,
		string(jsonData),
	)
}

func ErrSubmitProposal(error string) *sdkerrors.Error {
	jsonData, _ := json.Marshal(
		map[string]string{
			"code":      getCodeString(CodeSubmitProposal),
			"codespace": DefaultCodespace,
			"desc":      fmt.Sprintf("error submit proposal: %s", error),
			"error":     error,
		},
	)
	return sdkerrors.New(
		DefaultCodespace,
		CodeSubmitProposal,
		string(jsonData),
	)
}

func ErrStartBlock() *sdkerrors.Error {
	jsonData, _ := json.Marshal(
		map[string]string{
			"code":      getCodeString(CodeStartBlock),
			"codespace": DefaultCodespace,
			"desc":      "start block must greater then current block height",
		},
	)
	return sdkerrors.New(
		DefaultCodespace,
		CodeStartBlock,
		string(jsonData),
	)
}

func ErrDurationTooLong() *sdkerrors.Error {
	jsonData, _ := json.Marshal(
		map[string]string{
			"code":                getCodeString(CodeDurationTooLong),
			"codespace":           DefaultCodespace,
			"desc":                "start block must greater then current block height",
			"maxDurationInMonth":  fmt.Sprintf("%d", DurationInMonth),
			"maxDurationInBlocks": fmt.Sprintf("%d", DurationInBlocks),
		},
	)
	return sdkerrors.New(
		DefaultCodespace,
		CodeDurationTooLong,
		string(jsonData),
	)
}

func ErrNotAllowed() *sdkerrors.Error {
	jsonData, _ := json.Marshal(
		map[string]string{
			"code":      getCodeString(CodeNotAllowed),
			"codespace": DefaultCodespace,
			"desc":      fmt.Sprintf("not allowed to create the proposal from this address"),
		},
	)
	return sdkerrors.New(
		DefaultCodespace,
		CodeNotAllowed,
		string(jsonData),
	)
}

func getCodeString(code CodeType) string {
	return strconv.FormatInt(int64(code), 10)
}
