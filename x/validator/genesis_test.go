package validator

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/ed25519"

	sdk "github.com/cosmos/cosmos-sdk/types"

	decsdk "bitbucket.org/decimalteam/go-node/utils/types"
	val "bitbucket.org/decimalteam/go-node/x/validator/internal/keeper"
	"bitbucket.org/decimalteam/go-node/x/validator/internal/types"
)

func TestInitGenesis(t *testing.T) {
	ctx, _, keeper, supplyKeeper, _ := val.CreateTestInput(t, false, 1000)

	valTokens := sdk.TokensFromConsensusPower(1)

	params := keeper.GetParams(ctx)
	validators := make([]types.Validator, 2)
	delegations := make([]types.Delegation, 2)

	// initialize the validators
	validators[0].ValAddress = decsdk.ValAddress(val.Addrs[0])
	validators[0].PubKey = val.PKs[0]
	validators[0].Status = types.Bonded
	validators[0].Tokens = valTokens
	validators[1].ValAddress = decsdk.ValAddress(val.Addrs[1])
	validators[1].PubKey = val.PKs[1]
	validators[1].Status = types.Bonded
	validators[1].Tokens = valTokens

	delegations[0].ValidatorAddress = validators[0].ValAddress
	delegations[0].DelegatorAddress = decsdk.AccAddress(validators[0].ValAddress)
	delegations[0].Coin = sdk.NewCoin(keeper.BondDenom(ctx), valTokens)
	delegations[1].ValidatorAddress = validators[1].ValAddress
	delegations[1].DelegatorAddress = decsdk.AccAddress(validators[1].ValAddress)
	delegations[1].Coin = sdk.NewCoin(keeper.BondDenom(ctx), valTokens)

	genesisState := types.NewGenesisState(params, validators, delegations)
	vals := InitGenesis(ctx, keeper, supplyKeeper, genesisState)

	actualGenesis := ExportGenesis(ctx, keeper)
	require.Equal(t, genesisState.Params, actualGenesis.Params)
	require.Equal(t, genesisState.Delegations, actualGenesis.Delegations)
	require.EqualValues(t, keeper.GetAllValidators(ctx), actualGenesis.Validators)

	// now make sure the validators are bonded and intra-tx counters are correct
	resVal, err := keeper.GetValidator(ctx, decsdk.ValAddress(val.Addrs[0]))
	require.NoError(t, err)
	require.Equal(t, types.Bonded, resVal.Status)

	resVal, err = keeper.GetValidator(ctx, decsdk.ValAddress(val.Addrs[1]))
	require.NoError(t, err)
	require.Equal(t, types.Bonded, resVal.Status)

	abcivals := make([]abci.ValidatorUpdate, len(vals))
	for i, validator := range validators {
		abcivals[i] = validator.ABCIValidatorUpdate(validator.Tokens)
	}

	require.Equal(t, abcivals, vals)
}

func TestInitGenesisLargeValidatorSet(t *testing.T) {
	size := 200

	ctx, _, keeper, supplyKeeper, _ := val.CreateTestInput(t, false, 1000)

	params := keeper.GetParams(ctx)
	delegations := make([]types.Delegation, size)
	validators := make([]types.Validator, size)

	for i := range validators {
		validators[i] = types.NewValidator(decsdk.ValAddress(val.Addrs[i]), val.PKs[i], sdk.ZeroDec(), val.Addrs[i], types.Description{})

		validators[i].Status = types.Bonded

		tokens := sdk.TokensFromConsensusPower(1)
		if i < 100 {
			tokens = sdk.TokensFromConsensusPower(2)
		}
		delegations[i] = types.NewDelegation(decsdk.AccAddress(validators[i].ValAddress), validators[i].ValAddress, sdk.NewCoin(keeper.BondDenom(ctx), tokens))
	}

	genesisState := types.NewGenesisState(params, validators, delegations)
	vals := InitGenesis(ctx, keeper, supplyKeeper, genesisState)

	abcivals := make([]abci.ValidatorUpdate, 16)
	for i, validator := range validators[:16] {
		abcivals[i] = validator.ABCIValidatorUpdate(keeper.TotalStake(ctx, validator))
	}

	require.Equal(t, abcivals, vals)
}

func TestValidateGenesis(t *testing.T) {
	genValidators1 := make([]types.Validator, 1, 5)
	pk := ed25519.GenPrivKey().PubKey()
	genValidators1[0] = types.NewValidator(decsdk.ValAddress(pk.Address()), pk, sdk.ZeroDec(), decsdk.AccAddress(pk.Address()), types.Description{})
	genValidators1[0].Tokens = sdk.OneInt()
	genValidators1[0].DelegatorShares = sdk.OneDec()

	tests := []struct {
		name    string
		mutate  func(*types.GenesisState)
		wantErr bool
	}{
		{"default", func(*types.GenesisState) {}, false},
		// validate genesis validators
		{"duplicate validator", func(data *types.GenesisState) {
			data.Validators = genValidators1
			data.Validators = append(data.Validators, genValidators1[0])
		}, true},
		{"no delegator shares", func(data *types.GenesisState) {
			data.Validators = genValidators1
			data.Validators[0].DelegatorShares = sdk.ZeroDec()
		}, true},
		{"jailed and bonded validator", func(data *types.GenesisState) {
			data.Validators = genValidators1
			data.Validators[0].Jailed = true
			data.Validators[0].Status = types.Bonded
		}, true},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			genesisState := types.DefaultGenesisState()
			tt.mutate(&genesisState)
			if tt.wantErr {
				assert.Error(t, ValidateGenesis(genesisState))
			} else {
				assert.NoError(t, ValidateGenesis(genesisState))
			}
		})
	}
}
