package utils

import (
	ctx "github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authTypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

// Returns account for given address if exists
func GetAccount(ctx ctx.Context, addr sdk.AccAddress) (ctx.Account, error) {
	ar := authTypes.AccountRetriever{}
	account, _, err := ar.GetAccountWithHeight(ctx, addr)
	return account, err
}
