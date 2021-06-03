package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params/types"
)

type ParamSubspace interface {
	WithKeyTable(table types.KeyTable) types.Subspace
	Get(ctx sdk.Context, key []byte, ptr interface{})
	GetParamSet(ctx sdk.Context, ps types.ParamSet)
	SetParamSet(ctx sdk.Context, ps types.ParamSet)
}
