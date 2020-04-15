package types

// Local code type
type CodeType = uint32

const (
	// Default multisig codespace
	DefaultCodespace string = ModuleName

	// CodeInvalid      CodeType = 101
)

// TODO: Fill out some custom errors for the module
// You can see how they are constructed below:
// var (
//	ErrInvalid = sdkerrors.Register(ModuleName, 1, "custom error message")
// )
