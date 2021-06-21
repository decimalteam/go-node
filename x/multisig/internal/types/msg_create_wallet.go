package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"strconv"
)

var _ sdk.Msg = &MsgCreateWallet{}

// Multisignature wallet limitations.
const (
	MinOwnerCount = 2
	MaxOwnerCount = 16
	MinWeight     = 1
	MaxWeight     = 1024
)

// MsgCreateWallet defines a CreateWallet message to create new multisignature wallet.
type MsgCreateWallet struct {
	Sender    sdk.AccAddress   `json:"sender" yaml:"sender"`
	Owners    []sdk.AccAddress `json:"owners" yaml:"owners"`
	Weights   []uint           `json:"weights" yaml:"weights"`
	Threshold uint             `json:"threshold" yaml:"threshold"`
}

// NewMsgCreateWallet creates a new MsgCreateWallet instance.
func NewMsgCreateWallet(sender sdk.AccAddress, owners []sdk.AccAddress, weights []uint, threshold uint) MsgCreateWallet {
	return MsgCreateWallet{
		Sender:    sender,
		Owners:    owners,
		Weights:   weights,
		Threshold: threshold,
	}
}

const CreateWalletConst = "create_wallet"

// Route returns name of the route for the message.
func (msg MsgCreateWallet) Route() string { return RouterKey }

// Type returns the name of the type for the message.
func (msg MsgCreateWallet) Type() string { return CreateWalletConst }

// ValidateBasic performs basic validation of the message.
func (msg MsgCreateWallet) ValidateBasic() error {
	// Validate sender
	if msg.Sender.Empty() {
		return ErrInvalidSender()
	}
	// Validate owner count
	if len(msg.Owners) < MinOwnerCount {
		return ErrInvalidOwnerCount(false)
	}
	if len(msg.Owners) > MaxOwnerCount {
		return ErrInvalidOwnerCount(true)
	}
	// Validate weight count
	if len(msg.Owners) != len(msg.Weights) {
		return ErrInvalidWeightCount(strconv.Itoa(len(msg.Weights)), strconv.Itoa(len(msg.Owners)))
	}
	// Validate owners (ensure there are no duplicates)
	owners := make(map[string]bool, len(msg.Owners))
	for i, c := 0, len(msg.Owners); i < c; i++ {
		if msg.Owners[i].Empty() {
			return ErrInvalidOwner()
		}
		if owners[msg.Owners[i].String()] {
			return ErrDuplicateOwner(msg.Owners[i].String())
		}
		owners[msg.Owners[i].String()] = true
	}
	// Validate weights
	for i, c := 0, len(msg.Weights); i < c; i++ {
		if msg.Weights[i] < MinWeight {
			return ErrInvalidWeight(strconv.Itoa(MinWeight), "less")
		}
		if msg.Weights[i] > MaxWeight {
			return ErrInvalidWeight(strconv.Itoa(MaxWeight), "greater")
		}
	}
	return nil
}

// GetSignBytes returns the canonical byte representation of the message used to generate a signature.
func (msg MsgCreateWallet) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners returns the list of signers required to sign the message.
func (msg MsgCreateWallet) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}
