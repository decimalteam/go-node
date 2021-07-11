package types

import "github.com/cosmos/cosmos-sdk/codec"

// RegisterCodec registers concrete types on codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgHTLT{}, "swap/msg_htlt", nil)
	cdc.RegisterConcrete(MsgRedeem{}, "swap/msg_redeem", nil)
	cdc.RegisterConcrete(MsgRefund{}, "swap/msg_refund", nil)
}

// ModuleCdc defines the module codec
var ModuleCdc *codec.Codec

func init() {
	ModuleCdc = codec.New()
	RegisterCodec(ModuleCdc)
	codec.RegisterCrypto(ModuleCdc)
	ModuleCdc.Seal()
}
