package types

import (
	"regexp"

	"bitbucket.org/decimalteam/go-node/config"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GenesisState - all coin state that must be provided at genesis
type GenesisState struct {
	Title         string  `json:"title" yaml:"title"`   // Full coin title (Bitcoin)
	Symbol        string  `json:"symbol" yaml:"symbol"` // Short coin title (BTC)
	InitialVolume sdk.Int `json:"initial_volume" yaml:"initial_volume"`
	Coins         []Coin  `json:"coins" yaml:"coins"` // custom coins in store
}

// NewGenesisState creates a new GenesisState object
func NewGenesisState(title string, symbol string, initVolume sdk.Int, coins []Coin) GenesisState {
	return GenesisState{
		Title:         title,
		Symbol:        symbol,
		InitialVolume: initVolume,
		Coins:         coins,
	}
}

// DefaultGenesisState - default GenesisState used by Cosmos Hub
func DefaultGenesisState() GenesisState {
	return GenesisState{
		Title:         config.TitleBaseCoin,
		Symbol:        config.SymbolBaseCoin,
		InitialVolume: config.InitialVolumeBaseCoin,
		Coins:         []Coin{},
	}
}

// ValidateGenesis validates the coin genesis parameters
func ValidateGenesis(data GenesisState) error {
	// Check coin title maximum length
	if len(data.Title) > maxCoinNameBytes {
		return ErrInvalidCoinTitle(data.Title)
	}
	// Check coin symbol for correct regexp
	if match, _ := regexp.MatchString(allowedCoinSymbols, data.Symbol); !match {
		return ErrInvalidCoinSymbol(data.Symbol)
	}
	// Check coin initial volume to be correct
	if data.InitialVolume.LT(minCoinSupply) || data.InitialVolume.GT(maxCoinSupply) {
		return ErrInvalidCoinInitialVolume(data.InitialVolume.String())
	}
	return nil
}
