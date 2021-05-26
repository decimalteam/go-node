package types

import (
	"encoding/binary"
	"strconv"
	"time"

	"bitbucket.org/decimalteam/go-node/types"

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
	ValidatorSigningInfoKey          = 0x10
	ValidatorMissedBlockBitArrayKey  = 0x11
	AddrPubkeyRelationKey            = 0x12
	HistoricalInfoKey                = 0x13
	DelegationNFTKey                 = 0x14
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

// get the validator by power index.
// Power index is the key used in the power-store, and represents the relative
// power ranking of the validator.
// VALUE: validator operator address ([]byte)
func GetValidatorsByPowerIndexKey(validator Validator, power sdk.Int) []byte {
	consensusPower := types.TokensToConsensusPower(power)
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
func GetDelegationKey(delAddr sdk.AccAddress, valAddr sdk.ValAddress, coin string) []byte {
	return append(append(append([]byte{DelegationKey}, delAddr.Bytes()...), valAddr.Bytes()...), []byte(coin)...)
}

// gets the prefix for a delegator for all validators
func GetDelegationsKey(delAddr sdk.AccAddress) []byte {
	return append([]byte{DelegationKey}, delAddr.Bytes()...)
}

//______________________________________________________________________________

func GetDelegationNFTKey(delAddr sdk.AccAddress, valAddr sdk.ValAddress, tokenID, denom string) []byte {
	return append(append(append(append([]byte{DelegationNFTKey}, delAddr.Bytes()...), valAddr.Bytes()...), []byte(tokenID)...), []byte(denom)...)
}

func GetDelegationsNFTKey(delAddr sdk.AccAddress) []byte {
	return append([]byte{DelegationNFTKey}, delAddr.Bytes()...)
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
