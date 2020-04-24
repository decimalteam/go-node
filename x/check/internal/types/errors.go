package types

// import (
// 	sdk "github.com/cosmos/cosmos-sdk/types"
// )

// Local code type
type CodeType = uint32

const (
	// Default check codespace
	DefaultCodespace string = ModuleName

	InvalidVRS       CodeType = 101
	InvalidPublicKey CodeType = 102
	// CodeInvalid      CodeType = 101
)

// TODO: Fill out some custom errors for the module
// You can see how they are constructed below:
// var (
//	ErrInvalid = sdkerrors.Register(ModuleName, 1, "custom error message")
// )
