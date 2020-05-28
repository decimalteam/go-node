package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// RegisterCodec registers concrete types on codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgCreateCoin{}, "coin/CreateCoin", nil)
	cdc.RegisterConcrete(MsgBuyCoin{}, "coin/BuyCoin", nil)
	cdc.RegisterConcrete(MsgSellCoin{}, "coin/SellCoin", nil)
	cdc.RegisterConcrete(MsgSellAllCoin{}, "coin/SellAllCoin", nil)
	cdc.RegisterConcrete(MsgSendCoin{}, "coin/SendCoin", nil)
	cdc.RegisterConcrete(MsgMultiSendCoin{}, "coin/MultiSendCoin", nil)
	cdc.RegisterConcrete(MsgRedeemCheck{}, "coin/RedeemCheck", nil)
}

// ModuleCdc defines the module codec
var ModuleCdc *codec.Codec

func init() {
	ModuleCdc = codec.New()
	RegisterCodec(ModuleCdc)
	codec.RegisterCrypto(ModuleCdc)
	ModuleCdc.Seal()
}

type FeeCoinTx struct {
	sdk.Tx
	FeeCoin string
}
