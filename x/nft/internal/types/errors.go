package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	ErrInvalidCollection    = sdkerrors.Register(ModuleName, 1, "invalid NFT collection")
	ErrUnknownCollection    = sdkerrors.Register(ModuleName, 2, "unknown NFT collection")
	ErrInvalidNFT           = sdkerrors.Register(ModuleName, 3, "invalid NFT")
	ErrUnknownNFT           = sdkerrors.Register(ModuleName, 4, "unknown NFT")
	ErrNFTAlreadyExists     = sdkerrors.Register(ModuleName, 5, "NFT already exists")
	ErrEmptyMetadata        = sdkerrors.Register(ModuleName, 6, "NFT metadata can't be empty")
	ErrInvalidQuantity      = sdkerrors.Register(ModuleName, 7, "invalid NFT quantity")
	ErrInvalidReserve       = sdkerrors.Register(ModuleName, 8, "invalid NFT reserve")
	ErrNotAllowedBurn       = sdkerrors.Register(ModuleName, 9, "only the creator can burn a token")
	ErrNotAllowedMint       = sdkerrors.Register(ModuleName, 10, "only the creator can mint a token")
	ErrInvalidDenom         = sdkerrors.Register(ModuleName, 11, "invalid denom name")
	ErrInvalidTokenID       = sdkerrors.Register(ModuleName, 12, "invalid token name")
	ErrInvalidSubTokenID    = sdkerrors.Register(ModuleName, 13, "invalid subTokenID")
	ErrNotUniqueSubTokenIDs = sdkerrors.Register(ModuleName, 14, "subTokenIDs does not unique")
	ErrNotUniqueTokenURI    = sdkerrors.Register(ModuleName, 15, "tokenURI does not unique")
)
