package types

// coin module event types
const (
	EventTypeCreateCoin    = "CreateCoin"
	EventTypeBuyCoin       = "BuyCoin"
	EventTypeSellCoin      = "SellCoin"
	EventTypeSellAllCoin   = "SellAllCoin"
	EventTypeSendCoin      = "SendCoin"
	EventTypeMultiSendCoin = "MultiSendCoin"
	EventTypeRedeemCheck   = "RedeemCheck"
	// Create Coin
	AttributeSymbol      = "symbol"
	AttributeTitle       = "title"
	AttributeCRR         = "crr"
	AttributeInitVolume  = "initVolume"
	AttributeInitReserve = "initReserve"
	AttributeLimitVolume = "limitVolume"

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
