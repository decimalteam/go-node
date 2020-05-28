package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

// RegisterCodec registers concrete types on codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgDeclareCandidate{}, "validator/declare_candidate", nil)
	cdc.RegisterConcrete(MsgDelegate{}, "validator/delegate", nil)
	cdc.RegisterConcrete(MsgUnbond{}, "validator/unbond", nil)
	cdc.RegisterConcrete(MsgEditCandidate{}, "validator/edit-candidate", nil)
	cdc.RegisterConcrete(MsgSetOnline{}, "validator/set-online", nil)
	cdc.RegisterConcrete(MsgSetOffline{}, "validator/set-offline", nil)
	cdc.RegisterConcrete(StdTx{}, "decimal/StdTx", nil)
}

// ModuleCdc defines the module codec
var ModuleCdc *codec.Codec

func init() {
	ModuleCdc = codec.New()
	RegisterCodec(ModuleCdc)
	codec.RegisterCrypto(ModuleCdc)
	ModuleCdc.Seal()
}
