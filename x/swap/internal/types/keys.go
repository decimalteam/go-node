package types

import (
	"bitbucket.org/decimalteam/go-node/x/coin"
	"encoding/binary"
)

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
	SwapKey   = []byte{0x50, 0x01}
	SwapV2Key = []byte{0x50, 0x02}
	ChainKey  = []byte{0x50, 0x03}
)

func GetSwapKey(hash [32]byte) []byte {
	return append(SwapKey, hash[:]...)
}

func GetSwapV2Key(hash [32]byte) []byte {
	return append(SwapV2Key, hash[:]...)
}

func GetChainKey(chain int) []byte {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(chain))
	return append(ChainKey, buf...)
}
