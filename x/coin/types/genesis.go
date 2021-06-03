package types

import (
	"regexp"

	"bitbucket.org/decimalteam/go-node/config"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GenesisState - all coin state that must be provided at genesis

// NewGenesisState creates a new GenesisState object
func NewGenesisState(title string, symbol string, initVolume sdk.Int) GenesisState {
	return GenesisState{
		Title:         title,
		Symbol:        symbol,
		InitialVolume: initVolume,
	}
}

// DefaultGenesisState - default GenesisState used by Cosmos Hub
func DefaultGenesisState() GenesisState {
	return GenesisState{
		Title:         config.TitleBaseCoin,
		Symbol:        config.SymbolBaseCoin,
		InitialVolume: config.InitialVolumeBaseCoin,
	}
}

// ValidateGenesis validates the coin genesis parameters
func ValidateGenesis(data GenesisState) error {
	// Check coin title maximum length
	if len(data.Title) > maxCoinNameBytes {
		return ErrInvalidCoinTitle()
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
