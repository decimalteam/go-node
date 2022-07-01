package types

type GenesisState struct {
	Swaps  Swaps   `json:"swaps" yaml:"swaps"`
	Params Params  `json:"params" yaml:"params"`
	Chains []Chain `json:"chains" yaml:"chains"`
}

func NewGenesisState(params Params, swaps Swaps, chains []Chain) GenesisState {
	return GenesisState{
		Swaps:  swaps,
		Params: params,
		Chains: chains,
	}
}

// get raw genesis raw message for testing
func DefaultGenesisState() GenesisState {
	return GenesisState{
		Params: DefaultParams(),
	}
}
