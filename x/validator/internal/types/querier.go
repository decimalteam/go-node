package types

import (
	decsdk "bitbucket.org/decimalteam/go-node/utils/types"
)

// query endpoints supported by the validator Querier
const (
	QueryValidators                    = "validators"
	QueryValidator                     = "validator"
	QueryDelegatorDelegations          = "delegatorDelegations"
	QueryDelegatorUnbondingDelegations = "delegatorUnbondingDelegations"
	QueryRedelegations                 = "redelegations"
	QueryValidatorDelegations          = "validatorDelegations"
	QueryValidatorRedelegations        = "validatorRedelegations"
	QueryValidatorUnbondingDelegations = "validatorUnbondingDelegations"
	QueryDelegation                    = "delegation"
	QueryUnbondingDelegation           = "unbondingDelegation"
	QueryDelegatorValidators           = "delegatorValidators"
	QueryDelegatorValidator            = "delegatorValidator"
	QueryPool                          = "pool"
	QueryParameters                    = "parameters"
	QueryHistoricalInfo                = "historicalInfo"
)

// defines the params for the following queries:
// - 'custom/validator/delegatorDelegations'
// - 'custom/validator/delegatorUnbondingDelegations'
// - 'custom/validator/delegatorRedelegations'
// - 'custom/validator/delegatorValidators'
type QueryDelegatorParams struct {
	DelegatorAddr decsdk.AccAddress
}

func NewQueryDelegatorParams(delegatorAddr decsdk.AccAddress) QueryDelegatorParams {
	return QueryDelegatorParams{
		DelegatorAddr: delegatorAddr,
	}
}

// defines the params for the following queries:
// - 'custom/validator/validator'
// - 'custom/validator/validatorDelegations'
// - 'custom/validator/validatorUnbondingDelegations'
// - 'custom/validator/validatorRedelegations'
type QueryValidatorParams struct {
	ValidatorAddr decsdk.ValAddress
}

func NewQueryValidatorParams(validatorAddr decsdk.ValAddress) QueryValidatorParams {
	return QueryValidatorParams{
		ValidatorAddr: validatorAddr,
	}
}

// defines the params for the following queries:
// - 'custom/validator/delegation'
// - 'custom/validator/unbondingDelegation'
// - 'custom/validator/delegatorValidator'
type QueryBondsParams struct {
	DelegatorAddr decsdk.AccAddress
	ValidatorAddr decsdk.ValAddress
}

func NewQueryBondsParams(delegatorAddr decsdk.AccAddress, validatorAddr decsdk.ValAddress) QueryBondsParams {
	return QueryBondsParams{
		DelegatorAddr: delegatorAddr,
		ValidatorAddr: validatorAddr,
	}
}

// defines the params for the following queries:
// - 'custom/validator/redelegation'
type QueryRedelegationParams struct {
	DelegatorAddr    decsdk.AccAddress
	SrcValidatorAddr decsdk.ValAddress
	DstValidatorAddr decsdk.ValAddress
}

func NewQueryRedelegationParams(delegatorAddr decsdk.AccAddress,
	srcValidatorAddr, dstValidatorAddr decsdk.ValAddress) QueryRedelegationParams {

	return QueryRedelegationParams{
		DelegatorAddr:    delegatorAddr,
		SrcValidatorAddr: srcValidatorAddr,
		DstValidatorAddr: dstValidatorAddr,
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
