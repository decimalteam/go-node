package types

import (
	"fmt"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

type CodeType = uint32

const (
	// Default coin codespace
	DefaultCodespace string = ModuleName

	CodeInvalidCollection CodeType = 101
	CodeUnknownCollection CodeType = 102
	CodeInvalidNFT        CodeType = 103
	CodeUnknownNFT        CodeType = 104
	CodeNFTAlreadyExists  CodeType = 105
	CodeEmptyMetadata     CodeType = 106
	CodeInvalidQuantity   CodeType = 107
	CodeInvalidReserve    CodeType = 108
	CodeNotAllowedBurn    CodeType = 109
	CodeNotAllowedMint    CodeType = 110
	CodeInvalidDenom      CodeType = 111
	CodeInvalidTokenID    CodeType = 112
)

var (
//	ErrInvalidCollection = sdkerrors.Register(ModuleName, 1, "invalid NFT collection")
//ErrUnknownCollection = sdkerrors.Register(ModuleName, 2, "unknown NFT collection")
//ErrInvalidNFT        = sdkerrors.Register(ModuleName, 3, "invalid NFT")
//	ErrUnknownNFT        = sdkerrors.Register(ModuleName, 4, "unknown NFT")
//ErrNFTAlreadyExists  = sdkerrors.Register(ModuleName, 5, "NFT already exists")
//ErrEmptyMetadata     = sdkerrors.Register(ModuleName, 6, "NFT metadata can't be empty")
//ErrInvalidQuantity   = sdkerrors.Register(ModuleName, 7, "invalid NFT quantity")
//ErrInvalidReserve    = sdkerrors.Register(ModuleName, 8, "invalid NFT reserve")
//ErrNotAllowedBurn    = sdkerrors.Register(ModuleName, 9, "only the creator can burn a token")
//ErrNotAllowedMint    = sdkerrors.Register(ModuleName, 10, "only the creator can mint a token")
//	ErrInvalidDenom      = sdkerrors.Register(ModuleName, 11, "invalid denom name")
//	ErrInvalidTokenID    = sdkerrors.Register(ModuleName, 12, "invalid token name")
)

func ErrInvalidCollection(denom string) *sdkerrors.Error {
	return sdkerrors.New(
		DefaultCodespace,
		CodeInvalidCollection,
		fmt.Sprintf("invalid NFT collection: %s", denom),
	)
}

func ErrUnknownCollection() *sdkerrors.Error {
	return sdkerrors.New(
		DefaultCodespace,
		CodeUnknownCollection,
		fmt.Sprintf("unknown NFT collection"),
	)
}

func ErrInvalidNFT() *sdkerrors.Error {
	return sdkerrors.New(
		DefaultCodespace,
		CodeInvalidNFT,
		fmt.Sprintf("invalid NFT"),
	)
}

func ErrUnknownNFT() *sdkerrors.Error {
	return sdkerrors.New(
		DefaultCodespace,
		CodeUnknownNFT,
		fmt.Sprintf("unknown NFT"),
	)
}

func ErrNFTAlreadyExists() *sdkerrors.Error {
	return sdkerrors.New(
		DefaultCodespace,
		CodeNFTAlreadyExists,
		fmt.Sprintf("NFT already exists"),
	)
}

func ErrEmptyMetadata() *sdkerrors.Error {
	return sdkerrors.New(
		DefaultCodespace,
		CodeEmptyMetadata,
		fmt.Sprintf("NFT metadata can't be empty"),
	)
}

func ErrInvalidQuantity() *sdkerrors.Error {
	return sdkerrors.New(
		DefaultCodespace,
		CodeInvalidQuantity,
		fmt.Sprintf("invalid NFT quantity"),
	)
}

func ErrInvalidReserve() *sdkerrors.Error {
	return sdkerrors.New(
		DefaultCodespace,
		CodeInvalidReserve,
		fmt.Sprintf("invalid NFT reserve"),
	)
}

func ErrNotAllowedBurn() *sdkerrors.Error {
	return sdkerrors.New(
		DefaultCodespace,
		CodeNotAllowedBurn,
		fmt.Sprintf("only the creator can burn a token"),
	)
}

func ErrNotAllowedMint() *sdkerrors.Error {
	return sdkerrors.New(
		DefaultCodespace,
		CodeNotAllowedMint,
		fmt.Sprintf("only the creator can mint a token"),
	)
}

func ErrInvalidDenom() *sdkerrors.Error {
	return sdkerrors.New(
		DefaultCodespace,
		CodeInvalidDenom,
		fmt.Sprintf("invalid denom name"),
	)
}

func ErrInvalidTokenID() *sdkerrors.Error {
	return sdkerrors.New(
		DefaultCodespace,
		CodeInvalidTokenID,
		fmt.Sprintf("invalid token name"),
	)
}
