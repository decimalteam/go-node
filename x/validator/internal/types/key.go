package types

import (
	"encoding/binary"
	"math/big"
	"strconv"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// ModuleName is the name of the module
	ModuleName = "validator"

	// StoreKey to be used when creating the KVStore
	StoreKey = ModuleName

	// TStoreKey is the string transient store representation
	TStoreKey = "transient_" + ModuleName

	// RouterKey to be used for routing msgs
	RouterKey = ModuleName

	QuerierRoute = ModuleName

	LastTotalPowerKey                = 0x02
	ValidatorsByConsAddrKey          = 0x03
	ValidatorsByPowerIndexKey        = 0x04
	ValidatorQueueKey                = 0x05
	LastValidatorPowerKey            = 0x06
	ValidatorsKey                    = 0x07
	DelegationKey                    = 0x08
	UnbondingDelegationKey           = 0x09
	UnbondingDelegationByValIndexKey = 0x0a
	UnbondingQueueKey                = 0x0b
	RedelegationQueueKey             = 0x0c
	RedelegationKey                  = 0x0d
	RedelegationByValSrcIndexKey     = 0x0e
	RedelegationByValDstIndexKey     = 0x0f
	ValidatorSigningInfoKey          = 0x10
	ValidatorMissedBlockBitArrayKey  = 0x11
	AddrPubkeyRelationKey            = 0x12
	HistoricalInfoKey                = 0x13
)

func GetValidatorKey(addr sdk.ValAddress) []byte {
	return append([]byte{ValidatorsKey}, addr.Bytes()...)
}

// gets the key for the validator with pubkey
// VALUE: validator operator address ([]byte)
func GetValidatorByConsAddrKey(addr sdk.ConsAddress) []byte {
	return append([]byte{ValidatorsByConsAddrKey}, addr.Bytes()...)
}

// TokensToConsensusPower - convert input tokens to potential consensus-engine power
//func TokensToConsensusPower(tokens sdk.Int) (uint64, uint64) {
//	tokens = tokens.Quo(sdk.PowerReduction)
//	consensusPowerHigh := big.NewInt(0).Set(tokens.BigInt()).Rsh(tokens.BigInt(), 64).Uint64()
//	consensusPowerLow := big.NewInt(0).Set(tokens.BigInt()).And(tokens.BigInt(), big.NewInt(0).SetUint64(math.MaxUint64)).Uint64()
//	return consensusPowerHigh, consensusPowerLow
//}

// TokensFromConsensusPower - convert input power to tokens
//func TokensFromConsensusPower(powerHigh uint64, powerLow uint64) sdk.Int {
//	consensusPower := sdk.NewIntFromBigInt(big.NewInt(0).Lsh(big.NewInt(0).SetUint64(powerHigh), 64))
//	consensusPower = consensusPower.Add(sdk.NewIntFromUint64(powerLow))
//	return consensusPower.Mul(sdk.PowerReduction)
//}

// PowerReduction is the amount of staking tokens required for 1 unit of consensus-engine power
var PowerReduction = sdk.NewIntFromBigInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil))

// TokensToConsensusPower - convert input tokens to potential consensus-engine power
func TokensToConsensusPower(tokens sdk.Int) int64 {
	return (tokens.Quo(PowerReduction)).Int64()
}

// TokensFromConsensusPower - convert input power to tokens
func TokensFromConsensusPower(power int64) sdk.Int {
	return sdk.NewInt(power).Mul(PowerReduction)
}

