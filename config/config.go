package config

import sdk "github.com/cosmos/cosmos-sdk/types"

const (
	DecimalMainPrefix = "dx"

	// PrefixValidator is the prefix for validator keys
	PrefixValidator = "val"
	// PrefixConsensus is the prefix for consensus keys
	PrefixConsensus = "cons"
	// PrefixPublic is the prefix for public keys
	PrefixPublic = "pub"
	// PrefixOperator is the prefix for operator keys
	PrefixOperator = "oper"

	// DecimalPrefixAccAddr defines the Decimal prefix of an account's address
	DecimalPrefixAccAddr = DecimalMainPrefix
	// DecimalPrefixAccPub defines the Decimal prefix of an account's public key
	DecimalPrefixAccPub = DecimalMainPrefix + PrefixPublic
	// DecimalPrefixValAddr defines the Decimal prefix of a validator's operator address
	DecimalPrefixValAddr = DecimalMainPrefix + PrefixValidator + PrefixOperator
	// DecimalPrefixValPub defines the Decimal prefix of a validator's operator public key
	DecimalPrefixValPub = DecimalMainPrefix + PrefixValidator + PrefixOperator + PrefixPublic
	// DecimalPrefixConsAddr defines the Decimal prefix of a consensus node address
	DecimalPrefixConsAddr = DecimalMainPrefix + PrefixValidator + PrefixConsensus
	// DecimalPrefixConsPub defines the Decimal prefix of a consensus node public key
	DecimalPrefixConsPub = DecimalMainPrefix + PrefixValidator + PrefixConsensus + PrefixPublic

	// Change this params
	ChainID = "decimal-testnet"
	//

	TitleTestBaseCoin  = "Test decimal coin"
	SymbolTestBaseCoin = "tDEL"
	TitleBaseCoin      = "Decimal coin"
	SymbolBaseCoin     = "DEL"

	// test coin
	TitleTestCoin  = "Crypton coin"
	SymbolTestCoin = "CRT"

	// Check config
	DecimalCheckPrefix = "Dc"
)

var (
	InitialVolumeTestBaseCoin, _ = sdk.NewIntFromString("220000000000000000000000000")
	InitialVolumeBaseCoin, _     = sdk.NewIntFromString("220000000000000000000000000")

	// test params buy
	InitialReserveTestCoin, _ = sdk.NewIntFromString("120798840222697144373637")   //    120798.840222697144373637
	InitialVolumeTestCoin, _  = sdk.NewIntFromString("54363225921077956709926174") //  54363225.921077956709926174
	CRRTestCoin               = 88

	// test params sell
	//InitialReserveTestCoin, _ = sdk.NewIntFromString("86177720949431621141039204") //86177720.949431621141039204
	//InitialVolumeTestCoin, _  = sdk.NewIntFromString("19735598708313902262960810") //19735598.708313902262960810
	//CRRTestCoin               = 75
)

type Config struct {
	TitleBaseCoin         string  `json:"title" yaml:"title"`   // Full coin title (Bitcoin)
	SymbolBaseCoin        string  `json:"symbol" yaml:"symbol"` // Short coin title (BTC)
	InitialVolumeBaseCoin sdk.Int `json:"initial_volume" yaml:"initial_volume"`

	//test
	TitleTestCoin            string  `json:"title" yaml:"title"`   // Full coin title (Bitcoin)
	SymbolTestCoin           string  `json:"symbol" yaml:"symbol"` // Short coin title (BTC)
	InitialVolumeTestCoin    sdk.Int `json:"initial_volume" yaml:"initial_volume"`
	InitialReserveTestCoin   sdk.Int `json:"initial_volume" yaml:"initial_volume"`
	ConstantReserveRatioTest uint    `json:"constant_reserve_ratio" yaml:"constant_reserve_ratio"` // between 10 and 100
}

func GetDefaultConfig(chainId string) *Config {
	cnf := Config{}
	if chainId == "decimal-testnet" {
		cnf.TitleBaseCoin = TitleTestBaseCoin
		cnf.SymbolBaseCoin = SymbolTestBaseCoin
		cnf.InitialVolumeBaseCoin = InitialVolumeTestBaseCoin

		//test
		cnf.TitleTestCoin = TitleTestCoin
		cnf.SymbolTestCoin = SymbolTestCoin
		cnf.InitialVolumeTestCoin = InitialVolumeTestCoin
		cnf.InitialReserveTestCoin = InitialReserveTestCoin
		cnf.ConstantReserveRatioTest = uint(CRRTestCoin)

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
