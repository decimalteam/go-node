package types

import (
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"bitbucket.org/decimalteam/go-node/x/ibc/core/exported"
)

// RegisterInterfaces register the ibc interfaces submodule implementations to protobuf
// Any.
func RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	registry.RegisterImplementations(
		(*exported.ClientState)(nil),
		&ClientState{},
	)
}
