package types

import sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

// x/gov module sentinel errors
var (
	ErrUnknownProposal         = sdkerrors.Register(ModuleName, 100, "unknown proposal")
	ErrInactiveProposal        = sdkerrors.Register(ModuleName, 200, "inactive proposal")
	ErrAlreadyActiveProposal   = sdkerrors.Register(ModuleName, 300, "proposal already active")
	ErrInvalidProposalContent  = sdkerrors.Register(ModuleName, 400, "invalid proposal content")
	ErrInvalidProposalType     = sdkerrors.Register(ModuleName, 500, "invalid proposal type")
	ErrInvalidVote             = sdkerrors.Register(ModuleName, 600, "invalid vote option")
	ErrInvalidGenesis          = sdkerrors.Register(ModuleName, 700, "invalid genesis state")
	ErrNoProposalHandlerExists = sdkerrors.Register(ModuleName, 800, "no handler exists for proposal type")
	ErrInvalidStartEndBlocks   = sdkerrors.Register(ModuleName, 900, "invalid start or end blocks")
	ErrSubmitProposal          = sdkerrors.Register(ModuleName, 1000, "error submit proposal")
	ErrStartBlock              = sdkerrors.Register(ModuleName, 1100, "start block must greater then current block height")
	ErrDurationTooLong         = sdkerrors.Register(ModuleName, 1200, "duration too long. Max duration = 3 month (1296000 blocks)")
	ErrNotAllowed              = sdkerrors.Register(ModuleName, 1300, "not allowed to create the proposal from this address")
)
