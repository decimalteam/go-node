package types

import (
	"bitbucket.org/decimalteam/go-node/config"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"regexp"
)

// GenesisState - all coin state that must be provided at genesis
type GenesisState struct {
	Title         string  `json:"title" yaml:"title"`   // Full coin title (Bitcoin)
	Symbol        string  `json:"symbol" yaml:"symbol"` // Short coin title (BTC)
	InitialVolume sdk.Int `json:"initial_volume" yaml:"initial_volume"`
}

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
		return sdk.NewError(DefaultCodespace, InvalidCoinTitle, fmt.Sprintf("Coin name is invalid. Allowed up to %d bytes.", maxCoinNameBytes))
	}
	// Check coin symbol for correct regexp
	if match, _ := regexp.MatchString(allowedCoinSymbols, data.Symbol); !match {
		return sdk.NewError(DefaultCodespace, InvalidCoinSymbol, fmt.Sprintf("Invalid coin symbol. Should be %s", allowedCoinSymbols))
	}
	// Check coin initial volume to be correct
	if data.InitialVolume.LT(minCoinSupply) || data.InitialVolume.GT(maxCoinSupply) {
		return sdk.NewError(DefaultCodespace, InvalidCoinInitVolume, fmt.Sprintf("Coin initial volume should be between %s and %s. Given %s", minCoinSupply.String(), maxCoinSupply.String(), data.InitialVolume.String()))
	}
	return nil
}
