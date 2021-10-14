package coin

import (
	"bitbucket.org/decimalteam/go-node/x/coin/internal/keeper"
	"bitbucket.org/decimalteam/go-node/x/coin/internal/types"
)

const (
	// TODO: define constants that you would like exposed from the internal package

	ModuleName        = types.ModuleName
	RouterKey         = types.RouterKey
	StoreKey          = types.StoreKey
	DefaultParamspace = types.DefaultParamspace
	DefaultCodespace  = types.DefaultCodespace
	//QueryParams       = types.QueryParams
	QuerierRoute = types.QuerierRoute

	BuyCoinConst       = types.BuyCoinConst
	SellCoinConst      = types.SellCoinConst
	MultiSendCoinConst = types.MultiSendCoinConst
	RedeemCheckConst   = types.RedeemCheckConst
	SellAllConst       = types.SellAllCoinConst
	CreateCoinConst    = types.CreateCoinConst
	SendCoinConst      = types.SendCoinConst
)

var (
	// functions aliases
	NewKeeper           = keeper.NewKeeper
	NewQuerier          = keeper.NewQuerier
	RegisterCodec       = types.RegisterCodec
	NewGenesisState     = types.NewGenesisState
	DefaultGenesisState = types.DefaultGenesisState
	ValidateGenesis     = types.ValidateGenesis
	NewMsgSendCoin      = types.NewMsgSendCoin
	NewMsgBuyCoin       = types.NewMsgBuyCoin
	NewMsgSellCoin      = types.NewMsgSellCoin
	NewMsgCreateCoin    = types.NewMsgCreateCoin
	NewMsgSellAllCoin   = types.NewMsgSellAllCoin
	NewMsgMultiSendCoin = types.NewMsgMultiSendCoin
	NewMsgRedeemCheck   = types.NewMsgRedeemCheck
	NewMsgUpdateCoin    = types.NewMsgUpdateCoin

	MinCoinReserve = types.MinCoinReserve

	ErrTxBreaksMinReserveRule = types.ErrTxBreaksMinReserveRule

	// variable aliases
	ModuleCdc = types.ModuleCdc
	// TODO: Fill out variable aliases
)

type (
	Keeper           = keeper.Keeper
	CodeType         = types.CodeType
	GenesisState     = types.GenesisState
	Params           = types.Params
	Coin             = types.Coin
	MsgSendCoin      = types.MsgSendCoin
	MsgBuyCoin       = types.MsgBuyCoin
	MsgSellCoin      = types.MsgSellCoin
	MsgCreateCoin    = types.MsgCreateCoin
	MsgSellAllCoin   = types.MsgSellAllCoin
	MsgMultiSendCoin = types.MsgMultiSendCoin
	MsgRedeemCheck   = types.MsgRedeemCheck
	MsgUpdateCoin    = types.MsgUpdateCoin
	Send             = types.Send
)
