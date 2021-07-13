package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	codec2 "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

// RegisterCodec registers concrete types on codec
func RegisterCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(MsgHTLT{}, "swap/msg_htlt", nil)
	cdc.RegisterConcrete(MsgRedeem{}, "swap/msg_redeem", nil)
	cdc.RegisterConcrete(MsgRefund{}, "swap/msg_refund", nil)
}

func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgHTLT{},
		&MsgRedeem{},
		&MsgRefund{},
	)

	// todo implement msg_service
	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

// ModuleCdc defines the module codec
var ModuleCdc *codec.LegacyAmino

func init() {
	ModuleCdc = codec.NewLegacyAmino()
	RegisterCodec(ModuleCdc)
	codec2.RegisterCrypto(ModuleCdc)
	ModuleCdc.Seal()
}
