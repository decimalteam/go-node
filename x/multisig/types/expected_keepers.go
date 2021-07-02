package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params/types"
)

// ParamSubspace defines the expected Subspace interfacace
type ParamSubspace interface {
	WithKeyTable(table types.KeyTable) types.Subspace
	Get(ctx sdk.Context, key []byte, ptr interface{})
	GetParamSet(ctx sdk.Context, ps types.ParamSet)
	SetParamSet(ctx sdk.Context, ps types.ParamSet)
}

// TODO: Create interfaces of what you expect the other keepers to have to be able to use this module.