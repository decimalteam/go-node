package types

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	tmtypes "github.com/tendermint/tendermint/types"
)

func TestValidatorTestEquivalent(t *testing.T) {
	val1 := NewValidator(valAddr1, pk1, sdk.ZeroDec(), sdk.AccAddress(valAddr1), Description{})
	val2 := NewValidator(valAddr1, pk1, sdk.ZeroDec(), sdk.AccAddress(valAddr1), Description{})

	ok := val1.TestEquivalent(val2)
	require.True(t, ok)

	val2 = NewValidator(valAddr2, pk2, sdk.ZeroDec(), sdk.AccAddress(valAddr2), Description{})

	ok = val1.TestEquivalent(val2)
	require.False(t, ok)
}

func TestABCIValidatorUpdate(t *testing.T) {
	validator := NewValidator(valAddr1, pk1, sdk.ZeroDec(), sdk.AccAddress(valAddr1), Description{})

	abciVal := validator.ABCIValidatorUpdate(sdk.ZeroInt())
	require.Equal(t, tmtypes.TM2PB.PubKey(validator.PubKey), abciVal.PubKey)
	require.Equal(t, int64(0), abciVal.Power)
}

func TestABCIValidatorUpdateZero(t *testing.T) {
	validator := NewValidator(valAddr1, pk1, sdk.ZeroDec(), sdk.AccAddress(valAddr1), Description{})

	abciVal := validator.ABCIValidatorUpdateZero()
	require.Equal(t, tmtypes.TM2PB.PubKey(validator.PubKey), abciVal.PubKey)
	require.Equal(t, int64(0), abciVal.Power)
}

func TestUpdateStatus(t *testing.T) {
	validator := NewValidator(pk1.Address().Bytes(), pk1, sdk.ZeroDec(), pk1.Address().Bytes(), Description{})
	require.Equal(t, Unbonded, validator.Status)

	// Unbonded to Bonded
	validator = validator.UpdateStatus(Bonded)
	require.Equal(t, Bonded, validator.Status)

	// Bonded to Unbonding
	validator = validator.UpdateStatus(Unbonding)
	require.Equal(t, Unbonding, validator.Status)

	// Unbonding to Bonded
	validator = validator.UpdateStatus(Bonded)
	require.Equal(t, Bonded, validator.Status)
}

func TestValidatorMarshalUnmarshalJSON(t *testing.T) {
	validator := NewValidator(valAddr1, pk1, sdk.ZeroDec(), sdk.AccAddress(valAddr1), Description{})
	js, err := codec.Cdc.MarshalJSON(validator)
	require.NoError(t, err)
	require.NotEmpty(t, js)
	require.Contains(t, string(js), "\"pub_key\":\"cosmosvalconspub")
	got := &Validator{}
	err = codec.Cdc.UnmarshalJSON(js, got)
	assert.NoError(t, err)
	assert.Equal(t, validator, *got)
}

//func TestValidatorMarshalYAML(t *testing.T) {
//	validator := NewValidator(valAddr1, pk1, sdk.ZeroDec(), sdk.AccAddress(valAddr1))
//	bechifiedPub, err := sdk.Bech32ifyPubKey(Bech32PrefixConsPub, validator.PubKey)
//	require.NoError(t, err)
//	bs, err := yaml.Marshal(validator)
//	require.NoError(t, err)
//	want := fmt.Sprintf(`|
//  valaddress: %s
//  rewardaddress: %s
//  pubkey: %s
//  jailed: false
//  status: 0
//  tokens: "0"
//  delegatorshares: "0.000000000000000000"
//  description:
//    moniker: ""
//    identity: ""
//    website: ""
//    details: ""
//  unbondingheight: 0
//  unbondingcompletiontime: 1970-01-01T00:00:00Z
//  commission: "0.000000000000000000"
//  accumrewards: "0"
//  online: true
//`, validator.ValAddress.String(), validator.RewardAddress.String(), bechifiedPub)
//	require.Equal(t, want, string(bs))
//}
