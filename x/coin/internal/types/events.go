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
	// Create Coin
	AttributeTitle       = "title"
	AttributeSymbol      = "symbol"
	AttributeCRR         = "crr"
	AttributeInitVolume  = "initial_volume"
	AttributeInitReserve = "initial_reserve"
	AttributeLimitVolume = "limit_volume"

	// Buy/Sell Coin
	AttributeCoinToBuy        = "coinToBuy"
	AttributeCoinToSell       = "coinToSell"
	AttributeAmountToBuy      = "amountToBuy"
	AttributeAmountToSell     = "amountToSell"
	AttributeAmountInBaseCoin = "amountInBaseCoin"

	// Send/MultiSend Coin
	AttributeCoin     = "coin"
	AttributeAmount   = "amount"
	AttributeReceiver = "receiver"

	// Redeem Check
	AttributeIssuer     = "issuer"
	AttributeCheckNonce = "checkNonce"
	AttributeDueBlock   = "dueBlock"

	AttributeValueCategory = ModuleName
)
