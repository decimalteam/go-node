package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type DelegationNFT struct {
	DelegatorAddress sdk.AccAddress `json:"delegator_address" yaml:"delegator_address"`
	ValidatorAddress sdk.ValAddress `json:"validator_address" yaml:"validator_address"`
	Denom            string         `json:"denom" yaml:"denom"`
	TokenID          string         `json:"token_id" yaml:"token_id"`
	Quantity         sdk.Int        `json:"quantity" yaml:"quantity"`
	Coin             sdk.Coin       `json:"coin" yaml:"coin"`
}

func NewDelegationNFT(delegatorAddr sdk.AccAddress, validatorAddr sdk.ValAddress, tokenID, denom string, quantity sdk.Int, coin sdk.Coin) DelegationNFT {
	return DelegationNFT{
		DelegatorAddress: delegatorAddr,
		ValidatorAddress: validatorAddr,
		Denom:            denom,
		TokenID:          tokenID,
		Quantity:         quantity,
		Coin:             coin,
	}
}

func MustMarshalDelegationNFT(cdc *codec.Codec, delegation DelegationNFT) []byte {
	return cdc.MustMarshalBinaryLengthPrefixed(delegation)
}

// return the delegation
func MustUnmarshalDelegationNFT(cdc *codec.Codec, value []byte) DelegationNFT {
	delegation, err := UnmarshalDelegationNFT(cdc, value)
	if err != nil {
		panic(err)
	}
	return delegation
}

// return the delegation
func UnmarshalDelegationNFT(cdc *codec.Codec, value []byte) (delegation DelegationNFT, err error) {
	err = cdc.UnmarshalBinaryLengthPrefixed(value, &delegation)
	return delegation, err
}

func (d DelegationNFT) GetDelegatorAddr() sdk.AccAddress { return d.DelegatorAddress }
func (d DelegationNFT) GetValidatorAddr() sdk.ValAddress { return d.ValidatorAddress }
func (d DelegationNFT) GetCoin() sdk.Coin                { return d.Coin }
