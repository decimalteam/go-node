package genutil

import (
	"github.com/cosmos/cosmos-sdk/codec"
	codec2 "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	"bitbucket.org/decimalteam/go-node/x/validator"
)

// ModuleCdc defines a generic sealed codec to be used throughout this module
var ModuleCdc *codec.LegacyAmino

func init() {
	ModuleCdc = codec.NewLegacyAmino()
	validator.RegisterCodec(ModuleCdc)
	authtypes.RegisterLegacyAminoCodec(ModuleCdc)
	sdk.RegisterLegacyAminoCodec(ModuleCdc)
	codec2.RegisterCrypto(ModuleCdc)
	ModuleCdc.Seal()
}
