package types

// coin module event types
const (
	EventTypeCreateCoin = "CreateCoin"
	EventTypeBuyCoin    = "BuyCoin"
	// Create Coin
	AttributeSymbol      = "symbol"
	AttributeTitle       = "title"
	AttributeCRR         = "crr"
	AttributeInitVolume  = "initVolume"
	AttributeInitReserve = "initReserve"
	AttributeLimitVolume = "limitVolume"

	// Buy Coin
	AttributeCoinToBuy       = "coinToBuy"
	AttributeCoinToSell      = "coinToSell"
	AttributeAmountToBuy     = "amountToBuy"
	AttributeMaxAmountToSell = "maxAmountToSell"

	AttributeValueCategory = ModuleName
)
