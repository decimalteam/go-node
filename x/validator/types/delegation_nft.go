package types

import (
	"bitbucket.org/decimalteam/go-node/x/validator/exported"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"time"
)

//type DelegationNFT struct {
//	DelegatorAddress sdk.AccAddress `json:"delegator_address" yaml:"delegator_address"`
//	ValidatorAddress sdk.ValAddress `json:"validator_address" yaml:"validator_address"`
//	Denom            string         `json:"denom" yaml:"denom"`
//	TokenID          string         `json:"token_id" yaml:"token_id"`
//	Quantity         sdk.Int        `json:"quantity" yaml:"quantity"`
//	Coin             sdk.Coin       `json:"coin" yaml:"coin"`
//}

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

func MustMarshalDelegationNFT(cdc *codec.LegacyAmino, delegation DelegationNFT) []byte {
	return cdc.MustMarshalLengthPrefixed(delegation)
}

// return the delegation
func MustUnmarshalDelegationNFT(cdc *codec.LegacyAmino, value []byte) DelegationNFT {
	delegation, err := UnmarshalDelegationNFT(cdc, value)
	if err != nil {
		panic(err)
	}
	return delegation
}

// return the delegation
func UnmarshalDelegationNFT(cdc *codec.LegacyAmino, value []byte) (delegation DelegationNFT, err error) {
	err = cdc.UnmarshalLengthPrefixed(value, &delegation)
	return delegation, err
}

func (d DelegationNFT) GetDelegatorAddr() sdk.AccAddress             { return d.DelegatorAddress }
func (d DelegationNFT) GetValidatorAddr() sdk.ValAddress             { return d.ValidatorAddress }
func (d DelegationNFT) GetCoin() sdk.Coin                            { return d.Coin }
func (d DelegationNFT) GetTokensBase() sdk.Int                       { return d.Coin.Amount }
func (d DelegationNFT) SetTokensBase(_ sdk.Int) exported.DelegationI { return d }

type DelegationsNFT []DelegationNFT

//type UnbondingDelegationNFTEntry struct {
//	CreationHeight int64     `json:"creation_height" yaml:"creation_height"` // height which the unbonding took place
//	CompletionTime time.Time `json:"completion_time" yaml:"completion_time"` // time at which the unbonding delegation will complete
//	Denom          string    `json:"denom" yaml:"denom"`
//	TokenID        string    `json:"token_id" yaml:"token_id"`
//	Quantity       sdk.Int   `json:"quantity" yaml:"quantity"`
//	Balance        sdk.Coin  `json:"balance"`
//}

func NewUnbondingDelegationNFTEntry(creationHeight int64, completionTime time.Time, denom string, tokenID string, quantity sdk.Int, balance sdk.Coin) UnbondingDelegationNFTEntry {
	return UnbondingDelegationNFTEntry{
		CreationHeight: creationHeight,
		CompletionTime: completionTime,
		Denom:          denom,
		TokenID:        tokenID,
		Quantity:       quantity,
		Balance:        balance,
	}
}

func (u UnbondingDelegationNFTEntry) GetCreationHeight() int64     { return u.CreationHeight }
func (u UnbondingDelegationNFTEntry) GetCompletionTime() time.Time { return u.CompletionTime }

func (u UnbondingDelegationNFTEntry) GetBalance() sdk.Coin {
	return u.Balance
}

func (u UnbondingDelegationNFTEntry) GetInitialBalance() sdk.Coin {
	return u.Balance
}

func (u UnbondingDelegationNFTEntry) IsMature(currentTime time.Time) bool {
	return !u.CompletionTime.After(currentTime)
}

func (u UnbondingDelegationNFTEntry) GetEvent() sdk.Event {
	return sdk.NewEvent(
		EventTypeCompleteUnbondingNFT,
		sdk.NewAttribute(AttributeKeyDenom, u.Denom),
		sdk.NewAttribute(AttributeKeyID, u.TokenID),
		sdk.NewAttribute(AttributeKeyQuantity, u.Quantity.String()),
	)
}
