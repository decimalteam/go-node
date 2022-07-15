package types

// GenesisState - all multisig state that must be provided at genesis
type GenesisState struct {
	Wallets []Wallet      `json:"wallets"`
	Txs     []Transaction `json:"txs"`
}

// NewGenesisState creates a new GenesisState object
func NewGenesisState(wallets []Wallet, txs []Transaction) GenesisState {
	return GenesisState{
		Wallets: wallets,
		Txs:     txs,
	}
}

// DefaultGenesisState - default GenesisState used by Cosmos Hub
func DefaultGenesisState() GenesisState {
	return GenesisState{
		Wallets: make([]Wallet, 0),
		Txs:     make([]Transaction, 0),
	}
}

// ValidateGenesis validates the multisig genesis parameters
func ValidateGenesis(data GenesisState) error {
	// TODO: Create a sanity check to make sure the state conforms to the modules needs
	return nil
}
