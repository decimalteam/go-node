package types

// NFT module event types
var (
	EventTypeTransfer        = "transfer_nft"
	EventTypeEditNFTMetadata = "edit_nft_metadata"
	EventTypeMintNFT         = "mint_nft"
	EventTypeBurnNFT         = "burn_nft"

	AttributeValueCategory = ModuleName

	AttributeKeySender               = "sender"
	AttributeKeyRecipient            = "recipient"
	AttributeKeyOwner                = "owner"
	AttributeKeyNFTID                = "nft_id"
	AttributeKeyNFTTokenURI          = "token_uri"
	AttributeKeyDenom                = "denom"
	AttributeKeySubTokenIDStartRange = "sub_token_id_start_range"
)
