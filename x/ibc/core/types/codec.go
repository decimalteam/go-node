package types

import (
	clienttypes "bitbucket.org/decimalteam/go-node/x/ibc/core/02-client/types"
	connectiontypes "bitbucket.org/decimalteam/go-node/x/ibc/core/03-connection/types"
	channeltypes "bitbucket.org/decimalteam/go-node/x/ibc/core/04-channel/types"
	commitmenttypes "bitbucket.org/decimalteam/go-node/x/ibc/core/23-commitment/types"
	solomachinetypes "bitbucket.org/decimalteam/go-node/x/ibc/light-clients/06-solomachine/types"
	ibctmtypes "bitbucket.org/decimalteam/go-node/x/ibc/light-clients/07-tendermint/types"
	localhosttypes "bitbucket.org/decimalteam/go-node/x/ibc/light-clients/09-localhost/types"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
)

// RegisterInterfaces registers x/ibc interfaces into protobuf Any.
func RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	clienttypes.RegisterInterfaces(registry)
	connectiontypes.RegisterInterfaces(registry)
	channeltypes.RegisterInterfaces(registry)
	solomachinetypes.RegisterInterfaces(registry)
	ibctmtypes.RegisterInterfaces(registry)
	localhosttypes.RegisterInterfaces(registry)
	commitmenttypes.RegisterInterfaces(registry)
}
