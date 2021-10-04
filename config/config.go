package config

import (
	"fmt"
	"os"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (

	// DecimalVersion is integer version of the Decimal app.
	DecimalVersion = "1.3.9"
	// DecimalMainPrefix is the main prefix for all keys and addresses.
	DecimalMainPrefix = "dx"

	// PrefixValidator is the prefix for validator keys.
	PrefixValidator = "val"
	// PrefixConsensus is the prefix for consensus keys.
	PrefixConsensus = "cons"
	// PrefixPublic is the prefix for public keys.
	PrefixPublic = "pub"
	// PrefixOperator is the prefix for operator keys.
	PrefixOperator = "oper"

	// DecimalPrefixAccAddr defines the Decimal prefix of an account's address.
	DecimalPrefixAccAddr = DecimalMainPrefix
	// DecimalPrefixAccPub defines the Decimal prefix of an account's public key.
	DecimalPrefixAccPub = DecimalMainPrefix + PrefixPublic
	// DecimalPrefixValAddr defines the Decimal prefix of a validator's operator address.
	DecimalPrefixValAddr = DecimalMainPrefix + PrefixValidator + PrefixOperator
	// DecimalPrefixValPub defines the Decimal prefix of a validator's operator public key.
	DecimalPrefixValPub = DecimalMainPrefix + PrefixValidator + PrefixOperator + PrefixPublic
	// DecimalPrefixConsAddr defines the Decimal prefix of a consensus node address.
	DecimalPrefixConsAddr = DecimalMainPrefix + PrefixValidator + PrefixConsensus
	// DecimalPrefixConsPub defines the Decimal prefix of a consensus node public key.
	DecimalPrefixConsPub = DecimalMainPrefix + PrefixValidator + PrefixConsensus + PrefixPublic

	TitleTestBaseCoin  = "Test decimal coin"
	SymbolTestBaseCoin = "tdel"
	TitleBaseCoin      = "Decimal coin"
	SymbolBaseCoin     = "del"
)

const (
	SkipPlanName = "skip_plans.json"
	OneHour      = 660 // blocks
)

var (
	ConfigPath = fmt.Sprintf("%s/.decimal/daemon/config", os.Getenv("HOME"))
)

var (
	InitialVolumeTestBaseCoin, _ = sdk.NewIntFromString("340000000000000000000000000")
	InitialVolumeBaseCoin, _     = sdk.NewIntFromString("340000000000000000000000000")
)

var (
	ChainID                  = "decimal-testnet-06-09-13-00"
	WithoutSlashPeriodPrefix = []byte{0x60}
)

type Config struct {
	Initialized           bool    `json:"initialized" yaml:"initialized"`
	TitleBaseCoin         string  `json:"title" yaml:"title"`   // Full coin title (Bitcoin)
	SymbolBaseCoin        string  `json:"symbol" yaml:"symbol"` // Short coin title (BTC)
	InitialVolumeBaseCoin sdk.Int `json:"initial_volume" yaml:"initial_volume"`
}

func GetDefaultConfig(chainId string) *Config {
	cnf := Config{}
	if strings.HasPrefix(chainId, "decimal-testnet") {
		cnf.TitleBaseCoin = TitleTestBaseCoin
		cnf.SymbolBaseCoin = SymbolTestBaseCoin
		cnf.InitialVolumeBaseCoin = InitialVolumeTestBaseCoin
		return &cnf
	} else if strings.HasPrefix(chainId, "decimal") {
		cnf.TitleBaseCoin = TitleBaseCoin
		cnf.SymbolBaseCoin = SymbolBaseCoin
		cnf.InitialVolumeBaseCoin = InitialVolumeBaseCoin
		return &cnf
	} else {
		return &cnf
	}
}
