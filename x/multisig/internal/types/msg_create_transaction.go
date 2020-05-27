package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ sdk.Msg = &MsgCreateTransaction{}

// MsgCreateTransaction defines a CreateTransaction message to create new transaction for multisignature wallet.
type MsgCreateTransaction struct {
	Creator  sdk.AccAddress `json:"creator" yaml:"creator"`
	Wallet   sdk.AccAddress `json:"wallet" yaml:"wallet"`
	Receiver sdk.AccAddress `json:"receiver" yaml:"receiver"`
	Coins    sdk.Coins      `json:"coins" yaml:"coins"`
}

// NewMsgCreateTransaction creates a new MsgCreateTransaction instance.
func NewMsgCreateTransaction(creator sdk.AccAddress, wallet sdk.AccAddress, receiver sdk.AccAddress, coins sdk.Coins) MsgCreateTransaction {
	return MsgCreateTransaction{
		Creator:  creator,
		Wallet:   wallet,
		Receiver: receiver,
		Coins:    coins,
	}
}

// Route returns name of the route for the message.
func (msg MsgCreateTransaction) Route() string { return RouterKey }

// Type returns the name of the type for the message.
func (msg MsgCreateTransaction) Type() string { return "CreateTransaction" }

// ValidateBasic performs basic validation of the message.
func (msg MsgCreateTransaction) ValidateBasic() error {
	// TODO
	return nil
}

// GetSignBytes returns the canonical byte representation of the message used to generate a signature.
func (msg MsgCreateTransaction) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners returns the list of signers required to sign the message.
func (msg MsgCreateTransaction) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Creator}
}
