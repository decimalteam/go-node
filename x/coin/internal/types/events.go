package types

// coin module event types
const (
	EventTypeCreateCoin = "CreateCoin"
	EventTypeBuyCoin    = "BuyCoin"
	EventTypeSellCoin   = "SellCoin"
	// Create Coin
	AttributeSymbol      = "symbol"
	AttributeTitle       = "title"
	AttributeCRR         = "crr"
	AttributeInitVolume  = "initVolume"
	AttributeInitReserve = "initReserve"
	AttributeLimitVolume = "limitVolume"

	// Buy/Sell Coin
	AttributeCoinToBuy       = "coinToBuy"
	AttributeCoinToSell      = "coinToSell"
	AttributeAmountToBuy     = "amountToBuy"
	AttributeMaxAmountToSell = "maxAmountToSell"
	AttributeAmountToSell    = "amountToSell"
	AttributeMinAmountToBuy  = "minAmountToBuy"

	AttributeValueCategory = ModuleName
)
