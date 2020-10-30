package types

type GenesisState struct {
	Swaps  Swaps  `json:"swaps" yaml:"swaps"`
	Params Params `json:"params" yaml:"params"`
}

func NewGenesisState(params Params, swaps Swaps) GenesisState {
	return GenesisState{
		Swaps:  swaps,
		Params: params,
	}
}

// get raw genesis raw message for testing
func DefaultGenesisState() GenesisState {
	return GenesisState{
		Params: DefaultParams(),
	}
}
