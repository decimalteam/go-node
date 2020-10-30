package types

import (
	"encoding/hex"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Swap struct {
	TransferType TransferType   `json:"transfer_type"s`
	Hash         Hash           `json:"hash"`
	From         sdk.AccAddress `json:"from"`
	Recipient    string         `json:"recipient"`
	Amount       sdk.Coins      `json:"amount"`
	Timestamp    uint64         `json:"timestamp"`
	Claimed      bool           `json:"claimed"`
	Refunded     bool           `json:"refunded"`
}

func NewSwap(transferType TransferType, hash Hash, from sdk.AccAddress, recipient string, amount sdk.Coins, timestamp uint64) Swap {
	return Swap{TransferType: transferType, Hash: hash, From: from, Recipient: recipient, Amount: amount, Timestamp: timestamp, Claimed: false, Refunded: false}
}

type Swaps []Swap

type Hash [32]byte

func (h Hash) MarshalJSON() ([]byte, error) {
	return []byte("\"" + hex.EncodeToString(h[:]) + "\""), nil
}

func (h *Hash) UnmarshalJSON(b []byte) error {
	decoded, err := hex.DecodeString(string(b))
	if err != nil {
		return err
	}
	copy(h[:], decoded)
	return nil
}
