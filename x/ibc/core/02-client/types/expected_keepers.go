package types

import (
	"bitbucket.org/decimalteam/go-node/x/validator/types"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// ValidatorKeeper expected validator keeper
type ValidatorKeeper interface {
	GetHistoricalInfo(ctx sdk.Context, height int64) (types.HistoricalInfo, bool)
	UnbondingTime(ctx sdk.Context) time.Duration
}
