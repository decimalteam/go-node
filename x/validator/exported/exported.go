package exported

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto"
	"time"
)

// DelegationI delegation bond for a delegated proof of stake system
type DelegationI interface {
	GetDelegatorAddr() sdk.AccAddress // delegator sdk.AccAddress for the bond
	GetValidatorAddr() sdk.ValAddress // validator operator address
	GetCoin() sdk.Coin
	GetTokensBase() sdk.Int
	SetTokensBase(sdk.Int) DelegationI
}

type UnbondingDelegationEntryI interface {
	GetCreationHeight() int64
	GetCompletionTime() time.Time
	GetBalance() sdk.Coin
	GetInitialBalance() sdk.Coin
	IsMature(currentTime time.Time) bool
}

// ValidatorI expected validator functions
type ValidatorI interface {
	IsJailed() bool               // whether the validator is jailed
	GetMoniker() string           // moniker of the validator
	IsBonded() bool               // check if has a bonded status
	IsUnbonded() bool             // check if has status unbonded
	IsUnbonding() bool            // check if has status unbonding
	GetOperator() sdk.ValAddress  // operator address to receive/return validators coins
	GetConsPubKey() crypto.PubKey // validation consensus pubkey
	GetConsAddr() sdk.ConsAddress // validation consensus address
	GetTokens() sdk.Int           // validation tokens
	GetBondedTokens() sdk.Int     // validator bonded tokens
	GetConsensusPower() int64     // validation power in tendermint
	GetCommission() sdk.Dec       // validator commission rate
}
