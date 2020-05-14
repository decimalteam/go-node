package utils

import (
	ctx "github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"bitbucket.org/decimalteam/go-node/x/auth"
	"bitbucket.org/decimalteam/go-node/x/auth/exported"
)

// GetAccount returns account for given address if exists.
func GetAccount(cliCtx ctx.CLIContext, addr sdk.AccAddress) (exported.Account, error) {
	ar := auth.NewAccountRetriever(cliCtx)
	account, _, err := ar.GetAccountWithHeight(addr)
	return account, err
}
