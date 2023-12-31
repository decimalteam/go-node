package types

// coin module event types
const (
	EventTypeCreateCoin    = "create_coin"
	EventTypeBuyCoin       = "buy_coin"
	EventTypeSellCoin      = "sell_coin"
	EventTypeSellAllCoin   = "sell_all_coin"
	EventTypeSendCoin      = "send_coin"
	EventTypeMultiSendCoin = "multi_send_coin"
	EventTypeRedeemCheck   = "redeem_check"
	EventTypeUpdateCoin    = "update_coin"
	// Create Coin
	AttributeTitle                 = "title"
	AttributeSymbol                = "symbol"
	AttributeCRR                   = "crr"
	AttributeVolume                = "volume"
	AttributeReserve               = "reserve"
	AttributeInitVolume            = "initial_volume"
	AttributeInitReserve           = "initial_reserve"
	AttributeLimitVolume           = "limit_volume"
	AttributeCommissionCreateCoin  = "commission_create_coin"
	AttributeCommissionRedeemCheck = "commission_redeem_check"

	// Buy/Sell Coin
	AttributeCoinToBuy        = "coin_to_buy"
	AttributeCoinToSell       = "coin_to_sell"
	AttributeAmountToBuy      = "amount_to_buy"
	AttributeAmountToSell     = "amount_to_sell"
	AttributeAmountInBaseCoin = "amount_in_base_coin"

	// Send/MultiSend Coin
	AttributeCoin     = "coin"
	AttributeAmount   = "amount"
	AttributeReceiver = "receiver"

	// Redeem Check
	AttributeIssuer   = "issuer"
	AttributeNonce    = "nonce"
	AttributeDueBlock = "due_block"

	AttributeValueCategory = ModuleName
)
