package capability

import (
	"bitbucket.org/decimalteam/go-node/x/capability/internal/keeper"
	"bitbucket.org/decimalteam/go-node/x/capability/internal/types"
)

const (
	ModuleName        = types.ModuleName
	RouterKey         = types.RouterKey
	StoreKey          = types.StoreKey
	DefaultParamspace = types.DefaultParamspace
	DefaultCodespace  = types.DefaultCodespace
	//QueryParams       = types.QueryParams
	QuerierRoute = types.QuerierRoute
)

var (
	// functions aliases
	NewKeeper           = keeper.NewKeeper
	NewQuerier          = keeper.NewQuerier
	RegisterCodec       = types.RegisterCodec
	NewGenesisState     = types.NewGenesisState
	DefaultGenesisState = types.DefaultGenesisState
	ValidateGenesis     = types.ValidateGenesis

	// variable aliases
	ModuleCdc = types.ModuleCdc
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
