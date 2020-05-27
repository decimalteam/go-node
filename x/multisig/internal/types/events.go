package types

// multisig module event types
const (
	EventTypeCreateWallet      = "CreateWallet"
	EventTypeCreateTransaction = "CreateTransaction"
	EventTypeSignTransaction   = "SignTransaction"

	// Common
	AttributeKeyCreator     = "creator"
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
	AttributeKeySigner            = "signer"
	AttributeKeySignerWeight      = "signerWeight"
	AttributeKeyTransactionWeight = "transactionWeight"
	AttributeKeyApproved          = "approved"

	AttributeValueCategory = ModuleName
)
