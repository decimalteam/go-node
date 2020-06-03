package types

// validator module event types
const (
	EventTypeDeclareCandidate  = "declare_candidate"
	EventTypeEditCandidate     = "edit_candidate"
	EventTypeDelegate          = "delegate"
	EventTypeUnbond            = "unbond"
	EventTypeSetOnline         = "set_online"
	EventTypeSetOffline        = "set_offline"
	EventTypeCompleteUnbonding = "complete_unbonding"
	EventTypeProposerReward    = "proposer_reward"
	EventTypeCommissionReward  = "commission_reward"
	EventTypeSlash             = "slash"
	EventTypeEmission          = "emission"
	EventTypeLiveness          = "liveness"
	EventTypeUpdatesValidators = "updates_validator"

	AttributeKeyValidator                  = "validator"
	AttributeKeyDelegator                  = "delegator"
	AttributeKeyRewardAddress              = "reward_address"
	AttributeKeyCoin                       = "coin"
	AttributeKeyPubKey                     = "pub_key"
	AttributeKeyCompletionTime             = "completion_time"
	AttributeKeyAddress                    = "address"
	AttributeKeyReason                     = "reason"
	AttributeKeyPower                      = "power"
	AttributeKeySequence                   = "sequence"
	AttributeKeySlashAmount                = "slash_amount"
	AttributeKeyMissedBlocks               = "missed_blocks"
	AttributeKeyHeight                     = "height"
	AttributeKeyCommission                 = "commission"
	AttributeKeyDescriptionMoniker         = "moniker"
	AttributeKeyDescriptionIdentity        = "identity"
	AttributeKeyDescriptionWebsite         = "website"
	AttributeKeyDescriptionSecurityContact = "security_contact"
	AttributeKeyDescriptionDetails         = "details"

	AttributeValueDoubleSign       = "double_sign"
	AttributeValueMissingSignature = "missing_signature"

	AttributeValueCategory = ModuleName
)
