package types

import (
	"bytes"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GenesisState - all staking state that must be provided at genesis
type GenesisState struct {
	StartingProposalID uint64      `json:"starting_proposal_id" yaml:"starting_proposal_id"`
	Votes              Votes       `json:"votes" yaml:"votes"`
	Proposals          Proposals   `json:"proposals" yaml:"proposals"`
	TallyParams        TallyParams `json:"tally_params" yaml:"tally_params"`
}

// NewGenesisState creates a new genesis state for the governance module
func NewGenesisState(startingProposalID uint64, tp TallyParams) GenesisState {
	return GenesisState{
		StartingProposalID: startingProposalID,
		TallyParams:        tp,
	}
}

// DefaultGenesisState defines the default governance genesis state
func DefaultGenesisState() GenesisState {
	return NewGenesisState(
		DefaultStartingProposalID,
		DefaultTallyParams(),
	)
}

// Equal checks whether two gov GenesisState structs are equivalent
func (data GenesisState) Equal(data2 GenesisState) bool {
	b1 := ModuleCdc.MustMarshalBinaryBare(data)
	b2 := ModuleCdc.MustMarshalBinaryBare(data2)
	return bytes.Equal(b1, b2)
}

// IsEmpty returns true if a GenesisState is empty
func (data GenesisState) IsEmpty() bool {
	return data.Equal(GenesisState{})
}

// ValidateGenesis checks if parameters are within valid ranges
func ValidateGenesis(data GenesisState) error {
	threshold := data.TallyParams.Threshold
	if threshold.IsNegative() || threshold.GT(sdk.OneDec()) {
		return fmt.Errorf("governance vote threshold should be positive and less or equal to one, is %s",
			threshold.String())
	}

	return nil
}
