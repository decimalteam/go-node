package keeper

import (
	"bitbucket.org/decimalteam/go-node/x/validator"
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

	cdc *codec.Codec // The amino codec for binary encoding/decoding.

	supplyKeeper    supply.Keeper
	validatorKeeper validator.Keeper
}

// NewKeeper creates new instances of the nft Keeper
func NewKeeper(cdc *codec.Codec, storeKey sdk.StoreKey, supplyKeeper supply.Keeper, validatorKeeper validator.Keeper) Keeper {
	return Keeper{
		storeKey:        storeKey,
		cdc:             cdc,
		supplyKeeper:    supplyKeeper,
		validatorKeeper: validatorKeeper,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}
