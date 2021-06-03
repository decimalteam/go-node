package types

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/x/params/types"
)

// Default parameter namespace
const (
	DefaultParamspace = ModuleName
	// TODO: Define your default parameters
)

// Parameter store keys
var (
// TODO: Define your keys for the parameter store
// KeyParamName          = []byte("ParamName")
)

// ParamKeyTable for multisig module
func ParamKeyTable() types.KeyTable {
	return types.NewKeyTable().RegisterParamSet(&Params{})
}

// Params - used for initializing default parameter for multisig at genesis
type Params struct {
	// TODO: Add your Paramaters to the Paramter struct
	// KeyParamName string `json:"key_param_name"`
}

// NewParams creates a new Params object
func NewParams( /* TODO: Pass in the paramters*/ ) Params {

	return Params{
		// TODO: Create your Params Type
	}
}

// String implements the stringer interface for Params
func (p Params) String() string {
	return fmt.Sprintf(`
	// TODO: Return all the params as a string
	`)
}

// ParamSetPairs - Implements params.ParamSet
func (p *Params) ParamSetPairs() types.ParamSetPairs {
	return types.ParamSetPairs{
		// TODO: Pair your key with the param
		// params.NewParamSetPair(KeyParamName, &p.ParamName),
	}
}

// DefaultParams defines the parameters for this module
func DefaultParams() Params {
	return NewParams(
	// TODO: Pass in your default Params
	)
}