// get the validator by power index.
// Power index is the key used in the power-store, and represents the relative
// power ranking of the validator.
// VALUE: validator operator address ([]byte)
func GetValidatorsByPowerIndexKey(validator Validator, power sdk.Int) []byte {
	consensusPower := TokensToConsensusPower(power)
	consensusPowerBytes := make([]byte, 16)
	binary.BigEndian.PutUint64(consensusPowerBytes[:], uint64(consensusPower))

	powerBytes := consensusPowerBytes
	powerBytesLen := len(powerBytes) // 16

	// key is of format prefix || powerbytes || addrBytes
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

//______________________________________________________________________________

// gets the key for delegator bond with validator
// VALUE: staking/Delegation
func GetDelegationKey(delAddr sdk.AccAddress, valAddr sdk.ValAddress) []byte {
	return append(GetDelegationsKey(delAddr), valAddr.Bytes()...)
}

// gets the prefix for a delegator for all validators
func GetDelegationsKey(delAddr sdk.AccAddress) []byte {
	return append([]byte{DelegationKey}, delAddr.Bytes()...)
}

//______________________________________________________________________________

// gets the key for an unbonding delegation by delegator and validator addr
// VALUE: staking/UnbondingDelegation
func GetUBDKey(delAddr sdk.AccAddress, valAddr sdk.ValAddress) []byte {
	return append(
		GetUBDsKey(delAddr.Bytes()),
		valAddr.Bytes()...)
}

// gets the index-key for an unbonding delegation, stored by validator-index
// VALUE: none (key rearrangement used)
func GetUBDByValIndexKey(delAddr sdk.AccAddress, valAddr sdk.ValAddress) []byte {
	return append(GetUBDsByValIndexKey(valAddr), delAddr.Bytes()...)
}

// rearranges the ValIndexKey to get the UBDKey
func GetUBDKeyFromValIndexKey(IndexKey []byte) []byte {
	addrs := IndexKey[1:] // remove prefix bytes
	if len(addrs) != 2*sdk.AddrLen {
		panic("unexpected key length")
	}
	valAddr := addrs[:sdk.AddrLen]
	delAddr := addrs[sdk.AddrLen:]
	return GetUBDKey(delAddr, valAddr)
}

//______________

// gets the prefix for all unbonding delegations from a delegator
func GetUBDsKey(delAddr sdk.AccAddress) []byte {
	return append([]byte{UnbondingDelegationKey}, delAddr.Bytes()...)
}

// gets the prefix keyspace for the indexes of unbonding delegations for a validator
func GetUBDsByValIndexKey(valAddr sdk.ValAddress) []byte {
	return append([]byte{UnbondingDelegationByValIndexKey}, valAddr.Bytes()...)
}

// gets the prefix for all unbonding delegations from a delegator
func GetUnbondingDelegationTimeKey(timestamp time.Time) []byte {
	bz := sdk.FormatTimeBytes(timestamp)
	return append([]byte{UnbondingQueueKey}, bz...)
}

//________________________________________________________________________________

// gets the key for a redelegation
// VALUE: staking/RedelegationKey
func GetREDKey(delAddr sdk.AccAddress, valSrcAddr, valDstAddr sdk.ValAddress) []byte {
	key := make([]byte, 1+sdk.AddrLen*3)

	copy(key[0:sdk.AddrLen+1], GetREDsKey(delAddr.Bytes()))
	copy(key[sdk.AddrLen+1:2*sdk.AddrLen+1], valSrcAddr.Bytes())
	copy(key[2*sdk.AddrLen+1:3*sdk.AddrLen+1], valDstAddr.Bytes())

	return key
}

// gets the index-key for a redelegation, stored by source-validator-index
// VALUE: none (key rearrangement used)
func GetREDByValSrcIndexKey(delAddr sdk.AccAddress, valSrcAddr, valDstAddr sdk.ValAddress) []byte {
	REDSFromValsSrcKey := GetREDsFromValSrcIndexKey(valSrcAddr)
	offset := len(REDSFromValsSrcKey)

	// key is of the form REDSFromValsSrcKey || delAddr || valDstAddr
	key := make([]byte, len(REDSFromValsSrcKey)+2*sdk.AddrLen)
	copy(key[0:offset], REDSFromValsSrcKey)
	copy(key[offset:offset+sdk.AddrLen], delAddr.Bytes())
	copy(key[offset+sdk.AddrLen:offset+2*sdk.AddrLen], valDstAddr.Bytes())
	return key
}

// gets the index-key for a redelegation, stored by destination-validator-index
// VALUE: none (key rearrangement used)
func GetREDByValDstIndexKey(delAddr sdk.AccAddress, valSrcAddr, valDstAddr sdk.ValAddress) []byte {
	REDSToValsDstKey := GetREDsToValDstIndexKey(valDstAddr)
	offset := len(REDSToValsDstKey)

	// key is of the form REDSToValsDstKey || delAddr || valSrcAddr
	key := make([]byte, len(REDSToValsDstKey)+2*sdk.AddrLen)
	copy(key[0:offset], REDSToValsDstKey)
	copy(key[offset:offset+sdk.AddrLen], delAddr.Bytes())
	copy(key[offset+sdk.AddrLen:offset+2*sdk.AddrLen], valSrcAddr.Bytes())

	return key
}

// GetREDKeyFromValSrcIndexKey rearranges the ValSrcIndexKey to get the REDKey
func GetREDKeyFromValSrcIndexKey(indexKey []byte) []byte {
	// note that first byte is prefix byte
	if len(indexKey) != 3*sdk.AddrLen+1 {
		panic("unexpected key length")
	}
	valSrcAddr := indexKey[1 : sdk.AddrLen+1]
	delAddr := indexKey[sdk.AddrLen+1 : 2*sdk.AddrLen+1]
	valDstAddr := indexKey[2*sdk.AddrLen+1 : 3*sdk.AddrLen+1]

	return GetREDKey(delAddr, valSrcAddr, valDstAddr)
}

// GetREDKeyFromValDstIndexKey rearranges the ValDstIndexKey to get the REDKey
func GetREDKeyFromValDstIndexKey(indexKey []byte) []byte {
	// note that first byte is prefix byte
	if len(indexKey) != 3*sdk.AddrLen+1 {
		panic("unexpected key length")
	}
	valDstAddr := indexKey[1 : sdk.AddrLen+1]
	delAddr := indexKey[sdk.AddrLen+1 : 2*sdk.AddrLen+1]
	valSrcAddr := indexKey[2*sdk.AddrLen+1 : 3*sdk.AddrLen+1]
	return GetREDKey(delAddr, valSrcAddr, valDstAddr)
}

// gets the prefix for all unbonding delegations from a delegator
func GetRedelegationTimeKey(timestamp time.Time) []byte {
	bz := sdk.FormatTimeBytes(timestamp)
	return append([]byte{RedelegationQueueKey}, bz...)
}

//______________

// gets the prefix keyspace for redelegations from a delegator
func GetREDsKey(delAddr sdk.AccAddress) []byte {
	return append([]byte{RedelegationKey}, delAddr.Bytes()...)
}

// gets the prefix keyspace for all redelegations redelegating away from a source validator
func GetREDsFromValSrcIndexKey(valSrcAddr sdk.ValAddress) []byte {
	return append([]byte{RedelegationByValSrcIndexKey}, valSrcAddr.Bytes()...)
}

// gets the prefix keyspace for all redelegations redelegating towards a destination validator
func GetREDsToValDstIndexKey(valDstAddr sdk.ValAddress) []byte {
	return append([]byte{RedelegationByValDstIndexKey}, valDstAddr.Bytes()...)
}

// gets the prefix keyspace for all redelegations redelegating towards a destination validator
// from a particular delegator
func GetREDsByDelToValDstIndexKey(delAddr sdk.AccAddress, valDstAddr sdk.ValAddress) []byte {
	return append(
		GetREDsToValDstIndexKey(valDstAddr),
		delAddr.Bytes()...)
}

// stored by *Consensus* address (not operator address)
func GetValidatorSigningInfoKey(v sdk.ConsAddress) []byte {
	return append([]byte{ValidatorSigningInfoKey}, v.Bytes()...)
}

// stored by *Consensus* address (not operator address)
func GetValidatorMissedBlockBitArrayPrefixKey(v sdk.ConsAddress) []byte {
	return append([]byte{ValidatorMissedBlockBitArrayKey}, v.Bytes()...)
}

// stored by *Consensus* address (not operator address)
func GetValidatorMissedBlockBitArrayKey(v sdk.ConsAddress, i int64) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(i))
	return append(GetValidatorMissedBlockBitArrayPrefixKey(v), b...)
}

// parse the validators operator address from power rank key
func ParseValidatorPowerRankKey(key []byte) (operAddr []byte) {
	powerBytesLen := 16
	if len(key) != 1+powerBytesLen+sdk.AddrLen {
		panic("Invalid validator power rank key length")
	}
	operAddr = sdk.CopyBytes(key[powerBytesLen+1:])
	for i, b := range operAddr {
		operAddr[i] = ^b
	}
	return operAddr
}

//________________________________________________________________________________

// GetHistoricalInfoKey gets the key for the historical info
func GetHistoricalInfoKey(height int64) []byte {
	return append([]byte{HistoricalInfoKey}, []byte(strconv.FormatInt(height, 10))...)
}
