package multisig

import (
	"bitbucket.org/decimalteam/go-node/x/multisig/internal/keeper"
	"bitbucket.org/decimalteam/go-node/x/multisig/internal/types"
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
)

var (
	// functions aliases
	NewKeeper               = keeper.NewKeeper
	NewQuerier              = keeper.NewQuerier
	RegisterCodec           = types.RegisterCodec
	NewGenesisState         = types.NewGenesisState
	DefaultGenesisState     = types.DefaultGenesisState
	ValidateGenesis         = types.ValidateGenesis
	NewMsgCreateWallet      = types.NewMsgCreateWallet
	NewMsgCreateTransaction = types.NewMsgCreateTransaction
	NewMsgSignTransaction   = types.NewMsgSignTransaction
	NewWallet               = types.NewWallet
	NewTransaction          = types.NewTransaction

	// variable aliases
	ModuleCdc = types.ModuleCdc
	// TODO: Fill out variable aliases
)

type (
	Keeper               = keeper.Keeper
	CodeType             = types.CodeType
	GenesisState         = types.GenesisState
	Params               = types.Params
	MsgCreateWallet      = types.MsgCreateWallet
	MsgCreateTransaction = types.MsgCreateTransaction
	MsgSignTransaction   = types.MsgSignTransaction
	Wallet               = types.Wallet
	Transaction          = types.Transaction
)
