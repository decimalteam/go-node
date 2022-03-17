package types

import (
	"encoding/binary"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"bitbucket.org/decimalteam/go-node/utils/updates"
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
	ProposalsKeyPrefix          = []byte("gov/proposals/")          // []byte{0x00}
	ActiveProposalQueuePrefix   = []byte("gov/proposals/active/")   // []byte{0x01}
	InactiveProposalQueuePrefix = []byte("gov/proposals/inactive/") // []byte{0x02}
	ProposalIDKey               = []byte("gov/proposals/next")      // []byte{0x03}

	VotesKeyPrefix = []byte("gov/votes/") // []byte{0x10}

	PlanPrefix = []byte("gov/plan") // []byte{0x20}
	DonePrefix = []byte("gov/done") // []byte{0x21}
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
func ProposalKey(ctx sdk.Context, proposalID uint64) []byte {
	keyPrefix := ProposalsKeyPrefix
	if ctx.BlockHeight() < updates.Update14Block {
		keyPrefix = []byte{0x00}
	}
	return append(keyPrefix, GetProposalIDBytes(proposalID)...)
}

// ActiveProposalByTimeKey gets the active proposal queue key by endTime
func ActiveProposalByTimeKey(ctx sdk.Context, endBlock uint64) []byte {
	keyPrefix := ActiveProposalQueuePrefix
	if ctx.BlockHeight() < updates.Update14Block {
		keyPrefix = []byte{0x01}
	}
	return append(keyPrefix, GetBytesFromUint64(endBlock)...)
}

// ActiveProposalQueueKey returns the key for a proposalID in the activeProposalQueue
func ActiveProposalQueueKey(ctx sdk.Context, proposalID uint64, endBlock uint64) []byte {
	return append(ActiveProposalByTimeKey(ctx, endBlock), GetProposalIDBytes(proposalID)...)
}

// InactiveProposalByTimeKey gets the inactive proposal queue key by endTime
func InactiveProposalByTimeKey(ctx sdk.Context, endBlock uint64) []byte {
	keyPrefix := InactiveProposalQueuePrefix
	if ctx.BlockHeight() < updates.Update14Block {
		keyPrefix = []byte{0x02}
	}
	return append(keyPrefix, GetBytesFromUint64(endBlock)...)
}

// InactiveProposalQueueKey returns the key for a proposalID in the inactiveProposalQueue
func InactiveProposalQueueKey(ctx sdk.Context, proposalID uint64, endBlock uint64) []byte {
	return append(InactiveProposalByTimeKey(ctx, endBlock), GetProposalIDBytes(proposalID)...)
}

// VotesKey gets the first part of the votes key based on the proposalID
func VotesKey(ctx sdk.Context, proposalID uint64) []byte {
	keyPrefix := VotesKeyPrefix
	if ctx.BlockHeight() < updates.Update14Block {
		keyPrefix = []byte{0x10}
	}
	return append(keyPrefix, GetProposalIDBytes(proposalID)...)
}

// VoteKey key of a specific vote from the store
func VoteKey(ctx sdk.Context, proposalID uint64, voterAddr sdk.ValAddress) []byte {
	return append(VotesKey(ctx, proposalID), voterAddr.Bytes()...)
}

// Split keys function; used for iterators

// SplitProposalKey split the proposal key and returns the proposal id
func SplitProposalKey(ctx sdk.Context, key []byte) (proposalID uint64) {
	keyPrefix := ProposalsKeyPrefix
	if ctx.BlockHeight() < updates.Update14Block {
		keyPrefix = []byte{0x00}
	}
	tail := key[len(keyPrefix):]
	if len(tail) != 8 {
		panic(fmt.Sprintf("unexpected key length (%d)", len(key)))
	}
	return GetProposalIDFromBytes(tail)
}

// SplitActiveProposalQueueKey split the active proposal key and returns the proposal id and endBlock
func SplitActiveProposalQueueKey(ctx sdk.Context, key []byte) (proposalID uint64, endBlock uint64) {
	keyPrefix := ActiveProposalQueuePrefix
	if ctx.BlockHeight() < updates.Update14Block {
		keyPrefix = []byte{0x01}
	}
	tail := key[len(keyPrefix):]
	if len(tail) != 16 {
		panic(fmt.Sprintf("unexpected key length (%d)", len(key)))
	}
	endBlock = GetUint64FromBytes(tail[:8])
	proposalID = GetProposalIDFromBytes(tail[8:])
	return
}

// SplitInactiveProposalQueueKey split the inactive proposal key and returns the proposal id and endBlock
func SplitInactiveProposalQueueKey(ctx sdk.Context, key []byte) (proposalID uint64, endBlock uint64) {
	keyPrefix := InactiveProposalQueuePrefix
	if ctx.BlockHeight() < updates.Update14Block {
		keyPrefix = []byte{0x02}
	}
	tail := key[len(keyPrefix):]
	if len(tail) != 16 {
		panic(fmt.Sprintf("unexpected key length (%d)", len(key)))
	}
	endBlock = GetUint64FromBytes(tail[:8])
	proposalID = GetProposalIDFromBytes(tail[8:])
	return
}

// SplitKeyDeposit split the deposits key and returns the proposal id and depositor address
func SplitKeyDeposit(ctx sdk.Context, key []byte) (proposalID uint64, depositorAddr sdk.AccAddress) {
	return splitKeyWithAddress(ctx, key)
}

// SplitKeyVote split the votes key and returns the proposal id and voter address
func SplitKeyVote(ctx sdk.Context, key []byte) (proposalID uint64, voterAddr sdk.AccAddress) {
	return splitKeyWithAddress(ctx, key)
}

// private functions

func splitKeyWithAddress(ctx sdk.Context, key []byte) (proposalID uint64, addr sdk.AccAddress) {
	keyPrefix := VotesKeyPrefix
	if ctx.BlockHeight() < updates.Update14Block {
		keyPrefix = []byte{0x10}
	}
	tail := key[len(keyPrefix):]
	if len(tail) != 8+sdk.AddrLen {
		panic(fmt.Sprintf("unexpected key length (%d)", len(key)))
	}
	proposalID = GetProposalIDFromBytes(key[:8])
	addr = key[8:]
	return
}

// PlanKey is the key under which the current plan is saved
// We store PlanByte as a const to keep it immutable (unlike a []byte)
func PlanKey(ctx sdk.Context) []byte {
	keyPrefix := PlanPrefix
	if ctx.BlockHeight() < updates.Update14Block {
		keyPrefix = []byte{0x20}
	}
	return keyPrefix
}

// DoneKey is the key at which the given upgrade was executed
// We store DoneKey as a const to keep it immutable (unlike a []byte)
func DoneKey(ctx sdk.Context) []byte {
	keyPrefix := DonePrefix
	if ctx.BlockHeight() < updates.Update14Block {
		keyPrefix = []byte{0x21}
	}
	return keyPrefix
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
