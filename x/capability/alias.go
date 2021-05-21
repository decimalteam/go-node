package capability

import (
	"bitbucket.org/decimalteam/go-node/x/capability/internal/types"
	keeper2 "bitbucket.org/decimalteam/go-node/x/capability/keeper"
	types2 "bitbucket.org/decimalteam/go-node/x/capability/types"
)

const (
	ModuleName        = types2.ModuleName
	RouterKey         = types.RouterKey
	StoreKey          = types2.StoreKey
	DefaultParamspace = types.DefaultParamspace
	DefaultCodespace  = types.DefaultCodespace
	//QueryParams       = types.QueryParams
	QuerierRoute = types.QuerierRoute
)

var (
	// functions aliases
	NewKeeper           = keeper2.NewKeeper
	NewQuerier          = keeper2.NewQuerier
	RegisterCodec       = types2.RegisterCodec
	NewGenesisState     = types2.NewGenesisState
	DefaultGenesisState = types2.DefaultGenesisState
	ValidateGenesis     = types2.ValidateGenesis

	// variable aliases
	ModuleCdc = types2.ModuleCdc
)

type (
	Keeper           = keeper2.Keeper
	CodeType         = types.CodeType
	GenesisState     = types2.GenesisState
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
