package app

import (
	"encoding/json"
	"log"

	abci "github.com/tendermint/tendermint/abci/types"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	decsdk "bitbucket.org/decimalteam/go-node/utils/types"
	"bitbucket.org/decimalteam/go-node/x/validator"
)

// ExportAppStateAndValidators exports the state of the application for a genesis
// file.
func (app *newApp) ExportAppStateAndValidators(
	forZeroHeight bool, jailWhiteList []string,
) (appState json.RawMessage, validators []tmtypes.GenesisValidator, err error) {

	// as if they could withdraw from the start of the next block
	ctx := app.NewContext(true, abci.Header{Height: app.LastBlockHeight()})

	if forZeroHeight {
		app.prepForZeroHeightGenesis(ctx, jailWhiteList)
	}

	genState := app.mm.ExportGenesis(ctx)
	appState, err = codec.MarshalJSONIndent(app.cdc, genState)
	if err != nil {
		return nil, nil, err
	}

	validators = validator.WriteValidators(ctx, app.validatorKeeper)
	return appState, validators, nil
}

// prepare for fresh start at zero height
// NOTE zero height genesis is a temporary feature which will be deprecated
//      in favour of export at a block height
func (app *newApp) prepForZeroHeightGenesis(ctx sdk.Context, jailWhiteList []string) {
	applyWhiteList := false

	//Check if there is a whitelist
	if len(jailWhiteList) > 0 {
		applyWhiteList = true
	}

	whiteListMap := make(map[string]bool)

	for _, addr := range jailWhiteList {
		_, err := decsdk.ValAddressFromPrefixedHex(addr)
		if err != nil {
			log.Fatal(err)
		}
		whiteListMap[addr] = true
	}

	/* Handle fee distribution state. */

	// withdraw all validator commission

	// withdraw all delegator rewards

	// clear validator slash events
	// clear validator historical rewards
	// set context height to zero
	height := ctx.BlockHeight()
	ctx = ctx.WithBlockHeight(0)

	// reinitialize all validators

	// reinitialize all delegations

	// reset context height
	ctx = ctx.WithBlockHeight(height)

	/* Handle validator state. */

	// iterate through redelegations, reset creation height

	// iterate through unbonding delegations, reset creation height
	app.validatorKeeper.IterateUnbondingDelegations(ctx, func(_ int64, ubd validator.UnbondingDelegation) (stop bool) {
		for i := range ubd.Entries {
			ubd.Entries[i].CreationHeight = 0
		}
		app.validatorKeeper.SetUnbondingDelegation(ctx, ubd)
		return false
	})

	// Iterate through validators by power descending, reset bond heights, and
	// update bond intra-tx counters.
	store := ctx.KVStore(app.keys[validator.StoreKey])
	iter := sdk.KVStoreReversePrefixIterator(store, []byte{validator.ValidatorsKey})
	counter := int16(0)

	for ; iter.Valid(); iter.Next() {
		addr := decsdk.ValAddress(iter.Key()[1:])
		validator, err := app.validatorKeeper.GetValidator(ctx, addr)
		if err != nil {
			panic("expected validator, not found")
		}

		validator.UnbondingHeight = 0
		if applyWhiteList && !whiteListMap[addr.String()] {
			validator.Jailed = true
		}

		app.validatorKeeper.SetValidator(ctx, validator)
		counter++
	}

	iter.Close()

	_, err := app.validatorKeeper.ApplyAndReturnValidatorSetUpdates(ctx)
	if err != nil {
		panic(err)
	}

	/* Handle slashing state. */

	// reset start height on signing infos
}
