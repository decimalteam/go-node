package validator

import (
	vtypes "bitbucket.org/decimalteam/go-node/x/validator/internal/types"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking/types"
	abci "github.com/tendermint/tendermint/abci/types"
	tmtypes "github.com/tendermint/tendermint/types"
	"log"
)

// InitGenesis sets the pool and parameters for the provided keeper.  For each
// validator in data, it sets that validator in the keeper along with manually
// setting the indexes. In addition, it also sets any delegations found in
// data. Finally, it updates the bonded validators.
// Returns final validator set after applying all declaration and delegations
func InitGenesis(ctx sdk.Context, keeper Keeper,
	supplyKeeper types.SupplyKeeper, data GenesisState) []abci.ValidatorUpdate {

	var updates []abci.ValidatorUpdate
	bondedTokens := make(map[string]sdk.Int)
	notBondedTokens := make(map[string]sdk.Int)

	// We need to pretend to be "n blocks before genesis", where "n" is the
	// validator update delay, so that e.g. slashing periods are correctly
	// initialized for the validator set e.g. with a one-block offset - the
	// first TM block is at height 1, so state updates applied from
	// genesis.json are in block 0.
	ctx = ctx.WithBlockHeight(1 - sdk.ValidatorUpdateDelay)

	keeper.SetParams(ctx, data.Params)
	err := keeper.SetLastTotalPower(ctx, data.LastTotalPower)
	if err != nil {
		panic(fmt.Sprintln("Init genesis error: ", err))
	}

	for _, validator := range data.Validators {
		err := keeper.SetValidator(ctx, validator)
		if err != nil {
			log.Println("Init genesis error: ", err)
			continue
		}

		// Manually set indices for the first time
		keeper.SetValidatorByConsAddr(ctx, validator)
		keeper.SetValidatorByPowerIndex(ctx, validator)
		// update timeslice if necessary
		if validator.IsUnbonding() {
			err = keeper.InsertValidatorQueue(ctx, validator)
			if err != nil {
				log.Println("Init genesis error: ", err)
				continue
			}
		}

		switch validator.Status {
		case vtypes.Bonded:
			bondedTokens[vtypes.DefaultBondDenom] = validator.Tokens
		case vtypes.Unbonding, vtypes.Unbonded:
			notBondedTokens[vtypes.DefaultBondDenom] = validator.Tokens
		default:
			log.Println("Init genesis error: invalid validator status")
		}
	}

	//for _, delegation := range data.Delegations {
	//	// Call the before-creation hook if not exported
	//	if !data.Exported {
	//		keeper.BeforeDelegationCreated(ctx, delegation.DelegatorAddress, delegation.ValidatorAddress)
	//	}
	//	keeper.SetDelegation(ctx, delegation)
	//
	//	// Call the after-modification hook if not exported
	//	if !data.Exported {
	//		keeper.AfterDelegationModified(ctx, delegation.DelegatorAddress, delegation.ValidatorAddress)
	//	}
	//}
	//
	//for _, ubd := range data.UnbondingDelegations {
	//	keeper.SetUnbondingDelegation(ctx, ubd)
	//	for _, entry := range ubd.Entries {
	//		keeper.InsertUBDQueue(ctx, ubd, entry.CompletionTime)
	//		notBondedTokens = notBondedTokens.Add(entry.Balance)
	//	}
	//}
	//
	//for _, red := range data.Redelegations {
	//	keeper.SetRedelegation(ctx, red)
	//	for _, entry := range red.Entries {
	//		keeper.InsertRedelegationQueue(ctx, red, entry.CompletionTime)
	//	}
	//}

	bondedCoins := sdk.NewCoins()
	for denom, amount := range bondedTokens {
		bondedCoins = append(bondedCoins, sdk.NewCoin(denom, amount))
	}
	notBondedCoins := sdk.NewCoins()
	for denom, amount := range notBondedTokens {
		notBondedCoins = append(notBondedCoins, sdk.NewCoin(denom, amount))
	}

	// check if the unbonded and bonded pools accounts exists
	bondedPool := keeper.GetBondedPool(ctx)
	if bondedPool == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.BondedPoolName))
	}

	// TODO remove with genesis 2-phases refactor https://github.com/cosmos/cosmos-sdk/issues/2862
	// add coins if not provided on genesis
	if bondedPool.GetCoins().IsZero() {
		if err := bondedPool.SetCoins(bondedCoins); err != nil {
			panic(err)
		}
		supplyKeeper.SetModuleAccount(ctx, bondedPool)
	}

	notBondedPool := keeper.GetNotBondedPool(ctx)
	if notBondedPool == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.NotBondedPoolName))
	}

	if notBondedPool.GetCoins().IsZero() {
		if err := notBondedPool.SetCoins(notBondedCoins); err != nil {
			panic(err)
		}
		supplyKeeper.SetModuleAccount(ctx, notBondedPool)
	}

	// don't need to run Tendermint updates if we exported
	if data.Exported {
		for _, lv := range data.LastValidatorPowers {
			err = keeper.SetLastValidatorPower(ctx, lv.Address, lv.Power)
			if err != nil {
				panic(fmt.Sprintln("Init genesis error: ", err))
			}
			validator, err := keeper.GetValidator(ctx, lv.Address)
			if err != nil {
				panic(fmt.Sprintf("validator %s not found", lv.Address))
			}
			update := validator.ABCIValidatorUpdate(keeper.TotalStake(ctx, validator))
			update.Power = lv.Power // keep the next-val-set offset, use the last power for the first block
			updates = append(updates, update)
		}
	} else {
		updates, err = keeper.ApplyAndReturnValidatorSetUpdates(ctx)
		if err != nil {
			panic(fmt.Sprintln("Init genesis error: ", err))
		}
	}

	return updates
}

