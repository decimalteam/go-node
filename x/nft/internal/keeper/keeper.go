package keeper

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/x/supply"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"bitbucket.org/decimalteam/go-node/x/nft/internal/types"
)

// Keeper maintains the link to data storage and exposes getter/setter methods for the various parts of the state machine
type Keeper struct {
	storeKey sdk.StoreKey // Unexposed key to access store from sdk.Context

	cdc *codec.LegacyAmino // The amino codec for binary encoding/decoding.

	supplyKeeper supply.Keeper

	baseDenom string
}

// NewKeeper creates new instances of the nft Keeper
func NewKeeper(cdc *codec.LegacyAmino, storeKey sdk.StoreKey, supplyKeeper supply.Keeper, baseDenom string) Keeper {
	return Keeper{
		storeKey:     storeKey,
		cdc:          cdc,
		supplyKeeper: supplyKeeper,
		baseDenom:    baseDenom,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}
