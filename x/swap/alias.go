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

	MsgHTLTConst   = types.TypeMsgHTLT
	MsgRedeemConst = types.TypeMsgRedeem
	MsgRefundConst = types.TypeMsgRefund

	PoolName = types.PoolName
)

type (
	Keeper       = keeper.Keeper
	MsgHTLT      = types.MsgHTLT
	MsgRedeem    = types.MsgRedeem
	MsgRefund    = types.MsgRefund
	GenesisState = types.GenesisState
)

var (
	ModuleCdc = types.ModuleCdc

	RegisterCodec       = types.RegisterCodec
	DefaultGenesisState = types.DefaultGenesisState

	SwapServiceAccAddress = types.SwapServiceAccAddress
	SwapServiceAddress    = types.SwapServiceAddress

	NewKeeper = keeper.NewKeeper

	NewMsgRedeem = types.NewMsgRedeem
	NewMsgHTLT   = types.NewMsgHTLT
	NewMsgRefund = types.NewMsgRefund
)
