package config

import sdk "github.com/cosmos/cosmos-sdk/types"

const (
	// Change this params
	ChainID = "decimal"
	//

	TitleTestBaseCoin  = "Test decimal coin"
	SymbolTestBaseCoin = "tDCL"
	TitleBaseCoin      = "Decimal coin"
	SymbolBaseCoin     = "DCL"
)

var (
	InitialVolumeTestBaseCoin = sdk.NewInt(250000000000)
	InitialVolumeBaseCoin     = sdk.NewInt(1000000000000)
)

type Config struct {
	TitleBaseCoin         string  `json:"title" yaml:"title"`   // Full coin title (Bitcoin)
	SymbolBaseCoin        string  `json:"symbol" yaml:"symbol"` // Short coin title (BTC)
	InitialVolumeBaseCoin sdk.Int `json:"initial_volume" yaml:"initial_volume"`
}

func GetDefaultConfig(chainId string) *Config {
	cnf := Config{}
	if chainId == "decimal-testnet" {
		cnf.TitleBaseCoin = TitleTestBaseCoin
		cnf.SymbolBaseCoin = SymbolTestBaseCoin
		cnf.InitialVolumeBaseCoin = InitialVolumeTestBaseCoin
		return &cnf
	} else if chainId == "decimal" {
		cnf.TitleBaseCoin = TitleBaseCoin
		cnf.SymbolBaseCoin = SymbolBaseCoin
		cnf.InitialVolumeBaseCoin = InitialVolumeBaseCoin
		return &cnf
	} else {
		return &cnf
	}
}
