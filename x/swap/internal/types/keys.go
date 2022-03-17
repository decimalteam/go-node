package types

import (
	"encoding/binary"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"bitbucket.org/decimalteam/go-node/utils/updates"
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
)

var (
	SwapKey   = []byte("swap/")
	SwapV2Key = []byte("swap/v2/")
	ChainKey  = []byte("swap/chain/")
)

func GetSwapKey(ctx sdk.Context, hash [32]byte) []byte {
	keyPrefix := SwapKey
	if ctx.BlockHeight() < updates.Update14Block {
		keyPrefix = LegacySwapKey
	}
	return append(keyPrefix, hash[:]...)
}

func GetSwapV2Key(ctx sdk.Context, hash [32]byte) []byte {
	keyPrefix := SwapV2Key
	if ctx.BlockHeight() < updates.Update14Block {
		keyPrefix = LegacySwapV2Key
	}
	return append(keyPrefix, hash[:]...)
}

func GetChainKey(ctx sdk.Context, chain int) []byte {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(chain))
	keyPrefix := ChainKey
	if ctx.BlockHeight() < updates.Update14Block {
		keyPrefix = LegacyChainKey
	}
	return append(keyPrefix, buf...)
}
