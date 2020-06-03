package types

// CodeType defines the local code type.
type CodeType = uint32

// DefaultCodespace defines default multisig codespace.
const DefaultCodespace string = ModuleName

// Custom errors codes.
const (
	InvalidSender       CodeType = 101
	InvalidOwnerCount   CodeType = 102
	InvalidOwner        CodeType = 103
	InvalidWeightCount  CodeType = 104
	InvalidWeight       CodeType = 105
	InvalidCoinToSend   CodeType = 106
	InvalidAmountToSend CodeType = 107
)
