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

	AttributeKeyValidator      = "validator"
	AttributeKeyDelegator      = "delegator"
	AttributeKeyRewardAddress  = "reward_address"
	AttributeKeyDenom          = "denom"
	AttributeKeyCompletionTime = "completion_time"
	AttributeValueCategory     = ModuleName
)
