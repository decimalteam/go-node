package types

import (
	"encoding/binary"

	"bitbucket.org/decimalteam/go-node/x/coin"
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
	LegacySwapKey   = []byte{0x50, 0x01}
	LegacySwapV2Key = []byte{0x50, 0x02}
	LegacyChainKey  = []byte{0x50, 0x03}

	// This is special key used to determine if kv-records are migrated to keys with correct prefixes
	LegacyMigrationKey = []byte("swap/migrated")
)

var (
	SwapKey   = []byte("swap/v1/")
	SwapV2Key = []byte("swap/v2/")
	ChainKey  = []byte("swap/chain/")
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
