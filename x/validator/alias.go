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
	DefaultParamSpace = keeper.DefaultParamSpace
	DefaultCodespace  = types.DefaultCodespace
	QuerierRoute      = types.QuerierRoute
	NotBondedPoolName = types.NotBondedPoolName
	BondedPoolName    = types.BondedPoolName
)

var (
	// functions aliases
	NewKeeper           = keeper.NewKeeper
	NewQuerier          = keeper.NewQuerier
	RegisterCodec       = types.RegisterCodec
	NewGenesisState     = types.NewGenesisState
	DefaultGenesisState = types.DefaultGenesisState
	// TODO: Fill out function aliases

	// variable aliases
	ModuleCdc = types.ModuleCdc
	// TODO: Fill out variable aliases

	ErrValidatorOwnerExists  = types.ErrValidatorOwnerExists
	ErrValidatorPubKeyExists = types.ErrValidatorPubKeyExists
	ErrInvalidStruct         = types.ErrInvalidStruct
)

type (
	Keeper       = keeper.Keeper
	CodeType     = types.CodeType
	GenesisState = types.GenesisState
	Params       = types.Params

	// TODO: Fill out module types
)
