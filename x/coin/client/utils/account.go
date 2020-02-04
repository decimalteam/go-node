package utils

import (
	ctx "github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/exported"
)

func GetAccount(cliCtx ctx.CLIContext, addr sdk.AccAddress) (exported.Account, error) {
	ar := auth.NewAccountRetriever(cliCtx)
	account, _, err := ar.GetAccountWithHeight(addr)
	return account, err
}
