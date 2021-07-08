package app

import (
	"bitbucket.org/decimalteam/go-node/x/validator"
	"encoding/json"
	"github.com/cosmos/cosmos-sdk/codec/types"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	"log"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// ExportAppStateAndValidators exports the state of the application for a genesis
// file.
func (app *newApp) ExportAppStateAndValidators(
	forZeroHeight bool, jailWhiteList []string,
) (servertypes.ExportedApp, error) {

	// as if they could withdraw from the start of the next block
	ctx := app.NewContext(true, tmproto.Header{Height: app.LastBlockHeight()})

	height := app.LastBlockHeight() + 1
	if forZeroHeight {
		height = 0
		app.prepForZeroHeightGenesis(ctx, jailWhiteList)
	}

	genState := app.mm.ExportGenesis(ctx, app.appCodec)
	appState, err := json.MarshalIndent(genState, "", " ")
	if err != nil {
		return servertypes.ExportedApp{}, err
	}

	validators := validator.WriteValidators(ctx, app.validatorKeeper)

	return servertypes.ExportedApp{
		AppState:        appState,
		Validators:      validators,
		Height:          height,
		ConsensusParams: app.BaseApp.GetConsensusParams(ctx),
	}, nil
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
		_, err := sdk.ValAddressFromBech32(addr)
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
		for i, any := range ubd.Entries {
			entry := any.GetCachedValue().(validator.UnbondingDelegationEntry)

			delegationEntry := &validator.UnbondingDelegationEntry{
				CreationHeight: 0,
				CompletionTime: entry.GetCompletionTime(),
				InitialBalance: entry.GetInitialBalance(),
				Balance:        entry.GetBalance(),
			}
			entryAny, err := types.NewAnyWithValue(delegationEntry)
			if err != nil {
				continue
			}

			ubd.Entries[i] = entryAny
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
		addr := sdk.ValAddress(iter.Key()[1:])
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
