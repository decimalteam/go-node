package validator

import (
	"bitbucket.org/decimalteam/go-node/x/validator/internal/keeper"
	"bitbucket.org/decimalteam/go-node/x/validator/internal/types"
)

const (
	// TODO: define constants that you would like exposed from the internal package

	ModuleName        = types.ModuleName
	RouterKey         = types.RouterKey
	StoreKey          = types.StoreKey
	TStoreKey         = types.TStoreKey
	DefaultParamSpace = keeper.DefaultParamspace
	DefaultCodespace  = types.DefaultCodespace
	QuerierRoute      = types.QuerierRoute
	NotBondedPoolName = types.NotBondedPoolName
	BondedPoolName    = types.BondedPoolName
	DefaultBondDenom  = types.DefaultBondDenom

	ValidatorsKey = types.ValidatorsKey

	DeclareCandidateConst = types.DeclareCandidateConst
	DelegateConst         = types.DelegateConst
	SetOnlineConst        = types.SetOnlineConst
	SetOfflineConst       = types.SetOfflineConst
	UnbondConst           = types.UnbondConst
	EditCandidateConst    = types.EditCandidateConst
)

var (
	// functions aliases
	NewKeeper           = keeper.NewKeeper
	NewQuerier          = keeper.NewQuerier
	RegisterCodec       = types.RegisterCodec
	DefaultGenesisState = types.DefaultGenesisState

	ErrCalculateCommission             = types.ErrCalculateCommission
	ErrUpdateBalance                   = types.ErrUpdateBalance
	ErrInsufficientFunds               = types.ErrInsufficientFunds
	ErrInsufficientCoinToPayCommission = types.ErrInsufficientCoinToPayCommission

	// variable aliases
	ModuleCdc = types.ModuleCdc
)

type (
	Keeper       = keeper.Keeper
	CodeType     = types.CodeType
	GenesisState = types.GenesisState
	Params       = types.Params

	MsgDeclareCandidate = types.MsgDeclareCandidate
	MsgDelegate         = types.MsgDelegate
	UnbondingDelegation = types.UnbondingDelegation

	Validator = types.Validator
)
