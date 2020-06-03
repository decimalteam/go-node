package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
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

// Route returns name of the route for the message.
func (msg MsgCreateWallet) Route() string { return RouterKey }

// Type returns the name of the type for the message.
func (msg MsgCreateWallet) Type() string { return "create_wallet" }

// ValidateBasic performs basic validation of the message.
func (msg MsgCreateWallet) ValidateBasic() error {
	// Validate sender
	if msg.Sender.Empty() {
		return sdkerrors.New(
			DefaultCodespace,
			InvalidSender,
			"Invalid sender address: sender address cannot be empty",
		)
	}
	// Validate owner count
	if len(msg.Owners) < MinOwnerCount {
		return sdkerrors.New(
			DefaultCodespace,
			InvalidOwnerCount,
			fmt.Sprintf("Invalid owner count: need at least %d owners", MinOwnerCount),
		)
	}
	if len(msg.Owners) > MaxOwnerCount {
		return sdkerrors.New(
			DefaultCodespace,
			InvalidOwnerCount,
			fmt.Sprintf("Invalid owner count: allowed no more than %d owners", MaxOwnerCount),
		)
	}
	// Validate weight count
	if len(msg.Owners) != len(msg.Weights) {
		return sdkerrors.New(
			DefaultCodespace,
			InvalidWeightCount,
			fmt.Sprintf("Invalid weight count: weight count (%d) is not equal to owner count (%d)", len(msg.Weights), len(msg.Owners)),
		)
	}
	// Validate owners (ensure there are no duplicates)
	owners := make(map[string]bool, len(msg.Owners))
	for i, c := 0, len(msg.Owners); i < c; i++ {
		if msg.Owners[i].Empty() {
			return sdkerrors.New(
				DefaultCodespace,
				InvalidOwner,
				"Invalid owner address: owner address cannot be empty",
			)
		}
		if owners[msg.Owners[i].String()] {
			return sdkerrors.New(
				DefaultCodespace,
				InvalidOwner,
				fmt.Sprintf("Invalid owners: owner with address %s is duplicated", msg.Owners[i]),
			)
		}
		owners[msg.Owners[i].String()] = true
	}
	// Validate weights
	for i, c := 0, len(msg.Weights); i < c; i++ {
		if msg.Weights[i] < MinWeight {
			return sdkerrors.New(
				DefaultCodespace,
				InvalidWeight,
				fmt.Sprintf("Invalid weight: weight cannot be less than %d", MinWeight),
			)
		}
		if msg.Weights[i] > MaxWeight {
			return sdkerrors.New(
				DefaultCodespace,
				InvalidWeight,
				fmt.Sprintf("Invalid weight: weight cannot be greater than %d", MaxWeight),
			)
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
