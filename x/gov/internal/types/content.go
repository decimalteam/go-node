package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"strings"
)

// Constants pertaining to a Content object
const (
	MaxDescriptionLength int = 5000
	MaxTitleLength       int = 140
)

type Content struct {
	Title       string `json:"title" yaml:"title"`             // Proposal title
	Description string `json:"description" yaml:"description"` // Proposal description
}

func (c *Content) GetTitle() string       { return c.Title }
func (c *Content) GetDescription() string { return c.Description }

// Handler defines a function that handles a proposal after it has passed the
// governance process.
type Handler func(ctx sdk.Context, content Content) error

// Validate validates a proposal's abstract contents returning an error
// if invalid.
func Validate(c Content) error {
	title := c.GetTitle()
	if len(strings.TrimSpace(title)) == 0 {
		return ErrInvalidProposalContentTitleBlank()
	}
	if len(title) > MaxTitleLength {
		return ErrInvalidProposalContentTitleLong(MaxTitleLength)
	}

	description := c.GetDescription()
	if len(description) == 0 {
		return ErrInvalidProposalContentDescrBlank()
	}
	if len(description) > MaxDescriptionLength {
		return ErrInvalidProposalContentDescrLong(MaxDescriptionLength)
	}

	return nil
}
