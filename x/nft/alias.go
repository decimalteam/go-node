// nolint
// autogenerated code using github.com/rigelrozanski/multitool
// aliases generated for the following subdirectories:
// ALIASGEN: github.com/cosmos/modules/incubator/nft/keeper
// ALIASGEN: github.com/cosmos/modules/incubator/nft/types
package nft

import (
	"bitbucket.org/decimalteam/go-node/x/nft/internal/keeper"
	"bitbucket.org/decimalteam/go-node/x/nft/internal/types"
)

const (
	QuerySupply       = keeper.QuerySupply
	QueryOwner        = keeper.QueryOwner
	QueryOwnerByDenom = keeper.QueryOwnerByDenom
	QueryCollection   = keeper.QueryCollection
	QueryDenoms       = keeper.QueryDenoms
	QueryNFT          = keeper.QueryNFT
	ReservedPool      = types.ReservedPool
	ModuleName        = types.ModuleName
	StoreKey          = types.StoreKey
	QuerierRoute      = types.QuerierRoute
	RouterKey         = types.RouterKey
)

var (
	// functions aliases
	RegisterInvariants      = keeper.RegisterInvariants
	AllInvariants           = keeper.AllInvariants
	SupplyInvariant         = keeper.SupplyInvariant
	NewKeeper               = keeper.NewKeeper
	NewQuerier              = keeper.NewQuerier
	RegisterCodec           = types.RegisterCodec
	NewCollection           = types.NewCollection
	EmptyCollection         = types.EmptyCollection
	NewCollections          = types.NewCollections
	ErrInvalidCollection    = types.ErrInvalidCollection
	ErrUnknownCollection    = types.ErrUnknownCollection
	ErrInvalidNFT           = types.ErrInvalidNFT
	ErrNFTAlreadyExists     = types.ErrNFTAlreadyExists
	ErrUnknownNFT           = types.ErrUnknownNFT
	ErrEmptyMetadata        = types.ErrEmptyMetadata
	ErrNotAllowedBurn       = types.ErrNotAllowedBurn
	ErrNotAllowedUpdateRes  = types.ErrNotAllowedUpdateReserve
	ErrNotAllowedMint       = types.ErrNotAllowedMint
	ErrNotUniqueSubTokenIDs = types.ErrNotUniqueSubTokenIDs
	ErrNotUniqueTokenURI    = types.ErrNotUniqueTokenURI
	ErrNotUniqueTokenID     = types.ErrNotUniqueTokenID
	NewGenesisState         = types.NewGenesisState
	DefaultGenesisState     = types.DefaultGenesisState
	ValidateGenesis         = types.ValidateGenesis
	NewBaseNFT              = types.NewBaseNFT
	NewNFTs                 = types.NewNFTs
	NewMsgMintNFT           = types.NewMsgMintNFT
	NewMsgBurnNFT           = types.NewMsgBurnNFT
	NewMsgTranfserNFT       = types.NewMsgTransferNFT
	NewMsgEditNFTMetadata   = types.NewMsgEditNFTMetadata

	CheckUnique = types.CheckUnique

	// variable aliases
	ModuleCdc = types.ModuleCdc
)

type (
	Keeper                = keeper.Keeper
	Collection            = types.Collection
	Collections           = types.Collections
	CollectionJSON        = types.CollectionJSON
	GenesisState          = types.GenesisState
	MsgTransferNFT        = types.MsgTransferNFT
	MsgEditNFTMetadata    = types.MsgEditNFTMetadata
	MsgMintNFT            = types.MsgMintNFT
	MsgBurnNFT            = types.MsgBurnNFT
	BaseNFT               = types.BaseNFT
	NFTs                  = types.NFTs
	NFTJSON               = types.NFTJSON
	IDCollection          = types.IDCollection
	IDCollections         = types.IDCollections
	Owner                 = types.Owner
	TokenOwner            = types.TokenOwner
	QueryCollectionParams = types.QueryCollectionParams
	QueryBalanceParams    = types.QueryBalanceParams
	QueryNFTParams        = types.QueryNFTParams
	SortedIntArray        = types.SortedIntArray
)
