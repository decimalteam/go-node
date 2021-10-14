package types

import sdk "github.com/cosmos/cosmos-sdk/types"

// query endpoints supported by the validator Querier
const (
	QueryValidators                    = "validators"
	QueryValidator                     = "validator"
	QueryDelegatorDelegations          = "delegatorDelegations"
	QueryDelegatorUnbondingDelegations = "delegatorUnbondingDelegations"
	QueryValidatorDelegations          = "validatorDelegations"
	QueryValidatorUnbondingDelegations = "validatorUnbondingDelegations"
	QueryDelegation                    = "delegation"
	QueryUnbondingDelegation           = "unbondingDelegation"
	QueryDelegatorValidators           = "delegatorValidators"
	QueryDelegatorValidator            = "delegatorValidator"
	QueryPool                          = "pool"
	QueryParameters                    = "parameters"
	QueryHistoricalInfo                = "historicalInfo"
	QueryDelegatedCoins                = "delegatedCoins"
	QueryDelegatedCoin                 = "delegatedCoin"
)

// QueryDelegatorParams defines the params for the following queries:
// - 'custom/validator/delegatorDelegations'
// - 'custom/validator/delegatorUnbondingDelegations'
// - 'custom/validator/delegatorRedelegations'
// - 'custom/validator/delegatorValidators'
type QueryDelegatorParams struct {
	DelegatorAddr sdk.AccAddress
}

func NewQueryDelegatorParams(delegatorAddr sdk.AccAddress) QueryDelegatorParams {
	return QueryDelegatorParams{
		DelegatorAddr: delegatorAddr,
	}
}

// QueryValidatorParams defines the params for the following queries:
// - 'custom/validator/validator'
// - 'custom/validator/validatorDelegations'
// - 'custom/validator/validatorUnbondingDelegations'
// - 'custom/validator/validatorRedelegations'
type QueryValidatorParams struct {
	ValidatorAddr sdk.ValAddress
}

func NewQueryValidatorParams(validatorAddr sdk.ValAddress) QueryValidatorParams {
	return QueryValidatorParams{
		ValidatorAddr: validatorAddr,
	}
}

// QueryBondsParams defines the params for the following queries:
// - 'custom/validator/delegation'
// - 'custom/validator/unbondingDelegation'
// - 'custom/validator/delegatorValidator'
type QueryBondsParams struct {
	DelegatorAddr sdk.AccAddress
	ValidatorAddr sdk.ValAddress
	Coin          string
}

func NewQueryBondsParams(delegatorAddr sdk.AccAddress, validatorAddr sdk.ValAddress, coin string) QueryBondsParams {
	return QueryBondsParams{
		DelegatorAddr: delegatorAddr,
		ValidatorAddr: validatorAddr,
		Coin:          coin,
	}
}

// QueryValidatorsParams defines the params for the following queries:
// - 'custom/validator/validators'
type QueryValidatorsParams struct {
	Page, Limit int
	Status      string
}

func NewQueryValidatorsParams(page, limit int, status string) QueryValidatorsParams {
	return QueryValidatorsParams{page, limit, status}
}

// QueryHistoricalInfoParams defines the params for the following queries:
// - 'custom/validator/historicalInfo'
type QueryHistoricalInfoParams struct {
	Height int64
}

// NewQueryHistoricalInfoParams creates a new QueryHistoricalInfoParams instance
func NewQueryHistoricalInfoParams(height int64) QueryHistoricalInfoParams {
	return QueryHistoricalInfoParams{height}
}
