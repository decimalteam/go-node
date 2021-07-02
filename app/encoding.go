package app

import (
	"bitbucket.org/decimalteam/go-node/app/params"
	"github.com/cosmos/cosmos-sdk/std"
)

// MakeEncodingConfig creates an EncodingConfig for testing. This function
// should be used only in tests or when creating a new app instance (NewApp*()).
// App user shouldn't create new codecs - use the app.AppCodec instead.
// [DEPRECATED]
func MakeEncodingConfig() params.EncodingConfig {
	encodingConfig := params.NewEncodingConfig()
	std.RegisterLegacyAminoCodec(encodingConfig.Amino)
	std.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	ModuleBasics.RegisterLegacyAminoCodec(encodingConfig.Amino)
	ModuleBasics.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	return encodingConfig
}