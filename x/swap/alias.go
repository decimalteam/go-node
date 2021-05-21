package swap

import (
	keeper2 "bitbucket.org/decimalteam/go-node/x/swap/keeper"
	types2 "bitbucket.org/decimalteam/go-node/x/swap/types"
)

const (
	ModuleName        = types2.ModuleName
	StoreKey          = types2.StoreKey
	RouterKey         = types2.RouterKey
	QuerierRoute      = types2.QuerierRoute
	DefaultParamspace = keeper2.DefaultParamspace

	MsgHTLTConst   = types2.TypeMsgHTLT
	MsgRedeemConst = types2.TypeMsgRedeem
	MsgRefundConst = types2.TypeMsgRefund

	PoolName = types2.PoolName
)

type (
	Keeper       = keeper2.Keeper
	MsgHTLT      = types2.MsgHTLT
	MsgRedeem    = types2.MsgRedeem
	MsgRefund    = types2.MsgRefund
	GenesisState = types2.GenesisState
)

var (
	ModuleCdc = types2.ModuleCdc

	RegisterCodec       = types2.RegisterCodec
	DefaultGenesisState = types2.DefaultGenesisState

	SwapServiceAddress = types2.SwapServiceAddress

	NewKeeper = keeper2.NewKeeper

	NewMsgRedeem = types2.NewMsgRedeem
	NewMsgHTLT   = types2.NewMsgHTLT
	NewMsgRefund = types2.NewMsgRefund
)
