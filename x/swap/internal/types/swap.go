package types

import (
	"encoding/hex"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Swap struct {
	TransferType TransferType   `json:"transfer_type"`
	HashedSecret Hash           `json:"hashed_secret"`
	From         sdk.AccAddress `json:"from"`
	Recipient    string         `json:"recipient"`
	Amount       sdk.Coins      `json:"amount"`
	Timestamp    uint64         `json:"timestamp"`
	Redeemed     bool           `json:"redeemed"`
	Refunded     bool           `json:"refunded"`
}

func NewSwap(transferType TransferType, hash Hash, from sdk.AccAddress, recipient string, amount sdk.Coins, timestamp uint64) Swap {
	return Swap{TransferType: transferType, HashedSecret: hash, From: from, Recipient: recipient, Amount: amount, Timestamp: timestamp, Redeemed: false, Refunded: false}
}

type Swaps []Swap

type Hash [32]byte

func (h Hash) MarshalJSON() ([]byte, error) {
	return []byte("\"" + hex.EncodeToString(h[:]) + "\""), nil
}

func (h *Hash) UnmarshalJSON(b []byte) error {
	decoded, err := hex.DecodeString(string(b[1 : len(b)-1]))
	if err != nil {
		return err
	}
	copy(h[:], decoded)
	return nil
}

type Secret []byte

func (s Secret) MarshalJSON() ([]byte, error) {
	return []byte("\"" + hex.EncodeToString(s[:]) + "\""), nil
}

func (s *Secret) UnmarshalJSON(b []byte) error {
	decoded, err := hex.DecodeString(string(b[1 : len(b)-1]))
	if err != nil {
		return err
	}
	copy(b[:], decoded)
	return nil
}
