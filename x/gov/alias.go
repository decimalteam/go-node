package gov

import (
	"bitbucket.org/decimalteam/go-node/x/gov/internal/keeper"
	"bitbucket.org/decimalteam/go-node/x/gov/internal/types"
)

const (
	ModuleName         = types.ModuleName
	DefaultParamspace  = types.DefaultParamspace
	RouterKey          = types.RouterKey
	QuerierRoute       = types.QuerierRoute
	StoreKey           = types.StoreKey
	StatusNil          = types.StatusNil
	StatusWaiting      = types.StatusWaiting
	StatusVotingPeriod = types.StatusVotingPeriod
	StatusPassed       = types.StatusPassed
	StatusRejected     = types.StatusRejected
	StatusFailed       = types.StatusFailed
)

var (
	DefaultGenesisState  = types.DefaultGenesisState
	ValidateGenesis      = types.ValidateGenesis
	NewKeeper            = keeper.NewKeeper
	NewRouter            = types.NewRouter
	RegisterCodec        = types.RegisterCodec
	ModuleCdc            = types.ModuleCdc
	NewQuerier           = keeper.NewQuerier
	ParamKeyTable        = types.ParamKeyTable
	NewMsgSubmitProposal = types.NewMsgSubmitProposal
)

type (
	Keeper            = keeper.Keeper
	GenesisState      = types.GenesisState
	Proposal          = types.Proposal
	Vote              = types.Vote
	Votes             = types.Votes
	VoteOption        = types.VoteOption
	MsgSubmitProposal = types.MsgSubmitProposal
	MsgVote           = types.MsgVote
)
