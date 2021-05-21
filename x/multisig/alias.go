package multisig

import (
	keeper2 "bitbucket.org/decimalteam/go-node/x/multisig/keeper"
	types2 "bitbucket.org/decimalteam/go-node/x/multisig/types"
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

	CreateTransactionConst = types2.CreateTransactionConst
	CreateWalletConst      = types2.CreateWalletConst
	SignTransactionConst   = types2.SignTransactionConst
)

var (
	// functions aliases
	NewKeeper               = keeper2.NewKeeper
	NewQuerier              = keeper2.NewQuerier
	RegisterCodec           = types2.RegisterCodec
	NewGenesisState         = types2.NewGenesisState
	DefaultGenesisState     = types2.DefaultGenesisState
	ValidateGenesis         = types2.ValidateGenesis
	NewMsgCreateWallet      = types2.NewMsgCreateWallet
	NewMsgCreateTransaction = types2.NewMsgCreateTransaction
	NewMsgSignTransaction   = types2.NewMsgSignTransaction
	NewWallet               = types2.NewWallet
	NewTransaction          = types2.NewTransaction

	// variable aliases
	ModuleCdc = types2.ModuleCdc
	// TODO: Fill out variable aliases
)

type (
	Keeper               = keeper2.Keeper
	CodeType             = types2.CodeType
	GenesisState         = types2.GenesisState
	Params               = types2.Params
	MsgCreateWallet      = types2.MsgCreateWallet
	MsgCreateTransaction = types2.MsgCreateTransaction
	MsgSignTransaction   = types2.MsgSignTransaction
	QueryWallets         = types2.QueryWallets
	QueryTransactions    = types2.QueryTransactions
	Wallet               = types2.Wallet
	Transaction          = types2.Transaction
)
