package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// ValidatorGovInfo used for tallying
type ValidatorGovInfo struct {
	Address      sdk.ValAddress // address of the validator operator
	BondedTokens sdk.Int        // Power of a Validator
	Vote         VoteOption     // Vote of the validator
}

// NewValidatorGovInfo creates a ValidatorGovInfo instance
func NewValidatorGovInfo(address sdk.ValAddress, bondedTokens sdk.Int, vote VoteOption) ValidatorGovInfo {

	return ValidatorGovInfo{
		Address:      address,
		BondedTokens: bondedTokens,
		Vote:         vote,
	}
}

// TallyResult defines a standard tally for a proposal
type TallyResult struct {
	Yes     sdk.Int `json:"yes" yaml:"yes"`
	Abstain sdk.Int `json:"abstain" yaml:"abstain"`
	No      sdk.Int `json:"no" yaml:"no"`
}

// NewTallyResult creates a new TallyResult instance
func NewTallyResult(yes, abstain, no sdk.Int) TallyResult {
	return TallyResult{
		Yes:     yes,
		Abstain: abstain,
		No:      no,
	}
}

// NewTallyResultFromMap creates a new TallyResult instance from a Option -> Dec map
func NewTallyResultFromMap(results map[VoteOption]sdk.Dec) TallyResult {
	return NewTallyResult(
		results[OptionYes].TruncateInt(),
		results[OptionAbstain].TruncateInt(),
		results[OptionNo].TruncateInt(),
	)
}

// EmptyTallyResult returns an empty TallyResult.
func EmptyTallyResult() TallyResult {
	return NewTallyResult(sdk.ZeroInt(), sdk.ZeroInt(), sdk.ZeroInt())
}

// Equals returns if two proposals are equal.
func (tr TallyResult) Equals(comp TallyResult) bool {
	return tr.Yes.Equal(comp.Yes) &&
		tr.Abstain.Equal(comp.Abstain) &&
		tr.No.Equal(comp.No)
}

// String implements stringer interface
func (tr TallyResult) String() string {
	return fmt.Sprintf(`Tally Result:
  Yes:        %s
  Abstain:    %s
  No:         %s`, tr.Yes, tr.Abstain, tr.No)
}
