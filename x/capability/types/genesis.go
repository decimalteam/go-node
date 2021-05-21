package types

// GenesisState - all capability state that must be provided at genesis
type GenesisState struct {
	Index int64
}

// NewGenesisState creates a new GenesisState object
func NewGenesisState() GenesisState {

	return GenesisState{}
}

// DefaultGenesisState - default GenesisState used by Cosmos Hub
func DefaultGenesisState() GenesisState {
	return GenesisState{}
}

// ValidateGenesis validates the capability genesis parameters
func ValidateGenesis(data GenesisState) error {
	return nil
}
