package types

import "bitbucket.org/decimalteam/go-node/x/coin"

const (
	// ModuleName is the name of the module
	ModuleName = "swap"

	// RouterKey is the message route for swap
	RouterKey = ModuleName

	// StoreKey to be used when creating the KVStore
	StoreKey = coin.StoreKey

	QuerierRoute = ModuleName
)

var (
	SwapKey = []byte{0x50, 0x01}
)

func GetSwapKey(hash [32]byte) []byte {
	return append(SwapKey, hash[:]...)
}
