package swap

import (
	"bitbucket.org/decimalteam/go-node/x/swap/internal/keeper"
	"bitbucket.org/decimalteam/go-node/x/swap/internal/types"
)

const (
	ModuleName        = types.ModuleName
	StoreKey          = types.StoreKey
	RouterKey         = types.RouterKey
	QuerierRoute      = types.QuerierRoute
	DefaultParamspace = keeper.DefaultParamspace

	PoolName = types.PoolName
)

type (
	Keeper       = keeper.Keeper
	MsgHTLT      = types.MsgHTLT
	GenesisState = types.GenesisState
)

var (
	ModuleCdc = types.ModuleCdc

	RegisterCodec       = types.RegisterCodec
	DefaultGenesisState = types.DefaultGenesisState

	NewKeeper = keeper.NewKeeper
)
