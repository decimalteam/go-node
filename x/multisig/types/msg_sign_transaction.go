package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ sdk.Msg = &MsgSignTransaction{}

// MsgSignTransaction defines a SignTransaction message to sign existing transaction for multisignature wallet.
//type MsgSignTransaction struct {
//	Sender sdk.AccAddress `json:"sender" yaml:"sender"`
//	TxID   string         `json:"tx_id" yaml:"tx_id"`
//}

// NewMsgSignTransaction is a constructor function for MsgCreateTransaction
func NewMsgSignTransaction(sender sdk.AccAddress, txID string) MsgSignTransaction {
	return MsgSignTransaction{
		Sender: sender,
		TxID:   txID,
	}
}

const SignTransactionConst = "sign_transaction"

// Route returns name of the route for the message.
func (msg *MsgSignTransaction) Route() string { return RouterKey }

// Type returns the name of the type for the message.
func (msg *MsgSignTransaction) Type() string { return SignTransactionConst }

// ValidateBasic runs stateless checks on the message
func (msg *MsgSignTransaction) ValidateBasic() error {
	// TODO
	return nil
}

// GetSignBytes returns the canonical byte representation of the message used to generate a signature.
func (msg *MsgSignTransaction) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners returns the list of signers required to sign the message.
func (msg *MsgSignTransaction) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}
