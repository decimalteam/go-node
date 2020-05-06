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
	EventTypeSlash             = "slash"

	AttributeKeyValidator      = "validator"
	AttributeKeyDelegator      = "delegator"
	AttributeKeyRewardAddress  = "reward_address"
	AttributeKeyDenom          = "denom"
	AttributeKeyPubKey         = "pub_key"
	AttributeKeyCompletionTime = "completion_time"
	AttributeKeyAddress        = "address"
	AttributeKeyReason         = "reason"
	AttributeKeyPower          = "power"

	AttributeValueDoubleSign       = "double_sign"
	AttributeValueMissingSignature = "missing_signature"

	AttributeValueCategory = ModuleName
)
