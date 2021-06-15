package types

import (
	"bitbucket.org/decimalteam/go-node/utils/errors"
	"fmt"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"strconv"
)

type CodeType = uint32

const (
	// Default coin codespace
	DefaultCodespace string = ModuleName

	CodeInvalidCollection         CodeType = 101
	CodeUnknownCollection         CodeType = 102
	CodeInvalidNFT                CodeType = 103
	CodeUnknownNFT                CodeType = 104
	CodeNFTAlreadyExists          CodeType = 105
	CodeEmptyMetadata             CodeType = 106
	CodeInvalidQuantity           CodeType = 107
	CodeInvalidReserve            CodeType = 108
	CodeNotAllowedBurn            CodeType = 109
	CodeNotAllowedMint            CodeType = 110
	CodeInvalidDenom              CodeType = 111
	CodeInvalidTokenID            CodeType = 112
	CodeNotUniqueSubTokenIDs      CodeType = 113
	CodeNotUniqueTokenURI         CodeType = 114
	CodeOwnerDoesNotOwnSubTokenID CodeType = 115
)

func ErrInvalidCollection(denom string) *sdkerrors.Error {
	return sdkerrors.New(
		DefaultCodespace,
		CodeInvalidCollection,
		errors.Encode(
			DefaultCodespace,
			fmt.Sprintf("invalid NFT collection: %s", denom),
			errors.NewParam("denom", denom),
		),
	)
}

func ErrUnknownCollection(denom string) *sdkerrors.Error {
	return sdkerrors.New(
		DefaultCodespace,
		CodeUnknownCollection,
		errors.Encode(
			DefaultCodespace,
			fmt.Sprintf("unknown NFT collection: %s", denom),
			errors.NewParam("denom", denom),
		),
	)
}

func ErrInvalidNFT(id string) *sdkerrors.Error {
	return sdkerrors.New(
		DefaultCodespace,
		CodeInvalidNFT,
		errors.Encode(
			DefaultCodespace,
			fmt.Sprintf("invalid NFT: %s", id),
			errors.NewParam("id", id),
		),
	)
}

func ErrUnknownNFT(denom, id string) *sdkerrors.Error {
	return sdkerrors.New(
		DefaultCodespace,
		CodeUnknownNFT,
		errors.Encode(
			DefaultCodespace,
			"", /*fmt.Sprintf("unknown NFT: denom = %s, tokenID = %s", denom, id)*/
			errors.NewParam("id", id),
			errors.NewParam("denom", denom),
		),
	)
}

func ErrNFTAlreadyExists(id string) *sdkerrors.Error {
	return sdkerrors.New(
		DefaultCodespace,
		CodeNFTAlreadyExists,
		errors.Encode(
			DefaultCodespace,
			fmt.Sprintf("NFT with ID = %s already exists", id),
			errors.NewParam("id", id),
		),
	)
}

func ErrEmptyMetadata() *sdkerrors.Error {
	return sdkerrors.New(
		DefaultCodespace,
		CodeEmptyMetadata,
		errors.Encode(
			DefaultCodespace,
			"NFT metadata can't be empty",
		),
	)
}

func ErrInvalidQuantity(quantity string) *sdkerrors.Error {
	return sdkerrors.New(
		DefaultCodespace,
		CodeInvalidQuantity,
		errors.Encode(
			DefaultCodespace,
			fmt.Sprintf("invalid NFT quantity: %s", quantity),
			errors.NewParam("quantity", quantity),
		),
	)
}

func ErrInvalidReserve(reserve string) *sdkerrors.Error {
	return sdkerrors.New(
		DefaultCodespace,
		CodeInvalidReserve,
		errors.Encode(
			DefaultCodespace,
			fmt.Sprintf("invalid NFT reserve: %s", reserve),
			errors.NewParam("reserve", reserve),
		),
	)
}

func ErrNotAllowedBurn() *sdkerrors.Error {
	return sdkerrors.New(
		DefaultCodespace,
		CodeNotAllowedBurn,
		errors.Encode(
			DefaultCodespace,
			"only the creator can burn a token",
		),
	)
}

func ErrNotAllowedMint() *sdkerrors.Error {
	return sdkerrors.New(
		DefaultCodespace,
		CodeNotAllowedMint,
		errors.Encode(
			DefaultCodespace,
			"only the creator can mint a token",
		),
	)
}

func ErrInvalidDenom(denom string) *sdkerrors.Error {
	return sdkerrors.New(
		DefaultCodespace,
		CodeInvalidDenom,
		errors.Encode(
			DefaultCodespace,
			fmt.Sprintf("invalid denom name: %s", denom),
			errors.NewParam("denom", denom),
		),
	)
}

func ErrInvalidTokenID(name string) *sdkerrors.Error {
	return sdkerrors.New(
		DefaultCodespace,
		CodeInvalidTokenID,
		errors.Encode(
			DefaultCodespace,
			fmt.Sprintf("invalid token name: %s", name),
			errors.NewParam("name", name),
		),
	)
}

func ErrNotUniqueSubTokenIDs() *sdkerrors.Error {
	return sdkerrors.New(
		DefaultCodespace,
		CodeNotUniqueSubTokenIDs,
		errors.Encode(
			DefaultCodespace,
			"not unique subTokenIDs",
		),
	)
}

func ErrNotUniqueTokenURI() *sdkerrors.Error {
	return sdkerrors.New(
		DefaultCodespace,
		CodeNotUniqueTokenURI,
		errors.Encode(
			DefaultCodespace,
			"not unique tokenURI",
		),
	)
}

func ErrOwnerDoesNotOwnSubTokenID(owner string, subTokenID int64) *sdkerrors.Error {
	return sdkerrors.New(
		DefaultCodespace,
		CodeOwnerDoesNotOwnSubTokenID,
		errors.Encode(
			DefaultCodespace,
			fmt.Sprintf("owner %s does not own sub tokenID %d", owner, subTokenID),
			errors.NewParam("owner", owner),
			errors.NewParam("sub_token_id", strconv.FormatInt(subTokenID, 10)),
		),
	)
}
