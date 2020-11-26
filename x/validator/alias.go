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

	DAOAddress1 = keeper.DAOAddress1
	DAOAddress2 = keeper.DAOAddress2
	DAOAddress3 = keeper.DAOAddress3

	DevelopAddress1 = keeper.DevelopAddress1
	DevelopAddress2 = keeper.DevelopAddress2
	DevelopAddress3 = keeper.DevelopAddress3

	Unbonded  = types.Unbonded
	Unbonding = types.Unbonding
	Bonded    = types.Bonded
)

var (
	// functions aliases
	NewKeeper           = keeper.NewKeeper
	NewQuerier          = keeper.NewQuerier
	RegisterCodec       = types.RegisterCodec
	DefaultGenesisState = types.DefaultGenesisState

	NewMsgDeclareCandidate = types.NewMsgDeclareCandidate
	NewMsgEditCandidate    = types.NewMsgEditCandidate
	NewMsgDelegate         = types.NewMsgDelegate
	NewMsgUnbond           = types.NewMsgUnbond
	NewMsgSetOnline        = types.NewMsgSetOnline
	NewMsgSetOffline       = types.NewMsgSetOffline

	NewValidator = types.NewValidator

	ErrCalculateCommission             = types.ErrCalculateCommission
	ErrUpdateBalance                   = types.ErrUpdateBalance
	ErrInsufficientFunds               = types.ErrInsufficientFunds
	ErrInsufficientCoinToPayCommission = types.ErrInsufficientCoinToPayCommission

	DefaultParams = types.DefaultParams

	TokensFromConsensusPower = types.TokensFromConsensusPower
	TokensToConsensusPower   = types.TokensToConsensusPower

	// variable aliases
	ModuleCdc = types.ModuleCdc
)

type (
	Keeper       = keeper.Keeper
	CodeType     = types.CodeType
	GenesisState = types.GenesisState
	Params       = types.Params

	MsgDeclareCandidate = types.MsgDeclareCandidate
	MsgEditCandidate    = types.MsgEditCandidate
	MsgDelegate         = types.MsgDelegate
	MsgUnbond           = types.MsgUnbond
	MsgSetOnline        = types.MsgSetOnline
	MsgSetOffline       = types.MsgSetOffline

	UnbondingDelegation = types.UnbondingDelegation

	Validator = types.Validator

	Description = types.Description
)
