package utils

import "bitbucket.org/decimalteam/go-node/x/validator/internal/types"

var (
	NewTxBuilderFromCLI = types.NewTxBuilderFromCLI
	NewTxBuilder        = types.NewTxBuilder
)

type (
	TxBuilder = types.TxBuilder
	StdTx     = types.StdTx
)
