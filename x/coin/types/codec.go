package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	codec2 "github.com/cosmos/cosmos-sdk/crypto/codec"
)

// RegisterCodec registers concrete types on codec
func RegisterCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(MsgCreateCoin{}, "coin/create_coin", nil)
	cdc.RegisterConcrete(MsgUpdateCoin{}, "coin/update_coin", nil)
	cdc.RegisterConcrete(MsgBuyCoin{}, "coin/buy_coin", nil)
	cdc.RegisterConcrete(MsgSellCoin{}, "coin/sell_coin", nil)
	cdc.RegisterConcrete(MsgSellAllCoin{}, "coin/sell_all_coin", nil)
	cdc.RegisterConcrete(MsgSendCoin{}, "coin/send_coin", nil)
	cdc.RegisterConcrete(MsgMultiSendCoin{}, "coin/multi_send_coin", nil)
	cdc.RegisterConcrete(MsgRedeemCheck{}, "coin/redeem_check", nil)
}

// ModuleCdc defines the module codec
var ModuleCdc *codec.LegacyAmino

func init() {
	ModuleCdc = codec.NewLegacyAmino()
	RegisterCodec(ModuleCdc)
	codec2.RegisterCrypto(ModuleCdc)
	ModuleCdc.Seal()
}
