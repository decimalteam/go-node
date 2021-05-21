package coin

import (
	keeper2 "bitbucket.org/decimalteam/go-node/x/coin/keeper"
	types2 "bitbucket.org/decimalteam/go-node/x/coin/types"
)

const (
	// TODO: define constants that you would like exposed from the internal package

	ModuleName        = types2.ModuleName
	RouterKey         = types2.RouterKey
	StoreKey          = types2.StoreKey
	DefaultParamspace = types2.DefaultParamspace
	DefaultCodespace  = types2.DefaultCodespace
	//QueryParams       = types.QueryParams
	QuerierRoute = types2.QuerierRoute

	BuyCoinConst       = types2.BuyCoinConst
	SellCoinConst      = types2.SellCoinConst
	MultiSendCoinConst = types2.MultiSendCoinConst
	RedeemCheckConst   = types2.RedeemCheckConst
	SellAllConst       = types2.SellAllCoinConst
	CreateCoinConst    = types2.CreateCoinConst
	SendCoinConst      = types2.SendCoinConst
)

var (
	// functions aliases
	NewKeeper           = keeper2.NewKeeper
	NewQuerier          = keeper2.NewQuerier
	RegisterCodec       = types2.RegisterCodec
	NewGenesisState     = types2.NewGenesisState
	DefaultGenesisState = types2.DefaultGenesisState
	ValidateGenesis     = types2.ValidateGenesis
	NewMsgSendCoin      = types2.NewMsgSendCoin
	NewMsgBuyCoin       = types2.NewMsgBuyCoin
	NewMsgSellCoin      = types2.NewMsgSellCoin
	NewMsgCreateCoin    = types2.NewMsgCreateCoin
	NewMsgSellAllCoin   = types2.NewMsgSellAllCoin
	NewMsgMultiSendCoin = types2.NewMsgMultiSendCoin
	NewMsgRedeemCheck   = types2.NewMsgRedeemCheck
	NewMsgUpdateCoin    = types2.NewMsgUpdateCoin

	MinCoinReserve = types2.MinCoinReserve

	ErrTxBreaksMinReserveRule = types2.ErrTxBreaksMinReserveRule

	// variable aliases
	ModuleCdc = types2.ModuleCdc
	// TODO: Fill out variable aliases
)

type (
	Keeper           = keeper2.Keeper
	CodeType         = types2.CodeType
	GenesisState     = types2.GenesisState
	Params           = types2.Params
	Coin             = types2.Coin
	MsgSendCoin      = types2.MsgSendCoin
	MsgBuyCoin       = types2.MsgBuyCoin
	MsgSellCoin      = types2.MsgSellCoin
	MsgCreateCoin    = types2.MsgCreateCoin
	MsgSellAllCoin   = types2.MsgSellAllCoin
	MsgMultiSendCoin = types2.MsgMultiSendCoin
	MsgRedeemCheck   = types2.MsgRedeemCheck
	MsgUpdateCoin    = types2.MsgUpdateCoin
	Send             = types2.Send
)
