package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params/types"
)

// Default governance params
var (
	DefaultQuorum    = sdk.NewDecWithPrec(667, 3)
	DefaultThreshold = sdk.NewDecWithPrec(5, 1)
)

// Parameter store key
var (
	ParamStoreKeyTallyParams = []byte("tallyparams")
)

// ParamKeyTable - Key declaration for parameters
func ParamKeyTable() types.KeyTable {
	return types.NewKeyTable(
		types.NewParamSetPair(ParamStoreKeyTallyParams, TallyParams{}, validateTallyParams),
	)
}

// TallyParams defines the params around Tallying votes in governance
//type TallyParams struct {
//	Quorum    sdk.Dec `json:"quorum,omitempty" yaml:"quorum,omitempty"`       //  Minimum percentage of total stake needed to vote for a result to be considered valid
//	Threshold sdk.Dec `json:"threshold,omitempty" yaml:"threshold,omitempty"` //  Minimum proportion of Yes votes for proposal to pass. Initial value: 0.5
//}

// NewTallyParams creates a new TallyParams object
func NewTallyParams(quorum, threshold sdk.Dec) TallyParams {
	return TallyParams{
		Quorum:    quorum,
		Threshold: threshold,
	}
}

// DefaultTallyParams default parameters for tallying
func DefaultTallyParams() TallyParams {
	return NewTallyParams(DefaultQuorum, DefaultThreshold)
}

// String implements stringer insterface
func (tp TallyParams) String() string {
	return fmt.Sprintf(`Tally Params:
  Quorum:             %s
  Threshold:          %s`,
		tp.Quorum, tp.Threshold)
}

func validateTallyParams(i interface{}) error {
	v, ok := i.(TallyParams)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.Quorum.IsNegative() {
		return fmt.Errorf("quorom cannot be negative: %s", v.Quorum)
	}
	if v.Quorum.GT(sdk.OneDec()) {
		return fmt.Errorf("quorom too large: %s", v)
	}
	if !v.Threshold.IsPositive() {
		return fmt.Errorf("vote threshold must be positive: %s", v.Threshold)
	}
	if v.Threshold.GT(sdk.OneDec()) {
		return fmt.Errorf("vote threshold too large: %s", v)
	}

	return nil
}

// Params returns all of the governance params
type Params struct {
	TallyParams TallyParams `json:"tally_params" yaml:"tally_params"`
}

func (gp Params) String() string {
	return gp.TallyParams.String()
}

// NewParams creates a new gov Params instance
func NewParams(tp TallyParams) Params {
	return Params{
		TallyParams: tp,
	}
}

// DefaultParams default governance params
func DefaultParams() Params {
	return NewParams(DefaultTallyParams())
}
