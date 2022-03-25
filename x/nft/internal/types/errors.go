package types

import (
	"fmt"

	"bitbucket.org/decimalteam/go-node/utils/errors"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

type CodeType = uint32

const (
	// Default coin codespace
	DefaultCodespace string = ModuleName

	CodeInvalidCollection             CodeType = 101
	CodeUnknownCollection             CodeType = 102
	CodeInvalidNFT                    CodeType = 103
	CodeUnknownNFT                    CodeType = 104
	CodeNFTAlreadyExists              CodeType = 105
	CodeEmptyMetadata                 CodeType = 106
	CodeInvalidQuantity               CodeType = 107
	CodeInvalidReserve                CodeType = 108
	CodeNotAllowedBurn                CodeType = 109
	CodeNotAllowedMint                CodeType = 110
	CodeInvalidDenom                  CodeType = 111
	CodeInvalidTokenID                CodeType = 112
	CodeNotUniqueSubTokenIDs          CodeType = 113
	CodeNotUniqueTokenURI             CodeType = 114
	CodeOwnerDoesNotOwnSubTokenID     CodeType = 115
	CodeInvalidSenderAddress          CodeType = 116
	CodeInvalidRecipientAddress       CodeType = 117
	CodeForbiddenToTransferToYourself CodeType = 118
	CodeNotUniqueTokenID              CodeType = 119
	CodeNotAllowedUpdateNFTReserve    CodeType = 120
	CodeNotSetValueLowerNow           CodeType = 121
	CodeNotEnoughFunds                CodeType = 122
)

func ErrInvalidCollection(denom string) *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeInvalidCollection,
		fmt.Sprintf("invalid NFT collection: %s", denom),
	)
}

func ErrUnknownCollection(denom string) *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeUnknownCollection,
		fmt.Sprintf("unknown NFT collection: %s", denom),
	)
}

func ErrInvalidNFT(id string) *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeInvalidNFT,
		fmt.Sprintf("invalid NFT: %s", id),
	)
}

func ErrUnknownNFT(denom string, id string) *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeUnknownNFT,
		fmt.Sprintf("unknown NFT: denom = %s, tokenID = %s", denom, id),
	)
}

func ErrNFTAlreadyExists(id string) *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeNFTAlreadyExists,
		fmt.Sprintf("NFT with ID = %s already exists", id),
	)
}

func ErrEmptyMetadata() *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeEmptyMetadata,
		"NFT metadata can't be empty",
	)
}

func ErrInvalidQuantity(quantity string) *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeInvalidQuantity,
		fmt.Sprintf("invalid NFT quantity: %s", quantity),
	)
}

func ErrInvalidReserve(reserve string) *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeInvalidReserve,
		fmt.Sprintf("invalid NFT reserve: %s", reserve),
	)
}

func ErrNotAllowedBurn() *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeNotAllowedBurn,
		"only the creator can burn a token",
	)
}

func ErrNotSetValueLowerNow() *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeNotSetValueLowerNow,
		"Invalid new reserve",
	)
}

func ErrNotAllowedUpdateReserve() *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeNotAllowedUpdateNFTReserve,
		"only the creator can update reserve a token",
	)
}

func ErrNotAllowedMint() *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeNotAllowedMint,
		"only the creator can mint a token",
	)
}

func ErrInvalidDenom(denom string) *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeInvalidDenom,
		fmt.Sprintf("invalid denom name: %s", denom),
	)
}

func ErrInvalidTokenID(name string) *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeInvalidTokenID,
		fmt.Sprintf("invalid token name: %s", name),
	)
}

func ErrNotUniqueSubTokenIDs() *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeNotUniqueSubTokenIDs,
		"not unique subTokenIDs",
	)
}

func ErrNotUniqueTokenURI() *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeNotUniqueTokenURI,
		"not unique tokenURI",
	)
}

func ErrNotUniqueTokenID() *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeNotUniqueTokenID,
		"not unique token id",
	)
}

func ErrOwnerDoesNotOwnSubTokenID(owner string, subTokenID string) *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeOwnerDoesNotOwnSubTokenID,
		fmt.Sprintf("owner %s does not own sub tokenID %s", owner, subTokenID),
	)
}

func ErrInvalidSenderAddress(address string) *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeInvalidSenderAddress,
		fmt.Sprintf("invalid sender address: %s", address),
	)
}

func ErrInvalidRecipientAddress(address string) *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeInvalidRecipientAddress,
		fmt.Sprintf("invalid recipient address: %s", address),
	)
}

func ErrForbiddenToTransferToYourself() *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeForbiddenToTransferToYourself,
		"Forbidden to transfer to yourself",
	)
}

func ErrNotEnoughFunds(reserve string) *sdkerrors.Error {
	return errors.Encode(
		DefaultCodespace,
		CodeNotEnoughFunds,
		fmt.Sprintf("Insufficient funds are required: %s", reserve),
	)
}
