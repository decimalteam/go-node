package validator

import (
	val "bitbucket.org/decimalteam/go-node/x/validator/internal/keeper"
	"testing"

	"github.com/tendermint/tendermint/crypto/ed25519"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"

	"bitbucket.org/decimalteam/go-node/x/validator/internal/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestInitGenesis(t *testing.T) {
	ctx, _, keeper, supplyKeeper, _, _ := val.CreateTestInput(t, false, 1000)

	valTokens := TokensFromConsensusPower(1)

	params := keeper.GetParams(ctx)
	validators := make([]types.Validator, 2)
	delegations := make([]types.Delegation, 2)

	// initialize the validators
	validators[0].ValAddress = sdk.ValAddress(val.Addrs[0])
	validators[0].PubKey = val.PKs[0]
	validators[0].Status = types.Bonded
	validators[0].Online = true
	validators[0].Tokens = valTokens
	validators[1].ValAddress = sdk.ValAddress(val.Addrs[1])
	validators[1].PubKey = val.PKs[1]
	validators[1].Status = types.Bonded
	validators[1].Online = true
	validators[1].Tokens = valTokens

	delegations[0].ValidatorAddress = validators[0].ValAddress
	delegations[0].DelegatorAddress = sdk.AccAddress(validators[0].ValAddress)
	delegations[0].Coin = sdk.NewCoin(keeper.BondDenom(ctx), valTokens)
	delegations[0].TokensBase = valTokens
	delegations[1].ValidatorAddress = validators[1].ValAddress
	delegations[1].DelegatorAddress = sdk.AccAddress(validators[1].ValAddress)
	delegations[1].Coin = sdk.NewCoin(keeper.BondDenom(ctx), valTokens)
	delegations[1].TokensBase = valTokens

	genesisState := types.NewGenesisState(params, validators, delegations, nil)
	vals := InitGenesis(ctx, keeper, supplyKeeper, genesisState)

	actualGenesis := ExportGenesis(ctx, keeper)
	require.Equal(t, genesisState.Params, actualGenesis.Params)
	require.Equal(t, genesisState.Delegations, actualGenesis.Delegations)
	require.EqualValues(t, keeper.GetAllValidators(ctx), actualGenesis.Validators)

	// now make sure the validators are bonded and intra-tx counters are correct
	resVal, err := keeper.GetValidator(ctx, sdk.ValAddress(val.Addrs[0]))
	require.NoError(t, err)
	require.Equal(t, types.Bonded, resVal.Status)

	resVal, err = keeper.GetValidator(ctx, sdk.ValAddress(val.Addrs[1]))
	require.NoError(t, err)
	require.Equal(t, types.Bonded, resVal.Status)

	require.Equal(t, len(vals), len(validators))
	abcivals := make([]abci.ValidatorUpdate, len(vals))
	for i, validator := range validators {
		abcivals[i] = validator.ABCIValidatorUpdate()
	}

	require.Equal(t, abcivals, vals)
}

func TestValidateGenesis(t *testing.T) {
	genValidators1 := make([]types.Validator, 1, 5)
	pk := ed25519.GenPrivKey().PubKey()
	genValidators1[0] = types.NewValidator(sdk.ValAddress(pk.Address()), pk, sdk.ZeroDec(), sdk.AccAddress(pk.Address()), types.Description{})
	genValidators1[0].Tokens = sdk.OneInt()

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