// ExportGenesis returns a GenesisState for a given context and keeper. The
// GenesisState will contain the pool, params, validators, and bonds found in
// the keeper.
func ExportGenesis(ctx sdk.Context, keeper Keeper) vtypes.GenesisState {
	params := keeper.GetParams(ctx)
	lastTotalPower := keeper.GetLastTotalPower(ctx)
	validators := keeper.GetAllValidators(ctx)
	//delegations := keeper.GetAllDelegations(ctx)
	//var unbondingDelegations []types.UnbondingDelegation
	//keeper.IterateUnbondingDelegations(ctx, func(_ int64, ubd types.UnbondingDelegation) (stop bool) {
	//	unbondingDelegations = append(unbondingDelegations, ubd)
	//	return false
	//})
	//var redelegations []types.Redelegation
	//keeper.IterateRedelegations(ctx, func(_ int64, red types.Redelegation) (stop bool) {
	//	redelegations = append(redelegations, red)
	//	return false
	//})
	var lastValidatorPowers []vtypes.LastValidatorPower
	keeper.IterateLastValidatorPowers(ctx, func(addr sdk.ValAddress, power int64) (stop bool) {
		lastValidatorPowers = append(lastValidatorPowers, vtypes.LastValidatorPower{Address: addr, Power: power})
		return false
	})

	return vtypes.GenesisState{
		Params:              params,
		LastTotalPower:      lastTotalPower,
		LastValidatorPowers: lastValidatorPowers,
		Validators:          validators,
		Exported:            true,
	}
}

// WriteValidators returns a slice of bonded genesis validators.
func WriteValidators(ctx sdk.Context, keeper Keeper) (vals []tmtypes.GenesisValidator) {
	keeper.IterateLastValidators(ctx, func(_ int64, validator vtypes.Validator) (stop bool) {
		vals = append(vals, tmtypes.GenesisValidator{
			PubKey: validator.PubKey,
			Power:  validator.ConsensusPower(keeper.TotalStake(ctx, validator)),
			Name:   validator.Description.Moniker,
		})

		return false
	})

	return
}

// ValidateGenesis validates the provided staking genesis state to ensure the
// expected invariants holds. (i.e. params in correct bounds, no duplicate validators)
func ValidateGenesis(data vtypes.GenesisState) error {
	err := validateGenesisStateValidators(data.Validators)
	if err != nil {
		return err
	}
	err = data.Params.Validate()
	if err != nil {
		return err
	}

	return nil
}

func validateGenesisStateValidators(validators []vtypes.Validator) (err error) {
	addrMap := make(map[string]bool, len(validators))
	for i := 0; i < len(validators); i++ {
		val := validators[i]
		strKey := string(val.PubKey.Bytes())
		if _, ok := addrMap[strKey]; ok {
			return fmt.Errorf("duplicate validator in genesis state: moniker %v, address %v", val.Description.Moniker, val.ValAddress)
		}
		if val.Jailed && val.IsBonded() {
			return fmt.Errorf("validator is bonded and jailed in genesis state: moniker %v, address %v", val.Description.Moniker, val.ValAddress)
		}
		if val.DelegatorShares.IsZero() && !val.IsUnbonding() {
			return fmt.Errorf("bonded/unbonded genesis validator cannot have zero delegator shares, validator: %v", val)
		}
		addrMap[strKey] = true
	}
	return
}
