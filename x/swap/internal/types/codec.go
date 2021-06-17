package types

import "github.com/cosmos/cosmos-sdk/codec"

// RegisterCodec registers concrete types on codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgHTLT{}, "swap/msg_htlt", nil)
	cdc.RegisterConcrete(MsgRedeem{}, "swap/msg_redeem", nil)
	cdc.RegisterConcrete(MsgRefund{}, "swap/msg_refund", nil)
	cdc.RegisterConcrete(MsgSwapInitialize{}, "swap/msg_initialize", nil)
	cdc.RegisterConcrete(MsgRedeemV2{}, "swap/msg_redeem_v2", nil)
	cdc.RegisterConcrete(MsgChainActivate{}, "swap/msg_chain_activate", nil)
	cdc.RegisterConcrete(MsgChainDeactivate{}, "swap/msg_chain_deactivate", nil)
}

// ModuleCdc defines the module codec
var ModuleCdc *codec.Codec

func init() {
	ModuleCdc = codec.New()
	RegisterCodec(ModuleCdc)
	codec.RegisterCrypto(ModuleCdc)
	ModuleCdc.Seal()
}
