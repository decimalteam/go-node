package types

import (
	"encoding/binary"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"bitbucket.org/decimalteam/go-node/x/coin"
)

const (
	// ModuleName is the name of the module
	ModuleName = "gov"

	// StoreKey is the store key string for gov
	StoreKey = coin.StoreKey

	// RouterKey is the message route for gov
	RouterKey = ModuleName

	// QuerierRoute is the querier route for gov
	QuerierRoute = ModuleName

	// DefaultParamspace default name for parameter store
	DefaultParamspace = ModuleName

	KeyUpgradedIBCState = "upgradedIBCState"

	KeyUpgradedClient = "upgradedClient"

	KeyUpgradedConsState = "upgradedConsState"
)

var (
	LegacyProposalsKeyPrefix          = []byte{0x00}
	LegacyActiveProposalQueuePrefix   = []byte{0x01}
	LegacyInactiveProposalQueuePrefix = []byte{0x02}
	LegacyProposalIDKey               = []byte{0x03}

	LegacyVotesKeyPrefix = []byte{0x10}

	LegacyPlanPrefix = []byte{0x20}
	LegacyDonePrefix = []byte{0x21}

	// This is special key used to determine if kv-records are migrated to keys with correct prefixes
	LegacyMigrationKey = []byte("gov/migrated")
)

// Keys for governance store
// Items are stored with the following key: values
//
// - gov/proposals/<proposalID_Bytes>: Proposal
//
// - gov/proposals/active/<endTime_Bytes><proposalID_Bytes>: activeProposalID
//
// - gov/proposals/inactive/<endTime_Bytes><proposalID_Bytes>: inactiveProposalID
//
// - gov/proposals/next: nextProposalID
//
// - gov/votes/<proposalID_Bytes><voterAddr_Bytes>: Voter
var (
	ProposalsKeyPrefix          = []byte("gov/proposals/")
	ActiveProposalQueuePrefix   = []byte("gov/proposals/active/")
	InactiveProposalQueuePrefix = []byte("gov/proposals/inactive/")
	ProposalIDKey               = []byte("gov/proposals/next")

	VotesKeyPrefix = []byte("gov/votes/")

	PlanPrefix = []byte("gov/plan")
	DonePrefix = []byte("gov/done")
)

// GetProposalIDBytes returns the byte representation of the proposalID
func GetProposalIDBytes(proposalID uint64) (proposalIDBz []byte) {
	proposalIDBz = make([]byte, 8)
	binary.BigEndian.PutUint64(proposalIDBz, proposalID)
	return
}

func GetBytesFromUint64(i uint64) []byte {
	bytes := make([]byte, 8)
	binary.BigEndian.PutUint64(bytes, i)
	return bytes
}

func GetUint64FromBytes(b []byte) uint64 {
	return binary.BigEndian.Uint64(b)
}

// GetProposalIDFromBytes returns proposalID in uint64 format from a byte array
func GetProposalIDFromBytes(bz []byte) (proposalID uint64) {
	return binary.BigEndian.Uint64(bz)
}

// ProposalKey gets a specific proposal from the store
func ProposalKey(proposalID uint64) []byte {
	return append(ProposalsKeyPrefix, GetProposalIDBytes(proposalID)...)
}

// ActiveProposalByTimeKey gets the active proposal queue key by endTime
func ActiveProposalByTimeKey(endBlock uint64) []byte {
	return append(ActiveProposalQueuePrefix, GetBytesFromUint64(endBlock)...)
}

// ActiveProposalQueueKey returns the key for a proposalID in the activeProposalQueue
func ActiveProposalQueueKey(proposalID uint64, endBlock uint64) []byte {
	return append(ActiveProposalByTimeKey(endBlock), GetProposalIDBytes(proposalID)...)
}

// InactiveProposalByTimeKey gets the inactive proposal queue key by endTime
func InactiveProposalByTimeKey(endBlock uint64) []byte {
	return append(InactiveProposalQueuePrefix, GetBytesFromUint64(endBlock)...)
}

// InactiveProposalQueueKey returns the key for a proposalID in the inactiveProposalQueue
func InactiveProposalQueueKey(proposalID uint64, endBlock uint64) []byte {
	return append(InactiveProposalByTimeKey(endBlock), GetProposalIDBytes(proposalID)...)
}

// VotesKey gets the first part of the votes key based on the proposalID
func VotesKey(proposalID uint64) []byte {
	return append(VotesKeyPrefix, GetProposalIDBytes(proposalID)...)
}

// VoteKey key of a specific vote from the store
func VoteKey(proposalID uint64, voterAddr sdk.ValAddress) []byte {
	return append(VotesKey(proposalID), voterAddr.Bytes()...)
}

// Split keys function; used for iterators

// SplitProposalKey split the proposal key and returns the proposal id
func SplitProposalKey(key []byte) (proposalID uint64) {
	tail := key[len(ProposalsKeyPrefix):]
	if len(tail) != 8 {
		panic(fmt.Sprintf("unexpected key length (%d)", len(key)))
	}
	return GetProposalIDFromBytes(tail)
}

// SplitActiveProposalQueueKey split the active proposal key and returns the proposal id and endBlock
func SplitActiveProposalQueueKey(key []byte) (proposalID uint64, endBlock uint64) {
	tail := key[len(ActiveProposalQueuePrefix):]
	if len(tail) != 16 {
		panic(fmt.Sprintf("unexpected key length (%d)", len(key)))
	}
	endBlock = GetUint64FromBytes(tail[:8])
	proposalID = GetProposalIDFromBytes(tail[8:])
	return
}

// SplitInactiveProposalQueueKey split the inactive proposal key and returns the proposal id and endBlock
func SplitInactiveProposalQueueKey(key []byte) (proposalID uint64, endBlock uint64) {
	tail := key[len(InactiveProposalQueuePrefix):]
	if len(tail) != 16 {
		panic(fmt.Sprintf("unexpected key length (%d)", len(key)))
	}
	endBlock = GetUint64FromBytes(tail[:8])
	proposalID = GetProposalIDFromBytes(tail[8:])
	return
}

// PlanKey is the key under which the current plan is saved
// We store PlanByte as a const to keep it immutable (unlike a []byte)
func PlanKey() []byte {
	return PlanPrefix
}

// DoneKey is the key at which the given upgrade was executed
// We store DoneKey as a const to keep it immutable (unlike a []byte)
func DoneKey() []byte {
	return DonePrefix
}

func UpgradedClientKey(height int64) []byte {
	return []byte(fmt.Sprintf("%s/%d/%s", KeyUpgradedIBCState, height, KeyUpgradedClient))
}

// UpgradedConsStateKey is the key under which the upgraded consensus state is saved
// Connecting IBC chains can verify against the upgraded consensus state in this path before
// upgrading their clients.
func UpgradedConsStateKey(height int64) []byte {
	return []byte(fmt.Sprintf("%s/%d/%s", KeyUpgradedIBCState, height, KeyUpgradedConsState))
}
