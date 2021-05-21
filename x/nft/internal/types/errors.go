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

func ErrInvalidCollection(denom string) *sdkerrors.Error {
	return sdkerrors.New(
		DefaultCodespace,
		CodeInvalidCollection,
		fmt.Sprintf("invalid NFT collection: %s", denom),
	)
}

func ErrUnknownCollection(denom string) *sdkerrors.Error {
	return sdkerrors.New(
		DefaultCodespace,
		CodeUnknownCollection,
		fmt.Sprintf("unknown NFT collection: %s", denom),
	)
}

func ErrInvalidNFT() *sdkerrors.Error {
	return sdkerrors.New(
		DefaultCodespace,
		CodeInvalidNFT,
		fmt.Sprintf("invalid NFT"),
	)
}

func ErrUnknownNFT(denom string) *sdkerrors.Error {
	return sdkerrors.New(
		DefaultCodespace,
		CodeUnknownNFT,
		fmt.Sprintf("unknown NFT: %s", denom),
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

func ErrInvalidQuantity(quantity string) *sdkerrors.Error {
	return sdkerrors.New(
		DefaultCodespace,
		CodeInvalidQuantity,
		fmt.Sprintf("invalid NFT quantity: %s", quantity),
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
