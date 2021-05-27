package types

// validator module event types
const (
	EventTypeDeclareCandidate     = "declare_candidate"
	EventTypeEditCandidate        = "edit_candidate"
	EventTypeDelegate             = "delegate"
	EventTypeUnbond               = "unbond"
	EventTypeSetOnline            = "set_online"
	EventTypeSetOffline           = "set_offline"
	EventTypeCompleteUnbonding    = "complete_unbonding"
	EventTypeCompleteUnbondingNFT = "complete_unbonding_nft"
	EventTypeProposerReward       = "proposer_reward"
	EventTypeCommissionReward     = "commission_reward"
	EventTypeSlash                = "slash"
	EventTypeEmission             = "emission"
	EventTypeLiveness             = "liveness"
	EventTypeLivenessNFT          = "liveness_nft"
	EventTypeUpdatesValidators    = "updates_validator"
	EventTypeCalcStake            = "calc_stake"
	EventTypeDAOReward            = "dao_reward"
	EventTypeDevelopReward        = "develop_reward"

	AttributeKeyValidator                  = "validator"
	AttributeKeyDelegator                  = "delegator"
	AttributeKeyRewardAddress              = "reward_address"
	AttributeKeyCoin                       = "coin"
	AttributeKeyPubKey                     = "pub_key"
	AttributeKeyCompletionTime             = "completion_time"
	AttributeKeyAddress                    = "address"
	AttributeKeyReason                     = "reason"
	AttributeKeyPower                      = "power"
	AttributeKeyStake                      = "stake"
	AttributeKeyValidatorOdCandidate       = "status"
	AttributeKeySlashAmount                = "slash_amount"
	AttributeKeySlashSubTokenID            = "sub_token_id"
	AttributeKeySlashReserve               = "sub_token_id_reserve"
	AttributeKeyMissedBlocks               = "missed_blocks"
	AttributeKeyHeight                     = "height"
	AttributeKeyCommission                 = "commission"
	AttributeKeyDescriptionMoniker         = "moniker"
	AttributeKeyDescriptionIdentity        = "identity"
	AttributeKeyDescriptionWebsite         = "website"
	AttributeKeyDescriptionSecurityContact = "security_contact"
	AttributeKeyDescriptionDetails         = "details"
	AttributeKeyDAOAddress                 = "dao_address"
	AttributeKeyDevelopAddress             = "develop_address"
	AttributeKeyDenom                      = "denom"
	AttributeKeyID                         = "id"
	AttributeKeyQuantity                   = "quantity"
	AttributeKeySubTokenIDs                = "sub_token_ids"

	AttributeValueDoubleSign       = "double_sign"
	AttributeValueMissingSignature = "missing_signature"

	AttributeValueCategory = ModuleName
)
