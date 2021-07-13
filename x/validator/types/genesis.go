package types

import (
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
)

// GenesisState - all staking state that must be provided at genesis
//type GenesisState struct {
//	Params               Params                `json:"params" yaml:"params"`
//	LastTotalPower       sdk.Int               `json:"last_total_power" yaml:"last_total_power"`
//	LastValidatorPowers  []LastValidatorPower  `json:"last_validator_powers" yaml:"last_validator_powers"`
//	Validators           Validators            `json:"validators" yaml:"validators"`
//	Delegations          Delegations           `json:"delegations" yaml:"delegations"`
//	DelegationsNFT       DelegationsNFT        `json:"delegations_nft" yaml:"delegations_nft"`
//	UnbondingDelegations []UnbondingDelegation `json:"unbonding_delegations" yaml:"unbonding_delegations"`
//	Exported             bool                  `json:"exported" yaml:"exported"`
//}

// Last validator power, needed for validator set update logic
//type LastValidatorPower struct {
//	Address sdk.ValAddress
//	Power   int64
//}

func NewGenesisState(params Params, validators []Validator, delegations Delegations, delegationsNFT DelegationsNFT) *GenesisState {
	return &GenesisState{
		Params:         params,
		Validators:     validators,
		Delegations:    delegations,
		DelegationsNFT: delegationsNFT,
	}
}

// get raw genesis raw message for testing
func DefaultGenesisState() *GenesisState {
	return &GenesisState{
		Params: DefaultParams(),
	}
}

func (g GenesisState) UnpackInterfaces(c codectypes.AnyUnpacker) error {
	for i := range g.Validators {
		if err := g.Validators[i].UnpackInterfaces(c); err != nil {
			return err
		}
	}
	return nil
}
