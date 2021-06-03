package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	crypto "github.com/cosmos/cosmos-sdk/crypto/codec"
)

// RegisterCodec registers concrete types on codec
func RegisterCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(MsgCreateWallet{}, "multisig/create_wallet", nil)
	cdc.RegisterConcrete(MsgCreateTransaction{}, "multisig/create_transaction", nil)
	cdc.RegisterConcrete(MsgSignTransaction{}, "multisig/sign_transaction", nil)
}

// ModuleCdc defines the module codec
var ModuleCdc *codec.LegacyAmino

func init() {
	ModuleCdc = codec.NewLegacyAmino()
	RegisterCodec(ModuleCdc)
	crypto.RegisterCrypto(ModuleCdc)
	ModuleCdc.Seal()
}
