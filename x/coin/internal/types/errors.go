package types

// Local code type
type CodeType = uint32

const (
	// Default coin codespace
	DefaultCodespace string = ModuleName
	// Create coin
	DecodeError             CodeType = 101
	InvalidCRR              CodeType = 102
	InvalidCoinSymbol       CodeType = 103
	CoinAlreadyExists       CodeType = 104
	InvalidCoinTitle        CodeType = 105
	InvalidCoinInitVolume   CodeType = 106
	InvalidCoinInitReserve  CodeType = 107
	CodeInvalid             CodeType = 108
	InsufficientCoinReserve CodeType = 118

	// Buy/Sell Coin
	SameCoins               CodeType = 109
	CoinToBuyNotExists      CodeType = 110
	CoinToSellNotExists     CodeType = 111
	InsufficientCoinToSell  CodeType = 112
	TxBreaksVolumeLimit     CodeType = 113
	TxBreaksMinReserveLimit CodeType = 114
	UpdateBalanceError      CodeType = 115
	AmountBuyIsTooSmall     CodeType = 116
	// Send Coin
	InvalidAmount CodeType = 117
)
