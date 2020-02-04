package config

import sdk "github.com/cosmos/cosmos-sdk/types"

const (

	// Change this params
	ChainID = "decimal-testnet"

	TitleTestBaseCoin  = "tDCL"
	SymbolTestBaseCoin = "Test decimal coin"
	TitleBaseCoin      = "DCL"
	SymbolBaseCoin     = "Decimal coin"
)

var (
	InitialVolumeTestBaseCoin = sdk.NewInt(1000000000000)
	InitialVolumeBaseCoin     = sdk.NewInt(1000000000000)
)

type Config struct {
	TitleBaseCoin         string  `json:"title" yaml:"title"`   // Full coin title (Bitcoin)
	SymbolBaseCoin        string  `json:"symbol" yaml:"symbol"` // Short coin title (BTC)
	InitialVolumeBaseCoin sdk.Int `json:"initial_volume" yaml:"initial_volume"`
}

func DefaultConfig(cnf *Config) *Config {
	if ChainID == "testnet" {
		cnf.TitleBaseCoin = TitleTestBaseCoin
		cnf.SymbolBaseCoin = SymbolTestBaseCoin
		cnf.InitialVolumeBaseCoin = InitialVolumeTestBaseCoin
		return cnf
	} else {
		cnf.TitleBaseCoin = TitleBaseCoin
		cnf.SymbolBaseCoin = SymbolBaseCoin
		cnf.InitialVolumeBaseCoin = InitialVolumeBaseCoin
		return cnf
	}

}
