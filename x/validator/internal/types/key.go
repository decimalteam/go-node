package types

import (
	"encoding/binary"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"time"
)

const (
	// ModuleName is the name of the module
	ModuleName = "validator"

	// StoreKey to be used when creating the KVStore
	StoreKey = ModuleName

	// RouterKey to be used for routing msgs
	RouterKey = ModuleName

	QuerierRoute = ModuleName

	ValidatorKey              = 0x01
	LastTotalPowerKey         = 0x02
	ValidatorsByConsAddrKey   = 0x03
	ValidatorsByPowerIndexKey = 0x04
	ValidatorQueueKey         = 0x05
	LastValidatorPowerKey     = 0x06
	ValidatorsKey             = 0x07
)

func GetValidatorKey(addr sdk.ValAddress) []byte {
	return append([]byte{ValidatorKey}, addr.Bytes()...)
}

// gets the key for the validator with pubkey
// VALUE: validator operator address ([]byte)
func GetValidatorByConsAddrKey(addr sdk.ConsAddress) []byte {
	return append([]byte{ValidatorsByConsAddrKey}, addr.Bytes()...)
}

// get the validator by power index.
// Power index is the key used in the power-store, and represents the relative
// power ranking of the validator.
// VALUE: validator operator address ([]byte)
func GetValidatorsByPowerIndexKey(validator Validator, power sdk.Int) []byte {
	consensusPower := sdk.TokensToConsensusPower(power)
	consensusPowerBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(consensusPowerBytes, uint64(consensusPower))

	powerBytes := consensusPowerBytes
	powerBytesLen := len(powerBytes) // 8

	// key is of format prefix || powerBytes || addrBytes
	key := make([]byte, 1+powerBytesLen+sdk.AddrLen)

	key[0] = ValidatorsByPowerIndexKey
	copy(key[1:powerBytesLen+1], powerBytes)
	operAddrInvr := sdk.CopyBytes(validator.ValAddress)
	for i, b := range operAddrInvr {
		operAddrInvr[i] = ^b
	}
	copy(key[powerBytesLen+1:], operAddrInvr)

	return key
}

// gets the prefix for all unbonding delegations from a delegator
func GetValidatorQueueTimeKey(timestamp time.Time) []byte {
	bz := sdk.FormatTimeBytes(timestamp)
	return append([]byte{ValidatorQueueKey}, bz...)
}

// get the bonded validator index key for an operator address
func GetLastValidatorPowerKey(valAddress sdk.ValAddress) []byte {
	return append([]byte{LastValidatorPowerKey}, valAddress...)
}

// Get the validator operator address from LastValidatorPowerKey
func AddressFromLastValidatorPowerKey(key []byte) []byte {
	return key[1:] // remove prefix bytes
}
