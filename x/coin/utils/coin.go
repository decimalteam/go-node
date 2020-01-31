package utils

import (
	"bitbucket.org/decimalteam/go-node/x/coin/internal/types"
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
)

// Check if coin exists
func ExistsCoin(cliCtx client.CLIContext, symbol string) (bool, error) {
	res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%s", types.ModuleName, types.QueryGetCoin, symbol), nil)
	if err != nil {
		return res != nil, nil
	} else {
		return false, err
	}
}

// Return coin instance from State
func GetCoin(cliCtx client.CLIContext, symbol string) (types.Coin, error) {
	res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%s", types.ModuleName, types.QueryGetCoin, symbol), nil)
	coin := types.Coin{}
	if err = cliCtx.Codec.UnmarshalJSON(res, &coin); err != nil {
		return coin, err
	}
	return coin, err
}
