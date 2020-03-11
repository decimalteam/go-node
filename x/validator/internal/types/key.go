package types

import sdk "github.com/cosmos/cosmos-sdk/types"

const (
	// ModuleName is the name of the module
	ModuleName = "validator"

	// StoreKey to be used when creating the KVStore
	StoreKey = ModuleName

	// RouterKey to be used for routing msgs
	RouterKey = ModuleName

	QuerierRoute = ModuleName

	ValidatorKey = 0x01
)

func GetValidatorKey(addr sdk.ValAddress) []byte {
	return append([]byte(byte(ValidatorKey)), addr.Bytes()...)
}
