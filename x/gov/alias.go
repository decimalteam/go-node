package gov

import (
	keeper2 "bitbucket.org/decimalteam/go-node/x/gov/keeper"
	types2 "bitbucket.org/decimalteam/go-node/x/gov/types"
)

const (
	ModuleName         = types2.ModuleName
	DefaultParamspace  = types2.DefaultParamspace
	RouterKey          = types2.RouterKey
	QuerierRoute       = types2.QuerierRoute
	StoreKey           = types2.StoreKey
	StatusNil          = types2.StatusNil
	StatusWaiting      = types2.StatusWaiting
	StatusVotingPeriod = types2.StatusVotingPeriod
	StatusPassed       = types2.StatusPassed
	StatusRejected     = types2.StatusRejected
	StatusFailed       = types2.StatusFailed
)

var (
	DefaultGenesisState  = types2.DefaultGenesisState
	ValidateGenesis      = types2.ValidateGenesis
	NewKeeper            = keeper2.NewKeeper
	NewRouter            = types2.NewRouter
	RegisterCodec        = types2.RegisterCodec
	ModuleCdc            = types2.ModuleCdc
	NewQuerier           = keeper2.NewQuerier
	ParamKeyTable        = types2.ParamKeyTable
	NewMsgSubmitProposal = types2.NewMsgSubmitProposal
)

type (
	Keeper            = keeper2.Keeper
	GenesisState      = types2.GenesisState
	Proposal          = types2.Proposal
	Vote              = types2.Vote
	Votes             = types2.Votes
	VoteOption        = types2.VoteOption
	MsgSubmitProposal = types2.MsgSubmitProposal
	MsgVote           = types2.MsgVote
)
