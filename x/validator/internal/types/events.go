package types

// validator module event types
const (
	EventTypeDeclareCandidate  = "create_validator"
	EventTypeDelegate          = "delegate"
	EventTypeUnbond            = "unbond"
	EventTypeCompleteUnbonding = "complete_unbonding"

	AttributeKeyValidator      = "validator"
	AttributeKeyDelegator      = "delegator"
	AttributeKeyCompletionTime = "completion_time"
	AttributeValueCategory     = ModuleName
)
