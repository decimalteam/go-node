package types

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	decsdk "bitbucket.org/decimalteam/go-node/utils/types"
)

const (
	SignedBlocksWindow int64 = 24
	MinSignedPerWindow int64 = 12
)

var (
	SlashFractionDowntime   = sdk.NewDec(1).Quo(sdk.NewDec(100))
	SlashFractionDoubleSign = sdk.NewDec(1).Quo(sdk.NewDec(20))
)

// Signing info for a validator
type ValidatorSigningInfo struct {
	Address             decsdk.ConsAddress `json:"address" yaml:"address"`                             // validator consensus address
	StartHeight         int64              `json:"start_height" yaml:"start_height"`                   // height at which validator was first a candidate OR was unjailed
	IndexOffset         int64              `json:"index_offset" yaml:"index_offset"`                   // index offset into signed block bit array
	JailedUntil         time.Time          `json:"jailed_until" yaml:"jailed_until"`                   // timestamp validator cannot be unjailed until
	Tombstoned          bool               `json:"tombstoned" yaml:"tombstoned"`                       // whether or not a validator has been tombstoned (killed out of validator set)
	MissedBlocksCounter int64              `json:"missed_blocks_counter" yaml:"missed_blocks_counter"` // missed blocks counter (to avoid scanning the array every time)
}

// Construct a new `ValidatorSigningInfo` struct
func NewValidatorSigningInfo(
	condAddr decsdk.ConsAddress, startHeight, indexOffset int64,
	jailedUntil time.Time, tombstoned bool, missedBlocksCounter int64,
) ValidatorSigningInfo {

	return ValidatorSigningInfo{
		Address:             condAddr,
		StartHeight:         startHeight,
		IndexOffset:         indexOffset,
		JailedUntil:         jailedUntil,
		Tombstoned:          tombstoned,
		MissedBlocksCounter: missedBlocksCounter,
	}
}

// Return human readable signing info
func (i ValidatorSigningInfo) String() string {
	return fmt.Sprintf(`Validator Signing Info:
  Address:               %s
  Start Height:          %d
  Index Offset:          %d
  Jailed Until:          %v
  Tombstoned:            %t
  Missed Blocks Counter: %d`,
		i.Address, i.StartHeight, i.IndexOffset, i.JailedUntil,
		i.Tombstoned, i.MissedBlocksCounter)
}

// get pubkey relation key used to get the pubkey from the address
func GetAddrPubkeyRelationKey(address []byte) []byte {
	return append([]byte{AddrPubkeyRelationKey}, address...)
}
