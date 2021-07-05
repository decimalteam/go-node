package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ sdk.Msg = &MsgCreateTransaction{}

// MsgCreateTransaction defines a CreateTransaction message to create new transaction for multisignature wallet.
//type MsgCreateTransaction struct {
//	Sender   sdk.AccAddress `json:"sender" yaml:"sender"`
//	Wallet   sdk.AccAddress `json:"wallet" yaml:"wallet"`
//	Receiver sdk.AccAddress `json:"receiver" yaml:"receiver"`
//	Coins    sdk.Coins      `json:"coins" yaml:"coins"`
//}

// NewMsgCreateTransaction creates a new MsgCreateTransaction instance.
func NewMsgCreateTransaction(sender sdk.AccAddress, wallet sdk.AccAddress, receiver sdk.AccAddress, coins sdk.Coins) MsgCreateTransaction {
	return MsgCreateTransaction{
		Sender:   sender.String(),
		Wallet:   wallet.String(),
		Receiver: receiver.String(),
		Coins:    coins,
	}
}

const CreateTransactionConst = "create_transaction"

// Route returns name of the route for the message.
func (msg *MsgCreateTransaction) Route() string { return RouterKey }

// Type returns the name of the type for the message.
func (msg *MsgCreateTransaction) Type() string { return CreateTransactionConst }

// ValidateBasic performs basic validation of the message.
func (msg *MsgCreateTransaction) ValidateBasic() error {
	// TODO
	return nil
}

// GetSignBytes returns the canonical byte representation of the message used to generate a signature.
func (msg *MsgCreateTransaction) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners returns the list of signers required to sign the message.
func (msg *MsgCreateTransaction) GetSigners() []sdk.AccAddress {
	accAddr, _ := sdk.AccAddressFromBech32(msg.Sender)
	
	return []sdk.AccAddress{accAddr}
}
