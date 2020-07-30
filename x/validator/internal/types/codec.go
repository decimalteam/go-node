package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

// RegisterCodec registers concrete types on codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgDeclareCandidate{}, "validator/declare_candidate", nil)
	cdc.RegisterConcrete(MsgDelegate{}, "validator/delegate", nil)
	cdc.RegisterConcrete(MsgUnbond{}, "validator/unbond", nil)
	cdc.RegisterConcrete(MsgEditCandidate{}, "validator/edit_candidate", nil)
	cdc.RegisterConcrete(MsgSetOnline{}, "validator/set_online", nil)
	cdc.RegisterConcrete(MsgSetOffline{}, "validator/set_offline", nil)
}

// ModuleCdc defines the module codec
var ModuleCdc *codec.Codec

func init() {
	ModuleCdc = codec.New()
	RegisterCodec(ModuleCdc)
	codec.RegisterCrypto(ModuleCdc)
	ModuleCdc.Seal()
}
