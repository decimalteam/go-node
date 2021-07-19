package keeper

import (
	keeper2 "bitbucket.org/decimalteam/go-node/x/ibc/core/02-client/keeper"
	"bitbucket.org/decimalteam/go-node/x/ibc/core/02-client/types"
	"bitbucket.org/decimalteam/go-node/x/ibc/core/03-connection/keeper"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

	channelkeeper "bitbucket.org/decimalteam/go-node/x/ibc/core/04-channel/keeper"
	porttypes "bitbucket.org/decimalteam/go-node/x/ibc/core/05-port/types"
	"github.com/cosmos/cosmos-sdk/codec"
	capabilitykeeper "github.com/cosmos/cosmos-sdk/x/capability/keeper"
	portkeeper "bitbucket.org/decimalteam/go-node/x/ibc/core/05-port/keeper"
)

//var _ types.QueryServer = (*Keeper)(nil)

// Keeper defines each ICS keeper for IBC
type Keeper struct {
	// implements gRPC QueryServer interface
	//types.QueryServer

	cdc codec.BinaryMarshaler

	ClientKeeper     keeper2.Keeper
	ConnectionKeeper keeper.Keeper
	ChannelKeeper    channelkeeper.Keeper
	PortKeeper       portkeeper.Keeper
	Router           *porttypes.Router
}

// NewKeeper creates a new ibc Keeper
func NewKeeper(
	cdc codec.BinaryMarshaler, key sdk.StoreKey, paramSpace paramtypes.Subspace,
	validatorKeeper types.ValidatorKeeper, scopedKeeper capabilitykeeper.ScopedKeeper,
) *Keeper {
	clientKeeper := keeper2.NewKeeper(cdc, key, paramSpace, validatorKeeper)
	connectionKeeper := keeper.NewKeeper(cdc, key, clientKeeper)
	portKeeper := portkeeper.NewKeeper(scopedKeeper)
	channelKeeper := channelkeeper.NewKeeper(cdc, key, clientKeeper, connectionKeeper, portKeeper, scopedKeeper)

	return &Keeper{
		cdc:              cdc,
		ClientKeeper:     clientKeeper,
		ConnectionKeeper: connectionKeeper,
		ChannelKeeper:    channelKeeper,
		PortKeeper:       portKeeper,
	}
}

// Codec returns the IBC module codec.
func (k Keeper) Codec() codec.BinaryMarshaler {
	return k.cdc
}

// SetRouter sets the Router in IBC Keeper and seals it. The method panics if
// there is an existing router that's already sealed.
func (k *Keeper) SetRouter(rtr *porttypes.Router) {
	if k.Router != nil && k.Router.Sealed() {
		panic("cannot reset a sealed router")
	}
	k.Router = rtr
	k.Router.Seal()
}
