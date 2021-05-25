package types

import (
	"encoding/json"
	"fmt"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"strconv"
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
	jsonData, _ := json.Marshal(
		map[string]string{
			"code":      getCodeString(CodeInvalidCollection),
			"codespace": DefaultCodespace,
			"desc":      fmt.Sprintf("invalid NFT collection: %s", denom),
			"denom":     denom,
		},
	)
	return sdkerrors.New(
		DefaultCodespace,
		CodeInvalidCollection,
		string(jsonData),
	)
}

func ErrUnknownCollection(denom string) *sdkerrors.Error {
	jsonData, _ := json.Marshal(
		map[string]string{
			"code":      getCodeString(CodeUnknownCollection),
			"codespace": DefaultCodespace,
			"desc":      fmt.Sprintf("unknown NFT collection: %s", denom),
			"denom":     denom,
		},
	)
	return sdkerrors.New(
		DefaultCodespace,
		CodeUnknownCollection,
		string(jsonData),
	)
}

func ErrInvalidNFT() *sdkerrors.Error {
	jsonData, _ := json.Marshal(
		map[string]string{
			"code":      getCodeString(CodeInvalidNFT),
			"codespace": DefaultCodespace,
			"desc":      "invalid NFT",
		},
	)
	return sdkerrors.New(
		DefaultCodespace,
		CodeInvalidNFT,
		string(jsonData),
	)
}

func ErrUnknownNFT(denom string) *sdkerrors.Error {
	jsonData, _ := json.Marshal(
		map[string]string{
			"code":      getCodeString(CodeUnknownNFT),
			"codespace": DefaultCodespace,
			"desc":      fmt.Sprintf("unknown NFT: %s", denom),
			"denom":     denom,
		},
	)
	return sdkerrors.New(
		DefaultCodespace,
		CodeUnknownNFT,
		string(jsonData),
	)
}

func ErrNFTAlreadyExists() *sdkerrors.Error {
	jsonData, _ := json.Marshal(
		map[string]string{
			"code":      getCodeString(CodeNFTAlreadyExists),
			"codespace": DefaultCodespace,
			"desc":      fmt.Sprintf("NFT already exists"),
		},
	)
	return sdkerrors.New(
		DefaultCodespace,
		CodeNFTAlreadyExists,
		string(jsonData),
	)
}

func ErrEmptyMetadata() *sdkerrors.Error {
	jsonData, _ := json.Marshal(
		map[string]string{
			"code":      getCodeString(CodeUnknownNFT),
			"codespace": DefaultCodespace,
			"desc":      fmt.Sprintf("NFT metadata can't be empty"),
			"denom":     denom,
		},
	)
	return sdkerrors.New(
		DefaultCodespace,
		CodeEmptyMetadata,
		string(jsonData),
	)
}

func ErrInvalidQuantity(quantity string) *sdkerrors.Error {
	jsonData, _ := json.Marshal(
		map[string]string{
			"code":      getCodeString(CodeInvalidQuantity),
			"codespace": DefaultCodespace,
			"desc":      fmt.Sprintf("invalid NFT quantity: %s", quantity),
			"quantity":  quantity,
		},
	)
	return sdkerrors.New(
		DefaultCodespace,
		CodeInvalidQuantity,
		string(jsonData),
	)
}

func ErrInvalidReserve(reserve string) *sdkerrors.Error {
	jsonData, _ := json.Marshal(
		map[string]string{
			"code":      getCodeString(CodeInvalidReserve),
			"codespace": DefaultCodespace,
			"desc":      fmt.Sprintf("invalid NFT reserve: %s", reserve),
			"reserve":   reserve,
		},
	)
	return sdkerrors.New(
		DefaultCodespace,
		CodeInvalidReserve,
		string(jsonData),
	)
}

func ErrNotAllowedBurn() *sdkerrors.Error {
	jsonData, _ := json.Marshal(
		map[string]string{
			"code":      getCodeString(CodeNotAllowedBurn),
			"codespace": DefaultCodespace,
			"desc":      fmt.Sprintf("only the creator can burn a token"),
		},
	)
	return sdkerrors.New(
		DefaultCodespace,
		CodeNotAllowedBurn,
		string(jsonData),
	)
}

func ErrNotAllowedMint() *sdkerrors.Error {
	jsonData, _ := json.Marshal(
		map[string]string{
			"code":      getCodeString(CodeNotAllowedMint),
			"codespace": DefaultCodespace,
			"desc":      fmt.Sprintf("only the creator can mint a token"),
		},
	)
	return sdkerrors.New(
		DefaultCodespace,
		CodeNotAllowedMint,
		string(jsonData),
	)
}

func ErrInvalidDenom(denom string) *sdkerrors.Error {
	jsonData, _ := json.Marshal(
		map[string]string{
			"code":      getCodeString(CodeInvalidDenom),
			"codespace": DefaultCodespace,
			"desc":      fmt.Sprintf("invalid denom name: %s", denom),
			"denom":     denom,
		},
	)
	return sdkerrors.New(
		DefaultCodespace,
		CodeInvalidDenom,
		string(jsonData),
	)
}

func ErrInvalidTokenID(name string) *sdkerrors.Error {
	jsonData, _ := json.Marshal(
		map[string]string{
			"code":      getCodeString(CodeInvalidTokenID),
			"codespace": DefaultCodespace,
			"desc":      fmt.Sprintf("invalid token name: %s", name),
			"name":      name,
		},
	)
	return sdkerrors.New(
		DefaultCodespace,
		CodeInvalidTokenID,
		string(jsonData),
	)
}

func getCodeString(code CodeType) string {
	return strconv.FormatInt(int64(code), 10)
}
