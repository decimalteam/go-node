package types

const (
	// ModuleName is the name of the module
	ModuleName = "coin"

	// StoreKey to be used when creating the KVStore
	StoreKey = ModuleName

	// RouterKey to be used for routing msgs
	RouterKey = ModuleName

	QuerierRoute = ModuleName

	CoinPrefix  = "coin-"
	CheckPrefix = "check-"
)
