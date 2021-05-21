package types

// multisig module event types
const (
	EventTypeCreateWallet      = "create_wallet"
	EventTypeCreateTransaction = "create_transaction"
	EventTypeSignTransaction   = "sign_transaction"

	// Common
	AttributeKeySender      = "sender"
	AttributeKeyWallet      = "wallet"
	AttributeKeyTransaction = "transaction"

	// CreateWallet
	AttributeKeyOwners    = "owners"
	AttributeKeyWeights   = "weights"
	AttributeKeyThreshold = "threshold"

	// CreateTransaction
	AttributeKeyReceiver = "receiver"
	AttributeKeyCoins    = "coins"

	// SignTransaction
	AttributeKeySigner        = "signer"
	AttributeKeySignerWeight  = "signer_weight"
	AttributeKeyConfirmations = "confirmations"
	AttributeKeyConfirmed     = "confirmed"

	AttributeValueCategory = ModuleName
)
