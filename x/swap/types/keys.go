package types

const (
	// ModuleName is the name of the module
	ModuleName = "swap"

	// RouterKey is the message route for swap
	RouterKey = ModuleName

	// StoreKey to be used when creating the KVStore
	StoreKey = ModuleName

	QuerierRoute = ModuleName
)

var (
	SwapKey = []byte{0x01}
)

func GetSwapKey(hash [32]byte) []byte {
	return append(SwapKey, hash[:]...)
}
